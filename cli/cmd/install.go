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
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/keptn/keptn/cli/pkg/docker"
	"github.com/keptn/keptn/cli/pkg/file"
	"github.com/keptn/keptn/cli/pkg/kube"

	"github.com/keptn/keptn/cli/pkg/version"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/spf13/cobra"
)

type installCmdParams struct {
	ConfigFilePath            *string
	InstallerImage            *string
	KeptnVersion              *string
	PlatformIdentifier        *string
	GatewayInput              *string
	Gateway                   gateway
	Domain                    *string
	UseCaseInput              *string
	UseCase                   usecase
	IngressInstallOptionInput *string
	IngressInstallOption      ingressInstallOption
}

const keptnInstallerLogFileName = "keptn-installer.log"
const keptnInstallerErrorLogFileName = "keptn-installer-Err.log"

// KubeServerVersionConstraints the Kubernetes Cluster version's constraints is passed by ldflags
var KubeServerVersionConstraints string

var installParams installCmdParams

const openshift = "openshift"
const kubernetes = "kubernetes"

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
	Long: `The Keptn CLI allows installing Keptn on any Kubernetes derivate to which your kube config is pointing to, and on OpenShift.

For more information, please consult the following docs:

* https://keptn.sh/docs/develop/installation/setup-keptn/

`,
	Example: `keptn install # install on Kubernetes
keptn install --platform=openshift --use-case=continuous-delivery # install continuous-delivery on Openshift
keptn install --platform=kubernetes --gateway=NodePort # install on a Kubernetes instance with gateway NodePort`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {

		if installParams.KeptnVersion != nil && *installParams.KeptnVersion != "" {
			return errors.New("The flag --keptn-version is not supported anymore but you can specify " +
				"an image for the installer using the flag 'keptn-installer-image'")
		}

		// Parse IngressInstallOption
		if val, ok := ingressInstallOptionToID[*installParams.IngressInstallOptionInput]; ok {
			installParams.IngressInstallOption = val
		} else {
			return errors.New("value of 'ingress-install-option' flag is unknown. Supported values are " +
				"[" + StopIfInstalled.String() + "," + Reuse.String() + "," + Overwrite.String() + "]")
		}

		if insecureSkipTLSVerify {
			kubectlOptions = "--insecure-skip-tls-verify=true"
		}

		err := setPlatform()
		if err != nil {
			return err
		}

		if val, ok := gatewayToID[*installParams.GatewayInput]; ok {
			installParams.Gateway = val
		} else {
			return errors.New("value of 'gateway' flag is unknown. Supported values are " +
				"[" + NodePort.String() + "," + LoadBalancer.String() + "]")
		}

		if val, ok := usecaseToID[*installParams.UseCaseInput]; ok {
			installParams.UseCase = val
		} else {
			return errors.New("value of 'use-case' flag is unknown. Supported values are " +
				"[" + AllUseCases.String() + "," + QualityGates.String() + "]")
		}

		// Mark the quality-gates use case as deprecated - this is now the default option
		if *installParams.UseCaseInput == "quality-gates" {
			logging.PrintLog("NOTE: The --use-case=quality-gates option is now deprecated and is now a synonym for the default installation of Keptn.", logging.InfoLevel)
		}

		var image, tag string
		// Determine installer version
		if installParams.InstallerImage != nil && *installParams.InstallerImage != "" {
			image, tag = docker.SplitImageName(*installParams.InstallerImage)
		} else if version.IsOfficialKeptnVersion(Version) {
			image = "docker.io/keptn/installer"
			tag, err = version.GetOfficialKeptnVersion(Version)
			if err != nil {
				return fmt.Errorf("Error when parsing installer tag: %v", err)
			}
			*installParams.InstallerImage = image + ":" + tag
		} else {
			image = "docker.io/keptn/installer"
			tag = "latest"
			*installParams.InstallerImage = image + ":" + tag
		}

		err = docker.CheckImageAvailability(image, tag, nil)
		if err != nil {
			return fmt.Errorf("Installer image not found under: %v", err)
		}

		logging.PrintLog(fmt.Sprintf("Used Installer version: %s",
			*installParams.InstallerImage), logging.InfoLevel)

		if !mocking {
			if p.checkRequirements() != nil {
				return err
			}

			// Check whether kubectl is installed
			isKubAvailable, err := kube.IsKubectlAvailable()
			if err != nil || !isKubAvailable {
				return errors.New(`Keptn requires 'kubectl' but it is not available.
Please see https://kubernetes.io/docs/tasks/tools/install-kubectl/`)
			}
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

			// check if Kubernetes server version is compatible (except OpenShift)
			if *installParams.PlatformIdentifier != "openshift" {
				if err := kube.CheckKubeServerVersion(KubeServerVersionConstraints); err != nil {
					logging.PrintLog(err.Error(), logging.VerboseLevel)
					logging.PrintLog("See https://keptn.sh/docs/0.7.0/installation/k8s-support/ for details.", logging.VerboseLevel)
					return errors.New(`Keptn requires Kubernetes server version: ` + KubeServerVersionConstraints)
				}
			}

			return doInstallation()
		}
		fmt.Println("Skipping installation due to mocking flag")
		return nil
	},
}

