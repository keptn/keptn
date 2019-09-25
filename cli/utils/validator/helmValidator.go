package validator

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"

	"github.com/ghodss/yaml"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/renderutil"
	"k8s.io/helm/pkg/timeconv"
)

var reservedFileNameSuffixes = [...]string{"-istio-destinationrule.yaml", "-istio-virtualservice.yaml"}

// ValidateHelmChart validates keptn's requirements regarding
// the values, deployment, and service file
func ValidateHelmChart(ch *chart.Chart) (bool, error) {

	if resValues := validateValues(ch); !resValues {
		return false, nil
	}

	services, err := getRenderedServices(ch)
	if err != nil {
		return false, err
	}
	if resSvcBeforeRendering, err := validateServiceBeforeRendering(ch, services); !resSvcBeforeRendering || err != nil {
		return false, err
	}
	if resServices, err := validateServices(services); !resServices || err != nil {
		return false, err
	}

	deployments, err := getRenderedDeployments(ch)
	if err != nil {
		return false, err
	}
	if resDeplBeforeRendering, err := validateDeploymentBeforeRendering(ch, deployments); !resDeplBeforeRendering || err != nil {
		return false, err
	}
	if resDeployment, err := validateDeployments(deployments); !resDeployment || err != nil {
		return false, err
	}
	if !validateTemplateFileNames(ch) {
		return false, nil
	}
	return true, nil
}

func validateTemplateFileNames(ch *chart.Chart) bool {
	for _, t := range ch.Templates {
		for _, s := range reservedFileNameSuffixes {
			if strings.HasSuffix(t.Name, s) {
				logging.PrintLog(fmt.Sprintf("File %s has a reserved file name suffix", t.Name), logging.QuietLevel)
				return false
			}
		}
	}
	return true
}

func validateValues(ch *chart.Chart) bool {
	values := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(ch.Values.Raw), &values); err != nil {
		return false
	}
	_, containsImage := values["image"]
	if !containsImage {
		logging.PrintLog("Provided Helm chart does not contain \"image\" in values.yaml", logging.QuietLevel)
		return false
	}
	_, containsReplicaCount := values["replicas"]
	if !containsReplicaCount {
		logging.PrintLog("Provided Helm chart does not contain \"replicas\" in values.yaml", logging.QuietLevel)
		return false
	}
	return true
}

func validateServices(services map[*corev1.Service]string) (bool, error) {
	for svc := range services {
		if !validateService(svc) {
			return false, nil
		}
	}
	if len(services) == 0 {
		logging.PrintLog("Helm chart must contain at lease one service", logging.QuietLevel)
		return false, nil
	}
	return true, nil
}

func validateDeployments(deployments map[*appsv1.Deployment]string) (bool, error) {
	for depl := range deployments {
		if !validateDeployment(depl) {
			return false, nil
		}
	}
	if len(deployments) != 1 {
		logging.PrintLog("Helm chart must contain a single deployment", logging.QuietLevel)
		return false, nil
	}
	return true, nil
}

func validateServiceBeforeRendering(ch *chart.Chart, services map[*corev1.Service]string) (bool, error) {

	for k, v := range services {
		svcBeforeRendering := getTemplateContent(ch, v)

		reg, err := regexp.Compile(`name:\s+` + k.ObjectMeta.Name)
		if err != nil {
			log.Fatal(err)
		}
		if !reg.MatchString(svcBeforeRendering) {
			logging.PrintLog(fmt.Sprintf("For metadata.name no templates are allowed in file %s", v), logging.QuietLevel)
			return false, nil
		}

		if _, ok := k.Spec.Selector["app"]; ok {
			reg, err := regexp.Compile(`app:\s+\{\{.*`)
			if err != nil {
				return false, err
			}
			if reg.MatchString(svcBeforeRendering) {
				logging.PrintLog(fmt.Sprintf(
					"For \"app\" no template is allowed in file %s", v), logging.QuietLevel)
				return false, nil
			}
		}

		if _, ok := k.Spec.Selector["app.kubernetes.io/name"]; ok {
			reg, err := regexp.Compile(`app.kubernetes.io/name:\s+\{\{.*`)
			if err != nil {
				return false, err
			}
			if reg.MatchString(svcBeforeRendering) {
				logging.PrintLog(fmt.Sprintf(
					"For \"app.kubernetes.io/name\" no template is allowed in file %s", v), logging.QuietLevel)
				return false, nil
			}
		}
	}
	return true, nil
}

func validateDeploymentBeforeRendering(ch *chart.Chart, deployments map[*appsv1.Deployment]string) (bool, error) {

	for k, v := range deployments {
		dplBeforeRendering := getTemplateContent(ch, v)

		reg, err := regexp.Compile(`name:\s+` + k.Name)
		if err != nil {
			return false, err
		}
		if !reg.MatchString(dplBeforeRendering) {
			logging.PrintLog(fmt.Sprintf(
				"For \"name\" no template is allowed in file %s", v), logging.QuietLevel)
			return false, nil
		}

		_, okmLabelApp := k.Spec.Selector.MatchLabels["app"]
		_, okPodLabelApp := k.Spec.Template.ObjectMeta.Labels["app"]
		if okmLabelApp || okPodLabelApp {
			reg, err := regexp.Compile(`app:\s+\{\{.*`)
			if err != nil {
				return false, err
			}
			if reg.MatchString(dplBeforeRendering) {
				logging.PrintLog(fmt.Sprintf(
					"For \"app\" no template is allowed in file %s", v), logging.QuietLevel)
				return false, nil
			}
		}
		_, okmLabelAppk8sName := k.Spec.Selector.MatchLabels["app.kubernetes.io/name"]
		_, okPodLabelAppk8sName := k.Spec.Template.ObjectMeta.Labels["app.kubernetes.io/name"]
		if okmLabelAppk8sName || okPodLabelAppk8sName {
			reg, err := regexp.Compile(`app.kubernetes.io/name:\s+\{\{.*`)
			if err != nil {
				return false, err
			}
			if reg.MatchString(dplBeforeRendering) {
				logging.PrintLog(fmt.Sprintf(
					"For \"app.kubernetes.io/name\" no template is allowed in file %s", v), logging.QuietLevel)
				return false, nil
			}
		}
	}
	return true, nil
}

