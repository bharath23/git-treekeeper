package cmd

import (
	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/spf13/cobra"
)

func NewCloneCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clone <repo-url>",
		Short: "Clone a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return treekeeper.ErrMissingRepoURL
			}

			repo := args[0]
			defaultBranch, worktreePath, err := treekeeper.Clone(repo)
			if err != nil {
				return err
			}
			treekeeper.Info("Cloning repo %s", repo)
			treekeeper.Info("Default branch: %s", defaultBranch)
			treekeeper.Info("Worktree path: %s", worktreePath)
			return nil
		},
	}
}
