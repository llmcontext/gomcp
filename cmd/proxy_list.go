package main

import (
	"fmt"
	"os"

	"github.com/llmcontext/gomcp/providers/proxies"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/cobra"
)

var proxyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured proxies",
	Long:  `Display a list of all proxies that have been configured in the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		proxyRegistry, err := proxies.NewProxyRegistry()
		if err != nil {
			fmt.Println("Error getting proxy registry:", err)
			os.Exit(1)
		}
		leveledList := pterm.LeveledList{}
		proxies := proxyRegistry.GetProxies()
		for _, proxy := range proxies {
			leveledList = append(leveledList, pterm.LeveledListItem{
				Level: 0,
				Text:  proxy.ProxyId,
			})
			for _, tool := range proxy.Tools {
				leveledList = append(leveledList, pterm.LeveledListItem{
					Level: 1,
					Text:  fmt.Sprintf("%s - %s", tool.Name, tool.Description),
				})
			}
		}
		// Convert the leveled list into a tree structure.
		root := putils.TreeFromLeveledList(leveledList)
		root.Text = "gomcp proxies" // Set the root node text.

		// Render the tree structure using the default tree printer.
		pterm.DefaultTree.WithRoot(root).Render()
	},
}

func init() {
	proxyCmd.AddCommand(proxyListCmd)
}
