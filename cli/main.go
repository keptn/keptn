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

const (
	// DefaultKubeServerVersionConstraints is used when no version is passed by ldflags
	DefaultKubeServerVersionConstraints = ">= 1.14, <= 1.22"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func main() {
	cmd.KubeServerVersionConstraints = DefaultKubeServerVersionConstraints
	cmd.Execute()
}
