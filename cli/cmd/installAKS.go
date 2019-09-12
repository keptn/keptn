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

type aksCredentials struct {
	ClusterName        string `json:"clusterName"`
	AzureResourceGroup string `json:"azureResourceGroup"`
	AzureSubscription  string `json:"azureSubscription"`
}

type aksPlatform struct {
	creds *aksCredentials
}

func newAKSPlatform() *aksPlatform {
	return &aksPlatform{
		creds: &aksCredentials{},
	}
}

func (p aksPlatform) getCreds() interface{} {
	return p.creds
}

func (p aksPlatform) checkRequirements() error {
	_, err := getAzUser()
	return err
}

func (p aksPlatform) checkCreds() error {
	if p.creds.ClusterName == "" || p.creds.AzureResourceGroup == "" || p.creds.AzureSubscription == "" {
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

func (p aksPlatform) readCreds() {

	if p.creds.ClusterName == "" {
		p.creds.ClusterName = getAksClusterInfo()
	}

	connectionSuccessful := false
	for !connectionSuccessful {
		p.readClusterName()
		p.readAzureResourceGroup()
		p.readAzureSubscription()
		connectionSuccessful, _ = p.authenticateAtCluster()
	}
}

func (p aksPlatform) readClusterName() {
	readUserInput(&p.creds.ClusterName,
		"^(([a-zA-Z0-9]+-)*[a-zA-Z0-9]+)$",
		"Cluster Name",
		"Please enter a valid Cluster Name.",
	)
}

func (p aksPlatform) readAzureResourceGroup() {
	readUserInput(&p.creds.AzureResourceGroup,
		"^(([a-zA-Z0-9]+-)*[a-zA-Z0-9]+)$",
		"Azure Resource Group",
		"Please enter a valid Azure Resource Group.",
	)
}

func (p aksPlatform) readAzureSubscription() {
	readUserInput(&p.creds.AzureSubscription,
		"^(([a-zA-Z0-9]+-)*[a-zA-Z0-9]+)$",
		"Azure Subscription",
		"Please enter a valid Azure Subscription.",
	)
}

func (p aksPlatform) authenticateAtCluster() (bool, error) {
	_, err := keptnutils.ExecuteCommand("az", []string{
		"aks",
		"get-credentials",
		"--resource-group",
		p.creds.AzureResourceGroup,
		"--name",
		p.creds.ClusterName,
		"--subscription",
		p.creds.AzureSubscription,
		"--overwrite-existing",
	})

	if err != nil {
		fmt.Println("Could not connect to cluster. " +
			"Please verify that you have entered the correct information. Error: " + err.Error())
		return false, err
	}

	return true, nil
}

func getAksClusterInfo() string {
	// try to get current cluster from gcloud config
	out, err := keptnutils.ExecuteCommand("kubectl", []string{
		"config",
		"current-context",
	})

	if err != nil {
		return ""
	}
	return strings.TrimSpace(strings.Replace(string(out), "\r\n", "\n", -1))
}

func getAzUser() (string, error) {
	out, err := keptnutils.ExecuteCommand("az", []string{
		"account",
		"show",
		"--query=user.name",
		"--output=tsv",
	})

	if err != nil {
		return "", fmt.Errorf("Please configure az: %s", err)
	}
	// This command returns the account in the first line
	return strings.Split(strings.Replace(string(out), "\r\n", "\n", -1), "\n")[0], nil
}

func (p aksPlatform) printCreds() {
	fmt.Println("Cluster Name: " + p.creds.ClusterName)
	fmt.Println("Azure Resource Group: " + p.creds.AzureResourceGroup)
}
