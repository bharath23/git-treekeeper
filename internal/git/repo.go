package git

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func CloneBare(repoURL, gitDir string) error {
	_, err := Run("clone", "--bare", repoURL, gitDir)
	return err
}

func DefaultBranch(gitDir string) (string, error) {
	out, err := Run("--git-dir", gitDir, "symbolic-ref", "refs/remotes/origin/HEAD")
	if err == nil {
		return strings.TrimPrefix(out, "refs/remotes/origin/"), nil
	}

	out, err = Run("--git-dir", gitDir, "symbolic-ref", "HEAD")
	if err == nil {
		return strings.TrimPrefix(out, "refs/heads/"), nil
	}

	out, err = Run("--git-dir", gitDir, "rev-parse", "--abbrev-ref", "HEAD")
	if err == nil && out != "HEAD" {
		return out, nil
	}

	return "", err
}

func CommonDir(workDir string) (string, error) {
	out, err := RunInDir(workDir, "rev-parse", "--git-common-dir")
	if err != nil {
		return "", err
	}
	if filepath.IsAbs(out) {
		return out, nil
	}

	top, err := RunInDir(workDir, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	candidate := filepath.Clean(filepath.Join(top, out))
	if _, err := os.Stat(candidate); err == nil {
		return candidate, nil
	}

	candidate = filepath.Clean(filepath.Join(workDir, out))
	return candidate, nil
}

func CurrentBranch(workDir string) (string, error) {
	out, err := RunInDir(workDir, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	if out == "HEAD" {
		return "", errors.New("detached HEAD")
	}
	return out, nil
}
