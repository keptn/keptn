package handlers

import (
	"errors"
	"fmt"

	"github.com/go-openapi/swag"

	"github.com/go-openapi/runtime/middleware"
	keptnutils "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/configure"
	k8sutils "github.com/keptn/kubernetes-utils/pkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Ingress is an enum type for the ingress
type Ingress int

const (
	istio Ingress = iota
	nginx
)

func (i Ingress) String() string {
	return [...]string{"istio", "nginx"}[i]
}

const useInClusterConfig = false

func getIngressType() (Ingress, error) {

	clientset, err := k8sutils.GetKubeAPI(useInClusterConfig)
	if err != nil {
		return istio, err
	}

	nsList, err := clientset.Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return istio, err
	}

	for _, ns := range nsList.Items {
		if ns.Name == "istio-system" {
			return istio, nil
		} else if ns.Name == "ingress-nginx" {
			return nginx, nil
		}
	}
	return istio, errors.New("Cannot obtain type of ingress.")
}

func PostConfigureBridgeHandlerFunc(params configure.PostConfigureBridgeExposeParams, principal *models.Principal) middleware.Responder {

	l := keptnutils.NewLogger("", "", "api")
	l.Info("API received a configure Bridge request")

	ingress, err := getIngressType()
	if err != nil {
		errMsg := fmt.Sprintf("failed to retrieve ingress type: %v", err)
		l.Error(errMsg)
		return configure.NewPostConfigureBridgeExposeDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(errMsg)})
	}
	l.Info("Used ingress for configure Bridge: " + ingress.String())
	domain, err := k8sutils.GetKeptnDomain(useInClusterConfig)
	if err != nil {
		errMsg := fmt.Sprintf("failed to retrieve domain: %v", err)
		l.Error(errMsg)
		return configure.NewPostConfigureBridgeExposeDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(errMsg)})
	}
	l.Info("Used domain for configure Bridge: " + domain)

	if params.Expose {
		l.Info("Starting to expose bridge")

	} else {
		l.Info("Starting to remove exposure of bridge")

	}

	return configure.NewPostConfigureBridgeExposeOK()
}
