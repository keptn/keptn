package common

import (
	"fmt"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/mongodb-datastore/models"
)

func TransformEvaluationDoneEvent(keptnEvent *models.KeptnContextExtendedCE) error {
	eventMap := map[string]interface{}{}
	convertedEvent := &keptnv2.EvaluationFinishedEventData{}
	if err := keptnv2.Decode(keptnEvent.Data, &eventMap); err != nil {
		return fmt.Errorf("failed to transform evaluation-done event to evaluation.finished event %v", err)
	}
	if err := keptnv2.Decode(keptnEvent.Data, convertedEvent); err != nil {
		return fmt.Errorf("failed to transform evaluation-done event to evaluation.finished event %v", err)
	}
	if eventMap["evaluationdetails"] != nil {
		evaluationDetails := &keptnv2.EvaluationDetails{}
		if err := keptnv2.Decode(eventMap["evaluationdetails"], evaluationDetails); err != nil {
			return fmt.Errorf("failed to transform evaluationDetails of evaluation-done event: %v", err)
		}
		convertedEvent.Evaluation = *evaluationDetails
	}
	keptnEvent.Type = models.Type(keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName))
	keptnEvent.Data = convertedEvent
	keptnEvent.Specversion = "1.0"
	return nil
}
