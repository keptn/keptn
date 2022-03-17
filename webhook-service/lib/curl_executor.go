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
	requiredOptions     []string
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
		unAllowedCharacters: []string{"$", "|", ";", ">", "$(", " &", "&&", "`", "/var/run"},
		unAllowedOptions:    []string{"-o", "--output", "-F", "--form", "-T", "--upload-file", "-K", "--config"},
		requiredOptions:     []string{"--fail-with-body"},
		commandExecutor:     cmdExecutor,
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

	for _, unallowedCharacter := range ce.unAllowedCharacters {
		if strings.Contains(curlCmd, unallowedCharacter) {
			return nil, &CurlError{err: fmt.Errorf("curl command contains unallowed character '%s'", unallowedCharacter), reason: InvalidCommandError}
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

	// check if the curl command contains any of the disallowed URLs
	if err := ce.validateURL(curlCmd); err != nil {
		return nil, &CurlError{err: err, reason: UnallowedURLError}
	}

	args = ce.appendOptions(args)

	return args, nil
}

func (ce *CmdCurlExecutor) validateURL(curlCmd string) error {
	sanitizedCurlCmd := strings.ReplaceAll(curlCmd, "\\", "")
	for _, url := range ce.unAllowedURLs {
		if strings.Contains(sanitizedCurlCmd, url) {
			return fmt.Errorf("curl command contains invalid URL %s", url)
		}
	}
	return nil
}

func (ce *CmdCurlExecutor) validateCurlOptions(args []string) error {
	for i, arg := range args {
		for _, o := range ce.unAllowedOptions {
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

func BlacklistedKubeURLS(env map[string]string) []string {
	kubeAPIHostIP := env["KUBERNETES_SERVICE_HOST"]
	kubeAPIPort := env["KUBERNETES_SERVICE_PORT"]
	return []string{
		// Block access to Kubernetes API
		kubeAPIHostIP,
		kubeAPIHostIP + ":" + kubeAPIPort,
		"kubernetes",
		"kubernetes" + ":" + kubeAPIPort,
		"kubernetes.default",
		"kubernetes.default" + ":" + kubeAPIPort,
		"kubernetes.default.svc",
		"kubernetes.default.svc" + ":" + kubeAPIPort,
		"kubernetes.default.svc.cluster.local",
		"kubernetes.default.svc.cluster.local" + ":" + kubeAPIPort,
		// Block access to localhost
		"localhost",
		"127.0.0.1",
		"::1",
	}
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
