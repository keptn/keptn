package main

import (
	"github.com/hashicorp/go-version"
	"os"

	"github.com/keptn/keptn/cli/cmd"
	"github.com/keptn/keptn/cli/pkg/logging"
)

var (
	// Version information which is passed by ldflags
	Version = "develop"

	// DefaultKubeServerVersionConstraints is used when no version is passed by ldflags
	DefaultKubeServerVersionConstraints = ">= 1.14, <= 1.20"

	//KubeServerVersionConstraints the Kubernetes Cluster version's constraints is passed by ldflags
	KubeServerVersionConstraints = ""
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func main() {
	cmd.SetVersion(Version)

	if len(KubeServerVersionConstraints) > 0 {
		if _, err := version.NewConstraint(KubeServerVersionConstraints); err != nil {
			cmd.KubeServerVersionConstraints = DefaultKubeServerVersionConstraints
		} else {
			cmd.KubeServerVersionConstraints = KubeServerVersionConstraints
		}
	} else {
		cmd.KubeServerVersionConstraints = DefaultKubeServerVersionConstraints
	}

	cmd.Execute()
}
