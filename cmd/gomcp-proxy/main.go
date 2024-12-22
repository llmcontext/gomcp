package main

import (
	"fmt"
	"net"
	"os"

	"github.com/google/uuid"
	"github.com/llmcontext/gomcp/channels/proxy"
	"github.com/llmcontext/gomcp/defaults"
	"github.com/llmcontext/gomcp/logger"
	"github.com/llmcontext/gomcp/types"
	"github.com/llmcontext/gomcp/version"
	"github.com/spf13/cobra"
)

var (
	muxAddress string
	rootCmd    = &cobra.Command{
		Use:   "gomcp-proxy",
		Short: "A proxy server for MCP connections",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewTermLogger()

			// banner
			logger.Header(fmt.Sprintf("%s - %s", proxy.GomcpProxyClientName, version.Version))

			if len(args) == 0 {
				logger.Error("Please provide a program name as the first argument", types.LogArg{"args": args})
				os.Exit(1)
			}

			programName := args[0]
			args = args[1:]

			logger.Info("MCP Proxy is starting", types.LogArg{
				"address":     muxAddress,
				"programName": programName,
				"args":        args,
			})

			// check if address is valid
			if _, err := net.ResolveTCPAddr("tcp", muxAddress); err != nil {
				logger.Error("Invalid address for MCP Proxy", types.LogArg{"address": muxAddress, "error": err})
				os.Exit(1)
			}

			currentWorkingDirectory, err := os.Getwd()
			if err != nil {
				logger.Error("Failed to get current working directory", types.LogArg{"error": err})
				os.Exit(1)
			}

			// read the config file
			configPath := currentWorkingDirectory + "/" + defaults.DefaultProxyConfigPath
			proxyConfig, err := LoadProxyConfig(configPath)
			if err != nil {
				logger.Error("Failed to load proxy configuration file",
					types.LogArg{"error": err, "configPath": configPath})
				os.Exit(1)
			}

			if proxyConfig == nil {
				logger.Info("creating proxy configuration file", types.LogArg{"configPath": configPath})
				proxyConfig = &ProxyConfig{
					MuxAddress:  muxAddress,
					ProgramName: programName,
					ProgramArgs: args,
				}
			}

			// update the proxy config with the current values
			proxyConfig.MuxAddress = muxAddress
			proxyConfig.ProgramName = programName
			proxyConfig.ProgramArgs = args
			proxyConfig.WhatIsThat = defaults.DefaultProxyWhatIsThat
			proxyConfig.MoreInformation = defaults.DefaultProxyMoreInfo

			if proxyConfig.ProxyId == "" {
				proxyConfig.ProxyId = uuid.New().String()
				logger.Info("generated new proxy id", types.LogArg{"proxyId": proxyConfig.ProxyId})
			}

			// save the proxy config to the file
			err = SaveProxyConfig(configPath, proxyConfig)
			if err != nil {
				logger.Error("Failed to save proxy configuration file",
					types.LogArg{"error": err, "configPath": configPath})
				os.Exit(1)
			}

			proxyInformation := proxy.ProxyInformation{
				MuxAddress:              muxAddress,
				CurrentWorkingDirectory: currentWorkingDirectory,
				ProgramName:             programName,
				Args:                    args,
			}

			client := proxy.NewProxyClient(proxyInformation, logger)
			client.Start()
		},
	}
)

func init() {
	rootCmd.Flags().StringVarP(&muxAddress, "mux", "x", fmt.Sprintf(":%d", defaults.DefaultMultiplexerPort), "TCP address for the MCP multiplexer server (host:port)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
