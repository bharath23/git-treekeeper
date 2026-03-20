package treekeeper

import (
	"os"

	"github.com/bharath23/git-treekeeper/internal/git"
)

type PruneOptions struct {
	DryRun         bool
	MergedBranches bool
}

func Prune(options PruneOptions) (PruneResult, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return PruneResult{}, err
	}

	gitDir, _, err := resolveGitDir(workDir)
	if err != nil {
		return PruneResult{}, err
	}

	result := PruneResult{
		DryRun: options.DryRun,
	}

	worktrees, err := git.WorktreeList(gitDir)
	if err != nil {
		return result, err
	}

	for _, wt := range worktrees {
		if wt.Branch == "" {
			continue
		}
		if wt.Locked {
			result.SkippedWorktrees = append(result.SkippedWorktrees, SkippedWorktree{
				Branch: wt.Branch,
				Path:   wt.Path,
				Reason: lockedReason(wt),
			})
			continue
		}
		prunable := wt.Prunable
		if !prunable {
			if _, err := os.Stat(wt.Path); err == nil {
				continue
			} else if os.IsNotExist(err) {
				prunable = true
			} else {
				result.SkippedWorktrees = append(result.SkippedWorktrees, SkippedWorktree{
					Branch: wt.Branch,
					Path:   wt.Path,
					Reason: err.Error(),
				})
				continue
			}
		}
		if !prunable {
			continue
		}

		if options.DryRun {
			result.PrunedWorktrees = append(result.PrunedWorktrees, PrunedWorktree{
				Branch: wt.Branch,
				Path:   wt.Path,
			})
			continue
		}

		if _, err := git.Run("--git-dir", gitDir, "worktree", "remove", "--force", wt.Path); err != nil {
			return result, err
		}

		result.PrunedWorktrees = append(result.PrunedWorktrees, PrunedWorktree{
			Branch: wt.Branch,
			Path:   wt.Path,
		})
	}

	if !options.MergedBranches {
		return result, nil
	}

	worktrees, err = git.WorktreeList(gitDir)
	if err != nil {
		return result, err
	}

	active := map[string]struct{}{}
	for _, wt := range worktrees {
		if wt.Branch == "" {
			continue
		}
		active[wt.Branch] = struct{}{}
	}

	baseBranch := resolveBaseBranch(gitDir, workDir, false)

	branches, err := git.LocalBranches(gitDir)
	if err != nil {
		return result, err
	}

	for _, branch := range branches {
		if _, ok := active[branch]; ok {
			continue
		}
		if _, ok := protectedBranches[branch]; ok {
			result.SkippedBranches = append(result.SkippedBranches, SkippedBranch{
				Branch: branch,
				Reason: ErrProtectedBranch.Error(),
			})
			continue
		}

		mergeCheck, err := git.IsMerged(gitDir, branch, baseBranch)
		if err != nil {
			result.SkippedBranches = append(result.SkippedBranches, SkippedBranch{
				Branch: branch,
				Reason: err.Error(),
			})
			continue
		}
		if !mergeCheck.Merged {
			result.SkippedBranches = append(result.SkippedBranches, SkippedBranch{
				Branch: branch,
				Reason: ErrBranchNotMerged.Error(),
			})
			continue
		}

		if options.DryRun {
			result.PrunedBranches = append(result.PrunedBranches, PrunedBranch{
				Branch: branch,
			})
			continue
		}

		deleteFlag := "-d"
		if !mergeCheck.Ancestor {
			deleteFlag = "-D"
		}
		if _, err := git.Run("--git-dir", gitDir, "branch", deleteFlag, branch); err != nil {
			result.SkippedBranches = append(result.SkippedBranches, SkippedBranch{
				Branch: branch,
				Reason: err.Error(),
			})
			continue
		}

		result.PrunedBranches = append(result.PrunedBranches, PrunedBranch{
			Branch: branch,
		})
	}

	return result, nil
}

func lockedReason(wt git.Worktree) string {
	if wt.LockedReason != "" {
		return "locked: " + wt.LockedReason
	}
	return "locked"
}
