package validator

import (
	"fmt"
	"strings"

	"github.com/keptn/keptn/cli/pkg/logging"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"helm.sh/helm/v3/pkg/chart"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

var reservedFileNameSuffixes = [...]string{"-istio-destinationrule.yaml", "-istio-virtualservice.yaml"}

// ValidateHelmChart validates keptn's requirements regarding
// the values, deployment, and service file
func ValidateHelmChart(ch *chart.Chart) (bool, error) {

	if resValues := validateValues(ch); !resValues {
		return false, nil
	}

	services, err := keptnutils.GetRenderedServices(ch)
	if err != nil {
		return false, err
	}
	if resServices, err := validateServices(services); !resServices || err != nil {
		return false, err
	}

	deployments, err := keptnutils.GetRenderedDeployments(ch)
	if err != nil {
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
	// check image property
	if _, containsImage := ch.Values["image"]; !containsImage {
		logging.PrintLog("Provided Helm chart does not contain \"image\" in values.yaml", logging.QuietLevel)
		return false
	}
	// check replicas property
	if _, containsReplicas := ch.Values["replicaCount"]; !containsReplicas {
		logging.PrintLog("Provided Helm chart does not contain \"replicaCount\" in values.yaml", logging.QuietLevel)
		return false
	}
	return true
}

func validateServices(services []*corev1.Service) (bool, error) {
	for _, svc := range services {
		if !validateService(svc) {
			return false, nil
		}
	}
	if len(services) != 1 {
		logging.PrintLog("Helm chart must contain exact one service", logging.QuietLevel)
		return false, nil
	}
	return true, nil
}

func validateDeployments(deployments []*appsv1.Deployment) (bool, error) {
	for _, depl := range deployments {
		if !validateDeployment(depl) {
			return false, nil
		}
	}
	if len(deployments) != 1 {
		logging.PrintLog("Helm chart must contain exact one deployment", logging.QuietLevel)
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
