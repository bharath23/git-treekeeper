package treekeeper

import (
	"os"
	"path/filepath"

	"github.com/bharath23/git-treekeeper/internal/git"
)

func Clone(repoURL, destPath string) (string, string, error) {
	baseDir := destPath
	if baseDir == "" {
		baseDir = repoBaseFromURL(repoURL)
	}
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return "", "", err
	}
	if abs, err := filepath.Abs(baseDir); err == nil {
		baseDir = abs
	}

	gitDir := filepath.Join(baseDir, "repo.git")
	if err := git.CloneBare(repoURL, gitDir); err != nil {
		return "", "", err
	}

	defaultBranch, err := git.DefaultBranch(gitDir)
	if err != nil || defaultBranch == "" {
		defaultBranch = "main"
	}

	worktreesRoot := worktreeRoot(baseDir)
	if err := os.MkdirAll(worktreesRoot, 0o755); err != nil {
		return "", "", err
	}

	worktreePath := filepath.Join(worktreesRoot, defaultBranch)
	if err := git.AddWorktreeExisting(gitDir, worktreePath, defaultBranch); err != nil {
		return "", "", err
	}

	return defaultBranch, worktreePath, nil
}