func setPlatform() error {

	*installParams.PlatformIdentifier = strings.ToLower(*installParams.PlatformIdentifier)

	switch *installParams.PlatformIdentifier {
	case openshift:
		p = newOpenShiftPlatform()
		return nil
	case kubernetes:
		p = newKubernetesPlatform()
		return nil
	default:
		return errors.New("Unsupported platform '" + *installParams.PlatformIdentifier +
			"'. The following platforms are supported: openshift and kubernetes")
	}
}

func init() {
	rootCmd.AddCommand(installCmd)
	installParams = installCmdParams{}

	installParams.PlatformIdentifier = installCmd.Flags().StringP("platform", "p", "kubernetes",
		"The platform to run keptn on [kubernetes,openshift]")

	installParams.ConfigFilePath = installCmd.Flags().StringP("creds", "c", "",
		"Specify a JSON file containing cluster information needed for the installation (this allows skipping user prompts to execute a *silent* Keptn installation)")

	installParams.InstallerImage = installCmd.Flags().StringP("keptn-installer-image", "k",
		"", "The installer image which is used for the installation")
	installCmd.Flags().MarkHidden("keptn-installer-image")

	installParams.KeptnVersion = installCmd.Flags().StringP("keptn-version", "",
		"", "This flag is not supported anymore but you can specify an image for the installer"+
			"using the flag 'keptn-installer-image'")
	installCmd.Flags().MarkHidden("keptn-version")

	installParams.GatewayInput = installCmd.Flags().StringP("gateway", "g", LoadBalancer.String(),
		"The ingress-loadbalancer type ["+LoadBalancer.String()+","+NodePort.String()+"]")

	installParams.Domain = installCmd.Flags().StringP("domain", "d", "",
		"Experimental: Overwrite the ingress domain (e.g., 127.0.0.1.xip.io)")
	installCmd.Flags().MarkHidden("gateway-domain")

	installParams.UseCaseInput = installCmd.Flags().StringP("use-case", "u", "",
		"The use case to install Keptn for ["+AllUseCases.String()+","+QualityGates.String()+"]")

	installParams.IngressInstallOptionInput = installCmd.Flags().StringP("ingress-install-option", "",
		StopIfInstalled.String(), "Installation options for Ingress ["+StopIfInstalled.String()+","+
			Reuse.String()+","+Overwrite.String()+"]")
	installCmd.Flags().MarkHidden("ingress-install-option")

	installCmd.PersistentFlags().BoolVarP(&insecureSkipTLSVerify, "insecure-skip-tls-verify", "s",
		false, "Skip tls verification for kubectl commands")
}

func prepareInstallerManifest() (installerManifest string) {

	installerManifest = installerJob

	installerManifest = strings.ReplaceAll(installerManifest, "INSTALLER_IMAGE_PLACEHOLDER",
		*installParams.InstallerImage)

	installerManifest = strings.ReplaceAll(installerManifest, "PLATFORM_PLACEHOLDER", strings.ToLower(*installParams.PlatformIdentifier))
	installerManifest = strings.ReplaceAll(installerManifest, "GATEWAY_TYPE_PLACEHOLDER", installParams.Gateway.String())
	installerManifest = strings.ReplaceAll(installerManifest, "DOMAIN_PLACEHOLDER", *installParams.Domain)
	installerManifest = strings.ReplaceAll(installerManifest, "USE_CASE_PLACEHOLDER", installParams.UseCase.String())
	installerManifest = strings.ReplaceAll(installerManifest, "INGRESS_INSTALL_OPTION_PLACEHOLDER", installParams.IngressInstallOption.String())

	installerManifest = strings.ReplaceAll(installerManifest, "INGRESS_PLACEHOLDER", getIngress().String())
	return
}

