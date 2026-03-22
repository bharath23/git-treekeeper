package treekeeper

import "github.com/bharath23/git-treekeeper/internal/git"

const defaultOriginRemote = "origin"

func fetchRefspec(remote string) string {
	return "+refs/heads/*:refs/remotes/" + remote + "/*"
}

func ensureRemoteFetchRefspec(gitDir, remote string, dryRun bool) (bool, error) {
	refspecs, err := git.RemoteFetchRefspecs(gitDir, remote)
	if err != nil {
		return false, err
	}
	expected := fetchRefspec(remote)
	for _, refspec := range refspecs {
		if refspec == expected {
			return false, nil
		}
	}
	if dryRun {
		return true, nil
	}
	if err := git.AddRemoteFetchRefspec(gitDir, remote, expected); err != nil {
		return false, err
	}
	return true, nil
}

func ensureOriginFetchRefspec(gitDir string, dryRun bool) (bool, error) {
	return ensureRemoteFetchRefspec(gitDir, defaultOriginRemote, dryRun)
}

func ensureOriginUpstream(gitDir, branchName string) error {
	upstream, err := git.BranchUpstream(gitDir, branchName)
	if err != nil {
		return err
	}
	if upstream != "" {
		return nil
	}

	originExists, err := git.RemoteExists(gitDir, defaultOriginRemote)
	if err != nil || !originExists {
		return err
	}

	remoteRef := "refs/remotes/" + defaultOriginRemote + "/" + branchName
	remoteRefExists, err := git.RefExists(gitDir, remoteRef)
	if err != nil {
		return err
	}
	if !remoteRefExists {
		return nil
	}

	return git.SetBranchUpstream(gitDir, branchName, defaultOriginRemote+"/"+branchName)
}
