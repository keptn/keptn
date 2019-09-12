package helm

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	"github.com/keptn/keptn/helm-service/pkg/objectutils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/renderutil"
	"k8s.io/helm/pkg/timeconv"
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
func (c *GeneratedChartHandler) GenerateDuplicateManagedChart(ch *chart.Chart, project string, stageName string) ([]byte, error) {

	c.changeChartFile(ch)

	if err := c.changeTemplateContent(project, ch, stageName); err != nil {
		return nil, err
	}

	return keptnutils.PackageChart(ch)
}

// GenerateMeshChart generates a chart containing the required mesh setup
func (c *GeneratedChartHandler) GenerateMeshChart(ch *chart.Chart, project string, stageName string) ([]byte, error) {

	c.changeChartFile(ch)
	newTemplates := make([]*chart.Template, 0, 0)

	services, err := keptnutils.GetRenderedServices(ch)
	if err != nil {
		return nil, err
	}
	namespace := project + "-" + stageName

	for _, svc := range services {
		// Generate virtual service
		gws := []string{GetGatewayName(project, stageName) + "." + GetUmbrellaNamespace(project, stageName), "mesh"}
		hosts := []string{svc.Name + "." + namespace + "." + c.keptnDomain,
			svc.Name, svc.Name + "." + namespace}
		host := svc.Name + "." + namespace + ".svc.cluster.local"
		dest := mesh.HTTPRouteDestination{Host: host}
		httpRouteDestinations := []mesh.HTTPRouteDestination{dest}

		vs, err := c.mesh.GenerateVirtualService(svc.Name, gws, hosts, httpRouteDestinations)
		if err != nil {
			return nil, err
		}

		vsTemplate := chart.Template{Name: "templates/" + svc.Name + c.mesh.GetVirtualServiceSuffix(), Data: vs}
		newTemplates = append(newTemplates, &vsTemplate)
	}

	ch.Values.Raw = ""
	ch.Values.Values = make(map[string]*chart.Value)
	ch.Templates = newTemplates
	return keptnutils.PackageChart(ch)
}

func (c *GeneratedChartHandler) changeChartFile(ch *chart.Chart) {

	ch.Metadata.Name = ch.Metadata.Name + "-generated"
	ch.Metadata.Description = ch.Metadata.Description + " (generated)"
}

func (c *GeneratedChartHandler) changeTemplateContent(project string,
	ch *chart.Chart, stageName string) error {

	renderOpts := renderutil.Options{
		ReleaseOptions: chartutil.ReleaseOptions{
			Name:      ch.Metadata.Name,
			IsInstall: false,
			IsUpgrade: false,
			Time:      timeconv.Now(),
		},
	}

	renderedTemplates, err := renderutil.Render(ch, ch.Values, renderOpts)
	if err != nil {
		return err
	}

	newTemplates := make([]*chart.Template, 0, 0)

	for k, v := range renderedTemplates {
		dec := kyaml.NewYAMLToJSONDecoder(strings.NewReader(v))
		newContent := make([]byte, 0, 0)
		for {
			var document interface{}
			err := dec.Decode(&document)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			doc, err := json.Marshal(document)
			if err != nil {
				return err
			}

			newServiceTemplateContent, newServiceTemplates, err := c.handleService(doc, project, stageName)
			if err != nil {
				return err
			}
			if len(newServiceTemplateContent) > 0 {
				newContent = append(newContent, newServiceTemplateContent...)
				newTemplates = append(newTemplates, newServiceTemplates...)
				continue
			}

			newDeploymentTemplateContent, err := c.handleDeployment(doc)
			if err != nil {
				return err
			}
			if len(newDeploymentTemplateContent) > 0 {
				newContent = append(newContent, newDeploymentTemplateContent...)
				continue
			}

			if c.canaryLevelGen.IsK8sResourceDuplicated() {
				newContent, err = objectutils.AppendAsYaml(newContent, document)
				if err != nil {
					return err
				}
			}
		}

		if len(newContent) > 0 {
			newTemplates = append(newTemplates, &chart.Template{Name: k[len(ch.Metadata.Name)+1:], Data: newContent})
		}
	}

	ch.Templates = newTemplates
	return nil
}

