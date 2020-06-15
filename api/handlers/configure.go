package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/util/intstr"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/go-openapi/swag"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"

	"github.com/go-openapi/runtime/middleware"
	keptnutils "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/configure"
	k8sutils "github.com/keptn/kubernetes-utils/pkg"
	networking "k8s.io/api/networking/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Ingress is an enum type for the ingress
type Ingress int

const (
	istio Ingress = iota
	nginx
	ocroute
)

const keptnGateway = "public-gateway.istio-system"

func (i Ingress) String() string {
	return [...]string{"istio", "nginx", "openshift"}[i]
}

const useInClusterConfig = true

func getIngressType() (Ingress, error) {

	var config *rest.Config
	var err error
	config, err = rest.InClusterConfig()

	if err != nil {
		return istio, err
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return istio, err
	}

	nsList, err := k8sClient.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return istio, err
	}

	for _, ns := range nsList.Items {
		if ns.Name == "ingress-nginx" {
			ingresses, err := k8sClient.ExtensionsV1beta1().Ingresses("ingress-nginx").List(metav1.ListOptions{})
			if err != nil {
				return istio, err
			}
			for _, ingress := range ingresses.Items {
				if ingress.Name == "keptn-ingress" {
					return nginx, nil
				}
			}
		} else if ns.Name == "istio-system" {
			return istio, nil
		} else if ns.Name == "openshift" {
			// Note: The istio-system namespace was not found, hence OC routes are used
			return ocroute, nil
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
	l.Info("Used ingress for configuring Bridge: " + ingress.String())

	var exposeErr error
	var bridgeHost string
	if *params.ConfigureBridge.Expose {
		if params.ConfigureBridge.User == "" {
			errMsg := fmt.Sprintf("no user provided")
			l.Error(errMsg)
			return configure.NewPostConfigureBridgeExposeDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(errMsg)})
		}
		if params.ConfigureBridge.Password == "" {
			errMsg := fmt.Sprintf("no password provided")
			l.Error(errMsg)
			return configure.NewPostConfigureBridgeExposeDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(errMsg)})
		}

		err := createBridgeCredentials(params.ConfigureBridge.User, params.ConfigureBridge.Password, l)
		if err != nil {
			errMsg := fmt.Sprintf("failed to create bridge credentials: %v", err)
			l.Error(errMsg)
			return configure.NewPostConfigureBridgeExposeDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(errMsg)})
		}

		err = restartBridgePod(l)
		if err != nil {
			errMsg := fmt.Sprintf("failed to restart bridge pod: %v", err)
			l.Error(errMsg)
			return configure.NewPostConfigureBridgeExposeDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(errMsg)})
		}

		l.Info("Starting to retrieve Keptn domain")
		domain, err := k8sutils.GetKeptnDomain(useInClusterConfig)
		if err != nil {
			errMsg := fmt.Sprintf("failed to retrieve domain: %v", err)
			l.Error(errMsg)
			return configure.NewPostConfigureBridgeExposeDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(errMsg)})
		}
		bridgeHost = getHostForBridge(domain)
		l.Info("Used domain for configure Bridge: " + domain)
		l.Info("Used host for Bridge: " + getHostForBridge(domain))
		switch ingress {
		case istio:
			exposeErr = exposeBridgeIstio(domain, l)
		case nginx:
			exposeErr = exposeBridgeIngress(domain, l)
		case ocroute:
			exposeErr = sendOCRouteRequest(domain, true, l)
		}

	} else {
		l.Info("Starting to dispose bridge")
		switch ingress {
		case istio:
			exposeErr = disposeBridgeIstio(l)
		case nginx:
			exposeErr = disposeBridgeIngress(l)
		case ocroute:
			exposeErr = sendOCRouteRequest("", false, l)
		}
		bridgeHost = ""

		err = deleteBridgeCredentials(l)
		if err != nil {
			l.Error(err.Error())
			return configure.NewPostConfigureBridgeExposeDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}

		err = restartBridgePod(l)
		if err != nil {
			errMsg := fmt.Sprintf("failed to restart bridge pod: %v", err)
			l.Error(errMsg)
			return configure.NewPostConfigureBridgeExposeDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(errMsg)})
		}
	}
	if exposeErr != nil {
		l.Error(exposeErr.Error())
		return configure.NewPostConfigureBridgeExposeDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(exposeErr.Error())})
	}

	return configure.NewPostConfigureBridgeExposeOK().WithPayload(bridgeHost)
}

func deleteBridgeCredentials(l *keptnutils.Logger) error {
	l.Info("Deleting credentials for bridge")

	restConfig, _ := getRestConfig()

	k8s, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	return k8s.CoreV1().Secrets("keptn").Delete("bridge-credentials", &metav1.DeleteOptions{})
}

