//go:build !nokubectl
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
	"errors"
	"fmt"
	goutils "github.com/keptn/go-utils/pkg/common/httputils"
	"github.com/keptn/keptn/cli/pkg/common"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/helm"
	"github.com/keptn/keptn/cli/pkg/kube"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/pkg/platform"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/chart"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
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

var continuousDeliveryServiceCharts []*chart.Chart

const helmServiceName = "helm-service"
const jmeterServiceName = "jmeter-service"

var continuousDeliveryServices = []string{helmServiceName, jmeterServiceName}

// installCmd represents the version command
var installCmd = NewInstallCmd(helm.NewHelper(), kube.NewKubernetesUtilsKeptnNamespaceHandler(), common.NewUserInput())

func NewInstallCmd(helmHelper helm.IHelper, namespaceHandler kube.IKeptnNamespaceHandler, userInput common.IUserInput) *cobra.Command {
	installCmdHandler := &InstallCmdHandler{
		helmHelper:       helmHelper,
		namespaceHandler: namespaceHandler,
		userInput:        userInput,
	}
	return &cobra.Command{
		Use:   "install",
		Args:  cobra.NoArgs,
		Short: "Installs Keptn on a Kubernetes cluster",
		Long: `The Keptn CLI allows installing Keptn on any Kubernetes derivative to which your kube config is pointing, and on OpenShift.

For more information, please follow the installation guide [Install Keptn](https://keptn.sh/docs/` + getReleaseDocsURL() + `/operate/install/#install-keptn)
`,
		Example: `keptn install                                                          # install on Kubernetes

keptn install --platform=openshift --use-case=continuous-delivery      # install continuous delivery on OpenShift

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

			installPlatformManager, err := platform.NewPlatformManager(*installParams.PlatformIdentifier, credentialmanager.NewCredentialManager(assumeYes))
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
				err = installPlatformManager.ReadCreds(assumeYes)
				if err != nil {
					return err
				}
			}

			// check if Kubernetes server version is compatible (except OpenShift)
			if *installParams.PlatformIdentifier != platform.OpenShiftIdentifier {
				if isNewerVersion, err := kube.CheckKubeServerVersion(KubeServerVersionConstraints); err != nil {
					logging.PrintLog(err.Error(), logging.VerboseLevel)
					logging.PrintLog("See https://keptn.sh/docs/"+getReleaseDocsURL()+"/operate/k8s_support/ for details.", logging.VerboseLevel)
					return fmt.Errorf("Failed to check kubernetes server version: %w", err)
				} else if isNewerVersion {
					logging.PrintLog("The Kubernetes server version is higher than the one officially supported. This is not recommended and could have negative impacts on the stability of Keptn - use at your own risk.", logging.InfoLevel)
					userConfirmation := userInput.AskBool("Do you want to continue?", &common.UserInputOptions{assumeYes})

					if !userConfirmation {
						return fmt.Errorf("Stopping installation.")
					}
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
			return installCmdHandler.doInstallation(installParams)
		},
	}
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

////////
type InstallCmdHandler struct {
	helmHelper       helm.IHelper
	namespaceHandler kube.IKeptnNamespaceHandler
	userInput        common.IUserInput
}

func fetchCharts(helmHelper helm.IHelper, keptnChartRepoURL string) (*chart.Chart, error) {
	var err error
	var chart *chart.Chart
	if goutils.IsValidURL(keptnChartRepoURL) {
		chart, err = helmHelper.DownloadChart(keptnChartRepoURL)
		if err != nil {
			return nil, err
		}
	} else {
		chart, err = keptnutils.LoadChartFromPath(keptnChartRepoURL)
		if err != nil {
			return nil, err
		}
	}
	return chart, nil
}

func (i *InstallCmdHandler) doInstallation(installParams installCmdParams) error {
	keptnNamespace := namespace
	showFallbackConnectMessage := true
	var keptnChartRepoURL string

	if !isStringFlagSet(installParams.ChartRepoURL) {
		keptnChartRepoURL = getKeptnHelmChartRepoURL()
	} else {
		keptnChartRepoURL = *installParams.ChartRepoURL
	}

	var err error
	if keptnChart, err = fetchCharts(i.helmHelper, keptnChartRepoURL); err != nil {
		return err
	}

	if installParams.UseCase == ContinuousDelivery {
		if continuousDeliveryServiceCharts, err = fetchContinuousDeliveryCharts(i.helmHelper, installParams.ChartRepoURL); err != nil {
			return err
		}
	}

	logging.PrintLog(fmt.Sprintf("Helm Chart used for Keptn installation: %s", keptnChartRepoURL), logging.InfoLevel)

	namespaceExists, err := i.namespaceHandler.ExistsNamespace(false, keptnNamespace)
	if err != nil {
		return fmt.Errorf("Failed to check if namespace %s already exists: %v", keptnNamespace, err)
	}

	if namespaceExists {
		fmt.Printf("Existing Keptn installation found in namespace %s\n\n", keptnNamespace)
		userConfirmation := i.userInput.AskBool("Do you want to overwrite this installation?", &common.UserInputOptions{assumeYes})

		if !userConfirmation {
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
		if err := i.namespaceHandler.CreateNamespace(false, keptnNamespace, namespaceMetadata); err != nil {
			return fmt.Errorf("Failed to create Keptn namespace %s: %v", keptnNamespace, err)
		}
	}

	values := map[string]interface{}{
		"continuous-delivery": map[string]interface{}{
			"enabled": installParams.UseCase == ContinuousDelivery,
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

	if err := i.helmHelper.UpgradeChart(keptnChart, keptnReleaseName, keptnNamespace, values); err != nil {
		msg := fmt.Sprintf("Could not complete Keptn installation: %s \nFor troubleshooting, please check the status of the keptn deployment by executing the following command: \n\nkubectl get pods -n %s\n", err.Error(), keptnNamespace)
		return errors.New(msg)
	}

	logging.PrintLog("Keptn control plane has been successfully set up on your cluster.", logging.InfoLevel)

	if installParams.UseCase == ContinuousDelivery {
		logging.PrintLog("Installing execution plane services for continuous-delivery use case.", logging.InfoLevel)
		for _, serviceChart := range continuousDeliveryServiceCharts {
			if err := i.helmHelper.UpgradeChart(serviceChart, serviceChart.Name(), keptnNamespace, values); err != nil {
				msg := fmt.Sprintf("Could not complete Keptn installation: %s \nFor troubleshooting, please check the status of the keptn deployment by executing the following command: \n\nkubectl get pods -n %s\n", err.Error(), keptnNamespace)
				return errors.New(msg)
			}
		}
	}

	// Hide sensitive information like api-token and endpoint in post-installation output
	if *installParams.HideSensitiveData {
		return nil
	}

	logging.PrintLog("---------------------------------------------------", logging.InfoLevel)

	if installParams.EndPointServiceType.String() == "NodePort" || installParams.EndPointServiceType.String() == "LoadBalancer" {
		endpoint, err := getAPIEndpoint(keptnNamespace, installParams.EndPointServiceType.String())
		if err == nil {
			showFallbackConnectMessage = false
			common.PrintQuickAccessInstructions(keptnNamespace, getReleaseDocsURL(), endpoint)
		}
	}

	if showFallbackConnectMessage {
		common.PrintQuickAccessInstructions(keptnNamespace, getReleaseDocsURL(), "http://localhost:8080/api")
	}

	return nil
}

func fetchContinuousDeliveryCharts(helmHelper helm.IHelper, chartRepoURL *string) ([]*chart.Chart, error) {
	charts := []*chart.Chart{}
	for _, service := range continuousDeliveryServices {
		chartURL := getExecutionPlaneServiceChartRepoURL(chartRepoURL, service)
		serviceChart, err := fetchCharts(helmHelper, chartURL)
		if err != nil {
			return nil, err
		}
		charts = append(charts, serviceChart)
	}
	return charts, nil
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

	_, err = clientset.CoreV1().Namespaces().Get(rootCmd.Context(), "istio-system", metav1.GetOptions{})
	if err != nil {
		return err
	}

	_, err = clientset.CoreV1().Services("istio-system").Get(rootCmd.Context(), "istio-ingressgateway", metav1.GetOptions{})
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
		// Fetching the EXTERNAL-IP of the api-gateway-nginx loadbalancer service
		external, err := keptnutils.ExecuteCommand("kubectl", []string{"get", "svc", "api-gateway-nginx", "-n", keptnNamespace, "-o", "jsonpath='{.status.loadBalancer.ingress[0].ip}'"})
		if err != nil {
			return "", err
		}
		endpoint = strings.Trim(external, "'")
		return "http://" + endpoint + "/api", nil
	}
	return "", errors.New("Unknown service-type: " + serviceType)
}
