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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type installConfig struct {
	LogLevel *string `json:"logLevel"`
}

var installCfg installConfig

type installCredentials struct {
	ClusterName               string `json:"clusterName"`
	ClusterZone               string `json:"clusterZone"`
	JenkinsUser               string `json:"jenkinsUser"`
	JenkinsPassword           string `json:"jenkinsPassword"`
	GithubPersonalAccessToken string `json:"githubPersonalAccessToken"`
	GithubUserEmail           string `json:"githubUserEmail"`
	GithubOrg                 string `json:"githubOrg"`
}

type keptnAPITokenSecret struct {
	APIVersion string `json:"apiVersion"`
	Data       struct {
		KeptnAPIToken string `json:"keptn-api-token"`
	} `json:"data"`
	Kind     string `json:"kind"`
	Metadata struct {
		CreationTimestamp time.Time `json:"creationTimestamp"`
		Name              string    `json:"name"`
		Namespace         string    `json:"namespace"`
		ResourceVersion   string    `json:"resourceVersion"`
		SelfLink          string    `json:"selfLink"`
		UID               string    `json:"uid"`
	} `json:"metadata"`
	Type string `json:"type"`
}

// installCmd represents the version command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs keptn on your Kubernetes cluster",
	Long: `Installs keptn on your Kubernetes cluster

Example:
	keptn install --log-level=INFO`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Installing keptn...")
		fmt.Printf("LogLevel=%s\n", *installCfg.LogLevel)
		// var creds *installCredentials
		clusterName, clusterZone, _ := connectToCluster()

		var err error
		var creds *installCredentials
		creds, err = getInstallCredentials()

		// get the YAML for the installer pod
		fileURL := "https://raw.githubusercontent.com/keptn/keptn/install-keptn/install/manifests/installer/installer.yaml"

		if err := DownloadFile("installer.yaml", fileURL); err != nil {
			panic(err)
		}

		gcloudUser := getGcloudUser()
		clusterIPCIDR, servicesIPCIDR := getGcloudClusterIPCIDR(clusterName, clusterZone)
		setDeploymentFileKey("JENKINS_USER", "admin")
		setDeploymentFileKey("JENKINS_PASSWORD", "AiTx4u8VyUV8tCKk")
		setDeploymentFileKey("GITHUB_PERSONAL_ACCESS_TOKEN", creds.GithubPersonalAccessToken)
		setDeploymentFileKey("GITHUB_USER_EMAIL", creds.GithubUserEmail)
		setDeploymentFileKey("GITHUB_ORGANIZATION", creds.GithubOrg)
		setDeploymentFileKey("GCLOUD_USER", gcloudUser)
		setDeploymentFileKey("CLUSTER_IPV4_CIDR", clusterIPCIDR)
		setDeploymentFileKey("SERVICES_IPV4_CIDR", servicesIPCIDR)

		execCmd := exec.Command(
			"kubectl",
			"apply",
			"-f",
			"https://raw.githubusercontent.com/keptn/keptn/install-keptn/install/manifests/installer/rbac.yaml",
		)

		_, err = execCmd.Output()

		if err != nil {
			log.Fatalf("Error while applying RBAC: %s \n", err)
		}
		fmt.Println("Deploying keptn installer pod....")
		execCmd = exec.Command(
			"kubectl",
			"apply",
			"-f",
			"installer.yaml",
		)
		_, err = execCmd.Output()
		if err != nil {
			log.Fatalf("Error while deploying keptn installer pod: %s \nAborting installation. \n", err)
		}
		fmt.Println("Installer pod deployed successfully.")

		installerPodName := waitForInstallerPod()

		getInstallerLogs(installerPodName)
		// installation finished, get auth token and endpoint
		setupKeptnAuth()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCfg.LogLevel = installCmd.Flags().StringP("log-level", "l", "INFO", "Log level")
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

	// default creds for jenkins
	creds.JenkinsUser = "admin"
	creds.JenkinsPassword = "AiTx4u8VyUV8tCKk"

	readGithubUserEmail(creds)
	readGithubPersonalAccessToken(creds)
	readGithubOrg(creds)

	// TODO: check if the github org exists and the access token has the necessary permissions

	newCreds, _ := json.Marshal(creds)
	newCredsStr := strings.Replace(string(newCreds), "\n", "", -1)
	credentialmanager.SetInstallCreds(newCredsStr)
	return creds, err
}

