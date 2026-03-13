package cli

import (
	"errors"
	"os"
	"path/filepath"
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

func TestBranchCommandTooManyArgs(t *testing.T) {
	errOut := utils.CaptureStderr(func() {
		cmd.RootCmd.SetArgs([]string{"branch", "one", "two", "three"})
		err := cmd.RootCmd.Execute()
		if !errors.Is(err, treekeeper.ErrTooManyArgs) {
			t.Errorf("expected ErrTooManyArgs, got %v", err)
		}
	})
	if !strings.Contains(errOut, "Usage:") {
		t.Errorf("expected usage output, got: %q", errOut)
	}
}

func TestBranchCommandWithBase(t *testing.T) {
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

	out := utils.CaptureStdout(func() {
		cmd.RootCmd.SetArgs([]string{"branch", "feature1"})
		cmd.RootCmd.Execute()
	})

	if !strings.Contains(out, "Creating branch feature1 from main") {
		t.Errorf("unexpected output: %q", out)
	}
	expectedPath := utils.RealPath(t, filepath.Join(destRoot, "worktrees", "feature1"))
	if !strings.Contains(out, "Worktree path: "+expectedPath) {
		t.Errorf("expected output to contain worktree path, got: %q", out)
	}
}
