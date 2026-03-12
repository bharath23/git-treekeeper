package cli

import (
	"errors"
	"strings"
	"testing"

	"github.com/bharath23/git-treekeeper/cmd"
	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/bharath23/git-treekeeper/tests/utils"
)

func TestCheckoutCommandMissingBranchName(t *testing.T) {
	cmd.RootCmd.SetArgs([]string{"checkout"})
	err := cmd.RootCmd.Execute()
	if !errors.Is(err, treekeeper.ErrMissingCheckoutBranch) {
		t.Errorf("expected ErrMissingCheckoutBranch, got %v", err)
	}
}

func TestCheckoutCommandWithBranch(t *testing.T) {
	out := utils.CaptureStdout(func() {
		cmd.RootCmd.SetArgs([]string{"checkout", "feature1"})
		_ = cmd.RootCmd.Execute()
	})

	if !strings.Contains(out, "Checking out branch feature1") {
		t.Errorf("expected output to contain branch name, got: %q", out)
	}
	if !strings.Contains(out, "Worktree path: /tmp/feature1") {
		t.Errorf("expected output to contain worktree path, got: %q", out)
	}
}
