package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/ghodss/yaml"
	"github.com/google/uuid"
	configmodels "github.com/keptn/go-utils/pkg/api/models"
	configutils "github.com/keptn/go-utils/pkg/api/utils"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const remediationFileName = "remediation.yaml"
const configurationserviceconnection = "CONFIGURATION_SERVICE" //"localhost:6060" // "configuration-service:8080"
const datastoreConnection = "MONGODB_DATASTORE"
const remediationSpecVersion = "0.2.0"

// Handler handles incoming Keptn events
type Handler interface {
	HandleEvent() error
}

// NewHandler returns a new Handler for the incoming Keptn event
func NewHandler(event cloudevents.Event) (Handler, error) {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptn.NewLogger(shkeptncontext, event.Context.GetID(), "remediationStatus-service")
	logger.Debug("Received event for shkeptncontext:" + shkeptncontext)

	keptnHandler, err := keptn.NewKeptn(&event, keptn.KeptnOpts{})
	if err != nil {
		logger.Error("Could not initialize Keptn handler: " + err.Error())
		return nil, err
	}

	switch event.Type() {
	case keptn.ProblemOpenEventType:
		return &ProblemOpenEventHandler{
			KeptnHandler: keptnHandler,
			Logger:       logger,
			Event:        event,
			Remediation: &Remediation{
				Keptn:  keptnHandler,
				Logger: logger,
			},
		}, nil
	case keptn.EvaluationDoneEventType:
		return &EvaluationDoneEventHandler{
			KeptnHandler: keptnHandler,
			Logger:       logger,
			Event:        event,
			Remediation: &Remediation{
				Keptn:  keptnHandler,
				Logger: logger,
			},
		}, nil
	case keptn.ActionFinishedEventType:
		return &ActionFinishedEventHandler{
			KeptnHandler: keptnHandler,
			Logger:       logger,
			Event:        event,
			Remediation: &Remediation{
				Keptn:  keptnHandler,
				Logger: logger,
			},
		}, nil
	default:
		return nil, errors.New("no event handler found for type: " + event.Type())
	}
}

type remediationStatus struct {
	// Executed action
	Action string `json:"action,omitempty"`

	// ID of the event
	EventID string `json:"eventId,omitempty"`

	// Keptn Context ID of the event
	KeptnContext string `json:"keptnContext,omitempty"`

	// Time of the event
	Time string `json:"time,omitempty"`

	// Type of the event
	Type string `json:"type,omitempty"`
}

type remediationStatusList struct {
	// Pointer to next page, base64 encoded
	NextPageKey string `json:"nextPageKey,omitempty"`

	// Size of returned page
	PageSize float64 `json:"pageSize,omitempty"`

	// remediations
	Remediations []*remediationStatus `json:"remediations"`

	// Total number of stages
	TotalCount float64 `json:"totalCount,omitempty"`
}

// Remediation provides functions to access all resources related to the remediation workflow
type Remediation struct {
	Keptn  *keptn.Keptn
	Logger keptn.LoggerInterface
}

func (r *Remediation) getActionForProblemType(remediationData keptn.RemediationV02, problemType string, index int) *keptn.RemediationV02ActionsOnOpen {
	for _, remediation := range remediationData.Spec.Remediations {
		if strings.HasPrefix(problemType, remediation.ProblemType) {
			r.Logger.Info("Found remediation for problem type " + remediation.ProblemType)
			if len(remediation.ActionsOnOpen) > index {
				return &remediation.ActionsOnOpen[index]
			}
		}
	}
	return nil
}

func (r *Remediation) sendRemediationTriggeredEvent(problemDetails *keptn.ProblemEventData) error {
	source, _ := url.Parse("remediation-service")
	contentType := "application/json"

	remediationFinishedEventData := &keptn.RemediationTriggeredEventData{
		Project: r.Keptn.KeptnBase.Project,
		Service: r.Keptn.KeptnBase.Service,
		Stage:   r.Keptn.KeptnBase.Stage,
		Problem: keptn.ProblemDetails{
			State:          problemDetails.State,
			ProblemID:      problemDetails.ProblemID,
			ProblemTitle:   problemDetails.ProblemTitle,
			ProblemDetails: problemDetails.ProblemDetails,
			PID:            problemDetails.PID,
			ProblemURL:     problemDetails.ProblemURL,
			ImpactedEntity: problemDetails.ImpactedEntity,
			Tags:           problemDetails.Tags,
		},
		Labels: r.Keptn.KeptnBase.Labels,
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptn.RemediationTriggeredEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": r.Keptn.KeptnContext},
		}.AsV02(),
		Data: remediationFinishedEventData,
	}

	err := createRemediation(event.ID(), r.Keptn.KeptnContext, event.Time().String(), *r.Keptn.KeptnBase, keptn.RemediationTriggeredEventType, "")
	if err != nil {
		r.Logger.Error("Could not create remediation: " + err.Error())
		return err
	}
	err = r.Keptn.SendCloudEvent(event)
	if err != nil {
		r.Logger.Error("Could not send action.finished event: " + err.Error())
		return err
	}

	return nil
}

