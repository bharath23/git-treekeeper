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

func TestDoctorCommandTooManyArgs(t *testing.T) {
	errOut := utils.CaptureStderr(func() {
		cmd.RootCmd.SetArgs([]string{"doctor", "extra"})
		err := cmd.RootCmd.Execute()
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

	out := utils.CaptureStdout(func() {
		cmd.RootCmd.SetArgs([]string{"doctor"})
		_ = cmd.RootCmd.Execute()
	})

	if !strings.Contains(out, "branch") || !strings.Contains(out, "state") {
		t.Errorf("expected header, got: %q", out)
	}
	if !strings.Contains(out, "main") {
		t.Errorf("expected main branch, got: %q", out)
	}
	if !strings.Contains(out, "dirty") {
		t.Errorf("expected dirty state, got: %q", out)
	}
}