func readGithubUserEmail(creds *installCredentials) {
	readUserInput(
		&creds.GithubUserEmail,
		"^[a-z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,4}$",
		"GitHub User Email",
		"Please enter a valid email address.",
	)
}

func readGithubPersonalAccessToken(creds *installCredentials) {
	readUserInput(
		&creds.GithubPersonalAccessToken,
		"^[a-z0-9]{40}$",
		"GitHub Personal Access Token",
		"Please enter a valid GitHub Personal Access Token.",
	)
}

func readGithubOrg(creds *installCredentials) {
	readUserInput(
		&creds.GithubOrg,
		"^(([a-z0-9]+-)*[a-z0-9]+)$",
		"GitHub Organization",
		"Please enter a valid GitHub Organization.",
	)
}

func readUserInput(value *string, regex string, promptMessage string, regexViolationMessage string) {
	reader := bufio.NewReader(os.Stdin)
	var re *regexp.Regexp
	validateRegex := false
	if regex != "" {
		re = regexp.MustCompile(regex)
		validateRegex = true
	}
	keepAsking := true
	for ok := true; ok; ok = keepAsking {
		fmt.Printf("%s [%s]: ", promptMessage, *value)
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSuffix(userInput, "\n")
		if userInput != "" {
			if *value == "" && userInput == "" {
				fmt.Println("Please enter a value.")
			} else if validateRegex && !re.MatchString(userInput) {
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

func copyAndCapture(w io.Writer, r io.Reader, in io.WriteCloser) []byte {
	var log string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		log += scanner.Text() + "\n"
		if strings.HasPrefix(scanner.Text(), "[keptn|") {
			var reg = regexp.MustCompile(`\[keptn\|[a-zA-Z]+\]`)
			outputStr := reg.ReplaceAllString(scanner.Text(), "")
			fmt.Println(outputStr)
			if outputStr == "Installation of keptn complete." {
				cmd := exec.Command(
					"kubectl",
					"delete",
					"deployment",
					"installer",
				)
				cmd.Run()
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return []byte(log)
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func setDeploymentFileKey(key string, value string) {
	input, err := ioutil.ReadFile("installer.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	output := bytes.Replace(input, []byte("value: "+key), []byte("value: "+value), -1)

	if err = ioutil.WriteFile("installer.yaml", output, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func connectToCluster() (string, string, string) {
	reader := bufio.NewReader(os.Stdin)
	setupNewCluster := false
	clusterName, clusterZone, gkeProject := getClusterInfo()
	if clusterName != "" && clusterZone != "" && gkeProject != "" {
		fmt.Printf("Detected the following cluster:\n")
		fmt.Printf("Project: %s\n", gkeProject)
		fmt.Printf("Compute Zone: %s\n", clusterZone)
		fmt.Printf("Cluster Name: %s\n", clusterName)
		fmt.Printf("Would you like to use this cluster to set up keptn? [Y/n]")
		setupNewClusterPrompt, _ := reader.ReadString('\n')

		if setupNewClusterPrompt != "\n" && setupNewClusterPrompt != "y\n" && setupNewClusterPrompt != "Y\n" {
			setupNewCluster = true
		}
	} else {
		setupNewCluster = true
	}

	setupNewCluster = !authenticateAtCluster(clusterName, clusterZone, gkeProject)

	if setupNewCluster {
		connectionSuccessful := false
		for tryConnection := true; tryConnection; tryConnection = !connectionSuccessful {
			for ok := true; ok; ok = (gkeProject == "") {
				fmt.Printf("Please enter the GKE project [%s]:", gkeProject)
				newGkeProject, _ := reader.ReadString('\n')
				if newGkeProject != "\n" {
					gkeProject = strings.TrimSuffix(newGkeProject, "\n")
				}
			}

			for ok := true; ok; ok = (clusterName == "") {
				fmt.Printf("Please enter the GKE cluster name [%s]:", clusterName)
				newClusterName, _ := reader.ReadString('\n')
				if newClusterName != "\n" {
					clusterName = strings.TrimSuffix(newClusterName, "\n")
				}
			}

			for ok := true; ok; ok = (clusterZone == "") {
				fmt.Printf("Please enter the GKE cluster zone [%s]:", clusterZone)
				newClusterZone, _ := reader.ReadString('\n')
				if newClusterZone != "\n" {
					clusterZone = strings.TrimSuffix(newClusterZone, "\n")
				}
			}
			connectionSuccessful = authenticateAtCluster(clusterName, clusterZone, gkeProject)
		}
	}

	return clusterName, clusterZone, gkeProject
}

func authenticateAtCluster(clusterName string, clusterZone string, gkeProject string) bool {
	cmd := exec.Command(
		"gcloud",
		"container",
		"clusters",
		"get-credentials",
		clusterName,
		"--zone",
		clusterZone,
		"--project",
		gkeProject,
	)

	_, err := cmd.Output()
	if err != nil {
		fmt.Println("Could not connect to cluster. Please verify that you have entered the correct information.")
		return false
	}
	fmt.Println("Connection to cluster successful:")
	return true
}

func setGcloudConfig(key string, value string) {
	cmd := exec.Command(
		"gcloud",
		"config",
		"clusters",
		"set",
		key,
		value,
	)
	cmd.Run()
}

func getClusterInfo() (string, string, string) {
	fmt.Println("Trying to get GKE cluster info...")
	// try to get current cluster from gcloud config
	cmd := exec.Command("kubectl", "config", "current-context")
	out, err := cmd.Output()
	if err != nil {
		return "", "", ""
	}
	clusterInfo := strings.TrimSuffix(string(out), "\n")
	if !strings.HasPrefix(clusterInfo, "gke") {
		return "", "", ""
	}

	clusterInfoArray := strings.Split(clusterInfo, "_")
	if len(clusterInfoArray) < 4 {
		return "", "", ""
	}

	return clusterInfoArray[3], clusterInfoArray[2], clusterInfoArray[1]
}

func getGcloudUser() string {
	var err error

	cmd := exec.Command("gcloud", "config", "get-value", "account")
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return strings.Replace(string(out), "\n", "", -1)
}

func getGcloudClusterIPCIDR(clusterName string, clusterZone string) (string, string) {
	var err error
	var clusterDescription map[string]interface{}
	cmd := exec.Command("gcloud", "container", "clusters", "describe", clusterName, "--zone="+clusterZone)
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Could not get cluster Info. Aborting installation. \n")
	}

	err = yaml.Unmarshal([]byte(out), &clusterDescription)
	clusterIPCIDR := clusterDescription["clusterIpv4Cidr"].(string)
	servicesIPCIDR := clusterDescription["servicesIpv4Cidr"].(string)

	return clusterIPCIDR, servicesIPCIDR
}

func waitForInstallerPod() string {
	podName := ""
	podRunning := false
	for ok := true; ok; ok = !podRunning {
		time.Sleep(5 * time.Second)
		cmd := exec.Command(
			"kubectl",
			"get",
			"pods",
			"-l",
			"app=installer",
			"-ojson",
		)
		out, err := cmd.Output()
		if err != nil {
			log.Fatalf("Error while retrieving installer pod: %s\n. Aborting installation. \n", err)
		} else {
			var podInfo map[string]interface{}
			err = json.Unmarshal(out, &podInfo)
			if err == nil && podInfo != nil {
				podStatusArray := podInfo["items"].([]interface{})

				if len(podStatusArray) > 0 {
					podStatus := podStatusArray[0].(map[string]interface{})["status"].(map[string]interface{})["phase"].(string)
					if podStatus == "Running" {
						podName = podStatusArray[0].(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
						podRunning = true
					}
				}

			}
		}
	}
	return podName
}

func getInstallerLogs(podName string) {
	var stdout []byte
	var errStdout, errStderr error

	fmt.Printf("Getting logs of pod %s\n", podName)

	execCmd := exec.Command(
		"kubectl",
		"logs",
		podName,
		"-c",
		"keptn-installer",
		"-f",
	)

	stdoutIn, _ := execCmd.StdoutPipe()
	stderrIn, _ := execCmd.StderrPipe()
	stdIn, _ := execCmd.StdinPipe()
	err := execCmd.Start()
	if err != nil {
		log.Fatalf("Could not get installer pod logs: '%s'\n", err)
	}

	// cmd.Wait() should be called only after we finish reading
	// from stdoutIn and stderrIn.
	// wg ensures that we finish
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		stdout = copyAndCapture(os.Stdout, stdoutIn, stdIn)
		wg.Done()
	}()

	stdout = copyAndCapture(os.Stderr, stderrIn, stdIn)

	wg.Wait()

	err = execCmd.Wait()
	if err != nil {
		log.Fatalf("Could not get installer pod logs: '%s'\n", err)
	}
	if errStdout != nil || errStderr != nil {
		log.Fatal("Could not get installer pod logs.\n")
	}

	if err = ioutil.WriteFile("keptn-installer.log", stdout, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setupKeptnAuth() {
	cmd := exec.Command(
		"kubectl",
		"get",
		"secret",
		"keptn-api-token",
		"-n",
		"keptn",
		"-ojson",
	)

	out, err := cmd.Output()
	if err != nil {
		log.Fatal("Could not retrieve keptn API token.\n To manually set up your keptn CLI, please follow the instructions at https://keptn.sh/docs/0.2.0/reference/cli/.")
	}
	var secret keptnAPITokenSecret
	err = json.Unmarshal(out, &secret)
	if err != nil {
		log.Fatal("Could not retrieve keptn API token\n To manually set up your keptn CLI, please follow the instructions at https://keptn.sh/docs/0.2.0/reference/cli/.")
	}
	apiToken, err := base64.StdEncoding.DecodeString(secret.Data.KeptnAPIToken)
	if err != nil {
		log.Fatal("Could not retrieve keptn API token\n To manually set up your keptn CLI, please follow the instructions at https://keptn.sh/docs/0.2.0/reference/cli/.")
	}
	// $(kubectl get ksvc -n keptn control -o=yaml | yq r - status.domain)

	var keptnEndpoint string
	apiEndpointRetrieved := false
	retries := 0
	for tryGetAPIEndpoint := true; tryGetAPIEndpoint; tryGetAPIEndpoint = !apiEndpointRetrieved {
		cmd = exec.Command(
			"kubectl",
			"get",
			"ksvc",
			"-n",
			"keptn",
			"control",
			"-ojsonpath={.status.domain}",
		)

		out, err = cmd.Output()
		if err != nil {
			retries++
			if retries >= 15 {
				fmt.Println("API endpoint not yet available... trying again in 5s")
			}
		} else {
			retries = 0
		}
		keptnEndpoint = string(out)
		if keptnEndpoint == "" || !strings.Contains(keptnEndpoint, "xip.io") {
			retries++
			if retries >= 15 {
				fmt.Println("API endpoint not yet available... trying again in 5s")
			}
		} else {
			keptnEndpoint = "https://" + keptnEndpoint
			apiEndpointRetrieved = true
		}
		if !apiEndpointRetrieved {
			time.Sleep(5 * time.Second)
		}
	}
	fmt.Printf("Connecting to %s with API token %s\n", keptnEndpoint, apiToken)

	source, _ := url.Parse("https://github.com/keptn/keptn/cli#auth")
	contentType := "application/json"
	var data interface{}
	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Type:        "auth",
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
		}.AsV02(),
		Data: data,
	}

	u, err2 := url.Parse(keptnEndpoint)
	if err2 != nil {
		log.Fatal("Authentication at keptn API endpoint failed.\n To manually set up your keptn CLI, please follow the instructions at https://keptn.sh/docs/0.2.0/reference/cli/.")
	}

	authURL := *u
	authURL.Path = "auth"

	_, err = utils.Send(authURL, event, string(apiToken))
	if err != nil {
		fmt.Println("Authentication was unsuccessful")
		log.Fatal("Authentication at keptn API endpoint failed.\n To manually set up your keptn CLI, please follow the instructions at https://keptn.sh/docs/0.2.0/reference/cli/.")
	}

	fmt.Println("Successfully authenticated")
	credentialmanager.SetCreds(*u, string(apiToken))
	fmt.Println("You are now ready to use keptn.")

}
