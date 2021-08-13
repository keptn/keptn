package lib

import (
	"errors"
	"os/exec"
	"strings"
)

//go:generate moq  -pkg fake -out ./fake/curl_executor_mock.go . ICurlExecutor
type ICurlExecutor interface {
	Curl(curlCmd string) (string, error)
}

type CmdCurlExecutor struct{}

func (ce *CmdCurlExecutor) Curl(curlCmd string) (string, error) {
	cmdArr := strings.Split(curlCmd, " ")
	if len(cmdArr) == 0 {
		return "", errors.New("no command provided")
	}
	if cmdArr[0] != "curl" {
		return "", errors.New("only curl commands are allowed to be executed")
	}
	cmd := exec.Command(cmdArr[0], cmdArr[1:]...)

	output, err := cmd.Output()
	return string(output), err
}
