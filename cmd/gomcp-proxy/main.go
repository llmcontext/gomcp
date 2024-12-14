package main

import (
	"fmt"
	"net"
	"os"

	"github.com/llmcontext/gomcp/defaults"
	"github.com/llmcontext/gomcp/proxy"
	"github.com/llmcontext/gomcp/version"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	muxAddress string
	rootCmd    = &cobra.Command{
		Use:   "gomcp-proxy",
		Short: "A proxy server for MCP connections",
		Run: func(cmd *cobra.Command, args []string) {
			message := fmt.Sprintf("%s - %s", proxy.GomcpProxyClientName, version.Version)
			pterm.DefaultHeader.WithFullWidth().Println(message)
			pterm.Println()

			if len(args) == 0 {
				pterm.Error.Println("Please provide a program name as the first argument")
				os.Exit(1)
			}

			programName := args[0]
			args = args[1:]

			// check if address is valid
			if _, err := net.ResolveTCPAddr("tcp", muxAddress); err != nil {
				fmt.Printf("Invalid address for MCP Proxy: %s, err: %s\n", muxAddress, err)
				os.Exit(1)
			}

			// Print an informational message using PTerm's Info printer.
			// This message will stay in place while the area updates.
			pterm.Info.Println("MCP Proxy is starting")
			pterm.Info.Println(fmt.Sprintf("- ws address is: %s\n", muxAddress))
			pterm.Info.Println(fmt.Sprintf("- program name is: %s\n", programName))
			pterm.Info.Println(fmt.Sprintf("- program args are: %v\n", args))

			currentWorkingDirectory, err := os.Getwd()
			if err != nil {
				pterm.Error.Println(fmt.Sprintf("Failed to get current working directory: %s", err))
				os.Exit(1)
			}

			client := proxy.NewProxyClient(muxAddress, currentWorkingDirectory, programName, args)
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
