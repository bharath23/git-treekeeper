package cmd

import (
	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/spf13/cobra"
)

func NewCheckoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "checkout <branch>",
		Short: "Checkout a branch",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return treekeeper.ErrMissingCheckoutBranch
			}

			branch := args[0]
			path := "/tmp/" + branch // stub worktree path
			treekeeper.Info("Checking out branch %s", branch)
			treekeeper.Info("Worktree path: %s", path)
			return nil
		},
	}
}
