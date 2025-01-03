package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	shareCmd.Flags().StringP("name", "n", "", "the name of the share")
	shareCmd.Flags().StringP("path", "p", "", "the path to the directory to share (default: current directory)")
	rootCmd.AddCommand(shareCmd)
}

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "share directory with gomcp",
	Long:  "add a directory to the gomcp server to share with the client",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		// retrieve the name and path from the flags
		name, _ := cmd.Flags().GetString("name")
		path, _ := cmd.Flags().GetString("path")
		if name == "" {
			fmt.Println("name is required")
			return
		}
		if path == "" {
			path, err = os.Getwd()
			if err != nil {
				fmt.Println("error getting current directory:", err)
				os.Exit(1)
			}
		}
		path, err = convertRelativePathToAbsolute(path)
		if err != nil {
			fmt.Println("error converting relative path to absolute:", err)
			os.Exit(1)
		}

		// check if the path exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Println("path does not exist:", path)
			os.Exit(1)
		}
		fmt.Println("name:", name)
		fmt.Println("path:", path)
	},
}

func convertRelativePathToAbsolute(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	return filepath.Abs(path)
}
