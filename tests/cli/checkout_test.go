package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/bharath23/git-treekeeper/tests/utils"
)

func TestCheckoutCommandMissingBranchName(t *testing.T) {
	root := newRootCmd()
	root.SetArgs([]string{"checkout"})
	err := root.Execute()
	if !errors.Is(err, treekeeper.ErrMissingCheckoutBranch) {
		t.Errorf("expected ErrMissingCheckoutBranch, got %v", err)
	}
}

func TestCheckoutCommandTooManyArgs(t *testing.T) {
	root := newRootCmd()
	errOut := utils.CaptureStderr(func() {
		root.SetArgs([]string{"checkout", "one", "two"})
		err := root.Execute()
		if !errors.Is(err, treekeeper.ErrTooManyArgs) {
			t.Errorf("expected ErrTooManyArgs, got %v", err)
		}
	})
	if !strings.Contains(errOut, "Usage:") {
		t.Errorf("expected usage output, got: %q", errOut)
	}
}

func TestCheckoutCommandWithBranch(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	restore := utils.Chdir(t, destRoot)
	defer restore()

	utils.RunGit(t, destRoot, "clone", "--bare", srcRepo, filepath.Join(destRoot, "repo.git"))
	worktreePath := filepath.Join(destRoot, "worktrees", "main")
	if err := os.MkdirAll(filepath.Dir(worktreePath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	utils.RunGit(t, destRoot, "--git-dir", filepath.Join(destRoot, "repo.git"), "worktree", "add", worktreePath, "main")

	restoreWorktree := utils.Chdir(t, worktreePath)
	defer restoreWorktree()

	root := newRootCmd()
	out := utils.CaptureStdout(func() {
		root.SetArgs([]string{"checkout", "feature1"})
		_ = root.Execute()
	})

	if !strings.Contains(out, "Checking out branch feature1") {
		t.Errorf("expected output to contain branch name, got: %q", out)
	}
	expectedPath := utils.RealPath(t, filepath.Join(destRoot, "worktrees", "feature1"))
	if !strings.Contains(out, "Worktree path: "+expectedPath) {
		t.Errorf("expected output to contain worktree path, got: %q", out)
	}
}
