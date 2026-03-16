package treekeeper

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bharath23/git-treekeeper/internal/git"
)

type SetupOptions struct {
	Branch       string
	Upstream     string
	Origin       string
	UpstreamURL  string
	InstallHooks bool
	DryRun       bool
}

func Setup(options SetupOptions) (SetupResult, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return SetupResult{}, err
	}

	gitDir, baseDir, err := resolveGitDir(workDir)
	if err != nil {
		return SetupResult{}, err
	}

	defaultBranch, err := git.DefaultBranch(gitDir)
	if err != nil || defaultBranch == "" {
		defaultBranch = "main"
	}

	branchName := options.Branch
	if branchName == "" {
		branchName = defaultBranch
	}

	branchExists, err := git.BranchExists(gitDir, branchName)
	if err != nil {
		return SetupResult{}, err
	}
	if !branchExists {
		return SetupResult{}, ErrBranchNotFound
	}

	upstreamName := options.Upstream
	if upstreamName == "" {
		upstreamName = "upstream"
	}
	originName := options.Origin
	if originName == "" {
		originName = "origin"
	}

	result := SetupResult{
		Branch:       branchName,
		UpstreamName: upstreamName,
		OriginName:   originName,
		DryRun:       options.DryRun,
	}

	if options.UpstreamURL != "" {
		exists, err := git.RemoteExists(gitDir, upstreamName)
		if err != nil {
			return result, err
		}
		if exists {
			existingURL, err := git.RemoteURL(gitDir, upstreamName)
			if err != nil {
				return result, err
			}
			if existingURL != options.UpstreamURL {
				return result, ErrRemoteURLMismatch
			}
		} else {
			if !options.DryRun {
				if err := git.AddRemote(gitDir, upstreamName, options.UpstreamURL); err != nil {
					return result, err
				}
			}
			result.AddedUpstream = true
		}
		result.UpstreamURL = options.UpstreamURL
	} else {
		exists, err := git.RemoteExists(gitDir, upstreamName)
		if err != nil {
			return result, err
		}
		if !exists {
			return result, ErrRemoteNotFound
		}
		if url, err := git.RemoteURL(gitDir, upstreamName); err == nil {
			result.UpstreamURL = url
		}
	}

	originExists, err := git.RemoteExists(gitDir, originName)
	if err != nil {
		return result, err
	}
	if !originExists {
		return result, ErrOriginRemoteMissing
	}

	if !options.DryRun {
		if _, err := git.Run("--git-dir", gitDir, "fetch", upstreamName); err != nil {
			return result, err
		}
		if err := git.SetBranchUpstream(gitDir, branchName, upstreamName+"/"+branchName); err != nil {
			return result, err
		}
		if err := git.SetBranchPushRemote(gitDir, branchName, originName); err != nil {
			return result, err
		}
	}

	result.SetUpstream = true
	result.SetPushRemote = true

	if options.InstallHooks {
		hooksPath := filepath.Join(baseDir, ".githooks")
		if !options.DryRun {
			if err := os.MkdirAll(hooksPath, 0o755); err != nil {
				return result, err
			}
			hookFile := filepath.Join(hooksPath, "pre-commit")
			content := hookScript(branchName)
			if err := os.WriteFile(hookFile, []byte(content), 0o755); err != nil {
				return result, err
			}
			if err := git.SetHooksPath(gitDir, hooksPath); err != nil {
				return result, err
			}
		}
		result.HooksInstalled = true
		result.HooksPath = hooksPath
	}

	return result, nil
}

func hookScript(branchName string) string {
	return fmt.Sprintf(`#!/usr/bin/env bash
branch="$(git symbolic-ref --short HEAD 2>/dev/null)"
if [ "$branch" = "%s" ]; then
  echo "Refusing commit on %s. Use a feature branch." >&2
  exit 1
fi
`, branchName, branchName)
}
