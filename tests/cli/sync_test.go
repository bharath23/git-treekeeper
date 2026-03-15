package cli

import (
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
