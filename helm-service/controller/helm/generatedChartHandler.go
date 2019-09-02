package helm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	keptnevents "github.com/keptn/go-utils/pkg/events"
	"github.com/keptn/keptn/helm-service/controller/jsonutils"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type GeneratedChartHandler struct {
	mesh             mesh.Mesh
	canaryLevelGen   CanaryLevelGenerator
	keptnDomain      string
	configServiceURL string
}

func NewGeneratedChartHandler(mesh mesh.Mesh, canaryLevelGen CanaryLevelGenerator,
	keptnDomain string, configServiceURL string) *GeneratedChartHandler {
	return &GeneratedChartHandler{mesh: mesh, canaryLevelGen: canaryLevelGen, keptnDomain: keptnDomain, configServiceURL: configServiceURL}
}

// GenerateManagedChart generates a duplicated chart which is managed by keptn and used for
// b/g and canary releases
func (c *GeneratedChartHandler) GenerateManagedChart(event *keptnevents.ServiceCreateEventData, stageName string) ([]byte, error) {

	ch, err := LoadChart(event.HelmChart)
	if err != nil {
		return nil, err
	}

	c.changeChartFile(ch)

	if err := c.changeTemplateContent(event, ch, stageName); err != nil {
		return nil, err
	}

	return PackageChart(ch)
}

func (c *GeneratedChartHandler) changeChartFile(ch *chart.Chart) {

	ch.Metadata.Name = ch.Metadata.Name + "-generated"
	ch.Metadata.Description = ch.Metadata.Description + " (generated)"
}

func (c *GeneratedChartHandler) changeTemplateContent(event *keptnevents.ServiceCreateEventData,
	ch *chart.Chart, stageName string) error {

	newTemplates := make([]*chart.Template, 0, 0)

	for _, templateFile := range ch.Templates {
		dec := kyaml.NewYAMLToJSONDecoder(bytes.NewReader(templateFile.Data))
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

			newServiceTemplateContent, newServiceTemplates, err := c.handleService(doc, event, stageName)
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
				newContent, err = appendAsYaml(newContent, document)
				if err != nil {
					return err
				}
			}
		}

		templateFile.Data = newContent
	}

	ch.Templates = append(ch.Templates, newTemplates...)
	return nil
}

func appendAsYaml(content []byte, element interface{}) ([]byte, error) {

	jsonData, err := json.Marshal(element)
	if err != nil {
		return nil, err
	}
	yamlData, err := jsonutils.ToYAML(jsonData)
	if err != nil {
		return nil, err
	}
	content = append(content, []byte("---\n")...)
	return append(content, yamlData...), nil
}

