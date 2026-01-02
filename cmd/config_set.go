package cmd

import (
	"fmt"
	toml "github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"github.com/wisdom-plus/dottan/internal/config"
	"os"
	"path/filepath"
	"strings"
)

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		const appName = "dottan"
		key := args[0]
		value := args[1]

		path, err := config.ConfigPath(appName)
		if err != nil {
			return err
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read config: %w", err)
		}

		var m map[string]any
		if err := toml.Unmarshal(b, &m); err != nil {
			return err
		}
		if err := setByDotKey(m, key, value); err != nil {
			return err
		}

		out, err := toml.Marshal(m)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, out, 0o600); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "updated: %s = %s\n", key, value)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
}

func setByDotKey(root map[string]any, key string, value any) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return fmt.Errorf("key must not be empty")
	}

	parts := strings.Split(key, ".")
	cur := root
	for i := 0; i < len(parts); i++ {
		p := strings.TrimSpace(parts[i])
		if p == "" {
			return fmt.Errorf("invalid key: %q", key)
		}

		last := i == len(parts)-1
		if last {
			cur[p] = value
			return nil
		}

		next, ok := cur[p]
		if !ok {
			child := map[string]any{}
			cur[p] = child
			cur = child
			continue
		}

		child, ok := next.(map[string]any)
		if !ok {
			return fmt.Errorf("cannot set %q: %q is not a table", key, strings.Join(parts[:i+1], "."))
		}
		cur = child
	}
	return nil
}
