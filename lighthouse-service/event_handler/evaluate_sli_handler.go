package event_handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/ghodss/yaml"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptn "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type datastoreResult struct {
	NextPageKey string `json:"nextPageKey"`
	TotalCount  int    `json:"totalCount"`
	PageSize    int    `json:"pageSize"`
	Events      []struct {
		Data interface{} `json:"data"`
		ID   string      `json:"id"`
	}
}

type criteriaObject struct {
	Operator        string
	Value           float64
	CheckPercentage bool
	IsComparison    bool
	CheckIncrease   bool
}

type EvaluateSLIHandler struct {
	Event        cloudevents.Event
	HTTPClient   *http.Client
	KeptnHandler *keptnv2.Keptn
}

func (eh *EvaluateSLIHandler) HandleEvent() error {
	e := &keptnv2.GetSLIFinishedEventData{}

	var shkeptncontext string
	eh.Event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)
	err := eh.Event.DataAs(&e)

	triggeredEvents, err2 := eh.KeptnHandler.EventHandler.GetEvents(&keptnapi.EventFilter{
		Project:      e.Project,
		Stage:        e.Stage,
		Service:      e.Service,
		EventType:    keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName),
		KeptnContext: eh.KeptnHandler.KeptnContext,
	})
	if err2 != nil {
		msg := fmt.Sprintf("Could not retrieve evaluation.triggered event for context %s %v", eh.KeptnHandler.KeptnContext, err2)
		eh.KeptnHandler.Logger.Error(msg)
		return sendErroredFinishedEventWithMessage(shkeptncontext, "", msg, "", eh.KeptnHandler, e)
	}
	if triggeredEvents == nil || len(triggeredEvents) == 0 {
		msg := "Could not retrieve evaluation.triggered event for context " + eh.KeptnHandler.KeptnContext
		eh.KeptnHandler.Logger.Error(msg)
		return sendErroredFinishedEventWithMessage(shkeptncontext, "", msg, "", eh.KeptnHandler, e)
	}
	triggeredID := triggeredEvents[0].ID

	if err != nil {
		msg := "Could not parse event payload: " + err.Error()
		eh.KeptnHandler.Logger.Error(msg)
		return sendErroredFinishedEventWithMessage(shkeptncontext, "", msg, "", eh.KeptnHandler, e)
	}

	eh.KeptnHandler.Logger.Debug("Start to evaluate SLIs")
	// compare the results based on the evaluation strategy
	sloConfig, err := getSLOs(e.Project, e.Stage, e.Service)
	if err != nil {
		if err == ErrSLOFileNotFound {
			evaluationDetails := keptnv2.EvaluationDetails{
				IndicatorResults: nil,
				TimeStart:        e.GetSLI.Start,
				TimeEnd:          e.GetSLI.End,
				Result:           fmt.Sprintf("no evaluation performed by lighthouse because no SLO file configured for project %s", e.Project),
			}

			evaluationResult := keptnv2.EvaluationFinishedEventData{
				Evaluation: evaluationDetails,
				EventData: keptnv2.EventData{
					Result:  "pass",
					Project: e.Project,
					Service: e.Service,
					Stage:   e.Stage,
					Labels:  e.Labels,
				},
			}
			return sendEvent(shkeptncontext, triggeredID, keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName), eh.KeptnHandler, &evaluationResult)
		}
		return sendErroredFinishedEventWithMessage(shkeptncontext, triggeredID, err.Error(), "", eh.KeptnHandler, e)
	}

	var sloFileContent []byte
	// get the slo.yaml as a plain file to avoid confusion due to defaulted values (see https://github.com/keptn/keptn/issues/1495)
	sloFileContentTmp, err := eh.KeptnHandler.GetKeptnResource("slo.yaml")
	if err != nil {
		eh.KeptnHandler.Logger.Debug("Could not fetch slo.yaml from service repository: " + err.Error() + ". Will append internally used SLO object to evaluation.finished event.")
		sloFileContent, _ = yaml.Marshal(sloConfig)
	} else {
		sloFileContent = []byte(sloFileContentTmp)
	}

	// get results of previous evaluations from data store (mongodb-datastore)
	numberOfPreviousResults := 3
	if sloConfig.Comparison.CompareWith == "single_result" {
		numberOfPreviousResults = 1
	} else if sloConfig.Comparison.CompareWith == "several_results" {
		numberOfPreviousResults = sloConfig.Comparison.NumberOfComparisonResults
	}

	previousEvaluationEvents, comparisonEventIDs, err := eh.getPreviousEvaluations(e, numberOfPreviousResults, sloConfig.Comparison.IncludeResultWithScore)
	if err != nil {
		return sendErroredFinishedEventWithMessage(shkeptncontext, triggeredID, err.Error(), string(sloFileContent), eh.KeptnHandler, e)
	}

	var filteredPreviousEvaluationEvents []*keptnv2.EvaluationFinishedEventData

	// verify that we have enough evaluations
	for _, val := range previousEvaluationEvents {
		filteredPreviousEvaluationEvents = append(filteredPreviousEvaluationEvents, val)
	}

	evaluationResult, maximumAchievableScore, keySLIFailed := evaluateObjectives(e, sloConfig, filteredPreviousEvaluationEvents)
	evaluationResult.Labels = e.Labels
	evaluationResult.Evaluation.ComparedEvents = comparisonEventIDs

	// calculate the total score
	err = calculateScore(maximumAchievableScore, evaluationResult, sloConfig, keySLIFailed)
	if err != nil {
		return sendErroredFinishedEventWithMessage(shkeptncontext, triggeredID, err.Error(), string(sloFileContent), eh.KeptnHandler, e)
	}
	eh.KeptnHandler.Logger.Debug("Evaluation result: " + string(evaluationResult.Result))

	evaluationResult.Evaluation.SLOFileContent = base64.StdEncoding.EncodeToString(sloFileContent)

	// #1289: check if test execution that preceded the evaluation was successful or failed
	testsFinishedEvent, _ := eh.getPreviousTestExecutionResult(e)
	if testsFinishedEvent != nil {
		if testsFinishedEvent.Result == keptnv2.ResultFailed {
			eh.KeptnHandler.Logger.Debug("Setting evaluation result to 'fail' because of failed preceding test execution")
			evaluationResult.Result = keptnv2.ResultFailed
			evaluationResult.Status = keptnv2.StatusErrored
			evaluationResult.Evaluation.Result = "Setting evaluation result to 'fail' because of failed preceding test execution"
		}
	}

	return sendEvent(shkeptncontext, triggeredEvents[0].ID, keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName), eh.KeptnHandler, evaluationResult)
}

