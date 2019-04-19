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
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/spf13/cobra"
)

type installCredentials struct {
	ClusterName               string `json:"clusterName"`
	ClusterZone               string `json:"clusterZone"`
	JenkinsUser               string `json:"jenkinsUser"`
	JenkinsPassword           string `json:"jenkinsPassword"`
	GithubPersonalAccessToken string `json:"githubPersonalAccessToken"`
	GithubUserEmail           string `json:"githubUserEmail"`
	GithubOrg                 string `json:"githubOrg"`
}

// installCmd represents the version command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs keptn on your Kubernetes cluster",
	Long: `Installs keptn on your Kubernetes cluster

Example:
	keptn install`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Installing keptn...")
		// var creds *installCredentials
		var err error
		_, err = getInstallCredentials()
		execCmd := exec.Command("kubectl", "logs", "control-lxxt6-deployment-569499f6cb-gbzgp", "-n", "keptn", "-c", "user-container", "-f")

		var stdout, stderr []byte
		var errStdout, errStderr error
		stdoutIn, _ := execCmd.StdoutPipe()
		stderrIn, _ := execCmd.StderrPipe()
		err = execCmd.Start()
		if err != nil {
			log.Fatalf("cmd.Start() failed with '%s'\n", err)
		}

		// cmd.Wait() should be called only after we finish reading
		// from stdoutIn and stderrIn.
		// wg ensures that we finish
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			copyAndCapture(os.Stdout, stdoutIn)
			wg.Done()
		}()

		copyAndCapture(os.Stderr, stderrIn)

		wg.Wait()

		err = execCmd.Wait()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		if errStdout != nil || errStderr != nil {
			log.Fatal("failed to capture stdout or stderr\n")
		}
		outStr, errStr := string(stdout), string(stderr)
		fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func getInstallCredentials() (*installCredentials, error) {
	var err error
	var credsStr string
	var creds *installCredentials
	credsStr, err = credentialmanager.GetInstallCreds()
	if err != nil {
		credsStr = ""
	}
	err = json.Unmarshal([]byte(credsStr), &creds)
	if err != nil {
		creds = new(installCredentials)
	}

	fmt.Print("Please enter the following information (press enter to keep the old value):\n")

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("GKE cluster name [%s]: ", creds.ClusterName)
	clusterName, _ := reader.ReadString('\n')
	if clusterName != "\n" {
		creds.ClusterName = strings.TrimSuffix(clusterName, "\n")
	}

	fmt.Printf("GKE cluster zone [%s]: ", creds.ClusterZone)
	clusterZone, _ := reader.ReadString('\n')
	if clusterName != "\n" {
		creds.ClusterZone = strings.TrimSuffix(clusterZone, "\n")
	}

	// default creds for jenkins
	creds.JenkinsUser = "admin"
	creds.JenkinsPassword = "AiTx4u8VyUV8tCKk"

	fmt.Printf("GitHub User Email [%s]: ", creds.GithubUserEmail)
	githubUserEmail, _ := reader.ReadString('\n')
	if githubUserEmail != "\n" {
		creds.GithubUserEmail = strings.TrimSuffix(githubUserEmail, "\n")
	}

	fmt.Printf("GitHub Organization [%s]: ", creds.GithubOrg)
	githubOrg, _ := reader.ReadString('\n')
	if githubOrg != "\n" {
		creds.GithubOrg = strings.TrimSuffix(githubOrg, "\n")
	}

	fmt.Printf("GitHub Personal Access Token [%s]: ", creds.GithubPersonalAccessToken)
	githubPersonalAccessToken, _ := reader.ReadString('\n')
	if githubPersonalAccessToken != "\n" {
		creds.GithubPersonalAccessToken = strings.TrimSuffix(githubPersonalAccessToken, "\n")
	}

	newCreds, _ := json.Marshal(creds)
	newCredsStr := strings.Replace(string(newCreds), "\n", "", -1)
	fmt.Printf("new creds file content: %s", newCredsStr)
	credentialmanager.SetInstallCreds(newCredsStr)
	return nil, err
}

func copyAndCapture(w io.Writer, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := strings.Replace(scanner.Text(), "found", "", -1)
		fmt.Println(text)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
