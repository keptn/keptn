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
	DeniedURLError
	RequestError
)
const (
	KubernetesSvcHostEnvVar = "KUBERNETES_SERVICE_HOST"
	KubernetesAPIPortEnvVar = "KUBERNETES_SERVICE_PORT"
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

func IsDeniedURLError(err error) bool {
	var curlErr *CurlError
	if errors.As(err, &curlErr) {
		return curlErr.reason == DeniedURLError
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
	deniedCharacters []string
	deniedOptions    []string
	requiredOptions  []string
	commandExecutor  ICommandExecutor
}

type CmdCurlExecutorOption func(executor *CmdCurlExecutor)

func NewCmdCurlExecutor(cmdExecutor ICommandExecutor, opts ...CmdCurlExecutorOption) *CmdCurlExecutor {
	executor := &CmdCurlExecutor{
		deniedCharacters: []string{"$", "|", ";", ">", "$(", " &", "&&", "`", "/var/run"},
		deniedOptions:    []string{"-o", "--output", "-F", "--form", "-T", "--upload-file", "-K", "--config"},
		requiredOptions:  []string{"--fail-with-body"},
		commandExecutor:  cmdExecutor,
	}
	for _, o := range opts {
		o(executor)
	}
	return executor
}

func (ce *CmdCurlExecutor) Curl(curlCmd string) (string, error) {
	args, err := ce.parseArgs(curlCmd)
	if err != nil {
		return "", err
	}

	resp, err := ce.commandExecutor.ExecuteCommand("curl", args[1:]...)
	if err != nil {
		return "", &CurlError{err: fmt.Errorf("error during curl request execution: %s.\nResponse: \n%s", err.Error(), resp), reason: RequestError}
	}
	return resp, nil
}

func (ce *CmdCurlExecutor) parseArgs(curlCmd string) ([]string, error) {
	cmdArr := strings.Split(curlCmd, " ")
	if len(cmdArr) == 0 || len(cmdArr) == 1 && cmdArr[0] == "" {
		return nil, &CurlError{err: errors.New("no command provided"), reason: NoCommandError}
	}

	for _, char := range ce.deniedCharacters {
		if strings.Contains(curlCmd, char) {
			return nil, &CurlError{err: fmt.Errorf("curl command contains denied character '%s'", char), reason: InvalidCommandError}
		}
	}

	args, err := parseCommandLine(curlCmd)
	if err != nil {
		return nil, &CurlError{err: errors.New("could not parse curl command"), reason: InvalidCommandError}
	}

	if cmdArr[0] != "curl" {
		return nil, &CurlError{err: errors.New("only curl commands are allowed to be executed"), reason: InvalidCommandError}
	}

	if err := ce.validateCurlOptions(args); err != nil {
		return nil, &CurlError{err: err, reason: InvalidCommandError}
	}

	args = ce.appendOptions(args)

	return args, nil
}

func (ce *CmdCurlExecutor) validateCurlOptions(args []string) error {
	for i, arg := range args {
		for _, o := range ce.deniedOptions {
			if strings.HasPrefix(arg, o) {
				return fmt.Errorf("curl command contains invalid option '%s'", o)
			}
		}
		// disallow usage of @ inside --data for posting local files
		if (arg == "--data" || arg == "-d") && len(args) >= i+1 {
			dataArgValue := args[i+1]
			if strings.HasPrefix(dataArgValue, "@") {
				return fmt.Errorf("file uploads using @ in --data is not allowed")
			}
		}
	}
	return nil
}

func (ce *CmdCurlExecutor) appendOptions(args []string) []string {
	addOptions := []string{}
	for _, requiredOption := range ce.requiredOptions {
		optionFound := false
		for _, arg := range args {
			if arg == requiredOption {
				optionFound = true
				break
			}
		}
		if optionFound {
			continue
		}
		addOptions = append(addOptions, requiredOption)
	}
	args = append(args, addOptions...)
	return args
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

	return deleteEmpty(args), nil
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
