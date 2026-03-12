package tests

import (
	"testing"

	"github.com/bharath23/git-treekeeper/internal/treekeeper"
	"github.com/bharath23/git-treekeeper/tests/utils"
)

func TestCloneInfoMessage(t *testing.T) {
	out := utils.CaptureStdout(func() {
		_ = treekeeper.Clone("https://github.com/foo/bar")
	})
	expected := "Cloning https://github.com/foo/bar..."
	if out != expected+"\n" {
		t.Errorf("expected %q, got %q", expected, out)
	}
}

func TestCreateBranchInfoMessage(t *testing.T) {
	out := utils.CaptureStdout(func() {
		_ = treekeeper.CreateBranch("feature-x", "main")
	})
	expected := "Stub: creating branch feature-x from main"
	if out != expected+"\n" {
		t.Errorf("expected %q, got %q", expected, out)
	}
}

func TestCheckoutInfoMessage(t *testing.T) {
	out := utils.CaptureStdout(func() {
		_, _ = treekeeper.Checkout("feature-x")
	})
	expected := "Stub: checking out branch feature-x"
	if out != expected+"\n" {
		t.Errorf("expected %q, got %q", expected, out)
	}
}
