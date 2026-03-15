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

func TestBranchCommandMissingBranchName(t *testing.T) {
	root := newRootCmd()
	root.SetArgs([]string{"branch"})
	err := root.Execute()
	if !errors.Is(err, treekeeper.ErrMissingBranchName) {
		t.Errorf("expected ErrMissingBranchName, got %v", err)
	}
}

func TestBranchCommandTooManyArgs(t *testing.T) {
	root := newRootCmd()
	errOut := utils.CaptureStderr(func() {
		root.SetArgs([]string{"branch", "one", "two", "three"})
		err := root.Execute()
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

	root := newRootCmd()
	out := utils.CaptureStdout(func() {
		root.SetArgs([]string{"branch", "feature1"})
		_ = root.Execute()
	})

	if !strings.Contains(out, "Creating branch feature1 from main") {
		t.Errorf("unexpected output: %q", out)
	}
	expectedPath := utils.RealPath(t, filepath.Join(destRoot, "worktrees", "feature1"))
	if !strings.Contains(out, "Worktree path: "+expectedPath) {
		t.Errorf("expected output to contain worktree path, got: %q", out)
	}
}

func TestBranchCommandPathOnly(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	_, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	root := newRootCmd()
	out := utils.CaptureStdout(func() {
		root.SetArgs([]string{"branch", "--path-only", "feature-path"})
		_ = root.Execute()
	})

	expectedPath := utils.RealPath(t, filepath.Join(destPath, "worktrees", "feature-path"))
	if strings.TrimSpace(out) != expectedPath {
		t.Errorf("expected path-only output %q, got %q", expectedPath, out)
	}
}
