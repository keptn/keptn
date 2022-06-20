package common

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/keptn/go-utils/pkg/api/models"
	goutils "github.com/keptn/go-utils/pkg/api/utils"
	utils "github.com/keptn/go-utils/pkg/api/utils"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	appsv1 "k8s.io/api/apps/v1"
	typesv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apierr "k8s.io/apimachinery/pkg/api/errors"
)

//chartStorer  is able to store a helm chart
type chartStorer struct {
	resourceHandler *goutils.ResourceHandler
}

type StoreChartOptions struct {
	Project   string
	Service   string
	Stage     string
	ChartName string
	HelmChart []byte
}

//NewChartStorer creates a new chartStorer instance
func NewChartStorer(resourceHandler *goutils.ResourceHandler) *chartStorer {
	return &chartStorer{
		resourceHandler: resourceHandler,
	}
}

func getHelmChartURI(chartName string) string {
	return "helm/" + chartName + ".tgz"
}

//Store stores a chart in the configuration service
func (cs chartStorer) Store(storeChartOpts StoreChartOptions) (string, error) {

	uri := getHelmChartURI(storeChartOpts.ChartName)
	resource := models.Resource{ResourceURI: &uri, ResourceContent: string(storeChartOpts.HelmChart)}

	version, err := cs.resourceHandler.CreateServiceResources(storeChartOpts.Project, storeChartOpts.Stage, storeChartOpts.Service, []*models.Resource{&resource})
	if err != nil {
		return "", fmt.Errorf("Error when storing chart %s of service %s in project %s: %s",
			storeChartOpts.ChartName, storeChartOpts.Service, storeChartOpts.Project, err.Error())
	}
	return version, nil
}

//chartPackager is able to package a helm chart
type chartPackager struct {
}

//NewChartPackager creates a new chartPackager instance
func NewChartPackager() *chartPackager {
	return &chartPackager{}
}

//packages a helm chart into its byte representation
func (pc chartPackager) Package(ch *chart.Chart) ([]byte, error) {
	helmPackage, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, fmt.Errorf("Error when packaging chart: %s", err.Error())
	}
	defer os.RemoveAll(helmPackage)

	// Marshal values into values.yaml
	// This step is necessary as chartutil.Save uses the Raw content
	for _, f := range ch.Raw {
		if f.Name == chartutil.ValuesfileName {
			f.Data, err = yaml.Marshal(ch.Values)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	name, err := chartutil.Save(ch, helmPackage)
	if err != nil {
		return nil, fmt.Errorf("Error when packaging chart: %s", err.Error())
	}

	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("Error when packaging chart: %s", err.Error())
	}
	return data, nil
}

