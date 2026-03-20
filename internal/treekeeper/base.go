package treekeeper

import "github.com/bharath23/git-treekeeper/internal/git"

func resolveBaseBranch(gitDir, workDir string, preferCurrent bool) string {
	if preferCurrent {
		if branch, err := git.CurrentBranch(workDir); err == nil && branch != "" {
			return branch
		}
	}
	if branch, err := git.DefaultBranch(gitDir); err == nil && branch != "" {
		return branch
	}
	return "main"
}
