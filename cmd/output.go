package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/bharath23/git-treekeeper/internal/treekeeper"
)

type OutputFormat int

const (
	FormatHuman OutputFormat = iota
	FormatPorcelain
	FormatJSON
	FormatPathOnly
)

type RenderFunc func(out io.Writer, format OutputFormat, response treekeeper.Response) error

var renderers = map[treekeeper.ResponseKind]RenderFunc{
	treekeeper.ResponseBranchCreate: renderBranchCreate,
	treekeeper.ResponseBranchDelete: renderBranchDelete,
	treekeeper.ResponseCheckout:     renderCheckout,
	treekeeper.ResponseClone:        renderClone,
	treekeeper.ResponseList:         renderList,
	treekeeper.ResponseDoctor:       renderDoctor,
	treekeeper.ResponsePrune:        renderPrune,
	treekeeper.ResponseSync:         renderSync,
	treekeeper.ResponseSyncAll:      renderSyncAll,
	treekeeper.ResponseSetup:        renderSetup,
	treekeeper.ResponsePassThrough:  renderPassThrough,
}

func RenderResponse(out io.Writer, format OutputFormat, response treekeeper.Response) error {
	render, ok := renderers[response.Kind]
	if !ok {
		return fmt.Errorf("unknown response kind")
	}
	return render(out, format, response)
}

func renderList(out io.Writer, format OutputFormat, response treekeeper.Response) error {
	rows := make([][]string, 0, len(response.Worktrees))
	for _, wt := range response.Worktrees {
		rows = append(rows, []string{wt.Branch, wt.Path})
	}
	return renderTableOutput(out, format, []string{"branch", "path"}, rows, response.Worktrees)
}

func renderDoctor(out io.Writer, format OutputFormat, response treekeeper.Response) error {
	rows := make([][]string, 0, len(response.Doctor))
	for _, result := range response.Doctor {
		rows = append(rows, []string{result.Branch, result.State})
	}
	return renderTableOutput(out, format, []string{"branch", "state"}, rows, response.Doctor)
}

func renderTableOutput(out io.Writer, format OutputFormat, headers []string, rows [][]string, jsonValue any) error {
	switch format {
	case FormatHuman:
		return renderTable(out, headers, rows)
	case FormatPorcelain:
		return renderPorcelain(out, rows)
	case FormatJSON:
		return renderJSON(out, jsonValue)
	default:
		return fmt.Errorf("unknown output format")
	}
}

func renderTable(out io.Writer, headers []string, rows [][]string) error {
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, strings.Join(headers, "\t"))

	separators := make([]string, len(headers))
	for i, header := range headers {
		separators[i] = strings.Repeat("-", len(header))
	}
	fmt.Fprintln(tw, strings.Join(separators, "\t"))

	for _, row := range rows {
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}
	return tw.Flush()
}

func renderPorcelain(out io.Writer, rows [][]string) error {
	for _, row := range rows {
		fmt.Fprintln(out, strings.Join(row, "\t"))
	}
	return nil
}

func renderJSON(out io.Writer, jsonValue any) error {
	enc := json.NewEncoder(out)
	return enc.Encode(jsonValue)
}

func renderWorktreePath(out io.Writer, format OutputFormat, path string) bool {
	if format != FormatPathOnly {
		return false
	}
	fmt.Fprintln(out, path)
	return true
}

func renderBranchCreate(out io.Writer, format OutputFormat, response treekeeper.Response) error {
	if response.BranchCreate == nil {
		return fmt.Errorf("missing branch create payload")
	}
	result := *response.BranchCreate
	if renderWorktreePath(out, format, result.WorktreePath) {
		return nil
	}
	treekeeper.Info("Creating branch %s from %s", result.Branch, result.Base)
	treekeeper.Info("Worktree path: %s", result.WorktreePath)
	return nil
}

func renderBranchDelete(out io.Writer, format OutputFormat, response treekeeper.Response) error {
	if response.BranchDelete == nil {
		return fmt.Errorf("missing branch delete payload")
	}
	if format == FormatPathOnly {
		return nil
	}
	result := *response.BranchDelete
	if result.WorktreePath != "" {
		treekeeper.Info("Deleted workspace: %s", result.WorktreePath)
	}
	treekeeper.Info("Deleted branch: %s", result.Branch)
	if result.RemoteDeleted {
		treekeeper.Info("Deleted remote branch: %s/%s", result.RemoteName, result.Branch)
	}
	return nil
}

func renderClone(out io.Writer, format OutputFormat, response treekeeper.Response) error {
	if response.Clone == nil {
		return fmt.Errorf("missing clone payload")
	}
	result := *response.Clone
	if renderWorktreePath(out, format, result.WorktreePath) {
		return nil
	}
	treekeeper.Info("Cloning repo %s", result.RepoURL)
	treekeeper.Info("Default branch: %s", result.DefaultBranch)
	treekeeper.Info("Worktree path: %s", result.WorktreePath)
	return nil
}

