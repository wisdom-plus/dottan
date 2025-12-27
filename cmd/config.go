package cmd

import "github.com/spf13/cobra"

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage config file",
}

func init() {
	rootCmd.AddCommand(configCmd)
}
