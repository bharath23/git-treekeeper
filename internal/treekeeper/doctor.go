package treekeeper

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bharath23/git-treekeeper/internal/git"
)

type DoctorInfo struct {
	Branch   string `json:"branch"`
	State    string `json:"state"`
	Tracking string `json:"tracking"`
}

func Doctor() ([]DoctorInfo, error) {
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

	results := make([]DoctorInfo, 0)
	trackedPaths := make(map[string]bool)

	for _, wt := range worktrees {
		if wt.Branch == "" {
			continue
		}
		trackedPaths[wt.Path] = true
		state, err := worktreeState(wt.Path)
		if err != nil {
			// If we can't resolve git dir or something, it might be stale
			if os.IsNotExist(err) {
				results = append(results, DoctorInfo{
					Branch: wt.Branch,
					State:  "stale (directory missing)",
				})
				continue
			}
			return nil, err
		}
		tracking, err := worktreeTracking(ctx.GitDir, wt.Branch)
		if err != nil {
			return nil, err
		}
		results = append(results, DoctorInfo{
			Branch:   wt.Branch,
			State:    state,
			Tracking: tracking,
		})
	}

	// Check for orphaned directories in worktrees root
	wtRoot := ctx.WorktreesRoot
	entries, err := os.ReadDir(wtRoot)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			fullPath := filepath.Join(wtRoot, entry.Name())
			if !trackedPaths[fullPath] {
				results = append(results, DoctorInfo{
					Branch:   "(none) - " + entry.Name(),
					State:    "orphaned directory",
					Tracking: "n/a",
				})
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Branch < results[j].Branch
	})

	return results, nil
}

func worktreeTracking(gitDir, branchName string) (string, error) {
	upstream, err := git.BranchUpstream(gitDir, branchName)
	if err != nil {
		return "", err
	}
	if upstream == "" {
		return "none", nil
	}
	if strings.Contains(upstream, "/") {
		return upstream, nil
	}
	return upstream, nil
}

func worktreeState(worktreePath string) (string, error) {
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		return "stale (directory missing)", nil
	}

	dirty, inProgress, rebase, err := worktreeCheck(worktreePath)
	if err != nil {
		return "", err
	}
	if inProgress {
		if rebase {
			return "rebase in progress", nil
		}
		return "merge in progress", nil
	}
	if dirty {
		return "dirty", nil
	}
	return "clean", nil
}

func inProgress(worktreePath string) (bool, error) {
	_, inProgress, _, err := worktreeCheck(worktreePath)
	if err != nil {
		return false, err
	}
	return inProgress, nil
}

func isDirty(worktreePath string) (bool, error) {
	dirty, _, _, err := worktreeCheck(worktreePath)
	if err != nil {
		return false, err
	}
	return dirty, nil
}

func worktreeCheck(worktreePath string) (dirty bool, inProgress bool, rebase bool, err error) {
	gitDir, err := git.ResolveGitDir(worktreePath)
	if err != nil {
		return false, false, false, err
	}

	if exists(filepath.Join(gitDir, "rebase-apply")) || exists(filepath.Join(gitDir, "rebase-merge")) {
		return false, true, true, nil
	}
	if exists(filepath.Join(gitDir, "MERGE_HEAD")) {
		return false, true, false, nil
	}

	status, err := git.StatusPorcelain(worktreePath)
	if err != nil {
		return false, false, false, err
	}
	return status != "", false, false, nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
