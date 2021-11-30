package event_handler

import (
	"context"
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/url"
	"os"
	"strings"
	"sync"

	utils "github.com/keptn/go-utils/pkg/api/utils"
	keptn "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

const eventbroker = "EVENTBROKER"
const datastore = "MONGODB_DATASTORE"
const configurationServiceURL = "configuration-service:8080"

// Opaque key type used for graceful shutdown context value
type gracefulShutdownKeyType struct{}

var GracefulShutdownKey = gracefulShutdownKeyType{}

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

// ErrServiceNotFound godoc
var ErrConfigService = errors.New("could not checkout the SLO")

//go:generate moq -pkg event_handler_mock -skip-ensure -out ./fake/resource_handler_mock.go . ResourceHandler
type ResourceHandler interface {
	GetServiceResource(project string, stage string, service string, resourceURI string) (*keptnapimodels.Resource, error)
}

//go:generate moq -pkg event_handler_mock -skip-ensure -out ./fake/service_handler_mock.go . ServiceHandler
type ServiceHandler interface {
	GetService(project, stage, service string) (*keptnapimodels.Service, error)
}

//go:generate moq -pkg event_handler_mock -skip-ensure -out ./fake/event_store_mock.go . EventStore
type EventStore interface {
	GetEvents(filter *utils.EventFilter) ([]*keptnapimodels.KeptnContextExtendedCE, *keptnapimodels.Error)
}

type SLOFileRetriever struct {
	ResourceHandler ResourceHandler
	ServiceHandler  ServiceHandler
}

func (sr *SLOFileRetriever) GetSLOs(project, stage, service string) (*keptn.ServiceLevelObjectives, error) {
	sloFile, err := sr.ResourceHandler.GetServiceResource(project, stage, service, "slo.yaml")
	if err != nil {
		_, err2 := sr.ServiceHandler.GetService(project, stage, service)
		if err2 != nil {
			if strings.Contains(strings.ToLower(err2.Error()), "project not found") {
				return nil, ErrProjectNotFound
			} else if strings.Contains(strings.ToLower(err2.Error()), "stage not found") {
				return nil, ErrStageNotFound
			} else if strings.Contains(strings.ToLower(err2.Error()), "service not found") {
				return nil, ErrServiceNotFound
			}
		} else {
			if strings.Contains(strings.ToLower(err.Error()), "could not check out ") {
				return nil, ErrConfigService
			}
			return nil, ErrSLOFileNotFound
		}
	}
	if sloFile == nil || sloFile.ResourceContent == "" {
		return nil, ErrSLOFileNotFound
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

	objectives := []*keptn.SLO{}
	for _, objective := range slo.Objectives {
		if objective == nil {
			continue
		}
		if objective.Weight == 0 {
			objective.Weight = 1
		}
		objectives = append(objectives, objective)
	}
	slo.Objectives = objectives

	return slo, nil
}

func sendEvent(shkeptncontext string, triggeredID, eventType string, keptnHandler *keptnv2.Keptn, data interface{}) error {
	source, _ := url.Parse("lighthouse-service")

	event := cloudevents.NewEvent()
	event.SetType(eventType)
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", shkeptncontext)
	event.SetExtension("triggeredid", triggeredID)
	if data != nil {
		event.SetData(cloudevents.ApplicationJSON, data)
	}

	logger.Debug("Send event: " + eventType)
	return keptnHandler.SendCloudEvent(event)
}

func sendErroredFinishedEventWithMessage(shkeptncontext, triggeredID, message, sloFileContent string, keptnHandler *keptnv2.Keptn, incoming *keptnv2.GetSLIFinishedEventData) error {
	data := keptnv2.EvaluationFinishedEventData{
		EventData: keptnv2.EventData{
			Project: incoming.Project,
			Stage:   incoming.Stage,
			Service: incoming.Service,
			Labels:  incoming.Labels,
			Status:  keptnv2.StatusErrored,
			Result:  keptnv2.ResultFailed,
			Message: message,
		},
		Evaluation: keptnv2.EvaluationDetails{
			TimeStart:        incoming.GetSLI.Start,
			TimeEnd:          incoming.GetSLI.End,
			Result:           string(keptnv2.ResultFailed),
			Score:            0,
			SLOFileContent:   sloFileContent,
			IndicatorResults: nil,
			GitCommit:        "",
		},
	}
	return sendEvent(shkeptncontext, triggeredID, keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName), keptnHandler, data)
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
	kubeAPI, err := GetConfig().GetKubeAPI()
	if err != nil {
		return "", err
	}

	configMap, err := kubeAPI.CoreV1().ConfigMaps(namespace).Get(context.TODO(), "lighthouse-config", v1.GetOptions{})

	if err != nil {
		return "", errors.New("no default SLI provider specified")
	}

	sliProvider := configMap.Data["sli-provider"]

	return sliProvider, nil
}

// GetSLIProvider godoc
func (K8sSLIProviderConfig) GetSLIProvider(project string) (string, error) {
	kubeAPI, err := GetConfig().GetKubeAPI()
	if err != nil {
		return "", err
	}

	configMap, err := kubeAPI.CoreV1().ConfigMaps(namespace).Get(context.TODO(), "lighthouse-config-"+project, v1.GetOptions{})

	if err != nil {
		return "", errors.New("no SLI provider specified for project " + project)
	}

	sliProvider := configMap.Data["sli-provider"]

	return sliProvider, nil
}

type Config struct {
	GetKubeAPI KubeAPIConfigFunc
}

var config *Config
var configOnce sync.Once

func GetConfig() *Config {
	configOnce.Do(func() {
		config = &Config{GetKubeAPI: getInClusterKubeClient}
	})
	return config
}

type KubeAPIConfigFunc func() (kubernetes.Interface, error)

func getInClusterKubeClient() (kubernetes.Interface, error) {
	var config *rest.Config
	config, err := rest.InClusterConfig()

	if err != nil {
		return nil, err
	}

	kubeAPI, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return kubeAPI, nil
}
