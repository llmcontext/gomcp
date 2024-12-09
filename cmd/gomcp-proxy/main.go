package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	port    int
	rootCmd = &cobra.Command{
		Use:   "gomcp-proxy",
		Short: "A proxy server for MCP connections",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Starting gomcp-proxy on port %d\n", port)

			// display the arguments
			fmt.Printf("Arguments: %v\n", args)
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
