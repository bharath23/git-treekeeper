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
			if len(args) > 2 {
				_ = cmd.Usage()
				return treekeeper.ErrTooManyArgs
			}

			branchName := args[0]
			baseBranch := ""
			if len(args) > 1 {
				baseBranch = args[1]
			}

			result, err := treekeeper.CreateBranch(branchName, baseBranch)
			if err != nil {
				return err
			}
			treekeeper.Info("Creating branch %s from %s", branchName, result.Base)
			treekeeper.Info("Worktree path: %s", result.WorktreePath)
			return nil
		},
	}
}
