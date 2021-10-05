package lib

import "os/exec"

//go:generate moq  -pkg fake -out ./fake/cmd_executor_mock.go . ICommandExecutor
type ICommandExecutor interface {
	ExecuteCommand(cmd string, args ...string) (string, error)
}

type OSCmdExecutor struct{}

func (OSCmdExecutor) ExecuteCommand(cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)

	output, err := command.Output()
	return string(output), err
}
