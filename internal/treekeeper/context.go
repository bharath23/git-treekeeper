package treekeeper

import "os"

type RepoContext struct {
	WorkDir       string
	GitDir        string
	BaseDir       string
	WorktreesRoot string
}

func ResolveContext(workDir string) (RepoContext, error) {
	gitDir, baseDir, err := resolveGitDir(workDir)
	if err != nil {
		return RepoContext{}, err
	}
	return RepoContext{
		WorkDir:       workDir,
		GitDir:        gitDir,
		BaseDir:       baseDir,
		WorktreesRoot: worktreeRoot(baseDir),
	}, nil
}

func EnsureWorktreesRoot(ctx RepoContext) error {
	return os.MkdirAll(ctx.WorktreesRoot, 0o755)
}
