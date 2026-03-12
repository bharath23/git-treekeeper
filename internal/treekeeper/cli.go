package treekeeper

import (
	"fmt"
	"os"
)

func Fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "fatal: "+format+"\n", args...)
	os.Exit(1)
}

func Error(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
}

func Warning(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "warning: "+format+"\n", args...)
}

func Hint(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "hint: "+format+"\n", args...)
}

func Info(format string, args ...any) {
	fmt.Fprintf(os.Stdout, format+"\n", args...)
}
