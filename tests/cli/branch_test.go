package cli

import (
	"errors"
	"strings"
	"testing"

	"github.com/bharath23/git-treekeeper/cmd"
	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/bharath23/git-treekeeper/tests/utils"
)

func TestBranchCommandMissingBranchName(t *testing.T) {
	cmd.RootCmd.SetArgs([]string{"branch"})
	err := cmd.RootCmd.Execute()
	if !errors.Is(err, treekeeper.ErrMissingBranchName) {
		t.Errorf("expected ErrMissingBranchName, got %v", err)
	}
}

func TestBranchCommandWithBase(t *testing.T) {
	out := utils.CaptureStdout(func() {
		cmd.RootCmd.SetArgs([]string{"branch", "feature1", "develop"})
		cmd.RootCmd.Execute()
	})

	if !strings.Contains(out, "Creating branch feature1 from develop") {
		t.Errorf("unexpected output: %q", out)
	}
}
