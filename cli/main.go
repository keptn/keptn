package main

import (
	"os"

	"github.com/keptn/keptn/cli/cmd"
	"github.com/keptn/keptn/cli/pkg/logging"
)

var (
	// Version information which is passed by ldflags
	Version = "develop"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func main() {
	cmd.Version = Version
	cmd.Execute()
}
