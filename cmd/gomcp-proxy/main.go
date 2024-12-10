package main

import (
	"fmt"
	"os"

	"github.com/llmcontext/gomcp/proxy"
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
			for ix, arg := range args {
				fmt.Printf("Argument %d: %s\n", ix, arg)
			}

			if len(args) == 0 {
				fmt.Println("Please provide a program name as the first argument")
				os.Exit(1)
			}

			programName := args[0]
			args = args[1:]

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
