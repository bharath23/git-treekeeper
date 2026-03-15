package git

import (
	"errors"
	"os/exec"
	"strings"
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

func LocalBranches(gitDir string) ([]string, error) {
	out, err := Run("--git-dir", gitDir, "for-each-ref", "--format=%(refname:short)", "refs/heads")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	lines := strings.Split(out, "\n")
	branches := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		branches = append(branches, line)
	}
	return branches, nil
}

func BranchUpstream(gitDir, branchName string) (string, error) {
	out, err := Run("--git-dir", gitDir, "for-each-ref", "--format=%(upstream:short)", "refs/heads/"+branchName)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func SetBranchUpstream(gitDir, branchName, upstreamRef string) error {
	_, err := Run("--git-dir", gitDir, "branch", "--set-upstream-to="+upstreamRef, branchName)
	return err
}

func SetBranchPushRemote(gitDir, branchName, remoteName string) error {
	_, err := Run("--git-dir", gitDir, "config", "branch."+branchName+".pushRemote", remoteName)
	return err
}

func BranchPushRemote(gitDir, branchName string) (string, error) {
	cmd := exec.Command("git", "--git-dir", gitDir, "config", "--get", "branch."+branchName+".pushRemote")
	out, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(out)), nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return "", nil
	}

	return "", err
}