func evaluateObjectives(e *keptnv2.GetSLIFinishedEventData, sloConfig *keptn.ServiceLevelObjectives, previousEvaluationEvents []*keptnv2.EvaluationFinishedEventData) (*keptnv2.EvaluationFinishedEventData, float64, bool) {
	evaluationResult := &keptnv2.EvaluationFinishedEventData{
		EventData: keptnv2.EventData{
			Status:  "",
			Project: e.Project,
			Service: e.Service,
			Stage:   e.Stage,
		},
		Evaluation: keptnv2.EvaluationDetails{
			TimeStart: e.GetSLI.Start,
			TimeEnd:   e.GetSLI.End,
		},
	}
	var sliEvaluationResults []*keptnv2.SLIEvaluationResult
	maximumAchievableScore := 0.0
	keySLIFailed := false
	for _, objective := range sloConfig.Objectives {
		// only consider the SLI for the total score if pass criteria have been included
		if len(objective.Pass) > 0 {
			maximumAchievableScore += float64(objective.Weight)
		}
		sliEvaluationResult := &keptnv2.SLIEvaluationResult{}
		result := getSLIResult(e.GetSLI.IndicatorValues, objective.SLI)

		if result == nil {
			// no result available => fail the objective
			sliEvaluationResult.Value = &keptnv2.SLIResult{
				Metric:  objective.SLI,
				Success: false,
				Message: "no value received from SLI provider",
			}
			sliEvaluationResult.Status = "fail"
			sliEvaluationResult.Score = 0
			continue
		}
		sliEvaluationResult.Value = (*keptnv2.SLIResult)(result)

		// gather the previous results for the current SLI
		var previousSLIResults []*keptnv2.SLIEvaluationResult

		if previousEvaluationEvents != nil && len(previousEvaluationEvents) > 0 {
			for _, event := range previousEvaluationEvents {
				for _, prevSLIResult := range event.Evaluation.IndicatorResults {
					if strings.Compare(prevSLIResult.Value.Metric, objective.SLI) == 0 {
						previousSLIResults = append(previousSLIResults, prevSLIResult)
					}
				}
			}
		}

		var passTargets []*keptnv2.SLITarget
		var warningTargets []*keptnv2.SLITarget
		isPassed := true
		isWarning := true
		if objective.Pass != nil && len(objective.Pass) > 0 {
			isPassed, passTargets, _ = evaluateOrCombinedCriteria(sliEvaluationResult.Value, objective.Pass, previousSLIResults, sloConfig.Comparison)
			if isPassed {
				sliEvaluationResult.Score = float64(objective.Weight)
				sliEvaluationResult.Status = "pass"
			}
		} else {
			sliEvaluationResult.Status = "info"
		}

		if !isPassed {
			if objective.Warning != nil && len(objective.Warning) > 0 {
				isWarning, warningTargets, _ = evaluateOrCombinedCriteria(sliEvaluationResult.Value, objective.Warning, previousSLIResults, sloConfig.Comparison)
				if isWarning {
					sliEvaluationResult.Score = 0.5 * float64(objective.Weight)
					sliEvaluationResult.Status = "warning"
				}
			} else {
				isWarning = false
			}
		}

		sliEvaluationResult.Targets = append(warningTargets, passTargets...)

		if !isPassed && !isWarning {
			if objective.KeySLI {
				keySLIFailed = true
			}
			sliEvaluationResult.Status = "fail"
			sliEvaluationResult.Score = 0
		}

		sliEvaluationResults = append(sliEvaluationResults, sliEvaluationResult)
	}
	evaluationResult.Evaluation.IndicatorResults = sliEvaluationResults
	return evaluationResult, maximumAchievableScore, keySLIFailed
}

