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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/spf13/cobra"
)

var configFilePath *string
var installerVersion *string
var platform *string

const gke = "gke"
const aks = "aks"
const openshift = "openshift"
const kubernetes = "kubernetes"

const installerPrefixURL = "https://raw.githubusercontent.com/keptn/keptn/"
const installerSuffixPath = "/installer/manifests/installer/installer.yaml"
const rbacSuffixPath = "/installer/manifests/installer/rbac.yaml"

type installCredentials struct {
	GithubPersonalAccessToken string `json:"githubPersonalAccessToken"`
	GithubOrg                 string `json:"githubOrg"`
	GithubUserName            string `json:"githubUserName"`
	ClusterName               string `json:"clusterName"`
	ClusterZone               string `json:"clusterZone"`
	GkeProject                string `json:"gkeProject"`
	OpenshiftURL              string `json:"openshiftUrl"`
	OpenshiftUser             string `json:"openshiftUser"`
	OpenshiftPassword         string `json:"openshiftPassword"`
	ClusterIPCIDR             string `json:"clusterIPCIDR"`
	ServicesIPCIDR            string `json:"servicesIPCIDR"`
	AzureResourceGroup        string `json:"azureResourceGroup"`
	AzureSubscription         string `json:"azureSubscription"`
}

type placeholderReplacement struct {
	placeholderValue string
	desiredValue     string
}

// installCmd represents the version command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs keptn on your Kubernetes cluster",
	Long: `Installs keptn on your Kubernetes cluster

Example:
	keptn install`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {

		err := checkIfPlatformIsSupported()
		if err != nil {
			return err
		}

		isInstallerAvailable, err := checkInstallerAvailablity()
		if err != nil || !isInstallerAvailable {
			return errors.New("Installers not found under:\n" +
				getInstallerURL() + "\n" + getRbacURL())
		}

		if *platform == gke {
			// Check whether Gcloud user is configured
			_, err = getGcloudUser()
			if err != nil {
				return err
			}
		}
		if *platform == aks {
			_, err = getAzUser()
			if err != nil {
				return err
			}
		}

		// Check whether kubectl is installed
		isKubAvailable, err := isKubectlAvailable()
		if err != nil || !isKubAvailable {
			return errors.New(`keptn requires 'kubectl' but it is not available.
Please see https://kubernetes.io/docs/tasks/tools/install-kubectl/`)
		}

		if configFilePath != nil && *configFilePath != "" {
			// Config was provided in form of a file
			creds, err := parseConfig(*configFilePath)
			if err != nil {
				return err
			}
			// Verify the provided config
			// Check whether all data is provided
			if creds.GithubPersonalAccessToken == "" ||
				creds.GithubOrg == "" || creds.GithubUserName == "" {
				return errors.New("Incomplete credential file " + *configFilePath)
			}

			// Check whether the authentication at the cluster is valid
			var authenticated = false
			if *platform == gke {
				if creds.ClusterName == "" || creds.ClusterZone == "" {
					return errors.New("Incomplete credential file " + *configFilePath)
				}
				authenticated, err = authenticateAtGkeCluster(creds)
			} else if *platform == openshift {
				if creds.OpenshiftURL == "" || creds.OpenshiftUser == "" || creds.OpenshiftPassword == "" {
					return errors.New("Incomplete credential file " + *configFilePath)
				}
				authenticated, err = authenticateAtOpenshiftCluster(creds)
			} else if *platform == aks {
				if creds.ClusterName == "" || creds.AzureResourceGroup == "" || creds.AzureSubscription == "" {
					return errors.New("Incomplete credential file " + *configFilePath)
				}
				authenticated, err = authenticateAtAksCluster(creds)
			} else if *platform == kubernetes {
				if ctx, err := getKubeContext(); err == nil && ctx != "" {
					authenticated = true
				}
			}
			if err != nil {
				return err
			}
			if !authenticated {
				return errors.New("Cannot authenticate at cluster " + creds.ClusterName)
			}
			// Check GitHub token and org
			validScopeRes, err := utils.HasTokenRepoScope(creds.GithubPersonalAccessToken)
			if err != nil {
				return err
			}
			if !validScopeRes {
				return errors.New("Personal access token requies at least a 'repo'-scope")
			}
			validOrg, err := utils.IsOrgExisting(creds.GithubPersonalAccessToken, creds.GithubOrg)
			if err != nil {
				return err
			}
			if !validOrg {
				return errors.New("Provided organization " + creds.GithubOrg + " does not exist.")
			}
		} else if platform != nil && *platform == kubernetes {
			if ctx, err := getKubeContext(); err != nil || ctx == "" {
				return errors.New("kubectl is not properly configured. " +
					"Check your current context with 'kubectl config current-context'")
			}
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		utils.PrintLog("Installing keptn...", utils.InfoLevel)

		var creds installCredentials
		var err error
		if !mocking {
			if configFilePath != nil && *configFilePath != "" {
				creds, err = parseConfig(*configFilePath)
				if err != nil {
					return err
				}

			} else {
				err = getInstallCredentials(&creds)
				if err != nil {
					return err
				}
			}
			return doInstallation(creds)
		}
		fmt.Println("Skipping intallation due to mocking flag set to true")
		return nil
	},
}

