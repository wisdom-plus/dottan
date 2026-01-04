/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wisdom-plus/dottan/internal/config"
)

var force bool

// initCmd represents the init command
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Dottan config file",
	Long:  `Initialize Dottan config file in the user config directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		const appName = "dottan"
		path, err := config.ConfigPath(appName)
		if err != nil {
			return err
		}

		if config.Exists(path) && !force {
			return fmt.Errorf("%w: %s (use --force to overwrite)", config.ErrAlreadyExists, path)
		}
		cfg := &config.Config{
			DefaultProfile: "personal",
			Profiles: map[string]config.Profile{
				"personal": {
					GitHubURL: "https://github.com",
				},
			},
		}
		if err := config.Save(path, cfg); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "created: %s\n", path)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configInitCmd)
	configInitCmd.Flags().BoolVar(&force, "force", false, "overwrite existing config file")
}