func calculateScore(maximumAchievableScore float64, evaluationResult *keptnv2.EvaluationFinishedEventData, sloConfig *keptn.ServiceLevelObjectives, keySLIFailed bool) error {
	if maximumAchievableScore == 0 {
		evaluationResult.Evaluation.Result = "pass"
		evaluationResult.Result = keptnv2.ResultPass
		evaluationResult.Status = keptnv2.StatusSucceeded
		evaluationResult.Evaluation.Score = 100.0
		return nil
	}
	totalScore := 0.0
	for _, result := range evaluationResult.Evaluation.IndicatorResults {
		totalScore += result.Score
	}
	achievedPercentage := 100.0 * (totalScore / maximumAchievableScore)
	evaluationResult.Evaluation.Score = achievedPercentage
	if sloConfig.TotalScore == nil || sloConfig.TotalScore.Pass == "" {
		return errors.New("no target score defined")
	}
	passTargetPercentage, err := strconv.ParseFloat(strings.TrimSuffix(sloConfig.TotalScore.Pass, "%"), 64)
	if err != nil {
		return errors.New("could not parse pass target percentage")
	}
	if achievedPercentage >= passTargetPercentage && !keySLIFailed {
		evaluationResult.Evaluation.Result = "pass"
		evaluationResult.Result = keptnv2.ResultPass
		evaluationResult.Status = keptnv2.StatusSucceeded
	} else if sloConfig.TotalScore.Warning != "" && !keySLIFailed {
		warnTargetPercentage, err := strconv.ParseFloat(strings.TrimSuffix(sloConfig.TotalScore.Warning, "%"), 64)

		if err != nil {
			return errors.New("could not parse warning target percentage")
		}
		if achievedPercentage >= warnTargetPercentage {
			evaluationResult.Evaluation.Result = "warning"
			evaluationResult.Result = keptnv2.ResultWarning
			evaluationResult.Status = keptnv2.StatusSucceeded
		} else {
			evaluationResult.Evaluation.Result = "fail"
			evaluationResult.Result = keptnv2.ResultFailed
			evaluationResult.Status = keptnv2.StatusSucceeded
		}
	} else {
		evaluationResult.Evaluation.Result = "fail"
		evaluationResult.Result = keptnv2.ResultFailed
		evaluationResult.Status = keptnv2.StatusSucceeded
	}
	return nil
}

