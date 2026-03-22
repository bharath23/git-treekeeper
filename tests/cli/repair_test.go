package cli

import (
	"path/filepath"
	"testing"

	"github.com/bharath23/git-treekeeper/internal/git"
	"github.com/bharath23/git-treekeeper/tests/utils"
)

func TestRepairCommandSetsOriginTracking(t *testing.T) {
	srcRepo := utils.InitRepo(t)
	destRoot := t.TempDir()

	gitDir := filepath.Join(destRoot, "repo.git")
	utils.RunGit(t, destRoot, "clone", "--bare", srcRepo, gitDir)

	worktreePath := filepath.Join(destRoot, "worktrees", "main")
	utils.RunGit(t, destRoot, "--git-dir", gitDir, "worktree", "add", worktreePath, "main")

	restore := utils.Chdir(t, worktreePath)
	defer restore()

	upstreamRef, err := git.BranchUpstream(gitDir, "main")
	if err != nil {
		t.Fatalf("branch upstream: %v", err)
	}
	if upstreamRef != "" {
		t.Fatalf("expected no upstream, got %q", upstreamRef)
	}

	root := newRootCmd()
	_ = utils.CaptureStdout(func() {
		root.SetArgs([]string{"repair", "--apply"})
		_ = root.Execute()
	})

	upstreamRef, err = git.BranchUpstream(gitDir, "main")
	if err != nil {
		t.Fatalf("branch upstream: %v", err)
	}
	if upstreamRef != "origin/main" {
		t.Fatalf("expected origin/main, got %q", upstreamRef)
	}

	refspecs, err := git.RemoteFetchRefspecs(gitDir, "origin")
	if err != nil {
		t.Fatalf("origin fetch refspecs: %v", err)
	}
	found := false
	for _, refspec := range refspecs {
		if refspec == "+refs/heads/*:refs/remotes/origin/*" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected origin fetch refspec to be configured")
	}
}
