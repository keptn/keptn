package main

import (
	"os"

	"github.com/keptn/keptn/cli/cmd"
	"github.com/keptn/keptn/cli/pkg/logging"
)

var (
	// Version information which is passed by ldflags
	Version = "develop"

	// DefaultKubeServerVersionConstraints is used when no version is passed by ldflags
	DefaultKubeServerVersionConstraints = ">= 1.14, <= 1.19"

	//KubeServerVersionConstraints the Kubernetes Cluster version's constraints is passed by ldflags
	KubeServerVersionConstraints = ""
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func main() {
	cmd.SetVersion(Version)

	if len(KubeServerVersionConstraints) > 0 {
		cmd.KubeServerVersionConstraints = KubeServerVersionConstraints
	} else {
		cmd.KubeServerVersionConstraints = DefaultKubeServerVersionConstraints
	}

	cmd.Execute()
}
