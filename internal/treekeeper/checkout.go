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

	ctx, err := ResolveContext(workDir)
	if err != nil {
		return "", err
	}

	if err := EnsureWorktreesRoot(ctx); err != nil {
		return "", err
	}

	worktrees, err := git.WorktreeList(ctx.GitDir)
	if err != nil {
		return "", err
	}
	for _, wt := range worktrees {
		if wt.Branch == branchName {
			return wt.Path, nil
		}
	}

	worktreePath := filepath.Join(ctx.WorktreesRoot, branchName)
	branchExists, err := git.BranchExists(ctx.GitDir, branchName)
	if err != nil {
		return "", err
	}

	if branchExists {
		if err := git.AddWorktreeExisting(ctx.GitDir, worktreePath, branchName); err != nil {
			return "", err
		}
		if err := ensureOriginUpstream(ctx.GitDir, branchName); err != nil {
			return "", err
		}
		return worktreePath, nil
	}

	baseBranch := resolveBaseBranch(ctx.GitDir, workDir, false)
	if err := git.AddWorktreeNew(ctx.GitDir, worktreePath, branchName, baseBranch); err != nil {
		return "", err
	}
	if err := ensureOriginUpstream(ctx.GitDir, branchName); err != nil {
		return "", err
	}

	return worktreePath, nil
}
