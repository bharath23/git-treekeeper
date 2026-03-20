package treekeeper

import (
	"os"
	"path/filepath"

	"github.com/bharath23/git-treekeeper/internal/git"
)

type BranchResult struct {
	Base         string
	WorktreePath string
}

func CreateBranch(branch, base string) (BranchResult, error) {
	branchName := branch
	baseBranch := base
	workDir, err := os.Getwd()
	if err != nil {
		return BranchResult{}, err
	}

	ctx, err := ResolveContext(workDir)
	if err != nil {
		return BranchResult{}, err
	}

	resolvedBase := baseBranch
	if resolvedBase == "" {
		resolvedBase = resolveBaseBranch(ctx.GitDir, workDir, true)
	}

	if err := EnsureWorktreesRoot(ctx); err != nil {
		return BranchResult{}, err
	}

	worktreePath := filepath.Join(ctx.WorktreesRoot, branchName)
	if err := git.AddWorktreeNew(ctx.GitDir, worktreePath, branchName, resolvedBase); err != nil {
		return BranchResult{}, err
	}

	return BranchResult{
		Base:         resolvedBase,
		WorktreePath: worktreePath,
	}, nil
}
