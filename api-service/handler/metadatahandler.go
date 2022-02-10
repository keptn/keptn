package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	logger "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/api-service/model"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type IMetadataHandler interface {
	GetMetadata(c *gin.Context)
}

const defaultVersion = "N/A"

func NewMetadataHandler() *MetadataHandler {
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
	return &MetadataHandler{
		k8sClient:       clientSet,
		swaggerFilePath: "/swagger-ui/swagger.yaml",
	}
}

type swaggerFileProvider interface {
	getSwaggerFileContent() ([]byte, error)
}

type MetadataHandler struct {
	k8sClient kubernetes.Interface
	swaggerFileProvider
	swaggerFilePath string
}

func (eh *MetadataHandler) GetMetadata(c *gin.Context) {
	logger.Info("API received a GET metadata event")

	namespace := os.Getenv("POD_NAMESPACE")

	var payload model.Metadata
	payload.Namespace = namespace
	payload.Keptnversion = defaultVersion
	payload.Keptnlabel = "keptn"
	payload.Bridgeversion = defaultVersion
	payload.Shipyardversion = "0.2.0"

	if bridgeVersion, err := eh.getBridgeVersion(namespace); err != nil {
		logger.WithError(err).Error("Error getting bridge version")
	} else {
		payload.Bridgeversion = bridgeVersion
	}

	if keptnVersion, err := eh.getSwaggerKeptnVersion(); err != nil {
		logger.WithError(err).Error("Error getting swagger info")
	} else {
		payload.Keptnversion = keptnVersion
	}

	logger.Info("!!!!aom na konci a payload je ", payload)
	c.Status(http.StatusOK)
	c.JSON(http.StatusOK, payload)
}

func (h *MetadataHandler) getBridgeVersion(namespace string) (string, error) {

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

func (h *MetadataHandler) getSwaggerKeptnVersion() (string, error) {
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