func checkIfPlatformIsSupported() error {
	if platform == nil {
		*platform = gke
		return nil
	}
	*platform = strings.ToLower(*platform)
	if *platform != gke && *platform != aks && *platform != openshift && *platform != kubernetes {
		return errors.New("Unsupported platform '" + *platform + "'. The following platforms are supported: gke, aks, openshift, kubernetes")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(installCmd)
	configFilePath = installCmd.Flags().StringP("creds", "c", "", "The name of the creds file")
	installCmd.Flags().MarkHidden("creds")
	installerVersion = installCmd.Flags().StringP("keptn-version", "k", "master", "The branch or tag of the version which is installed")
	installCmd.Flags().MarkHidden("keptn-version")
	platform = installCmd.Flags().StringP("platform", "p", "gke", "The platform to run keptn on [gke,openshift,aks,kubernetes]")
}

func checkInstallerAvailablity() (bool, error) {

	resp, err := http.Get(getInstallerURL())
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	resp, err = http.Get(getRbacURL())
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}
	return true, nil
}

func getInstallerURL() string {
	return installerPrefixURL + *installerVersion + installerSuffixPath
}

func getRbacURL() string {
	return installerPrefixURL + *installerVersion + rbacSuffixPath
}

// Preconditions: 1. Already authenticated against the cluster; 2. Github credentials are checked
func doInstallation(creds installCredentials) error {

	path, err := keptnutils.GetKeptnDirectory()
	if err != nil {
		return err
	}
	installerPath := path + "installer.yaml"

	// get the YAML for the installer pod
	if err := downloadFile(installerPath, getInstallerURL()); err != nil {
		return err
	}

	if err := setDeploymentFileKey(installerPath,
		placeholderReplacement{"PLATFORM", *platform}); err != nil {
		return err
	}

	if *platform == gke || *platform == aks || *platform == kubernetes {
		_, err := keptnutils.ExecuteCommand("kubectl", []string{
			"apply",
			"-f",
			getRbacURL(),
		})

		if err != nil {
			return fmt.Errorf("Error while applying RBAC for installer pod: %s \nAborting installation", err.Error())
		}
	}

	utils.PrintLog("Deploying keptn installer pod...", utils.InfoLevel)

	_, err = keptnutils.ExecuteCommand("kubectl", []string{
		"apply",
		"-f",
		installerPath,
	})

	if err != nil {
		return fmt.Errorf("Error while deploying keptn installer pod: %s \nAborting installation", err.Error())
	}

	utils.PrintLog("Installer pod deployed successfully.", utils.InfoLevel)

	installerPodName, err := waitForInstallerPod()
	if err != nil {
		return err
	}

	err = getInstallerLogs(installerPodName)
	if err != nil {
		return err
	}
	// installation finished, get auth token and endpoint
	err = setupKeptnAuthAndConfigure(creds)
	if err != nil {
		return err
	}

	return os.Remove(installerPath)
}

func parseConfig(configFile string) (installCredentials, error) {
	data, err := utils.ReadFile(configFile)
	if err != nil {
		return installCredentials{}, err
	}
	var installCreds installCredentials
	json.Unmarshal([]byte(data), &installCreds)
	return installCreds, nil
}

