package event_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnmodelsv2 "github.com/keptn/go-utils/pkg/models/v2"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
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
		if sloConfig.Comparison.IncludeResultWithScore == "all" {
			// always include
			filteredPreviousEvaluationEvents = append(filteredPreviousEvaluationEvents, val)
		} else if sloConfig.Comparison.IncludeResultWithScore == "pass_or_warn" {
			// only include warnings and passes
			if val.Result == "warning" || val.Result == "pass" {
				filteredPreviousEvaluationEvents = append(filteredPreviousEvaluationEvents, val)
			}
		} else if sloConfig.Comparison.IncludeResultWithScore == "pass" {
			// only include passes
			if val.Result == "pass" {
				filteredPreviousEvaluationEvents = append(filteredPreviousEvaluationEvents, val)
			}
		}
	}

	fmt.Println(sloConfig)

	evaluationResult, maximumAchievableScore, keySLIFailed := evaluateObjectives(e, sloConfig, filteredPreviousEvaluationEvents)

	// calculate the total score
	err = calculateScore(maximumAchievableScore, evaluationResult, sloConfig, keySLIFailed)
	if err != nil {
		return err
	}

	// send the evaluation-done-event
	var shkeptncontext string
	eh.Event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)
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
		maximumAchievableScore += float64(objective.Weight)
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

		var passedViolations []*keptnevents.SLIViolation
		var warningViolations []*keptnevents.SLIViolation
		isPassed := true
		isWarning := true
		if objective.Pass != nil {
			isPassed, passedViolations, _ = evaluateOrCombinedCriteria(sliEvaluationResult.Value, objective.Pass, previousSLIResults, sloConfig.Comparison)

			sliEvaluationResult.Violations = passedViolations
			if isPassed {
				sliEvaluationResult.Score = float64(objective.Weight)
				sliEvaluationResult.Status = "pass"
			}
		}

		if !isPassed {
			if objective.Warning != nil {
				isWarning, warningViolations, _ = evaluateOrCombinedCriteria(sliEvaluationResult.Value, objective.Warning, previousSLIResults, sloConfig.Comparison)
				sliEvaluationResult.Violations = warningViolations
				if isWarning {
					sliEvaluationResult.Score = 0.5 * float64(objective.Weight)
					sliEvaluationResult.Status = "warning"
				}
			}
		}

		sliEvaluationResult.Violations = append(warningViolations, passedViolations...)

		if !isPassed && !isWarning {
			if objective.KeySLI {
				keySLIFailed = true
			}
			sliEvaluationResult.Status = "failed"
			sliEvaluationResult.Score = 0
		}

		sliEvaluationResults = append(sliEvaluationResults, sliEvaluationResult)
	}
	evaluationResult.EvaluationDetails.IndicatorResults = sliEvaluationResults
	return evaluationResult, maximumAchievableScore, keySLIFailed
}

