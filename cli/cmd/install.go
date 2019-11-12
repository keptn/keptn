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
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/spf13/cobra"
)

type installCmdParams struct {
	ConfigFilePath     *string
	InstallerVersion   *string
	PlatformIdentifier *string
	GatewayType        *string
	UseCase            *string
}

var installParams *installCmdParams

const gke = "gke"
const aks = "aks"
const eks = "eks"
const pks = "pks"
const openshift = "openshift"
const kubernetes = "kubernetes"

const installerPrefixURL = "https://raw.githubusercontent.com/keptn/keptn/"
const installerSuffixPath = "/installer/manifests/installer/installer.yaml"
const rbacSuffixPath = "/installer/manifests/installer/rbac.yaml"

type platform interface {
	checkRequirements() error
	getCreds() interface{}
	checkCreds() error
	readCreds()
	printCreds()
}

var p platform

// installCmd represents the version command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs Keptn on a Kubernetes cluster",
	Long: `Installs Keptn on a Kubernetes cluster

Example:
	keptn install`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {

		if insecureSkipTLSVerify {
			kubectlOptions = "--insecure-skip-tls-verify=true"
		}

		err := setPlatform()
		if err != nil {
			return err
		}

		if !checkIfGatewayTypeIsSupported() {
			return errors.New(`Keptn currently supports: 'LoadBalancer' and 'NodePort'`)
		}

		if !checkIfUseCaseIsSUpported() {
			return errors.New(`Keptn currently supports use case: 'quality-gates' and 'all'`)
		}

		isInstallerAvailable, err := checkInstallerAvailability()
		if err != nil || !isInstallerAvailable {
			return errors.New("Installers not found under:\n" +
				getInstallerURL() + "\n" + getRbacURL())
		}

		if p.checkRequirements() != nil {
			return err
		}

		// Check whether kubectl is installed
		isKubAvailable, err := utils.IsKubectlAvailable()
		if err != nil || !isKubAvailable {
			return errors.New(`Keptn requires 'kubectl' but it is not available.
Please see https://kubernetes.io/docs/tasks/tools/install-kubectl/`)
		}

		if installParams.ConfigFilePath != nil && *installParams.ConfigFilePath != "" {
			// Config was provided in form of a file
			err = parseConfig(*installParams.ConfigFilePath)
			if err != nil {
				return err
			}

			// Check whether the authentication at the cluster is valid
			err = p.checkCreds()
			if err != nil {
				return err
			}
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		logging.PrintLog("Installing Keptn ...", logging.InfoLevel)

		var err error
		if !mocking {
			if installParams.ConfigFilePath == nil || *installParams.ConfigFilePath == "" {
				err = readCreds()
				if err != nil {
					return err
				}
			}
			return doInstallation()
		}
		fmt.Println("Skipping installation due to mocking flag")
		return nil
	},
}

func checkIfGatewayTypeIsSupported() bool {
	return *installParams.GatewayType == "NodePort" || *installParams.GatewayType == "LoadBalancer"
}

func checkIfUseCaseIsSUpported() bool {
	return *installParams.UseCase == "quality-gates" || *installParams.GatewayType == "all"
}

func setPlatform() error {

	*installParams.PlatformIdentifier = strings.ToLower(*installParams.PlatformIdentifier)

	switch *installParams.PlatformIdentifier {
	case gke:
		p = newGKEPlatform()
		return nil
	case aks:
		p = newAKSPlatform()
		return nil
	case eks:
		p = newEKSPlatform()
		return nil
	case pks:
		p = newPKSPlatform()
		return nil
	case openshift:
		p = newOpenShiftPlatform()
		return nil
	case kubernetes:
		p = newKubernetesPlatform()
		return nil
	default:
		return errors.New("Unsupported platform '" + *installParams.PlatformIdentifier +
			"'. The following platforms are supported: aks, eks, gke, pks, openshift, and kubernetes")
	}
}

func init() {
	rootCmd.AddCommand(installCmd)
	installParams = &installCmdParams{}

	installParams.PlatformIdentifier = installCmd.Flags().StringP("platform", "p", "gke",
		"The platform to run keptn on [aks,eks,gke,pks,openshift,kubernetes]")

	installParams.ConfigFilePath = installCmd.Flags().StringP("creds", "c", "",
		"The name of the creds file")
	installCmd.Flags().MarkHidden("creds")

	installParams.InstallerVersion = installCmd.Flags().StringP("keptn-version", "k",
		"master", "The branch or tag of the version which is installed")
	installCmd.Flags().MarkHidden("keptn-version")

	installParams.GatewayType = installCmd.Flags().StringP("gateway", "g", "LoadBalancer",
		"The ingress-loadbalancer type [LoadBalancer,NodePort]")
	installCmd.Flags().MarkHidden("gateway")

	installParams.UseCase = installCmd.Flags().StringP("use-case", "u", "all",
		"The use case to install Keptn for [quality-gates,all]")
	installCmd.Flags().MarkHidden("use-case")

	installCmd.PersistentFlags().BoolVarP(&insecureSkipTLSVerify, "insecure-skip-tls-verify", "s",
		false, "Skip tls verification for kubectl commands")
}

func checkInstallerAvailability() (bool, error) {

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
	return installerPrefixURL + *installParams.InstallerVersion + installerSuffixPath
}

func getRbacURL() string {
	return installerPrefixURL + *installParams.InstallerVersion + rbacSuffixPath
}

