package main

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/llmcontext/gomcp/channels/proxy"
	"github.com/llmcontext/gomcp/config"
	"github.com/llmcontext/gomcp/defaults"
	"github.com/llmcontext/gomcp/logger"
	"github.com/llmcontext/gomcp/types"
	"github.com/llmcontext/gomcp/version"
	"github.com/spf13/cobra"
)

var (
	debug   bool
	rootCmd = &cobra.Command{
		Use:   "gomcp-proxy",
		Short: "A proxy server for MCP connections",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewTermLogger(debug)

			// banner
			logger.Header(fmt.Sprintf("%s - %s", proxy.GomcpProxyClientName, version.Version))

			// get the current working directory
			currentWorkingDirectory, err := os.Getwd()
			if err != nil {
				logger.Error("Failed to get current working directory", types.LogArg{"error": err})
				os.Exit(1)
			}

			// read the config file
			proxyConfig, err := config.LoadProxyConfiguration(currentWorkingDirectory)
			if err != nil {
				logger.Error("Failed to load proxy configuration file",
					types.LogArg{"error": err, "configPath": proxyConfig.ConfigurationFilePath})
				os.Exit(1)
			}

			var invalidArgs bool
			var programName string
			var programArgs []string

			if len(args) == 0 {
				// that's ok if we have a config file
				if proxyConfig == nil {
					invalidArgs = true
				} else {
					programName = proxyConfig.ProgramName
					programArgs = proxyConfig.ProgramArgs
					invalidArgs = false
				}
			} else {
				programName = args[0]
				programArgs = args[1:]
				invalidArgs = false
			}

			if invalidArgs {
				logger.Error("Please provide a program name as the first argument, and optionally arguments", types.LogArg{"args": args})
				os.Exit(1)
			}

			// load the hub configuration
			hubConfig, err := config.LoadHubConfiguration()
			if err != nil {
				logger.Error("Failed to load hub configuration file",
					types.LogArg{
						"error":      err,
						"configPath": config.GetDefaultHubConfigurationPath(),
					})
				os.Exit(1)
			}

			// check if we have a proxy configuration
			if hubConfig.Proxy == nil || !hubConfig.Proxy.Enabled {
				logger.Info("proxy is not enabled in the hub configuration",
					types.LogArg{"configPath": config.GetDefaultHubConfigurationPath()})
				os.Exit(0)
			}

			logger.Info("MCP Proxy is starting", types.LogArg{
				"address":     hubConfig.Proxy.ListenAddress,
				"programName": programName,
				"programArgs": programArgs,
			})

			if proxyConfig == nil {
				logger.Info("creating proxy configuration file", types.LogArg{"configPath": proxyConfig.ConfigurationFilePath})
				proxyConfig = &config.ProxyConfiguration{}
			}

			// update the proxy config with the current values
			proxyConfig.ProgramName = programName
			proxyConfig.ProgramArgs = programArgs
			proxyConfig.WhatIsThat = defaults.DefaultProxyWhatIsThat
			proxyConfig.MoreInformation = defaults.DefaultProxyMoreInfo

			if proxyConfig.ProxyId == "" {
				proxyConfig.ProxyId = uuid.New().String()
				logger.Info("generated new proxy id", types.LogArg{"proxyId": proxyConfig.ProxyId})
			}

			// save the proxy config to the file
			err = config.SaveProxyConfiguration(proxyConfig)
			if err != nil {
				logger.Error("Failed to save proxy configuration file",
					types.LogArg{"error": err, "configPath": proxyConfig.ConfigurationFilePath})
				os.Exit(1)
			}

			proxyInformation := proxy.ProxyInformation{
				ProxyId:                 proxyConfig.ProxyId,
				MuxAddress:              hubConfig.Proxy.ListenAddress,
				CurrentWorkingDirectory: currentWorkingDirectory,
				ProgramName:             programName,
				Args:                    programArgs,
			}

			client := proxy.NewProxyClient(proxyInformation, debug, logger)
			client.Start()
		},
	}
)

func init() {
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
