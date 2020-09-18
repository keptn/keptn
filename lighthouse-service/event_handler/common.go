package event_handler

import (
	"errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	utils "github.com/keptn/go-utils/pkg/api/utils"
	keptn "github.com/keptn/go-utils/pkg/lib"
)

const eventbroker = "EVENTBROKER"
const datastore = "MONGODB_DATASTORE"
const configurationServiceURL = "configuration-service:8080"

func getDatastoreURL() string {
	if os.Getenv(datastore) != "" {
		return "http://" + os.Getenv(datastore)
	}
	return "http://mongodb-datastore:8080"
}

// ErrSLOFileNotFound godoc
var ErrSLOFileNotFound = errors.New("no slo file available")

// ErrProjectNotFound godoc
var ErrProjectNotFound = errors.New("project not found")

// ErrStageNotFound godoc
var ErrStageNotFound = errors.New("stage not found")

// ErrServiceNotFound godoc
var ErrServiceNotFound = errors.New("service not found")

func getSLOs(project string, stage string, service string) (*keptn.ServiceLevelObjectives, error) {
	resourceHandler := utils.NewResourceHandler(configurationServiceURL)
	sloFile, err := resourceHandler.GetServiceResource(project, stage, service, "slo.yaml")
	if err != nil {
		// check if service/stage/project actually exist
		serviceHandler := utils.NewServiceHandler(configurationServiceURL)
		_, err = serviceHandler.GetService(project, stage, service)
		if strings.Contains(strings.ToLower(err.Error()), "project not found") {
			return nil, ErrProjectNotFound
		} else if strings.Contains(strings.ToLower(err.Error()), "stage not found") {
			return nil, ErrStageNotFound
		} else if strings.Contains(strings.ToLower(err.Error()), "service not found") {
			return nil, ErrServiceNotFound
		} else {
			return nil, ErrSLOFileNotFound
		}
	}

	slo, err := parseSLO([]byte(sloFile.ResourceContent))

	if err != nil {
		return nil, errors.New("Could not parse SLO file for service " + service + " in stage " + stage + " in project " + project)
	}

	return slo, nil
}

func parseSLO(input []byte) (*keptn.ServiceLevelObjectives, error) {
	slo := &keptn.ServiceLevelObjectives{}
	err := yaml.Unmarshal([]byte(input), &slo)

	if err != nil {
		return nil, err
	}

	if slo.Comparison == nil {
		slo.Comparison = &keptn.SLOComparison{
			CompareWith:               "single_result",
			IncludeResultWithScore:    "all",
			NumberOfComparisonResults: 1,
			AggregateFunction:         "avg",
		}
	}

	if slo.Comparison != nil {
		if slo.Comparison.IncludeResultWithScore == "" {
			slo.Comparison.IncludeResultWithScore = "all"
		}
		if slo.Comparison.NumberOfComparisonResults == 0 {
			slo.Comparison.NumberOfComparisonResults = 3
		}
		if slo.Comparison.AggregateFunction == "" {
			slo.Comparison.AggregateFunction = "avg"
		}
	}

	for _, objective := range slo.Objectives {
		if objective.Weight == 0 {
			objective.Weight = 1
		}
	}

	return slo, nil
}

// SLIProviderConfig godoc
type SLIProviderConfig interface {
	GetDefaultSLIProvider() (string, error)
	GetSLIProvider(project string) (string, error)
}

// K8sSLIProviderConfig godoc
type K8sSLIProviderConfig struct{}

// GetDefaultSLIProvider godoc
func (K8sSLIProviderConfig) GetDefaultSLIProvider() (string, error) {
	kubeAPI, err := getKubeAPI()
	if err != nil {
		return "", err
	}

	configMap, err := kubeAPI.CoreV1().ConfigMaps(namespace).Get("lighthouse-config", v1.GetOptions{})

	if err != nil {
		return "", errors.New("no default SLI provider specified")
	}

	sliProvider := configMap.Data["sli-provider"]

	return sliProvider, nil
}

// GetSLIProvider godoc
func (K8sSLIProviderConfig) GetSLIProvider(project string) (string, error) {
	kubeAPI, err := getKubeAPI()
	if err != nil {
		return "", err
	}

	configMap, err := kubeAPI.CoreV1().ConfigMaps(namespace).Get("lighthouse-config-"+project, v1.GetOptions{})

	if err != nil {
		return "", errors.New("no SLI provider specified for project " + project)
	}

	sliProvider := configMap.Data["sli-provider"]

	return sliProvider, nil
}
