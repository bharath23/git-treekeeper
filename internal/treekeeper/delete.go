package treekeeper

import (
	"fmt"
	"os"

	"github.com/bharath23/git-treekeeper/internal/git"
)

type DeleteResult struct {
	WorktreePath  string
	RemoteDeleted bool
	RemoteName    string
}

var protectedBranches = map[string]struct{}{
	"main":   {},
	"master": {},
	"trunk":  {},
}

func DeleteBranch(branchName string, deleteRemote bool, force bool) (DeleteResult, error) {
	if _, ok := protectedBranches[branchName]; ok {
		return DeleteResult{}, ErrProtectedBranch
	}

	workDir, err := os.Getwd()
	if err != nil {
		return DeleteResult{}, err
	}

	gitDir, _, err := resolveGitDir(workDir)
	if err != nil {
		return DeleteResult{}, err
	}

	worktreePath, ok := worktreeForBranch(gitDir, branchName)

	if ok {
		inProgress, err := inProgress(worktreePath)
		if err != nil {
			return DeleteResult{}, err
		}
		if inProgress {
			return DeleteResult{}, ErrInProgress
		}

		dirty, err := isDirty(worktreePath)
		if err != nil {
			return DeleteResult{}, err
		}
		if dirty {
			return DeleteResult{}, ErrDirtyWorktree
		}
	}

	if !force {
		baseBranch, err := git.CurrentBranch(workDir)
		if err != nil || baseBranch == "" {
			baseBranch, err = git.DefaultBranch(gitDir)
			if err != nil || baseBranch == "" {
				baseBranch = "main"
			}
		}
		merged, err := git.IsMerged(gitDir, branchName, baseBranch)
		if err != nil {
			return DeleteResult{}, err
		}
		if !merged {
			return DeleteResult{}, ErrBranchNotMerged
		}
	}

	if ok {
		if _, err := git.Run("--git-dir", gitDir, "worktree", "remove", worktreePath); err != nil {
			return DeleteResult{}, err
		}
	}

	deleteArgs := []string{"--git-dir", gitDir, "branch"}
	if force {
		deleteArgs = append(deleteArgs, "-D")
	} else {
		deleteArgs = append(deleteArgs, "-d")
	}
	deleteArgs = append(deleteArgs, branchName)
	if _, err := git.Run(deleteArgs...); err != nil {
		return DeleteResult{}, err
	}

	result := DeleteResult{
		WorktreePath: worktreePath,
		RemoteName:   "origin",
	}

	if deleteRemote {
		if _, err := git.Run("--git-dir", gitDir, "push", "origin", "--delete", branchName); err != nil {
			return DeleteResult{}, fmt.Errorf("delete remote branch: %w", err)
		}
		result.RemoteDeleted = true
	}

	return result, nil
}

func worktreeForBranch(gitDir, branchName string) (string, bool) {
	worktrees, err := git.WorktreeList(gitDir)
	if err != nil {
		return "", false
	}
	for _, wt := range worktrees {
		if wt.Branch == branchName {
			return wt.Path, true
		}
	}
	return "", false
}
