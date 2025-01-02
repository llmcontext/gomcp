package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var proxyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured proxies",
	Long:  `Display a list of all proxies that have been configured in the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement actual proxy listing logic
		fmt.Println("Listing all configured proxies...")
	},
}

func init() {
	proxyCmd.AddCommand(proxyListCmd)
}
