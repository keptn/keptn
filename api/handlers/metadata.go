package handlers

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"gopkg.in/yaml.v2"

	keptnutils "github.com/keptn/go-utils/pkg/lib/keptn"
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
	logger := keptnutils.NewLogger("", "", "api")
	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Debug(fmt.Sprintf("Could not get InClusterConfig, will skip k8s-deployments: %s", err.Error()))
	} else {
		// creates the clientset
		clientSet, err = kubernetes.NewForConfig(config)
		if err != nil {
			logger.Debug("Could not create kubernetes client, will skip k8s-deployments.")
		}
	}
	return metadataHandler{
		k8sClient:       clientSet,
		logger:          logger,
		swaggerFilePath: "/swagger-ui/swagger.yaml",
	}
}

type swaggerFileProvider interface {
	getSwaggerFileContent() ([]byte, error)
}

type metadataHandler struct {
	k8sClient kubernetes.Interface
	swaggerFileProvider
	logger          keptnutils.LoggerInterface
	swaggerFilePath string
}

func (h *metadataHandler) getMetadata() middleware.Responder {
	h.logger.Info("API received a GET metadata event")

	var namespace string
	namespace = os.Getenv("POD_NAMESPACE")

	var payload models.Metadata
	payload.Namespace = namespace
	payload.Keptnversion = "N/A"
	payload.Keptnlabel = "keptn"
	payload.Bridgeversion = "N/A"

	if h.k8sClient != nil {
		deploymentsClient := h.k8sClient.AppsV1().Deployments(namespace)
		bridgeDeployment, err := deploymentsClient.Get("bridge", metav1.GetOptions{})
		if err != nil {
			// log the error, but continue
			h.logger.Error(fmt.Sprintf("Error getting deployment info: %s", err.Error()))
		} else {
			payload.Bridgeversion = strings.TrimPrefix(bridgeDeployment.Spec.Template.Spec.Containers[0].Image, "keptn/")
		}

		// Load swagger.yaml from /swagger-ui/swagger.yaml
		mapSwagger := make(map[interface{}]interface{})
		yamlFile, err := ioutil.ReadFile(h.swaggerFilePath)

		if err != nil {
			fmt.Printf("yamlFile.Get err   #%v ", err)
		}
		err = yaml.Unmarshal(yamlFile, &mapSwagger)
		if err != nil {
			fmt.Printf("Unmarshal: %v", err)
		}
		info := mapSwagger["info"].(map[interface{}]interface{})

		for k, v := range info {
			if k == "version" {
				payload.Keptnversion = fmt.Sprintf("%v", v)
			}
		}
	}

	return metadata.NewMetadataOK().WithPayload(&payload)
}
