package handler

import (
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/ghodss/yaml"
	configmodels "github.com/keptn/go-utils/pkg/api/models"
	configutils "github.com/keptn/go-utils/pkg/api/utils"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"os"
	"strings"
)

type ProblemOpenEventHandler struct {
	KeptnHandler   *keptn.Keptn
	Logger         keptn.LoggerInterface
	Event          cloudevents.Event
	RemediationLog *RemediationLog
}

func (eh *ProblemOpenEventHandler) HandleEvent() error {
	var problemEvent *keptn.ProblemEventData
	if eh.Event.Type() == keptn.ProblemOpenEventType {
		eh.Logger.Debug("Received problem notification")
		problemEvent = &keptn.ProblemEventData{}
		if err := eh.Event.DataAs(problemEvent); err != nil {
			return err
		}
	}

	if !isProjectAndStageAvailable(problemEvent) {
		deriveFromTags(problemEvent)
	}
	if !isProjectAndStageAvailable(problemEvent) {
		return errors.New("Cannot derive project and stage from tags nor impacted entity")
	}

	eh.Logger.Debug("Received problem event with state " + problemEvent.State)

	// check if remediation should be performed
	resourceHandler := configutils.NewResourceHandler(os.Getenv(configurationserviceconnection))
	autoRemediate, err := eh.isRemediationEnabled()
	if err != nil {
		eh.Logger.Error(fmt.Sprintf("Failed to check if remediation is enabled: %s", err.Error()))
		return err
	}

	if autoRemediate {
		eh.Logger.Info(fmt.Sprintf("Remediation enabled for project %s in stage %s", problemEvent.Project, problemEvent.Stage))
	} else {
		eh.Logger.Info(fmt.Sprintf("Remediation disabled for project %s in stage %s", problemEvent.Project, problemEvent.Stage))
		return nil
	}

	// get remediation.yaml
	var resource *configmodels.Resource
	if problemEvent.Service != "" {
		resource, err = resourceHandler.GetServiceResource(problemEvent.Project, problemEvent.Stage,
			problemEvent.Service, remediationFileName)
	} else {
		resource, err = resourceHandler.GetStageResource(problemEvent.Project, problemEvent.Stage, remediationFileName)
	}

	if err != nil {
		msg := "remediation file not configured"
		eh.Logger.Error(msg)
		_ = eh.RemediationLog.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return err
	}
	eh.Logger.Debug("remediation.yaml for service found")

	// get remediation action from remediation.yaml
	var remediationData keptn.RemediationV02
	err = yaml.Unmarshal([]byte(resource.ResourceContent), &remediationData)
	if err != nil {
		msg := "could not parse remediation.yaml"
		eh.Logger.Error(msg + ": " + err.Error())
		_ = eh.RemediationLog.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return err
	}

	if remediationData.Version != remediationSpecVersion {
		msg := "remediation.yaml file does not conform to remediation spec v0.2.0"
		eh.Logger.Error(msg)
		_ = eh.RemediationLog.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return err
	}

	err = eh.RemediationLog.sendRemediationTriggeredEvent(problemEvent)
	if err != nil {
		msg := "could not send remediation.triggered event"
		eh.Logger.Error(msg + ": " + err.Error())
		_ = eh.RemediationLog.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return err
	}

	problemType := problemEvent.ProblemTitle

	action := eh.RemediationLog.getActionForProblemType(remediationData, problemType, 0)
	if action == nil {
		action = eh.RemediationLog.getActionForProblemType(remediationData, "*", 0)
	}

	if action != nil {
		err := eh.RemediationLog.sendRemediationStatusChangedEvent(action, 0)
		if err != nil {
			msg := "could not send remediation.status.changed event"
			eh.Logger.Error(msg + ": " + err.Error())
			_ = eh.RemediationLog.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
			return err
		}

		actionTriggeredEventData, err := eh.RemediationLog.getActionTriggeredEventData(problemEvent, action)
		if err != nil {
			msg := "could not create action.triggered event"
			eh.Logger.Error(msg + ": " + err.Error())
			_ = eh.RemediationLog.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
			return err
		}

		if err := eh.RemediationLog.sendActionTriggeredEvent(eh.Event, actionTriggeredEventData); err != nil {
			msg := "could not send action.triggered event"
			eh.Logger.Error(msg + ": " + err.Error())
			_ = eh.RemediationLog.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
			return err
		}
	} else {
		msg := "No remediation configured for problem type " + problemType
		eh.Logger.Info(msg)
		_ = eh.RemediationLog.sendRemediationFinishedEvent(keptn.RemediationStatusSucceeded, keptn.RemediationResultPass, "triggered all actions")
		err = deleteRemediation(eh.KeptnHandler.KeptnContext, *eh.KeptnHandler.KeptnBase)
		if err != nil {
			eh.Logger.Error("Could not close remediation: " + err.Error())
			return err
		}
	}

	return nil
}

func (eh *ProblemOpenEventHandler) isRemediationEnabled() (bool, error) {
	shipyard, err := eh.KeptnHandler.GetShipyard()
	if err != nil {
		return false, err
	}
	for _, s := range shipyard.Stages {
		if s.Name == eh.KeptnHandler.KeptnBase.Stage && s.RemediationStrategy == "automated" {
			return true, nil
		}
	}

	return false, nil
}

func isProjectAndStageAvailable(problem *keptn.ProblemEventData) bool {
	return problem.Project != "" && problem.Stage != ""
}

// deriveFromTags allows to derive project, stage, and service information from tags
// Input example: "Tags:":"keptn_service:carts, keptn_stage:dev, keptn_stage:sockshop"
func deriveFromTags(problem *keptn.ProblemEventData) {

	tags := strings.Split(problem.Tags, ", ")

	for _, tag := range tags {
		if strings.HasPrefix(tag, "keptn_service:") {
			problem.Service = tag[len("keptn_service:"):]
		} else if strings.HasPrefix(tag, "keptn_stage:") {
			problem.Stage = tag[len("keptn_stage:"):]
		} else if strings.HasPrefix(tag, "keptn_project:") {
			problem.Project = tag[len("keptn_project:"):]
		}
	}
}
