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

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/keptn/keptn/cli/pkg/helm"
	"github.com/keptn/keptn/cli/pkg/kube"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/pkg/platform"

	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/strvals"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var installParams installUpgradeParams
var keptnInstallChart *chart.Chart
var installValues = getDefaultValues()


// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs Keptn on a Kubernetes cluster",
	Long: `The Keptn CLI allows installing Keptn on any Kubernetes derivate to which your kube config is pointing to, and on OpenShift.

For more information, please follow the installation guide [Install Keptn](https://keptn.sh/docs/` + keptnReleaseDocsURL + `/operate/install/#install-keptn)
`,
	Example: `keptn install --values=./my-values.yaml`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {

		chartRepoURL := getChartRepoURL(installParams.ChartRepoURL)

		var err error
		if keptnInstallChart, err = helm.NewHelper().DownloadChart(chartRepoURL); err != nil {
			return err
		}

		logging.PrintLog(fmt.Sprintf("Helm Chart used for Keptn installation: %s", chartRepoURL), logging.InfoLevel)

		var valuesFile map[string]interface{}
		if installParams.ValuesFile != nil && *installParams.ValuesFile != "" {
			err = getAndParseYaml(*installParams.ValuesFile, &valuesFile)
			if err != nil {
				return fmt.Errorf("Failed to read and parse values file - %s", err.Error())
			}
		}
		installValues = mergeMaps(installValues, valuesFile)

		if installParams.Values != nil {
			for _, value := range *installParams.Values {
				if err := strvals.ParseInto(value, installValues); err != nil {
					return fmt.Errorf("Failed to parse --set data - %s", err.Error())
				}
			}
		}

		return checkInput(installParams, installValues)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		logging.PrintLog("Installing Keptn ...", logging.InfoLevel)
		if !mocking {
			return doInstallation()
		}
		fmt.Println("Skipping installation due to mocking flag")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installParams = installUpgradeParams{}

	installParams.ConfigFilePath = installCmd.Flags().StringP("creds", "c", "",
		"Specify a JSON file containing cluster information needed for the installation. This allows skipping user prompts to execute a *silent* Keptn installation.")

	installParams.ChartRepoURL = installCmd.Flags().StringP("chart-repo", "",
		"", "URL of the Keptn Helm Chart repository")
	installCmd.Flags().MarkHidden("chart-repo")

	installCmd.PersistentFlags().BoolVarP(&insecureSkipTLSVerify, "insecure-skip-tls-verify", "s",
		false, "Skip tls verification for kubectl commands")

	installParams.Namespace = installCmd.Flags().StringP("namespace", "n", "keptn",
		"Specify the namespace where Keptn should be installed in (default keptn).")

	installParams.ValuesFile = installCmd.Flags().StringP("values", "f", "",
		"Specify values in a YAML file or a URL.")

	installParams.Values = installCmd.Flags().StringArrayP("set", "", []string{},
		"Set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
}

// Preconditions: 1. Already authenticated against the cluster.
func doInstallation() error {
	keptnNamespace := *installParams.Namespace

	res, err := keptnutils.ExistsNamespace(false, keptnNamespace)
	if err != nil {
		return fmt.Errorf("Failed to check if namespace %s already exists: %v", keptnNamespace, err)
	}

	if res {
		fmt.Printf("Existing Keptn installation found in namespace %s\n", keptnNamespace)
		fmt.Println()
		fmt.Println("Do you want to overwrite this installation? (y/n)")

		reader := bufio.NewReader(os.Stdin)
		in, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		in = strings.ToLower(strings.TrimSpace(in))
		if !(in == "y" || in == "yes") {
			return fmt.Errorf("Stopping installation.")
		}
	} else {
		if err := keptnutils.CreateNamespace(false, keptnNamespace); err != nil {
			return fmt.Errorf("Failed to create Keptn namespace %s: %v", keptnNamespace, err)
		}
	}

	if err := helm.NewHelper().UpgradeChart(keptnInstallChart, keptnReleaseName, keptnNamespace, installValues); err != nil {
		logging.PrintLog("Could not complete Keptn installation: "+err.Error(), logging.InfoLevel)
		return err
	}

	logging.PrintLog("Keptn has been successfully set up on your cluster.", logging.InfoLevel)
	logging.PrintLog("---------------------------------------------------", logging.InfoLevel)
	fmt.Println("* To quickly access Keptn, you can use a port-forward and then authenticate your Keptn CLI (in a Linux shell):\n" +
		" - kubectl -n " + keptnNamespace + " port-forward service/api-gateway-nginx 8080:80\n" +
		" - keptn auth --endpoint=http://localhost:8080/api --api-token=$(kubectl get secret keptn-api-token -n " + keptnNamespace + " -ojsonpath={.data.keptn-api-token} | base64 --decode)\n")
	fmt.Println("* To expose Keptn on a public endpoint, please continue with the installation guidelines provided at:\n" +
		" - https://keptn.sh/docs/" + keptnReleaseDocsURL + "/operate/install#install-keptn\n")
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

func checkInput(params installUpgradeParams, values map[string]interface{}) error {
	isOpenShiftEnabled := values["continuous-delivery"].(map[string]interface{})["openshift"].(map[string]interface{})["enabled"].(bool)
	isContinuousDeliveryEnabled := values["continuous-delivery"].(map[string]interface{})["enabled"].(bool)

	platformManager := platform.NewPlatformManager(isOpenShiftEnabled)
	if !mocking {
		if err := platformManager.CheckRequirements(); err != nil {
			return err
		}
	}

	if params.ConfigFilePath != nil && *params.ConfigFilePath != "" {
		// Config was provided in form of a file
		if err := platformManager.ParseConfig(*params.ConfigFilePath); err != nil {
			return err
		}

		// Check whether the authentication at the cluster is valid
		if err := platformManager.CheckCreds(); err != nil {
			return err
		}
	} else {
		err := platformManager.ReadCreds()
		if err != nil {
			return err
		}
	}

	// check if Kubernetes server version is compatible (except OpenShift)
	if !isOpenShiftEnabled {
		if err := kube.CheckKubeServerVersion(KubeServerVersionConstraints); err != nil {
			logging.PrintLog(err.Error(), logging.VerboseLevel)
			logging.PrintLog("See https://keptn.sh/docs/"+keptnReleaseDocsURL+"/operate/k8s_support/ for details.", logging.VerboseLevel)
			return fmt.Errorf("Failed to check kubernetes server version: %w", err)
		}
	}

	// check if istio-system namespace and istio-ingressgateway is available (full installation only)
	if isContinuousDeliveryEnabled {
		if err := checkIstioInstallation(); err != nil {
			logging.PrintLog("Istio is required for values.continuous-delivery.enabled: true, "+
				"but could not be found in your cluster in namespace istio-system.", logging.InfoLevel)
			logging.PrintLog("Please install Istio as described "+
				"in the Istio docs: https://istio.io/latest/docs/setup/getting-started/", logging.InfoLevel)
		}
	}

	return nil
}

func getDefaultValues() map[string]interface{} {
	return map[string]interface{}{
		"continuous-delivery": map[string]interface{}{
			"enabled": false,
			"openshift": map[string]interface{}{
				"enabled": false,
			},
		},
	}
}

// mergeMaps merges two map[string]interface{} together
// source: https://github.com/helm/helm/blob/9bc7934/pkg/cli/values/options.go#L88
func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}
