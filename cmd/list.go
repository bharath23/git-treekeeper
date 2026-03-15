package cmd

import (
	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	var porcelain bool
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List worktrees",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				_ = cmd.Usage()
				return treekeeper.ErrTooManyArgs
			}
			if porcelain && jsonOutput {
				return treekeeper.ErrOutputFormatConflict
			}

			worktrees, err := treekeeper.ListWorktrees()
			if err != nil {
				return err
			}

			format := FormatHuman
			if porcelain {
				format = FormatPorcelain
			} else if jsonOutput {
				format = FormatJSON
			}

			return RenderResponse(cmd.OutOrStdout(), format, treekeeper.Response{
				Kind:      treekeeper.ResponseList,
				Worktrees: worktrees,
			})
		},
	}

	cmd.Flags().BoolVar(&porcelain, "porcelain", false, "Machine-readable output")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "JSON output")
	return cmd
}
