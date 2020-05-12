package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

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
	logger *keptnevents.Logger
}

func NewApprovalFinishedEventHandler(l *keptnevents.Logger) *ApprovalFinishedEventHandler {
	return &ApprovalFinishedEventHandler{logger: l}
}

func (a *ApprovalFinishedEventHandler) IsTypeHandled(event cloudevents.Event) bool {
	return event.Type() == keptnevents.ApprovalFinishedEventType
}

func (a *ApprovalFinishedEventHandler) Handle(event cloudevents.Event, keptnHandler *keptnevents.Keptn, shipyard *keptnevents.Shipyard) {
	data := &keptnevents.ApprovalFinishedEventData{}
	if err := event.DataAs(data); err != nil {
		a.logger.Error(fmt.Sprintf("failed to parse ApprovalTriggeredEventData: %v", err))
	}

	outgoingEvents := a.handleApprovalFinishedEvent(*data, keptnHandler.KeptnContext, *shipyard)
	sendEvents(keptnHandler, outgoingEvents, a.logger)
}

func (a *ApprovalFinishedEventHandler) handleApprovalFinishedEvent(inputEvent keptnevents.ApprovalFinishedEventData, shkeptncontext string,
	shipyard keptnevents.Shipyard) []cloudevents.Event {

	outgoingEvents := make([]cloudevents.Event, 0)
	if inputEvent.Approval.Status != SucceededResult {
		a.logger.Info(fmt.Sprintf("Approval finished with failed status for "+
			"image %s for service %s of project %s and current stage %s received",
			inputEvent.Image, inputEvent.Service, inputEvent.Project, inputEvent.Stage))
	} else {
		if inputEvent.Approval.Result == PassResult {
			a.logger.Info(fmt.Sprintf("Approval for image %s for service %s of project %s and current stage %s received",
				inputEvent.Image, inputEvent.Service, inputEvent.Project, inputEvent.Stage))

			// TODO: Check image using inputEvent.Approval.TriggeredID
			openApproval, err := getOpenApproval(inputEvent)
			if err != nil {
				a.logger.Error("Could not retrieve open Approval with EventID " + inputEvent.Approval.TriggeredID + ": " + err.Error())
				return outgoingEvents
			}
			if openApproval.Image != inputEvent.Image {
				a.logger.Error(fmt.Sprintf("Image of approval-finished event %s does not match with image of open approval: %s != %s\n", openApproval.EventID, openApproval.Image, inputEvent.Image))
				return outgoingEvents
			}
			if openApproval.Tag != inputEvent.Tag {
				a.logger.Error(fmt.Sprintf("Tag of approval-finished event %s does not match with image of open approval: %s != %s\n", openApproval.EventID, openApproval.Image, inputEvent.Image))
				return outgoingEvents
			}
			image := inputEvent.Image
			if inputEvent.Tag != "" {
				image += ":" + inputEvent.Tag
			}
			if event := getPromotionEvent(inputEvent.Project, inputEvent.Stage, inputEvent.Service, image,
				shkeptncontext, inputEvent.Labels, shipyard, a.logger); event != nil {
				outgoingEvents = append(outgoingEvents, *event)
			}
		} else {
			a.logger.Info(fmt.Sprintf("Rejection for image %s for service %s of project %s and current stage %s received",
				inputEvent.Image, inputEvent.Service, inputEvent.Project, inputEvent.Stage))
		}
	}

	return outgoingEvents
}

func getOpenApproval(inputEvent keptnevents.ApprovalFinishedEventData) (*approval, error) {
	configurationServiceEndpoint, err := keptnevents.GetServiceEndpoint(configService)
	if err != nil {
		return nil, errors.New("could not retrieve configuration-service URL")
	}

	queryURL := getApprovalsEndpoint(configurationServiceEndpoint, inputEvent)
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

func getApprovalsEndpoint(configurationServiceEndpoint url.URL, inputEvent keptnevents.ApprovalFinishedEventData) string {
	return fmt.Sprintf("%s://%s/v1/project/%s/stage/%s/service/%s/approval/%s", configurationServiceEndpoint.Scheme, configurationServiceEndpoint.Host, inputEvent.Project, inputEvent.Stage, inputEvent.Service, inputEvent.Approval.TriggeredID)
}

func closeOpenApproval(inputEvent keptnevents.ApprovalFinishedEventData) error {
	configurationServiceEndpoint, err := keptnevents.GetServiceEndpoint(configService)
	if err != nil {
		return errors.New("could not retrieve configuration-service URL")
	}

	queryURL := getApprovalsEndpoint(configurationServiceEndpoint, inputEvent)
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

	if resp.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	return nil
}
