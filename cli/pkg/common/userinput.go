package common

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/user_input.go . IUserInput
type IUserInput interface {
	AskBool(question string, opts *UserInputOptions) bool
}

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
	fmt.Fprintf(ui.Writer, "%s (y/n)\n", question)

	if opts.AssumeYes {
		return true
	}

	line, err := ui.read()
	if err != nil {
		log.Fatal(err.Error())
	}

	if input := strings.TrimSpace(strings.ToLower(line)); input == "y" || input == "yes" {
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
