package helm

import (
	"bytes"
	"encoding/json"
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

// GenerateManagedChart generates a duplicated chart which is managed by keptn and used for
// b/g and canary releases
func GenerateManagedChart(event *keptnevents.ServiceCreateEventData, stageName string,
	mesh mesh.Mesh, domain string) ([]byte, error) {

	ch, err := LoadChart(event.HelmChart)
	if err != nil {
		return nil, err
	}

	changeChartFile(ch)

	if err := changeTemplateContent(event, ch, mesh, stageName, domain); err != nil {
		return nil, err
	}

	return PackageChart(ch)
}

func changeChartFile(ch *chart.Chart) {

	ch.Metadata.Name = ch.Metadata.Name + "-generated"
	ch.Metadata.Description = ch.Metadata.Description + " (generated)"
}

func changeTemplateContent(event *keptnevents.ServiceCreateEventData,
	ch *chart.Chart, meshHandler mesh.Mesh, stageName string, domain string) error {

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

			newServiceTemplateContent, newServiceTemplates, err := handleService(doc, event, stageName, meshHandler, domain)
			if err != nil {
				return err
			}
			if len(newServiceTemplateContent) > 0 {
				newContent = append(newContent, newServiceTemplateContent...)
				newTemplates = append(newTemplates, newServiceTemplates...)
				continue
			}

			newDeploymentTemplateContent, err := handleDeployment(doc)
			if err != nil {
				return err
			}
			if len(newDeploymentTemplateContent) > 0 {
				newContent = append(newContent, newDeploymentTemplateContent...)
				continue
			}

			newContent, err = appendAsYaml(newContent, document)
			if err != nil {
				return err
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

func handleService(document []byte, event *keptnevents.ServiceCreateEventData, stageName string, meshHandler mesh.Mesh,
	domain string) ([]byte, []*chart.Template, error) {

	var svc corev1.Service
	if err := json.Unmarshal(document, &svc); err != nil {
		// Ignore unmarshaling error
		return nil, nil, nil
	}

	newTemplateContent := make([]byte, 0, 0)
	newTemplates := make([]*chart.Template, 0, 0)

	if isService(svc) {

		var err error
		serviceCanary := svc.DeepCopy()
		serviceCanary.Name = serviceCanary.Name + "-canary"
		newTemplateContent, err = appendAsYaml(newTemplateContent, serviceCanary)
		if err != nil {
			return nil, nil, err
		}

		// Generate destination rule for canary service
		hostCanary := serviceCanary.Name + "." + GetNamespace(event.Project, stageName, true) + ".svc.cluster.local"
		destinationRuleCanary, err := meshHandler.GenerateDestinationRule(serviceCanary.Name, hostCanary)
		if err != nil {
			return nil, nil, err
		}

		d1 := chart.Template{Name: "templates/" + serviceCanary.Name + meshHandler.GetDestinationRuleSuffix(), Data: destinationRuleCanary}
		newTemplates = append(newTemplates, &d1)

		servicePrimary := svc.DeepCopy()
		servicePrimary.Name = servicePrimary.Name + "-primary"
		servicePrimary.Spec.Selector["app"] = servicePrimary.Spec.Selector["app"] + "-primary"
		newTemplateContent, err = appendAsYaml(newTemplateContent, servicePrimary)
		if err != nil {
			return nil, nil, err
		}

		// Generate destination rule for primary service
		hostPrimary := servicePrimary.Name + "." + GetNamespace(event.Project, stageName, true) + ".svc.cluster.local"
		destinationRulePrimary, err := meshHandler.GenerateDestinationRule(servicePrimary.Name, hostPrimary)
		if err != nil {
			return nil, nil, err
		}

		d2 := chart.Template{Name: "templates/" + servicePrimary.Name + meshHandler.GetDestinationRuleSuffix(), Data: destinationRulePrimary}
		newTemplates = append(newTemplates, &d2)

		gws := []string{GetGatwayName(event.Project, stageName), "mesh"}
		hosts := []string{svc.Name + "." + GetNamespace(event.Project, stageName, false) + "." + domain,
			svc.Name, svc.Name + "." + GetNamespace(event.Project, stageName, false)}
		destCanary := mesh.HTTPRouteDestination{Host: hostCanary, Weight: 0}
		destPrimary := mesh.HTTPRouteDestination{Host: hostPrimary, Weight: 100}
		httpRouteDestinations := []mesh.HTTPRouteDestination{destCanary, destPrimary}

		vs, err := meshHandler.GenerateVirtualService(svc.Name, gws, hosts, httpRouteDestinations)
		if err != nil {
			return nil, nil, err
		}

		gw := chart.Template{Name: "templates/" + svc.Name + meshHandler.GetVirtualServiceSuffix(), Data: vs}
		newTemplates = append(newTemplates, &gw)
	}
	return newTemplateContent, newTemplates, nil
}

func handleDeployment(document []byte) ([]byte, error) {
	// Try to unmarshal Deployment
	var depl appsv1.Deployment
	if err := json.Unmarshal(document, &depl); err != nil {
		// Ignore unmarshaling error
		return nil, nil
	}

	newTemplateContent := make([]byte, 0, 0)
	if isDeployment(depl) {
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

func isService(svc corev1.Service) bool {
	return strings.ToLower(svc.Kind) == "service"
}

func isDeployment(dpl appsv1.Deployment) bool {
	return strings.ToLower(dpl.Kind) == "deployment"
}
