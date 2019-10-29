package event_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/go-openapi/strfmt"
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

type eventBody struct {
	// contenttype
	Contenttype string `json:"contenttype,omitempty"`

	// data
	Data interface{} `json:"data,omitempty"`

	// extensions
	Extensions interface{} `json:"extensions,omitempty"`

	// id
	// Required: true
	ID *string `json:"id"`

	// source
	// Required: true
	Source *string `json:"source"`

	// specversion
	// Required: true
	Specversion *string `json:"specversion"`

	// time
	// Format: date-time
	Time strfmt.DateTime `json:"time,omitempty"`

	// type
	// Required: true
	Type *string `json:"type"`

	// shkeptncontext
	Shkeptncontext string `json:"shkeptncontext,omitempty"`
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

	// get results of previous evaluations from data store (mongodb-datastore.keptn-datastore.svc.cluster.local)
	previousEvaluations, err := eh.getPreviousEvaluations(e)
	if err != nil {
		return err
	}

	// compare the results based on the evaluation strategy
	objectives, err := getSLOs(e.Project, e.Stage, e.Service)
	if err != nil {
		return err
	}

	fmt.Println(objectives)
	if len(previousEvaluations) == 0 {

	}

	// send the evaluation-done-event
	return nil
}

func evaluateSingleCriteria(sliResult *keptnevents.SLIResult, criteria string, previousResults []*keptnevents.SLIEvaluationResult, comparison *keptnmodelsv2.SLOComparison) (bool, error) {
	if !sliResult.Success {
		return false, errors.New("cannot evaluate invalid SLI result")
	}
	co, err := parseCriteriaString(criteria)

	if err != nil {
		return false, err
	}

	if !co.IsComparison {
		// do a fixed threshold comparison
		return evaluateFixedThreshold(sliResult, co)
	}

	return evaluateComparison(sliResult, co, previousResults, comparison)
}

func evaluateComparison(sliResult *keptnevents.SLIResult, co *criteriaObject, previousResults []*keptnevents.SLIEvaluationResult, comparison *keptnmodelsv2.SLOComparison) (bool, error) {
	// aggregate previous results
	var aggregatedValue float64
	var targetValue float64
	var previousValues []float64

	for _, val := range previousResults {
		if !val.Value.Success {
			continue
		}
		if comparison.IncludeResultWithScore == "all" {
			previousValues = append(previousValues, val.Value.Value)
		} else if comparison.IncludeResultWithScore == "pass_or_warn" {
			if val.Status == "warning" || val.Status == "pass" {
				previousValues = append(previousValues, val.Value.Value)
			}
		} else if comparison.IncludeResultWithScore == "pass" {
			if val.Status == "pass" {
				previousValues = append(previousValues, val.Value.Value)
			}
		}
	}
	// aggregate the previous values based on the passed aggregation function
	switch comparison.AggregateFunction {
	case "avg":
		aggregatedValue = calculateAverage(previousValues)
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

	// compare!
	return evaluateValue(sliResult.Value, targetValue, co.Operator)

	return false, nil
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

func evaluateFixedThreshold(sliResult *keptnevents.SLIResult, co *criteriaObject) (bool, error) {
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
	regex := `^([<|<=|=|>|>=]{1,2})([+|-]{0,1}\\d*\.?\d*)([%]{0,1})`
	var re *regexp.Regexp
	re = regexp.MustCompile(regex)

	// remove whitespaces
	criteria = strings.Trim(criteria, " ")

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
		c.CheckIncrease = true
	}

	return c, nil
}

func (eh *EvaluateSLIHandler) getPreviousEvaluations(e *keptnevents.InternalGetSLIDoneEventData) ([]*keptnevents.EvaluationDoneEventData, error) {
	queryString := "http://mongodb-datastore.keptn-datastore:8080/event?type=" + keptnevents.EvaluationDoneEventType + "&project=" + e.Project + "&stage=" + e.Stage + "&service=" + e.Service + "&pageSize=20"
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

func (eh *EvaluateSLIHandler) sendEvaluationDoneEvent(shkeptncontext string, project string,
	service string, stage string) error {

	source, _ := url.Parse("lighthouse-service")
	contentType := "application/json"

	configChangedEvent := keptnevents.EvaluationDoneEventData{
		Project: project,
		Service: service,
		Stage:   stage,
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptnevents.EvaluationDoneEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: configChangedEvent,
	}

	return sendEvent(event)
}
