package main

import (
	"os"

	"github.com/keptn/keptn/cli/cmd"
	"github.com/keptn/keptn/cli/utils"
)

func main() {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	cmd.Execute()
}
