package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bharath23/git-treekeeper/internal/git"
	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/bharath23/git-treekeeper/tests/utils"
)

func TestGCDryRunPrunesMergedOldBranches(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	_, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	oldDate := time.Now().AddDate(0, 0, -40).UTC().Format(time.RFC3339)
	env := map[string]string{
		"GIT_AUTHOR_DATE":    oldDate,
		"GIT_COMMITTER_DATE": oldDate,
	}

	utils.RunGit(t, worktreePath, "checkout", "-b", "feature-old")
	if err := os.WriteFile(filepath.Join(worktreePath, "old.txt"), []byte("old"), 0o644); err != nil {
		t.Fatalf("write old file: %v", err)
	}
	utils.RunGit(t, worktreePath, "add", ".")
	utils.RunGitEnv(t, worktreePath, env, "-c", "user.name=git-tk", "-c", "user.email=git-tk@example.com", "commit", "-m", "old commit")
	utils.RunGit(t, worktreePath, "checkout", "main")
	utils.RunGit(t, worktreePath, "merge", "--squash", "feature-old")
	utils.RunGit(t, worktreePath, "-c", "user.name=git-tk", "-c", "user.email=git-tk@example.com", "commit", "-m", "squash merge")

	root := newRootCmd()
	out := utils.CaptureStdout(func() {
		root.SetArgs([]string{"gc"})
		_ = root.Execute()
	})

	if !strings.Contains(out, "Would gc branch: feature-old") {
		t.Errorf("expected gc dry-run output, got: %q", out)
	}

	gitDir := filepath.Join(destPath, "repo.git")
	exists, err := git.RefExists(gitDir, "refs/heads/feature-old")
	if err != nil {
		t.Fatalf("check branch exists: %v", err)
	}
	if !exists {
		t.Fatalf("expected branch to remain after dry-run")
	}
}

func TestGCApplyDeletesMergedOldBranches(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	_, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	oldDate := time.Now().AddDate(0, 0, -40).UTC().Format(time.RFC3339)
	env := map[string]string{
		"GIT_AUTHOR_DATE":    oldDate,
		"GIT_COMMITTER_DATE": oldDate,
	}

	utils.RunGit(t, worktreePath, "checkout", "-b", "feature-old")
	if err := os.WriteFile(filepath.Join(worktreePath, "old.txt"), []byte("old"), 0o644); err != nil {
		t.Fatalf("write old file: %v", err)
	}
	utils.RunGit(t, worktreePath, "add", ".")
	utils.RunGitEnv(t, worktreePath, env, "-c", "user.name=git-tk", "-c", "user.email=git-tk@example.com", "commit", "-m", "old commit")
	utils.RunGit(t, worktreePath, "checkout", "main")
	utils.RunGit(t, worktreePath, "merge", "--squash", "feature-old")
	utils.RunGit(t, worktreePath, "-c", "user.name=git-tk", "-c", "user.email=git-tk@example.com", "commit", "-m", "squash merge")

	root := newRootCmd()
	out := utils.CaptureStdout(func() {
		root.SetArgs([]string{"gc", "--apply"})
		_ = root.Execute()
	})

	if !strings.Contains(out, "GC branch: feature-old") {
		t.Errorf("expected gc output, got: %q", out)
	}

	gitDir := filepath.Join(destPath, "repo.git")
	exists, err := git.RefExists(gitDir, "refs/heads/feature-old")
	if err != nil {
		t.Fatalf("check branch exists: %v", err)
	}
	if exists {
		t.Fatalf("expected branch deleted")
	}
}

func TestGCSkipsActiveWorktrees(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()
	destPath := filepath.Join(destRoot, "repo")
	_, worktreePath, err := treekeeper.Clone(srcRepo, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	branchResult, err := treekeeper.CreateBranch("feature-active", "")
	if err != nil {
		t.Fatalf("create branch: %v", err)
	}

	oldDate := time.Now().AddDate(0, 0, -40).UTC().Format(time.RFC3339)
	env := map[string]string{
		"GIT_AUTHOR_DATE":    oldDate,
		"GIT_COMMITTER_DATE": oldDate,
	}
	if err := os.WriteFile(filepath.Join(branchResult.WorktreePath, "active.txt"), []byte("active"), 0o644); err != nil {
		t.Fatalf("write active file: %v", err)
	}
	utils.RunGit(t, branchResult.WorktreePath, "add", ".")
	utils.RunGitEnv(t, branchResult.WorktreePath, env, "-c", "user.name=git-tk", "-c", "user.email=git-tk@example.com", "commit", "-m", "old commit")

	root := newRootCmd()
	_ = utils.CaptureStdout(func() {
		root.SetArgs([]string{"gc", "--apply"})
		_ = root.Execute()
	})

	gitDir := filepath.Join(destPath, "repo.git")
	exists, err := git.RefExists(gitDir, "refs/heads/feature-active")
	if err != nil {
		t.Fatalf("check branch exists: %v", err)
	}
	if !exists {
		t.Fatalf("expected active branch to remain")
	}
}
