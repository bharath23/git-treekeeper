package treekeeper

// Checkout is a stub for Stage 1
func Checkout(branch string) (string, error) {
	return "/tmp/" + branch, nil
}
