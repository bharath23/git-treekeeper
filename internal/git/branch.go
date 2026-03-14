package git

import (
	"errors"
	"os/exec"
)

func IsMerged(gitDir, branchName, baseBranch string) (bool, error) {
	cmd := exec.Command("git", "--git-dir", gitDir, "merge-base", "--is-ancestor", branchName, baseBranch)
	err := cmd.Run()
	if err == nil {
		return true, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return false, nil
	}

	return false, err
}

func RefExists(gitDir, ref string) (bool, error) {
	cmd := exec.Command("git", "--git-dir", gitDir, "show-ref", "--verify", "--quiet", ref)
	err := cmd.Run()
	if err == nil {
		return true, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return false, nil
	}

	return false, err
}