func getTemplateContent(ch *chart.Chart, fileName string) string {

	for _, t := range ch.Templates {
		if strings.HasSuffix(fileName, t.Name) {
			return string(t.Data)
		}
	}
	return ""
}

func getRenderedTemplates(ch *chart.Chart) (map[string]string, error) {
	renderOpts := renderutil.Options{
		ReleaseOptions: chartutil.ReleaseOptions{
			Name:      ch.Metadata.Name,
			IsInstall: false,
			IsUpgrade: false,
			Time:      timeconv.Now(),
		},
	}

	return renderutil.Render(ch, ch.Values, renderOpts)
}

func getRenderedServices(ch *chart.Chart) (map[*corev1.Service]string, error) {

	renderedTemplates, err := getRenderedTemplates(ch)
	if err != nil {
		return nil, err
	}

	services := make(map[*corev1.Service]string)

	for k, v := range renderedTemplates {
		dec := kyaml.NewYAMLToJSONDecoder(strings.NewReader(v))
		for {
			var svc corev1.Service
			err := dec.Decode(&svc)
			if err == io.EOF {
				break
			}
			if err != nil {
				continue
			}

			if keptnutils.IsService(&svc) {
				services[&svc] = k
			}
		}
	}

	return services, nil

}

func getRenderedDeployments(ch *chart.Chart) (map[*appsv1.Deployment]string, error) {

	renderedTemplates, err := getRenderedTemplates(ch)
	if err != nil {
		return nil, err
	}

	deployments := make(map[*appsv1.Deployment]string)

	for k, v := range renderedTemplates {
		dec := kyaml.NewYAMLToJSONDecoder(strings.NewReader(v))
		for {
			var depl appsv1.Deployment
			err := dec.Decode(&depl)
			if err == io.EOF {
				break
			}
			if err != nil {
				continue
			}

			if keptnutils.IsDeployment(&depl) {
				deployments[&depl] = k
			}
		}
	}

	return deployments, nil
}

func validateService(svc *corev1.Service) bool {
	if !keptnutils.IsService(svc) {
		logging.PrintLog(fmt.Sprintf("Service %s does not have kind \"service\"", svc.Name), logging.QuietLevel)
		return false
	}
	if svc.Spec.Selector == nil {
		logging.PrintLog(fmt.Sprintf("Service %s does not contain \"selector\"", svc.Name), logging.QuietLevel)
		return false
	}
	valApp, okApp := svc.Spec.Selector["app"]
	valAppk8sName, okAppk8sName := svc.Spec.Selector["app.kubernetes.io/name"]
	if (!okApp || valApp == "") && (!okAppk8sName || valAppk8sName == "") {
		logging.PrintLog(fmt.Sprintf("Service %s does not have \"spec.selector.app\"", svc.Name), logging.QuietLevel)
		return false
	}
	return true
}

func validateDeployment(depl *appsv1.Deployment) bool {
	if !keptnutils.IsDeployment(depl) {
		logging.PrintLog(fmt.Sprintf("Deployment %s does not have kind \"deployment\"", depl.Name), logging.QuietLevel)
		return false
	}
	if depl.Spec.Selector == nil {
		logging.PrintLog(fmt.Sprintf("Deployment %s does not contain \"selector\"", depl.Name), logging.QuietLevel)
		return false
	}
	if depl.Spec.Selector.MatchLabels == nil {
		logging.PrintLog(fmt.Sprintf("Deployment %s does not contain \"selector.matchLabels\"", depl.Name), logging.QuietLevel)
		return false
	}
	if depl.Spec.Template.ObjectMeta.Labels == nil {
		logging.PrintLog(fmt.Sprintf("Deployment %s does not contain \"spec.template.metadata.labels\"", depl.Name), logging.QuietLevel)
		return false
	}
	mLabelApp, okmLabelApp := depl.Spec.Selector.MatchLabels["app"]
	mLabelAppk8sName, okmLabelAppk8sName := depl.Spec.Selector.MatchLabels["app.kubernetes.io/name"]
	if (!okmLabelApp || mLabelApp == "") && (!okmLabelAppk8sName || mLabelAppk8sName == "") {
		logging.PrintLog(fmt.Sprintf("Deployment %s does not contain \"spec.selector.matchLabels.app\"", depl.Name), logging.QuietLevel)
		return false
	}
	podLabelApp, okPodLabelApp := depl.Spec.Template.ObjectMeta.Labels["app"]
	podLabelAppk8sName, okPodLabelAppk8sName := depl.Spec.Template.ObjectMeta.Labels["app.kubernetes.io/name"]
	if (!okPodLabelApp || podLabelApp == "") && (!okPodLabelAppk8sName || podLabelAppk8sName == "") {
		logging.PrintLog(fmt.Sprintf("Deployment %s does not contain \"spec.template.metadata.labels.app\"", depl.Name), logging.QuietLevel)
		return false
	}
	return true
}
