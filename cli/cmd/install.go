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
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/keptn/keptn/cli/pkg/version"

	"github.com/keptn/keptn/cli/pkg/file"
	"github.com/keptn/keptn/cli/pkg/kube"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	kubeclient "helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

type installCmdParams struct {
	ConfigFilePath      *string
	KeptnVersion        *string
	PlatformIdentifier  *string
	UseCaseInput        *string
	UseCase             usecase
	ApiServiceTypeInput *string
	ApiServiceType      apiServiceType
	ChartRepoURL        *string
}

const keptnInstallerLogFileName = "keptn-installer.log"
const keptnInstallerErrorLogFileName = "keptn-installer-Err.log"

const keptnInstallerHelmRepoURL = "https://storage.googleapis.com/keptn-installer/"

const keptnReleaseDocsURL = "0.7.x"

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

var installChart *chart.Chart

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
keptn install --platform=kubernetes --keptn-api-service-type=NodePort # install on a Kubernetes instance with gateway NodePort`,
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
		if *installParams.UseCaseInput == QualityGates.String() {
			logging.PrintLog("NOTE: The --use-case=quality-gates option is now deprecated and is now a synonym for the default installation of Keptn.", logging.InfoLevel)
		}

		var chartRepoURL string
		// Determine installer version
		if installParams.ChartRepoURL != nil && *installParams.ChartRepoURL != "" {
			chartRepoURL = *installParams.ChartRepoURL
		} else if version.IsOfficialKeptnVersion(Version) {
			version, _ := version.GetOfficialKeptnVersion(Version)
			chartRepoURL = keptnInstallerHelmRepoURL + "keptn-" + version + ".tgz"
		} else {
			chartRepoURL = keptnInstallerHelmRepoURL + "latest/keptn-0.1.0.tgz"
		}

		if err := downloadChart(chartRepoURL); err != nil {
			return err
		}

		logging.PrintLog(fmt.Sprintf("Helm Chart used for Keptn installation: %s", chartRepoURL), logging.InfoLevel)

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
					logging.PrintLog("See https://keptn.sh/docs/"+keptnReleaseDocsURL+"/installation/k8s-support/ for details.", logging.VerboseLevel)
					return errors.New(`Keptn requires Kubernetes server version: ` + KubeServerVersionConstraints)
				}
			}

			// check if istio-system namespace and istio-ingressgateway is available (full installation only)
			if installParams.UseCase == ContinuousDelivery {
				if err := checkIstioInstallation(); err != nil {
					logging.PrintLog("Istio is required for the continuous-delivery use case, "+
						"but could not be found in your cluster in namespace istio-system.", logging.InfoLevel)
					logging.PrintLog("Please install Istio as described "+
						"in the Istio docs: https://istio.io/latest/docs/setup/getting-started/", logging.InfoLevel)
				}
			}

			return doInstallation()
		}
		fmt.Println("Skipping installation due to mocking flag")
		return nil
	},
}

func downloadChart(chartRepoURL string) error {

	resp, err := http.Get(chartRepoURL)
	if err != nil {
		return errors.New("error retrieving Keptn Helm Chart at " + chartRepoURL + ": " + err.Error())
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("error retrieving Keptn Helm Chart at " + chartRepoURL + ": " + err.Error())
	}

	installChart, err = keptnutils.LoadChart(bytes)
	if err != nil {
		return errors.New("error retrieving Keptn Helm Chart at " + chartRepoURL + ": " + err.Error())
	}
	return nil
}

func checkIstioInstallation() error {
	clientset, err := keptnutils.GetClientset(false)
	if err != nil {
		return err
	}
	_, err = clientset.CoreV1().Namespaces().Get("istio-system", metav1.GetOptions{})
	if err != nil {
		return err
	}

	_, err = clientset.CoreV1().Services("istio-system").Get("istio-ingressgateway", metav1.GetOptions{})
	if err != nil {
		return err
	}

	return nil
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
		"The platform to run keptn on ["+kubernetes+","+openshift+"]")

	installParams.ConfigFilePath = installCmd.Flags().StringP("creds", "c", "",
		"Specify a JSON file containing cluster information needed for the installation (this allows skipping user prompts to execute a *silent* Keptn installation)")

	installParams.UseCaseInput = installCmd.Flags().StringP("use-case", "u", "",
		"The use case to install Keptn for ["+ContinuousDelivery.String()+","+QualityGates.String()+"]")

	installParams.ApiServiceTypeInput = installCmd.Flags().StringP("keptn-api-service-type", "",
		ClusterIP.String(), "Installation options for the api-service type ["+ClusterIP.String()+","+
			LoadBalancer.String()+","+NodePort.String()+"]")

	installParams.ChartRepoURL = installCmd.Flags().StringP("chart-repo", "",
		"", "URL of the Keptn Helm Chart repository")
	installCmd.Flags().MarkHidden("chart-repo")

	installCmd.PersistentFlags().BoolVarP(&insecureSkipTLSVerify, "insecure-skip-tls-verify", "s",
		false, "Skip tls verification for kubectl commands")
}

func createKeptnNamespace(keptnNamespace string) error {

	res, err := keptnutils.ExistsNamespace(false, keptnNamespace)
	if err != nil {
		return fmt.Errorf("Failed to check if namespace %s already exists: %v", keptnNamespace, err)
	}
	if res {
		return fmt.Errorf("Existing Keptn installation found in namespace %s.", keptnNamespace)
	}

	err = keptnutils.CreateNamespace(false, keptnNamespace)
	if err != nil {
		return fmt.Errorf("Failed to create Keptn namespace %s: %v", keptnNamespace, err)
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
	if err := createKeptnNamespace("keptn"); err != nil {
		return err
	}

	if len(ch.Templates) > 0 {
		logging.PrintLog(fmt.Sprintf("Start upgrading Helm Chart %s in namespace: %s", releaseName, namespace), logging.InfoLevel)
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
			return fmt.Errorf("Error when installing/upgrading Helm Chart %s in namespace %s: %s",
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
		logging.PrintLog(fmt.Sprintf("Finished upgrading Helm Chart %s in namespace %s", releaseName, namespace), logging.InfoLevel)
	} else {
		logging.PrintLog("Upgrade not done since this is an empty Helm Chart", logging.InfoLevel)
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

	values := map[string]interface{}{
		"continuous-delivery": map[string]interface{}{
			"enabled": installParams.UseCase == ContinuousDelivery,
			"openshift": map[string]interface{}{
				"enabled": *installParams.PlatformIdentifier == "openshift",
			},
		},
		"control-plane": map[string]interface{}{
			"enabled": true,
			"apiNginxGateway": map[string]interface{}{
				"type": installParams.ApiServiceType.String(),
			},
		},
	}

	if err := upgradeChart(installChart, "keptn", "keptn", values); err != nil {
		logging.PrintLog("Could not complete Keptn installation: "+err.Error(), logging.InfoLevel)
		return err
	}

	logging.PrintLog("Keptn has been successfully set up on your cluster.", logging.InfoLevel)
	logging.PrintLog("To connect the Keptn CLI with the Keptn API on your cluster, "+
		"please refer to the instructions at https://keptn.sh/docs/"+keptnReleaseDocsURL+"/operate", logging.InfoLevel)
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

	for {
		p.readCreds()

		fmt.Println()
		fmt.Println("Please confirm that the provided cluster information is correct: ")

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
