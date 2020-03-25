package cmd

import (
	"bytes"
	"io"
	"os"

	"github.com/mattn/go-shellwords"
)

const unexpectedErrMsg = "unexpected error, got '%v'"

func executeActionCommandC(cmd string) (string, error) {
	args, err := shellwords.Parse(cmd)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)

	rootCmd.SetOut(buf)
	rootCmd.SetArgs(args)
	err = rootCmd.Execute()
	return buf.String(), err
}

type redirector struct {
	originalStdOut *os.File
	r              *os.File
	w              *os.File
}

func newRedirector() *redirector {
	return &redirector{}
}

func (r *redirector) redirectStdOut() {
	r.originalStdOut = os.Stdout
	r.r, r.w, _ = os.Pipe()
	os.Stdout = r.w
}

func (r *redirector) revertStdOut() string {

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r.r)
		outC <- buf.String()
	}()

	// back to normal state
	r.w.Close()
	os.Stdout = r.originalStdOut
	return <-outC
}