// Preconditions: 1. Already authenticated against the cluster.
func doInstallation() error {

	path, err := keptnutils.GetKeptnDirectory()
	if err != nil {
		return err
	}
	installerPath := path + "installer.yaml"

	// get the YAML for the installer pod
	if err := utils.DownloadFile(installerPath, getInstallerURL()); err != nil {
		return err
	}

	if err := utils.Replace(installerPath,
		utils.PlaceholderReplacement{PlaceholderValue: "PLATFORM_PLACEHOLDER",
			DesiredValue: *installParams.PlatformIdentifier}); err != nil {
		return err
	}
	if err := utils.Replace(installerPath,
		utils.PlaceholderReplacement{PlaceholderValue: "GATEWAY_TYPE_PLACEHOLDER",
			DesiredValue: *installParams.GatewayType}); err != nil {
		return err
	}

	// use case specific Keptn installation
	ingress := "istio"
	if *installParams.UseCase == "quality-gates" {
		ingress = "ngnix"
	}
	if err := utils.Replace(installerPath,
		utils.PlaceholderReplacement{PlaceholderValue: "INGRESS_PLACEHOLDER",
			DesiredValue: ingress}); err != nil {
		return err
	}

	_, aks := p.(*aksPlatform)
	_, eks := p.(*eksPlatform)
	_, gke := p.(*gkePlatform)
	_, pks := p.(*pksPlatform)
	_, k8s := p.(*kubernetesPlatform)
	if gke || aks || k8s || eks || pks {
		options := options{"apply", "-f", getRbacURL()}
		options.appendIfNotEmpty(kubectlOptions)
		_, err = keptnutils.ExecuteCommand("kubectl", options)

		if err != nil {
			return fmt.Errorf("Error while applying RBAC for installer pod: %s \nAborting installation", err.Error())
		}
	}

	logging.PrintLog("Deploying Keptn installer pod ...", logging.InfoLevel)

	o := options{"apply", "-f", installerPath}
	o.appendIfNotEmpty(kubectlOptions)
	if _, err := keptnutils.ExecuteCommand("kubectl", o); err != nil {
		return fmt.Errorf("Error while deploying keptn installer pod: %s \nAborting installation", err.Error())
	}

	logging.PrintLog("Installer pod deployed successfully.", logging.InfoLevel)

	installerPodName, err := waitForInstallerPod()
	if err != nil {
		return err
	}

	if err := getInstallerLogs(installerPodName); err != nil {
		return err
	}

	if err := os.Remove(installerPath); err != nil {
		return err
	}

	// use case specific Keptn modification
	if *installParams.UseCase == "quality-gates" {
		o = options{"delete", "deployment", "gatekeeper-service-evaluation-done-distributor", "-n", "keptn"}
		o.appendIfNotEmpty(kubectlOptions)
		_, err = keptnutils.ExecuteCommand("kubectl", o)
		if err != nil {
			return err
		}
	}

	o = options{"delete", "job", "installer", "-n", "default"}
	o.appendIfNotEmpty(kubectlOptions)
	_, err = keptnutils.ExecuteCommand("kubectl", o)
	if err != nil {
		return err
	}

	if eks {
		o = options{"get", "svc", "istio-ingressgateway", "-n", "istio-system",
			"-ojsonpath={.status.loadBalancer.ingress[0].hostname}"}
		o.appendIfNotEmpty(kubectlOptions)
		hostname, err := keptnutils.ExecuteCommand("kubectl", o)
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println("Please create a Route53 Hosted Zone with a wildcard record set for " + hostname)
		fmt.Println("Afterwards, call 'keptn configure domain YOUR_ROUTE53_DOMAIN'")
	} else {
		// installation finished, get auth token and endpoint
		if err := authUsingKube(); err != nil {
			return err
		}
	}
	return nil
}

func parseConfig(configFile string) error {
	data, err := utils.ReadFile(configFile)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), p.getCreds())
}

func readCreds() error {

	credsStr, err := credentialmanager.GetInstallCreds()
	if err != nil {
		credsStr = ""
	}
	// Ignore unmarshaling error
	json.Unmarshal([]byte(credsStr), p.getCreds())

	fmt.Print("Please enter the following information or press enter to keep the old value:\n")

	for {
		p.readCreds()

		fmt.Println()
		fmt.Println("Please confirm that the provided information is correct: ")

		p.printCreds()

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

	newCreds, _ := json.Marshal(p.getCreds())
	newCredsStr := strings.Replace(string(newCreds), "\r\n", "\n", -1)
	newCredsStr = strings.Replace(newCredsStr, "\n", "", -1)
	return credentialmanager.SetInstallCreds(newCredsStr)
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

func waitForInstallerPod() (string, error) {
	podName := ""
	podRunning := false
	for ok := true; ok; ok = !podRunning {
		time.Sleep(5 * time.Second)

		options := options{"get",
			"pods",
			"-l",
			"app=installer",
			"-n",
			"default",
			"-ojson"}
		options.appendIfNotEmpty(kubectlOptions)
		out, err := keptnutils.ExecuteCommand("kubectl", options)

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

	options := options{"logs",
		podName,
		"-c",
		"keptn-installer",
		"-n",
		"default",
		"-f"}
	options.appendIfNotEmpty(kubectlOptions)

	execCmd := exec.Command(
		"kubectl", options...,
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
		return errors.New("Keptn installation was unsuccessful")
	}
	return nil
}

func copyAndCapture(r io.Reader, fileName string) (bool, error) {

	var file *os.File

	errorOccured := false
	installSuccessful := true
	firstRead := true

	const successMsg = "Installation of Keptn complete."

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

			logging.PrintLogStringLevel(outputStr, msgLogLevel)
			if logging.GetLogLevel(msgLogLevel) == logging.QuietLevel {
				logging.PrintLog(fmt.Sprintf("Error occured with message: %s", txt), logging.InfoLevel)
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
