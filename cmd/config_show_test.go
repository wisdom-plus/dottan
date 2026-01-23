package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wisdom-plus/dottan/internal/config"
)

func TestConfigShow_PrintsWholeFile_AddsTrailingNewlineIfMissing(t *testing.T) {
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

	seed := `foo = "bar"` // ← 末尾改行なし
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
	rootCmd.SetArgs([]string{"config", "show"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute error: %v\nstderr: %s", err, errOut.String())
	}

	got := out.String()
	want := seed + "\n"
	if got != want {
		t.Fatalf("output got %q, want %q", got, want)
	}
}

func TestConfigShow_PrintsStringValueByKey(t *testing.T) {
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
	rootCmd.SetArgs([]string{"config", "show", "foo.bar"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute error: %v\nstderr: %s", err, errOut.String())
	}

	got := strings.TrimSpace(out.String())
	want := "baz"
	if got != want {
		t.Fatalf("output got %q, want %q", got, want)
	}
}

func TestConfigShow_KeyNotFound_ReturnsError(t *testing.T) {
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
	rootCmd.SetArgs([]string{"config", "show", "foo.missing"})

	if err := rootCmd.Execute(); err == nil {
		t.Fatalf("expected error, got nil")
	}
}
