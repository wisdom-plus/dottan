/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
)

// versionCmd represents the version command
func newVersionCmd(out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the version of Dottan CLI",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(out, version)
			return err
		},
	}
}

func init() {
	stdout := rootCmd.OutOrStdout()
	rootCmd.AddCommand(newVersionCmd(stdout))
}
