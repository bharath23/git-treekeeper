package treekeeper

import (
	"os"
	"path/filepath"

	"github.com/bharath23/git-treekeeper/internal/git"
)

func Checkout(branch string) (string, error) {
	branchName := branch
	workDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	gitDir, err := git.CommonDir(workDir)
	if err != nil {
		return "", err
	}

	worktreesRoot := worktreeRoot(filepath.Dir(gitDir))
	if err := os.MkdirAll(worktreesRoot, 0o755); err != nil {
		return "", err
	}

	worktrees, err := git.WorktreeList(gitDir)
	if err != nil {
		return "", err
	}
	for _, wt := range worktrees {
		if wt.Branch == branchName {
			return wt.Path, nil
		}
	}

	worktreePath := filepath.Join(worktreesRoot, branchName)
	branchExists, err := git.BranchExists(gitDir, branchName)
	if err != nil {
		return "", err
	}

	if branchExists {
		if err := git.AddWorktreeExisting(gitDir, worktreePath, branchName); err != nil {
			return "", err
		}
		return worktreePath, nil
	}

	baseBranch, err := git.DefaultBranch(gitDir)
	if err != nil || baseBranch == "" {
		baseBranch = "main"
	}

	if err := git.AddWorktreeNew(gitDir, worktreePath, branchName, baseBranch); err != nil {
		return "", err
	}

	return worktreePath, nil
}
