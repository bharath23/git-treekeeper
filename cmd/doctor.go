package cmd

import (
	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/spf13/cobra"
)

func NewDoctorCmd() *cobra.Command {
	var porcelain bool
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check worktree health",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				_ = cmd.Usage()
				return treekeeper.ErrTooManyArgs
			}
			if porcelain && jsonOutput {
				return treekeeper.ErrOutputFormatConflict
			}

			results, err := treekeeper.Doctor()
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
				Kind:   treekeeper.ResponseDoctor,
				Doctor: results,
			})
		},
	}

	cmd.Flags().BoolVar(&porcelain, "porcelain", false, "Machine-readable output")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "JSON output")
	return cmd
}
