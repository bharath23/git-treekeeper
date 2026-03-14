package git

import (
	"os"
	"path/filepath"
	"strings"
)

func StatusPorcelain(worktreePath string) (string, error) {
	return RunInDir(worktreePath, "status", "--porcelain")
}

func ResolveGitDir(worktreePath string) (string, error) {
	gitPath := filepath.Join(worktreePath, ".git")
	info, err := os.Stat(gitPath)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return gitPath, nil
	}

	content, err := os.ReadFile(gitPath)
	if err != nil {
		return "", err
	}
	line := strings.TrimSpace(string(content))
	const prefix = "gitdir:"
	if !strings.HasPrefix(line, prefix) {
		return gitPath, nil
	}

	gitDir := strings.TrimSpace(strings.TrimPrefix(line, prefix))
	if filepath.IsAbs(gitDir) {
		return gitDir, nil
	}
	return filepath.Clean(filepath.Join(worktreePath, gitDir)), nil
}
