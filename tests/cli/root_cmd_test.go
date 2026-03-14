package cli

import (
	"github.com/bharath23/git-treekeeper/cmd"
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	return cmd.NewRootCmd()
}
