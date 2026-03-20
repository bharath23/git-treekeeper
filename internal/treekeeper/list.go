package treekeeper

import (
	"os"
	"sort"

	"github.com/bharath23/git-treekeeper/internal/git"
)

type WorktreeInfo struct {
	Branch string `json:"branch"`
	Path   string `json:"path"`
}

func ListWorktrees() ([]WorktreeInfo, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	ctx, err := ResolveContext(workDir)
	if err != nil {
		return nil, err
	}

	worktrees, err := git.WorktreeList(ctx.GitDir)
	if err != nil {
		return nil, err
	}

	var results []WorktreeInfo
	for _, wt := range worktrees {
		if wt.Branch == "" {
			continue
		}
		results = append(results, WorktreeInfo{
			Branch: wt.Branch,
			Path:   wt.Path,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Branch < results[j].Branch
	})

	return results, nil
}
