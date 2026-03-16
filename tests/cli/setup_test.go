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

func TestSetupConfiguresUpstreamAndPushRemote(t *testing.T) {
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
		root.SetArgs([]string{"setup", "--upstream-url", upstream})
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

func TestSetupInstallHooks(t *testing.T) {
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
		root.SetArgs([]string{"setup", "--upstream-url", upstream, "--install-hooks"})
		_ = root.Execute()
	})

	gitDir := filepath.Join(destPath, "repo.git")
	hooksPath := utils.RealPath(t, filepath.Join(destPath, ".githooks"))
	configPath := strings.TrimSpace(utils.RunGit(t, destPath, "--git-dir", gitDir, "config", "--get", "core.hooksPath"))
	if configPath == "" {
		t.Fatalf("expected hooksPath config")
	}
	actualPath := utils.RealPath(t, configPath)
	if actualPath != hooksPath {
		t.Fatalf("expected hooksPath %q, got %q", hooksPath, actualPath)
	}

	hookFile := filepath.Join(hooksPath, "pre-commit")
	info, err := os.Stat(hookFile)
	if err != nil {
		t.Fatalf("hook file missing: %v", err)
	}
	if info.Mode()&0o111 == 0 {
		t.Fatalf("expected hook to be executable")
	}
}
