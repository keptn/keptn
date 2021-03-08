// +build !nokubectl

// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package platform

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/keptn/keptn/cli/pkg/logging"

	keptnutils "github.com/keptn/kubernetes-utils/pkg"
)

type openShiftCredentials struct {
	OpenshiftURL      string `json:"openshiftUrl"`
	OpenshiftUser     string `json:"openshiftUser"`
	OpenshiftPassword string `json:"openshiftPassword"`
}

type openShiftPlatform struct {
	creds *openShiftCredentials
}

func newOpenShiftPlatform() *openShiftPlatform {
	return &openShiftPlatform{
		creds: &openShiftCredentials{},
	}
}

func (p openShiftPlatform) getCreds() interface{} {
	return p.creds
}

func (p openShiftPlatform) checkRequirements() error {
	return nil
}

func (p openShiftPlatform) checkCreds() error {
	if p.creds.OpenshiftURL == "" || p.creds.OpenshiftUser == "" || p.creds.OpenshiftPassword == "" {
		return errors.New("Incomplete credentials")
	}

	authenticated, err := p.authenticateAtCluster()
	if err != nil {
		return err
	}
	if !authenticated {
		return errors.New("Cannot authenticate at cluster " + p.creds.OpenshiftURL + ": " + err.Error())
	}
	return nil
}

func (p openShiftPlatform) readCreds() {

	fmt.Print("Please enter the following information or press enter to keep the old value:\n")

	connectionSuccessful := false
	for !connectionSuccessful {
		p.readOpenshiftClusterURL()
		p.readOpenshiftUser()
		p.readOpenshiftPassword()
		connectionSuccessful, _ = p.authenticateAtCluster()
	}
}

func (p openShiftPlatform) readOpenshiftClusterURL() {
	readUserInput(&p.creds.OpenshiftURL,
		"",
		"OpenShift Server URL",
		"Please enter a valid Server URL.",
	)
}

func (p openShiftPlatform) readOpenshiftUser() {
	readUserInput(&p.creds.OpenshiftUser,
		"",
		"OpenShift User",
		"Please enter a valid user name.",
	)
}

func (p openShiftPlatform) readOpenshiftPassword() {
	readUserInput(&p.creds.OpenshiftPassword,
		"",
		"OpenShift Password",
		"Please enter a valid password.",
	)
}

func readUserInput(value *string, regex string, promptMessage string, regexViolationMessage string) {
	var re *regexp.Regexp
	validateRegex := false
	if regex != "" {
		re = regexp.MustCompile(regex)
		validateRegex = true
	}
	keepAsking := true
	reader := bufio.NewReader(os.Stdin)
	for keepAsking {
		fmt.Printf("%s [%s]: ", promptMessage, *value)
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(strings.Replace(userInput, "\r\n", "\n", -1))
		if userInput != "" || *value == "" {
			if validateRegex && !re.MatchString(userInput) {
				fmt.Println(regexViolationMessage)
			} else {
				*value = userInput
				keepAsking = false
			}
		} else {
			keepAsking = false
		}
	}
}

func (p openShiftPlatform) authenticateAtCluster() (bool, error) {
	logging.PrintLog("Authenticating at OpenShift cluster: oc login "+p.creds.OpenshiftURL, logging.VerboseLevel)
	out, err := keptnutils.ExecuteCommand("oc", []string{
		"login",
		p.creds.OpenshiftURL,
		"-u=" + p.creds.OpenshiftUser,
		"-p=" + p.creds.OpenshiftPassword,
		"--insecure-skip-tls-verify=true",
	})

	logging.PrintLog("Result: "+out, logging.VerboseLevel)

	if err != nil {
		fmt.Println("Could not connect to cluster. Please verify that you have entered the correct information.")
		return false, err
	}

	return true, nil
}

func (p openShiftPlatform) printCreds() {
	fmt.Println("OpenShift Server URL: " + p.creds.OpenshiftURL)
	fmt.Println("OpenShift User: " + p.creds.OpenshiftUser)
}
