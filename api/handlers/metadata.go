package handlers

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	logger "github.com/sirupsen/logrus"

	"github.com/go-openapi/runtime/middleware"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/metadata"
)

const defaultVersion = "N/A"

// Swagger Structure

// GetMetadataHandlerFunc returns metadata of the keptn installation
func GetMetadataHandlerFunc(params metadata.MetadataParams, principal *models.Principal) middleware.Responder {

	handler := newMetadataHandler()

	return handler.getMetadata()
}

func newMetadataHandler() metadataHandler {
	var clientSet kubernetes.Interface

	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Debugf("Could not get InClusterConfig, will skip k8s-deployments: %s", err.Error())
	} else {
		// creates the clientset
		clientSet, err = kubernetes.NewForConfig(config)
		if err != nil {
			logger.Debug("Could not create kubernetes client, will skip k8s-deployments.")
		}
	}
	return metadataHandler{
		k8sClient:       clientSet,
		swaggerFilePath: "/swagger-ui/swagger.yaml",
	}
}

type swaggerFileProvider interface {
	getSwaggerFileContent() ([]byte, error)
}

type metadataHandler struct {
	k8sClient kubernetes.Interface
	swaggerFileProvider
	swaggerFilePath string
}

func (h *metadataHandler) getMetadata() middleware.Responder {
	logger.Info("API received a GET metadata event")

	namespace := os.Getenv("POD_NAMESPACE")
	automaticProv := os.Getenv("AUTOMATIC_PROVISIONING_URL")

	var payload models.Metadata
	payload.Namespace = namespace
	payload.Keptnversion = defaultVersion
	payload.Keptnlabel = "keptn"
	payload.Bridgeversion = defaultVersion
	payload.Shipyardversion = "0.2.0"
	payload.Automaticprovisioning = automaticProv != ""

	if bridgeVersion, err := h.getBridgeVersion(namespace); err != nil {
		logger.WithError(err).Error("Error getting bridge version")
	} else {
		payload.Bridgeversion = bridgeVersion
	}

	if keptnVersion, err := h.getSwaggerKeptnVersion(); err != nil {
		logger.WithError(err).Error("Error getting swagger info")
	} else {
		payload.Keptnversion = keptnVersion
	}

	return metadata.NewMetadataOK().WithPayload(&payload)
}

func (h *metadataHandler) getBridgeVersion(namespace string) (string, error) {

	if h.k8sClient == nil {
		return "", fmt.Errorf("unable to get bridge version")
	}
	bridgeDeployment, err := h.k8sClient.AppsV1().Deployments(namespace).Get(context.TODO(), "bridge", metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("unable to get bridge version %w", err)
	}

	v := strings.Split(bridgeDeployment.Spec.Template.Spec.Containers[0].Image, ":")
	if len(v) >= 2 {
		return v[1], nil
	}

	return "", fmt.Errorf("unable to get bridge version")
}

func (h *metadataHandler) getSwaggerKeptnVersion() (string, error) {
	// Load swagger.yaml from /swagger-ui/swagger.yaml
	mapSwagger := make(map[interface{}]interface{})
	yamlFile, err := ioutil.ReadFile(h.swaggerFilePath)

	if err != nil {
		return defaultVersion, err
	}
	err = yaml.Unmarshal(yamlFile, &mapSwagger)
	if err != nil {
		return defaultVersion, err
	}
	info := mapSwagger["info"].(map[interface{}]interface{})
	for k, v := range info {
		if k == "version" {
			return fmt.Sprintf("%v", v), nil
		}
	}
	return defaultVersion, nil
}
