package helm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/controller/jsonutils"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
)

// GenerateManagedChart generates a duplicated chart which is managed by keptn and used for
// b/g and canary releases
func GenerateManagedChart(event *keptnevents.ServiceCreateEventData, stageName string,
	mesh mesh.Mesh, domain string) ([]byte, error) {

	workingPath, err := ioutil.TempDir("", "helm")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(workingPath)

	if err := keptnutils.Untar(workingPath, bytes.NewReader(event.HelmChart)); err != nil {
		return nil, err
	}

	helmCharts, err := ioutil.ReadDir(workingPath)
	if err != nil {
		return nil, err
	}
	if len(helmCharts) != 1 {
		return nil, errors.New("Multiple helm charts are found wihin tar")
	}

	chartFolder := filepath.Join(workingPath, helmCharts[0].Name())

	if err := changeChartFile(chartFolder); err != nil {
		return nil, err
	}

	if err := changeTemplateContent(event, filepath.Join(chartFolder, "templates"),
		mesh, stageName, domain); err != nil {
		return nil, err
	}

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	if err := keptnutils.Tar(chartFolder, writer); err != nil {
		return nil, err
	}
	writer.Flush()

	return b.Bytes(), os.RemoveAll(workingPath)
}

func changeChartFile(chartPath string) error {

	chartFiles, err := keptnutils.GetFiles(chartPath, "Chart.yaml", "Chart.yml")
	if len(chartFiles) == 0 {
		return errors.New("No Chart.yaml file can be found for chart " + chartPath)
	} else if len(chartFiles) > 1 {
		return errors.New("Multiple Chart.yaml files found for chart " + chartPath)
	}

	dat, err := ioutil.ReadFile(chartFiles[0])
	if err != nil {
		return err
	}

	dec := kyaml.NewYAMLToJSONDecoder(bytes.NewReader(dat))
	var chart Chart
	if err := dec.Decode(&chart); err != nil {
		return err
	}
	chart.Name = chart.Name + "-generated"
	chart.Description = chart.Description + " (generated)"

	jsonBytes, err := json.Marshal(chart)
	if err != nil {
		return err
	}
	yamlBytes, err := jsonutils.ToYAML(jsonBytes)
	return ioutil.WriteFile(chartFiles[0], yamlBytes, 0644)
}

func getNamespace(projectName string, stageName string) string {
	return projectName + "-" + stageName
}

func changeTemplateContent(event *keptnevents.ServiceCreateEventData,
	templatesPath string, meshHandler mesh.Mesh, stageName string, domain string) error {

	templateFiles, err := keptnutils.GetFiles(templatesPath, ".yml", ".yaml")
	if err != nil {
		return err
	}

	for _, file := range templateFiles {
		dat, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		dec := kyaml.NewYAMLToJSONDecoder(bytes.NewReader(dat))
		elements := make([]interface{}, 0, 1)
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

			newServiceElements, err := handleService(doc, event, stageName, meshHandler, templatesPath, domain)
			if err != nil {
				return err
			}
			newDeploymentElement, err := handleDeployment(doc)
			if err != nil {
				return err
			}

			if len(newServiceElements) > 0 {
				elements = append(elements, newServiceElements...)
			} else if newDeploymentElement != nil {
				elements = append(elements, newDeploymentElement)
			} else {
				elements = append(elements, document)
			}
		}

		// Create new yaml file with changed services and deployments
		if err := writeNewContent(file, elements); err != nil {
			return err
		}
	}
	return nil
}

func writeNewContent(file string, elements []interface{}) error {
	newFileContent := ""
	for _, element := range elements {
		jsonData, err := json.Marshal(element)
		if err != nil {
			return err
		}
		yamlData, err := jsonutils.ToYAML(jsonData)
		if err != nil {
			return err
		}
		newFileContent = newFileContent + "---\n" + string(yamlData)
	}
	return ioutil.WriteFile(file, []byte(newFileContent), 0644)
}

func handleService(document []byte, event *keptnevents.ServiceCreateEventData, stageName string, meshHandler mesh.Mesh,
	templatesPath string, domain string) ([]interface{}, error) {
	var svc corev1.Service
	if err := json.Unmarshal(document, &svc); err != nil {
		// Ignore unmarshaling error
		return nil, nil
	}

	elements := make([]interface{}, 0, 0)
	if isService(svc) {

		serviceCanary := svc.DeepCopy()
		serviceCanary.Name = serviceCanary.Name + "-canary"
		elements = append(elements, serviceCanary)

		// Generate destination rule for canary service
		hostCanary := serviceCanary.Name + "." + getNamespace(event.Project, stageName) + ".svc.cluster.local"
		destinationRuleCanary, err := meshHandler.GenerateDestinationRule(serviceCanary.Name, hostCanary)
		if err != nil {
			return nil, err
		}
		if err := ioutil.WriteFile(filepath.Join(templatesPath, serviceCanary.Name+"-istio-destinationrule.yaml"),
			destinationRuleCanary, 0644); err != nil {
			return nil, err
		}

		servicePrimary := svc.DeepCopy()
		servicePrimary.Name = servicePrimary.Name + "-primary"
		servicePrimary.Spec.Selector["app"] = servicePrimary.Spec.Selector["app"] + "-primary"
		elements = append(elements, servicePrimary)

		// Generate destination rule for primary service
		hostPrimary := servicePrimary.Name + "." + getNamespace(event.Project, stageName) + ".svc.cluster.local"
		destinationRulePrimary, err := meshHandler.GenerateDestinationRule(servicePrimary.Name, hostPrimary)
		if err != nil {
			return nil, err
		}
		if err := ioutil.WriteFile(filepath.Join(templatesPath, servicePrimary.Name+"-istio-destinationrule.yaml"),
			destinationRulePrimary, 0644); err != nil {
			return nil, err
		}

		gws := []string{getGatwayName(event.Project, stageName)}
		hosts := []string{svc.Name + "." + getNamespace(event.Project, stageName) + "." + domain}
		destCanary := mesh.HTTPRouteDestination{Host: hostCanary, Weight: 80}
		destPrimary := mesh.HTTPRouteDestination{Host: hostPrimary, Weight: 20}
		httpRouteDestinations := []mesh.HTTPRouteDestination{destCanary, destPrimary}

		vs, err := meshHandler.GenerateVirtualService(svc.Name, gws, hosts, httpRouteDestinations)
		if err != nil {
			return nil, err
		}
		if err := ioutil.WriteFile(filepath.Join(templatesPath, svc.Name+"-istio-virtualservice.yaml"), vs, 0644); err != nil {
			return nil, err
		}
	}
	return elements, nil
}

func handleDeployment(document []byte) (interface{}, error) {
	// Try to unmarshal Deployment
	var depl appsv1.Deployment
	if err := json.Unmarshal(document, &depl); err != nil {
		// Ignore unmarshaling error
		return nil, nil
	}

	if isDeployment(depl) {
		depl.Name = depl.Name + "-primary"
		depl.Spec.Selector.MatchLabels["app"] = depl.Spec.Selector.MatchLabels["app"] + "-primary"
		depl.Spec.Template.ObjectMeta.Labels["app"] = depl.Spec.Template.ObjectMeta.Labels["app"] + "-primary"
		return depl, nil
	}
	return nil, nil
}

func isService(svc corev1.Service) bool {
	return strings.ToLower(svc.Kind) == "service"
}

func isDeployment(dpl appsv1.Deployment) bool {
	return strings.ToLower(dpl.Kind) == "deployment"
}
