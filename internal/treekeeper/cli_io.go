package treekeeper

import (
	"io"
	"os"
)

var (
	Stdout io.Writer = os.Stdout
	Stderr io.Writer = os.Stderr
)

func SetOutput(out, err io.Writer) {
	Stdout = out
	Stderr = err
}

func ResetOutput() {
	Stdout = os.Stdout
	Stderr = os.Stderr
}
