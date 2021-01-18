package helm

import (
	"errors"
	"fmt"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"net/url"
	"strings"

	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"

	keptnevents "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/helm-service/pkg/mesh"
	"helm.sh/helm/v3/pkg/chart"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

// ChartGenerator ...
type ChartGenerator interface {
	// GenerateDuplicateChart generates a duplicated chart which is managed by keptn and used for
	// b/g and canary releases
	GenerateDuplicateChart(helmManifest string, project string, stageName string, service string) (*chart.Chart, error)
	// GenerateMeshChart generates a chart containing the required mesh setup
	GenerateMeshChart(helmManifest string, project string, stageName string, service string) (*chart.Chart, error)
}

// GeneratedChartGenerator allows to generate the generated-chart
type GeneratedChartGenerator struct {
	mesh   mesh.Mesh
	logger keptncommon.LoggerInterface
}

// NewGeneratedChartGenerator creates a new GeneratedChartGenerator
func NewGeneratedChartGenerator(mesh mesh.Mesh, logger keptncommon.LoggerInterface) *GeneratedChartGenerator {
	return &GeneratedChartGenerator{
		mesh:   mesh,
		logger: logger,
	}
}

// GenerateDuplicateChart generates a duplicated chart which is managed by keptn and used for
// b/g and canary releases
func (c *GeneratedChartGenerator) GenerateDuplicateChart(helmManifest string, project string, stageName string, service string) (*chart.Chart, error) {
	meta := &chart.Metadata{
		APIVersion: "v2",
		Name:       service + "-generated",
		Keywords:   []string{"deployment_strategy=" + keptnevents.Duplicate.String()},
		Version:    "0.1.0",
	}
	ch := chart.Chart{Metadata: meta}

	svcs := GetServices(helmManifest)
	depls := GetDeployments(helmManifest)

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

func resetService(svc *corev1.Service) {
	svc.Kind = "Service"
	svc.APIVersion = "v1"
	svc.Namespace = ""
	svc.Spec.ClusterIP = ""
	svc.Spec.LoadBalancerIP = ""
	svc.Spec.ExternalIPs = nil
	for idx := range svc.Spec.Ports {
		svc.Spec.Ports[idx].NodePort = 0
	}
	svc.Status = corev1.ServiceStatus{}
}

func resetDeployment(depl *appsv1.Deployment) {
	depl.Kind = "Deployment"
	depl.APIVersion = "apps/v1"
	depl.Namespace = ""
	depl.ResourceVersion = ""
	depl.Status = appsv1.DeploymentStatus{}
}

func (*GeneratedChartGenerator) getNamespace(project string, stage string) string {
	return project + "-" + stage
}

func (c *GeneratedChartGenerator) generateServices(svc *corev1.Service, project string, stageName string) ([]*chart.File, error) {

	templates := make([]*chart.File, 0, 0)

	serviceCanary := svc.DeepCopy()
	serviceCanary.Name = serviceCanary.Name + "-canary"

	c.logger.Info("Generating canary service for " + svc.Name + ": " + serviceCanary.Name)
	resetService(serviceCanary)
	data, err := yaml.Marshal(serviceCanary)
	if err != nil {
		return nil, err
	}
	templates = append(templates, &chart.File{Name: "templates/" + serviceCanary.Name + "-service" + ".yaml", Data: data})

	// Generate destination rule for canary service
	c.logger.Info("Generating destination rule for canary service " + serviceCanary.Name)
	hostCanary := serviceCanary.Name + "." + c.getNamespace(project, stageName) + ".svc.cluster.local"
	destinationRuleCanary, err := c.mesh.GenerateDestinationRule(serviceCanary.Name, hostCanary)
	if err != nil {
		c.logger.Error("Error while generating destination rule for canary service " + serviceCanary.Name + ": " + err.Error())
		return nil, err
	}
	templates = append(templates, &chart.File{Name: "templates/" + serviceCanary.Name + c.mesh.GetDestinationRuleSuffix(), Data: destinationRuleCanary})

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
	templates = append(templates, &chart.File{Name: "templates/" + servicePrimary.Name + "-service" + ".yaml", Data: data})

	// Generate destination rule for primary service
	c.logger.Info("Generating destination rule for primary service " + svc.Name)
	hostPrimary := servicePrimary.Name + "." + c.getNamespace(project, stageName) + ".svc.cluster.local"
	destinationRulePrimary, err := c.mesh.GenerateDestinationRule(servicePrimary.Name, hostPrimary)
	if err != nil {
		c.logger.Error("Error while generating destination rule for primary service " + svc.Name + ": " + err.Error())
		return nil, err
	}
	templates = append(templates, &chart.File{Name: "templates/" + servicePrimary.Name + c.mesh.GetDestinationRuleSuffix(), Data: destinationRulePrimary})

	// get the public hostname based on what has been configured in HOSTNAME_TEMPLATE and INGRESS_HOSTNAME_SUFFIX
	publicHostName, err := getVirtualServicePublicHost(svc.Name, project, stageName)
	if err != nil {
		return nil, err
	}

	// Generate virtual service
	gws := []string{mesh.GetIngressGateway(), "mesh"}
	hosts := []string{
		publicHostName, // service_name.dev.123.45.67.89.xip.io
		svc.Name,       // service-name
	}
	destCanary := mesh.HTTPRouteDestination{Host: hostCanary, Weight: 0}
	destPrimary := mesh.HTTPRouteDestination{Host: hostPrimary, Weight: 100}
	httpRouteDestinations := []mesh.HTTPRouteDestination{destCanary, destPrimary}

	c.logger.Info("Generating VirtualService for service " + svc.Name + ". URL = " + mesh.GetIngressProtocol() +
		"://" + hosts[0] + ":" + mesh.GetIngressPort())
	vs, err := c.mesh.GenerateVirtualService(svc.Name, gws, hosts, httpRouteDestinations)
	if err != nil {
		c.logger.Error("Error while generating VirtualService for service " + svc.Name + ": " + err.Error())
		return nil, err
	}

	templates = append(templates, &chart.File{Name: "templates/" + svc.Name + c.mesh.GetVirtualServiceSuffix(), Data: vs})

	return templates, nil
}

func getVirtualServicePublicHost(serviceName, projectName, stageName string) (string, error) {
	publicURI := mesh.GetPublicDeploymentURI(keptnv2.EventData{
		Project: projectName,
		Stage:   stageName,
		Service: serviceName,
	})

	if len(publicURI) == 0 {
		return "", errors.New("could not determine public host name")
	}
	url, err := url.Parse(publicURI[0])
	if err != nil {
		return "", fmt.Errorf("could not parse hostname: " + err.Error())
	}
	if url.Hostname() == "" && url.Path != "" {
		return "", errors.New("Missing leading protocol (e.g. http://) in HOSTNAME_TEMPLATE environment variable")
	}
	return url.Hostname(), nil
}

func (c *GeneratedChartGenerator) generateDeployment(depl *appsv1.Deployment) (*chart.File, error) {
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
	return &chart.File{Name: "templates/" + primaryDeployment.Name + "-deployment" + ".yaml", Data: []byte(yamlString)}, nil
}

// GenerateMeshChart generates a chart containing the required mesh setup
func (c *GeneratedChartGenerator) GenerateMeshChart(helmManifest string, project string, stageName string,
	service string) (*chart.Chart, error) {

	meta := &chart.Metadata{
		APIVersion: "v2",
		Name:       service + "-generated",
		Keywords:   []string{"deployment_strategy=" + keptnevents.Direct.String()},
		Version:    "0.1.0",
	}
	ch := chart.Chart{Metadata: meta}

	svcs := GetServices(helmManifest)

	for _, svc := range svcs {
		// get the public hostname based on what has been configured in HOSTNAME_TEMPLATE and INGRESS_HOSTNAME_SUFFIX
		publicHostName, err := getVirtualServicePublicHost(svc.Name, project, stageName)
		if err != nil {
			return nil, err
		}

		// Generate virtual service for external access
		gws := []string{mesh.GetIngressGateway(), "mesh"}
		hosts := []string{
			publicHostName,
			svc.Name,
		}
		host := svc.Name + "." + c.getNamespace(project, stageName) + ".svc.cluster.local"
		dest := mesh.HTTPRouteDestination{Host: host}
		httpRouteDestinations := []mesh.HTTPRouteDestination{dest}

		vs, err := c.mesh.GenerateVirtualService(svc.Name, gws, hosts, httpRouteDestinations)
		if err != nil {
			return nil, err
		}

		vsTemplate := chart.File{Name: "templates/" + svc.Name + c.mesh.GetVirtualServiceSuffix(), Data: vs}
		ch.Templates = append(ch.Templates, &vsTemplate)

		dr, err := c.mesh.GenerateDestinationRule(svc.Name, host)
		if err != nil {
			return nil, err
		}
		drTemplate := chart.File{Name: "templates/" + svc.Name + c.mesh.GetDestinationRuleSuffix(), Data: dr}
		ch.Templates = append(ch.Templates, &drTemplate)
	}

	return &ch, nil
}
