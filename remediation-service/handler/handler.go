package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/keptn/go-utils/pkg/lib/v0_1_4"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/ghodss/yaml"
	configmodels "github.com/keptn/go-utils/pkg/api/models"
	configutils "github.com/keptn/go-utils/pkg/api/utils"
	keptn "github.com/keptn/go-utils/pkg/lib"
)

const remediationFileName = "remediation.yaml"
const configurationserviceconnection = "CONFIGURATION_SERVICE" //"localhost:6060" // "configuration-service:8080"
const datastoreConnection = "MONGODB_DATASTORE"
const remediationSpecVersion = "spec.keptn.sh/0.1.4"

// Handler handles incoming Keptn events
type Handler interface {
	HandleEvent() error
}

// NewHandler returns a new Handler for the incoming Keptn event
func NewHandler(event cloudevents.Event) (Handler, error) {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	serviceName := "remediation-service"
	keptnHandler, err := keptnv2.NewKeptn(&event, keptncommon.KeptnOpts{
		LoggingOptions: &keptncommon.LoggingOpts{
			ServiceName: &serviceName,
		},
	})
	if err != nil {
		fmt.Println("Could not initialize Keptn handler: " + err.Error())
		return nil, err
	}

	keptnHandler.Logger.Debug("Received event for shkeptncontext:" + shkeptncontext)

	switch event.Type() {
	case keptn.ProblemOpenEventType:
		return &ProblemOpenEventHandler{
			KeptnHandler: keptnHandler,
			Event:        event,
			Remediation: &Remediation{
				Keptn: keptnHandler,
			},
		}, nil
	case keptn.EvaluationDoneEventType:
		return &EvaluationDoneEventHandler{
			KeptnHandler: keptnHandler,
			Event:        event,
			Remediation: &Remediation{
				Keptn: keptnHandler,
			},
		}, nil
	case keptn.ActionFinishedEventType:
		return &ActionFinishedEventHandler{
			KeptnHandler: keptnHandler,
			Event:        event,
			Remediation: &Remediation{
				Keptn: keptnHandler,
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
	Keptn *keptnv2.Keptn
}

func (r *Remediation) getActionForProblemType(remediationData v0_1_4.Remediation, problemType string, index int) *v0_1_4.RemediationActionsOnOpen {
	for _, remediation := range remediationData.Spec.Remediations {
		if strings.HasPrefix(problemType, remediation.ProblemType) {
			r.Keptn.Logger.Info("Found remediation for problem type " + remediation.ProblemType)
			if len(remediation.ActionsOnOpen) > index {
				return &remediation.ActionsOnOpen[index]
			}
		}
	}
	return nil
}

func (r *Remediation) sendRemediationTriggeredEvent(problemDetails *keptn.ProblemEventData) error {
	source, _ := url.Parse("remediation-service")

	remediationFinishedEventData := &keptn.RemediationTriggeredEventData{
		Project: r.Keptn.KeptnBase.Event.GetProject(),
		Service: r.Keptn.KeptnBase.Event.GetService(),
		Stage:   r.Keptn.KeptnBase.Event.GetStage(),
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
		Labels: r.Keptn.KeptnBase.Event.GetLabels(),
	}

	event := cloudevents.NewEvent()
	event.SetType(keptn.RemediationTriggeredEventType)
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", r.Keptn.KeptnContext)
	event.SetData(cloudevents.ApplicationJSON, remediationFinishedEventData)

	err := createRemediation(event.ID(), r.Keptn.KeptnContext, event.Time().String(), r.Keptn.KeptnBase.Event, keptn.RemediationTriggeredEventType, "")
	if err != nil {
		r.Keptn.Logger.Error("Could not create remediation: " + err.Error())
		return err
	}
	err = r.Keptn.SendCloudEvent(event)
	if err != nil {
		r.Keptn.Logger.Error("Could not send action.finished event: " + err.Error())
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

func createRemediation(eventID, keptnContext, time string, keptnBase keptncommon.EventProperties, remediationEventType, action string) error {
	configurationServiceEndpoint, err := keptncommon.GetServiceEndpoint(configurationserviceconnection)
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

	queryURL := getRemediationsEndpoint(configurationServiceEndpoint, keptnBase.GetProject(), keptnBase.GetStage(), keptnBase.GetService(), "")
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

func deleteRemediation(keptnContext string, keptnBase keptncommon.EventProperties) error {
	configurationServiceEndpoint, err := keptncommon.GetServiceEndpoint(configurationserviceconnection)
	if err != nil {
		return errors.New("could not retrieve configuration-service URL")
	}

	queryURL := getRemediationsEndpoint(configurationServiceEndpoint, keptnBase.GetProject(), keptnBase.GetStage(), keptnBase.GetService(), keptnContext)
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

func getRemediationsByContext(keptnContext string, keptnBase keptncommon.EventProperties) ([]*remediationStatus, error) {
	configurationServiceEndpoint, err := keptncommon.GetServiceEndpoint(configurationserviceconnection)
	if err != nil {
		return nil, errors.New("could not retrieve configuration-service URL")
	}

	remediations := []*remediationStatus{}

	nextPageKey := ""

	for {
		queryURL := getRemediationsEndpoint(configurationServiceEndpoint, keptnBase.GetProject(), keptnBase.GetStage(), keptnBase.GetService(), keptnContext)

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

	waitTime := getWaitTime()
	startEvaluationEventData := &keptn.StartEvaluationEventData{
		Project:      r.Keptn.Event.GetProject(),
		Service:      r.Keptn.Event.GetService(),
		Stage:        r.Keptn.Event.GetStage(),
		Labels:       r.Keptn.Event.GetLabels(),
		Start:        time.Now().Add(-waitTime).Format(time.RFC3339),
		End:          time.Now().Format(time.RFC3339),
		TestStrategy: "real-user",
	}

	event := cloudevents.NewEvent()
	event.SetType(keptn.StartEvaluationEventType)
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", r.Keptn.KeptnContext)
	event.SetData(cloudevents.ApplicationJSON, startEvaluationEventData)

	err := r.Keptn.SendCloudEvent(event)
	if err != nil {
		r.Keptn.Logger.Error("Could not send astart-evaluation event: " + err.Error())
		return err
	}
	return nil
}

func (r *Remediation) sendRemediationFinishedEvent(status keptn.RemediationStatusType, result keptn.RemediationResultType, message string) error {
	source, _ := url.Parse("remediation-service")

	remediationFinishedEventData := &keptn.RemediationFinishedEventData{
		Project: r.Keptn.Event.GetProject(),
		Service: r.Keptn.Event.GetService(),
		Stage:   r.Keptn.Event.GetStage(),
		Problem: keptn.ProblemDetails{},
		Labels:  r.Keptn.Event.GetLabels(),
		Remediation: keptn.RemediationFinished{
			Status:  status,
			Result:  result,
			Message: message,
		},
	}

	event := cloudevents.NewEvent()
	event.SetType(keptn.RemediationFinishedEventType)
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", r.Keptn.KeptnContext)
	event.SetData(cloudevents.ApplicationJSON, remediationFinishedEventData)

	err := deleteRemediation(r.Keptn.KeptnContext, r.Keptn.Event)
	if err != nil {
		r.Keptn.Logger.Error("Could not close remediation: " + err.Error())
	}

	err = r.Keptn.SendCloudEvent(event)
	if err != nil {
		r.Keptn.Logger.Error("Could not send action.finished event: " + err.Error())
		return err
	}
	return nil
}

func (r *Remediation) getActionTriggeredEventData(problemDetails keptn.ProblemDetails, action *v0_1_4.RemediationActionsOnOpen) (keptn.ActionTriggeredEventData, error) {
	return keptn.ActionTriggeredEventData{
		Project: r.Keptn.Event.GetProject(),
		Service: r.Keptn.Event.GetService(),
		Stage:   r.Keptn.Event.GetStage(),
		Action: keptn.ActionInfo{
			Name:        action.Name,
			Action:      action.Action,
			Description: action.Description,
			Value:       action.Value,
		},
		Problem: problemDetails,
		Labels:  r.Keptn.Event.GetLabels(),
	}, nil
}

func (r *Remediation) sendActionTriggeredEvent(actionTriggeredEventData keptn.ActionTriggeredEventData) error {

	source, _ := url.Parse("remediation-service")

	event := cloudevents.NewEvent()
	event.SetType(keptn.ActionTriggeredEventType)
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", r.Keptn.KeptnContext)
	event.SetData(cloudevents.ApplicationJSON, actionTriggeredEventData)

	err := r.Keptn.SendCloudEvent(event)
	if err != nil {
		r.Keptn.Logger.Error("Could not send action.finished event: " + err.Error())
		return err
	}
	return nil
}

func (r *Remediation) sendRemediationStatusChangedEvent(action *v0_1_4.RemediationActionsOnOpen, actionIndex int) error {

	remediationStatusChangedEventData := &keptn.RemediationStatusChangedEventData{
		Project: r.Keptn.Event.GetProject(),
		Service: r.Keptn.Event.GetService(),
		Stage:   r.Keptn.Event.GetStage(),
		Labels:  r.Keptn.Event.GetLabels(),
		Remediation: keptn.RemediationStatusChanged{
			Status: keptn.RemediationStatusSucceeded,
			Result: keptn.RemediationResult{
				ActionIndex: actionIndex,
				ActionName:  action.Action,
			},
		},
	}

	source, _ := url.Parse("remediation-service")

	event := cloudevents.NewEvent()
	event.SetType(keptn.RemediationStatusChangedEventType)
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", r.Keptn.KeptnContext)
	event.SetData(cloudevents.ApplicationJSON, remediationStatusChangedEventData)

	err := createRemediation(event.ID(), r.Keptn.KeptnContext, event.Time().String(), r.Keptn.Event, keptn.RemediationStatusChangedEventType, action.Action)
	if err != nil {
		r.Keptn.Logger.Error("Could not create remediation: " + err.Error())
		return err
	}
	err = r.Keptn.SendCloudEvent(event)
	if err != nil {
		r.Keptn.Logger.Error("Could not send remediation.status.changed event: " + err.Error())
		return err
	}
	return nil
}

func (r *Remediation) getRemediationFile() (*configmodels.Resource, error) {
	resourceHandler := configutils.NewResourceHandler(os.Getenv(configurationserviceconnection))
	var resource *configmodels.Resource
	var err error
	if r.Keptn.Event.GetService() != "" {
		resource, err = resourceHandler.GetServiceResource(r.Keptn.Event.GetProject(), r.Keptn.Event.GetStage(),
			r.Keptn.Event.GetService(), remediationFileName)
	} else {
		resource, err = resourceHandler.GetStageResource(r.Keptn.Event.GetProject(), r.Keptn.Event.GetStage(), remediationFileName)
	}

	if err != nil {
		var msg string
		if strings.Contains(strings.ToLower(err.Error()), "service not found") {
			msg = "Could not execute remediation action because service is not available"
		} else {
			msg = "Could not execute remediation action because no remediation file available"
		}
		r.Keptn.Logger.Error(msg)
		_ = r.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return nil, err
	}
	r.Keptn.Logger.Debug("remediation.yaml for service found")
	return resource, nil
}

func (r *Remediation) getRemediation(resource *configmodels.Resource) (*v0_1_4.Remediation, error) {
	remediationData := &v0_1_4.Remediation{}
	err := yaml.Unmarshal([]byte(resource.ResourceContent), remediationData)
	if err != nil {
		msg := "could not parse remediation.yaml"
		r.Keptn.Logger.Error(msg + ": " + err.Error())
		_ = r.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return nil, err
	}

	if remediationData.ApiVersion != remediationSpecVersion {
		msg := "remediation.yaml file does not conform to remediation spec " + remediationSpecVersion
		r.Keptn.Logger.Error(msg)
		_ = r.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return nil, errors.New(msg)
	}
	return remediationData, nil
}

func (r *Remediation) triggerAction(action *v0_1_4.RemediationActionsOnOpen, actionIndex int, problemDetails keptn.ProblemDetails) error {
	err := r.sendRemediationStatusChangedEvent(action, actionIndex)
	if err != nil {
		msg := "could not send remediation.status.changed event"
		r.Keptn.Logger.Error(msg + ": " + err.Error())
		_ = r.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return err
	}

	actionTriggeredEventData, err := r.getActionTriggeredEventData(problemDetails, action)
	if err != nil {
		msg := "could not create action.triggered event"
		r.Keptn.Logger.Error(msg + ": " + err.Error())
		_ = r.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return err
	}

	if err := r.sendActionTriggeredEvent(actionTriggeredEventData); err != nil {
		msg := "could not send action.triggered event"
		r.Keptn.Logger.Error(msg + ": " + err.Error())
		_ = r.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return err
	}
	return nil
}
