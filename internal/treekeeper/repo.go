package treekeeper

import (
	"os"
	"path/filepath"

	"github.com/bharath23/git-treekeeper/internal/git"
)

func resolveGitDir(workDir string) (string, string, error) {
	gitDir, err := git.CommonDir(workDir)
	if err == nil {
		return gitDir, filepath.Dir(gitDir), nil
	}

	dir := workDir
	for {
		candidate := filepath.Join(dir, "repo.git")
		info, statErr := os.Stat(candidate)
		if statErr == nil && info.IsDir() {
			return candidate, dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", "", err
}