func getInstallCredentials(creds *installCredentials) error {

	credsStr, err := credentialmanager.GetInstallCreds()
	if err != nil {
		credsStr = ""
	}
	// Ignore unmarshaling error
	json.Unmarshal([]byte(credsStr), &creds)

	fmt.Print("Please enter the following information or press enter to keep the old value:\n")

	for {
		connectToCluster(creds)

		readGithubUserName(creds)

		// Check if the access token has the necessary permissions and the github org exists
		validScopeRes := false
		for !validScopeRes {
			readGithubPersonalAccessToken(creds)
			validScopeRes, err = utils.HasTokenRepoScope(creds.GithubPersonalAccessToken)
			if err != nil {
				return err
			}
			if !validScopeRes {
				fmt.Println("GitHub Personal Access Token requies at least a 'repo'-scope")
				creds.GithubPersonalAccessToken = ""
			}
		}
		validOrg := false
		for !validOrg {
			readGithubOrg(creds)
			validOrg, err = utils.IsOrgExisting(creds.GithubPersonalAccessToken, creds.GithubOrg)
			if err != nil {
				return err
			}
			if !validOrg {
				fmt.Println("Provided GitHub Organization " + creds.GithubOrg + " does not exist.")
				creds.GithubOrg = ""
			}
		}

		fmt.Println()
		fmt.Println("Please confirm that the provided information is correct: ")

		if *platform == gke {
			fmt.Println("Cluster Name: " + creds.ClusterName)
			fmt.Println("Cluster Zone: " + creds.ClusterZone)
			fmt.Println("GKE Project: " + creds.GkeProject)
		} else if *platform == openshift {
			fmt.Println("Openshift Server URL: " + creds.OpenshiftURL)
			fmt.Println("Openshift User: " + creds.OpenshiftUser)
		} else if *platform == aks {
			fmt.Println("Cluster Name: " + creds.ClusterName)
			fmt.Println("Azure Resource Group: " + creds.AzureResourceGroup)
		} else if *platform == kubernetes {
			ctx, _ := getKubeContext()
			fmt.Println("Cluster: " + strings.TrimSpace(ctx))
		}

		fmt.Println("GitHub User Name: " + creds.GithubUserName)
		fmt.Println("GitHub Personal Access Token: " + creds.GithubPersonalAccessToken)
		fmt.Println("GitHub Organization: " + creds.GithubOrg)

		fmt.Println()
		fmt.Println("Is this all correct? (y/n)")

		reader := bufio.NewReader(os.Stdin)
		in, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		in = strings.TrimSpace(in)
		if in == "y" || in == "yes" {
			break
		}
	}

	newCreds, _ := json.Marshal(creds)
	newCredsStr := strings.Replace(string(newCreds), "\r\n", "\n", -1)
	newCredsStr = strings.Replace(newCredsStr, "\n", "", -1)
	return credentialmanager.SetInstallCreds(newCredsStr)
}

func connectToCluster(creds *installCredentials) {
	if *platform == gke {
		if creds.ClusterName == "" || creds.ClusterZone == "" || creds.GkeProject == "" {
			creds.ClusterName, creds.ClusterZone, creds.GkeProject = getGkeClusterInfo()
		}

		connectionSuccessful := false
		for !connectionSuccessful {
			readClusterName(creds)
			readClusterZone(creds)
			readGkeProject(creds)
			connectionSuccessful, _ = authenticateAtGkeCluster(*creds)
		}
	} else if *platform == openshift {
		connectionSuccessful := false
		for !connectionSuccessful {
			readOpenshiftClusterURL(creds)
			readOpenshiftUser(creds)
			readOpenshiftPassword(creds)
			connectionSuccessful, _ = authenticateAtOpenshiftCluster(*creds)
		}
	} else if *platform == aks {
		if creds.ClusterName == "" {
			creds.ClusterName = getAksClusterInfo()
		}

		connectionSuccessful := false
		for !connectionSuccessful {
			readClusterName(creds)
			readAzureResourceGroup(creds)
			readAzureSubscription(creds)
			connectionSuccessful, _ = authenticateAtAksCluster(*creds)
		}
	}
}

func readOpenshiftClusterURL(creds *installCredentials) {
	readUserInput(&creds.OpenshiftURL,
		"",
		"Openshift Server URL",
		"Please enter a valid Server URL.",
	)
}

