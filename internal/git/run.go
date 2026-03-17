package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Run(args ...string) (string, error) {
	return runCmd(exec.Command("git", args...))
}

func RunInDir(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	return runCmd(cmd)
}

func RawRun(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCmd(cmd *exec.Cmd) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	out := strings.TrimSpace(stdout.String())
	errOut := strings.TrimSpace(stderr.String())
	if err != nil {
		if errOut == "" {
			errOut = out
		}
		if errOut != "" {
			return "", fmt.Errorf("git %s: %w: %s", strings.Join(cmd.Args[1:], " "), err, errOut)
		}
		return "", fmt.Errorf("git %s: %w", strings.Join(cmd.Args[1:], " "), err)
	}

	return out, nil
}