// GetClientset returns the kubernetes Clientset
func GetClientset(useInClusterConfig bool) (*kubernetes.Clientset, error) {

	var config *rest.Config
	var err error
	if useInClusterConfig {
		config, err = rest.InClusterConfig()
	} else {
		var kubeconfig string
		if os.Getenv("KUBECONFIG") != "" {
			kubeconfig = ExpandTilde(os.Getenv("KUBECONFIG"))
		} else {
			kubeconfig = filepath.Join(
				UserHomeDir(), ".kube", "config",
			)
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

// UserHomeDir returns the HOME directory by taking into account the operating system
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// ExpandTilde expands ~ to HOME
func ExpandTilde(fileName string) string {
	if fileName == "~" {
		return UserHomeDir()
	} else if strings.HasPrefix(fileName, "~/") {
		return filepath.Join(UserHomeDir(), fileName[2:])
	}
	return fileName
}

//chartRetriever is able to store a helm chart
type chartRetriever struct {
	resourceHandler *goutils.ResourceHandler
}

//RetrieveChartOptions are the parameters to obtain a chart
type RetrieveChartOptions struct {
	Project   string
	Service   string
	Stage     string
	ChartName string
	CommitID  string
}

//NewChartRetriever creates a new chartRetriever instance
func NewChartRetriever(resourceHandler *goutils.ResourceHandler) *chartRetriever {
	return &chartRetriever{
		resourceHandler: resourceHandler,
	}
}

func (cs chartRetriever) Retrieve(chartOpts RetrieveChartOptions) (*chart.Chart, string, error) {
	option := url.Values{}
	if chartOpts.CommitID != "" {
		option.Add("gitCommitID", chartOpts.CommitID)
	}
	resource, err := cs.resourceHandler.GetResource(
		*goutils.NewResourceScope().
			Project(chartOpts.Project).
			Service(chartOpts.Service).
			Resource(getHelmChartURI(chartOpts.ChartName)).
			Stage(chartOpts.Stage),
		goutils.AppendQuery(option))

	if err != nil {
		return nil, "", fmt.Errorf("Error when reading chart %s from project %s: %s",
			chartOpts.ChartName, chartOpts.Project, err.Error())
	}
	ch, err := LoadChart([]byte(resource.ResourceContent))
	if err != nil {
		return nil, "", fmt.Errorf("Error when reading chart %s from project %s: %s",
			chartOpts.ChartName, chartOpts.Project, err.Error())
	}
	if chartOpts.CommitID == "" {
		return ch, resource.Metadata.Version, nil
	}

	return ch, chartOpts.CommitID, nil
}

// LoadChart converts a byte array into a Chart
func LoadChart(data []byte) (*chart.Chart, error) {
	return loader.LoadArchive(bytes.NewReader(data))
}

// StoreChart stores a chart in the configuration service
//Deprecated: StoreChart is deprecated, use chartStorer.Store instead
func StoreChart(project string, service string, stage string, chartName string, helmChart []byte, configServiceURL string) (string, error) {

	cs := chartStorer{
		resourceHandler: utils.NewResourceHandler(configServiceURL),
	}

	opts := StoreChartOptions{
		Project:   project,
		Service:   service,
		Stage:     stage,
		ChartName: chartName,
		HelmChart: helmChart,
	}
	return cs.Store(opts)

}

// PackageChart packages the chart and returns it
//Deprecated: PackageChart is deprecated, use chartPackager.Package instead
func PackageChart(ch *chart.Chart) ([]byte, error) {
	cp := chartPackager{}
	return cp.Package(ch)
}

// IsDeployment tests whether the provided struct is a deployment
func IsDeployment(dpl *appsv1.Deployment) bool {
	return strings.ToLower(dpl.Kind) == "deployment"
}

// IsService tests whether the provided struct is a service
func IsService(svc *typesv1.Service) bool {
	return strings.ToLower(svc.Kind) == "service"
}

// WaitForDeploymentToBeRolledOut waits until the deployment is Available
func WaitForDeploymentToBeRolledOut(useInClusterConfig bool, deploymentName string, namespace string) error {
	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return err
	}

	const maxWaitForDeploymentRetries = 90
	deployment, err := getDeployment(clientset, namespace, deploymentName)
	retries := 0
	for {

		var cond *appsv1.DeploymentCondition

		for i := range deployment.Status.Conditions {
			c := deployment.Status.Conditions[i]
			if c.Type == appsv1.DeploymentProgressing {
				cond = &c
				break
			}
		}

		if cond != nil && cond.Reason == "ProgressDeadlineExceeded" {
			return fmt.Errorf("Deployment %q exceeded its progress deadline", deployment.Name)
		}
		if !(deployment.Spec.Replicas != nil && deployment.Status.UpdatedReplicas < *deployment.Spec.Replicas ||
			deployment.Status.Replicas > deployment.Status.UpdatedReplicas ||
			deployment.Status.AvailableReplicas < deployment.Status.UpdatedReplicas) {
			return nil
		}

		time.Sleep(2 * time.Second)
		deployment, err = getDeployment(clientset, namespace, deploymentName)
		if err != nil {
			return err
		}
		retries = retries + 1
		if retries >= maxWaitForDeploymentRetries {
			return fmt.Errorf("Timed out waiting for deployment %q", deployment.Name)
		}
	}
}

func getDeployment(clientset *kubernetes.Clientset, namespace string, deploymentName string) (*appsv1.Deployment, error) {
	dep, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil &&
		strings.Contains(err.Error(), "the object has been modified; please apply your changes to the latest version and try again") {
		time.Sleep(10 * time.Second)
		return clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	}
	return dep, nil
}

// CreateNamespace creates a new Kubernetes namespace with the provided name
func CreateNamespace(useInClusterConfig bool, namespace string, namespaceMetadata ...metav1.ObjectMeta) error {

	var buildNamespaceMetadata metav1.ObjectMeta
	if len(namespaceMetadata) > 0 {
		buildNamespaceMetadata = namespaceMetadata[0]
	}

	buildNamespaceMetadata.Name = namespace

	ns := &typesv1.Namespace{ObjectMeta: buildNamespaceMetadata}
	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return err
	}
	_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	return err
}

// ExistsNamespace checks whether a namespace with the provided name exists
func ExistsNamespace(useInClusterConfig bool, namespace string) (bool, error) {
	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return false, err
	}
	_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*apierr.StatusError); ok && statusErr.ErrStatus.Reason == metav1.StatusReasonNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetKubeAPI returns the CoreV1Interface
func GetKubeAPI(useInClusterConfig bool) (v1.CoreV1Interface, error) {

	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1(), nil
}