func restartBridgePod(l *keptnutils.Logger) error {
	l.Info("Restarting bridge for credentials to take effect")

	restConfig, _ := getRestConfig()

	k8s, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	return k8s.CoreV1().Pods("keptn").DeleteCollection(&metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: "run=bridge",
	})
}

func createBridgeCredentials(user string, password string, l *keptnutils.Logger) error {
	l.Info("Creating or updating credentials for bridge")

	restConfig, _ := getRestConfig()

	k8s, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	l.Info("Checking for existing secret")
	bridgeCredentials, err := k8s.CoreV1().Secrets("keptn").Get("bridge-credentials", metav1.GetOptions{})
	if bridgeCredentials != nil {
		// update existing secret
		l.Info("Existing secret found. Updating with new values for user and password")
		newSecret := getBridgeCredentials(user, password)
		bridgeCredentials.Data = newSecret.Data
		_, err = k8s.CoreV1().Secrets("keptn").Update(newSecret)
		if err != nil {
			l.Error("could not update secret: " + err.Error())
			return err
		}
	} else {
		l.Info("Creating a new secret")
		newSecret := getBridgeCredentials(user, password)
		_, err = k8s.CoreV1().Secrets("keptn").Create(newSecret)
		if err != nil {
			l.Error("could not create new secret: " + err.Error())
			return err
		}
	}
	return nil
}

func getBridgeCredentials(user string, password string) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bridge-credentials",
			Namespace: "keptn",
		},
		Data: map[string][]byte{
			"BASIC_AUTH_USERNAME": []byte(user),
			"BASIC_AUTH_PASSWORD": []byte(password),
		},
		Type: "Opaque",
	}
}

func getDestinationRule() *v1alpha3.DestinationRule {

	return &v1alpha3.DestinationRule{
		TypeMeta:   metav1.TypeMeta{APIVersion: "networking.istio.io/v1alpha3", Kind: "DestinationRule"},
		ObjectMeta: metav1.ObjectMeta{Namespace: "keptn", Name: "bridge"},
		Spec:       networkingv1alpha3.DestinationRule{Host: "bridge.keptn.svc.cluster.local"},
	}
}

func getVirtualService(keptnDomain string) *v1alpha3.VirtualService {

	return &v1alpha3.VirtualService{
		TypeMeta: metav1.TypeMeta{APIVersion: "networking.istio.io/v1alpha3", Kind: "VirtualService"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bridge",
			Namespace: "keptn",
		},
		Spec: networkingv1alpha3.VirtualService{
			Hosts:    []string{"bridge.keptn." + keptnDomain},
			Gateways: []string{keptnGateway},
			Http: []*networkingv1alpha3.HTTPRoute{{
				Route: []*networkingv1alpha3.HTTPRouteDestination{
					{Destination: &networkingv1alpha3.Destination{Host: "bridge.keptn.svc.cluster.local"}}},
			}},
		},
	}
}

func getRestConfig() (config *rest.Config, err error) {
	if useInClusterConfig {
		config, err = rest.InClusterConfig()
	} else {
		kubeconfig := filepath.Join(
			k8sutils.UserHomeDir(), ".kube", "config",
		)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return
}

func exposeBridgeIstio(keptnDomain string, l *keptnutils.Logger) error {
	l.Info("Expose Bridge using a VirtualService and DestinationRule")

	restConfig, _ := getRestConfig()
	ic, err := versionedclient.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	_, err = ic.NetworkingV1alpha3().VirtualServices("keptn").Create(getVirtualService(keptnDomain))
	if k8serrors.IsAlreadyExists(err) {
		l.Info("VirtualService already exists")
	} else if err != nil {
		return err
	} else {
		l.Info("VirtualService for Bridge created")
	}

	_, err = ic.NetworkingV1alpha3().DestinationRules("keptn").Create(getDestinationRule())
	if k8serrors.IsAlreadyExists(err) {
		l.Info("DestinationRule already exists")
	} else if err != nil {
		return err
	} else {
		l.Info("DestinationRule for Bridge created")
	}
	return nil
}

func exposeBridgeIngress(keptnDomain string, l *keptnutils.Logger) error {
	l.Info("Expose Bridge using the keptn-ingress")

	clientset, err := k8sutils.GetClientset(useInClusterConfig)
	if err != nil {
		return err
	}
	ing, err := clientset.NetworkingV1beta1().Ingresses("keptn").Get("keptn-ingress", metav1.GetOptions{})
	if err != nil {
		return err
	}
	l.Info("keptn-ingress retrieved")
	addBridgeToIngress(keptnDomain, ing)

	_, err = clientset.NetworkingV1beta1().Ingresses("keptn").Update(ing)
	if err != nil {
		return err
	}
	l.Info("Rule for Bridge successfully added to Keptn-ingress")
	return nil
}

func sendOCRouteRequest(keptnDomain string, expose bool, l *keptnutils.Logger) error {
	l.Info("Expose Bridge using the openshift-route-service")

	url, err := keptnutils.GetServiceEndpoint("OPENSHIFT_ROUTE_SERVICE_URI")
	if err != nil {
		return err
	}
	url.Path = "/configure/bridge/expose"
	if err != nil {
		return err
	}
	var jsonStr = []byte(`{"expose": ` + strconv.FormatBool(expose) + `, "keptnDomain": "` + keptnDomain + `"}`)
	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var actionPast string
	var actionPres string
	if expose {
		l.Info("Used url of openshift-route-service: " + url.String())
		actionPast = "added"
		actionPres = "add"
	} else {
		actionPast = "deleted"
		actionPast = "delete"
	}

	if resp.StatusCode == 200 {
		l.Info("Route for Bridge successfully " + actionPast)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			l.Error(fmt.Sprintf("Failed to "+actionPres+" route with status %d and error %s", resp.StatusCode, string(body)))
		} else {
			l.Error(fmt.Sprintf("Failed to "+actionPres+" route with status %d ", resp.StatusCode))
		}
	}
	return nil
}

