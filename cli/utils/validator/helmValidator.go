package validator

import (
	"fmt"
	"strings"

	"k8s.io/helm/pkg/proto/hapi/chart"

	"github.com/ghodss/yaml"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

var reservedFileNameSuffixes = [...]string{"-istio-destinationrule.yaml", "-istio-virtualservice.yaml"}

// ValidateHelmChart validates keptn's requirements regarding
// the values, deployment, and service file
func ValidateHelmChart(helmChart []byte) (bool, error) {

	ch, err := keptnutils.LoadChart(helmChart)
	if err != nil {
		return false, err
	}
	resValues := validateValues(ch)
	resServiecs, err := validateServices(ch)
	if err != nil {
		return false, err
	}
	resDeployment, err := validateDeployments(ch)
	if err != nil {
		return false, err
	}
	return resValues && resServiecs && resDeployment, nil
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
	yaml.Unmarshal([]byte(ch.Values.Raw), &values)
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

func validateServices(ch *chart.Chart) (bool, error) {
	services, err := keptnutils.GetRenderedServices(ch)
	if err != nil {
		logging.PrintLog("Error rendering Helm chart", logging.QuietLevel)
		return false, err
	}
	for _, svc := range services {
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

func validateDeployments(ch *chart.Chart) (bool, error) {
	deployments, err := keptnutils.GetRenderedDeployments(ch)
	if err != nil {
		logging.PrintLog("Error rendering Helm chart", logging.QuietLevel)
		return false, err
	}
	for _, depl := range deployments {
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

func validateService(svc *corev1.Service) bool {
	if !keptnutils.IsService(svc) {
		logging.PrintLog(fmt.Sprintf("Service %s does not have kind \"service\"", svc.Name), logging.QuietLevel)
		return false
	}
	if svc.Spec.Selector == nil {
		logging.PrintLog(fmt.Sprintf("Service %s does not contain \"selector\"", svc.Name), logging.QuietLevel)
		return false
	}
	val, ok := svc.Spec.Selector["app"]
	if !ok || val == "" {
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
	mLabel, okmLabel := depl.Spec.Selector.MatchLabels["app"]
	if !okmLabel || mLabel == "" {
		logging.PrintLog(fmt.Sprintf("Deployment %s does not contain \"spec.selector.matchLabels.app\"", depl.Name), logging.QuietLevel)
		return false
	}
	podLabel, okPodLabel := depl.Spec.Template.ObjectMeta.Labels["app"]
	if !okPodLabel || podLabel == "" {
		logging.PrintLog(fmt.Sprintf("Deployment %s does not contain \"spec.template.metadata.labels.app\"", depl.Name), logging.QuietLevel)
		return false
	}
	return true
}
