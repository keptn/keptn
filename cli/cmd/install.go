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
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	kubeclient "helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"io"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/keptn/keptn/cli/pkg/file"
	"github.com/keptn/keptn/cli/pkg/kube"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

type installCmdParams struct {
	ConfigFilePath      *string
	InstallerImage      *string
	KeptnVersion        *string
	PlatformIdentifier  *string
	UseCaseInput        *string
	UseCase             usecase
	ApiServiceTypeInput *string
	ApiServiceType      apiServiceType
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

		// Parse api service type
		if val, ok := apiServiceTypeToID[*installParams.ApiServiceTypeInput]; ok {
			installParams.ApiServiceType = val
		} else {
			return errors.New("value of 'keptn-api-service-type' flag is unknown. Supported values are " +
				"[" + ClusterIP.String() + "," + LoadBalancer.String() + "," + NodePort.String() + "]")
		}

		if insecureSkipTLSVerify {
			kubectlOptions = "--insecure-skip-tls-verify=true"
		}

		err := setPlatform()
		if err != nil {
			return err
		}

		if val, ok := usecaseToID[*installParams.UseCaseInput]; ok {
			installParams.UseCase = val
		} else {
			return errors.New("value of 'use-case' flag is unknown. Supported values are " +
				"[" + ContinuousDelivery.String() + "," + QualityGates.String() + "]")
		}

		// Mark the quality-gates use case as deprecated - this is now the default option
		if installParams.UseCase == QualityGates {
			logging.PrintLog("NOTE: The --use-case=quality-gates option is now deprecated and is now a synonym for the default installation of Keptn.", logging.InfoLevel)
		}

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

	installParams.UseCaseInput = installCmd.Flags().StringP("use-case", "u", "",
		"The use case to install Keptn for ["+ContinuousDelivery.String()+","+QualityGates.String()+"]")

	installParams.ApiServiceTypeInput = installCmd.Flags().StringP("keptn-api-service-type", "",
		ClusterIP.String(), "Installation options for the api-service type ["+ClusterIP.String()+","+
			LoadBalancer.String()+","+NodePort.String()+"]")

	installCmd.PersistentFlags().BoolVarP(&insecureSkipTLSVerify, "insecure-skip-tls-verify", "s",
		false, "Skip tls verification for kubectl commands")
}