func readOpenshiftUser(creds *installCredentials) {
	readUserInput(&creds.OpenshiftUser,
		"",
		"Openshift User",
		"Please enter a valid user name.",
	)
}

func readOpenshiftPassword(creds *installCredentials) {
	readUserInput(&creds.OpenshiftPassword,
		"",
		"Openshift Password",
		"Please enter a valid password.",
	)
}

func readClusterName(creds *installCredentials) {
	readUserInput(&creds.ClusterName,
		"^(([a-zA-Z0-9]+-)*[a-zA-Z0-9]+)$",
		"Cluster Name",
		"Please enter a valid Cluster Name.",
	)
}

func readAzureResourceGroup(creds *installCredentials) {
	readUserInput(&creds.AzureResourceGroup,
		"^(([a-zA-Z0-9]+-)*[a-zA-Z0-9]+)$",
		"Azure Resource Group",
		"Please enter a valid Azure Resource Group.",
	)
}

func readAzureSubscription(creds *installCredentials) {
	readUserInput(&creds.AzureSubscription,
		"^(([a-zA-Z0-9]+-)*[a-zA-Z0-9]+)$",
		"Azure Subscription",
		"Please enter a valid Azure Subscription.",
	)
}

func readClusterZone(creds *installCredentials) {
	readUserInput(&creds.ClusterZone,
		"^(([a-zA-Z0-9]+-)*[a-zA-Z0-9]+)$",
		"Cluster Zone",
		"Please enter a valid Cluster Zone.",
	)
}

func readGkeProject(creds *installCredentials) {
	readUserInput(&creds.GkeProject,
		"^(([a-zA-Z0-9]+-)*[a-zA-Z0-9]+)$",
		"GKE Project",
		"Please enter a valid GKE Project.",
	)
}

func readGithubUserName(creds *installCredentials) {
	readUserInput(&creds.GithubUserName,
		"^(([a-zA-Z0-9]+-)*[a-zA-Z0-9]+)$",
		"GitHub User Name",
		"Please enter a valid GitHub User Name.",
	)
}

func readGithubPersonalAccessToken(creds *installCredentials) {
	readUserInput(&creds.GithubPersonalAccessToken,
		"^[a-z0-9]{40}$",
		"GitHub Personal Access Token",
		"Please enter a valid GitHub Personal Access Token.",
	)
}

