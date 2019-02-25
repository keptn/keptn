package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/keptn/keptn/cli/cmd"
	"github.com/keptn/keptn/cli/utils"
)

func init() {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
}

const logoFileName = "logo.txt"

func main() {
	dat, err := ioutil.ReadFile(logoFileName)
	if err == nil {
		fmt.Println(string(dat))
	}

	cmd.Execute()
}
