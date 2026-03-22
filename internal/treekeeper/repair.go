package treekeeper

import (
	"os"

	"github.com/bharath23/git-treekeeper/internal/git"
)

type RepairOptions struct {
	Branch   string
	DryRun   bool
	Tracking bool
	Remote   string
}

type RepairResult struct {
	DryRun            bool
	Remote            string
	AddedFetchRefspec bool
	Fetched           bool
	UpdatedUpstreams  []string
	SkippedBranches   []SkippedBranch
}

func Repair(options RepairOptions) (RepairResult, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return RepairResult{}, err
	}

	ctx, err := ResolveContext(workDir)
	if err != nil {
		return RepairResult{}, err
	}

	result := RepairResult{
		DryRun: options.DryRun,
		Remote: defaultOriginRemote,
	}

	if !options.Tracking {
		return result, nil
	}

	remoteName := options.Remote
	if remoteName == "" {
		remoteName = defaultOriginRemote
	}
	result.Remote = remoteName

	originExists, err := git.RemoteExists(ctx.GitDir, remoteName)
	if err != nil {
		return result, err
	}
	if !originExists {
		return result, ErrOriginRemoteMissing
	}

	addedRefspec, err := ensureRemoteFetchRefspec(ctx.GitDir, remoteName, options.DryRun)
	if err != nil {
		return result, err
	}
	result.AddedFetchRefspec = addedRefspec

	if !options.DryRun {
		if _, err := git.Run("--git-dir", ctx.GitDir, "fetch", remoteName); err != nil {
			return result, err
		}
		result.Fetched = true
	}

	branches := []string{}
	if options.Branch != "" {
		branches = []string{options.Branch}
	} else {
		worktrees, err := git.WorktreeList(ctx.GitDir)
		if err != nil {
			return result, err
		}
		seen := make(map[string]bool)
		for _, wt := range worktrees {
			if wt.Branch == "" || seen[wt.Branch] {
				continue
			}
			branches = append(branches, wt.Branch)
			seen[wt.Branch] = true
		}
	}

	for _, branchName := range branches {
		exists, err := git.BranchExists(ctx.GitDir, branchName)
		if err != nil {
			return result, err
		}
		if !exists {
			if options.Branch != "" {
				return result, ErrBranchNotFound
			}
			result.SkippedBranches = append(result.SkippedBranches, SkippedBranch{
				Branch: branchName,
				Reason: ErrBranchNotFound.Error(),
			})
			continue
		}

		upstream, err := git.BranchUpstream(ctx.GitDir, branchName)
		if err != nil {
			return result, err
		}
		if upstream != "" {
			result.SkippedBranches = append(result.SkippedBranches, SkippedBranch{
				Branch: branchName,
				Reason: "upstream already set",
			})
			continue
		}

		remoteRef := "refs/remotes/" + remoteName + "/" + branchName
		remoteRefExists, err := git.RefExists(ctx.GitDir, remoteRef)
		if err != nil {
			return result, err
		}
		if !remoteRefExists {
			if options.DryRun {
				exists, err := git.RemoteBranchExists(ctx.GitDir, remoteName, branchName)
				if err != nil {
					return result, err
				}
				if exists {
					remoteRefExists = true
				}
			}
		}

		if !remoteRefExists {
			result.SkippedBranches = append(result.SkippedBranches, SkippedBranch{
				Branch: branchName,
				Reason: "origin branch not found",
			})
			continue
		}

		if !options.DryRun {
			if err := git.SetBranchUpstream(ctx.GitDir, branchName, remoteName+"/"+branchName); err != nil {
				return result, err
			}
		}
		result.UpdatedUpstreams = append(result.UpdatedUpstreams, branchName)
	}

	return result, nil
}
