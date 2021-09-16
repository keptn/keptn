package lib

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

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
		return "", errors.New("no command provided")
	}
	if cmdArr[0] != "curl" {
		return "", errors.New("only curl commands are allowed to be executed")
	}

	// check if the curl command contains any of the disallowed URLs
	for _, url := range ce.unAllowedURLs {
		if strings.Contains(curlCmd, url) {
			return "", fmt.Errorf("requests to %s are disallowed", url)
		}
	}

	cmd := exec.Command("/bin/sh", "-c", curlCmd)

	output, err := cmd.Output()
	return string(output), err
}
