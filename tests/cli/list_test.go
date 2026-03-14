package cli

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bharath23/git-treekeeper/cmd"
	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/bharath23/git-treekeeper/tests/utils"
)

func TestListCommandTooManyArgs(t *testing.T) {
	errOut := utils.CaptureStderr(func() {
		cmd.RootCmd.SetArgs([]string{"list", "extra"})
		err := cmd.RootCmd.Execute()
		if !errors.Is(err, treekeeper.ErrTooManyArgs) {
			t.Errorf("expected ErrTooManyArgs, got %v", err)
		}
	})
	if !strings.Contains(errOut, "Usage:") {
		t.Errorf("expected usage output, got: %q", errOut)
	}
}

func TestListCommandOutput(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	_, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	expectedPath := utils.RealPath(t, filepath.Join(destPath, "worktrees", "main"))

	restore := utils.Chdir(t, worktreePath)
	out := utils.CaptureStdout(func() {
		cmd.RootCmd.SetArgs([]string{"list"})
		_ = cmd.RootCmd.Execute()
	})
	restore()

	if !strings.Contains(out, "branch") || !strings.Contains(out, "path") {
		t.Errorf("expected header, got: %q", out)
	}
	if !strings.Contains(out, "main") {
		t.Errorf("expected main branch, got: %q", out)
	}
	if !strings.Contains(out, expectedPath) {
		t.Errorf("expected worktree path, got: %q", out)
	}

	restoreBase := utils.Chdir(t, destPath)
	out = utils.CaptureStdout(func() {
		cmd.RootCmd.SetArgs([]string{"list"})
		_ = cmd.RootCmd.Execute()
	})
	restoreBase()

	if !strings.Contains(out, expectedPath) {
		t.Errorf("expected worktree path from base dir, got: %q", out)
	}
}
