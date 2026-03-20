package treekeeper

import (
	"fmt"
	"os"
	"time"

	"github.com/bharath23/git-treekeeper/internal/git"
)

type GCOptions struct {
	DryRun  bool
	AgeDays int
}

func GC(options GCOptions) (GCResult, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return GCResult{}, err
	}

	ctx, err := ResolveContext(workDir)
	if err != nil {
		return GCResult{}, err
	}

	ageDays := options.AgeDays
	if ageDays <= 0 {
		ageDays = 30
	}
	cutoff := time.Now().AddDate(0, 0, -ageDays)

	result := GCResult{
		DryRun:          options.DryRun,
		AgeDays:         ageDays,
		PrunedBranches:  make([]PrunedBranch, 0),
		SkippedBranches: make([]SkippedBranch, 0),
	}

	worktrees, err := git.WorktreeList(ctx.GitDir)
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

	baseBranch := resolveBaseBranch(ctx.GitDir, workDir, false)

	branches, err := git.LocalBranches(ctx.GitDir)
	if err != nil {
		return result, err
	}

	for _, branch := range branches {
		if _, ok := protectedBranches[branch]; ok {
			result.SkippedBranches = append(result.SkippedBranches, SkippedBranch{
				Branch: branch,
				Reason: ErrProtectedBranch.Error(),
			})
			continue
		}
		if _, ok := active[branch]; ok {
			result.SkippedBranches = append(result.SkippedBranches, SkippedBranch{
				Branch: branch,
				Reason: "active worktree",
			})
			continue
		}

		commitTime, err := git.BranchCommitTime(ctx.GitDir, branch)
		if err != nil {
			result.SkippedBranches = append(result.SkippedBranches, SkippedBranch{
				Branch: branch,
				Reason: err.Error(),
			})
			continue
		}
		if commitTime.After(cutoff) {
			result.SkippedBranches = append(result.SkippedBranches, SkippedBranch{
				Branch: branch,
				Reason: fmt.Sprintf("recent (<%dd)", ageDays),
			})
			continue
		}

		mergeCheck, err := git.IsMerged(ctx.GitDir, branch, baseBranch)
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
		if _, err := git.Run("--git-dir", ctx.GitDir, "branch", deleteFlag, branch); err != nil {
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
