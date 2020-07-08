// Inspired by `hugo gen doc`  - see https://github.com/gohugoio/hugo/blob/release-0.69.0/commands/gendoc.go
package cmd

import (
	"archive/zip"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/keptn/keptn/cli/pkg/exechelper"

	"github.com/keptn/go-utils/pkg/api/models"

	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"

	"github.com/spf13/cobra"
)

var generateSupportArchiveParams *generateCmdParams

type errorableStringResult struct {
	Result string
	Err    error `json:",omitempty"`
}

type errorableBoolResult struct {
	Result bool
	Err    error `json:",omitempty"`
}

type errorableProjectResult struct {
	Result []*models.Project
	Err    error `json:",omitempty"`
}

type errorableMetadataResult struct {
	Result *models.Metadata
	Err    error `json:",omitempty"`
}

type errorableIPInfoResult struct {
	Result []*ipInfo
	Err    error `json:",omitempty"`
}

type ipInfo struct {
	Name       string `json:",omitempty"`
	InternalIP string `json:",omitempty"`
	ExternalIP string `json:",omitempty"`
}

var namespaces = [...]string{"keptn", "keptn-datastore"} //"istio-system"

type metaData struct {
	OperatingSystem                 string
	KeptnCLIVersion                 string
	KeptnAPIMetadata                *errorableMetadataResult `json:",omitempty"`
	KeptnAPIUrl                     *errorableStringResult  `json:",omitempty"`
	KeptnAPIReachable               *errorableBoolResult    `json:",omitempty"`
	Projects                        *errorableProjectResult `json:",omitempty"`
	KubectlVersion                  *errorableStringResult  `json:",omitempty"`
	KubeContextPointsToKeptnCluster *errorableBoolResult    `json:",omitempty"`
	IngressHostnameSuffix           *errorableStringResult  `json:",omitempty"`
	IstioSystemInstalled            *errorableBoolResult    `json:",omitempty"`
	ConfiguredGatewayAvailable      *errorableBoolResult    `json:",omitempty"`
	IngressPort                     *errorableStringResult  `json:",omitempty"`
	IngressProtocol                 *errorableStringResult  `json:",omitempty"`
	IngressGateway                  *errorableStringResult  `json:",omitempty"`
}

