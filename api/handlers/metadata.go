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

	var namespace string
	namespace = os.Getenv("POD_NAMESPACE")

	var payload models.Metadata
	payload.Namespace = namespace
	payload.Keptnversion = "N/A"
	payload.Keptnlabel = "keptn"
	payload.Bridgeversion = "N/A"
	payload.Shipyardversion = "0.2.0"

	if h.k8sClient != nil {
		deploymentsClient := h.k8sClient.AppsV1().Deployments(namespace)

		if bridgeDeployment, err := deploymentsClient.Get(context.TODO(), "bridge", metav1.GetOptions{}); err == nil {
			v := strings.Split(bridgeDeployment.Spec.Template.Spec.Containers[0].Image, ":")
			if len(v) >= 2 {
				payload.Bridgeversion = v[1]
			}
		} else {
			logger.WithError(err).Error("Error getting deployment info")
		}

		if info, err := h.getSwaggerInfo(); err == nil {
			for k, v := range info {
				if k == "version" {
					payload.Keptnversion = fmt.Sprintf("%v", v)
				}
			}
		} else {
			logger.WithError(err).Error("Error getting swagger info")
		}
	}

	return metadata.NewMetadataOK().WithPayload(&payload)
}

func (h *metadataHandler) getSwaggerInfo() (map[interface{}]interface{}, error) {
	// Load swagger.yaml from /swagger-ui/swagger.yaml
	mapSwagger := make(map[interface{}]interface{})
	yamlFile, err := ioutil.ReadFile(h.swaggerFilePath)

	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &mapSwagger)
	if err != nil {
		return nil, err
	}
	info := mapSwagger["info"].(map[interface{}]interface{})
	return info, nil
}