func getSLIResult(results []*keptnv2.SLIResult, sli string) *keptnv2.SLIResult {
	for _, sliResult := range results {
		if sliResult.Metric == sli {
			return sliResult
		}
	}
	return nil
}

func evaluateOrCombinedCriteria(result *keptnv2.SLIResult, sloCriteria []*keptn.SLOCriteria, previousResults []*keptnv2.SLIEvaluationResult, comparison *keptn.SLOComparison) (bool, []*keptnv2.SLITarget, error) {
	var satisfied bool
	satisfied = false
	var sliTargets []*keptnv2.SLITarget
	for _, crit := range sloCriteria {
		criteriaSatisfied, evaluatedTargets, _ := evaluateCriteriaSet(result, crit, previousResults, comparison)
		if criteriaSatisfied {
			// one matching criteria set is sufficient to satisfy the evaluation. Other criteria sets are evaluated nevertheless, to get potential violations
			satisfied = true
		}
		for _, evaluatedTarget := range evaluatedTargets {
			sliTargets = append(sliTargets, evaluatedTarget)
		}
	}
	return satisfied, sliTargets, nil
}

// evaluateCriteria evaluates a set of criteria strings. Per definition, all criteria clauses within a SLOCriteria object have to be fulfilled to satisfy the SLOCriteria
func evaluateCriteriaSet(result *keptnv2.SLIResult, sloCriteria *keptn.SLOCriteria, previousResults []*keptnv2.SLIEvaluationResult, comparison *keptn.SLOComparison) (bool, []*keptnv2.SLITarget, error) {
	satisfied := true
	var sliTargets []*keptnv2.SLITarget
	for _, criteria := range sloCriteria.Criteria {
		target := &keptnv2.SLITarget{
			Criteria: criteria,
		}
		criteriaSatisfied, _ := evaluateSingleCriteria(result, criteria, previousResults, comparison, target)
		if !criteriaSatisfied {
			target.Violated = true
			satisfied = false
		} else {
			target.Violated = false
		}
		sliTargets = append(sliTargets, target)
	}
	return satisfied, sliTargets, nil
}

func evaluateSingleCriteria(sliResult *keptnv2.SLIResult, criteria string, previousResults []*keptnv2.SLIEvaluationResult, comparison *keptn.SLOComparison, violation *keptnv2.SLITarget) (bool, error) {
	if !sliResult.Success {
		return false, errors.New("cannot evaluate invalid SLI result")
	}

	co, err := parseCriteriaString(criteria)

	if err != nil {
		return false, err
	}

	if !co.IsComparison {
		// do a fixed threshold comparison
		return evaluateFixedThreshold(sliResult, co, violation)
	}

	return evaluateComparison(sliResult, co, previousResults, comparison, violation)
}

