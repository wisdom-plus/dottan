package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	old := version
	version = "1.2.3"
	t.Cleanup(func() { version = old })

	buf := new(bytes.Buffer)

	cmd := newVersionCmd()
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	got := strings.TrimSpace(buf.String())
	want := "1.2.3"

	if got != want {
		t.Fatalf("output got %q, want %q", got, want)
	}

}

func TestVersionCommand_RejectArgs(t *testing.T) {
	cmd := newVersionCmd()

	cmd.SetArgs([]string{"extra"})
	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestVersionCommand_Call_RootCmd(t *testing.T) {
	old := version
	version = "1.2.3"
	t.Cleanup(func() { version = old })

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
	rootCmd.SetArgs([]string{"version"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute() error: %v\nstderr: %s", err, errOut.String())
	}

	got := strings.TrimSpace(out.String())
	want := "1.2.3"
	if got != want {
		t.Fatalf("output got %q, want %q", got, want)
	}
}
