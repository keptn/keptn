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
	"errors"
	"fmt"
	"strings"

	"github.com/keptn/keptn/cli/pkg/logging"

	keptnutils "github.com/keptn/kubernetes-utils/pkg"
)

type kubernetesCredentials struct{}

type kubernetesPlatform struct {
	creds *kubernetesCredentials
}

func newKubernetesPlatform() *kubernetesPlatform {
	return &kubernetesPlatform{
		creds: &kubernetesCredentials{},
	}
}

func (p kubernetesPlatform) getCreds() interface{} {
	return p.creds
}

func (p kubernetesPlatform) checkRequirements() error {
	if ctx, err := GetKubeContext(); err != nil || ctx == "" {
		return errors.New("kubectl is not properly configured. " +
			"Check your current context with 'kubectl config current-context'")
	}
	return nil
}

func (p kubernetesPlatform) readCreds() {

}

func (p kubernetesPlatform) authenticateAtCluster() (bool, error) {
	return false, nil
}

func (p kubernetesPlatform) checkCreds() error {
	if ctx, err := GetKubeContext(); err == nil && ctx != "" {
		return nil
	}
	return errors.New("Kubectl is not correctly configured")
}

func GetKubeContext() (string, error) {
	logging.PrintLog("Checking current Kubernetes context: kubectl config current-context", logging.VerboseLevel)
	out, err := keptnutils.ExecuteCommand("kubectl", []string{
		"config",
		"current-context",
	})
	logging.PrintLog("Result: "+out, logging.VerboseLevel)
	return out, err
}

func (p kubernetesPlatform) printCreds() {
	ctx, _ := GetKubeContext()
	fmt.Println("Cluster: " + strings.TrimSpace(ctx))
}
