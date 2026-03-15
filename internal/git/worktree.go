package git

import (
	"errors"
	"os/exec"
	"strings"
)

type Worktree struct {
	Path           string
	Branch         string
	Locked         bool
	LockedReason   string
	Prunable       bool
	PrunableReason string
}

func AddWorktreeExisting(gitDir, worktreePath, branchName string) error {
	_, err := Run("--git-dir", gitDir, "worktree", "add", worktreePath, branchName)
	return err
}

func AddWorktreeNew(gitDir, worktreePath, branchName, baseBranch string) error {
	_, err := Run("--git-dir", gitDir, "worktree", "add", "-b", branchName, worktreePath, baseBranch)
	return err
}

func WorktreeList(gitDir string) ([]Worktree, error) {
	out, err := Run("--git-dir", gitDir, "worktree", "list", "--porcelain")
	if err != nil {
		return nil, err
	}

	var worktrees []Worktree
	var current *Worktree

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "worktree ") {
			if current != nil {
				worktrees = append(worktrees, *current)
			}
			current = &Worktree{Path: strings.TrimPrefix(line, "worktree ")}
			continue
		}
		if current == nil {
			continue
		}
		if strings.HasPrefix(line, "branch ") {
			ref := strings.TrimPrefix(line, "branch ")
			current.Branch = strings.TrimPrefix(ref, "refs/heads/")
			continue
		}
		if strings.HasPrefix(line, "locked") {
			current.Locked = true
			current.LockedReason = strings.TrimSpace(strings.TrimPrefix(line, "locked"))
			continue
		}
		if strings.HasPrefix(line, "prunable") {
			current.Prunable = true
			current.PrunableReason = strings.TrimSpace(strings.TrimPrefix(line, "prunable"))
			continue
		}
	}

	if current != nil {
		worktrees = append(worktrees, *current)
	}

	return worktrees, nil
}

func BranchExists(gitDir, branchName string) (bool, error) {
	cmd := exec.Command("git", "--git-dir", gitDir, "show-ref", "--verify", "--quiet", "refs/heads/"+branchName)
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
