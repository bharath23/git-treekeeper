package git

func SetHooksPath(gitDir, hooksPath string) error {
	_, err := Run("--git-dir", gitDir, "config", "core.hooksPath", hooksPath)
	return err
}
