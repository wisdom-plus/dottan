package cmd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	toml "github.com/pelletier/go-toml/v2"
	"github.com/wisdom-plus/dottan/internal/config"
)

func TestConfigInit_CreatesConfigFile(t *testing.T) {
	// force はグローバルなので必ず初期化して戻す
	oldForce := force
	force = false
	t.Cleanup(func() { force = oldForce })

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))

	const appName = "dottan"
	path, err := config.ConfigPath(appName)
	if err != nil {
		t.Fatalf("ConfigPath error: %v", err)
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
	rootCmd.SetArgs([]string{"config", "init"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute error: %v\nstderr: %s", err, errOut.String())
	}

	// 出力メッセージ
	if !strings.Contains(out.String(), "created: ") {
		t.Fatalf("expected created message, got: %q", out.String())
	}

	// ファイル作成確認
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected config file to exist at %s: %v", path, err)
	}

	var got config.Config
	if err := toml.Unmarshal(b, &got); err != nil {
		t.Fatalf("toml.Unmarshal error: %v\ncontent:\n%s", err, string(b))
	}

	if got.DefaultProfile != "personal" {
		t.Fatalf("DefaultProfile got %q, want %q", got.DefaultProfile, "personal")
	}

	p, ok := got.Profiles["personal"]
	if !ok {
		t.Fatalf("Profiles[personal] not found. got: %#v", got.Profiles)
	}

	if p.GitHubURL != "https://github.com" {
		t.Fatalf("GitHubURL got %q, want %q", p.GitHubURL, "https://github.com")
	}
}

func TestConfigInit_WhenAlreadyExists_WithoutForce_ReturnsErrAlreadyExists(t *testing.T) {
	oldForce := force
	force = false
	t.Cleanup(func() { force = oldForce })

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))

	const appName = "dottan"
	path, err := config.ConfigPath(appName)
	if err != nil {
		t.Fatalf("ConfigPath error: %v", err)
	}

	// 既存ファイルを作る
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll error: %v", err)
	}
	if err := os.WriteFile(path, []byte("default_profile = \"existing\"\n"), 0o600); err != nil {
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
	rootCmd.SetArgs([]string{"config", "init"})

	err = rootCmd.Execute()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, config.ErrAlreadyExists) {
		t.Fatalf("expected errors.Is(err, config.ErrAlreadyExists) == true, got err: %v", err)
	}
}

func TestConfigInit_WhenAlreadyExists_WithForce_Overwrites(t *testing.T) {
	oldForce := force
	force = false
	t.Cleanup(func() { force = oldForce })

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))

	const appName = "dottan"
	path, err := config.ConfigPath(appName)
	if err != nil {
		t.Fatalf("ConfigPath error: %v", err)
	}

	// 既存ファイルを作る
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll error: %v", err)
	}
	if err := os.WriteFile(path, []byte("default_profile = \"existing\"\n"), 0o600); err != nil {
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
	rootCmd.SetArgs([]string{"config", "init", "--force"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute error: %v\nstderr: %s", err, errOut.String())
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if strings.Contains(string(b), `default_profile = "existing"`) {
		t.Fatalf("expected existing config to be overwritten, but found old content:\n%s", string(b))
	}
}
