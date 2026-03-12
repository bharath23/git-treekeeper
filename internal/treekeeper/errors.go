package treekeeper

import "errors"

var (
	ErrMissingBranchName     = errors.New("branch name required")
	ErrMissingCheckoutBranch = errors.New("checkout branch required")
	ErrMissingRepoURL        = errors.New("repository URL required")
)