func (c *GeneratedChartHandler) handleService(document []byte, project string, stageName string) ([]byte, []*chart.Template, error) {

	var svc corev1.Service
	if err := json.Unmarshal(document, &svc); err != nil {
		// Ignore unmarshaling error
		return nil, nil, nil
	}

	newTemplateContent := make([]byte, 0, 0)
	newTemplates := make([]*chart.Template, 0, 0)

	if keptnutils.IsService(&svc) {
		var err error

		serviceCanary := c.canaryLevelGen.GetCanaryService(svc, project, stageName)

		newTemplateContent, err = objectutils.AppendAsYaml(newTemplateContent, serviceCanary)
		if err != nil {
			return nil, nil, err
		}

		// Generate destination rule for canary service
		hostCanary := serviceCanary.Name + "." + c.canaryLevelGen.GetNamespace(project, stageName, true) + ".svc.cluster.local"
		destinationRuleCanary, err := c.mesh.GenerateDestinationRule(serviceCanary.Name, hostCanary)
		if err != nil {
			return nil, nil, err
		}

		d1 := chart.Template{Name: "templates/" + serviceCanary.Name + c.mesh.GetDestinationRuleSuffix(), Data: destinationRuleCanary}
		newTemplates = append(newTemplates, &d1)

		servicePrimary := svc.DeepCopy()
		servicePrimary.Name = servicePrimary.Name + "-primary"
		servicePrimary.Spec.Selector["app"] = servicePrimary.Spec.Selector["app"] + "-primary"
		newTemplateContent, err = objectutils.AppendAsYaml(newTemplateContent, servicePrimary)
		if err != nil {
			return nil, nil, err
		}

		// Generate destination rule for primary service
		hostPrimary := servicePrimary.Name + "." + c.canaryLevelGen.GetNamespace(project, stageName, true) + ".svc.cluster.local"
		destinationRulePrimary, err := c.mesh.GenerateDestinationRule(servicePrimary.Name, hostPrimary)
		if err != nil {
			return nil, nil, err
		}

		d2 := chart.Template{Name: "templates/" + servicePrimary.Name + c.mesh.GetDestinationRuleSuffix(), Data: destinationRulePrimary}
		newTemplates = append(newTemplates, &d2)

		// Generate virtual service
		gws := []string{GetGatewayName(project, stageName) + "." + GetUmbrellaNamespace(project, stageName), "mesh"}
		hosts := []string{svc.Name + "." + c.canaryLevelGen.GetNamespace(project, stageName, false) + "." + c.keptnDomain,
			svc.Name, svc.Name + "." + c.canaryLevelGen.GetNamespace(project, stageName, true)}
		destCanary := mesh.HTTPRouteDestination{Host: hostCanary, Weight: 0}
		destPrimary := mesh.HTTPRouteDestination{Host: hostPrimary, Weight: 100}
		httpRouteDestinations := []mesh.HTTPRouteDestination{destCanary, destPrimary}

		vs, err := c.mesh.GenerateVirtualService(svc.Name, gws, hosts, httpRouteDestinations)
		if err != nil {
			return nil, nil, err
		}

		vsTemplate := chart.Template{Name: "templates/" + svc.Name + c.mesh.GetVirtualServiceSuffix(), Data: vs}
		newTemplates = append(newTemplates, &vsTemplate)
	}
	return newTemplateContent, newTemplates, nil
}

func (c *GeneratedChartHandler) handleDeployment(document []byte) ([]byte, error) {
	// Try to unmarshal Deployment
	var depl appsv1.Deployment
	if err := json.Unmarshal(document, &depl); err != nil {
		// Ignore unmarshaling error
		return nil, nil
	}

	newTemplateContent := make([]byte, 0, 0)
	if keptnutils.IsDeployment(&depl) {
		depl.Name = depl.Name + "-primary"
		depl.Spec.Selector.MatchLabels["app"] = depl.Spec.Selector.MatchLabels["app"] + "-primary"
		depl.Spec.Template.ObjectMeta.Labels["app"] = depl.Spec.Template.ObjectMeta.Labels["app"] + "-primary"
		var err error
		newTemplateContent, err = objectutils.AppendAsYaml(newTemplateContent, depl)
		if err != nil {
			return nil, err
		}
	}
	return newTemplateContent, nil
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
