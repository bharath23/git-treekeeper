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

func TestBranchDeleteUnmergedRefused(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	_, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	branchResult, err := treekeeper.CreateBranch("feature-x", "")
	if err != nil {
		t.Fatalf("create branch: %v", err)
	}

	restoreFeature := utils.Chdir(t, branchResult.WorktreePath)
	featureFile := filepath.Join(branchResult.WorktreePath, "feature.txt")
	if err := os.WriteFile(featureFile, []byte("feature"), 0o644); err != nil {
		t.Fatalf("write feature file: %v", err)
	}
	utils.RunGit(t, branchResult.WorktreePath, "add", ".")
	utils.RunGit(t, branchResult.WorktreePath, "-c", "user.name=git-tk", "-c", "user.email=git-tk@example.com", "commit", "-m", "feature")
	restoreFeature()

	root := newRootCmd()
	root.SetArgs([]string{"branch", "-d", "feature-x"})
	err = root.Execute()
	if !errors.Is(err, treekeeper.ErrBranchNotMerged) {
		t.Errorf("expected ErrBranchNotMerged, got %v", err)
	}
}

func TestBranchDeleteDirtyRefused(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	_, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	branchResult, err := treekeeper.CreateBranch("feature-x", "")
	if err != nil {
		t.Fatalf("create branch: %v", err)
	}

	if err := os.WriteFile(filepath.Join(branchResult.WorktreePath, "dirty.txt"), []byte("dirty"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	root := newRootCmd()
	root.SetArgs([]string{"branch", "-d", "feature-x"})
	err = root.Execute()
	if !errors.Is(err, treekeeper.ErrDirtyWorktree) {
		t.Errorf("expected ErrDirtyWorktree, got %v", err)
	}
}

func TestBranchDeleteForce(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	_, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	_, err = treekeeper.CreateBranch("feature-x", "")
	if err != nil {
		t.Fatalf("create branch: %v", err)
	}

	root := newRootCmd()
	out := utils.CaptureStdout(func() {
		root.SetArgs([]string{"branch", "-D", "--yes", "feature-x"})
		_ = root.Execute()
	})

	if !strings.Contains(out, "Deleted branch: feature-x") {
		t.Errorf("expected delete output, got: %q", out)
	}
}

func TestBranchDeleteMissingBranch(t *testing.T) {
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
	root.SetArgs([]string{"branch", "-d", "missing-branch"})
	err = root.Execute()
	if !errors.Is(err, treekeeper.ErrBranchNotFound) {
		t.Errorf("expected ErrBranchNotFound, got %v", err)
	}
}

func TestBranchDeleteRemoteMissing(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	_, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	_, err = treekeeper.CreateBranch("feature-remote", "")
	if err != nil {
		t.Fatalf("create branch: %v", err)
	}

	root := newRootCmd()
	root.SetArgs([]string{"branch", "-d", "--remote", "--yes", "feature-remote"})
	err = root.Execute()
	if !errors.Is(err, treekeeper.ErrRemoteBranchNotFound) {
		t.Errorf("expected ErrRemoteBranchNotFound, got %v", err)
	}
}

func TestBranchDeleteCheckedOutRefused(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	_, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	restoreMain := utils.Chdir(t, worktreePath)
	defer restoreMain()

	branchResult, err := treekeeper.CreateBranch("feature-checked", "")
	if err != nil {
		t.Fatalf("create branch: %v", err)
	}

	restoreFeature := utils.Chdir(t, branchResult.WorktreePath)
	defer restoreFeature()

	root := newRootCmd()
	root.SetArgs([]string{"branch", "-d", "feature-checked"})
	err = root.Execute()
	if !errors.Is(err, treekeeper.ErrBranchCheckedOut) {
		t.Errorf("expected ErrBranchCheckedOut, got %v", err)
	}
}

func TestBranchDeleteMergedAllows(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	_, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	branchResult, err := treekeeper.CreateBranch("feature-merge", "")
	if err != nil {
		t.Fatalf("create branch: %v", err)
	}

	restoreFeature := utils.Chdir(t, branchResult.WorktreePath)
	filePath := filepath.Join(branchResult.WorktreePath, "merged.txt")
	if err := os.WriteFile(filePath, []byte("merged"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	utils.RunGit(t, branchResult.WorktreePath, "add", ".")
	utils.RunGit(t, branchResult.WorktreePath, "-c", "user.name=git-tk", "-c", "user.email=git-tk@example.com", "commit", "-m", "merge me")
	restoreFeature()

	utils.RunGit(t, worktreePath, "merge", "feature-merge")

	root := newRootCmd()
	out := utils.CaptureStdout(func() {
		root.SetArgs([]string{"branch", "-d", "feature-merge"})
		_ = root.Execute()
	})
	if !strings.Contains(out, "Deleted branch: feature-merge") {
		t.Errorf("expected delete output, got: %q", out)
	}
}
