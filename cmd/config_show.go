package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"github.com/wisdom-plus/dottan/internal/config"
)

var configShowCmd = &cobra.Command{
	Use:   "show [key]",
	Short: "Show config file or specific key",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		const appName = "dottan"

		path, err := config.ConfigPath(appName)
		if err != nil {
			return err
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			fmt.Fprint(cmd.OutOrStdout(), string(b))
			if len(b) > 0 && b[len(b)-1] != '\n' {
				fmt.Fprintln(cmd.OutOrStdout())
			}
			return nil
		}

		key := args[0]

		var m map[string]any
		if err := toml.Unmarshal(b, &m); err != nil {
			return err
		}

		v, ok := lookupByDotKey(m, key)
		if !ok {
			return fmt.Errorf("key not found: %s", key)
		}

		switch vv := v.(type) {
		case string:
			fmt.Fprintln(cmd.OutOrStdout(), vv)
		default:
			out, err := toml.Marshal(v)
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf("%v", v))
				return nil
			}
			fmt.Fprint(cmd.OutOrStdout(), fmt.Sprintf("%v", v))
			if len(out) > 0 && out[len(out)-1] != '\n' {
				fmt.Fprintln(cmd.OutOrStdout())
			}
		}
		return nil
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
}

func lookupByDotKey(root map[string]any, key string) (any, bool) {
	if key == "" {
		return nil, false
	}

	cur := any(root)
	parts := strings.Split(key, ".")
	for _, p := range parts {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil, false
		}
		next, ok := m[p]
		if !ok {
			return nil, false
		}
		cur = next
	}
	return cur, true
}

var _ = errors.New
