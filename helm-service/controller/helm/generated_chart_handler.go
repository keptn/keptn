package helm

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"sigs.k8s.io/yaml"
)

type GeneratedChartHandler struct {
	mesh           mesh.Mesh
	canaryLevelGen CanaryLevelGenerator
	keptnDomain    string
}

func NewGeneratedChartHandler(mesh mesh.Mesh, canaryLevelGen CanaryLevelGenerator,
	keptnDomain string) *GeneratedChartHandler {
	return &GeneratedChartHandler{mesh: mesh, canaryLevelGen: canaryLevelGen, keptnDomain: keptnDomain}
}

// GenerateDuplicateManagedChart generates a duplicated chart which is managed by keptn and used for
// b/g and canary releases
func (c *GeneratedChartHandler) GenerateDuplicateManagedChart(helmUpgradeMsg string, project string, stageName string, service string) (*chart.Chart, error) {

	if _, ok := c.canaryLevelGen.(*CanaryOnDeploymentGenerator); ok {

		meta := &chart.Metadata{
			Name:     service + "-generated",
			Keywords: []string{"deployment_strategy=" + keptnevents.Duplicate.String()},
			Version:  "0.1.0",
		}
		ch := chart.Chart{Metadata: meta}

		svcs, err := getServices(helmUpgradeMsg, project, stageName)
		if err != nil {
			return nil, err
		}
		depls, err := getDeployments(helmUpgradeMsg, project, stageName)
		if err != nil {
			return nil, err
		}

		for _, svc := range svcs {
			templates, err := c.generateServices(svc, project, stageName)
			if err != nil {
				return nil, err
			}
			ch.Templates = append(ch.Templates, templates...)
		}

		for _, depl := range depls {
			template, err := c.generateDeployment(depl)
			if err != nil {
				return nil, err
			}
			ch.Templates = append(ch.Templates, template)
		}

		return &ch, nil
	}
	log.Fatal("Currently canary is only supported on a deployment-level")
	return nil, nil
}

