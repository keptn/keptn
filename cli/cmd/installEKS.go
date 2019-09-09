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
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
)

type eksCredentials struct {
	ClusterName string `json:"clusterName"`
	AwsRegion   string `json:"awsRegion"`
}

type eksPlatform struct {
	creds *eksCredentials
}

func newEKSPlatform() *eksPlatform {
	return &eksPlatform{
		creds: &eksCredentials{},
	}
}

func (p eksPlatform) getCreds() interface{} {
	return p.creds
}

func (p eksPlatform) checkRequirements() error {
	err := checkAWSCliVersion()
	if err != nil {
		return err
	}
	err = checkIfAWSIsConfigured()
	if err != nil {
		return err
	}
	return checkStsCallerIdentity()
}

func (p eksPlatform) checkCreds() error {
	if p.creds.ClusterName == "" || p.creds.AwsRegion == "" {
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

func (p eksPlatform) readCreds() {

	connectionSuccessful := false
	for !connectionSuccessful {
		p.readClusterName()
		p.readAwsRegion()
		connectionSuccessful, _ = p.authenticateAtCluster()
	}
}

func (p eksPlatform) readClusterName() {
	readUserInput(&p.creds.ClusterName,
		"^(([a-zA-Z0-9]+-)*[a-zA-Z0-9]+)$",
		"Cluster Name",
		"Please enter a valid Cluster Name.",
	)
}

func (p eksPlatform) readAwsRegion() {
	readUserInput(&p.creds.AwsRegion,
		"^(([a-zA-Z0-9]+-)*[a-zA-Z0-9]+)$",
		"AWS Region",
		"Please enter a valid AWS Region.",
	)
}

func (p eksPlatform) authenticateAtCluster() (bool, error) {
	_, err := keptnutils.ExecuteCommand("aws", []string{
		"eks",
		"--region",
		p.creds.AwsRegion,
		"update-kubeconfig",
		"--name",
		p.creds.ClusterName,
	})

	if err != nil {
		fmt.Println("Could not connect to cluster. " +
			"Please verify that you have entered the correct information. Error: " + err.Error())
		return false, err
	}

	return true, nil
}

type stsCallerIdentity struct {
	Account string `json:"Account"`
	UserID  string `json:"UserId"`
	Arn     string `json:"Arn"`
}

func checkStsCallerIdentity() error {
	out, err := keptnutils.ExecuteCommand("aws", []string{
		"sts",
		"get-caller-identity",
	})

	if err != nil {
		return fmt.Errorf("Please configure your sts caller-identity: %s", err)
	}

	callerIdentity := &stsCallerIdentity{}
	err = json.Unmarshal([]byte(out), &callerIdentity)
	if err != nil {
		return err
	}
	if callerIdentity.Account == "" || callerIdentity.UserID == "" || callerIdentity.Arn != "" {
		return errors.New("Please check your sts -caller-idenity")
	}
	return nil
}

func checkIfAWSIsConfigured() error {

	out, err := keptnutils.ExecuteCommand("aws", []string{
		"configure", "get", "aws_access_key_id",
	})

	if err != nil || out == "" {
		return fmt.Errorf("Please configure your aws CLI: %s", err)
	}
	return nil
}

func checkAWSCliVersion() error {

	out, err := keptnutils.ExecuteCommand("aws", []string{
		"--version",
	})

	if err != nil {
		return err
	}

	re := regexp.MustCompile("aws-cli/(.*?) ")
	version := strings.TrimSpace(re.FindString(out))
	if version == "" {
		return errors.New("Please install the aws CLI at least in version 1.16.156")
	}
	vNo := strings.Split(version[len("aws-cli/"):], ".")
	v1, err := strconv.Atoi(vNo[0])
	if err != nil {
		return err
	}
	v2, err := strconv.Atoi(vNo[1])
	if err != nil {
		return err
	}
	v3, err := strconv.Atoi(vNo[2])
	if err != nil {
		return err
	}

	if v1 >= 1 && v2 >= 16 && v3 >= 156 {
		return nil
	}
	return errors.New("Please install the AWS CLI at least in version 1.16.156")
}

func (p eksPlatform) printCreds() {
	fmt.Println("Cluster Name: " + p.creds.ClusterName)
	fmt.Println("AWS Region: " + p.creds.AwsRegion)
}
