package cmd

import (
	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/spf13/cobra"
)

func NewBranchCmd() *cobra.Command {
	var deleteBranch bool
	var deleteRemote bool
	var forceDelete bool

	cmd := &cobra.Command{
		Use:   "branch <name> [base]",
		Short: "Create or delete a branch",
		RunE: func(cmd *cobra.Command, args []string) error {
			if forceDelete {
				deleteBranch = true
			}
			if len(args) < 1 {
				return treekeeper.ErrMissingBranchName
			}
			if len(args) > 2 && !deleteBranch {
				_ = cmd.Usage()
				return treekeeper.ErrTooManyArgs
			}
			if len(args) > 1 && deleteBranch {
				_ = cmd.Usage()
				return treekeeper.ErrTooManyArgs
			}

			branchName := args[0]
			if deleteBranch {
				result, err := treekeeper.DeleteBranch(branchName, deleteRemote, forceDelete)
				if err != nil {
					return err
				}
				if result.WorktreePath != "" {
					treekeeper.Info("Deleted workspace: %s", result.WorktreePath)
				}
				treekeeper.Info("Deleted branch: %s", branchName)
				if result.RemoteDeleted {
					treekeeper.Info("Deleted remote branch: %s/%s", result.RemoteName, branchName)
				}
				return nil
			}

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
	cmd.Flags().BoolVarP(&deleteBranch, "delete", "d", false, "Delete branch and worktree")
	cmd.Flags().BoolVarP(&forceDelete, "force", "D", false, "Force delete branch even if unmerged")
	cmd.Flags().BoolVar(&deleteRemote, "remote", false, "Delete remote branch")
	return cmd
}
