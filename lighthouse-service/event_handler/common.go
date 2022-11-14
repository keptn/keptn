package event_handler

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	utils "github.com/keptn/go-utils/pkg/api/utils"
	keptn "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

const datastore = "MONGODB_DATASTORE"
const configurationServiceURL = "resource-service:8080"

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
	GetResource(scope utils.ResourceScope, options ...utils.URIOption) (*apimodels.Resource, error)
}

//go:generate moq -pkg event_handler_mock -skip-ensure -out ./fake/service_handler_mock.go . ServiceHandler
type ServiceHandler interface {
	GetService(project, stage, service string) (*apimodels.Service, error)
}

//go:generate moq -pkg event_handler_mock -skip-ensure -out ./fake/event_store_mock.go . EventStore
type EventStore interface {
	GetEvents(filter *utils.EventFilter) ([]*apimodels.KeptnContextExtendedCE, *apimodels.Error)
}

type SLOFileRetriever struct {
	ResourceHandler ResourceHandler
	ServiceHandler  ServiceHandler
}

func (sr *SLOFileRetriever) GetSLOs(project, stage, service, commitID string) (*keptn.ServiceLevelObjectives, []byte, error) {
	commitOption := url.Values{}
	if commitID != "" {
		commitOption.Add("gitCommitID", commitID)
	}
	resourceScope := *utils.NewResourceScope().Project(project).Stage(stage).Service(service).Resource("slo.yaml")
	sloFile, err := sr.ResourceHandler.GetResource(resourceScope, utils.AppendQuery(commitOption))
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return nil, nil, err
		}
		_, serviceErr := sr.ServiceHandler.GetService(project, stage, service)
		if serviceErr != nil {
			return nil, nil, checkNotFound(serviceErr, err)
		}
	}
	if sloFile == nil || sloFile.ResourceContent == "" {
		return nil, nil, ErrSLOFileNotFound
	}

	slo, err := parseSLO([]byte(sloFile.ResourceContent))

	if err != nil {
		return nil, nil, errors.New("Could not parse SLO file for service " + service + " in stage " + stage + " in project " + project)
	}
	// return also slo.yaml as a plain file to avoid confusion due to defaulted values (see https://github.com/keptn/keptn/issues/1495)
	return slo, []byte(sloFile.ResourceContent), nil
}

func checkNotFound(notFound, checkOut error) error {
	if strings.Contains(strings.ToLower(notFound.Error()), "project not found") {
		return ErrProjectNotFound
	} else if strings.Contains(strings.ToLower(notFound.Error()), "stage not found") {
		return ErrStageNotFound
	} else if strings.Contains(strings.ToLower(notFound.Error()), "service not found") {
		return ErrServiceNotFound
	} else {
		if strings.Contains(strings.ToLower(checkOut.Error()), "could not check out ") {
			return ErrConfigService
		}
		return checkOut
	}
}

func parseSLO(input []byte) (*keptn.ServiceLevelObjectives, error) {
	slo := &keptn.ServiceLevelObjectives{}
	err := yaml.Unmarshal([]byte(input), &slo)

	if err != nil {
		return nil, fmt.Errorf(err.Error())
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

func sendEvent(shkeptncontext string, triggeredID, eventType, commitID string, keptnHandler *keptnv2.Keptn, data interface{}) error {
	source, _ := url.Parse("lighthouse-service")

	event := cloudevents.NewEvent()
	event.SetType(eventType)
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetTime(time.Now().UTC())
	event.SetID(uuid.New().String())
	event.SetExtension("shkeptncontext", shkeptncontext)
	event.SetExtension("triggeredid", triggeredID)
	event.SetExtension("gitcommitid", commitID)
	if data != nil {
		_ = event.SetData(cloudevents.ApplicationJSON, data)
	}

	logger.Debug("Send event: " + eventType)
	return keptnHandler.SendCloudEvent(event)
}

func sendErroredFinishedEventWithMessage(shkeptncontext, triggeredID, commitID, message, sloFileContent string, keptnHandler *keptnv2.Keptn, incoming *keptnv2.GetSLIFinishedEventData) error {
	encodedSLOFileContent := base64.StdEncoding.EncodeToString([]byte(sloFileContent))
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
			SLOFileContent:   encodedSLOFileContent,
			IndicatorResults: nil,
		},
	}
	return sendEvent(shkeptncontext, triggeredID, keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName), commitID, keptnHandler, data)
}

// SLIProviderConfig godoc
type SLIProviderConfig interface {
	GetDefaultSLIProvider() (string, error)
	GetSLIProvider(project string) (string, error)
}

// K8sSLIProviderConfig godoc
type K8sSLIProviderConfig struct {
	KubeAPI kubernetes.Interface
}

func NewSLIProviderConfig(kubeAPI kubernetes.Interface) K8sSLIProviderConfig {
	return K8sSLIProviderConfig{
		KubeAPI: kubeAPI,
	}
}

// GetDefaultSLIProvider godoc
func (k K8sSLIProviderConfig) GetDefaultSLIProvider() (string, error) {

	configMap, err := k.KubeAPI.CoreV1().ConfigMaps(namespace).Get(context.TODO(), "lighthouse-config", v1.GetOptions{})

	if err != nil {
		return "", errors.New("no default SLI provider specified")
	}

	sliProvider := configMap.Data["sli-provider"]

	return sliProvider, nil
}

// GetSLIProvider godoc
func (k K8sSLIProviderConfig) GetSLIProvider(project string) (string, error) {

	configMap, err := k.KubeAPI.CoreV1().ConfigMaps(namespace).Get(context.TODO(), "lighthouse-config-"+project, v1.GetOptions{})

	if err != nil {
		return "", errors.New("no SLI provider specified for project " + project)
	}

	sliProvider := configMap.Data["sli-provider"]

	return sliProvider, nil
}
