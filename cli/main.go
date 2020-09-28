package main

import (
	"os"

	"github.com/keptn/keptn/cli/cmd"
	"github.com/keptn/keptn/cli/pkg/logging"
)

var (
	// Version information which is passed by ldflags
	Version = "develop"

	// KubeServerVersionConstraints the Kubernetes Cluster version's constraints is passed by ldflags
	KubeServerVersionConstraints = ">= 1.14, <= 1.19"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func main() {
	cmd.Version = Version
	cmd.KubeServerVersionConstraints = KubeServerVersionConstraints
	cmd.Execute()
}
