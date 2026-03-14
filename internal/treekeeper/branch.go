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

	gitDir, baseDir, err := resolveGitDir(workDir)
	if err != nil {
		return BranchResult{}, err
	}

	resolvedBase := baseBranch
	if resolvedBase == "" {
		resolvedBase, err = git.CurrentBranch(workDir)
		if err != nil || resolvedBase == "" {
			resolvedBase, err = git.DefaultBranch(gitDir)
			if err != nil || resolvedBase == "" {
				resolvedBase = "main"
			}
		}
	}

	worktreesRoot := worktreeRoot(baseDir)
	if err := os.MkdirAll(worktreesRoot, 0o755); err != nil {
		return BranchResult{}, err
	}

	worktreePath := filepath.Join(worktreesRoot, branchName)
	if err := git.AddWorktreeNew(gitDir, worktreePath, branchName, resolvedBase); err != nil {
		return BranchResult{}, err
	}

	return BranchResult{
		Base:         resolvedBase,
		WorktreePath: worktreePath,
	}, nil
}
