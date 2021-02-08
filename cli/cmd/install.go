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
	"os"
	"strings"

	"github.com/keptn/keptn/cli/pkg/common"

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
	HideSensitiveData        *bool
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
	Example: `keptn install                                                          # install on Kubernetes

keptn install --platform=openshift --use-case=continuous-delivery      # install continuous delivery on Openshift

keptn install --platform=kubernetes --endpoint-service-type=NodePort   # install on Kubernetes with gateway NodePort

keptn install --hide-sensitive-data                                    # install on Kubernetes and hides sensitive data like api-token and endpoint in post-installation output
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
		"Use --use-case=continuous-delivery to install the execution plane for continuous delivery. Without this flag, your Keptn is capable of the quality gate and automated remediations use-case.")

	installParams.EndPointServiceTypeInput = installCmd.Flags().StringP("endpoint-service-type", "",
		ClusterIP.String(), "Installation options for the endpoint-service type ["+ClusterIP.String()+","+
			LoadBalancer.String()+","+NodePort.String()+"]")

	installParams.ChartRepoURL = installCmd.Flags().StringP("chart-repo", "",
		"", "URL of the Keptn Helm Chart repository")
	installCmd.Flags().MarkHidden("chart-repo")

	installCmd.PersistentFlags().BoolVarP(&insecureSkipTLSVerify, "insecure-skip-tls-verify", "s",
		false, "Skip tls verification for kubectl commands")
	installParams.HideSensitiveData = installCmd.Flags().BoolP("hide-sensitive-data", "", false,
		"Hide the sensitive data like api-tokens and endpoints in post-installation output.")

}

// Preconditions: 1. Already authenticated against the cluster.
func doInstallation() error {
	keptnNamespace := namespace
	showFallbackConnectMessage := true

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
		namespaceMetadata := metav1.ObjectMeta{
			Annotations: map[string]string{
				"keptn.sh/managed-by": "keptn",
			},
			Labels: map[string]string{
				"keptn.sh/managed-by": "keptn",
			},
		}
		if err := keptnutils.CreateNamespace(false, keptnNamespace, namespaceMetadata); err != nil {
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
			"bridge": map[string]interface{}{
				"installationType": getInstallationTypeEnvVar(*installParams.UseCaseInput),
			},
		},
	}

	if err := helm.NewHelper().UpgradeChart(keptnChart, keptnReleaseName, keptnNamespace, values); err != nil {
		msg := fmt.Sprintf("Could not complete Keptn installation: %s \nFor troubleshooting, please check the status of the keptn deployment by executing the following command: \n\nkubectl get pods -n %s\n", err.Error(), keptnNamespace)
		return errors.New(msg)
	}

	logging.PrintLog("Keptn has been successfully set up on your cluster.", logging.InfoLevel)

	// Hide sensitive information like api-token and endpoint in post-installation output
	if *installParams.HideSensitiveData {
		return nil
	}

	logging.PrintLog("---------------------------------------------------", logging.InfoLevel)

	if installParams.EndPointServiceType.String() == "NodePort" || installParams.EndPointServiceType.String() == "LoadBalancer" {
		endpoint, err := getAPIEndpoint(keptnNamespace, installParams.EndPointServiceType.String())
		if err == nil {
			showFallbackConnectMessage = false
			common.PrintQuickAccessInstructions(keptnNamespace, keptnReleaseDocsURL, endpoint)
		}
	}

	if showFallbackConnectMessage {
		common.PrintQuickAccessInstructions(keptnNamespace, keptnReleaseDocsURL, "http://localhost:8080/api")
	}

	return nil
}

func getInstallationTypeEnvVar(useCase string) string {
	if useCase == ContinuousDelivery.String() {
		return "QUALITY_GATES,CONTINUOUS_OPERATIONS,CONTINUOUS_DELIVERY"
	}
	return "QUALITY_GATES,CONTINUOUS_OPERATIONS"
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

func getAPIEndpoint(keptnNamespace string, serviceType string) (string, error) {
	var endpoint, port string
	switch serviceType {
	case "NodePort":
		// Fetching external and internal node IP
		external, err := keptnutils.ExecuteCommand("kubectl", []string{"get", "nodes", "-o", "jsonpath='{ $.items[0].status.addresses[?(@.type==\"ExternalIP\")].address }'"})
		internal, err := keptnutils.ExecuteCommand("kubectl", []string{"get", "nodes", "-o", "jsonpath='{ $.items[0].status.addresses[?(@.type==\"InternalIP\")].address }'"})
		if err != nil {
			return "", err
		}
		endpoint = strings.Trim(external, "'")
		internal = strings.Trim(internal, "'")
		// Fetching mapped port of the api-gateway-nginx nodeport service
		port, _ = keptnutils.ExecuteCommand("kubectl", []string{"get", "svc", "api-gateway-nginx", "-n", keptnNamespace, "-o", "jsonpath='{.spec.ports[?(@.name==\"http\")].nodePort}'"})
		port = strings.Trim(port, "'")
		if endpoint == "" {
			endpoint = internal
		}
		return "http://" + endpoint + ":" + port + "/api", nil
	case "LoadBalancer":
		// Fetching the EXTERNAL-IP of the api-gateway-ngix loadbalancer service
		external, err := keptnutils.ExecuteCommand("kubectl", []string{"get", "svc", "api-gateway-nginx", "-n", keptnNamespace, "-o", "jsonpath='{.status.loadBalancer.ingress[0].ip}'"})
		if err != nil {
			return "", err
		}
		endpoint = strings.Trim(external, "'")
		return "http://" + endpoint + "/api", nil
	}
	return "", errors.New("Unknown service-type: " + serviceType)
}
