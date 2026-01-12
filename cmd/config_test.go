package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCommand_HasConfigSubcommand(t *testing.T) {
	c, _, err := rootCmd.Find([]string{"config"})
	if err != nil {
		t.Fatalf("rootCmd.Find() error: %v", err)
	}

	if c == nil {
		t.Fatalf("config command not found")
	}

	if c.Use != "config" {
		t.Fatalf("Use go %q, want %q", c.Use, "config")
	}

	if c.Short != "Manage config file" {
		t.Fatalf("Short got %q, want %q", c.Short, "Manage config file")
	}
}

func TestRootCommand_ConfigHelp(t *testing.T) {
	out := new(bytes.Buffer)
	errOut := new(bytes.Buffer)

	// rootCmd はグローバルなのでテスト中の変更を戻す
	origOut := rootCmd.OutOrStdout()
	origErr := rootCmd.ErrOrStderr()
	t.Cleanup(func() {
		rootCmd.SetOut(origOut)
		rootCmd.SetErr(origErr)
		rootCmd.SetArgs(nil)
	})

	rootCmd.SetOut(out)
	rootCmd.SetErr(errOut)
	rootCmd.SetArgs([]string{"config", "--help"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute() error: %v\nstderr: %s", err, errOut.String())
	}

	help := out.String()

	// Cobra の help 出力は環境や設定で少し変わるので、「含まれること」を確認するのが安定
	if !strings.Contains(help, "Manage config file") {
		t.Fatalf("help does not contain Short text. help:\n%s", help)
	}
	if !strings.Contains(help, "config") {
		t.Fatalf("help does not contain command name. help:\n%s", help)
	}
}
