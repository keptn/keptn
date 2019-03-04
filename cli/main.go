package main

import (
	"os"

	"github.com/keptn/keptn/cli/cmd"
	"github.com/keptn/keptn/cli/utils"
)

func init() {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
}

func main() {
	cmd.Execute()
}
