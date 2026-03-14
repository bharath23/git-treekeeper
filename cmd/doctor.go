package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/spf13/cobra"
)

func NewDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check worktree health",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				_ = cmd.Usage()
				return treekeeper.ErrTooManyArgs
			}

			results, err := treekeeper.Doctor()
			if err != nil {
				return err
			}

			tw := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(tw, "branch\tstate")
			fmt.Fprintln(tw, "------\t-----")
			for _, result := range results {
				fmt.Fprintf(tw, "%s\t%s\n", result.Branch, result.State)
			}
			return tw.Flush()
		},
	}
}
