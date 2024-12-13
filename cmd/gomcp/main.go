package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	port    int
	rootCmd = &cobra.Command{
		Use:   "gomcp",
		Short: "A MCP multiplexer server that enables multiple MCP proxy client connections",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("gomcp")
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
