package treekeeper

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/bharath23/git-treekeeper/internal/git"
)

type DoctorInfo struct {
	Branch string
	State  string
}

func Doctor() ([]DoctorInfo, error) {
	worktrees, err := ListWorktrees()
	if err != nil {
		return nil, err
	}

	results := make([]DoctorInfo, 0, len(worktrees))
	for _, wt := range worktrees {
		state, err := worktreeState(wt.Path)
		if err != nil {
			return nil, err
		}
		results = append(results, DoctorInfo{
			Branch: wt.Branch,
			State:  state,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Branch < results[j].Branch
	})

	return results, nil
}

func worktreeState(worktreePath string) (string, error) {
	gitDir, err := git.ResolveGitDir(worktreePath)
	if err != nil {
		return "", err
	}

	if exists(filepath.Join(gitDir, "rebase-apply")) || exists(filepath.Join(gitDir, "rebase-merge")) {
		return "rebase in progress", nil
	}
	if exists(filepath.Join(gitDir, "MERGE_HEAD")) {
		return "merge in progress", nil
	}

	status, err := git.StatusPorcelain(worktreePath)
	if err != nil {
		return "", err
	}
	if status != "" {
		return "dirty", nil
	}

	return "clean", nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
