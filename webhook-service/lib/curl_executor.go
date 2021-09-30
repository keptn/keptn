package lib

import (
	"errors"
	"fmt"
	"strings"
)

type errType int

const (
	NoCommandError errType = iota
	InvalidCommandError
	UnallowedURLError
	RequestError
)

type CurlError struct {
	err    error
	reason errType
}

func (c *CurlError) Error() string {
	return c.err.Error()
}

func NewCurlError(err error, reason errType) *CurlError {
	return &CurlError{
		err:    err,
		reason: reason,
	}
}

func IsNoCommandError(err error) bool {
	var curlErr *CurlError
	if errors.As(err, &curlErr) {
		return curlErr.reason == NoCommandError
	}
	return false
}

func IsInvalidCommandError(err error) bool {
	var curlErr *CurlError
	if errors.As(err, &curlErr) {
		return curlErr.reason == InvalidCommandError
	}
	return false
}

func IsUnallowedURLError(err error) bool {
	var curlErr *CurlError
	if errors.As(err, &curlErr) {
		return curlErr.reason == UnallowedURLError
	}
	return false
}

func IsRequestError(err error) bool {
	var curlErr *CurlError
	if errors.As(err, &curlErr) {
		return curlErr.reason == RequestError
	}
	return false
}

//go:generate moq  -pkg fake -out ./fake/curl_executor_mock.go . ICurlExecutor
type ICurlExecutor interface {
	Curl(curlCmd string) (string, error)
}

type CmdCurlExecutor struct {
	unAllowedURLs       []string
	unAllowedCharacters []string
	unAllowedOptions    []string
	commandExecutor     ICommandExecutor
}

type CmdCurlExecutorOption func(executor *CmdCurlExecutor)

func WithUnAllowedURLs(urls []string) CmdCurlExecutorOption {
	return func(executor *CmdCurlExecutor) {
		executor.unAllowedURLs = urls
	}
}

func NewCmdCurlExecutor(cmdExecutor ICommandExecutor, opts ...CmdCurlExecutorOption) *CmdCurlExecutor {
	executor := &CmdCurlExecutor{
		unAllowedCharacters: []string{"$", "|", ";", ">", "$(", " &", "&&", "`"},
		unAllowedOptions:    []string{"-o", "--output"},
		commandExecutor:     cmdExecutor,
	}
	for _, o := range opts {
		o(executor)
	}
	return executor
}

func (ce *CmdCurlExecutor) Curl(curlCmd string) (string, error) {
	cmdArr := strings.Split(curlCmd, " ")
	if len(cmdArr) == 0 || len(cmdArr) == 1 && cmdArr[0] == "" {
		return "", &CurlError{err: errors.New("no command provided"), reason: NoCommandError}
	}

	for _, unallowedCharacter := range ce.unAllowedCharacters {
		if strings.Contains(curlCmd, unallowedCharacter) {
			return "", &CurlError{err: fmt.Errorf("curl command contains unallowed character '%s'", unallowedCharacter), reason: InvalidCommandError}
		}
	}

	args, err := parseCommandLine(curlCmd)
	if err != nil {
		return "", &CurlError{err: errors.New("could not parse curl command"), reason: InvalidCommandError}
	}

	if cmdArr[0] != "curl" {
		return "", &CurlError{err: errors.New("only curl commands are allowed to be executed"), reason: InvalidCommandError}
	}

	for _, arg := range args {
		for _, o := range ce.unAllowedOptions {
			if strings.HasPrefix(arg, o) {
				return "", &CurlError{err: fmt.Errorf("curl command contains invalid option '%s'", o), reason: InvalidCommandError}
			}
		}
	}

	// check if the curl command contains any of the disallowed URLs
	for _, url := range ce.unAllowedURLs {
		if strings.Contains(curlCmd, url) {
			return "", &CurlError{err: fmt.Errorf("curl command contains invalid URL %s", url), reason: UnallowedURLError}
		}
	}

	resp, err := ce.commandExecutor.ExecuteCommand("curl", args[1:]...)
	if err != nil {
		return "", &CurlError{err: fmt.Errorf("error during curl request execution"), reason: RequestError}
	}
	return resp, nil
}

func parseCommandLine(command string) ([]string, error) {
	var args []string
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true
	for i := 0; i < len(command); i++ {
		c := command[i]

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, errors.New("unclosed quote in command")
	}

	if current != "" {
		args = append(args, current)
	}

	return args, nil
}
