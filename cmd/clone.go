package cmd

import (
	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/spf13/cobra"
)

func NewCloneCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clone <repo-url> [path]",
		Short: "Clone a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return treekeeper.ErrMissingRepoURL
			}
			if len(args) > 2 {
				_ = cmd.Usage()
				return treekeeper.ErrTooManyArgs
			}

			repoURL := args[0]
			destPath := ""
			if len(args) > 1 {
				destPath = args[1]
			}
			defaultBranch, worktreePath, err := treekeeper.Clone(repoURL, destPath)
			if err != nil {
				return err
			}
			treekeeper.Info("Cloning repo %s", repoURL)
			treekeeper.Info("Default branch: %s", defaultBranch)
			treekeeper.Info("Worktree path: %s", worktreePath)
			return nil
		},
	}
}
