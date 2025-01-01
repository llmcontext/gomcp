package main

import (
	"fmt"

	"github.com/llmcontext/gomcp/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of gomcp",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gomcp version:", version.Version)
	},
}
