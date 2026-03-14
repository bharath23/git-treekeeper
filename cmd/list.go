package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List worktrees",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				_ = cmd.Usage()
				return treekeeper.ErrTooManyArgs
			}

			worktrees, err := treekeeper.ListWorktrees()
			if err != nil {
				return err
			}

			tw := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(tw, "branch\tpath")
			fmt.Fprintln(tw, "------\t----")
			for _, wt := range worktrees {
				fmt.Fprintf(tw, "%s\t%s\n", wt.Branch, wt.Path)
			}
			return tw.Flush()
		},
	}
}
