package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bharath23/git-treekeeper/internal/git"
	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/bharath23/git-treekeeper/tests/utils"
)

func TestPruneRemovesStaleWorktree(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	_, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	branchResult, err := treekeeper.CreateBranch("feature-stale", "")
	if err != nil {
		t.Fatalf("create branch: %v", err)
	}
	if err := os.RemoveAll(branchResult.WorktreePath); err != nil {
		t.Fatalf("remove worktree: %v", err)
	}

	root := newRootCmd()
	out := utils.CaptureStdout(func() {
		root.SetArgs([]string{"prune"})
		_ = root.Execute()
	})

	expectedPath := utils.RealPath(t, branchResult.WorktreePath)
	if !strings.Contains(out, "Pruned worktree: "+expectedPath) {
		t.Errorf("expected prune output, got: %q", out)
	}

	gitDir := filepath.Join(destPath, "repo.git")
	worktrees, err := git.WorktreeList(gitDir)
	if err != nil {
		t.Fatalf("worktree list: %v", err)
	}
	for _, wt := range worktrees {
		if wt.Branch == "feature-stale" {
			t.Fatalf("expected worktree entry removed, got branch %s", wt.Branch)
		}
	}
}

func TestPruneDryRunKeepsWorktree(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	_, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	branchResult, err := treekeeper.CreateBranch("feature-dry", "")
	if err != nil {
		t.Fatalf("create branch: %v", err)
	}
	if err := os.RemoveAll(branchResult.WorktreePath); err != nil {
		t.Fatalf("remove worktree: %v", err)
	}

	root := newRootCmd()
	out := utils.CaptureStdout(func() {
		root.SetArgs([]string{"prune", "--dry-run"})
		_ = root.Execute()
	})

	expectedPath := utils.RealPath(t, branchResult.WorktreePath)
	if !strings.Contains(out, "Would prune worktree: "+expectedPath) {
		t.Errorf("expected dry-run output, got: %q", out)
	}

	gitDir := filepath.Join(destPath, "repo.git")
	worktrees, err := git.WorktreeList(gitDir)
	if err != nil {
		t.Fatalf("worktree list: %v", err)
	}
	found := false
	for _, wt := range worktrees {
		if wt.Branch == "feature-dry" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected worktree entry to remain after dry-run")
	}
}

func TestPruneMergedBranches(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	_, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	utils.RunGit(t, worktreePath, "checkout", "-b", "feature-merged")
	if err := os.WriteFile(filepath.Join(worktreePath, "merged.txt"), []byte("merged"), 0o644); err != nil {
		t.Fatalf("write merged file: %v", err)
	}
	utils.RunGit(t, worktreePath, "add", ".")
	utils.RunGit(t, worktreePath, "-c", "user.name=git-tk", "-c", "user.email=git-tk@example.com", "commit", "-m", "merge me")
	utils.RunGit(t, worktreePath, "checkout", "main")
	utils.RunGit(t, worktreePath, "merge", "feature-merged")

	root := newRootCmd()
	out := utils.CaptureStdout(func() {
		root.SetArgs([]string{"prune", "--merged-branches"})
		_ = root.Execute()
	})

	if !strings.Contains(out, "Pruned branch: feature-merged") {
		t.Errorf("expected prune branch output, got: %q", out)
	}

	gitDir := filepath.Join(destPath, "repo.git")
	exists, err := git.RefExists(gitDir, "refs/heads/feature-merged")
	if err != nil {
		t.Fatalf("check branch exists: %v", err)
	}
	if exists {
		t.Fatalf("expected merged branch deleted")
	}
}
