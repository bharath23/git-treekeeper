package treekeeper

import "errors"

var (
	ErrMissingBranchName     = errors.New("branch name required")
	ErrMissingCheckoutBranch = errors.New("checkout branch required")
	ErrMissingRepoURL        = errors.New("repository URL required")
	ErrTooManyArgs           = errors.New("too many arguments")
	ErrProtectedBranch       = errors.New("refusing to delete protected branch")
	ErrDirtyWorktree         = errors.New("refusing to delete dirty worktree")
	ErrInProgress            = errors.New("refusing to delete worktree with in-progress operation")
	ErrBranchNotMerged       = errors.New("refusing to delete unmerged branch")
	ErrBranchNotCheckedOut   = errors.New("branch is not checked out as a worktree")
	ErrBranchCheckedOut      = errors.New("refusing to delete branch checked out in current worktree")
	ErrDeleteAborted         = errors.New("delete aborted")
	ErrBranchNotFound        = errors.New("branch not found")
	ErrRemoteBranchNotFound  = errors.New("remote branch not found")
)
