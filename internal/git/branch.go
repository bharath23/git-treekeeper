package git

import "strings"

func IsMerged(gitDir, branchName, baseBranch string) (bool, error) {
	out, err := Run("--git-dir", gitDir, "branch", "--merged", baseBranch)
	if err != nil {
		return false, err
	}

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(strings.TrimPrefix(line, "*"))
		if line == branchName {
			return true, nil
		}
	}

	return false, nil
}
