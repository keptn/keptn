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
	ClusterEndpoint string `json:"clusterEndpoint"`
	ClusterName     string `json:"clusterName"`
	PksUser         string `json:"pksUser"`
	PksPassword     string `json:"pksUser"`
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
	_, err := getGcloudUser()
	return err
}

func (p pksPlatform) checkCreds() error {
	if p.creds.ClusterName == "" || p.creds.ClusterEndpoint == "" {
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

	if p.creds.ClusterEndpoint == "" || p.creds.ClusterName == "" || p.creds.PksUser == "" || p.creds.PksPassword == "" {
		p.creds.ClusterEndpoint, p.creds.ClusterName, p.creds.PksUser, p.creds.PksPassword = getPksClusterInfo()
	}

	connectionSuccessful := false
	for !connectionSuccessful {
		p.readClusterName()
		p.readClusterEndpoint()
		p.readGkeProject()
		connectionSuccessful, _ = p.authenticateAtCluster()
	}
}

func (p pksPlatform) readClusterName() {
	readUserInput(&p.creds.ClusterName,
		"^(([a-zA-Z0-9]+-)*[a-zA-Z0-9]+)$",
		"Cluster Name",
		"Please enter a valid Cluster Name.",
	)
}

func (p pksPlatform) readClusterEndpoint() {
	readUserInput(&p.creds.ClusterEndpoint,
		"^(([a-zA-Z0-9]+-)*[a-zA-Z0-9]+)$",
		"Cluster Endpoint",
		"Please enter a valid Cluster Endpoint.",
	)
}

func (p pksPlatform) readGkeProject() {
	readUserInput(&p.creds.GkeProject,
		"^(([a-zA-Z0-9]+-)*[a-zA-Z0-9]+)$",
		"GKE Project",
		"Please enter a valid GKE Project.",
	)
}

func (p pksPlatform) authenticateAtCluster() (bool, error) {
	_, err := keptnutils.ExecuteCommand("gcloud", []string{
		"container",
		"clusters",
		"get-credentials",
		p.creds.ClusterName,
		"--zone",
		p.creds.ClusterZone,
		"--project",
		p.creds.GkeProject,
	})

	if err != nil {
		fmt.Println("Could not connect to cluster. " +
			"Please verify that you have entered the correct information. Error: " + err.Error())
		return false, err
	}

	return true, nil
}

func getPksClusterInfo() (string, string, string, string) {
	// try to get current cluster from gcloud config
	out, err := getKubeContext()

	if err != nil {
		return "", "", "", ""
	}
	clusterInfo := strings.TrimSpace(strings.Replace(string(out), "\r\n", "\n", -1))
	if !strings.HasPrefix(clusterInfo, pks) {
		return "", "", "", ""
	}

	clusterInfoArray := strings.Split(clusterInfo, "_")
	if len(clusterInfoArray) < 4 {
		return "", "", "", ""
	}

	return "", clusterInfoArray[3], clusterInfoArray[2], clusterInfoArray[1]
}

func getGcloudUser() (string, error) {
	out, err := keptnutils.ExecuteCommand("gcloud", []string{
		"config",
		"get-value",
		"account",
	})

	if err != nil {
		return "", fmt.Errorf("Please configure your gcloud: %s", err)
	}
	// This command returns the account in the first line
	return strings.Split(strings.Replace(string(out), "\r\n", "\n", -1), "\n")[0], nil
}

func (p pksPlatform) printCreds() {
	fmt.Println("Cluster endpoint: " + p.creds.ClusterEndpoint)
	fmt.Println("Cluster name: " + p.creds.ClusterName)
	fmt.Println("PKS user: " + p.creds.PksUser)
}
