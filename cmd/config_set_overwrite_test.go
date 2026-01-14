package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	toml "github.com/pelletier/go-toml/v2"
	"github.com/wisdom-plus/dottan/internal/config"
)

func TestConfigSet_OverwriteValue(t *testing.T) {
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
	// 初回 set のために空ファイルを用意（実装が ReadFile するため）
	if err := os.WriteFile(path, []byte(""), 0o600); err != nil {
		t.Fatalf("WriteFile (seed) error: %v", err)
	}

	// rootCmd はグローバルなので毎回リセット
	origOut := rootCmd.OutOrStdout()
	origErr := rootCmd.ErrOrStderr()
	t.Cleanup(func() {
		rootCmd.SetOut(origOut)
		rootCmd.SetErr(origErr)
		rootCmd.SetArgs(nil)
	})

	runSet := func(key, value string) {
		out := new(bytes.Buffer)
		errOut := new(bytes.Buffer)

		rootCmd.SetOut(out)
		rootCmd.SetErr(errOut)
		rootCmd.SetArgs([]string{"config", "set", key, value})

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("Execute error (set %s=%s): %v\nstderr: %s", key, value, err, errOut.String())
		}
		// メッセージ確認は任意（壊れやすいので contains 程度）
		if !strings.Contains(out.String(), "updated:") {
			t.Fatalf("expected stdout to contain %q, got: %q", "updated:", out.String())
		}
	}

	// 1回目
	runSet("foo.bar", "old")
	// 2回目（上書き）
	runSet("foo.bar", "new")

	// ファイルを読んで最終状態を検証
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}

	var m map[string]any
	if err := toml.Unmarshal(b, &m); err != nil {
		t.Fatalf("toml.Unmarshal error: %v\ncontent:\n%s", err, string(b))
	}

	foo, ok := m["foo"].(map[string]any)
	if !ok {
		t.Fatalf("expected [foo] table, got: %#v", m["foo"])
	}

	got, _ := foo["bar"].(string)
	want := "new"
	if got != want {
		t.Fatalf("foo.bar got %q, want %q\ncontent:\n%s", got, want, string(b))
	}
}