func renderCheckout(out io.Writer, format OutputFormat, response treekeeper.Response) error {
	if response.Checkout == nil {
		return fmt.Errorf("missing checkout payload")
	}
	result := *response.Checkout
	if renderWorktreePath(out, format, result.WorktreePath) {
		return nil
	}
	treekeeper.Info("Checking out branch %s", result.Branch)
	treekeeper.Info("Worktree path: %s", result.WorktreePath)
	return nil
}

func renderPrune(out io.Writer, format OutputFormat, response treekeeper.Response) error {
	if response.Prune == nil {
		return fmt.Errorf("missing prune payload")
	}
	result := *response.Prune
	for _, wt := range result.PrunedWorktrees {
		if result.DryRun {
			treekeeper.Info("Would prune worktree: %s", wt.Path)
		} else {
			treekeeper.Info("Pruned worktree: %s", wt.Path)
		}
	}
	for _, branch := range result.PrunedBranches {
		if result.DryRun {
			treekeeper.Info("Would prune branch: %s", branch.Branch)
		} else {
			treekeeper.Info("Pruned branch: %s", branch.Branch)
		}
	}
	for _, wt := range result.SkippedWorktrees {
		treekeeper.Verbose("Skipped worktree %s: %s", wt.Path, wt.Reason)
	}
	for _, branch := range result.SkippedBranches {
		treekeeper.Verbose("Skipped branch %s: %s", branch.Branch, branch.Reason)
	}
	return nil
}

func renderSync(out io.Writer, format OutputFormat, response treekeeper.Response) error {
	if response.Sync == nil {
		return fmt.Errorf("missing sync payload")
	}
	return renderSyncResult(out, *response.Sync)
}

func renderSyncAll(out io.Writer, format OutputFormat, response treekeeper.Response) error {
	if response.SyncAll == nil {
		return fmt.Errorf("missing sync all payload")
	}
	result := *response.SyncAll

	for i, res := range result.Results {
		if i > 0 {
			fmt.Fprintln(out)
		}
		if err := renderSyncResult(out, res); err != nil {
			return err
		}
	}

	if len(result.Skipped) > 0 {
		fmt.Fprintln(out)
		treekeeper.Warning("Skipped branches:")
		for _, skip := range result.Skipped {
			treekeeper.Warning("- %s: %s", skip.Branch, skip.Reason)
		}
	}

	return nil
}

func renderSyncResult(out io.Writer, result treekeeper.SyncResult) error {
	treekeeper.Info("Syncing %s from %s/%s", result.Branch, result.Remote, result.RemoteBranch)

	if result.AddedUpstream {
		if result.DryRun {
			treekeeper.Info("Would add upstream remote %s (%s)", result.UpstreamName, result.UpstreamURL)
		} else {
			treekeeper.Info("Added upstream remote %s (%s)", result.UpstreamName, result.UpstreamURL)
		}
	}
	if result.SetUpstream {
		if result.DryRun {
			treekeeper.Info("Would set upstream for %s to %s/%s", result.Branch, result.UpstreamName, result.Branch)
			if result.PushRemote != "" {
				treekeeper.Info("Would set push remote for %s to %s", result.Branch, result.PushRemote)
			}
		} else {
			treekeeper.Info("Set upstream for %s to %s/%s", result.Branch, result.UpstreamName, result.Branch)
			if result.PushRemote != "" {
				treekeeper.Info("Set push remote for %s to %s", result.Branch, result.PushRemote)
			}
		}
	}

	if result.DryRun {
		treekeeper.Info("Would fetch %s", result.Remote)
		treekeeper.Info("Would fast-forward %s from %s/%s", result.Branch, result.Remote, result.RemoteBranch)
		return nil
	}

	for _, line := range result.FetchOutput {
		treekeeper.Info("%s", line)
	}
	for _, line := range result.MergeOutput {
		treekeeper.Info("%s", line)
	}
	return nil
}

func renderSetup(out io.Writer, format OutputFormat, response treekeeper.Response) error {
	if response.Setup == nil {
		return fmt.Errorf("missing setup payload")
	}
	result := *response.Setup

	if result.AddedUpstream {
		if result.DryRun {
			treekeeper.Info("Would add upstream remote %s (%s)", result.UpstreamName, result.UpstreamURL)
		} else {
			treekeeper.Info("Added upstream remote %s (%s)", result.UpstreamName, result.UpstreamURL)
		}
	}

	if result.SetUpstream {
		if result.DryRun {
			treekeeper.Info("Would set upstream for %s to %s/%s", result.Branch, result.UpstreamName, result.Branch)
		} else {
			treekeeper.Info("Set upstream for %s to %s/%s", result.Branch, result.UpstreamName, result.Branch)
		}
	}

	if result.SetPushRemote {
		if result.DryRun {
			treekeeper.Info("Would set push remote for %s to %s", result.Branch, result.OriginName)
		} else {
			treekeeper.Info("Set push remote for %s to %s", result.Branch, result.OriginName)
		}
	}

	if result.HooksInstalled {
		if result.DryRun {
			treekeeper.Info("Would install hooks in %s", result.HooksPath)
		} else {
			treekeeper.Info("Installed hooks in %s", result.HooksPath)
		}
	}

	return nil
}

func renderPassThrough(out io.Writer, format OutputFormat, response treekeeper.Response) error {
	return nil
}
