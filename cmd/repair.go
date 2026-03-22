package cmd

import (
	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/spf13/cobra"
)

func NewRepairCmd() *cobra.Command {
	var branch string
	dryRun := true
	var tracking bool
	var apply bool
	var remote string

	cmd := &cobra.Command{
		Use:   "repair",
		Short: "Repair origin tracking for worktrees",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				_ = cmd.Usage()
				return treekeeper.ErrTooManyArgs
			}

			if apply {
				dryRun = false
			} else {
				dryRun = true
			}

			result, err := treekeeper.Repair(treekeeper.RepairOptions{
				Branch:   branch,
				DryRun:   dryRun,
				Tracking: tracking,
				Remote:   remote,
			})
			if err != nil {
				return err
			}

			return RenderResponse(cmd.OutOrStdout(), FormatHuman, treekeeper.Response{
				Kind:   treekeeper.ResponseRepair,
				Repair: &result,
			})
		},
	}

	cmd.Flags().StringVar(&branch, "branch", "", "Repair only the specified branch")
	cmd.Flags().BoolVar(&dryRun, "dry-run", true, "Show what would change without making updates")
	cmd.Flags().BoolVar(&tracking, "tracking", true, "Fix origin tracking (default)")
	cmd.Flags().BoolVar(&apply, "apply", false, "Apply changes (disables dry-run)")
	cmd.Flags().StringVar(&remote, "remote", "origin", "Remote to use for tracking fixes")

	_ = cmd.RegisterFlagCompletionFunc("branch", treekeeper.CompleteBranches)
	_ = cmd.RegisterFlagCompletionFunc("remote", treekeeper.CompleteRemotes)

	return cmd
}
