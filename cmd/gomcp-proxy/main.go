package main

import (
	"fmt"
	"os"

	"github.com/llmcontext/gomcp/proxy"
	"github.com/llmcontext/gomcp/version"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	port    int
	rootCmd = &cobra.Command{
		Use:   "gomcp-proxy",
		Short: "A proxy server for MCP connections",
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Println(fmt.Sprintf("%s, version %s", version.Version, proxy.GomcpProxyClientName))
			if len(args) == 0 {
				pterm.Error.Println("Please provide a program name as the first argument")
				os.Exit(1)
			}

			programName := args[0]
			args = args[1:]

			// Print an informational message using PTerm's Info printer.
			// This message will stay in place while the area updates.
			pterm.Info.Println("MCP Proxy is starting")
			pterm.Info.Println(fmt.Sprintf("- ws port is: %d\n", port))
			pterm.Info.Println(fmt.Sprintf("- program name is: %s\n", programName))
			pterm.Info.Println(fmt.Sprintf("- program args are: %v\n", args))

			client := proxy.NewClient(programName, args)
			client.Start()
		},
	}
)

func init() {
	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port number for the WebSocket server")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
