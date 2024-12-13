package main

import (
	"fmt"
	"os"

	"github.com/llmcontext/gomcp/config"
	"github.com/spf13/cobra"
)

var (
	configFile string
	rootCmd    = &cobra.Command{
		Use:   "gomcp",
		Short: "A MCP multiplexer server that enables multiple MCP proxy client connections",
		Run: func(cmd *cobra.Command, args []string) {
			// we load the config file
			config, err := config.LoadConfig(configFile)
			if err != nil {
				fmt.Printf("failed to load config file %s: %v", configFile, err)
				os.Exit(1)
			}

			fmt.Printf("config: %+v\n", config)

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
