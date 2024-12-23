package main

import (
	"fmt"
	"os"

	"github.com/llmcontext/gomcp/channels/hub"
	"github.com/spf13/cobra"
)

var (
	debug   bool
	rootCmd = &cobra.Command{
		Use:   "gomcp",
		Short: "A MCP multiplexer server that enables multiple MCP proxy client connections",
		Run: func(cmd *cobra.Command, args []string) {
			// we create the MCP server
			mcp, err := hub.NewHubModelContextProtocolServer(debug)
			if err != nil {
				fmt.Println("Error creating MCP server:", err)
				os.Exit(1)
			}

			// start the server
			transport := mcp.StdioTransport()
			mcp.Start(transport)
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
