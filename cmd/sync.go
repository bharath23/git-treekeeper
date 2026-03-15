package cmd

import (
	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/spf13/cobra"
)

func NewSyncCmd() *cobra.Command {
	var branch string
	var useDefault bool
	var remote string
	var upstream string
	var addUpstream string
	var setUpstream bool
	var origin string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync a branch with its remote",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				_ = cmd.Usage()
				return treekeeper.ErrTooManyArgs
			}
			if useDefault && branch != "" {
				return treekeeper.ErrSyncBranchConflict
			}

			result, err := treekeeper.Sync(treekeeper.SyncOptions{
				Branch:         branch,
				DefaultBranch:  useDefault,
				Remote:         remote,
				Upstream:       upstream,
				AddUpstreamURL: addUpstream,
				SetUpstream:    setUpstream,
				Origin:         origin,
				DryRun:         dryRun,
			})
			if err != nil {
				return err
			}

			return RenderResponse(cmd.OutOrStdout(), FormatHuman, treekeeper.Response{
				Kind: treekeeper.ResponseSync,
				Sync: &result,
			})
		},
	}

	cmd.Flags().BoolVar(&useDefault, "default", false, "Sync the default branch")
	cmd.Flags().StringVar(&branch, "branch", "", "Sync the specified branch")
	cmd.Flags().StringVar(&remote, "remote", "", "Remote to fetch from")
	cmd.Flags().StringVar(&upstream, "upstream", "upstream", "Upstream remote name")
	cmd.Flags().StringVar(&addUpstream, "add-upstream", "", "Add upstream remote with the given URL")
	cmd.Flags().BoolVar(&setUpstream, "set-upstream", false, "Set branch upstream to the upstream remote")
	cmd.Flags().StringVar(&origin, "origin", "origin", "Origin remote name for push defaults")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be synced without making changes")
	return cmd
}