func addBridgeToIngress(keptnDomain string, ingress *networking.Ingress) {

	containsBridge := false
	for i, x := range ingress.Spec.Rules {
		if strings.HasPrefix(x.Host, "bridge.keptn") {
			containsBridge = true
			ingress.Spec.Rules[i] = getBridgeRule(keptnDomain)
		}
	}
	if !containsBridge {
		ingress.Spec.Rules = append(ingress.Spec.Rules, getBridgeRule(keptnDomain))
	}
}

func getBridgeRule(keptnDomain string) networking.IngressRule {
	return networking.IngressRule{
		Host: getHostForBridge(keptnDomain),
		IngressRuleValue: networking.IngressRuleValue{
			HTTP: &networking.HTTPIngressRuleValue{
				Paths: []networking.HTTPIngressPath{
					{
						Backend: networking.IngressBackend{
							ServiceName: "bridge",
							ServicePort: intstr.IntOrString{IntVal: 8080},
						},
					},
				},
			},
		},
	}
}

func disposeBridgeIstio(l *keptnutils.Logger) error {
	l.Info("Dispose Bridge by deleting the VirtualService and DestinationRule")

	restConfig, _ := getRestConfig()
	ic, err := versionedclient.NewForConfig(restConfig)
	if err != nil {
		l.Error(fmt.Sprintf("Failed to create istio client: %s", err))
	}

	err = ic.NetworkingV1alpha3().VirtualServices("keptn").Delete("bridge", &metav1.DeleteOptions{})
	if k8serrors.IsNotFound(err) {
		l.Info("VirtualService did not exist")
	} else if err != nil {
		return err
	} else {
		l.Info("VirtualService for Bridge deleted")
	}
	err = ic.NetworkingV1alpha3().DestinationRules("keptn").Delete("bridge", &metav1.DeleteOptions{})
	if k8serrors.IsNotFound(err) {
		l.Info("DestinationRule did not exist")
	} else if err != nil {
		return err
	} else {
		l.Info("DestinationRule for Bridge deleted")
	}
	return nil
}

func disposeBridgeIngress(l *keptnutils.Logger) error {
	l.Info("Dispose Bridge of keptn-ingress")

	clientset, err := k8sutils.GetClientset(useInClusterConfig)
	if err != nil {
		return err
	}
	ing, err := clientset.NetworkingV1beta1().Ingresses("keptn").Get("keptn-ingress", metav1.GetOptions{})
	if err != nil {
		return err
	}
	l.Info("keptn-ingress retrieved")
	removeBridgeFromIngress(ing)

	_, err = clientset.NetworkingV1beta1().Ingresses("keptn").Update(ing)
	if err != nil {
		return err
	}
	l.Info("Rule of Bridge successfully deleted from Keptn-ingress")
	return nil
}

func removeBridgeFromIngress(ingress *networking.Ingress) {

	rules := ingress.Spec.Rules[:0]
	for _, x := range ingress.Spec.Rules {
		if !strings.HasPrefix(x.Host, "bridge.keptn") {
			rules = append(rules, x)
		}
	}
	ingress.Spec.Rules = rules
}

func getHostForBridge(keptnDomain string) string {
	// check if the domain contains a port. If yes, only the first part without the port will be used
	split := strings.Split(keptnDomain, ":")
	return "bridge.keptn." + split[0]
}