func prepareInstallerManifest() (installerManifest string) {

	installerManifest = installerJob

	installerManifest = strings.ReplaceAll(installerManifest, "INSTALLER_IMAGE_PLACEHOLDER",
		*installParams.InstallerImage)

	installerManifest = strings.ReplaceAll(installerManifest, "PLATFORM_PLACEHOLDER", strings.ToLower(*installParams.PlatformIdentifier))
	installerManifest = strings.ReplaceAll(installerManifest, "API_SERVICE_TYPE_PLACEHOLDER", installParams.ApiServiceType.String())
	installerManifest = strings.ReplaceAll(installerManifest, "USE_CASE_PLACEHOLDER", installParams.UseCase.String())
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

func newActionConfig(config *rest.Config, namespace string) (*action.Configuration, error) {

	logFunc := func(format string, v ...interface{}) {
		fmt.Sprintf(format, v)
	}

	restClientGetter := newConfigFlags(config, namespace)
	kubeClient := &kubeclient.Client{
		Factory: cmdutil.NewFactory(restClientGetter),
		Log:     logFunc,
	}
	client, err := kubeClient.Factory.KubernetesClientSet()
	if err != nil {
		return nil, err
	}

	s := driver.NewSecrets(client.CoreV1().Secrets(namespace))
	s.Log = logFunc

	return &action.Configuration{
		RESTClientGetter: restClientGetter,
		Releases:         storage.Init(s),
		KubeClient:       kubeClient,
		Log:              logFunc,
	}, nil
}

func newConfigFlags(config *rest.Config, namespace string) *genericclioptions.ConfigFlags {
	return &genericclioptions.ConfigFlags{
		Namespace:   &namespace,
		APIServer:   &config.Host,
		CAFile:      &config.CAFile,
		BearerToken: &config.BearerToken,
	}
}

func upgradeChart(ch *chart.Chart, releaseName, namespace string, vals map[string]interface{}) error {

	if len(ch.Templates) > 0 {
		logging.PrintLog(fmt.Sprintf("Start upgrading chart %s in namespace %s", releaseName, namespace), logging.InfoLevel)
		var kubeconfig string
		if os.Getenv("KUBECONFIG") != "" {
			kubeconfig = keptnutils.ExpandTilde(os.Getenv("KUBECONFIG"))
		} else {
			kubeconfig = filepath.Join(
				keptnutils.UserHomeDir(), ".kube", "config",
			)
		}
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return err
		}

		cfg, err := newActionConfig(config, namespace)
		if err != nil {
			return err
		}

		histClient := action.NewHistory(cfg)
		var release *release.Release

		if _, err = histClient.Run(releaseName); err == driver.ErrReleaseNotFound {
			iCli := action.NewInstall(cfg)
			iCli.Namespace = namespace
			iCli.ReleaseName = releaseName
			iCli.Wait = true
			release, err = iCli.Run(ch, vals)
		} else {
			iCli := action.NewUpgrade(cfg)
			iCli.Namespace = namespace
			iCli.Wait = true
			iCli.ResetValues = true
			release, err = iCli.Run(releaseName, ch, vals)
		}
		if err != nil {
			return fmt.Errorf("Error when installing/upgrading chart %s in namespace %s: %s",
				releaseName, namespace, err.Error())
		}
		if release != nil {
			logging.PrintLog(release.Manifest, logging.VerboseLevel)
			if err := waitForDeploymentsOfHelmRelease(release.Manifest); err != nil {
				return err
			}
		} else {
			logging.PrintLog("Release is nil", logging.InfoLevel)
		}
		logging.PrintLog(fmt.Sprintf("Finished upgrading chart %s in namespace %s", releaseName, namespace), logging.InfoLevel)
	} else {
		logging.PrintLog("Upgrade not done as this is an empty chart", logging.InfoLevel)
	}
	return nil
}

func getDeployments(helmManifest string) []*appsv1.Deployment {

	deployments := []*appsv1.Deployment{}
	dec := kyaml.NewYAMLToJSONDecoder(strings.NewReader(helmManifest))
	for {
		var dpl appsv1.Deployment
		err := dec.Decode(&dpl)
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		if keptnutils.IsDeployment(&dpl) {
			deployments = append(deployments, &dpl)
		}
	}
	return deployments
}

func waitForDeploymentsOfHelmRelease(helmManifest string) error {
	depls := getDeployments(helmManifest)
	for _, depl := range depls {
		if err := keptnutils.WaitForDeploymentToBeRolledOut(false, depl.Name, depl.Namespace); err != nil {
			return fmt.Errorf("Error when waiting for deployment %s in namespace %s: %s", depl.Name, depl.Namespace, err.Error())
		}
	}
	return nil
}

// Preconditions: 1. Already authenticated against the cluster.
func doInstallation() error {
	url := "https://storage.googleapis.com/keptn-installer/latest/keptn-0.1.0.tgz"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	ch, err := keptnutils.LoadChart(bytes)
	if err != nil {
		return err
	}

	values := map[string]interface{}{
		"continuous-delivery": map[string]interface{}{
			"enabled": "true",
		},
	}

	return upgradeChart(ch, "keptn", "keptn", values)
	///////
	/*
		if err := createKeptnNamespace("keptn", "keptn-datastore"); err != nil {
			return err
		}

		rbacSet := make(map[string]bool)
		rbacSet[installerServiceAccount] = true

		rbacSet[getAdminRoleForNamespace("keptn")] = true
		rbacSet[getAdminRoleForNamespace("keptn-datastore")] = true

		// Add RBACs for NATS
		rbacSet[natsClusterRole] = true
		rbacSet[natsOperatorRoles] = true
		rbacSet[natsOperatorServer] = true

		err := applyRbac(rbacSet)
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

	*/
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
