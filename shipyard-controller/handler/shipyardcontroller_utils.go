package handler

import (
	"encoding/json"
	"fmt"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
)

func GetTaskSequenceInStage(stageName, taskSequenceName string, shipyard *keptnv2.Shipyard) (*keptnv2.Sequence, error) {
	for _, stage := range shipyard.Spec.Stages {
		if stage.Name == stageName {
			for _, taskSequence := range stage.Sequences {
				if taskSequence.Name == taskSequenceName {
					log.Infof("Found matching task sequence %s in stage %s", taskSequence.Name, stage.Name)
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
	}
	return nil, fmt.Errorf("no stage with name %s", stageName)
}

func GetNextTaskOfSequence(taskSequence *keptnv2.Sequence, previousTask *models.TaskSequenceEvent, eventScope *models.EventScope, eventHistory []interface{}) *models.Task {
	if previousTask != nil {
		for _, e := range eventHistory {
			eventData := keptnv2.EventData{}
			_ = keptnv2.Decode(e, &eventData)

			// if one of the tasks has failed previously, no further task should be executed
			if eventData.Status == keptnv2.StatusErrored || eventData.Result == keptnv2.ResultFailed {
				eventScope.Status = eventData.Status
				eventScope.Result = eventData.Result
				return nil
			}
		}
	}

	if len(taskSequence.Tasks) == 0 {
		log.Infof("Task sequence %s does not contain any tasks", taskSequence.Name)
		return nil
	}
	if previousTask == nil {
		log.Infof("Returning first task of task sequence %s", taskSequence.Name)
		return &models.Task{
			Task:      taskSequence.Tasks[0],
			TaskIndex: 0,
		}
	}

	log.Infof("Getting task that should be executed after task %s", previousTask.Task.Name)

	nextIndex := previousTask.Task.TaskIndex + 1
	if len(taskSequence.Tasks) > nextIndex && taskSequence.Tasks[nextIndex-1].Name == previousTask.Task.Name {
		log.Infof("found next task: %s", taskSequence.Tasks[nextIndex].Name)
		return &models.Task{
			Task:      taskSequence.Tasks[nextIndex],
			TaskIndex: nextIndex,
		}
	}

	log.Info("No further tasks detected")
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

func GetMergedPayloadForSequenceTriggeredEvent(inputEvent *models.Event, eventPayload map[string]interface{}, eventHistory []interface{}) (interface{}, error) {
	var mergedPayload interface{}
	if inputEvent != nil {
		marshal, err := json.Marshal(inputEvent.Data)
		if err != nil {
			return nil, fmt.Errorf("could not marshal input event: %s ", err.Error())
		}
		tmp := map[string]interface{}{}
		if err := json.Unmarshal(marshal, &tmp); err != nil {
			return nil, fmt.Errorf("could not convert input event: %s ", err.Error())
		}
		mergedPayload = common.Merge(eventPayload, tmp)
	}
	if eventHistory != nil {
		for index := range eventHistory {
			if mergedPayload == nil {
				mergedPayload = common.Merge(eventPayload, eventHistory[index])
			} else {
				mergedPayload = common.Merge(mergedPayload, eventHistory[index])
			}
		}
	}
	return mergedPayload, nil
}

func ObjToJSON(obj interface{}) string {
	indent, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return ""
	}
	return string(indent)
}