func getRemediationsEndpoint(configurationServiceEndpoint url.URL, project, stage, service, keptnContext string) string {
	if keptnContext == "" {
		return fmt.Sprintf("%s://%s/v1/project/%s/stage/%s/service/%s/remediation", configurationServiceEndpoint.Scheme, configurationServiceEndpoint.Host, project, stage, service)
	}
	return fmt.Sprintf("%s://%s/v1/project/%s/stage/%s/service/%s/remediation/%s", configurationServiceEndpoint.Scheme, configurationServiceEndpoint.Host, project, stage, service, keptnContext)
}

func createRemediation(eventID, keptnContext, time string, keptnBase keptn.KeptnBase, remediationEventType, action string) error {
	configurationServiceEndpoint, err := keptn.GetServiceEndpoint(configurationserviceconnection)
	if err != nil {
		return errors.New("could not retrieve configuration-service URL")
	}

	newRemediation := &remediationStatus{
		Action:       action,
		EventID:      eventID,
		KeptnContext: keptnContext,
		Time:         time,
		Type:         remediationEventType,
	}

	queryURL := getRemediationsEndpoint(configurationServiceEndpoint, keptnBase.Project, keptnBase.Stage, keptnBase.Service, "")
	client := &http.Client{}
	payload, err := json.Marshal(newRemediation)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", queryURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.New(string(body))
	}

	return nil
}

func deleteRemediation(keptnContext string, keptnBase keptn.KeptnBase) error {
	configurationServiceEndpoint, err := keptn.GetServiceEndpoint(configurationserviceconnection)
	if err != nil {
		return errors.New("could not retrieve configuration-service URL")
	}

	queryURL := getRemediationsEndpoint(configurationServiceEndpoint, keptnBase.Project, keptnBase.Stage, keptnBase.Service, keptnContext)
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", queryURL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.New(string(body))
	}

	return nil
}

func getRemediationsByContext(keptnContext string, keptnBase keptn.KeptnBase) ([]*remediationStatus, error) {
	configurationServiceEndpoint, err := keptn.GetServiceEndpoint(configurationserviceconnection)
	if err != nil {
		return nil, errors.New("could not retrieve configuration-service URL")
	}

	remediations := []*remediationStatus{}

	nextPageKey := ""

	for {
		queryURL := getRemediationsEndpoint(configurationServiceEndpoint, keptnBase.Project, keptnBase.Stage, keptnBase.Service, keptnContext)

		url, err := url.Parse(queryURL)
		if err != nil {
			return nil, err
		}
		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
			url.RawQuery = q.Encode()
		}
		client := &http.Client{}

		req, err := http.NewRequest("GET", url.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("content-type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			return nil, errors.New(string(body))
		}

		remediationList := &remediationStatusList{}

		err = json.Unmarshal(body, remediationList)
		if err != nil {
			return nil, err
		}

		remediations = append(remediations, remediationList.Remediations...)

		if remediationList.NextPageKey == "" || remediationList.NextPageKey == "0" {
			break
		}
		nextPageKey = remediationList.NextPageKey
	}
	return remediations, nil

}

func (r *Remediation) sendStartEvaluationEvent() error {
	source, _ := url.Parse("remediation-service")
	contentType := "application/json"

	startEvaluationEventData := &keptn.StartEvaluationEventData{
		Project:      r.Keptn.KeptnBase.Project,
		Service:      r.Keptn.KeptnBase.Service,
		Stage:        r.Keptn.KeptnBase.Stage,
		Labels:       r.Keptn.KeptnBase.Labels,
		TestStrategy: "real-user",
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptn.StartEvaluationEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": r.Keptn.KeptnContext},
		}.AsV02(),
		Data: startEvaluationEventData,
	}

	err := r.Keptn.SendCloudEvent(event)
	if err != nil {
		r.Logger.Error("Could not send astart-evaluation event: " + err.Error())
		return err
	}
	return nil
}

func (r *Remediation) sendRemediationFinishedEvent(status keptn.RemediationStatusType, result keptn.RemediationResultType, message string) error {
	source, _ := url.Parse("remediation-service")
	contentType := "application/json"

	remediationFinishedEventData := &keptn.RemediationFinishedEventData{
		Project: r.Keptn.KeptnBase.Project,
		Service: r.Keptn.KeptnBase.Service,
		Stage:   r.Keptn.KeptnBase.Stage,
		Problem: keptn.ProblemDetails{},
		Labels:  r.Keptn.KeptnBase.Labels,
		Remediation: keptn.RemediationFinished{
			Status:  status,
			Result:  result,
			Message: message,
		},
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptn.RemediationFinishedEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": r.Keptn.KeptnContext},
		}.AsV02(),
		Data: remediationFinishedEventData,
	}

	err := r.Keptn.SendCloudEvent(event)
	if err != nil {
		r.Logger.Error("Could not send action.finished event: " + err.Error())
		return err
	}
	return nil
}

