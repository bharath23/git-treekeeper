package cmd

import (
	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/spf13/cobra"
)

func NewGCCmd() *cobra.Command {
	var apply bool
	var ageDays int

	cmd := &cobra.Command{
		Use:   "gc",
		Short: "Garbage collect merged branches",
		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun := !apply
			result, err := treekeeper.GC(treekeeper.GCOptions{
				DryRun:  dryRun,
				AgeDays: ageDays,
			})
			if err != nil {
				return err
			}
			return RenderResponse(cmd.OutOrStdout(), FormatHuman, treekeeper.Response{
				Kind: treekeeper.ResponseGC,
				GC:   &result,
			})
		},
	}

	cmd.Flags().BoolVar(&apply, "apply", false, "Delete branches (default is dry-run)")
	cmd.Flags().IntVar(&ageDays, "age-days", 30, "Only consider branches older than this many days")
	return cmd
}
