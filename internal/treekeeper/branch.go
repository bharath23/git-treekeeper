package treekeeper

// CreateBranch is a stub for Stage 1
func CreateBranch(branch, base string) (string, error) {
	if base == "" {
		base = "main"
	}
	return "/tmp/" + branch, nil
}