func readGithubOrg(creds *installCredentials) {
	readUserInput(&creds.GithubOrg,
		"^(([a-zA-Z0-9]+-)*[a-zA-Z0-9]+)$",
		"GitHub Organization",
		"Please enter a valid GitHub Organization.",
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

// downloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func downloadFile(filepath string, url string) error {

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

func setDeploymentFileKey(installerPath string, replacements ...placeholderReplacement) error {
	content, err := utils.ReadFile(installerPath)
	if err != nil {
		return err
	}
	for _, replacement := range replacements {
		content = strings.ReplaceAll(content, "value: "+replacement.placeholderValue, "value: "+replacement.desiredValue)
	}

	return ioutil.WriteFile(installerPath, []byte(content), 0666)
}

func authenticateAtOpenshiftCluster(creds installCredentials) (bool, error) {
	_, err := keptnutils.ExecuteCommand("oc", []string{
		"login",
		creds.OpenshiftURL,
		"-u=" + creds.OpenshiftUser,
		"-p=" + creds.OpenshiftPassword,
		"--insecure-skip-tls-verify=true",
	})

	if err != nil {
		fmt.Println("Could not connect to cluster. Please verify that you have entered the correct information.")
		return false, err
	}

	return true, nil
}

func authenticateAtGkeCluster(creds installCredentials) (bool, error) {
	_, err := keptnutils.ExecuteCommand("gcloud", []string{
		"container",
		"clusters",
		"get-credentials",
		creds.ClusterName,
		"--zone",
		creds.ClusterZone,
		"--project",
		creds.GkeProject,
	})

	if err != nil {
		fmt.Println("Could not connect to cluster. Please verify that you have entered the correct information.")
		return false, err
	}

	return true, nil
}

func authenticateAtAksCluster(creds installCredentials) (bool, error) {
	_, err := keptnutils.ExecuteCommand("az", []string{
		"aks",
		"get-credentials",
		"--resource-group",
		creds.AzureResourceGroup,
		"--name",
		creds.ClusterName,
		"--subscription",
		creds.AzureSubscription,
		"--overwrite-existing",
	})

	if err != nil {
		fmt.Println("Could not connect to cluster. Please verify that you have entered the correct information.")
		return false, err
	}

	return true, nil
}

func getGkeClusterInfo() (string, string, string) {
	// try to get current cluster from gcloud config
	out, err := getKubeContext()

	if err != nil {
		return "", "", ""
	}
	clusterInfo := strings.TrimSpace(strings.Replace(string(out), "\r\n", "\n", -1))
	if !strings.HasPrefix(clusterInfo, gke) {
		return "", "", ""
	}

	clusterInfoArray := strings.Split(clusterInfo, "_")
	if len(clusterInfoArray) < 4 {
		return "", "", ""
	}

	return clusterInfoArray[3], clusterInfoArray[2], clusterInfoArray[1]
}

func getKubeContext() (string, error) {
	return keptnutils.ExecuteCommand("kubectl", []string{
		"config",
		"current-context",
	})
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

func getAzUser() (string, error) {
	out, err := keptnutils.ExecuteCommand("az", []string{
		"account",
		"show",
		"--query=user.name",
		"--output=tsv",
	})

	if err != nil {
		return "", fmt.Errorf("Please configure your gcloud: %s", err)
	}
	// This command returns the account in the first line
	return strings.Split(strings.Replace(string(out), "\r\n", "\n", -1), "\n")[0], nil
}

func isKubectlAvailable() (bool, error) {

	_, err := keptnutils.ExecuteCommand("kubectl", []string{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func waitForInstallerPod() (string, error) {
	podName := ""
	podRunning := false
	for ok := true; ok; ok = !podRunning {
		time.Sleep(5 * time.Second)
		out, err := keptnutils.ExecuteCommand("kubectl", []string{
			"get",
			"pods",
			"-l",
			"app=installer",
			"-ojson",
		})

		if err != nil {
			return "", fmt.Errorf("Error while retrieving installer pod: %s\n. Aborting installation", err)
		}

		var podInfo map[string]interface{}
		err = json.Unmarshal([]byte(out), &podInfo)
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
	return podName, nil
}

func getInstallerLogs(podName string) error {

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
	err := execCmd.Start()
	if err != nil {
		return fmt.Errorf("Could not get installer pod logs: '%s'", err)
	}

	// cmd.Wait() should be called only after we finish reading from stdoutIn and stderrIn.
	cRes := make(chan bool)
	cErr := make(chan error)

	go func() {
		res, err := copyAndCapture(stdoutIn, "keptn-installer.log")
		cRes <- res
		cErr <- err
	}()

	installSuccessfulStdErr, errStdErr := copyAndCapture(stderrIn, "keptn-installer-err.log")
	installSuccessfulStdOut := <-cRes
	errStdOut := <-cErr

	if errStdErr != nil {
		return errStdErr
	}
	if errStdOut != nil {
		return errStdOut
	}

	err = execCmd.Wait()
	if err != nil {
		return fmt.Errorf("Could not get installer pod logs: '%s'", err)
	}

	if !installSuccessfulStdErr || !installSuccessfulStdOut {
		return errors.New("keptn installation was unsuccessful")
	}
	return nil
}

func copyAndCapture(r io.Reader, fileName string) (bool, error) {

	var file *os.File

	errorOccured := false
	installSuccessful := true
	firstRead := true

	const successMsg = "Installation of keptn complete."

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {

		if firstRead {
			// If something is read from the provided stream (stdin or stderr),
			// the data of the stream has to contain the 'successMsg'
			// for considering the keptn installation successful.
			installSuccessful = false
			firstRead = false
		}

		if file == nil {
			// Only create file on-demand
			var err error
			file, err = createFileInKeptnDirectory(fileName)
			if err != nil {
				return false, fmt.Errorf("Could not write logs into file: '%s", err)
			}
			defer file.Close()
		}
		file.WriteString(scanner.Text() + "\n")
		file.Sync()

		var reg = regexp.MustCompile(`\[keptn\|[a-zA-Z]+\]`)
		txt := scanner.Text()
		matches := reg.FindStringSubmatch(txt)
		if len(matches) == 1 {
			msgLogLevel := matches[0]
			msgLogLevel = strings.TrimPrefix(msgLogLevel, "[keptn|")
			msgLogLevel = strings.TrimSuffix(msgLogLevel, "]")
			msgLogLevel = strings.TrimSpace(msgLogLevel)
			var fullSufixReg = regexp.MustCompile(`\[keptn\|[a-zA-Z]+\]\s+\[.*\]`)
			outputStr := strings.TrimSpace(fullSufixReg.ReplaceAllString(txt, ""))

			utils.PrintLogStringLevel(outputStr, msgLogLevel)
			if utils.GetLogLevel(msgLogLevel) == utils.QuietLevel {
				errorOccured = true
			}
			if outputStr == successMsg {
				installSuccessful = true
			}
		}
	}
	return !errorOccured && installSuccessful, nil
}

func createFileInKeptnDirectory(fileName string) (*os.File, error) {
	path, err := keptnutils.GetKeptnDirectory()
	if err != nil {
		return nil, err
	}

	return os.Create(path + fileName)
}

func setupKeptnAuthAndConfigure(creds installCredentials) error {

	utils.PrintLog("Starting to configure your keptn CLI...", utils.InfoLevel)

	out, err := keptnutils.ExecuteCommand("kubectl", []string{
		"get",
		"secret",
		"keptn-api-token",
		"-n",
		"keptn",
		"-ojsonpath={.data.keptn-api-token}",
	})

	const errorMsg = `Could not retrieve keptn API token: %s
To manually set up your keptn CLI, please follow the instructions at https://keptn.sh/docs/0.2.0/reference/cli/.`

	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	apiToken, err := base64.StdEncoding.DecodeString(out)
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}
	var keptnEndpoint string
	apiEndpointRetrieved := false
	retries := 0
	for tryGetAPIEndpoint := true; tryGetAPIEndpoint; tryGetAPIEndpoint = !apiEndpointRetrieved {

		out, err := keptnutils.ExecuteCommand("kubectl", []string{
			"get",
			"virtualservice",
			"api",
			"-n",
			"keptn",
			"-ojsonpath={.spec.hosts[0]}",
		})

		if err != nil {
			retries++
			if retries >= 15 {
				utils.PrintLog("API endpoint not yet available... trying again in 5s", utils.InfoLevel)
			}
		} else {
			retries = 0
		}
		keptnEndpoint = strings.TrimSpace(string(out))
		if keptnEndpoint == "" {
			retries++
			if retries >= 15 {
				utils.PrintLog("API endpoint not yet available... trying again in 5s", utils.InfoLevel)
			}
		} else {
			keptnEndpoint = "https://" + keptnEndpoint
			apiEndpointRetrieved = true
		}
		if !apiEndpointRetrieved {
			time.Sleep(5 * time.Second)
		}
	}
	err = authenticate(keptnEndpoint, string(apiToken))
	if err != nil {
		return err
	}
	err = configure(creds)
	if err != nil {
		return err
	}
	utils.PrintLog("Your CLI is now successfully configured. You are now ready to use keptn.", utils.InfoLevel)
	return nil
}

func authenticate(endPoint string, apiToken string) error {
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"auth",
		fmt.Sprintf("--endpoint=%s", endPoint),
		fmt.Sprintf("--api-token=%s", apiToken),
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err != nil {
		return fmt.Errorf("Authentication at keptn failed: %s\n"+
			"To manually set up your keptn CLI, please follow the instructions at https://keptn.sh/docs/0.2.0/reference/cli/.", err)
	}
	return nil
}

func configure(creds installCredentials) error {

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"configure",
		fmt.Sprintf("--org=%s", creds.GithubOrg),
		fmt.Sprintf("--user=%s", creds.GithubUserName),
		fmt.Sprintf("--token=%s", creds.GithubPersonalAccessToken),
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err != nil {
		return fmt.Errorf("Configuration failed: %s\n"+
			"To manually set up your keptn CLI, please follow the instructions at https://keptn.sh/docs/0.2.0/reference/cli/.", err)
	}
	return nil
}
