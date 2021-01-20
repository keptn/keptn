package controller

import (
	"encoding/base64"
	"errors"
	"fmt"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
)

func CheckAndSetServiceName(event *keptnv2.ServiceCreateFinishedEventData) error {

	errorMsg := "Service name contains upper case letter(s) or special character(s).\n " +
		"Keptn relies on the following conventions: " +
		"start with a lower case letter, then lower case letters, numbers, and hyphens are allowed."

	helmChartData, err := base64.StdEncoding.DecodeString(event.Helm.Chart)
	if err != nil {
		return fmt.Errorf("Error when decoding the Helm Chart: %v", err)
	}
	ch, err := keptnutils.LoadChart(helmChartData)
	if err != nil {
		return fmt.Errorf("Error when loading Helm Chart: %v", err)
	}
	services, err := keptnutils.GetRenderedServices(ch)
	if err != nil {
		return fmt.Errorf("Error when rendering services: %v", err)
	}
	if len(services) != 1 {
		return fmt.Errorf("Helm Chart has to contain exactly one Kubernetes service, but it contains %d services", len(services))
	}
	k8sServiceName := services[0].Name
	if !keptncommon.ValidateKeptnEntityName(k8sServiceName) {
		return errors.New(errorMsg)
	}
	if event.Service == "" {
		// Set service name in event
		event.Service = k8sServiceName
	}
	if k8sServiceName != event.Service {
		return fmt.Errorf("Provided Keptn service name \"%s\" "+
			"does not match Kubernetes service name \"%s\"", event.Service, k8sServiceName)
	}
	return nil
}
