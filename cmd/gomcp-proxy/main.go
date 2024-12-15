package main

import (
	"fmt"
	"net"
	"os"

	"github.com/llmcontext/gomcp/defaults"
	"github.com/llmcontext/gomcp/logger"
	"github.com/llmcontext/gomcp/proxy"
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
	rootCmd.Flags().StringVarP(&muxAddress, "address", "a", fmt.Sprintf(":%d", defaults.DefaultMultiplexerPort), "TCP address for the MCP multiplexer server (host:port)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
