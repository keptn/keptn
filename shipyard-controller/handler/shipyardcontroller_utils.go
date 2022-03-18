package handler

import (
	"encoding/json"
	"fmt"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
)

func GetTaskSequenceInStage(stageName, taskSequenceName string, shipyard *keptnv2.Shipyard) (*keptnv2.Sequence, error) {
	stage := GetStageFromShipyard(stageName, shipyard)
	if stage == nil {
		return nil, fmt.Errorf("no stage with name %s", stageName)
	}

	for _, taskSequence := range stage.Sequences {
		if taskSequence.Name == taskSequenceName {
			log.Infof("Found matching task sequence %s in stage %s", taskSequence.Name, stage.Name)
			if len(taskSequence.Tasks) == 0 {
				return nil, fmt.Errorf("task sequence %s does not contain any tasks", taskSequenceName)
			}
			return &taskSequence, nil
		}
	}
	// provide built-int task sequence for evaluation
	if taskSequenceName == keptnv2.EvaluationTaskName {
		return &keptnv2.Sequence{
			Name:        "evaluation",
			TriggeredOn: nil,
			Tasks: []keptnv2.Task{
				{
					Name: keptnv2.EvaluationTaskName,
				},
			},
		}, nil
	}
	return nil, fmt.Errorf("no task sequence with name %s found in stage %s", taskSequenceName, stageName)

}

func GetStageFromShipyard(stageName string, shipyard *keptnv2.Shipyard) *keptnv2.Stage {
	for _, stage := range shipyard.Spec.Stages {
		if stage.Name == stageName {
			return &stage
		}
	}
	return nil
}

func GetTaskSequencesByTrigger(eventScope models.EventScope, completedTaskSequence string, shipyard *keptnv2.Shipyard, previousTask string) []NextTaskSequence {
	var result []NextTaskSequence

	for _, stage := range shipyard.Spec.Stages {
		for tsIndex, taskSequence := range stage.Sequences {
			for _, trigger := range taskSequence.TriggeredOn {
				if trigger.Event == eventScope.Stage+"."+completedTaskSequence+".finished" {
					appendSequence := false
					// default behavior if no selector is available: 'pass', as well as 'warning' results trigger this sequence
					if trigger.Selector.Match == nil {
						if eventScope.Result == keptnv2.ResultPass || eventScope.Result == keptnv2.ResultWarning {
							appendSequence = true
						}
					} else {
						// if a selector is there, compare the 'result' property
						if string(eventScope.Result) == trigger.Selector.Match["result"] {
							appendSequence = true
						} else if string(eventScope.Result) == trigger.Selector.Match[previousTask+".result"] {
							appendSequence = true
						}
					}
					if appendSequence {
						result = append(result, NextTaskSequence{
							Sequence:  stage.Sequences[tsIndex],
							StageName: stage.Name,
						})
					}
				}
			}
		}
	}
	return result
}

func ObjToJSON(obj interface{}) string {
	indent, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return ""
	}
	return string(indent)
}

func ExtractEventKind(event apimodels.KeptnContextExtendedCE) (string, error) {
	eventData := &keptnv2.EventData{}
	err := keptnv2.Decode(event.Data, eventData)
	if err != nil {
		log.Errorf("Could not parse event data: %v", err)
		return "", err
	}

	statusType, err := keptnv2.ParseEventKind(*event.Type)
	if err != nil {
		return "", err
	}
	return statusType, nil
}
