package exechelper

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mattn/go-shellwords"
)

// ExecuteCommand executes the command using the args
func ExecuteCommand(command string, argStr string) (string, error) {
	args, err := shellwords.Parse(argStr)
	if err != nil {
		return "", err
	}
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Error executing command %s %s: %s\n%s", command, strings.Join(args, " "), err.Error(), string(out))
	}
	return string(out), nil
}