// generateSupportArchiveCmd implements the generate support-archive command
var generateSupportArchiveCmd = &cobra.Command{
	Use:   "support-archive",
	Short: "Generates a support archive.",
	Long:  `Generates a support archive containing information of the Keptn installation.`,
	Example: `keptn generate support-archive
keptn generate support-archive --dir=/some/directory`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		outputDir := "./support-archive"
		if *generateSupportArchiveParams.Directory != "" {
			outputDir = *generateSupportArchiveParams.Directory
		}

		// check if output directory exists
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			// outputDir does not exist
			return errors.New(fmt.Sprintf("Error trying to access directory %s. Please make sure the directory exists.", outputDir))
		}

		tmpDir, err := ioutil.TempDir("", "keptn-support-archive")
		if err != nil {
			return fmt.Errorf("Error when creating a temporary directory: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		s := &metaData{}
		s.OperatingSystem = runtime.GOOS
		s.KeptnCLIVersion = Version

		if !mocking {
			s.KeptnAPIUrl = getKeptnAPIUrl()
			if s.KeptnAPIUrl.Err == nil {
				s.KeptnAPIReachable = getKeptnAPIReachable()
				if s.KeptnAPIReachable.Err == nil && s.KeptnAPIReachable.Result {
					s.Projects = getProjects()
					s.KeptnAPIMetadata = getKeptnMetadata()
				}
			}
			writeKeptnInstallerLog(keptnInstallerLogFileName, tmpDir)
			writeKeptnInstallerLog(keptnInstallerErrorLogFileName, tmpDir)

			s.KubectlVersion = getKubectlVersion()
			if s.KubectlVersion.Err == nil {
				s.KubeContextPointsToKeptnCluster = getKubeContextPointsToKeptnCluster()

				if s.KubeContextPointsToKeptnCluster.Err == nil && s.KubeContextPointsToKeptnCluster.Result {
					ctx, _ := getKubeContext()
					fmt.Println("Retrieving logs from cluster " + strings.TrimSpace(ctx))
					s.IngressHostnameSuffix = getIngressHostnameSuffix()
					s.IngressPort = getIngressPort()
					s.IngressProtocol = getIngressProtocol()
					s.IngressGateway = getIngressGateway()
					writeNamespaces(tmpDir)

					s.IstioSystemInstalled = isIstioSystemInstalled()
					if s.IstioSystemInstalled.Result {
						k8sNSFilePath := filepath.Join(tmpDir, "istio-system")
						err := os.MkdirAll(k8sNSFilePath, os.ModePerm)
						if err == nil {
							writeServices("istio-system", k8sNSFilePath)
							// list of all available ingress gateways
							writeIstioIngressGateways(tmpDir)
						} else {
							fmt.Printf("Error making directory %s: %v\n", k8sNSFilePath, err)
						}
					}
					s.ConfiguredGatewayAvailable = isConfiguredIngressGatewayAvailable()

					// get nodes + their internal and external IPs
					writeClusterNodes(tmpDir)

					// check if there are any services with type==LoadBalancer in the cluster
					writeLoadBalancerServices(tmpDir)

					for _, ns := range namespaces {
						k8sNSFilePath := filepath.Join(tmpDir, ns)
						err := os.MkdirAll(k8sNSFilePath, os.ModePerm)
						if err != nil {
							fmt.Printf("Error making directory %s: %v\n", k8sNSFilePath, err)
							continue
						}
						writeConfigMaps(ns, k8sNSFilePath)
						writeSecrets(ns, k8sNSFilePath)
						writeDeployments(ns, k8sNSFilePath)
						writePods(ns, k8sNSFilePath)
						writeServices(ns, k8sNSFilePath)
						writeVirtualServices(ns, k8sNSFilePath)
						writeIngresses(ns, k8sNSFilePath)
						writePodLogs(ns, k8sNSFilePath)
						writePodDescriptions(ns, k8sNSFilePath)
						writeDeploymentDescriptions(ns, k8sNSFilePath)
					}
				} else {
					fmt.Println("Your kube context does not point to a Keptn cluster!")
				}
			}
		}

		supportData, err := json.Marshal(s)
		if err != nil {
			return fmt.Errorf("Error marshalling the suppport data: %v", err)
		}

		err = ioutil.WriteFile(filepath.Join(tmpDir, "metadata.json"), supportData, 0644)
		if err != nil {
			return fmt.Errorf("Error writing file: %v", err)
		}

		supportArchive := filepath.Join(outputDir, "keptn-support-archive-"+strconv.FormatInt(time.Now().Unix(), 10)+".zip")
		if err := recursiveZip(tmpDir, supportArchive); err != nil {
			return fmt.Errorf("Error writing zip: %v", err)
		}

		fmt.Println("The support archive is available here " + supportArchive)
		fmt.Println("This support archive potentially contains sensitive data. Therefore, please first review it before distributing.")
		fmt.Println("If you need help, please use the #help channel in the Keptn Slack workspace https://slack.keptn.sh")
		return nil
	},
}

func isConfiguredIngressGatewayAvailable() *errorableBoolResult {
	res, err := exechelper.ExecuteCommand("kubectl", "get cm -n keptn ingress-config -ojsonpath='{.data.ingress_gateway}'")
	if err != nil {
		return newErrorableBoolResult(false, err)
	}
	var gatewayFullName string
	if res == "" {
		gatewayFullName = "public-gateway.istio-system"
	} else {
		gatewayFullName = res
	}

	var gatewayName, gatewayNamespace string
	gwSplit := strings.Split(gatewayFullName, ".")

	gatewayName = gwSplit[0]
	if len(gwSplit) > 1 {
		gatewayNamespace = gwSplit[1]
	} else {
		gatewayNamespace = "keptn"
	}

	res, err = exechelper.ExecuteCommand("kubectl", "get gateway -n "+gatewayNamespace+" "+gatewayName)
	if err != nil {
		return newErrorableBoolResult(false, err)
	}
	return newErrorableBoolResult(true, nil)
}

