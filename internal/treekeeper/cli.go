package treekeeper

import (
	"fmt"
	"os"
)

var outputQuiet bool
var outputVerbose bool

func SetOutputMode(quiet, verbose bool) {
	outputQuiet = quiet
	outputVerbose = verbose
}

func Fatal(format string, args ...any) {
	fmt.Fprintf(Stderr, "fatal: "+format+"\n", args...)
	os.Exit(1)
}

func Error(format string, args ...any) {
	fmt.Fprintf(Stderr, "error: "+format+"\n", args...)
}

func Warning(format string, args ...any) {
	fmt.Fprintf(Stderr, "warning: "+format+"\n", args...)
}

func Hint(format string, args ...any) {
	fmt.Fprintf(Stderr, "hint: "+format+"\n", args...)
}

func Info(format string, args ...any) {
	if outputQuiet {
		return
	}
	fmt.Fprintf(Stdout, format+"\n", args...)
}

func Verbose(format string, args ...any) {
	if !outputVerbose {
		return
	}
	fmt.Fprintf(Stdout, format+"\n", args...)
}