func getServices(helmUpgradeMsg string, project string, stageName string) ([]*corev1.Service, error) {

	namespace := project + "-" + stageName
	serviceNames, err := getServiceNames(helmUpgradeMsg)
	if err != nil {
		return nil, err
	}
	useInClusterConfig := false
	if os.Getenv("ENVIRONMENT") == "production" {
		useInClusterConfig = true
	}
	clientset, err := keptnutils.GetClientset(useInClusterConfig)
	if err != nil {
		return nil, err
	}
	services := []*corev1.Service{}
	for _, serviceName := range serviceNames {
		svc, err := clientset.CoreV1().Services(namespace).Get(serviceName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		services = append(services, svc)
	}
	return services, nil
}

func getServiceNames(helmUpgradeMsg string) ([]string, error) {
	serviceNames := []string{}
	startIdx := strings.Index(helmUpgradeMsg, "==> v1/Service")
	if startIdx > 0 {
		endIdx := strings.Index(helmUpgradeMsg[startIdx:], "\n\n")
		serviceBlock := strings.TrimSpace(helmUpgradeMsg[startIdx : startIdx+endIdx])
		lines := strings.Split(serviceBlock, "\n")
		if len(lines) < 3 {
			return nil, errors.New("Unexpected format of helm upgrade message")
		}
		for i := 2; i < len(lines); i++ {
			parts := strings.Split(lines[i], " ")
			serviceNames = append(serviceNames, strings.TrimSpace(parts[0]))
		}
	}
	return serviceNames, nil
}

func getDeployments(helmUpgradeMsg string, project string, stageName string) ([]*appsv1.Deployment, error) {
	namespace := project + "-" + stageName
	deploymentNames, err := getDeploymentNames(helmUpgradeMsg)
	if err != nil {
		return nil, err
	}
	useInClusterConfig := false
	if os.Getenv("ENVIRONMENT") == "production" {
		useInClusterConfig = true
	}
	clientset, err := keptnutils.GetClientset(useInClusterConfig)
	if err != nil {
		return nil, err
	}
	deployments := []*appsv1.Deployment{}
	for _, deplName := range deploymentNames {
		depl, err := clientset.AppsV1().Deployments(namespace).Get(deplName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		deployments = append(deployments, depl)
	}
	return deployments, nil
}

func getDeploymentNames(helmUpgradeMsg string) ([]string, error) {
	deploymentNames := []string{}
	startIdx := strings.Index(helmUpgradeMsg, "==> v1/Deployment")
	if startIdx > 0 {
		endIdx := strings.Index(helmUpgradeMsg[startIdx:], "\n\n")
		deploymentBlock := strings.TrimSpace(helmUpgradeMsg[startIdx : startIdx+endIdx])
		lines := strings.Split(deploymentBlock, "\n")
		if len(lines) < 3 {
			return nil, errors.New("Unexpected format of helm upgrade message")
		}
		for i := 2; i < len(lines); i++ {
			parts := strings.Split(lines[i], " ")
			deploymentNames = append(deploymentNames, strings.TrimSpace(parts[0]))
		}
	}
	return deploymentNames, nil
}

func resetService(svc *corev1.Service) {
	svc.Kind = "Service"
	svc.APIVersion = "v1"
	svc.Namespace = ""
	svc.Spec.ClusterIP = ""
	svc.Spec.LoadBalancerIP = ""
	svc.Spec.ExternalIPs = nil
	svc.Status = corev1.ServiceStatus{}
}

func resetDeployment(depl *appsv1.Deployment) {
	depl.Kind = "Deployment"
	depl.APIVersion = "apps/v1"
	depl.Namespace = ""
	depl.Status = appsv1.DeploymentStatus{}
}

func (c *GeneratedChartHandler) generateServices(svc *corev1.Service, project string, stageName string) ([]*chart.Template, error) {

	templates := make([]*chart.Template, 0, 0)

	serviceCanary := c.canaryLevelGen.GetCanaryService(*svc, project, stageName)
	resetService(serviceCanary)
	data, err := yaml.Marshal(serviceCanary)
	if err != nil {
		return nil, err
	}
	templates = append(templates, &chart.Template{Name: "templates/" + serviceCanary.Name + "-service" + ".yaml", Data: data})

	// Generate destination rule for canary service
	hostCanary := serviceCanary.Name + "." + c.canaryLevelGen.GetNamespace(project, stageName, true) + ".svc.cluster.local"
	destinationRuleCanary, err := c.mesh.GenerateDestinationRule(serviceCanary.Name, hostCanary)
	if err != nil {
		return nil, err
	}
	templates = append(templates, &chart.Template{Name: "templates/" + serviceCanary.Name + c.mesh.GetDestinationRuleSuffix(), Data: destinationRuleCanary})

	servicePrimary := svc.DeepCopy()
	servicePrimary.Name = servicePrimary.Name + "-primary"
	if _, ok := servicePrimary.Spec.Selector["app"]; ok {
		servicePrimary.Spec.Selector["app"] = servicePrimary.Spec.Selector["app"] + "-primary"
	}
	if _, ok := servicePrimary.Spec.Selector["app.kubernetes.io/name"]; ok {
		servicePrimary.Spec.Selector["app.kubernetes.io/name"] = servicePrimary.Spec.Selector["app.kubernetes.io/name"] + "-primary"
	}
	resetService(servicePrimary)
	data, err = yaml.Marshal(servicePrimary)
	if err != nil {
		return nil, err
	}
	templates = append(templates, &chart.Template{Name: "templates/" + servicePrimary.Name + "-service" + ".yaml", Data: data})

	// Generate destination rule for primary service
	hostPrimary := servicePrimary.Name + "." + c.canaryLevelGen.GetNamespace(project, stageName, true) + ".svc.cluster.local"
	destinationRulePrimary, err := c.mesh.GenerateDestinationRule(servicePrimary.Name, hostPrimary)
	if err != nil {
		return nil, err
	}
	templates = append(templates, &chart.Template{Name: "templates/" + servicePrimary.Name + c.mesh.GetDestinationRuleSuffix(), Data: destinationRulePrimary})

	// Generate virtual service
	gws := []string{GetGatewayName(project, stageName) + "." + GetUmbrellaNamespace(project, stageName), "mesh"}
	hosts := []string{
		svc.Name + "." + c.canaryLevelGen.GetNamespace(project, stageName, false) + "." + c.keptnDomain, // service_name.dev.123.45.67.89.xip.io
		svc.Name, // service-name
	}
	destCanary := mesh.HTTPRouteDestination{Host: hostCanary, Weight: 0}
	destPrimary := mesh.HTTPRouteDestination{Host: hostPrimary, Weight: 100}
	httpRouteDestinations := []mesh.HTTPRouteDestination{destCanary, destPrimary}

	vs, err := c.mesh.GenerateVirtualService(svc.Name, gws, hosts, httpRouteDestinations)
	if err != nil {
		return nil, err
	}

	templates = append(templates, &chart.Template{Name: "templates/" + svc.Name + c.mesh.GetVirtualServiceSuffix(), Data: vs})

	return templates, nil
}

func (c *GeneratedChartHandler) generateDeployment(depl *appsv1.Deployment) (*chart.Template, error) {
	primaryDeployment := depl.DeepCopy()

	primaryDeployment.Name = primaryDeployment.Name + "-primary"
	if _, ok := primaryDeployment.Spec.Selector.MatchLabels["app"]; ok {
		primaryDeployment.Spec.Selector.MatchLabels["app"] = primaryDeployment.Spec.Selector.MatchLabels["app"] + "-primary"
	}
	if _, ok := primaryDeployment.Spec.Selector.MatchLabels["app.kubernetes.io/name"]; ok {
		primaryDeployment.Spec.Selector.MatchLabels["app.kubernetes.io/name"] = primaryDeployment.Spec.Selector.MatchLabels["app.kubernetes.io/name"] + "-primary"
	}
	if _, ok := primaryDeployment.Spec.Template.ObjectMeta.Labels["app"]; ok {
		primaryDeployment.Spec.Template.ObjectMeta.Labels["app"] = primaryDeployment.Spec.Template.ObjectMeta.Labels["app"] + "-primary"
	}
	if _, ok := primaryDeployment.Spec.Template.ObjectMeta.Labels["app.kubernetes.io/name"]; ok {
		primaryDeployment.Spec.Template.ObjectMeta.Labels["app.kubernetes.io/name"] = primaryDeployment.Spec.Template.ObjectMeta.Labels["app.kubernetes.io/name"] + "-primary"
	}
	resetDeployment(primaryDeployment)
	data, err := yaml.Marshal(primaryDeployment)
	if err != nil {
		return nil, err
	}
	// Set the keptn_deployment to primary
	yamlString := strings.ReplaceAll(string(data), "keptn_deployment=canary", "keptn_deployment=primary")
	return &chart.Template{Name: "templates/" + primaryDeployment.Name + "-deployment" + ".yaml", Data: []byte(yamlString)}, nil
}

// GenerateMeshChart generates a chart containing the required mesh setup
func (c *GeneratedChartHandler) GenerateMeshChart(helmUpgradeMsg string, project string, stageName string,
	service string) (*chart.Chart, error) {

	namespace := project + "-" + stageName

	if _, ok := c.canaryLevelGen.(*CanaryOnDeploymentGenerator); ok {

		meta := &chart.Metadata{
			Name:     service + "-generated",
			Keywords: []string{"deployment_strategy=" + keptnevents.Direct.String()},
			Version:  "0.1.0",
		}
		ch := chart.Chart{Metadata: meta}

		svcs, err := getServices(helmUpgradeMsg, project, stageName)
		if err != nil {
			return nil, err
		}

		for _, svc := range svcs {
			// Generate virtual service for external access
			gws := []string{GetGatewayName(project, stageName) + "." + GetUmbrellaNamespace(project, stageName), "mesh"}
			hosts := []string{
				svc.Name + "." + namespace + "." + c.keptnDomain,
				svc.Name,
			}
			host := svc.Name + "." + namespace + ".svc.cluster.local"
			dest := mesh.HTTPRouteDestination{Host: host}
			httpRouteDestinations := []mesh.HTTPRouteDestination{dest}

			vs, err := c.mesh.GenerateVirtualService(svc.Name, gws, hosts, httpRouteDestinations)
			if err != nil {
				return nil, err
			}

			vsTemplate := chart.Template{Name: "templates/" + svc.Name + c.mesh.GetVirtualServiceSuffix(), Data: vs}
			ch.Templates = append(ch.Templates, &vsTemplate)
		}

		return &ch, nil
	}
	log.Fatal("Currently canary is only supported on a deployment-level")
	return nil, nil
}

// UpdateCanaryWeight updates the provided traffic weight in the VirtualService contained in the chart
func (c *GeneratedChartHandler) UpdateCanaryWeight(ch *chart.Chart, canaryWeight int32) error {

	// Set weights in all virtualservices
	for _, template := range ch.Templates {
		if strings.HasPrefix(template.Name, "templates/") &&
			strings.HasSuffix(template.Name, c.mesh.GetVirtualServiceSuffix()) {

			vs, err := c.mesh.UpdateWeights(template.Data, canaryWeight)
			if err != nil {
				return fmt.Errorf("Error when setting new weights in VirtualService %s: %s",
					template.Name, err.Error())
			}
			template.Data = vs
			break
		}
	}
	return nil
}

// GenerateEmptyChart generates an empty chart
func (c *GeneratedChartHandler) GenerateEmptyChart(project string, stageName string, service string,
	strategy keptnevents.DeploymentStrategy) *chart.Chart {

	if _, ok := c.canaryLevelGen.(*CanaryOnDeploymentGenerator); ok {

		meta := &chart.Metadata{
			Name:     service + "-generated",
			Keywords: []string{"deployment_strategy=" + strategy.String()},
			Version:  "0.1.0",
		}
		return &chart.Chart{Metadata: meta}
	}
	log.Fatal("Currently canary is only supported on a deployment-level")
	return nil
}
