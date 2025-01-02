package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/llmcontext/gomcp/providers/proxies"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var proxyDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a proxy",
	Long:  `Delete a proxy from the gomcp service.`,
	Run: func(cmd *cobra.Command, args []string) {
		proxyRegistry, err := proxies.NewProxyRegistry()
		if err != nil {
			fmt.Println("Error getting proxy registry:", err)
			os.Exit(1)
		}
		// Initialize an empty slice to hold the options
		var options []string

		proxies := proxyRegistry.GetProxies()
		for _, proxy := range proxies {
			lstTools := []string{}
			for _, tool := range proxy.Tools {
				lstTools = append(lstTools, tool.Name)
			}
			cli := fmt.Sprintf("%s: %s", proxy.ProxyId, strings.Join(lstTools, ","))
			options = append(options, cli)
		}

		// Use PTerm's interactive select feature to present the options to the user and capture their selection
		// The Show() method displays the options and waits for the user's input
		selectedOption, _ := pterm.DefaultInteractiveSelect.
			WithDefaultText("Select a proxy to delete").
			WithOptions(options).
			Show()

		// extract proxyId from selectedOption
		proxyId := strings.Split(selectedOption, ":")[0]

		pterm.Println()
		newHeader := pterm.HeaderPrinter{
			TextStyle:       pterm.NewStyle(pterm.FgBlack),
			BackgroundStyle: pterm.NewStyle(pterm.BgRed),
			Margin:          20,
		}
		newHeader.Println(selectedOption)
		pterm.Println()

		// ask for confirmation
		confirm, _ := pterm.DefaultInteractiveConfirm.WithDefaultText("Are you sure you want to delete this proxy?").Show()
		if !confirm {
			fmt.Println("Aborted")
			os.Exit(1)
		}

		pterm.Info.Printfln("Deleting proxy %s", pterm.Green(proxyId))

	},
}

func init() {
	proxyCmd.AddCommand(proxyDeleteCmd)
}
