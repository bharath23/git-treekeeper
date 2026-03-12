package utils

import (
	"bytes"
	"os"
)

// CaptureStdout captures standard output during f()
func CaptureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	os.Stdout = old
	return buf.String()
}

// CaptureStderr captures standard error during f()
func CaptureStderr(f func()) string {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f()

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	os.Stderr = old
	return buf.String()
}
