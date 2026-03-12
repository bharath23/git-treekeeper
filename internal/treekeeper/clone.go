package treekeeper

// Clone is a stub for Stage 1
func Clone(repoURL string) (string, string, error) {
	defaultBranch := "main"
	worktreePath := "/tmp/" + defaultBranch
	return defaultBranch, worktreePath, nil
}
