package helm

import (
	"fmt"
	"strings"

	keptnevents "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	"helm.sh/helm/v3/pkg/chart"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type GeneratedChartHandler struct {
	mesh        mesh.Mesh
	keptnDomain string
}

func NewGeneratedChartHandler(mesh mesh.Mesh, keptnDomain string) *GeneratedChartHandler {
	return &GeneratedChartHandler{mesh: mesh, keptnDomain: keptnDomain}
}

// GenerateDuplicateManagedChart generates a duplicated chart which is managed by keptn and used for
// b/g and canary releases
func (c *GeneratedChartHandler) GenerateDuplicateManagedChart(helmManifest string, project string, stageName string, service string) (*chart.Chart, error) {
	meta := &chart.Metadata{
		Name:     service + "-generated",
		Keywords: []string{"deployment_strategy=" + keptnevents.Duplicate.String()},
		Version:  "0.1.0",
	}
	ch := chart.Chart{Metadata: meta}

	svcs := GetServices(helmManifest)
	depls := GetDeployments(helmManifest)

	for _, svc := range svcs {
		templates, err := c.generateServices(svc, project, stageName)
		if err != nil {
			continue
		}
		ch.Templates = append(ch.Templates, templates...)
	}

	for _, depl := range depls {
		template, err := c.generateDeployment(depl)
		if err != nil {
			fmt.Println(err)
			continue
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
	svc.ObjectMeta = metav1.ObjectMeta{Name: svc.Name}
	svc.Status = corev1.ServiceStatus{}
}

func resetDeployment(depl *appsv1.Deployment) {
	depl.Kind = "Deployment"
	depl.APIVersion = "apps/v1"
	depl.Namespace = ""
	depl.ResourceVersion = ""
	depl.ObjectMeta = metav1.ObjectMeta{Name: depl.Name}
	depl.Status = appsv1.DeploymentStatus{}
}

func (*GeneratedChartHandler) getNamespace(project string, stage string) string {
	return project + "-" + stage
}

func (c *GeneratedChartHandler) generateServices(svc *corev1.Service, project string, stageName string) ([]*chart.Template, error) {

	templates := make([]*chart.File, 0, 0)

	serviceCanary := svc.DeepCopy()
	serviceCanary.Name = serviceCanary.Name + "-canary"

	resetService(serviceCanary)
	data, err := yaml.Marshal(serviceCanary)
	if err != nil {
		return nil, err
	}
	templates = append(templates, &chart.File{Name: "templates/" + serviceCanary.Name + "-service" + ".yaml", Data: data})

	// Generate destination rule for canary service
	hostCanary := serviceCanary.Name + "." + c.getNamespace(project, stageName) + ".svc.cluster.local"
	destinationRuleCanary, err := c.mesh.GenerateDestinationRule(serviceCanary.Name, hostCanary)
	if err != nil {
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
	hostPrimary := servicePrimary.Name + "." + c.getNamespace(project, stageName) + ".svc.cluster.local"
	destinationRulePrimary, err := c.mesh.GenerateDestinationRule(servicePrimary.Name, hostPrimary)
	if err != nil {
		return nil, err
	}
	templates = append(templates, &chart.File{Name: "templates/" + servicePrimary.Name + c.mesh.GetDestinationRuleSuffix(), Data: destinationRulePrimary})

	// Generate virtual service
	gws := []string{"public-gateway.istio-system", "mesh"}
	hosts := []string{
		svc.Name + "." + c.getNamespace(project, stageName) + "." + c.keptnDomain, // service_name.dev.123.45.67.89.xip.io
		svc.Name, // service-name
	}
	destCanary := mesh.HTTPRouteDestination{Host: hostCanary, Weight: 0}
	destPrimary := mesh.HTTPRouteDestination{Host: hostPrimary, Weight: 100}
	httpRouteDestinations := []mesh.HTTPRouteDestination{destCanary, destPrimary}

	vs, err := c.mesh.GenerateVirtualService(svc.Name, gws, hosts, httpRouteDestinations)
	if err != nil {
		return nil, err
	}

	templates = append(templates, &chart.File{Name: "templates/" + svc.Name + c.mesh.GetVirtualServiceSuffix(), Data: vs})

	return templates, nil
}

func (c *GeneratedChartHandler) generateDeployment(depl *appsv1.Deployment) (*chart.File, error) {
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
func (c *GeneratedChartHandler) GenerateMeshChart(helmManifest string, project string, stageName string,
	service string) (*chart.Chart, error) {

	meta := &chart.Metadata{
		Name:     service + "-generated",
		Keywords: []string{"deployment_strategy=" + keptnevents.Direct.String()},
		Version:  "0.1.0",
	}
	ch := chart.Chart{Metadata: meta}

	svcs := GetServices(helmManifest)

	for _, svc := range svcs {
		// Generate virtual service for external access
		gws := []string{"public-gateway.istio-system", "mesh"}
		hosts := []string{
			svc.Name + "." + c.getNamespace(project, stageName) + "." + c.keptnDomain,
			svc.Name,
		}
		host := svc.Name + "." + c.getNamespace(project, stageName) + ".svc.cluster.local"
		dest := mesh.HTTPRouteDestination{Host: host}
		httpRouteDestinations := []mesh.HTTPRouteDestination{dest}

		vs, err := c.mesh.GenerateVirtualService(svc.Name, gws, hosts, httpRouteDestinations)
		if err != nil {
			return nil, err
		}

		vsTemplate := chart.Template{Name: "templates/" + svc.Name + c.mesh.GetVirtualServiceSuffix(), Data: vs}
		ch.Templates = append(ch.Templates, &vsTemplate)

		dr, err := c.mesh.GenerateDestinationRule(svc.Name, host)
		if err != nil {
			return nil, err
		}
		drTemplate := chart.Template{Name: "templates/" + svc.Name + c.mesh.GetDestinationRuleSuffix(), Data: dr}
		ch.Templates = append(ch.Templates, &drTemplate)
	}

	return &ch, nil
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
			APIVersion: "v2",
			Name:       service + "-generated",
			Keywords:   []string{"deployment_strategy=" + strategy.String()},
			Version:    "0.1.0",
		}
		return &chart.Chart{Metadata: meta}
	}
	return &chart.Chart{Metadata: meta}

}
