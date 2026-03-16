package cmd

import (
	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/spf13/cobra"
)

func NewSetupCmd() *cobra.Command {
	var branch string
	var upstream string
	var origin string
	var upstreamURL string
	var installHooks bool
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Configure upstream and workflow defaults",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				_ = cmd.Usage()
				return treekeeper.ErrTooManyArgs
			}

			result, err := treekeeper.Setup(treekeeper.SetupOptions{
				Branch:       branch,
				Upstream:     upstream,
				Origin:       origin,
				UpstreamURL:  upstreamURL,
				InstallHooks: installHooks,
				DryRun:       dryRun,
			})
			if err != nil {
				return err
			}

			return RenderResponse(cmd.OutOrStdout(), FormatHuman, treekeeper.Response{
				Kind:  treekeeper.ResponseSetup,
				Setup: &result,
			})
		},
	}

	cmd.Flags().StringVar(&branch, "branch", "", "Branch to configure (defaults to the default branch)")
	cmd.Flags().StringVar(&upstream, "upstream", "upstream", "Upstream remote name")
	cmd.Flags().StringVar(&origin, "origin", "origin", "Origin remote name for push defaults")
	cmd.Flags().StringVar(&upstreamURL, "upstream-url", "", "Add upstream remote with the given URL")
	cmd.Flags().BoolVar(&installHooks, "install-hooks", false, "Install hooks that block commits on the branch")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be configured without making changes")
	return cmd
}
