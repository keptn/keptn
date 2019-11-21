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

package cmd

import (
	"errors"
	"fmt"
	"strings"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
)

type pksCredentials struct {
	PksAPI          string `json:"pksAPI"`
	ClusterName     string `json:"clusterName"`
	PksUser         string `json:"pksUser"`
	PksUserPassword string `json:"pksUserPassword"`
}

type pksPlatform struct {
	creds *pksCredentials
}

func newPKSPlatform() *pksPlatform {
	return &pksPlatform{
		creds: &pksCredentials{},
	}
}

func (p pksPlatform) getCreds() interface{} {
	return p.creds
}

func (p pksPlatform) checkRequirements() error {
	_, err := getPksUser()
	return err
}

func (p pksPlatform) checkCreds() error {
	if p.creds.PksAPI == "" || p.creds.ClusterName == "" || p.creds.PksUser == "" || p.creds.PksUserPassword == "" {
		return errors.New("Incomplete credentials")
	}

	authenticated, err := p.authenticateAtCluster()
	if err != nil {
		return err
	}
	if !authenticated {
		return errors.New("Cannot authenticate at cluster " + p.creds.ClusterName)
	}
	return nil
}

func (p pksPlatform) readCreds() {

	if p.creds.PksAPI == "" || p.creds.ClusterName == "" || p.creds.PksUser == "" || p.creds.PksUserPassword == "" {
		p.creds.ClusterName = getPksClusterInfo()
	}

	connectionSuccessful := false
	for !connectionSuccessful {
		p.readClusterEndpoint()
		p.readClusterName()
		p.readPksUser()
		p.readPksUserPassword()
		connectionSuccessful, _ = p.authenticateAtCluster()
	}
}

func (p pksPlatform) readClusterEndpoint() {
	readUserInput(&p.creds.PksAPI,
		"",
		"PKS API",
		"Please enter a valid PKS API.",
	)
}

func (p pksPlatform) readClusterName() {
	readUserInput(&p.creds.ClusterName,
		"^(([a-zA-Z0-9]+-)*[a-zA-Z0-9]+)$",
		"Cluster Name",
		"Please enter a valid Cluster Name.",
	)
}

func (p pksPlatform) readPksUser() {
	readUserInput(&p.creds.PksUser,
		"^(([a-zA-Z0-9]+-)*[a-zA-Z0-9]+)$",
		"PKS User",
		"Please enter your PKS User.",
	)
}

func (p pksPlatform) readPksUserPassword() {
	readUserInput(&p.creds.PksUserPassword,
		"",
		"PKS User Password",
		"Please enter your PKS User Password.",
	)
}

func (p pksPlatform) authenticateAtCluster() (bool, error) {
	_, err := keptnutils.ExecuteCommand("pks", []string{
		"login",
		"clusters",
		"--api",
		p.creds.PksAPI,
		"--username",
		p.creds.PksUser,
		"--password",
		p.creds.PksUserPassword,
		"--skip-ssl-verification",
	})

	if err != nil {
		fmt.Println("Could not connect to cluster. " +
			"Please verify that you have entered the correct information. Error: " + err.Error())
		return false, err
	}

	_, err = keptnutils.ExecuteCommand("pks", []string{
		"get-credentials",
		p.creds.ClusterName,
	})

	if err != nil {
		fmt.Println("Could not update a kubeconfig file. Error: " + err.Error())
		return false, err
	}

	return true, nil
}

func getPksClusterInfo() string {
	// try to get current cluster from kubectl config
	out, err := keptnutils.ExecuteCommand("kubectl", []string{
		"config",
		"current-context",
	})

	if err != nil {
		return ""
	}
	return strings.TrimSpace(strings.Replace(string(out), "\r\n", "\n", -1))
}

func getPksUser() (string, error) {
	// the PKS CLI does not support it to retrieve the user
	return "Not implemented yet", nil
}

func (p pksPlatform) printCreds() {
	fmt.Println("PKS API: " + p.creds.PksAPI)
	fmt.Println("Cluster Name: " + p.creds.ClusterName)
	fmt.Println("PKS User: " + p.creds.PksUser)
}