func (r *Remediation) getActionTriggeredEventData(problemDetails keptn.ProblemDetails, action *keptn.RemediationV02ActionsOnOpen) (keptn.ActionTriggeredEventData, error) {
	return keptn.ActionTriggeredEventData{
		Project: r.Keptn.KeptnBase.Project,
		Service: r.Keptn.KeptnBase.Service,
		Stage:   r.Keptn.KeptnBase.Stage,
		Action: keptn.ActionInfo{
			Name:        action.Name,
			Action:      action.Action,
			Description: action.Description,
			Value:       action.Value,
		},
		Problem: problemDetails,
		Labels:  nil,
	}, nil
}

func (r *Remediation) sendActionTriggeredEvent(actionTriggeredEventData keptn.ActionTriggeredEventData) error {

	source, _ := url.Parse("remediation-service")
	contentType := "application/json"

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptn.ActionTriggeredEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": r.Keptn.KeptnContext},
		}.AsV02(),
		Data: actionTriggeredEventData,
	}

	err := r.Keptn.SendCloudEvent(event)
	if err != nil {
		r.Logger.Error("Could not send action.finished event: " + err.Error())
		return err
	}
	return nil
}

func (r *Remediation) sendRemediationStatusChangedEvent(action *keptn.RemediationV02ActionsOnOpen, actionIndex int) error {

	remediationStatusChangedEventData := &keptn.RemediationStatusChangedEventData{
		Project: r.Keptn.KeptnBase.Project,
		Service: r.Keptn.KeptnBase.Service,
		Stage:   r.Keptn.KeptnBase.Stage,
		Labels:  r.Keptn.KeptnBase.Labels,
		Remediation: keptn.RemediationStatusChanged{
			Status: keptn.RemediationStatusSucceeded,
			Result: keptn.RemediationResult{
				ActionIndex: actionIndex,
				ActionName:  action.Action,
			},
		},
	}

	source, _ := url.Parse("remediation-service")
	contentType := "application/json"

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptn.RemediationStatusChangedEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": r.Keptn.KeptnContext},
		}.AsV02(),
		Data: remediationStatusChangedEventData,
	}

	err := createRemediation(event.ID(), r.Keptn.KeptnContext, event.Time().String(), *r.Keptn.KeptnBase, keptn.RemediationStatusChangedEventType, action.Action)
	if err != nil {
		r.Logger.Error("Could not create remediation: " + err.Error())
		return err
	}
	err = r.Keptn.SendCloudEvent(event)
	if err != nil {
		r.Logger.Error("Could not send remediation.status.changed event: " + err.Error())
		return err
	}
	return nil
}

func (r *Remediation) getRemediationFile() (*configmodels.Resource, error) {
	resourceHandler := configutils.NewResourceHandler(os.Getenv(configurationserviceconnection))
	var resource *configmodels.Resource
	var err error
	if r.Keptn.KeptnBase.Service != "" {
		resource, err = resourceHandler.GetServiceResource(r.Keptn.KeptnBase.Project, r.Keptn.KeptnBase.Stage,
			r.Keptn.KeptnBase.Service, remediationFileName)
	} else {
		resource, err = resourceHandler.GetStageResource(r.Keptn.KeptnBase.Project, r.Keptn.KeptnBase.Stage, remediationFileName)
	}

	if err != nil {
		msg := "remediation file not configured"
		r.Logger.Error(msg)
		_ = r.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return nil, err
	}
	r.Logger.Debug("remediation.yaml for service found")
	return resource, nil
}

func (r *Remediation) getRemediation(resource *configmodels.Resource) (*keptn.RemediationV02, error) {
	remediationData := &keptn.RemediationV02{}
	err := yaml.Unmarshal([]byte(resource.ResourceContent), remediationData)
	if err != nil {
		msg := "could not parse remediation.yaml"
		r.Logger.Error(msg + ": " + err.Error())
		_ = r.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return nil, err
	}

	if remediationData.Version != remediationSpecVersion {
		msg := "remediation.yaml file does not conform to remediation spec v0.2.0"
		r.Logger.Error(msg)
		_ = r.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return nil, errors.New(msg)
	}
	return remediationData, nil
}

func (r *Remediation) triggerAction(action *keptn.RemediationV02ActionsOnOpen, actionIndex int, problemDetails keptn.ProblemDetails) error {
	err := r.sendRemediationStatusChangedEvent(action, actionIndex)
	if err != nil {
		msg := "could not send remediation.status.changed event"
		r.Logger.Error(msg + ": " + err.Error())
		_ = r.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return err
	}

	actionTriggeredEventData, err := r.getActionTriggeredEventData(problemDetails, action)
	if err != nil {
		msg := "could not create action.triggered event"
		r.Logger.Error(msg + ": " + err.Error())
		_ = r.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return err
	}

	if err := r.sendActionTriggeredEvent(actionTriggeredEventData); err != nil {
		msg := "could not send action.triggered event"
		r.Logger.Error(msg + ": " + err.Error())
		_ = r.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return err
	}
	return nil
}
