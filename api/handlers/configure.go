package handlers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	keptnutils "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/configuration"
	k8sutils "github.com/keptn/kubernetes-utils/pkg"
)

const (
	useInClusterConfig      = true
	bridgeCredentialsSecret = "bridge-credentials"
)

var namespace = os.Getenv("POD_NAMESPACE")

// PostConfigureBridgeHandlerFunc handler function for POST requests
func PostConfigureBridgeHandlerFunc(params configuration.PostConfigBridgeParams, principal *models.Principal) middleware.Responder {
	l := keptnutils.NewLogger("", "", "api")
	l.Info("API received a configure Bridge POST request")

	if params.ConfigureBridge.User == nil || *params.ConfigureBridge.User == "" {
		errMsg := fmt.Sprintf("no user provided")
		l.Error(errMsg)
		return configuration.NewPostConfigBridgeDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(errMsg)})
	}
	if params.ConfigureBridge.Password == nil || *params.ConfigureBridge.Password == "" {
		errMsg := fmt.Sprintf("no password provided")
		l.Error(errMsg)
		return configuration.NewPostConfigBridgeDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(errMsg)})
	}

	err := createBridgeCredentials(*params.ConfigureBridge.User, *params.ConfigureBridge.Password, l)
	if err != nil {
		errMsg := fmt.Sprintf("failed to create bridge credentials: %v", err)
		l.Error(errMsg)
		return configuration.NewPostConfigBridgeDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(errMsg)})
	}

	err = restartBridgePod(l)
	if err != nil {
		errMsg := fmt.Sprintf("failed to restart bridge pod: %v", err)
		l.Error(errMsg)
		return configuration.NewPostConfigBridgeDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(errMsg)})
	}

	return configuration.NewPostConfigBridgeOK()
}

// GetConfigureBridgeHandlerFunc handler function for GET requests
func GetConfigureBridgeHandlerFunc(params configuration.GetConfigBridgeParams, principal *models.Principal) middleware.Responder {
	l := keptnutils.NewLogger("", "", "api")
	l.Info("API received a configuration Bridge GET request")

	k8s, err := getK8sClient()
	if err != nil {
		l.Error("Could not create k8s client: " + err.Error())
		return configuration.NewGetConfigBridgeDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not read bridge credentials")})
	}

	l.Info("Checking for existing secret")
	bridgeCredentials, err := k8s.CoreV1().Secrets(namespace).Get(bridgeCredentialsSecret, metav1.GetOptions{})

	if err != nil {
		l.Error(err.Error())
		return configuration.NewGetConfigBridgeDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not read bridge credentials")})
	}

	if bridgeCredentials.Data["BASIC_AUTH_PASSWORD"] == nil || len(bridgeCredentials.Data["BASIC_AUTH_PASSWORD"]) == 0 {
		return configuration.NewGetConfigBridgeDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not read bridge credentials: no password found")})
	}

	if bridgeCredentials.Data["BASIC_AUTH_USERNAME"] == nil || len(bridgeCredentials.Data["BASIC_AUTH_USERNAME"]) == 0 {
		return configuration.NewGetConfigBridgeDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not read bridge credentials: no username found")})
	}

	/*
		userDecoded, err := b64.StdEncoding.DecodeString(string(bridgeCredentials.Data["BASIC_AUTH_USERNAME"]))
		if err != nil {
			l.Error("Could not decode username: " + err.Error())
			return configuration.NewGetConfigBridgeDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not red bridge credentials: could not decode username")})
		}
		passwordDecoded, err := b64.StdEncoding.DecodeString(string(bridgeCredentials.Data["BASIC_AUTH_PASSWORD"]))
		if err != nil {
			l.Error("Could not decode password: " + err.Error())
			return configuration.NewGetConfigBridgeDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not red bridge credentials: could not decode password")})
		}

	*/

	creds := &models.ConfigureBridge{
		Password: swag.String(string(bridgeCredentials.Data["BASIC_AUTH_PASSWORD"])),
		User:     swag.String(string(bridgeCredentials.Data["BASIC_AUTH_USERNAME"])),
	}

	return configuration.NewGetConfigBridgeOK().WithPayload(creds)
}

func getK8sClient() (*kubernetes.Clientset, error) {
	restConfig, _ := getRestConfig()

	k8s, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return k8s, nil
}

func restartBridgePod(l *keptnutils.Logger) error {
	l.Info("Restarting bridge for credentials to take effect")

	k8s, err := getK8sClient()
	if err != nil {
		return err
	}

	return k8s.CoreV1().Pods(namespace).DeleteCollection(&metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=bridge",
	})
}

func createBridgeCredentials(user string, password string, l *keptnutils.Logger) error {
	l.Info("Creating or updating credentials for bridge")

	k8s, err := getK8sClient()
	if err != nil {
		return err
	}

	l.Info("Checking for existing secret")
	bridgeCredentials, err := k8s.CoreV1().Secrets(namespace).Get(bridgeCredentialsSecret, metav1.GetOptions{})
	if err == nil && bridgeCredentials != nil {
		// update existing secret
		l.Info("Existing secret found. Updating with new values for user and password")
		newSecret := getBridgeCredentials(user, password)
		bridgeCredentials.Data = newSecret.Data
		_, err = k8s.CoreV1().Secrets(namespace).Update(newSecret)
		if err != nil {
			l.Error("could not update secret: " + err.Error())
			return err
		}
	} else {
		l.Info("Creating a new secret")
		newSecret := getBridgeCredentials(user, password)
		_, err = k8s.CoreV1().Secrets(namespace).Create(newSecret)
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
			Name:      bridgeCredentialsSecret,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"BASIC_AUTH_USERNAME": []byte(user),
			"BASIC_AUTH_PASSWORD": []byte(password),
		},
		Type: "Opaque",
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