func getIngress() (ingress Ingress) {
	ingress = istio
	if installParams.UseCase == QualityGates {
		ingress = nginx
	}
	return
}

func createKeptnNamespace(keptnNamespace, keptnDatastoreNamespace string) error {

	for _, ns := range [...]string{keptnNamespace, keptnDatastoreNamespace} {
		res, err := keptnutils.ExistsNamespace(false, ns)
		if err != nil {
			return fmt.Errorf("Failed to check if namespace %s already exists: %v", ns, err)
		}
		if res {
			return fmt.Errorf("Existing Keptn installation found in namespace %s.", ns)
		}
	}
	for _, ns := range [...]string{keptnNamespace, keptnDatastoreNamespace} {
		err := keptnutils.CreateNamespace(false, ns)
		if err != nil {
			return fmt.Errorf("Failed to create Keptn namespace %s: %v", ns, err)
		}
	}
	return nil
}

func prepareIngressControllerInstall(ingressControllerNamespace string) (string, error) {

	ingressControllerAvailable, err := keptnutils.ExistsNamespace(false, ingressControllerNamespace)
	if err != nil {
		return "", fmt.Errorf("Failed to check if namespace %s already exists: %v", ingressControllerNamespace, err)
	}

	if ingressControllerAvailable && installParams.IngressInstallOption == Reuse {
		fmt.Printf("Ingress controller in namespace %s is reused but its compatibility is not checked\n", ingressControllerNamespace)
		return getGetRoleForNamespace(ingressControllerNamespace), nil
	} else if ingressControllerAvailable && installParams.IngressInstallOption == StopIfInstalled {
		return "", fmt.Errorf("Namespace %s is already available but its Ingress controller is not used due to unknown compatibility",
			ingressControllerNamespace)
	} else {
		if !ingressControllerAvailable {
			err := keptnutils.CreateNamespace(false, ingressControllerNamespace)
			if err != nil {
				return "", fmt.Errorf("Failed to create namespace %s: %v", ingressControllerNamespace, err)
			}
		}
		return installerClusterAdminClusterRoleBinding, nil
	}
}

// Preconditions: 1. Already authenticated against the cluster.
func doInstallation() error {

	if err := createKeptnNamespace("keptn", "keptn-datastore"); err != nil {
		return err
	}

	rbacSet := make(map[string]bool)
	rbacSet[installerServiceAccount] = true

	// Prepare Ingress
	var ingressRBAC string
	var err error
	if ingressRBAC, err = prepareIngressControllerInstall(getIngress().getDefaultNamespace()); err != nil {
		return err
	}
	rbacSet[ingressRBAC] = true

	if getIngress() == istio {
		// Full installation: cluster-admin rights are needed for helm-service
		rbacSet[installerClusterAdminClusterRoleBinding] = true
	} else {
		// QG only installation
		rbacSet[getAdminRoleForNamespace("keptn")] = true
		rbacSet[getAdminRoleForNamespace("keptn-datastore")] = true
	}

	// Add RBACs for NATS
	rbacSet[natsClusterRole] = true
	rbacSet[natsOperatorRoles] = true
	rbacSet[natsOperatorServer] = true

	err = applyRbac(rbacSet)
	defer deleteRbac(rbacSet)
	if err != nil {
		return err
	}

	logging.PrintLog("Deploying Keptn installer job ...", logging.InfoLevel)

	o := options{"apply", "-f", "-"}
	o.appendIfNotEmpty(kubectlOptions)
	logging.PrintLog("Executing: kubectl "+strings.Join(o, " "), logging.VerboseLevel)
	cmd := exec.Command("kubectl", o...)
	cmd.Stdin = strings.NewReader(prepareInstallerManifest())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error while applying installer job: %s \n%s\nAborting installation", err.Error(), string(out))
	}
	logging.PrintLog("Result: "+string(out), logging.VerboseLevel)

	logging.PrintLog("Waiting for installer pod to be started ...", logging.InfoLevel)
	installerPodName, err := waitForInstallerPod()
	if err != nil {
		return err
	}

	logging.PrintLog("Installer pod started successfully.", logging.InfoLevel)

	if err := getInstallerLogs(installerPodName); err != nil {
		return err
	}

	o = options{"delete", "job", "installer", "-n", "keptn"}
	o.appendIfNotEmpty(kubectlOptions)
	logging.PrintLog("Executing: kubectl "+strings.Join(o, " "), logging.VerboseLevel)
	result, err := keptnutils.ExecuteCommand("kubectl", o)
	if err != nil {
		logging.PrintLog("Error deleting installer job : kubectl "+err.Error(), logging.QuietLevel)
		return err
	}
	logging.PrintLog("Result: "+result, logging.VerboseLevel)

	fmt.Println("Trying to get auth-token and API endpoint from Kubernetes.")
	// installation finished, get auth token and endpoint
	if err := authUsingKube(); err != nil {
		return err
	}

	return nil
}