func writeIstioIngressGateways(dir string) {
	writeErrorableStringResult(newErrorableStringResult(exechelper.ExecuteCommand("kubectl", "get gateway --all-namespaces -oyaml")),
		filepath.Join(dir, "ingress-gateways.txt"))
}

func writeClusterNodes(dir string) {
	writeErrorableStringResult(newErrorableStringResult(exechelper.ExecuteCommand("kubectl", "get nodes -owide")),
		filepath.Join(dir, "nodes.txt"))
}

func writeLoadBalancerServices(dir string) {
	res, err := exechelper.ExecuteCommand("kubectl", "get svc --all-namespaces -owide")
	if err != nil {
		writeErrorableStringResult(newErrorableStringResult("", err), "loadbalancer-services.txt")
	}
	output := ""
	split := strings.Split(res, "\n")
	output = output + split[0] + "\n"

	if len(split) > 1 {
		for _, line := range split[1:] {
			if strings.Contains(line, "LoadBalancer") {
				output = output + line + "\n"
			}
		}
	}
	writeErrorableStringResult(newErrorableStringResult(output, nil),
		filepath.Join(dir, "loadbalancer-services.txt"))
}

func isIstioSystemInstalled() *errorableBoolResult {
	fmt.Println("Checking availability of istio-system namespace")
	_, err := exechelper.ExecuteCommand("kubectl",
		`get namespace istio-system`)
	if err != nil {
		return newErrorableBoolResult(false, err)
	}
	return newErrorableBoolResult(true, nil)
}

func recursiveZip(pathToZip, destPath string) error {
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destFile)
	defer myZip.Close()
	return filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))
		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
}

func writeErrorableStringResult(res *errorableStringResult, file string) {
	var data []byte
	if res.Err != nil {
		data = []byte(res.Err.Error())
	} else {
		data = []byte(res.Result)
	}
	err := ioutil.WriteFile(file, data, 0644)
	if err != nil {
		fmt.Printf("Error writing file %s: %v\n", file, err)
	}
}

func newErrorableStringResult(result string, err error) *errorableStringResult {
	return &errorableStringResult{
		Result: result,
		Err:    err,
	}
}

func newErrorableBoolResult(result bool, err error) *errorableBoolResult {
	return &errorableBoolResult{
		Result: result,
		Err:    err,
	}
}

func newErrorableProjectResult(result []*models.Project, err error) *errorableProjectResult {
	return &errorableProjectResult{
		Result: result,
		Err:    err,
	}
}

func newErrorableMetadataResult(result *models.Metadata, err error) *errorableMetadataResult {
	return &errorableMetadataResult{
		Result: result,
		Err:    err,
	}
}

func getKeptnAPIUrl() *errorableStringResult {
	fmt.Println("Retrieving Keptn API")
	endPoint, _, err := credentialmanager.NewCredentialManager().GetCreds()
	return newErrorableStringResult(endPoint.String(), err)
}

