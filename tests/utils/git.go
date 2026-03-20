package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func RunGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, string(out))
	}
	return string(out)
}

func RunGitEnv(t *testing.T, dir string, env map[string]string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = os.Environ()
	for key, value := range env {
		cmd.Env = append(cmd.Env, key+"="+value)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, string(out))
	}
	return string(out)
}

func InitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	RunGit(t, dir, "init", "-b", "main")

	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("test"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	RunGit(t, dir, "add", ".")
	RunGit(t, dir, "-c", "user.name=git-tk", "-c", "user.email=git-tk@example.com", "commit", "-m", "init")
	return dir
}

func Chdir(t *testing.T, dir string) func() {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	return func() {
		if err := os.Chdir(wd); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	}
}
