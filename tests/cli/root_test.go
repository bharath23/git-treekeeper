package cli

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/bharath23/git-treekeeper/tests/utils"
)

func TestPassThroughCommand(t *testing.T) {
	bin := buildBinary(t)
	srcRepo := utils.InitRepo(t)
	tmp := t.TempDir()

	err := exec.Command("git", "clone", srcRepo, tmp).Run()
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	cmd := exec.Command(bin, "status")
	cmd.Dir = tmp
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git-tk status failed: %v, output: %s", err, string(out))
	}

	if !strings.Contains(string(out), "On branch main") {
		t.Errorf("expected git status output, got: %q", string(out))
	}
}

func TestPassThroughFlag(t *testing.T) {
	bin := buildBinary(t)
	srcRepo := utils.InitRepo(t)
	tmp := t.TempDir()
	err := exec.Command("git", "clone", srcRepo, tmp).Run()
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	// 'branch -v' is not supported by git-tk, should pass through
	cmd := exec.Command(bin, "branch", "-v")
	cmd.Dir = tmp
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git-tk branch -v failed: %v, output: %s", err, string(out))
	}

	if !strings.Contains(string(out), "main") || !strings.Contains(string(out), "init") {
		t.Errorf("expected git branch -v output containing 'main' and 'init', got:\n%s", string(out))
	}
}

func TestPassThroughUnknownFlagOnRoot(t *testing.T) {
	bin := buildBinary(t)
	srcRepo := utils.InitRepo(t)
	tmp := t.TempDir()
	err := exec.Command("git", "clone", srcRepo, tmp).Run()
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	// '--version' is a git flag, should pass through
	cmd := exec.Command(bin, "--version")
	cmd.Dir = tmp
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git-tk --version failed: %v, output: %s", err, string(out))
	}

	if !strings.Contains(string(out), "git version") {
		t.Errorf("expected git version output, got: %q", string(out))
	}
}

func buildBinary(t *testing.T) string {
	t.Helper()

	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "..", "..")

	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "git-tk")
	if runtime.GOOS == "windows" {
		binPath += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", binPath, root)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build git-tk: %v, output: %s", err, string(out))
	}

	return binPath
}
