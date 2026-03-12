package cli

import (
	"errors"
	"strings"
	"testing"

	"github.com/bharath23/git-treekeeper/cmd"
	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/bharath23/git-treekeeper/tests/utils"
)

func TestCloneCommandMissingRepoURL(t *testing.T) {
	cmd.RootCmd.SetArgs([]string{"clone"})
	err := cmd.RootCmd.Execute()
	if !errors.Is(err, treekeeper.ErrMissingRepoURL) {
		t.Errorf("expected ErrMissingRepoURL, got %v", err)
	}
}

func TestCloneCommandWithRepo(t *testing.T) {
	out := utils.CaptureStdout(func() {
		cmd.RootCmd.SetArgs([]string{"clone", "https://github.com/bharath23/git-treekeeper.git"})
		_ = cmd.RootCmd.Execute()
	})

	if !strings.Contains(out, "Cloning repo https://github.com/bharath23/git-treekeeper.git") {
		t.Errorf("expected clone message, got: %q", out)
	}
	if !strings.Contains(out, "Default branch: main") {
		t.Errorf("expected default branch info, got: %q", out)
	}
	if !strings.Contains(out, "Worktree path: main") {
		t.Errorf("expected worktree path info, got: %q", out)
	}
}
