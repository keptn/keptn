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

type kubernetesPlatform struct {
}

func newKubernetesPlatform() *kubernetesPlatform {
	return &kubernetesPlatform{}
}

func (p kubernetesPlatform) getCreds() interface{} {
	return nil
}

func (p kubernetesPlatform) checkRequirements() error {
	if ctx, err := getKubeContext(); err != nil || ctx == "" {
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
	if ctx, err := getKubeContext(); err == nil && ctx != "" {
		return nil
	}
	return errors.New("Kubectl is not correctly configured")
}

func getKubeContext() (string, error) {
	return keptnutils.ExecuteCommand("kubectl", []string{
		"config",
		"current-context",
	})
}

func (p kubernetesPlatform) printCreds() {
	ctx, _ := getKubeContext()
	fmt.Println("Cluster: " + strings.TrimSpace(ctx))
}
