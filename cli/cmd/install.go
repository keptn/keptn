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
	"errors"
	"fmt"
	"github.com/keptn/keptn/cli/pkg/common"
	"os"
	"strings"

	"helm.sh/helm/v3/pkg/chart"

	"github.com/keptn/keptn/cli/pkg/platform"

	"github.com/keptn/keptn/cli/pkg/helm"

	"github.com/keptn/keptn/cli/pkg/kube"

	"github.com/keptn/keptn/cli/pkg/logging"

	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type installCmdParams struct {
	installUpgradeParams
	UseCaseInput             *string
	UseCase                  usecase
	EndPointServiceTypeInput *string
	EndPointServiceType      endpointServiceType
}

var installParams installCmdParams
var keptnChart *chart.Chart

// installCmd represents the version command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs Keptn on a Kubernetes cluster",
	Long: `The Keptn CLI allows installing Keptn on any Kubernetes derivate to which your kube config is pointing to, and on OpenShift.

For more information, please follow the installation guide [Install Keptn](https://keptn.sh/docs/` + keptnReleaseDocsURL + `/operate/install/#install-keptn)
`,
	Example: `keptn install                                                        # install on Kubernetes

keptn install --platform=openshift --use-case=continuous-delivery    # install continuous delivery on Openshift

keptn install --platform=kubernetes --endpoint-service-type=NodePort # install on Kubernetes with gateway NodePort
`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {

		// Parse api service type
		if val, ok := apiServiceTypeToID[*installParams.EndPointServiceTypeInput]; ok {
			installParams.EndPointServiceType = val
		} else {
			return errors.New("value of 'endpoint-service-type' flag is unknown. Supported values are " +
				"[" + ClusterIP.String() + "," + LoadBalancer.String() + "," + NodePort.String() + "]")
		}

		if val, ok := usecaseToID[*installParams.UseCaseInput]; ok {
			installParams.UseCase = val
		} else {
			return errors.New("value of 'use-case' flag is unknown. Supported values are " +
				"[" + ContinuousDelivery.String() + "]")
		}

		// Mark the quality-gates use case as deprecated - this is now the default option
		if *installParams.UseCaseInput == QualityGates.String() {
			logging.PrintLog("Note: The --use-case=quality-gates option is now deprecated and is now a synonym for the default installation of Keptn.", logging.InfoLevel)
		}

		chartRepoURL := getChartRepoURL(installParams.ChartRepoURL)

		var err error
		if keptnChart, err = helm.NewHelper().DownloadChart(chartRepoURL); err != nil {
			return err
		}

		logging.PrintLog(fmt.Sprintf("Helm Chart used for Keptn installation: %s", chartRepoURL), logging.InfoLevel)

		installPlatformManager, err := platform.NewPlatformManager(*installParams.PlatformIdentifier)
		if err != nil {
			return err
		}

		if !mocking {
			if err := installPlatformManager.CheckRequirements(); err != nil {
				return err
			}
		}

		if installParams.ConfigFilePath != nil && *installParams.ConfigFilePath != "" {
			// Config was provided in form of a file
			if err := installPlatformManager.ParseConfig(*installParams.ConfigFilePath); err != nil {
				return err
			}

			// Check whether the authentication at the cluster is valid
			if err := installPlatformManager.CheckCreds(); err != nil {
				return err
			}
		} else {
			err = installPlatformManager.ReadCreds()
			if err != nil {
				return err
			}
		}

		// check if Kubernetes server version is compatible (except OpenShift)
		if *installParams.PlatformIdentifier != platform.OpenShiftIdentifier {
			if err := kube.CheckKubeServerVersion(KubeServerVersionConstraints); err != nil {
				logging.PrintLog(err.Error(), logging.VerboseLevel)
				logging.PrintLog("See https://keptn.sh/docs/"+keptnReleaseDocsURL+"/operate/k8s_support/ for details.", logging.VerboseLevel)
				return fmt.Errorf("Failed to check kubernetes server version: %w", err)
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

		return nil
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
	installParams = installCmdParams{}

	installParams.PlatformIdentifier = installCmd.Flags().StringP("platform", "p", "kubernetes",
		"The platform to run Keptn on ["+platform.KubernetesIdentifier+","+platform.OpenShiftIdentifier+"]")

	installParams.ConfigFilePath = installCmd.Flags().StringP("creds", "c", "",
		"Specify a JSON file containing cluster information needed for the installation. This allows skipping user prompts to execute a *silent* Keptn installation.")

	installParams.UseCaseInput = installCmd.Flags().StringP("use-case", "u", "",
		"The use case to install Keptn for ["+ContinuousDelivery.String()+"]")

	installParams.EndPointServiceTypeInput = installCmd.Flags().StringP("endpoint-service-type", "",
		ClusterIP.String(), "Installation options for the endpoint-service type ["+ClusterIP.String()+","+
			LoadBalancer.String()+","+NodePort.String()+"]")

	installParams.ChartRepoURL = installCmd.Flags().StringP("chart-repo", "",
		"", "URL of the Keptn Helm Chart repository")
	installCmd.Flags().MarkHidden("chart-repo")

	installCmd.PersistentFlags().BoolVarP(&insecureSkipTLSVerify, "insecure-skip-tls-verify", "s",
		false, "Skip tls verification for kubectl commands")

	installParams.Namespace = installCmd.Flags().StringP("namespace", "n", "keptn",
		"Specify the namespace where Keptn should be installed in (default keptn).")
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

	values := map[string]interface{}{
		"continuous-delivery": map[string]interface{}{
			"enabled": installParams.UseCase == ContinuousDelivery,
			"openshift": map[string]interface{}{
				"enabled": *installParams.PlatformIdentifier == "openshift",
			},
		},
		"control-plane": map[string]interface{}{
			"enabled": true,
			"apiGatewayNginx": map[string]interface{}{
				"type": installParams.EndPointServiceType.String(),
			},
		},
	}

	if err := helm.NewHelper().UpgradeChart(keptnChart, keptnReleaseName, keptnNamespace, values); err != nil {
		logging.PrintLog("Could not complete Keptn installation: "+err.Error(), logging.InfoLevel)
		return err
	}

	logging.PrintLog("Keptn has been successfully set up on your cluster.", logging.InfoLevel)
	logging.PrintLog("---------------------------------------------------", logging.InfoLevel)

	common.PrintQuickAccessInstructions(keptnNamespace, keptnReleaseDocsURL)

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
