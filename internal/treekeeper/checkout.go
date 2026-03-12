package treekeeper

// Checkout is a stub for Stage 1
func Checkout(branch string) (string, error) {
	Info("Stub: checking out branch %s", branch)
	return "/tmp/worktrees/" + branch, nil
}
