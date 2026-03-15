package treekeeper

import (
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

	if options.DryRun {
		return result, nil
	}

	fetchOutput, err := git.Run("--git-dir", gitDir, "fetch", remoteName)
	if err != nil {
		return result, err
	}
	result.FetchOutput = splitLines(fetchOutput)

	mergeOutput, err := git.RunInDir(worktreePath, "merge", "--ff-only", remoteName+"/"+remoteBranch)
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
