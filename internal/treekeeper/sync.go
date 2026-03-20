package treekeeper

import (
	"fmt"
	"os"
	"strings"

	"github.com/bharath23/git-treekeeper/internal/git"
)

type SyncOptions struct {
	Branch         string
	DefaultBranch  bool
	Remote         string
	Upstream       string
	AddUpstreamURL string
	SetUpstream    bool
	Origin         string
	DryRun         bool
}

func Sync(options SyncOptions) (SyncResult, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return SyncResult{}, err
	}

	gitDir, _, err := resolveGitDir(workDir)
	if err != nil {
		return SyncResult{}, err
	}

	result, err := prepareSync(gitDir, workDir, options)
	if err != nil {
		return result, err
	}

	return runSyncFetchMerge(gitDir, result, options, false)
}

func SyncAll(options SyncOptions) (SyncAllResult, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return SyncAllResult{}, err
	}

	gitDir, _, err := resolveGitDir(workDir)
	if err != nil {
		return SyncAllResult{}, err
	}

	worktrees, err := git.WorktreeList(gitDir)
	if err != nil {
		return SyncAllResult{}, err
	}

	result := SyncAllResult{
		Results: make([]SyncResult, 0),
		Skipped: make([]SkippedSync, 0),
	}

	fetchedRemotes := make(map[string]bool)

	for _, wt := range worktrees {
		if wt.Branch == "" {
			continue
		}

		state, err := worktreeState(wt.Path)
		if err != nil || state != "clean" {
			reason := state
			if err != nil {
				reason = err.Error()
			}
			result.Skipped = append(result.Skipped, SkippedSync{
				Branch: wt.Branch,
				Path:   wt.Path,
				Reason: reason,
			})
			continue
		}

		branchOptions := options
		branchOptions.Branch = wt.Branch

		res, err := prepareSync(gitDir, wt.Path, branchOptions)
		if err != nil {
			result.Skipped = append(result.Skipped, SkippedSync{
				Branch: wt.Branch,
				Path:   wt.Path,
				Reason: err.Error(),
			})
			continue
		}

		skipFetch := fetchedRemotes[res.Remote]
		res, err = runSyncFetchMerge(gitDir, res, branchOptions, skipFetch)
		if err != nil {
			result.Skipped = append(result.Skipped, SkippedSync{
				Branch: wt.Branch,
				Path:   wt.Path,
				Reason: err.Error(),
			})
			continue
		}

		result.Results = append(result.Results, res)
		fetchedRemotes[res.Remote] = true
	}

	return result, nil
}

