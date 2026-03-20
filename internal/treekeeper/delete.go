package treekeeper

import (
	"fmt"
	"os"
	"path/filepath"

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

	branchExists, err := git.RefExists(gitDir, "refs/heads/"+branchName)
	if err != nil {
		return DeleteResult{}, err
	}
	if !branchExists {
		return DeleteResult{}, ErrBranchNotFound
	}

	worktreePath, ok := worktreeForBranch(gitDir, branchName)

	if ok {
		currentRoot, err := git.TopLevel(workDir)
		if err == nil && samePath(currentRoot, worktreePath) {
			return DeleteResult{}, ErrBranchCheckedOut
		}

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

	useForceDelete := force
	if !force {
		baseBranch, err := git.CurrentBranch(workDir)
		if err != nil || baseBranch == "" {
			baseBranch, err = git.DefaultBranch(gitDir)
			if err != nil || baseBranch == "" {
				baseBranch = "main"
			}
		}
		mergeCheck, err := git.IsMerged(gitDir, branchName, baseBranch)
		if err != nil {
			return DeleteResult{}, err
		}
		if !mergeCheck.Merged {
			return DeleteResult{}, ErrBranchNotMerged
		}
		if !mergeCheck.Ancestor {
			useForceDelete = true
		}

		if deleteRemote {
			remoteRef := "refs/remotes/origin/" + baseBranch
			remoteExists, err := git.RefExists(gitDir, remoteRef)
			if err != nil {
				return DeleteResult{}, err
			}
			if remoteExists {
				mergedRemote, err := git.IsMerged(gitDir, branchName, "origin/"+baseBranch)
				if err != nil {
					return DeleteResult{}, err
				}
				if !mergedRemote.Merged {
					return DeleteResult{}, ErrBranchNotMerged
				}
			}
		}
	}

	if ok {
		if _, err := git.Run("--git-dir", gitDir, "worktree", "remove", worktreePath); err != nil {
			return DeleteResult{}, err
		}
	}

	deleteArgs := []string{"--git-dir", gitDir, "branch"}
	if useForceDelete {
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
		remoteExists, err := git.RemoteBranchExists(gitDir, result.RemoteName, branchName)
		if err != nil {
			return DeleteResult{}, err
		}
		if !remoteExists {
			return DeleteResult{}, ErrRemoteBranchNotFound
		}
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

func samePath(a, b string) bool {
	if resolved, err := filepath.EvalSymlinks(a); err == nil {
		a = resolved
	}
	if resolved, err := filepath.EvalSymlinks(b); err == nil {
		b = resolved
	}
	return filepath.Clean(a) == filepath.Clean(b)
}
