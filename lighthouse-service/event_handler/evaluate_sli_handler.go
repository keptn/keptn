package event_handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/ghodss/yaml"
	"github.com/google/uuid"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnmodelsv2 "github.com/keptn/go-utils/pkg/models/v2"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
)

type datastoreResult struct {
	NextPageKey string `json:"nextPageKey"`
	TotalCount  int    `json:"totalCount"`
	PageSize    int    `json:"pageSize"`
	Events      []struct {
		Data interface{} `json:"data"`
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
	Logger     *keptnutils.Logger
	Event      cloudevents.Event
	HTTPClient *http.Client
}

func (eh *EvaluateSLIHandler) HandleEvent() error {
	e := &keptnevents.InternalGetSLIDoneEventData{}

	err := eh.Event.DataAs(&e)

	if err != nil {
		eh.Logger.Error("Could not parse event payload: " + err.Error())
		return err
	}

	eh.Logger.Debug("Start to evaluate SLIs")
	// compare the results based on the evaluation strategy
	sloConfig, err := getSLOs(e.Project, e.Stage, e.Service)
	if err != nil {
		return err
	}

	// get results of previous evaluations from data store (mongodb-datastore.keptn-datastore.svc.cluster.local)
	numberOfPreviousResults := 3
	if sloConfig.Comparison.CompareWith == "single_result" {
		numberOfPreviousResults = 1
	} else if sloConfig.Comparison.CompareWith == "several_results" {
		numberOfPreviousResults = sloConfig.Comparison.NumberOfComparisonResults
	}
	previousEvaluationEvents, err := eh.getPreviousEvaluations(e, numberOfPreviousResults)
	if err != nil {
		return err
	}

	var filteredPreviousEvaluationEvents []*keptnevents.EvaluationDoneEventData

	// verify that we have enough evaluations
	for _, val := range previousEvaluationEvents {
		filteredPreviousEvaluationEvents = append(filteredPreviousEvaluationEvents, val)
	}

	evaluationResult, maximumAchievableScore, keySLIFailed := evaluateObjectives(e, sloConfig, filteredPreviousEvaluationEvents)
	evaluationResult.Labels = e.Labels

	// calculate the total score
	err = calculateScore(maximumAchievableScore, evaluationResult, sloConfig, keySLIFailed)
	if err != nil {
		return err
	}
	eh.Logger.Debug("Evaluation result: " + evaluationResult.Result)

	sloFileContent, _ := yaml.Marshal(sloConfig)
	base64.StdEncoding.EncodeToString(sloFileContent)
	evaluationResult.EvaluationDetails.SLOFileContent = base64.StdEncoding.EncodeToString(sloFileContent)

	// send the evaluation-done-event
	var shkeptncontext string
	eh.Event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	// #1289: check if test execution that preceded the evaluation was successful or failed
	testsFinishedEvent, _ := eh.getPreviousTestExecutionResult(e, shkeptncontext)
	if testsFinishedEvent != nil {
		if testsFinishedEvent.Result == "fail" {
			eh.Logger.Debug("Setting evaluation result to 'fail' because of failed preceding test execution")
			evaluationResult.Result = "fail"
			evaluationResult.EvaluationDetails.Result = "fail"
		}
	}

	err = eh.sendEvaluationDoneEvent(shkeptncontext, evaluationResult)
	return err
}

func evaluateObjectives(e *keptnevents.InternalGetSLIDoneEventData, sloConfig *keptnmodelsv2.ServiceLevelObjectives, previousEvaluationEvents []*keptnevents.EvaluationDoneEventData) (*keptnevents.EvaluationDoneEventData, float64, bool) {
	evaluationResult := &keptnevents.EvaluationDoneEventData{
		Result:  "",
		Project: e.Project,
		Service: e.Service,
		Stage:   e.Stage,
		EvaluationDetails: &keptnevents.EvaluationDetails{
			TimeStart: e.Start,
			TimeEnd:   e.End,
		},
		TestStrategy:       e.TestStrategy,
		DeploymentStrategy: e.DeploymentStrategy,
	}
	var sliEvaluationResults []*keptnevents.SLIEvaluationResult
	maximumAchievableScore := 0.0
	keySLIFailed := false
	for _, objective := range sloConfig.Objectives {
		// only consider the SLI for the total score if pass criteria have been included
		if len(objective.Pass) > 0 {
			maximumAchievableScore += float64(objective.Weight)
		}
		sliEvaluationResult := &keptnevents.SLIEvaluationResult{}
		result := getSLIResult(e.IndicatorValues, objective.SLI)

		if result == nil {
			// no result available => fail the objective
			sliEvaluationResult.Value = &keptnevents.SLIResult{
				Metric:  objective.SLI,
				Success: false,
				Message: "no value received from SLI provider",
			}
			sliEvaluationResult.Status = "fail"
			sliEvaluationResult.Score = 0
			continue
		}
		sliEvaluationResult.Value = result

		// gather the previous results for the current SLI
		var previousSLIResults []*keptnevents.SLIEvaluationResult

		if previousEvaluationEvents != nil && len(previousEvaluationEvents) > 0 {
			for _, event := range previousEvaluationEvents {
				for _, prevSLIResult := range event.EvaluationDetails.IndicatorResults {
					if strings.Compare(prevSLIResult.Value.Metric, objective.SLI) == 0 {
						previousSLIResults = append(previousSLIResults, prevSLIResult)
					}
				}
			}
		}

		var passTargets []*keptnevents.SLITarget
		var warningTargets []*keptnevents.SLITarget
		isPassed := true
		isWarning := true
		if objective.Pass != nil {
			isPassed, passTargets, _ = evaluateOrCombinedCriteria(sliEvaluationResult.Value, objective.Pass, previousSLIResults, sloConfig.Comparison)
			if isPassed {
				sliEvaluationResult.Score = float64(objective.Weight)
				sliEvaluationResult.Status = "pass"
			}
		} else {
			sliEvaluationResult.Status = "info"
		}

		if !isPassed {
			if objective.Warning != nil {
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
	evaluationResult.EvaluationDetails.IndicatorResults = sliEvaluationResults
	return evaluationResult, maximumAchievableScore, keySLIFailed
}

func calculateScore(maximumAchievableScore float64, evaluationResult *keptnevents.EvaluationDoneEventData, sloConfig *keptnmodelsv2.ServiceLevelObjectives, keySLIFailed bool) error {
	if maximumAchievableScore == 0 {
		evaluationResult.EvaluationDetails.Result = "pass"
		evaluationResult.Result = evaluationResult.EvaluationDetails.Result
		evaluationResult.EvaluationDetails.Score = 100.0
		return nil
	}
	totalScore := 0.0
	for _, result := range evaluationResult.EvaluationDetails.IndicatorResults {
		totalScore += result.Score
	}
	achievedPercentage := 100.0 * (totalScore / maximumAchievableScore)
	evaluationResult.EvaluationDetails.Score = achievedPercentage
	if sloConfig.TotalScore == nil || sloConfig.TotalScore.Pass == "" {
		return errors.New("no target score defined")
	}
	passTargetPercentage, err := strconv.ParseFloat(strings.TrimSuffix(sloConfig.TotalScore.Pass, "%"), 64)
	if err != nil {
		return errors.New("could not parse pass target percentage")
	}
	if achievedPercentage >= passTargetPercentage && !keySLIFailed {
		evaluationResult.EvaluationDetails.Result = "pass"
	} else if sloConfig.TotalScore.Warning != "" && !keySLIFailed {
		warnTargetPercentage, err := strconv.ParseFloat(strings.TrimSuffix(sloConfig.TotalScore.Warning, "%"), 64)

		if err != nil {
			return errors.New("could not parse warning target percentage")
		}
		if achievedPercentage >= warnTargetPercentage {
			evaluationResult.EvaluationDetails.Result = "warning"
		} else {
			evaluationResult.EvaluationDetails.Result = "fail"
		}
	} else {
		evaluationResult.EvaluationDetails.Result = "fail"
	}
	evaluationResult.Result = evaluationResult.EvaluationDetails.Result
	return nil
}

func getSLIResult(results []*keptnevents.SLIResult, sli string) *keptnevents.SLIResult {
	for _, sliResult := range results {
		if sliResult.Metric == sli {
			return sliResult
		}
	}
	return nil
}

func evaluateOrCombinedCriteria(result *keptnevents.SLIResult, sloCriteria []*keptnmodelsv2.SLOCriteria, previousResults []*keptnevents.SLIEvaluationResult, comparison *keptnmodelsv2.SLOComparison) (bool, []*keptnevents.SLITarget, error) {
	var satisfied bool
	satisfied = false
	var sliTargets []*keptnevents.SLITarget
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
func evaluateCriteriaSet(result *keptnevents.SLIResult, sloCriteria *keptnmodelsv2.SLOCriteria, previousResults []*keptnevents.SLIEvaluationResult, comparison *keptnmodelsv2.SLOComparison) (bool, []*keptnevents.SLITarget, error) {
	satisfied := true
	var sliTargets []*keptnevents.SLITarget
	for _, criteria := range sloCriteria.Criteria {
		target := &keptnevents.SLITarget{
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

func evaluateSingleCriteria(sliResult *keptnevents.SLIResult, criteria string, previousResults []*keptnevents.SLIEvaluationResult, comparison *keptnmodelsv2.SLOComparison, violation *keptnevents.SLITarget) (bool, error) {
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

func evaluateComparison(sliResult *keptnevents.SLIResult, co *criteriaObject, previousResults []*keptnevents.SLIEvaluationResult, comparison *keptnmodelsv2.SLOComparison, violation *keptnevents.SLITarget) (bool, error) {
	// aggregate previous results
	var aggregatedValue float64
	var targetValue float64
	var previousValues []float64

	if len(previousResults) == 0 {
		// if no comparison values are available, the evaluation passes
		return true, nil
	}

	for _, val := range previousResults {
		if comparison.IncludeResultWithScore == "all" {
			if val.Value.Success == true {
				// always include
				previousValues = append(previousValues, val.Value.Value)
			}
		} else if comparison.IncludeResultWithScore == "pass_or_warn" {
			// only include warnings and passes
			if (val.Status == "warning" || val.Status == "pass") && val.Value.Success == true {
				previousValues = append(previousValues, val.Value.Value)
			}
		} else if comparison.IncludeResultWithScore == "pass" {
			// only include passes
			if val.Status == "pass" && val.Value.Success == true {
				previousValues = append(previousValues, val.Value.Value)
			}
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

func evaluateFixedThreshold(sliResult *keptnevents.SLIResult, co *criteriaObject, violation *keptnevents.SLITarget) (bool, error) {
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

// gets previous evaluation-done events from mongodb-datastore
func (eh *EvaluateSLIHandler) getPreviousEvaluations(e *keptnevents.InternalGetSLIDoneEventData, numberOfPreviousResults int) ([]*keptnevents.EvaluationDoneEventData, error) {
	// previous results are fetched from mongodb datastore with source=lighthouse-service
	queryString := fmt.Sprintf(getDatastoreURL()+"/event?type=%s&source=%s&project=%s&stage=%s&service=%s&pageSize=%d",
		keptnevents.EvaluationDoneEventType, "lighthouse-service",
		e.Project, e.Stage, e.Service, numberOfPreviousResults)

	req, err := http.NewRequest("GET", queryString, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := eh.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, errors.New("could not retrieve previous evaluation-done events")
	}
	previousEvents := &datastoreResult{}
	err = json.Unmarshal(body, previousEvents)
	if err != nil {
		return nil, err
	}
	var evaluationDoneEvents []*keptnevents.EvaluationDoneEventData

	// iterate over previous events
	for _, event := range previousEvents.Events {
		bytes, err := json.Marshal(event.Data)
		if err != nil {
			continue
		}
		var evaluationDoneEvent keptnevents.EvaluationDoneEventData
		err = json.Unmarshal(bytes, &evaluationDoneEvent)

		if err != nil {
			continue
		}
		evaluationDoneEvents = append(evaluationDoneEvents, &evaluationDoneEvent)
	}
	return evaluationDoneEvents, nil
}

func (eh *EvaluateSLIHandler) getPreviousTestExecutionResult(e *keptnevents.InternalGetSLIDoneEventData, keptnContext string) (*keptnevents.TestsFinishedEventData, error) {
	queryString := fmt.Sprintf(getDatastoreURL()+"/event?type=%s&project=%s&stage=%s&service=%s&keptnContext=%s&pageSize=%d",
		keptnevents.TestsFinishedEventType,
		e.Project, e.Stage, e.Service, keptnContext, 1)

	req, err := http.NewRequest("GET", queryString, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := eh.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, errors.New("could not retrieve previous evaluation-done events")
	}
	previousEvents := &datastoreResult{}
	err = json.Unmarshal(body, previousEvents)
	if err != nil {
		return nil, err
	}
	if len(previousEvents.Events) == 0 {
		return nil, nil
	}

	bytes, err := json.Marshal(previousEvents.Events[0].Data)
	if err != nil {
		return nil, err
	}

	testsFinishedEvent := &keptnevents.TestsFinishedEventData{}
	err = json.Unmarshal(bytes, &testsFinishedEvent)
	if err != nil {
		return nil, err
	}
	return testsFinishedEvent, nil

}

func (eh *EvaluateSLIHandler) sendEvaluationDoneEvent(shkeptncontext string, data *keptnevents.EvaluationDoneEventData) error {

	source, _ := url.Parse("lighthouse-service")
	contentType := "application/json"

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptnevents.EvaluationDoneEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: data,
	}

	eh.Logger.Debug("Send event: " + keptnevents.EvaluationDoneEventType)
	return sendEvent(event)
}
