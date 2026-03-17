package git

import (
	"errors"
	"os/exec"
	"strings"
)

func RemoteBranchExists(gitDir, remoteName, branchName string) (bool, error) {
	cmd := exec.Command("git", "--git-dir", gitDir, "ls-remote", "--exit-code", "--heads", remoteName, branchName)
	err := cmd.Run()
	if err == nil {
		return true, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 2 {
		return false, nil
	}

	return false, err
}

func RemoteExists(gitDir, remoteName string) (bool, error) {
	cmd := exec.Command("git", "--git-dir", gitDir, "remote", "get-url", remoteName)
	err := cmd.Run()
	if err == nil {
		return true, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 2 {
		return false, nil
	}

	return false, err
}

func RemoteURL(gitDir, remoteName string) (string, error) {
	return Run("--git-dir", gitDir, "remote", "get-url", remoteName)
}

func AddRemote(gitDir, remoteName, remoteURL string) error {
	_, err := Run("--git-dir", gitDir, "remote", "add", remoteName, remoteURL)
	return err
}

func Remotes(gitDir string) ([]string, error) {
	out, err := Run("--git-dir", gitDir, "remote")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}
