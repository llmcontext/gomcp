package main

import (
	"fmt"
	"os"

	"github.com/llmcontext/gomcp/config"
	"github.com/llmcontext/gomcp/modelcontextprotocol/mcpserver"
	"github.com/spf13/cobra"
)

var (
	debug   bool
	rootCmd = &cobra.Command{
		Use:   "gomcp",
		Short: "A Model Context Protocol server multiplexer",
		Run: func(cmd *cobra.Command, args []string) {
			// we read the configuration file
			conf, err := config.LoadHubConfiguration()
			if err != nil {
				fmt.Println("Error reading configuration file:", err)
				os.Exit(1)
			}

			// we create the MCP server
			mcpServer, err := mcpserver.NewMcpServer(conf.ServerInfo, conf.Logging, debug)
			if err != nil {
				fmt.Println("Error creating MCP server:", err)
				os.Exit(1)
			}

			// start the server
			transport := mcpServer.StdioTransport()
			err = mcpServer.Start(transport)
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