func prepareSync(gitDir, workDir string, options SyncOptions) (SyncResult, error) {
	defaultBranch, err := git.DefaultBranch(gitDir)
	if err != nil || defaultBranch == "" {
		defaultBranch = "main"
	}

	currentBranch := ""
	if _, err := git.TopLevel(workDir); err == nil {
		currentBranch, _ = git.CurrentBranch(workDir)
	}

	targetBranch := options.Branch
	if targetBranch == "" {
		if options.DefaultBranch || currentBranch == "" {
			targetBranch = defaultBranch
		} else {
			targetBranch = currentBranch
		}
	}

	result := SyncResult{
		Branch: targetBranch,
		DryRun: options.DryRun,
	}

	branchExists, err := git.BranchExists(gitDir, targetBranch)
	if err != nil {
		return result, err
	}
	if !branchExists {
		return result, ErrBranchNotFound
	}

	worktreePath, ok := worktreeForBranch(gitDir, targetBranch)
	if !ok {
		return result, ErrBranchNotCheckedOut
	}
	result.WorktreePath = worktreePath

	originName := options.Origin
	if originName == "" {
		originName = "origin"
	}

	upstreamName := options.Upstream
	if upstreamName == "" {
		upstreamName = "upstream"
	}

	if options.AddUpstreamURL != "" {
		exists, err := git.RemoteExists(gitDir, upstreamName)
		if err != nil {
			return result, err
		}
		if exists {
			existingURL, err := git.RemoteURL(gitDir, upstreamName)
			if err != nil {
				return result, err
			}
			if existingURL != options.AddUpstreamURL {
				return result, ErrRemoteURLMismatch
			}
		} else {
			if !options.DryRun {
				if err := git.AddRemote(gitDir, upstreamName, options.AddUpstreamURL); err != nil {
					return result, err
				}
			}
			result.AddedUpstream = true
		}
		result.UpstreamName = upstreamName
		result.UpstreamURL = options.AddUpstreamURL
	}

	if options.SetUpstream {
		exists, err := git.RemoteExists(gitDir, upstreamName)
		if err != nil {
			return result, err
		}
		if !exists {
			if options.DryRun && options.AddUpstreamURL != "" {
				exists = true
			} else {
				return result, ErrRemoteNotFound
			}
		}

		originExists, err := git.RemoteExists(gitDir, originName)
		if err != nil {
			return result, err
		}
		if !originExists {
			return result, ErrOriginRemoteMissing
		}

		upstreamRef := upstreamName + "/" + targetBranch
		if !options.DryRun {
			if _, err := git.Run("--git-dir", gitDir, "fetch", upstreamName); err != nil {
				return result, err
			}
			if err := git.SetBranchUpstream(gitDir, targetBranch, upstreamRef); err != nil {
				return result, err
			}
			if err := git.SetBranchPushRemote(gitDir, targetBranch, originName); err != nil {
				return result, err
			}
		}
		result.SetUpstream = true
		result.UpstreamName = upstreamName
		if result.UpstreamURL == "" {
			if url, err := git.RemoteURL(gitDir, upstreamName); err == nil {
				result.UpstreamURL = url
			}
		}
		result.PushRemote = originName
	}

	upstreamRef, err := git.BranchUpstream(gitDir, targetBranch)
	if err != nil {
		return result, err
	}

	remoteName := options.Remote
	remoteBranch := targetBranch
	if remoteName == "" {
		if upstreamRef != "" {
			remoteName, remoteBranch = splitRemoteRef(upstreamRef, targetBranch)
		} else {
			remoteName = originName
		}
	}

	exists, err := git.RemoteExists(gitDir, remoteName)
	if err != nil {
		return result, err
	}
	if !exists {
		return result, ErrRemoteNotFound
	}

	result.Remote = remoteName
	result.RemoteBranch = remoteBranch

	return result, nil
}

func runSyncFetchMerge(gitDir string, result SyncResult, options SyncOptions, skipFetch bool) (SyncResult, error) {
	if options.DryRun {
		return result, nil
	}

	if !skipFetch {
		fetchOutput, err := git.Run("--git-dir", gitDir, "fetch", result.Remote)
		if err != nil {
			return result, err
		}
		result.FetchOutput = splitLines(fetchOutput)
	}

	remoteRef := fmt.Sprintf("refs/remotes/%s/%s", result.Remote, result.RemoteBranch)
	remoteRefExists, err := git.RefExists(gitDir, remoteRef)
	if err != nil {
		return result, err
	}
	if !remoteRefExists {
		refspec := fmt.Sprintf("+refs/heads/%s:%s", result.RemoteBranch, remoteRef)
		fallbackOutput, err := git.Run("--git-dir", gitDir, "fetch", result.Remote, refspec)
		if err != nil {
			return result, err
		}
		if lines := splitLines(fallbackOutput); len(lines) > 0 {
			result.FetchOutput = append(result.FetchOutput, lines...)
		}
	}

	mergeOutput, err := git.RunInDir(result.WorktreePath, "merge", "--ff-only", result.Remote+"/"+result.RemoteBranch)
	if err != nil {
		return result, err
	}
	result.MergeOutput = splitLines(mergeOutput)

	return result, nil
}

func splitRemoteRef(ref string, fallbackBranch string) (string, string) {
	parts := strings.SplitN(ref, "/", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return ref, fallbackBranch
}

func splitLines(output string) []string {
	if output == "" {
		return nil
	}
	return strings.Split(output, "\n")
}
