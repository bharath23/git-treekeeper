package treekeeper

import (
	"path"
	"path/filepath"
	"strings"
)

func repoBaseFromURL(repoURL string) string {
	trimmed := strings.TrimSuffix(strings.TrimSuffix(repoURL, "/"), ".git")
	trimmed = strings.ReplaceAll(trimmed, ":", "/")
	return path.Base(trimmed)
}

func worktreeRoot(baseDir string) string {
	return filepath.Join(baseDir, "worktrees")
}
