package common

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type UserInput struct {
	Writer io.Writer
	Reader *bufio.Reader
}

type UserInputOptions struct {
	AssumeYes bool
}

func NewUserInput() *UserInput {
	return &UserInput{
		Writer: os.Stdout,
		Reader: bufio.NewReader(os.Stdin),
	}
}

func (ui *UserInput) AskBool(question string, opts *UserInputOptions) bool {
	fmt.Fprintf(ui.Writer, "%s (y/n)", question)

	if opts.AssumeYes {
		return true
	}

	line, err := ui.read()
	if err != nil {
		log.Fatal(err.Error())
	}

	if line == "y" || line == "Y" || line == "Yes" || line == "yes" {
		return true
	}

	return false

}

func (ui *UserInput) read() (string, error) {
	line, err := ui.Reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	resultStr := strings.TrimSuffix(line, "\n")
	return resultStr, err
}
