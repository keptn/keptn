package lib

import (
	"errors"
	"os/exec"
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

func IsNoCommandError(err error) bool {
	if curlErr, ok := err.(*CurlError); ok {
		return curlErr.reason == NoCommandError
	}
	return false
}

func IsInvalidCommandError(err error) bool {
	if curlErr, ok := err.(*CurlError); ok {
		return curlErr.reason == InvalidCommandError
	}
	return false
}

func IsUnallowedURLError(err error) bool {
	if curlErr, ok := err.(*CurlError); ok {
		return curlErr.reason == UnallowedURLError
	}
	return false
}

//go:generate moq  -pkg fake -out ./fake/curl_executor_mock.go . ICurlExecutor
type ICurlExecutor interface {
	Curl(curlCmd string) (string, error)
}

type CmdCurlExecutor struct {
	unAllowedURLs []string
}

type CmdCurlExecutorOption func(executor *CmdCurlExecutor)

func WithUnAllowedURLs(urls []string) CmdCurlExecutorOption {
	return func(executor *CmdCurlExecutor) {
		executor.unAllowedURLs = urls
	}
}

func NewCmdCurlExecutor(opts ...CmdCurlExecutorOption) *CmdCurlExecutor {
	executor := &CmdCurlExecutor{}
	for _, o := range opts {
		o(executor)
	}
	return executor
}

func (ce *CmdCurlExecutor) Curl(curlCmd string) (string, error) {
	cmdArr := strings.Split(curlCmd, " ")
	if len(cmdArr) == 0 {
		return "", &CurlError{err: errors.New("no command provided"), reason: NoCommandError}
	}
	if cmdArr[0] != "curl" {
		return "", &CurlError{err: errors.New("only curl commands are allowed to be executed"), reason: InvalidCommandError}
	}

	// check if the curl command contains any of the disallowed URLs
	for _, url := range ce.unAllowedURLs {
		if strings.Contains(curlCmd, url) {
			return "", &CurlError{err: errors.New("only curl commands are allowed to be executed"), reason: UnallowedURLError}
		}
	}

	cmd := exec.Command("/bin/sh", "-c", curlCmd)

	output, err := cmd.Output()
	return string(output), err
}