func evaluateComparison(sliResult *keptnv2.SLIResult, co *criteriaObject, previousResults []*keptnv2.SLIEvaluationResult, comparison *keptn.SLOComparison, violation *keptnv2.SLITarget) (bool, error) {
	// aggregate previous results
	var aggregatedValue float64
	var targetValue float64
	var previousValues []float64

	if len(previousResults) == 0 {
		// if no comparison values are available, the evaluation passes
		return true, nil
	}

	for _, val := range previousResults {
		if val.Value.Success == true {
			// always include
			previousValues = append(previousValues, val.Value.Value)
		}
	}

	if len(previousValues) == 0 {
		// if no comparison values are available, the evaluation passes
		return true, nil
	}

	// aggregate the previous values based on the passed aggregation function
	switch comparison.AggregateFunction {
	case "avg":
		aggregatedValue = calculateAverage(previousValues)
	case "p50":
		aggregatedValue = calculatePercentile(sort.Float64Slice(previousValues), 0.5)
	case "p90":
		aggregatedValue = calculatePercentile(sort.Float64Slice(previousValues), 0.9)
	case "p95":
		aggregatedValue = calculatePercentile(sort.Float64Slice(previousValues), 0.95)
	default:
		break
	}

	// calculate the comparison value
	if co.CheckPercentage && co.CheckIncrease {
		targetValue = (aggregatedValue * (100.0 + co.Value)) / 100.0
	} else if co.CheckPercentage && !co.CheckIncrease {
		targetValue = (aggregatedValue * (100.0 - co.Value)) / 100.0
	} else if !co.CheckPercentage && co.CheckIncrease {
		targetValue = aggregatedValue + co.Value
	} else if !co.CheckPercentage && !co.CheckIncrease {
		targetValue = aggregatedValue - co.Value
	}
	violation.TargetValue = targetValue
	// compare!
	return evaluateValue(sliResult.Value, targetValue, co.Operator)
}

func calculateAverage(values []float64) float64 {
	sum := 0.0

	for _, value := range values {
		sum += value
	}
	if len(values) > 0 {
		return sum / float64(len(values))
	}
	return 0.0
}

func calculatePercentile(values sort.Float64Slice, perc float64) float64 {
	if len(values) == 0 {
		return 0.0
	}
	ps := []float64{perc}

	scores := make([]float64, len(ps))
	size := len(values)
	if size > 0 {
		sort.Sort(values)
		for i, p := range ps {
			pos := p * float64(size+1) //ALTERNATIVELY, DROP THE +1
			if pos < 1.0 {
				scores[i] = float64(values[0])
			} else if pos >= float64(size) {
				scores[i] = float64(values[size-1])
			} else {
				lower := float64(values[int(pos)-1])
				upper := float64(values[int(pos)])
				scores[i] = lower + (pos-math.Floor(pos))*(upper-lower)
			}
		}
	}
	return scores[0]
}

func evaluateFixedThreshold(sliResult *keptnv2.SLIResult, co *criteriaObject, violation *keptnv2.SLITarget) (bool, error) {
	violation.TargetValue = co.Value
	return evaluateValue(sliResult.Value, co.Value, co.Operator)
}

func evaluateValue(measured float64, expected float64, operator string) (bool, error) {
	switch operator {
	case "<":
		return measured < expected, nil
	case "<=":
		return measured <= expected, nil
	case "=":
		return measured == expected, nil
	case ">=":
		return measured >= expected, nil
	case ">":
		return measured > expected, nil
	default:
		return false, errors.New("no operator set")
	}
}

func parseCriteriaString(criteria string) (*criteriaObject, error) {
	// example values: <+15%, <500, >-8%, =0
	// possible operators: <, <=, =, >, >=
	// regex: ^([<|<=|=|>|>=]{1,2})([+|-]{0,1}\\d*\.?\d*)([%]{0,1})
	regex := `^([<|<=|=|>|>=]{1,2})([+|-]{0,1}\d*\.?\d*)([%]{0,1})`
	var re *regexp.Regexp
	re = regexp.MustCompile(regex)

	// remove whitespaces
	criteria = strings.Replace(criteria, " ", "", -1)

	if !re.MatchString(criteria) {
		return nil, errors.New("invalid criteria string")
	}

	c := &criteriaObject{}

	operators := []string{"<=", "<", "=", ">=", ">"}

	for _, operator := range operators {
		if strings.HasPrefix(criteria, operator) {
			c.Operator = operator
			criteria = strings.TrimPrefix(criteria, operator)
			break
		}
	}

	if strings.HasSuffix(criteria, "%") {
		c.CheckPercentage = true
		c.IsComparison = true // Issue #1498: criteria containing '%' is always a comparison
		c.CheckIncrease = true
		criteria = strings.TrimSuffix(criteria, "%")
	}

	if strings.HasPrefix(criteria, "-") {
		c.IsComparison = true
		c.CheckIncrease = false
		criteria = strings.TrimPrefix(criteria, "-")
	} else if strings.HasPrefix(criteria, "+") {
		c.IsComparison = true
		c.CheckIncrease = true
		criteria = strings.TrimPrefix(criteria, "+")
	}

	floatValue, err := strconv.ParseFloat(criteria, 64)
	if err != nil {
		return nil, errors.New("could not parse criteria target value")
	}
	c.Value = floatValue

	return c, nil
}

