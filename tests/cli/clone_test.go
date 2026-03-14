package cli

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/bharath23/git-treekeeper/tests/utils"
)

func TestCloneCommandMissingRepoURL(t *testing.T) {
	root := newRootCmd()
	root.SetArgs([]string{"clone"})
	err := root.Execute()
	if !errors.Is(err, treekeeper.ErrMissingRepoURL) {
		t.Errorf("expected ErrMissingRepoURL, got %v", err)
	}
}

func TestCloneCommandTooManyArgs(t *testing.T) {
	root := newRootCmd()
	errOut := utils.CaptureStderr(func() {
		root.SetArgs([]string{"clone", "repo", "path", "extra"})
		err := root.Execute()
		if !errors.Is(err, treekeeper.ErrTooManyArgs) {
			t.Errorf("expected ErrTooManyArgs, got %v", err)
		}
	})
	if !strings.Contains(errOut, "Usage:") {
		t.Errorf("expected usage output, got: %q", errOut)
	}
}

func TestCloneCommandWithRepo(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	restore := utils.Chdir(t, destRoot)
	defer restore()

	root := newRootCmd()
	out := utils.CaptureStdout(func() {
		root.SetArgs([]string{"clone", srcRepo})
		_ = root.Execute()
	})

	if !strings.Contains(out, "Cloning repo "+srcRepo) {
		t.Errorf("expected clone message, got: %q", out)
	}
	if !strings.Contains(out, "Default branch: main") {
		t.Errorf("expected default branch info, got: %q", out)
	}
	repoName := filepath.Base(srcRepo)
	expectedPath := utils.RealPath(t, filepath.Join(destRoot, repoName, "worktrees", "main"))
	if !strings.Contains(out, "Worktree path: "+expectedPath) {
		t.Errorf("expected worktree path info, got: %q", out)
	}
}
