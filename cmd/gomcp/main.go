package main

import (
	"fmt"
	"os"

	"github.com/llmcontext/gomcp"
	"github.com/spf13/cobra"
)

var (
	configFile string
	rootCmd    = &cobra.Command{
		Use:   "gomcp",
		Short: "A MCP multiplexer server that enables multiple MCP proxy client connections",
		Run: func(cmd *cobra.Command, args []string) {
			// we create the MCP server
			mcp, err := gomcp.NewModelContextProtocolServer(configFile)
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
	rootCmd.Flags().StringVarP(&configFile, "configFile", "f", "", "Path to configuration file")
	rootCmd.MarkFlagRequired("configFile")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