// gets previous evaluation.finished events from mongodb-datastore
func (eh *EvaluateSLIHandler) getPreviousEvaluations(e *keptnv2.GetSLIFinishedEventData, numberOfPreviousResults int, includeResult string) ([]*keptnv2.EvaluationFinishedEventData, []string, error) {
	var evaluationDoneEvents []*keptnv2.EvaluationFinishedEventData
	var eventIDs []string

	// previous results are fetched from mongodb datastore with source=lighthouse-service
	queryString := fmt.Sprintf("source=%s&limit=%d&excludeInvalidated=true&",
		"lighthouse-service", numberOfPreviousResults)

	includeResult = strings.ToLower(includeResult)

	filter := "filter=data.project:" + e.Project + "%20AND%20data.stage:" + e.Stage + "%20AND%20data.service:" + e.Service
	switch includeResult {
	case "pass":
		filter = filter + "%20AND%20data.result:pass"
		break
	case "pass_or_warn":
		filter = filter + "%20AND%20data.result:pass,warning"
		break
	default:
		break
	}

	queryString = queryString + filter

	req, err := http.NewRequest("GET", getDatastoreURL()+"/event/type/"+keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName)+"?"+queryString, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := eh.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, nil, errors.New("could not retrieve previous evaluation.finished events")
	}
	previousEvents := &datastoreResult{}
	err = json.Unmarshal(body, previousEvents)
	if err != nil {
		return nil, nil, err
	}

	// iterate over previous events
	for _, event := range previousEvents.Events {
		bytes, err := json.Marshal(event.Data)
		if err != nil {
			continue
		}
		var evaluationDoneEvent keptnv2.EvaluationFinishedEventData
		err = json.Unmarshal(bytes, &evaluationDoneEvent)

		if err != nil {
			continue
		}
		evaluationDoneEvents = append(evaluationDoneEvents, &evaluationDoneEvent)
		eventIDs = append(eventIDs, event.ID)
		if len(evaluationDoneEvents) == numberOfPreviousResults {
			return evaluationDoneEvents, eventIDs, nil
		}
	}

	return evaluationDoneEvents, eventIDs, nil
}

func (eh *EvaluateSLIHandler) getPreviousTestExecutionResult(e *keptnv2.GetSLIFinishedEventData) (*keptnv2.TestFinishedEventData, error) {
	events, _ := eh.KeptnHandler.EventHandler.GetEvents(&keptnapi.EventFilter{
		Project:      e.Project,
		Stage:        e.Stage,
		Service:      e.Service,
		EventType:    keptnv2.GetFinishedEventType(keptnv2.TestTaskName),
		KeptnContext: eh.KeptnHandler.KeptnContext,
	})
	if events == nil || len(events) == 0 {
		msg := "Could not retrieve test.finished event for context " + eh.KeptnHandler.KeptnContext
		eh.KeptnHandler.Logger.Error(msg)
		return nil, errors.New(msg)
	}

	testsFinishedEvent := &keptnv2.TestFinishedEventData{}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Squash: true,
		Result: testsFinishedEvent,
	})
	err = decoder.Decode(events[0].Data)
	if err != nil {
		eh.KeptnHandler.Logger.Error("Cannot decode approval.triggered event: " + err.Error())
		return nil, err
	}
	return testsFinishedEvent, nil

}
