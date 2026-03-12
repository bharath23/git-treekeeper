package tests

import (
	"testing"

	"github.com/bharath23/git-treekeeper/internal/treekeeper"
)

func TestCloneInfoMessage(t *testing.T) {
	defaultBranch, worktreePath, err := treekeeper.Clone("https://github.com/foo/bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if defaultBranch != "main" {
		t.Errorf("expected default branch %q, got %q", "main", defaultBranch)
	}
	if worktreePath != "/tmp/main" {
		t.Errorf("expected worktree path %q, got %q", "/tmp/main", worktreePath)
	}
}

func TestCreateBranchInfoMessage(t *testing.T) {
	worktreePath, err := treekeeper.CreateBranch("feature-x", "main")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "/tmp/feature-x"
	if worktreePath != expected {
		t.Errorf("expected worktree path %q, got %q", expected, worktreePath)
	}
}

func TestCheckoutInfoMessage(t *testing.T) {
	path, err := treekeeper.Checkout("feature-x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "/tmp/feature-x"
	if path != expected {
		t.Errorf("expected %q, got %q", expected, path)
	}
}
