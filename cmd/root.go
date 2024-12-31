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
		Short: "A Model Context Protocol server multiplexer",
		Run: func(cmd *cobra.Command, args []string) {
			// we create the MCP server
			mcp, err := hub.NewHubModelContextProtocolServer(debug)
			if err != nil {
				fmt.Println("Error creating MCP server:", err)
				os.Exit(1)
			}

			// start the server
			transport := mcp.StdioTransport()
			err = mcp.Start(transport)
			if err != nil {
				fmt.Println("Error starting MCP server:", err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
