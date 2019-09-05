package validator

import (
	"k8s.io/helm/pkg/proto/hapi/chart"

	"github.com/ghodss/yaml"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// ValidateHelmChart validates keptn's requirements regarding
// the values, deployment, and service file
func ValidateHelmChart(helmChart []byte) (bool, error) {

	ch, err := keptnutils.LoadChart(helmChart)
	if err != nil {
		return false, err
	}
	return validateValues(ch) && validateServices(ch) && validateDeployments(ch), nil
}

func validateValues(ch *chart.Chart) bool {

	values := make(map[string]interface{})
	yaml.Unmarshal([]byte(ch.Values.Raw), &values)
	_, containsImage := values["image"]
	_, containReplicaCount := values["replicas"]
	return containsImage && containReplicaCount
}

func validateServices(ch *chart.Chart) bool {
	services := keptnutils.GetServices(ch)
	for _, svc := range services {
		if !validateService(svc) {
			return false
		}
	}
	return len(services) > 0
}

func validateDeployments(ch *chart.Chart) bool {
	deployments := keptnutils.GetDeployments(ch)
	for _, depl := range deployments {
		if !validateDeployment(depl) {
			return false
		}
	}
	return len(deployments) > 0
}

func validateService(svc *corev1.Service) bool {

	val, ok := svc.Spec.Selector["app"]
	return keptnutils.IsService(svc) && ok && val != ""
}

func validateDeployment(depl *appsv1.Deployment) bool {

	mLabel, ok1 := depl.Spec.Selector.MatchLabels["app"]
	podLabel, ok2 := depl.Spec.Template.ObjectMeta.Labels["app"]
	return keptnutils.IsDeployment(depl) && ok1 && ok2 && mLabel != "" && podLabel != ""
}