func getKeptnAPIReachable() *errorableBoolResult {
	fmt.Println("Checking availability of Keptn API")
	endPoint, _, err := credentialmanager.NewCredentialManager().GetCreds()
	if err != nil {
		return newErrorableBoolResult(false, err)
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	endPoint.Path = "swagger-ui"
	resp, err := http.Get(endPoint.String())
	if err != nil {
		return newErrorableBoolResult(false, err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return newErrorableBoolResult(true, nil)
	}
	return newErrorableBoolResult(false, errors.New(resp.Status))
}

func getKubeContextPointsToKeptnCluster() *errorableBoolResult {
	fmt.Println("Checking whether kube context points to Keptn cluster")
	keptnDomain, err := getEndpointUsingKube()
	if err != nil {
		return newErrorableBoolResult(false, err)
	}
	endPoint, _, err := credentialmanager.NewCredentialManager().GetCreds()
	if err != nil {
		return newErrorableBoolResult(false, err)
	}
	return newErrorableBoolResult(strings.Contains(endPoint.String(), keptnDomain), nil)
}

func getKubectlVersion() *errorableStringResult {
	fmt.Println("Retrieving kubectl version")
	return newErrorableStringResult(exechelper.ExecuteCommand("kubectl", "version"))
}

func getIngressHostnameSuffix() *errorableStringResult {
	fmt.Println("Retrieving Keptn domain")
	return newErrorableStringResult(exechelper.ExecuteCommand("kubectl", "get cm ingress-config -n keptn -ojsonpath='{.data.ingress_hostname_suffix}'"))
}

func getIngressPort() *errorableStringResult {
	res, err := exechelper.ExecuteCommand("kubectl", "get cm ingress-config -n keptn -ojsonpath='{.data.ingress_port}'")
	if err != nil {
		return newErrorableStringResult("", err)
	}
	if res == "" {
		return newErrorableStringResult("80", nil)
	}
	return newErrorableStringResult(res, nil)
}

func getIngressProtocol() *errorableStringResult {
	res, err := exechelper.ExecuteCommand("kubectl", "get cm ingress-config -n keptn -ojsonpath='{.data.ingress_protocol}'")
	if err != nil {
		return newErrorableStringResult("", err)
	}
	if res == "" {
		return newErrorableStringResult("http", nil)
	}
	return newErrorableStringResult(res, nil)
}

func getIngressGateway() *errorableStringResult {
	res, err := exechelper.ExecuteCommand("kubectl", "get cm ingress-config -n keptn -ojsonpath='{.data.ingress_gateway}'")
	if err != nil {
		return newErrorableStringResult("", err)
	}
	if res == "" {
		return newErrorableStringResult("public-gateway.istio-system", nil)
	}
	return newErrorableStringResult(res, nil)
}

func writeNamespaces(dir string) {
	fmt.Println("Retrieving namespaces")
	writeErrorableStringResult(newErrorableStringResult(exechelper.ExecuteCommand("kubectl", "get namespaces")),
		filepath.Join(dir, "namespaces.txt"))
}

func writeConfigMaps(namespace, dir string) {
	fmt.Println("Retrieving list of config maps in " + namespace)
	writeErrorableStringResult(newErrorableStringResult(exechelper.ExecuteCommand("kubectl", "get cm -n "+namespace)),
		filepath.Join(dir, "configmap.txt"))
}

func writeSecrets(namespace, dir string) {
	fmt.Println("Retrieving list of secrets in " + namespace)
	writeErrorableStringResult(newErrorableStringResult(exechelper.ExecuteCommand("kubectl", "get secrets -n "+namespace)),
		filepath.Join(dir, "secrets.txt"))
}

func writeDeployments(namespace, dir string) {
	fmt.Println("Retrieving list of deployments in " + namespace)
	writeErrorableStringResult(newErrorableStringResult(exechelper.ExecuteCommand("kubectl", "get deployments -owide -n "+namespace)),
		filepath.Join(dir, "deployments.txt"))
}

func writePods(namespace, dir string) {
	fmt.Println("Retrieving list of pods in " + namespace)
	writeErrorableStringResult(newErrorableStringResult(exechelper.ExecuteCommand("kubectl", "get pods -owide  -n "+namespace)),
		filepath.Join(dir, "pods.txt"))
}

func writeServices(namespace, dir string) {
	fmt.Println("Retrieving list of services in " + namespace)
	writeErrorableStringResult(newErrorableStringResult(exechelper.ExecuteCommand("kubectl", "get services -owide -n "+namespace)),
		filepath.Join(dir, "services.txt"))
}

func writeVirtualServices(namespace, dir string) {
	fmt.Println("Retrieving list of virtual services in " + namespace)
	writeErrorableStringResult(newErrorableStringResult(exechelper.ExecuteCommand("kubectl", "get vs -owide -n "+namespace)),
		filepath.Join(dir, "virtualservices.txt"))

}

func writeIngresses(namespace, dir string) {
	fmt.Println("Retrieving list of ingresses in " + namespace)
	writeErrorableStringResult(newErrorableStringResult(exechelper.ExecuteCommand("kubectl", "get ingress -owide -n "+namespace)),
		filepath.Join(dir, "ingresses.txt"))
}

func writePodDescriptions(namespace, dir string) {
	fmt.Println("Retrieving pod descriptions in " + namespace)
	res, err := exechelper.ExecuteCommand("kubectl",
		`get pods --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' -n `+namespace)
	if err != nil {
		writeErrorableStringResult(newErrorableStringResult("", err),
			filepath.Join(dir, "poddescriptions.txt"))
		return
	}
	for _, pod := range strings.Split(strings.TrimSpace(res), "\n") {
		res := newErrorableStringResult(exechelper.ExecuteCommand("kubectl",
			"describe pod "+pod+" -n "+namespace))
		writeErrorableStringResult(res, filepath.Join(dir, pod+"_description.txt"))
	}
}

func writeDeploymentDescriptions(namespace, dir string) {
	fmt.Println("Retrieving deployment descriptions in " + namespace)
	res, err := exechelper.ExecuteCommand("kubectl",
		`get deployments --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' -n `+namespace)
	if err != nil {
		writeErrorableStringResult(newErrorableStringResult("", err),
			filepath.Join(dir, "deploymentdescriptions.txt"))
		return
	}
	for _, deployment := range strings.Split(strings.TrimSpace(res), "\n") {
		res := newErrorableStringResult(exechelper.ExecuteCommand("kubectl",
			"describe deployment "+deployment+" -n "+namespace))
		writeErrorableStringResult(res, filepath.Join(dir, deployment+"_description.txt"))
	}
}

func writePodLogs(namespace, dir string) {
	fmt.Println("Retrieving pod logs in " + namespace)
	res, err := exechelper.ExecuteCommand("kubectl",
		`get pods --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' -n `+namespace)
	if err != nil {
		writeErrorableStringResult(newErrorableStringResult("", err),
			filepath.Join(dir, "podlogs.txt"))
		return
	}
	for _, pod := range strings.Split(strings.TrimSpace(res), "\n") {
		res := newErrorableStringResult(exechelper.ExecuteCommand("kubectl", "logs "+pod+" --all-containers=true -n "+namespace))
		writeErrorableStringResult(res, filepath.Join(dir, pod+"_log.txt"))
	}
}

func getProjects() *errorableProjectResult {
	fmt.Println("Retrieving list of Keptn projects")
	endPoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
	if err != nil {
		return newErrorableProjectResult(nil, err)
	}
	projectHandler := apiutils.NewAuthenticatedProjectHandler(endPoint.String(), apiToken, "x-token", nil, *scheme)
	return newErrorableProjectResult(projectHandler.GetAllProjects())
}

func writeKeptnInstallerLog(logFileName string, dir string) {
	fmt.Println("Retrieving Keptn installer log " + logFileName)
	path, err := keptnutils.GetKeptnDirectory()
	if err != nil {
		writeErrorableStringResult(newErrorableStringResult("", err), filepath.Join(dir, logFileName))
		return
	}
	installerLog := filepath.Join(path, logFileName)
	res, err := ioutil.ReadFile(installerLog)
	writeErrorableStringResult(newErrorableStringResult(string(res), err), filepath.Join(dir, logFileName))
}

func getKeptnMetadata() *errorableMetadataResult {
	fmt.Println("Retrieving Keptn Metadata from API Service")
	endPoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
	if err != nil {
		return newErrorableMetadataResult(nil, err)
	}
	metadataHandler := apiutils.NewAuthenticatedAPIHandler(endPoint.String(), apiToken, "x-token", nil, *scheme)
	metadataData, errMetadata := metadataHandler.GetMetadata()
	if errMetadata != nil {
		err = errors.New("Error occured with response code " + string(errMetadata.Code) + " with message " + *errMetadata.Message)
	}
	return newErrorableMetadataResult(metadataData, err)
}

func init() {
	generateCmd.AddCommand(generateSupportArchiveCmd)

	generateSupportArchiveParams = &generateCmdParams{}
	generateSupportArchiveParams.Directory = generateSupportArchiveCmd.Flags().StringP("dir", "",
		"./support-archive", "directory where the docs should be written to")
}
