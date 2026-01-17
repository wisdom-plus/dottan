package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wisdom-plus/dottan/internal/config"
)

func TestConfigShow_Table_PrintsToml(t *testing.T) {
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

	seed := `
[foo]
bar = "baz"
num = 1
`
	if err := os.WriteFile(path, []byte(seed), 0o600); err != nil {
		t.Fatalf("WriteFile seed error: %v", err)
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
	rootCmd.SetArgs([]string{"config", "show", "foo"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute error: %v\nstderr: %s", err, errOut.String())
	}

	got := out.String()

	// 期待値は「TOMLとしての内容が含まれる」までにしておく（順序差で落ちにくくする）
	if !strings.Contains(got, `bar = 'baz'`) {
		t.Fatalf("output missing bar. got:\n%s", got)
	}
	if !strings.Contains(got, `num = 1`) {
		t.Fatalf("output missing num. got:\n%s", got)
	}
	// TOML寄せなので map[...] は出ない
	if strings.Contains(got, "map[") {
		t.Fatalf("output looks like Go map, expected TOML. got:\n%s", got)
	}
}
