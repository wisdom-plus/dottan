package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	toml "github.com/pelletier/go-toml/v2"
	"github.com/wisdom-plus/dottan/internal/config"
)

func TestConfigSetCommand_WritesToml(t *testing.T) {
	home := t.TempDir()

	t.Setenv("HOME", home)

	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))

	const appName = "dottan"
	path, err := config.ConfigPath(appName)
	if err != nil {
		t.Fatalf("ConfigPath error: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll error: %v", err)
	}
	if err := os.WriteFile(path, []byte(""), 0o600); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	out := new(bytes.Buffer)
	errOut := new(bytes.Buffer)

	origOut := rootCmd.OutOrStdout()
	origErr := rootCmd.ErrOrStderr()
	t.Cleanup(func() {
		rootCmd.SetOut(origOut)
		rootCmd.SetErr(origErr)
		rootCmd.SetArgs(nil)
	})

	rootCmd.SetOut(out)
	rootCmd.SetErr(errOut)
	rootCmd.SetArgs([]string{"config", "set", "foo.bar", "baz"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute error: %v\nstderr: %s", err, errOut.String())
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}

	var m map[string]any
	if err := toml.Unmarshal(b, &m); err != nil {
		t.Fatalf("toml.Unmarshal error: %v\ncontent:\n%s", err, string(b))
	}

	// m["foo"] は map[string]any のはず
	foo, ok := m["foo"].(map[string]any)
	if !ok {
		t.Fatalf("expected [foo] table, got: %#v", m["foo"])
	}
	if foo["bar"] != "baz" {
		t.Fatalf("expected foo.bar == %q, got %#v", "baz", foo["bar"])
	}

	// stdout メッセージも軽く確認（完全一致より contains の方が安定）
	if got := out.String(); got == "" {
		t.Fatalf("expected stdout, got empty")
	}
}
