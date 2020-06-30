package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
)

const datastore = "MONGODB_DATASTORE"

type approval struct {
	EventID      string `json:"eventId"`
	Image        string `json:"image"`
	KeptnContext string `json:"keptnContext"`
	Tag          string `json:"tag"`
	Time         string `json:"time"`
}

type ApprovalFinishedEventHandler struct {
	keptn *keptnevents.Keptn
}

// NewApprovalFinishedEventHandler returns a new approval.finished event handler
func NewApprovalFinishedEventHandler(keptn *keptnevents.Keptn) *ApprovalFinishedEventHandler {
	return &ApprovalFinishedEventHandler{keptn: keptn}
}

func (a *ApprovalFinishedEventHandler) IsTypeHandled(event cloudevents.Event) bool {
	return event.Type() == keptnevents.ApprovalFinishedEventType
}

func (a *ApprovalFinishedEventHandler) Handle(event cloudevents.Event, keptnHandler *keptnevents.Keptn, shipyard *keptnevents.Shipyard) {
	data := &keptnevents.ApprovalFinishedEventData{}
	if err := event.DataAs(data); err != nil {
		a.keptn.Logger.Error(fmt.Sprintf("failed to parse ApprovalTriggeredEventData: %v", err))
		return
	}

	var triggeredID string
	if err := event.Context.ExtensionAs("triggeredid", &triggeredID); err != nil {
		a.keptn.Logger.Error(fmt.Sprintf("triggeredid is missing: %v", err))
		return
	}
	outgoingEvents := a.handleApprovalFinishedEvent(*data, keptnHandler.KeptnContext, triggeredID, *shipyard)
	sendEvents(keptnHandler, outgoingEvents, a.keptn.Logger)
}

func (a *ApprovalFinishedEventHandler) handleApprovalFinishedEvent(inputEvent keptnevents.ApprovalFinishedEventData, shkeptncontext string,
	triggeredID string, shipyard keptnevents.Shipyard) []cloudevents.Event {

	outgoingEvents := make([]cloudevents.Event, 0)
	if inputEvent.Approval.Status != SucceededResult {
		a.keptn.Logger.Info(fmt.Sprintf("Approval finished with failed status for "+
			"image %s for service %s of project %s and current stage %s received",
			inputEvent.Image, inputEvent.Service, inputEvent.Project, inputEvent.Stage))
	} else {
		if inputEvent.Approval.Result == PassResult {
			a.keptn.Logger.Info(fmt.Sprintf("Approval for image %s for service %s of project %s and current stage %s received",
				inputEvent.Image, inputEvent.Service, inputEvent.Project, inputEvent.Stage))

			openApproval, err := getOpenApproval(inputEvent, triggeredID)
			if err != nil {
				a.keptn.Logger.Error("Could not retrieve open Approval with EventID " + triggeredID + ": " + err.Error())
				return outgoingEvents
			}
			if openApproval.Image != inputEvent.Image {
				a.keptn.Logger.Error(fmt.Sprintf("Image of approval-finished event %s does not match with image of open approval: %s != %s\n", openApproval.EventID, openApproval.Image, inputEvent.Image))
				return outgoingEvents
			}
			if openApproval.Tag != inputEvent.Tag {
				a.keptn.Logger.Error(fmt.Sprintf("Tag of approval-finished event %s does not match with image of open approval: %s != %s\n", openApproval.EventID, openApproval.Image, inputEvent.Image))
				return outgoingEvents
			}
			image := inputEvent.Image
			if inputEvent.Tag != "" {
				image += ":" + inputEvent.Tag
			}
			if event := getConfigurationChangeEventForCanary(
				inputEvent.Project, inputEvent.Service, inputEvent.Stage, image, shkeptncontext, inputEvent.Labels); event != nil {
				outgoingEvents = append(outgoingEvents, *event)
			}
		} else {
			a.keptn.Logger.Info(fmt.Sprintf("Rejection for image %s for service %s of project %s and current stage %s received",
				inputEvent.Image, inputEvent.Service, inputEvent.Project, inputEvent.Stage))
		}
		if err := closeOpenApproval(inputEvent, triggeredID); err != nil {
			a.keptn.Logger.Error(fmt.Sprintf("failed to close open approvals in materialized view: %v", err))
			return outgoingEvents
		}
	}

	return outgoingEvents
}

func getOpenApproval(inputEvent keptnevents.ApprovalFinishedEventData, triggeredID string) (*approval, error) {
	configurationServiceEndpoint, err := keptnevents.GetServiceEndpoint(configService)
	if err != nil {
		return nil, errors.New("could not retrieve configuration-service URL")
	}

	queryURL := getApprovalsEndpoint(configurationServiceEndpoint, inputEvent.Project, inputEvent.Stage, inputEvent.Service, triggeredID)
	client := &http.Client{}
	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		//
		<-time.After(5 * time.Second)
		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	approval := &approval{}
	err = json.Unmarshal(body, approval)
	if err != nil {
		return nil, err
	}

	return approval, nil
}

func getApprovalsEndpoint(configurationServiceEndpoint url.URL, project, stage, service, approvalTriggeredID string) string {
	if approvalTriggeredID == "" {
		return fmt.Sprintf("%s://%s/v1/project/%s/stage/%s/service/%s/approval", configurationServiceEndpoint.Scheme, configurationServiceEndpoint.Host, project, stage, service)
	}
	return fmt.Sprintf("%s://%s/v1/project/%s/stage/%s/service/%s/approval/%s", configurationServiceEndpoint.Scheme, configurationServiceEndpoint.Host, project, stage, service, approvalTriggeredID)
}

func closeOpenApproval(inputEvent keptnevents.ApprovalFinishedEventData, triggeredID string) error {
	configurationServiceEndpoint, err := keptnevents.GetServiceEndpoint(configService)
	if err != nil {
		return errors.New("could not retrieve configuration-service URL")
	}

	queryURL := getApprovalsEndpoint(configurationServiceEndpoint, inputEvent.Project, inputEvent.Stage, inputEvent.Service, triggeredID)
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

	if resp.StatusCode == http.StatusNotFound {
		<-time.After(5 * time.Second)
		resp, err = client.Do(req)
		if err != nil {
			return err
		}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	return nil
}
