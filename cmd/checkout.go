package cmd

import (
	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/spf13/cobra"
)

func NewCheckoutCmd() *cobra.Command {
	var pathOnly bool

	cmd := &cobra.Command{
		Use:   "checkout <branch>",
		Short: "Checkout a branch",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return treekeeper.ErrMissingCheckoutBranch
			}
			if len(args) > 1 {
				_ = cmd.Usage()
				return treekeeper.ErrTooManyArgs
			}

			branchName := args[0]
			path, err := treekeeper.Checkout(branchName)
			if err != nil {
				return err
			}
			format := FormatHuman
			if pathOnly {
				format = FormatPathOnly
			}
			return RenderResponse(cmd.OutOrStdout(), format, treekeeper.Response{
				Kind: treekeeper.ResponseCheckout,
				Checkout: &treekeeper.CheckoutOutput{
					Branch:       branchName,
					WorktreePath: path,
				},
			})
		},
	}
	cmd.Flags().BoolVar(&pathOnly, "path-only", false, "Print only the worktree path")
	return cmd
}
