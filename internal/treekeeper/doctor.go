package treekeeper

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/bharath23/git-treekeeper/internal/git"
)

type DoctorInfo struct {
	Branch string `json:"branch"`
	State  string `json:"state"`
}

func Doctor() ([]DoctorInfo, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	gitDir, baseDir, err := resolveGitDir(workDir)
	if err != nil {
		return nil, err
	}

	worktrees, err := git.WorktreeList(gitDir)
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
		results = append(results, DoctorInfo{
			Branch: wt.Branch,
			State:  state,
		})
	}

	// Check for orphaned directories in worktrees root
	wtRoot := worktreeRoot(baseDir)
	entries, err := os.ReadDir(wtRoot)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			fullPath := filepath.Join(wtRoot, entry.Name())
			if !trackedPaths[fullPath] {
				results = append(results, DoctorInfo{
					Branch: "(none) - " + entry.Name(),
					State:  "orphaned directory",
				})
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Branch < results[j].Branch
	})

	return results, nil
}

func worktreeState(worktreePath string) (string, error) {
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		return "stale (directory missing)", nil
	}

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

func inProgress(worktreePath string) (bool, error) {
	gitDir, err := git.ResolveGitDir(worktreePath)
	if err != nil {
		return false, err
	}

	if exists(filepath.Join(gitDir, "rebase-apply")) || exists(filepath.Join(gitDir, "rebase-merge")) {
		return true, nil
	}
	if exists(filepath.Join(gitDir, "MERGE_HEAD")) {
		return true, nil
	}

	return false, nil
}

func isDirty(worktreePath string) (bool, error) {
	status, err := git.StatusPorcelain(worktreePath)
	if err != nil {
		return false, err
	}
	return status != "", nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
