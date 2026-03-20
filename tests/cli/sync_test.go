package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bharath23/git-treekeeper/internal/git"
	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/bharath23/git-treekeeper/tests/utils"
)

func TestSyncSetsUpstreamAndPushRemote(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	tmp := t.TempDir()

	origin := filepath.Join(tmp, "origin.git")
	utils.RunGit(t, tmp, "clone", "--bare", srcRepo, origin)

	upstream := filepath.Join(tmp, "upstream.git")
	utils.RunGit(t, tmp, "clone", "--bare", srcRepo, upstream)

	destPath := filepath.Join(tmp, "repo")
	_, worktreePath, err := treekeeper.Clone(origin, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	root := newRootCmd()
	_ = utils.CaptureStdout(func() {
		root.SetArgs([]string{"sync", "--default", "--add-upstream", upstream, "--set-upstream"})
		_ = root.Execute()
	})

	gitDir := filepath.Join(destPath, "repo.git")
	upstreamRef, err := git.BranchUpstream(gitDir, "main")
	if err != nil {
		t.Fatalf("branch upstream: %v", err)
	}
	if upstreamRef != "upstream/main" {
		t.Fatalf("expected upstream/main, got %q", upstreamRef)
	}

	pushRemote, err := git.BranchPushRemote(gitDir, "main")
	if err != nil {
		t.Fatalf("branch pushRemote: %v", err)
	}
	if pushRemote != "origin" {
		t.Fatalf("expected pushRemote origin, got %q", pushRemote)
	}
}

func TestSyncUpdatesFromUpstream(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	tmp := t.TempDir()

	origin := filepath.Join(tmp, "origin.git")
	utils.RunGit(t, tmp, "clone", "--bare", srcRepo, origin)

	upstream := filepath.Join(tmp, "upstream.git")
	utils.RunGit(t, tmp, "clone", "--bare", srcRepo, upstream)

	destPath := filepath.Join(tmp, "repo")
	_, worktreePath, err := treekeeper.Clone(origin, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	root := newRootCmd()
	_ = utils.CaptureStdout(func() {
		root.SetArgs([]string{"sync", "--default", "--add-upstream", upstream, "--set-upstream"})
		_ = root.Execute()
	})

	utils.RunGit(t, srcRepo, "remote", "add", "upstream", upstream)
	utils.RunGit(t, srcRepo, "checkout", "main")
	utils.RunGit(t, srcRepo, "-c", "user.name=git-tk", "-c", "user.email=git-tk@example.com", "commit", "--allow-empty", "-m", "upstream change")
	utils.RunGit(t, srcRepo, "push", "upstream", "main")

	_ = utils.CaptureStdout(func() {
		root.SetArgs([]string{"sync", "--default"})
		_ = root.Execute()
	})

	upstreamHead := strings.TrimSpace(utils.RunGit(t, tmp, "--git-dir", upstream, "rev-parse", "refs/heads/main"))
	localHead := strings.TrimSpace(utils.RunGit(t, worktreePath, "rev-parse", "HEAD"))
	if upstreamHead != localHead {
		t.Fatalf("expected main to match upstream, got %q vs %q", localHead, upstreamHead)
	}
}

func TestSyncAllSkipsDirty(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	tmp := t.TempDir()

	origin := filepath.Join(tmp, "origin.git")
	utils.RunGit(t, tmp, "clone", "--bare", srcRepo, origin)

	destPath := filepath.Join(tmp, "repo")
	_, worktreePath, err := treekeeper.Clone(origin, destPath)
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	// 1. Create a clean feature branch
	_, err = treekeeper.CreateBranch("feature-clean", "")
	if err != nil {
		t.Fatalf("create clean branch: %v", err)
	}
	// Push it to origin so it can be fetched back
	utils.RunGit(t, filepath.Join(destPath, "worktrees", "feature-clean"), "push", "-u", "origin", "feature-clean")

	// 2. Create a dirty feature branch
	branchResult, err := treekeeper.CreateBranch("feature-dirty", "")
	if err != nil {
		t.Fatalf("create dirty branch: %v", err)
	}
	if err := os.WriteFile(filepath.Join(branchResult.WorktreePath, "dirty.txt"), []byte("dirty"), 0o644); err != nil {
		t.Fatalf("write dirty file: %v", err)
	}

	// Add a commit to origin to make sure there is something to sync
	utils.RunGit(t, srcRepo, "-c", "user.name=git-tk", "-c", "user.email=git-tk@example.com", "commit", "--allow-empty", "-m", "remote change")
	utils.RunGit(t, srcRepo, "push", origin, "main")

	root := newRootCmd()
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	root.SetOut(outBuf)
	root.SetErr(errBuf)
	treekeeper.SetOutput(outBuf, errBuf)
	defer treekeeper.ResetOutput()

	root.SetArgs([]string{"sync", "--all"})
	_ = root.Execute()

	out := outBuf.String()
	errOut := errBuf.String()

	// Clean branches should sync
	if !strings.Contains(out, "Syncing main from origin/main") {
		t.Errorf("expected main branch sync, got: %q", out)
	}
	if !strings.Contains(out, "Syncing feature-clean from origin/feature-clean") {
		t.Errorf("expected feature-clean sync, got: %q", out)
	}
	// Dirty branches should be skipped and reported in stderr
	if !strings.Contains(errOut, "Skipped branches:") || !strings.Contains(errOut, "feature-dirty: dirty") {
		t.Errorf("expected dirty branch skip warning in stderr, got: %q", errOut)
	}
}
