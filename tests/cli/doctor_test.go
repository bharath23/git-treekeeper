package cli

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/bharath23/git-treekeeper/tests/utils"
)

func TestDoctorCommandTooManyArgs(t *testing.T) {
	root := newRootCmd()
	errOut := utils.CaptureStderr(func() {
		root.SetArgs([]string{"doctor", "extra"})
		err := root.Execute()
		if !errors.Is(err, treekeeper.ErrTooManyArgs) {
			t.Errorf("expected ErrTooManyArgs, got %v", err)
		}
	})
	if !strings.Contains(errOut, "Usage:") {
		t.Errorf("expected usage output, got: %q", errOut)
	}
}

func TestDoctorCommandDirty(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	_, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	if err := os.WriteFile(filepath.Join(worktreePath, "dirty.txt"), []byte("dirty"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	root := newRootCmd()
	out := utils.CaptureStdout(func() {
		root.SetArgs([]string{"doctor"})
		_ = root.Execute()
	})

	if !strings.Contains(out, "branch") || !strings.Contains(out, "state") || !strings.Contains(out, "tracking") {
		t.Errorf("expected header, got: %q", out)
	}
	if !strings.Contains(out, "main") {
		t.Errorf("expected main branch, got: %q", out)
	}
	if !strings.Contains(out, "dirty") {
		t.Errorf("expected dirty state, got: %q", out)
	}
}

func TestDoctorCommandPorcelain(t *testing.T) {
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
		root.SetArgs([]string{"doctor", "--porcelain"})
		_ = root.Execute()
	})

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d: %q", len(lines), out)
	}
	if lines[0] != "main\tclean\torigin/main" {
		t.Errorf("unexpected porcelain output: %q", out)
	}
}

func TestDoctorCommandJSON(t *testing.T) {
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
		root.SetArgs([]string{"doctor", "--json"})
		_ = root.Execute()
	})

	var got []treekeeper.DoctorInfo
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, out)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if got[0].Branch != "main" || got[0].State != "clean" {
		t.Errorf("unexpected json output: %+v", got[0])
	}
	if got[0].Tracking != "origin/main" {
		t.Errorf("expected tracking origin/main, got: %q", got[0].Tracking)
	}
}

func TestDoctorCommandStaleAndOrphaned(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	_, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	// 1. Create a stale worktree (tracked by git, directory gone)
	_, err = treekeeper.CreateBranch("stale-branch", "")
	if err != nil {
		t.Fatalf("create stale branch: %v", err)
	}
	stalePath := filepath.Join(destPath, "worktrees", "stale-branch")
	if err := os.RemoveAll(stalePath); err != nil {
		t.Fatalf("remove stale dir: %v", err)
	}

	// 2. Create an orphaned directory (not tracked by git)
	orphanPath := filepath.Join(destPath, "worktrees", "orphan-dir")
	if err := os.MkdirAll(orphanPath, 0o755); err != nil {
		t.Fatalf("create orphan dir: %v", err)
	}

	root := newRootCmd()
	out := utils.CaptureStdout(func() {
		root.SetArgs([]string{"doctor"})
		_ = root.Execute()
	})

	if !strings.Contains(out, "stale-branch") || !strings.Contains(out, "stale (directory missing)") {
		t.Errorf("expected stale branch detection, got: %q", out)
	}
	if !strings.Contains(out, "orphan-dir") || !strings.Contains(out, "orphaned directory") {
		t.Errorf("expected orphaned directory detection, got: %q", out)
	}
}

func TestDoctorCommandFormatConflict(t *testing.T) {
	root := newRootCmd()
	root.SetArgs([]string{"doctor", "--porcelain", "--json"})
	err := root.Execute()
	if !errors.Is(err, treekeeper.ErrOutputFormatConflict) {
		t.Errorf("expected ErrOutputFormatConflict, got %v", err)
	}
}
