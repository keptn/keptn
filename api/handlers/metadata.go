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
func GetMetadataHandlerFunc(params metadata.MetadataParams, pricipal *models.Principal) middleware.Responder {

	logger := keptnutils.NewLogger("", "", "api")
	logger.Info("API received a GET metadata event")

	var namespace string
	namespace = os.Getenv("POD_NAMESPACE")

	var checkDeployment bool
	checkDeployment = true

	var payload models.Metadata
	payload.Namespace = namespace
	payload.Keptnversion = "N/A"
	payload.Keptnlabel = "keptn"

	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Debug("Could not get InClusterConfig, will skip k8s-deployments.")
		checkDeployment = false
	} else {
		// creates the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			logger.Debug("Could not create kubernetes client, will skip k8s-deployments.")
			payload.Bridgeversion = "N/A"
			checkDeployment = false
		} else {
			deploymentsClient := clientset.AppsV1().Deployments(namespace)
			bridgeDeployment, err := deploymentsClient.Get("bridge", metav1.GetOptions{})
			if err != nil {
				logger.Error(fmt.Sprintf("Error getting deployment info %s", err.Error()))
				return nil
			}
			payload.Bridgeversion = strings.TrimPrefix(bridgeDeployment.Spec.Template.Spec.Containers[0].Image, "keptn/")
		}
	}

	if checkDeployment {
		// Load swagger.yaml from /swagger-ui/swagger.yaml
		mapSwagger := make(map[interface{}]interface{})
		yamlFile, err := ioutil.ReadFile("/swagger-ui/swagger.yaml")

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
	} else {
		payload.Keptnversion = "N/A"
	}

	return metadata.NewMetadataOK().WithPayload(&payload)
}