func calculateScore(maximumAchievableScore float64, evaluationResult *keptnevents.EvaluationDoneEventData, sloConfig *keptnmodelsv2.ServiceLevelObjectives, keySLIFailed bool) error {
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

func evaluateOrCombinedCriteria(result *keptnevents.SLIResult, sloCriteria []*keptnmodelsv2.SLOCriteria, previousResults []*keptnevents.SLIEvaluationResult, comparison *keptnmodelsv2.SLOComparison) (bool, []*keptnevents.SLIViolation, error) {
	var satisfied bool
	satisfied = false
	var violations []*keptnevents.SLIViolation
	for _, crit := range sloCriteria {
		criteriaSatisfied, v, _ := evaluateCriteriaSet(result, crit, previousResults, comparison)
		if criteriaSatisfied {
			// one matching criteria set is sufficient to satisfy the evaluation. Other criteria sets are evaluated nevertheless, to get potential violations
			satisfied = true
		} else {
			for _, violation := range v {
				violations = append(violations, violation)
			}
		}
	}
	return satisfied, violations, nil
}

// evaluateCriteria evaluates a set of criteria strings. Per definition, all criteria clauses within a SLOCriteria object have to be fulfilled to satisfy the SLOCriteria
func evaluateCriteriaSet(result *keptnevents.SLIResult, sloCriteria *keptnmodelsv2.SLOCriteria, previousResults []*keptnevents.SLIEvaluationResult, comparison *keptnmodelsv2.SLOComparison) (bool, []*keptnevents.SLIViolation, error) {
	satisfied := true
	var violations []*keptnevents.SLIViolation
	for _, criteria := range sloCriteria.Criteria {
		violation := &keptnevents.SLIViolation{
			Criteria: criteria,
		}
		criteriaSatisfied, _ := evaluateSingleCriteria(result, criteria, previousResults, comparison, violation)
		if !criteriaSatisfied {
			satisfied = false
			violations = append(violations, violation)
		}
	}
	return satisfied, violations, nil
}

func evaluateSingleCriteria(sliResult *keptnevents.SLIResult, criteria string, previousResults []*keptnevents.SLIEvaluationResult, comparison *keptnmodelsv2.SLOComparison, violation *keptnevents.SLIViolation) (bool, error) {
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

func evaluateComparison(sliResult *keptnevents.SLIResult, co *criteriaObject, previousResults []*keptnevents.SLIEvaluationResult, comparison *keptnmodelsv2.SLOComparison, violation *keptnevents.SLIViolation) (bool, error) {
	// aggregate previous results
	var aggregatedValue float64
	var targetValue float64
	var previousValues []float64

	if len(previousResults) == 0 {
		// if no comparison values are available, the evaluation passes
		return true, nil
	}

	for _, val := range previousResults {
		previousValues = append(previousValues, val.Value.Value)
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

func evaluateFixedThreshold(sliResult *keptnevents.SLIResult, co *criteriaObject, violation *keptnevents.SLIViolation) (bool, error) {
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

	if strings.HasSuffix(criteria, "%") {
		c.CheckPercentage = true
		criteria = strings.TrimSuffix(criteria, "%")
	}

	operators := []string{"<=", "<", "=", ">=", ">"}

	for _, operator := range operators {
		if strings.HasPrefix(criteria, operator) {
			c.Operator = operator
			criteria = strings.TrimPrefix(criteria, operator)
			break
		}
	}

	if strings.HasPrefix(criteria, "-") {
		c.IsComparison = true
		c.CheckIncrease = false
		criteria = strings.TrimPrefix(criteria, "-")
	} else if strings.HasPrefix(criteria, "+") {
		c.IsComparison = true
		c.CheckIncrease = true
		criteria = strings.TrimPrefix(criteria, "+")
	} else {
		c.IsComparison = false
		c.CheckIncrease = false
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
	queryString := fmt.Sprintf("http://mongodb-datastore.keptn-datastore:8080/event?type=%s&source=%s&project=%s&stage=%s&service=%s&pageSize=%d",
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

func extractSLIEvaluationResult(inMap []interface{}) (*keptnevents.SLIEvaluationResult, error) {
	result := &keptnevents.SLIEvaluationResult{
		Score:      0,
		Value:      nil,
		Violations: nil,
		Status:     "",
	}

	for _, value := range inMap {
		tmp := value.(map[string]interface{})
		if tmp["Key"] == "score" {
			result.Score = tmp["Value"].(float64)
		}
		if tmp["Key"] == "value" {
			sliResult, err := extractSLIResult(tmp["Value"].([]interface{}))
			if err != nil {
				return nil, err
			}
			result.Value = sliResult
		}
		if tmp["Key"] == "status" {
			result.Status = tmp["Value"].(string)
		}
	}

	return result, nil
}

func extractSLIResult(inMap []interface{}) (*keptnevents.SLIResult, error) {
	sliResult := &keptnevents.SLIResult{
		Metric:  "",
		Value:   0,
		Success: false,
		Message: "",
	}

	for _, value := range inMap {
		tmp := value.(map[string]interface{})
		if tmp["Key"] == "metric" {
			sliResult.Metric = tmp["Value"].(string)
		}
		if tmp["Key"] == "value" {
			sliResult.Value = tmp["Value"].(float64)
		}
		if tmp["Key"] == "success" {
			sliResult.Success = tmp["Value"].(bool)
		}
		if tmp["Key"] == "message" {
			sliResult.Message = tmp["Value"].(string)
		}
	}

	return sliResult, nil
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

	return sendEvent(event)
}
