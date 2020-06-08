package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const remediationFileName = "remediation.yaml"
const configurationserviceconnection = "CONFIGURATION_SERVICE" //"localhost:6060" // "configuration-service:8080"
const remediationSpecVersion = "0.2.0"

type Handler interface {
	HandleEvent() error
}

func NewHandler(event cloudevents.Event) (Handler, error) {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptn.NewLogger(shkeptncontext, event.Context.GetID(), "remediation-service")
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
			RemediationLog: &RemediationLog{
				Keptn:  keptnHandler,
				Logger: logger,
			},
		}, nil
	case keptn.EvaluationDoneEventType:
		return &EvaluationDoneEventHandler{
			KeptnHandler: keptnHandler,
			Logger:       logger,
			Event:        event,
			RemediationLog: &RemediationLog{
				Keptn:  keptnHandler,
				Logger: logger,
			},
		}, nil
	case keptn.ActionFinishedEventType:
		return &ActionFinishedEventHandler{
			KeptnHandler: keptnHandler,
			Logger:       logger,
			Event:        event,
			RemediationLog: &RemediationLog{
				Keptn:  keptnHandler,
				Logger: logger,
			},
		}, nil
	default:
		return nil, errors.New("no event handler found for type: " + event.Type())
	}
}

type remediation struct {
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

type RemediationLog struct {
	Keptn  *keptn.Keptn
	Logger keptn.LoggerInterface
}

func (rl *RemediationLog) getActionForProblemType(remediationData keptn.RemediationV02, problemType string, index int) *keptn.RemediationV02ActionsOnOpen {
	for _, remediation := range remediationData.Spec.Remediations {
		if strings.HasPrefix(problemType, remediation.ProblemType) {
			rl.Logger.Info("Found remediation for problem type " + remediation.ProblemType)
			if len(remediation.ActionsOnOpen) > index {
				return &remediation.ActionsOnOpen[index]
			}
		}
	}
	return nil
}

func (rl *RemediationLog) sendRemediationTriggeredEvent(problemDetails *keptn.ProblemEventData) error {
	source, _ := url.Parse("remediation-service")
	contentType := "application/json"

	remediationFinishedEventData := &keptn.RemediationTriggeredEventData{
		Project: rl.Keptn.KeptnBase.Project,
		Service: rl.Keptn.KeptnBase.Service,
		Stage:   rl.Keptn.KeptnBase.Stage,
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
		Labels: rl.Keptn.KeptnBase.Labels,
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptn.RemediationTriggeredEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": rl.Keptn.KeptnContext},
		}.AsV02(),
		Data: remediationFinishedEventData,
	}

	err := createRemediation(event.ID(), rl.Keptn.KeptnContext, event.Time().String(), *rl.Keptn.KeptnBase, keptn.RemediationTriggeredEventType, "")
	if err != nil {
		rl.Logger.Error("Could not create remediation: " + err.Error())
		return err
	}
	err = rl.Keptn.SendCloudEvent(event)
	if err != nil {
		rl.Logger.Error("Could not send action.finished event: " + err.Error())
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

	newRemediation := &remediation{
		Action:       action,
		EventID:      eventID,
		KeptnContext: keptnContext,
		Time:         time,
		Type:         remediationEventType,
	}

	queryURL := getRemediationsEndpoint(configurationServiceEndpoint, keptnBase.Project, keptnBase.Stage, keptnBase.Service, keptnContext)
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

func (rl *RemediationLog) sendStartEvaluationEvent() error {
	source, _ := url.Parse("remediation-service")
	contentType := "application/json"

	startEvaluationEventData := &keptn.StartEvaluationEventData{
		Project:      rl.Keptn.KeptnBase.Project,
		Service:      rl.Keptn.KeptnBase.Service,
		Stage:        rl.Keptn.KeptnBase.Stage,
		Labels:       rl.Keptn.KeptnBase.Labels,
		TestStrategy: "real-user",
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptn.StartEvaluationEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": rl.Keptn.KeptnContext},
		}.AsV02(),
		Data: startEvaluationEventData,
	}

	err := rl.Keptn.SendCloudEvent(event)
	if err != nil {
		rl.Logger.Error("Could not send astart-evaluation event: " + err.Error())
		return err
	}
	return nil
}

func (rl *RemediationLog) sendRemediationFinishedEvent(status keptn.RemediationStatusType, result keptn.RemediationResultType, message string) error {
	source, _ := url.Parse("remediation-service")
	contentType := "application/json"

	remediationFinishedEventData := &keptn.RemediationFinishedEventData{
		Project: rl.Keptn.KeptnBase.Project,
		Service: rl.Keptn.KeptnBase.Service,
		Stage:   rl.Keptn.KeptnBase.Stage,
		Problem: keptn.ProblemDetails{},
		Labels:  rl.Keptn.KeptnBase.Labels,
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
			Extensions:  map[string]interface{}{"shkeptncontext": rl.Keptn.KeptnContext},
		}.AsV02(),
		Data: remediationFinishedEventData,
	}

	err := rl.Keptn.SendCloudEvent(event)
	if err != nil {
		rl.Logger.Error("Could not send action.finished event: " + err.Error())
		return err
	}
	return nil
}

func (rl *RemediationLog) getActionTriggeredEventData(problemEvent *keptn.ProblemEventData, action *keptn.RemediationV02ActionsOnOpen) (keptn.ActionTriggeredEventData, error) {
	problemDetails := keptn.ProblemDetails{}
	if err := json.Unmarshal(problemEvent.ProblemDetails, &problemDetails); err != nil {
		rl.Logger.Error("Could not unmarshal ProblemDetails: " + err.Error())
		return keptn.ActionTriggeredEventData{}, err
	}

	return keptn.ActionTriggeredEventData{
		Project: problemEvent.Project,
		Service: problemEvent.Service,
		Stage:   problemEvent.Stage,
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

func (rl *RemediationLog) sendActionTriggeredEvent(ce cloudevents.Event, actionTriggeredEventData keptn.ActionTriggeredEventData) error {

	source, _ := url.Parse("remediation-service")
	contentType := "application/json"

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptn.ActionTriggeredEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": rl.Keptn.KeptnContext},
		}.AsV02(),
		Data: actionTriggeredEventData,
	}

	err := rl.Keptn.SendCloudEvent(event)
	if err != nil {
		rl.Logger.Error("Could not send action.finished event: " + err.Error())
		return err
	}
	return nil
}

func (rl *RemediationLog) sendRemediationStatusChangedEvent(action *keptn.RemediationV02ActionsOnOpen, actionIndex int) error {

	remediationStatusChangedEventData := &keptn.RemediationStatusChangedEventData{
		Project: rl.Keptn.KeptnBase.Project,
		Service: rl.Keptn.KeptnBase.Service,
		Stage:   rl.Keptn.KeptnBase.Stage,
		Labels:  rl.Keptn.KeptnBase.Labels,
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
			Extensions:  map[string]interface{}{"shkeptncontext": rl.Keptn.KeptnContext},
		}.AsV02(),
		Data: remediationStatusChangedEventData,
	}

	err := createRemediation(event.ID(), rl.Keptn.KeptnContext, event.Time().String(), *rl.Keptn.KeptnBase, keptn.RemediationStatusChangedEventType, action.Action)
	if err != nil {
		rl.Logger.Error("Could not create remediation: " + err.Error())
		return err
	}
	err = rl.Keptn.SendCloudEvent(event)
	if err != nil {
		rl.Logger.Error("Could not send remediation.status.changed event: " + err.Error())
		return err
	}
	return nil
}