func applyRbac(rbac map[string]bool) error {
	logging.PrintLog("Applying RBAC rules for installer", logging.VerboseLevel)
	var merged string
	for k := range rbac {
		merged += k
	}
	o := options{"apply", "-f", "-"}
	o.appendIfNotEmpty(kubectlOptions)
	logging.PrintLog("Executing: kubectl "+strings.Join(o, " "), logging.VerboseLevel)
	cmd := exec.Command("kubectl", o...)
	cmd.Stdin = strings.NewReader(merged)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error while applying RBAC: %s \n%s\nAborting installation", err.Error(), string(out))
	}
	logging.PrintLog("Result: "+string(out), logging.VerboseLevel)
	return nil
}

func deleteRbac(rbac map[string]bool) error {
	logging.PrintLog("Deleting RBAC rules for installer", logging.VerboseLevel)
	var merged string
	for k := range rbac {
		merged += k
	}
	o := options{"delete", "-f", "-"}
	o.appendIfNotEmpty(kubectlOptions)

	logging.PrintLog("Executing: kubectl "+strings.Join(o, " "), logging.VerboseLevel)
	cmd := exec.Command("kubectl", o...)
	cmd.Stdin = strings.NewReader(merged)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error while deleting RBAC: %s \n%s\nAborting installation", err.Error(), string(out))
	}
	logging.PrintLog("Result: "+string(out), logging.VerboseLevel)
	return nil
}

func parseConfig(configFile string) error {
	data, err := file.ReadFile(configFile)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), p.getCreds())
}

func readCreds() error {

	cm := credentialmanager.NewCredentialManager()
	credsStr, err := cm.GetInstallCreds()
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
	return cm.SetInstallCreds(newCredsStr)
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
	for true {
		options := options{"get",
			"pods",
			"-l",
			"app=installer",
			"-n",
			"keptn",
			"-ojson"}
		options.appendIfNotEmpty(kubectlOptions)
		logging.PrintLog("Executing: kubectl "+strings.Join(options, " "), logging.VerboseLevel)
		out, err := keptnutils.ExecuteCommand("kubectl", options)

		if err != nil {
			return "", fmt.Errorf("Error while retrieving installer pod: %s\n. Aborting installation", err)
		}
		logging.PrintLog("Result: "+out, logging.VerboseLevel)

		var podInfo map[string]interface{}
		err = json.Unmarshal([]byte(out), &podInfo)
		if err == nil && podInfo != nil {
			podStatusArray := podInfo["items"].([]interface{})

			if len(podStatusArray) > 0 {
				podStatus := podStatusArray[0].(map[string]interface{})["status"].(map[string]interface{})["phase"].(string)
				if podStatus == "Running" {
					return podStatusArray[0].(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string), nil
				} else if podStatus == "ImagePullBackOff" {
					return "", errors.New("Installer pod could not be deployed (Status: ImagePullBackOff). " +
						"Please verify that your Kubernetes cluster has a working Internet connection.")
				} else if podStatus == "Failed" {
					return "", errors.New("Installer pod ran into failure. " +
						"Please check if a failed installer job exits: \"kubectl get jobs installer -n keptn\"." +
						"Please also check the log output: \"kubectl logs jobs/installer -n keptn\".")
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
	return "", nil
}

func getInstallerLogs(podName string) error {

	fmt.Printf("Getting logs of pod %s\n", podName)

	options := options{"logs",
		podName,
		"-c",
		"keptn-installer",
		"-n",
		"keptn",
		"-f"}
	options.appendIfNotEmpty(kubectlOptions)

	logging.PrintLog("Executing: kubectl "+strings.Join(options, " "), logging.VerboseLevel)
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
		res, err := copyAndCapture(stdoutIn, keptnInstallerLogFileName)
		cRes <- res
		cErr <- err
	}()

	installSuccessfulStdErr, errStdErr := copyAndCapture(stderrIn, keptnInstallerErrorLogFileName)
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
