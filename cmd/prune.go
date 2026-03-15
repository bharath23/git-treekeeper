package cmd

import (
	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/spf13/cobra"
)

func NewPruneCmd() *cobra.Command {
	var dryRun bool
	var mergedBranches bool

	cmd := &cobra.Command{
		Use:   "prune",
		Short: "Prune stale worktrees",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				_ = cmd.Usage()
				return treekeeper.ErrTooManyArgs
			}

			result, err := treekeeper.Prune(treekeeper.PruneOptions{
				DryRun:         dryRun,
				MergedBranches: mergedBranches,
			})
			if err != nil {
				return err
			}

			return RenderResponse(cmd.OutOrStdout(), FormatHuman, treekeeper.Response{
				Kind:  treekeeper.ResponsePrune,
				Prune: &result,
			})
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be pruned without making changes")
	cmd.Flags().BoolVar(&mergedBranches, "merged-branches", false, "Also delete merged branches without worktrees")
	return cmd
}
