package cli

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/bharath23/git-treekeeper/tests/utils"
)

func TestListCommandTooManyArgs(t *testing.T) {
	root := newRootCmd()
	errOut := utils.CaptureStderr(func() {
		root.SetArgs([]string{"list", "extra"})
		err := root.Execute()
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

	root := newRootCmd()
	restore := utils.Chdir(t, worktreePath)
	out := utils.CaptureStdout(func() {
		root.SetArgs([]string{"list"})
		_ = root.Execute()
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
		root.SetArgs([]string{"list"})
		_ = root.Execute()
	})
	restoreBase()

	if !strings.Contains(out, expectedPath) {
		t.Errorf("expected worktree path from base dir, got: %q", out)
	}
}

func TestListCommandPorcelain(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	_, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	expectedPath := utils.RealPath(t, filepath.Join(destPath, "worktrees", "main"))

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	root := newRootCmd()
	out := utils.CaptureStdout(func() {
		root.SetArgs([]string{"list", "--porcelain"})
		_ = root.Execute()
	})

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d: %q", len(lines), out)
	}
	if lines[0] != "main\t"+expectedPath {
		t.Errorf("unexpected porcelain output: %q", out)
	}
}

func TestListCommandJSON(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	_, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	expectedPath := utils.RealPath(t, filepath.Join(destPath, "worktrees", "main"))

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	root := newRootCmd()
	out := utils.CaptureStdout(func() {
		root.SetArgs([]string{"list", "--json"})
		_ = root.Execute()
	})

	var got []treekeeper.WorktreeInfo
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, out)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 worktree, got %d", len(got))
	}
	if got[0].Branch != "main" || got[0].Path != expectedPath {
		t.Errorf("unexpected json output: %+v", got[0])
	}
}

func TestListCommandFormatConflict(t *testing.T) {
	root := newRootCmd()
	root.SetArgs([]string{"list", "--porcelain", "--json"})
	err := root.Execute()
	if !errors.Is(err, treekeeper.ErrOutputFormatConflict) {
		t.Errorf("expected ErrOutputFormatConflict, got %v", err)
	}
}
