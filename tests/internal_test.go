package tests

import (
	"path/filepath"
	"testing"

	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/bharath23/git-treekeeper/tests/utils"
)

func TestCloneInfoMessage(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	defaultBranch, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if defaultBranch != "main" {
		t.Errorf("expected default branch %q, got %q", "main", defaultBranch)
	}
	expectedPath := utils.RealPath(t, filepath.Join(destPath, "worktrees", "main"))
	worktreePath = utils.RealPath(t, worktreePath)
	if worktreePath != expectedPath {
		t.Errorf("expected worktree path %q, got %q", expectedPath, worktreePath)
	}
}

func TestCreateBranchInfoMessage(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	_, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	result, err := treekeeper.CreateBranch("feature-x", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Base != "main" {
		t.Errorf("expected base %q, got %q", "main", result.Base)
	}
	expected := utils.RealPath(t, filepath.Join(destPath, "worktrees", "feature-x"))
	if result.WorktreePath != expected {
		t.Errorf("expected worktree path %q, got %q", expected, result.WorktreePath)
	}
}

func TestCheckoutInfoMessage(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	_, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	path, err := treekeeper.Checkout("feature-x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := utils.RealPath(t, filepath.Join(destPath, "worktrees", "feature-x"))
	if path != expected {
		t.Errorf("expected %q, got %q", expected, path)
	}
}
