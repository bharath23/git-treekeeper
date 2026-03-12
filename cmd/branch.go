package cmd

import (
	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/spf13/cobra"
)

func NewBranchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "branch <name> [base]",
		Short: "Create a new branch",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return treekeeper.ErrMissingBranchName
			}

			branch := args[0]
			base := "main"
			if len(args) > 1 {
				base = args[1]
			}

			treekeeper.Info("Creating branch %s from %s", branch, base)
			// TODO: implement real git worktree creation
			return nil
		},
	}
}
