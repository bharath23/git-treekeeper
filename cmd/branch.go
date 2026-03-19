package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/spf13/cobra"
)

func NewBranchCmd() *cobra.Command {
	var deleteBranch bool
	var deleteRemote bool
	var forceDelete bool
	var assumeYes bool
	var pathOnly bool

	cmd := &cobra.Command{
		Use:   "branch <name> [base]",
		Short: "Create or delete a branch",
		RunE: func(cmd *cobra.Command, args []string) error {
			if forceDelete {
				deleteBranch = true
			}
			if len(args) < 1 {
				return treekeeper.ErrMissingBranchName
			}
			if len(args) > 2 && !deleteBranch {
				_ = cmd.Usage()
				return treekeeper.ErrTooManyArgs
			}
			if len(args) > 1 && deleteBranch {
				_ = cmd.Usage()
				return treekeeper.ErrTooManyArgs
			}

			branchName := args[0]
			format := FormatHuman
			if pathOnly {
				format = FormatPathOnly
			}

			if deleteBranch {
				if err := confirmDelete(branchName, forceDelete, deleteRemote, assumeYes); err != nil {
					return err
				}
				result, err := treekeeper.DeleteBranch(branchName, deleteRemote, forceDelete)
				if err != nil {
					return err
				}
				return RenderResponse(cmd.OutOrStdout(), format, treekeeper.Response{
					Kind: treekeeper.ResponseBranchDelete,
					BranchDelete: &treekeeper.BranchDeleteOutput{
						Branch:        branchName,
						WorktreePath:  result.WorktreePath,
						RemoteDeleted: result.RemoteDeleted,
						RemoteName:    result.RemoteName,
					},
				})
			}

			baseBranch := ""
			if len(args) > 1 {
				baseBranch = args[1]
			}

			result, err := treekeeper.CreateBranch(branchName, baseBranch)
			if err != nil {
				return err
			}
			return RenderResponse(cmd.OutOrStdout(), format, treekeeper.Response{
				Kind: treekeeper.ResponseBranchCreate,
				BranchCreate: &treekeeper.BranchCreateOutput{
					Branch:       branchName,
					Base:         result.Base,
					WorktreePath: result.WorktreePath,
				},
			})
		},
		ValidArgsFunction: treekeeper.CompleteBranches,
	}
	cmd.Flags().BoolVarP(&deleteBranch, "delete", "d", false, "Delete branch and worktree")
	cmd.Flags().BoolVarP(&forceDelete, "force", "D", false, "Force delete branch even if unmerged")
	cmd.Flags().BoolVar(&deleteRemote, "remote", false, "Delete remote branch")
	cmd.Flags().BoolVar(&assumeYes, "yes", false, "Skip delete confirmation")
	cmd.Flags().BoolVar(&pathOnly, "path-only", false, "Print only the worktree path")
	return cmd
}

func confirmDelete(branchName string, forceDelete bool, deleteRemote bool, assumeYes bool) error {
	if assumeYes || (!forceDelete && !deleteRemote) {
		return nil
	}

	flags := []string{}
	if forceDelete {
		flags = append(flags, "force")
	}
	if deleteRemote {
		flags = append(flags, "remote")
	}

	note := ""
	if len(flags) > 0 {
		note = " (" + strings.Join(flags, ", ") + ")"
	}
	fmt.Fprintf(os.Stderr, "Confirm delete of branch %s%s [y/N]: ", branchName, note)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	resp := strings.ToLower(strings.TrimSpace(line))
	if resp == "y" || resp == "yes" {
		return nil
	}
	return treekeeper.ErrDeleteAborted
}
