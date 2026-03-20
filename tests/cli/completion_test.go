package cli

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/bharath23/git-treekeeper/tests/utils"
)

func TestCompletionSubcommands(t *testing.T) {
	root := newRootCmd()

	out := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	root.SetOut(out)
	root.SetErr(errBuf)
	root.SetArgs([]string{"__complete", "che"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "checkout") {
		t.Errorf("expected 'checkout' in completions, got: %q", output)
	}
}

func TestCompletionBranches(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	tmp := t.TempDir()

	err := exec.Command("git", "clone", srcRepo, tmp).Run()
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	restore := utils.Chdir(t, tmp)
	defer restore()

	// Create a feature branch
	err = exec.Command("git", "checkout", "-b", "feature-abc").Run()
	if err != nil {
		t.Fatalf("create branch: %v", err)
	}

	root := newRootCmd()

	out := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	root.SetOut(out)
	root.SetErr(errBuf)
	root.SetArgs([]string{"__complete", "checkout", "feat"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "feature-abc") {
		t.Errorf("expected 'feature-abc' in branch completions, got: %q", output)
	}
}

func TestCompletionRemotes(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	tmp := t.TempDir()

	err := exec.Command("git", "clone", srcRepo, tmp).Run()
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	restore := utils.Chdir(t, tmp)
	defer restore()

	// Add a custom remote
	err = exec.Command("git", "remote", "add", "upstream", srcRepo).Run()
	if err != nil {
		t.Fatalf("add remote: %v", err)
	}

	root := newRootCmd()

	out := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	root.SetOut(out)
	root.SetErr(errBuf)
	root.SetArgs([]string{"__complete", "sync", "--remote", "up"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "upstream") {
		t.Errorf("expected 'upstream' in remote completions, got: %q", output)
	}
}
