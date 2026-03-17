package treekeeper

import (
	"os"
	"strings"

	"github.com/bharath23/git-treekeeper/internal/git"
	"github.com/spf13/cobra"
)

func CompleteBranches(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	gitDir, _, err := resolveGitDir(workDir)
	if err != nil {
		// Not a git repo or something wrong, fallback to file completion might be confusing
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	branches, err := git.LocalBranches(gitDir)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var matches []string
	for _, branch := range branches {
		if strings.HasPrefix(branch, toComplete) {
			matches = append(matches, branch)
		}
	}

	return matches, cobra.ShellCompDirectiveNoFileComp
}

func CompleteRemotes(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	gitDir, _, err := resolveGitDir(workDir)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	remotes, err := git.Remotes(gitDir)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var matches []string
	for _, remote := range remotes {
		if strings.HasPrefix(remote, toComplete) {
			matches = append(matches, remote)
		}
	}

	return matches, cobra.ShellCompDirectiveNoFileComp
}