func (c *GeneratedChartHandler) handleService(document []byte, event *keptnevents.ServiceCreateEventData, stageName string) ([]byte, []*chart.Template, error) {

	var svc corev1.Service
	if err := json.Unmarshal(document, &svc); err != nil {
		// Ignore unmarshaling error
		return nil, nil, nil
	}

	newTemplateContent := make([]byte, 0, 0)
	newTemplates := make([]*chart.Template, 0, 0)

	if IsService(&svc) {
		var err error

		serviceCanary := c.canaryLevelGen.GetCanaryService(svc, event, stageName)

		newTemplateContent, err = appendAsYaml(newTemplateContent, serviceCanary)
		if err != nil {
			return nil, nil, err
		}

		// Generate destination rule for canary service
		hostCanary := serviceCanary.Name + "." + c.canaryLevelGen.GetNamespace(event.Project, stageName, true) + ".svc.cluster.local"
		destinationRuleCanary, err := c.mesh.GenerateDestinationRule(serviceCanary.Name, hostCanary)
		if err != nil {
			return nil, nil, err
		}

		d1 := chart.Template{Name: "templates/" + serviceCanary.Name + c.mesh.GetDestinationRuleSuffix(), Data: destinationRuleCanary}
		newTemplates = append(newTemplates, &d1)

		servicePrimary := svc.DeepCopy()
		servicePrimary.Name = servicePrimary.Name + "-primary"
		servicePrimary.Spec.Selector["app"] = servicePrimary.Spec.Selector["app"] + "-primary"
		newTemplateContent, err = appendAsYaml(newTemplateContent, servicePrimary)
		if err != nil {
			return nil, nil, err
		}

		// Generate destination rule for primary service
		hostPrimary := servicePrimary.Name + "." + c.canaryLevelGen.GetNamespace(event.Project, stageName, true) + ".svc.cluster.local"
		destinationRulePrimary, err := c.mesh.GenerateDestinationRule(servicePrimary.Name, hostPrimary)
		if err != nil {
			return nil, nil, err
		}

		d2 := chart.Template{Name: "templates/" + servicePrimary.Name + c.mesh.GetDestinationRuleSuffix(), Data: destinationRulePrimary}
		newTemplates = append(newTemplates, &d2)

		gws := []string{GetGatwayName(event.Project, stageName), "mesh"}
		hosts := []string{svc.Name + "." + c.canaryLevelGen.GetNamespace(event.Project, stageName, false) + "." + c.keptnDomain,
			svc.Name, svc.Name + "." + c.canaryLevelGen.GetNamespace(event.Project, stageName, false)}
		destCanary := mesh.HTTPRouteDestination{Host: hostCanary, Weight: 0}
		destPrimary := mesh.HTTPRouteDestination{Host: hostPrimary, Weight: 100}
		httpRouteDestinations := []mesh.HTTPRouteDestination{destCanary, destPrimary}

		vs, err := c.mesh.GenerateVirtualService(svc.Name, gws, hosts, httpRouteDestinations)
		if err != nil {
			return nil, nil, err
		}

		gw := chart.Template{Name: "templates/" + svc.Name + c.mesh.GetVirtualServiceSuffix(), Data: vs}
		newTemplates = append(newTemplates, &gw)
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
	if IsDeployment(&depl) {
		depl.Name = depl.Name + "-primary"
		depl.Spec.Selector.MatchLabels["app"] = depl.Spec.Selector.MatchLabels["app"] + "-primary"
		depl.Spec.Template.ObjectMeta.Labels["app"] = depl.Spec.Template.ObjectMeta.Labels["app"] + "-primary"
		var err error
		newTemplateContent, err = appendAsYaml(newTemplateContent, depl)
		if err != nil {
			return nil, err
		}
	}
	return newTemplateContent, nil
}

func (c *GeneratedChartHandler) SetCanaryWeight(e *keptnevents.ConfigurationChangeEventData, canaryWeight int32) error {

	// Read chart
	chart, err := GetChart(e.Project, e.Service, e.Stage, GetChartName(e.Service, true), c.configServiceURL)
	if err != nil {
		return fmt.Errorf("Error when reading chart %s from project %s: %s", GetChartName(e.Service, true), e.Project, err.Error())
	}

	// Set weights in all virtualservices
	for _, template := range chart.Templates {
		if strings.HasPrefix(template.Name, "templates/") &&
			strings.HasSuffix(template.Name, c.mesh.GetVirtualServiceSuffix()) {

			vs, err := c.mesh.UpdateWeights(template.Data, canaryWeight)
			if err != nil {
				return fmt.Errorf("Error when setting new weights in VirtualService %s from chart %s and project %s: %s",
					template.Name, GetChartName(e.Service, true), e.Project, err.Error())
			}
			template.Data = vs
		}
	}

	// Store chart
	chartData, err := PackageChart(chart)
	if err != nil {
		return fmt.Errorf("Error when packaging modified chart %s from project %s: %s", GetChartName(e.Service, true), e.Project, err.Error())
	}
	err = StoreChart(e.Project, e.Service, e.Stage, GetChartName(e.Service, true), chartData, c.configServiceURL)
	if err != nil {
		return fmt.Errorf("Error when storing modified chart %s from project %s: %s", GetChartName(e.Service, true), e.Project, err.Error())
	}
	return nil
}

func (c *GeneratedChartHandler) CopyCanaryToPrimary(e *keptnevents.ConfigurationChangeEventData) error {

	return nil
}
