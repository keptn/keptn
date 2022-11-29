package event_handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/common/strutils"
	apimodelsv2 "github.com/keptn/go-utils/pkg/lib"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	keptnfake "github.com/keptn/go-utils/pkg/lib/v0_2_0/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	event_handler_mock "github.com/keptn/keptn/lighthouse-service/event_handler/fake"
)

type operatorParserTest struct {
	Criteria               string
	ExpectedCriteriaObject *criteriaObject
}

func TestParseCriteriaString(t *testing.T) {
	tests := []*operatorParserTest{
		{
			Criteria: "<10",
			ExpectedCriteriaObject: &criteriaObject{
				Operator:        "<",
				Value:           10,
				CheckPercentage: false,
				IsComparison:    false,
				CheckIncrease:   false,
			},
		}, {
			Criteria: "<=10",
			ExpectedCriteriaObject: &criteriaObject{
				Operator:        "<=",
				Value:           10,
				CheckPercentage: false,
				IsComparison:    false,
				CheckIncrease:   false,
			},
		}, {
			Criteria: "=10",
			ExpectedCriteriaObject: &criteriaObject{
				Operator:        "=",
				Value:           10,
				CheckPercentage: false,
				IsComparison:    false,
				CheckIncrease:   false,
			},
		}, {
			Criteria: ">=10",
			ExpectedCriteriaObject: &criteriaObject{
				Operator:        ">=",
				Value:           10,
				CheckPercentage: false,
				IsComparison:    false,
				CheckIncrease:   false,
			},
		}, {
			Criteria: ">10",
			ExpectedCriteriaObject: &criteriaObject{
				Operator:        ">",
				Value:           10,
				CheckPercentage: false,
				IsComparison:    false,
				CheckIncrease:   false,
			},
		}, {
			Criteria: ">-10%",
			ExpectedCriteriaObject: &criteriaObject{
				Operator:        ">",
				Value:           10,
				CheckPercentage: true,
				IsComparison:    true,
				CheckIncrease:   false,
			},
		}, {
			Criteria: "<=+10.5%",
			ExpectedCriteriaObject: &criteriaObject{
				Operator:        "<=",
				Value:           10.5,
				CheckPercentage: true,
				IsComparison:    true,
				CheckIncrease:   true,
			},
		}, {
			Criteria: "<=+10",
			ExpectedCriteriaObject: &criteriaObject{
				Operator:        "<=",
				Value:           10,
				CheckPercentage: false,
				IsComparison:    true,
				CheckIncrease:   true,
			},
		},
		{
			Criteria: "  <=+10   %",
			ExpectedCriteriaObject: &criteriaObject{
				Operator:        "<=",
				Value:           10,
				CheckPercentage: true,
				IsComparison:    true,
				CheckIncrease:   true,
			},
		},
		{
			Criteria: "  <=10%",
			ExpectedCriteriaObject: &criteriaObject{
				Operator:        "<=",
				Value:           10,
				CheckPercentage: true,
				IsComparison:    true,
				CheckIncrease:   true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Criteria, func(t *testing.T) {
			co, _ := parseCriteriaString(test.Criteria)
			assert.EqualValues(t, test.ExpectedCriteriaObject.Operator, co.Operator)
			assert.EqualValues(t, test.ExpectedCriteriaObject.Value, co.Value)
			assert.EqualValues(t, test.ExpectedCriteriaObject.CheckPercentage, co.CheckPercentage)
			assert.EqualValues(t, test.ExpectedCriteriaObject.IsComparison, co.IsComparison)
			assert.EqualValues(t, test.ExpectedCriteriaObject.CheckIncrease, co.CheckIncrease)
		})
	}
}

type evaluateValueTestObject struct {
	Name           string
	MeasuredValue  float64
	ExpectedValue  float64
	Operator       string
	ExpectedResult bool
	ExpectedError  error
}

func TestEvaluateValue(t *testing.T) {
	tests := []*evaluateValueTestObject{
		{
			Name:           "10 > 9 should return true",
			MeasuredValue:  10.0,
			ExpectedValue:  9.0,
			Operator:       ">",
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name:           "10 >= 10 should return true",
			MeasuredValue:  10.0,
			ExpectedValue:  10.0,
			Operator:       ">=",
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name:           "10 <= 10 should return true",
			MeasuredValue:  10.0,
			ExpectedValue:  10.0,
			Operator:       "<=",
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name:           "10 < 10 should return false",
			MeasuredValue:  10.0,
			ExpectedValue:  10.0,
			Operator:       "<",
			ExpectedResult: false,
			ExpectedError:  nil,
		},
		{
			Name:           "10 > 10 should return false",
			MeasuredValue:  10.0,
			ExpectedValue:  10.0,
			Operator:       ">",
			ExpectedResult: false,
			ExpectedError:  nil,
		},
		{
			Name:           "10 ? 10 should return an error",
			MeasuredValue:  10.0,
			ExpectedValue:  10.0,
			Operator:       "?",
			ExpectedResult: false,
			ExpectedError:  errors.New("no operator set"),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := evaluateValue(test.MeasuredValue, test.ExpectedValue, test.Operator)
			assert.EqualValues(t, test.ExpectedResult, result)
			assert.EqualValues(t, test.ExpectedError, err)
		})
	}
}

type evaluateFixedThresholdTestObject struct {
	Name             string
	InSLIResult      *keptnv2.SLIResult
	InCriteriaObject *criteriaObject
	InTarget         *keptnv2.SLITarget
	ExpectedResult   bool
	ExpectedError    error
}

func TestEvaluateFixedThreshold(t *testing.T) {
	tests := []*evaluateFixedThresholdTestObject{
		{
			Name: "10.0 > 9.0 should return true",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteriaObject: &criteriaObject{
				Operator:        ">",
				Value:           9.0,
				CheckPercentage: false,
				IsComparison:    false,
				CheckIncrease:   false,
			},
			InTarget: &keptnv2.SLITarget{
				Criteria: ">9.0",
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name: "10.0 = 9.0 should return false",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteriaObject: &criteriaObject{
				Operator:        "=",
				Value:           9.0,
				CheckPercentage: false,
				IsComparison:    false,
				CheckIncrease:   false,
			},
			InTarget: &keptnv2.SLITarget{
				Criteria: "=9.0",
			},
			ExpectedResult: false,
			ExpectedError:  nil,
		},
		{
			Name: "10.0 ? 9.0 should return an error",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteriaObject: &criteriaObject{
				Operator:        "?",
				Value:           9.0,
				CheckPercentage: false,
				IsComparison:    false,
				CheckIncrease:   false,
			},
			InTarget: &keptnv2.SLITarget{
				Criteria: "?9.0",
			},
			ExpectedResult: false,
			ExpectedError:  errors.New("no operator set"),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := evaluateFixedThreshold(test.InSLIResult, test.InCriteriaObject, test.InTarget)

			assert.EqualValues(t, test.ExpectedResult, result)
			assert.EqualValues(t, test.ExpectedError, err)
			assert.EqualValues(t, test.InTarget.TargetValue, test.InCriteriaObject.Value)
		})
	}
}

type calculatePercentileTestObject struct {
	Name          string
	InValue       sort.Float64Slice
	InPercentile  float64
	ExpectedValue float64
}

func TestCalculatePercentile(t *testing.T) {
	tests := []*calculatePercentileTestObject{
		{
			Name:          "Should return 5.0",
			InValue:       []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			InPercentile:  0.5,
			ExpectedValue: 5.0,
		},
		{
			Name:          "Should return 9.0",
			InValue:       []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			InPercentile:  0.9,
			ExpectedValue: 9.8,
		},
		{
			Name:          "Should return 10.0",
			InValue:       []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			InPercentile:  0.95,
			ExpectedValue: 10.0,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			percentile := calculatePercentile(test.InValue, test.InPercentile)
			assert.EqualValues(t, test.ExpectedValue, percentile)
		})
	}
}

type evaluateComparisonTestObject struct {
	Name              string
	InSLIResult       *keptnv2.SLIResult
	InCriteriaObject  *criteriaObject
	InPreviousResults []*keptnv2.SLIEvaluationResult
	InComparison      *apimodelsv2.SLOComparison
	InTarget          *keptnv2.SLITarget
	ExpectedResult    bool
	ExpectedError     error
}

func TestEvaluateComparison(t *testing.T) {
	tests := []*evaluateComparisonTestObject{
		{
			Name: "Expect true for 10.0 <= avg([10.0, 10.0]) + 10%",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteriaObject: &criteriaObject{
				Operator:        "<=",
				Value:           10.0,
				CheckPercentage: true,
				IsComparison:    true,
				CheckIncrease:   true,
			},
			InPreviousResults: []*keptnv2.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
			},
			InComparison: &apimodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnv2.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name: "Expect false for 11.01 <= avg([10.0, 10.0]) + 10%",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   11.01,
				Success: true,
				Message: "",
			},
			InCriteriaObject: &criteriaObject{
				Operator:        "<=",
				Value:           10.0,
				CheckPercentage: true,
				IsComparison:    true,
				CheckIncrease:   true,
			},
			InPreviousResults: []*keptnv2.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
			},
			InComparison: &apimodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnv2.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: false,
			ExpectedError:  nil,
		},
		{
			Name: "Expect false for 10.0 < avg([10.0, 10.0]) + 0%",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteriaObject: &criteriaObject{
				Operator:        "<",
				Value:           0.0,
				CheckPercentage: true,
				IsComparison:    true,
				CheckIncrease:   true,
			},
			InPreviousResults: []*keptnv2.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
			},
			InComparison: &apimodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnv2.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: false,
			ExpectedError:  nil,
		},
		{
			Name: "Expect false for 10.0 > avg([10.0, 10.0]) + 0%",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteriaObject: &criteriaObject{
				Operator:        ">",
				Value:           0.0,
				CheckPercentage: true,
				IsComparison:    true,
				CheckIncrease:   true,
			},
			InPreviousResults: []*keptnv2.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
			},
			InComparison: &apimodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnv2.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: false,
			ExpectedError:  nil,
		},
		{
			Name: "Expect true for 11.0 = avg([10.0, 10.0]) + 1.0",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   11.0,
				Success: true,
				Message: "",
			},
			InCriteriaObject: &criteriaObject{
				Operator:        "=",
				Value:           1.0,
				CheckPercentage: false,
				IsComparison:    true,
				CheckIncrease:   true,
			},
			InPreviousResults: []*keptnv2.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
			},
			InComparison: &apimodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnv2.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name: "Expect true for 9.0 = avg([10.0, 10.0]) - 1.0",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   9.0,
				Success: true,
				Message: "",
			},
			InCriteriaObject: &criteriaObject{
				Operator:        "=",
				Value:           1.0,
				CheckPercentage: false,
				IsComparison:    true,
				CheckIncrease:   false,
			},
			InPreviousResults: []*keptnv2.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
			},
			InComparison: &apimodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnv2.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name: "Expect true for 10.0 <= p50([10.0, 5.0]) + 10.0%",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteriaObject: &criteriaObject{
				Operator:        ">=",
				Value:           10.0,
				CheckPercentage: true,
				IsComparison:    true,
				CheckIncrease:   false,
			},
			InPreviousResults: []*keptnv2.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "fail",
				},
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
			},
			InComparison: &apimodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "p50",
			},
			InTarget: &keptnv2.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := evaluateComparison(test.InSLIResult, test.InCriteriaObject, test.InPreviousResults, test.InComparison, test.InTarget)
			assert.EqualValues(t, test.ExpectedResult, result)
			assert.EqualValues(t, test.ExpectedError, err)
		})
	}
}

type evaluateSingleCriteriaTestObject struct {
	Name              string
	InSLIResult       *keptnv2.SLIResult
	InCriteria        string
	InPreviousResults []*keptnv2.SLIEvaluationResult
	InComparison      *apimodelsv2.SLOComparison
	InTarget          *keptnv2.SLITarget
	ExpectedResult    bool
	ExpectedError     error
}

func TestEvaluateSingleCriteria(t *testing.T) {
	tests := []*evaluateSingleCriteriaTestObject{
		{
			Name: "Expect true for 10.0 <= avg([10.0, 10.0]) + 10%",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteria: "<=+10%",
			InPreviousResults: []*keptnv2.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
			},
			InComparison: &apimodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnv2.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name: "Expect false for 11.01 <= avg([10.0, 10.0]) + 10%",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   11.01,
				Success: true,
				Message: "",
			},
			InCriteria: "<=+10%",
			InPreviousResults: []*keptnv2.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
			},
			InComparison: &apimodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnv2.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: false,
			ExpectedError:  nil,
		},
		{
			Name: "Expect false for 10.0 < avg([10.0, 10.0]) + 0%",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteria: "<+0%",
			InPreviousResults: []*keptnv2.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
			},
			InComparison: &apimodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnv2.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: false,
			ExpectedError:  nil,
		},
		{
			Name: "Expect false for 10.0 > avg([10.0, 10.0]) + 0%",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteria: ">+0%",
			InPreviousResults: []*keptnv2.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
			},
			InComparison: &apimodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnv2.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: false,
			ExpectedError:  nil,
		},
		{
			Name: "Expect true for 11.0 = avg([10.0, 10.0]) + 1.0",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   11.0,
				Success: true,
				Message: "",
			},
			InCriteria: "=+1",
			InPreviousResults: []*keptnv2.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
			},
			InComparison: &apimodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnv2.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name: "Expect true for 9.0 = avg([10.0, 10.0]) - 1.0",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   9.0,
				Success: true,
				Message: "",
			},
			InCriteria: "=-1",
			InPreviousResults: []*keptnv2.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
			},
			InComparison: &apimodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnv2.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name: "Expect true for 10.0 <= p50([10.0, 5.0]) + 10.0%",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteria: "<=+10%",
			InPreviousResults: []*keptnv2.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
			},
			InComparison: &apimodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "p50",
			},
			InTarget: &keptnv2.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name: "Expect true for 10.0 <= 10.0",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteria: "<=10",
			InTarget: &keptnv2.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := evaluateSingleCriteria(test.InSLIResult, test.InCriteria, test.InPreviousResults, test.InComparison, test.InTarget)
			assert.EqualValues(t, test.ExpectedResult, result)
			assert.EqualValues(t, test.ExpectedError, err)
		})
	}
}

type evaluateCriteriaSetTestObject struct {
	Name              string
	InSLIResult       *keptnv2.SLIResult
	InCriteriaSet     *apimodelsv2.SLOCriteria
	InPreviousResults []*keptnv2.SLIEvaluationResult
	InComparison      *apimodelsv2.SLOComparison
	ExpectedTargets   []*keptnv2.SLITarget
	ExpectedResult    bool
	ExpectedError     error
}

func TestEvaluateCriteriaSet(t *testing.T) {
	tests := []*evaluateCriteriaSetTestObject{
		{
			Name: "Expect true for (10.0 <= avg([10.0, 10.0]) + 10%) && (10.0 <= 10.0)",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteriaSet: &apimodelsv2.SLOCriteria{
				Criteria: []string{"<=+10%", "<=10.0"},
			},
			InPreviousResults: []*keptnv2.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
			},
			InComparison: &apimodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			ExpectedTargets: []*keptnv2.SLITarget{
				{
					Criteria:    "<=+10%",
					TargetValue: 11,
					Violated:    false,
				},
				{
					Criteria:    "<=10.0",
					TargetValue: 10,
					Violated:    false,
				},
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name: "Expect false for (11.0 <= avg([10.0, 10.0]) + 10%) && (10.0 <= 10.0)",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   11.0,
				Success: true,
				Message: "",
			},
			InCriteriaSet: &apimodelsv2.SLOCriteria{
				Criteria: []string{"<=+10%", "<=10.0"},
			},
			InPreviousResults: []*keptnv2.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
			},
			InComparison: &apimodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			ExpectedTargets: []*keptnv2.SLITarget{
				{
					Criteria:    "<=+10%",
					TargetValue: 11,
					Violated:    false,
				},
				{
					Criteria:    "<=10.0",
					TargetValue: 10.0,
					Violated:    true,
				},
			},
			ExpectedResult: false,
			ExpectedError:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, violations, err := evaluateCriteriaSet(test.InSLIResult, test.InCriteriaSet, test.InPreviousResults, test.InComparison)
			assert.EqualValues(t, test.ExpectedResult, result)
			assert.EqualValues(t, test.ExpectedTargets, violations)
			assert.EqualValues(t, test.ExpectedError, err)
		})
	}
}

type evaluateOrCombinedCriteriaTestObject struct {
	Name              string
	InSLIResult       *keptnv2.SLIResult
	InCriteriaSets    []*apimodelsv2.SLOCriteria
	InPreviousResults []*keptnv2.SLIEvaluationResult
	InComparison      *apimodelsv2.SLOComparison
	ExpectedTargets   []*keptnv2.SLITarget
	ExpectedResult    bool
	ExpectedError     error
}

func TestEvaluateOrCombinedCriteria(t *testing.T) {
	tests := []*evaluateOrCombinedCriteriaTestObject{
		{
			Name: "Expect true for (10.0 <= avg([10.0, 10.0]) + 10%) || (10.0 <= 10.0)",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteriaSets: []*apimodelsv2.SLOCriteria{
				{
					Criteria: []string{"<=10.0"},
				},
				{
					Criteria: []string{"<=+10%"},
				},
			},
			InPreviousResults: []*keptnv2.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
			},
			InComparison: &apimodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			ExpectedTargets: []*keptnv2.SLITarget{
				{
					Criteria:    "<=10.0",
					TargetValue: 10.0,
					Violated:    false,
				},
				{
					Criteria:    "<=+10%",
					TargetValue: 11,
					Violated:    false,
				},
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name: "Expect true for (11.0 <= avg([10.0, 10.0]) + 10%) || (10.0 <= 10.0)",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   11.0,
				Success: true,
				Message: "",
			},
			InCriteriaSets: []*apimodelsv2.SLOCriteria{
				{
					Criteria: []string{"<=10.0"},
				},
				{
					Criteria: []string{"<=+10%"},
				},
			},
			InPreviousResults: []*keptnv2.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
			},
			InComparison: &apimodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			ExpectedTargets: []*keptnv2.SLITarget{
				{
					Criteria:    "<=10.0",
					TargetValue: 10.0,
					Violated:    true,
				},
				{
					Criteria:    "<=+10%",
					TargetValue: 11,
					Violated:    false,
				},
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name: "Expect false for (20.0 <= avg([10.0, 10.0]) + 10%) || (10.0 <= 10.0)",
			InSLIResult: &keptnv2.SLIResult{
				Metric:  "my-test-metric",
				Value:   20.0,
				Success: true,
				Message: "",
			},
			InCriteriaSets: []*apimodelsv2.SLOCriteria{
				{
					Criteria: []string{"<=10.0"},
				},
				{
					Criteria: []string{"<=+10%"},
				},
			},
			InPreviousResults: []*keptnv2.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
				{
					Score: 2,
					Value: &keptnv2.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					PassTargets:    nil,
					WarningTargets: nil,
					KeySLI:         false,
					Status:         "pass",
				},
			},
			InComparison: &apimodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			ExpectedTargets: []*keptnv2.SLITarget{
				{
					Criteria:    "<=10.0",
					TargetValue: 10.0,
					Violated:    true,
				},
				{
					Criteria:    "<=+10%",
					TargetValue: 11.0,
					Violated:    true,
				},
			},
			ExpectedResult: false,
			ExpectedError:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Run(test.Name, func(t *testing.T) {
				result, violations, err := evaluateOrCombinedCriteria(test.InSLIResult, test.InCriteriaSets, test.InPreviousResults, test.InComparison)
				assert.EqualValues(t, test.ExpectedResult, result)
				assert.EqualValues(t, test.ExpectedTargets, violations)
				assert.EqualValues(t, test.ExpectedError, err)
			})
		})
	}
}

type evaluateObjectivesTestObject struct {
	Name                       string
	InGetSLIDoneEvent          *keptnv2.GetSLIFinishedEventData
	InSLOConfig                *apimodelsv2.ServiceLevelObjectives
	InPreviousEvaluationEvents []*keptnv2.EvaluationFinishedEventData
	ExpectedEvaluationResult   *keptnv2.EvaluationFinishedEventData
	ExpectedMaximumScore       float64
	ExpectedKeySLIFailed       bool
}

func TestEvaluateObjectives(t *testing.T) {
	tests := []*evaluateObjectivesTestObject{
		{
			Name: "Simple comparison evaluation",
			InGetSLIDoneEvent: &keptnv2.GetSLIFinishedEventData{
				EventData: keptnv2.EventData{
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
				},
				GetSLI: keptnv2.GetSLIFinished{
					Start: "2019-10-20T07:57:27.152330783Z",
					End:   "2019-10-22T08:57:27.152330783Z",
					IndicatorValues: []*keptnv2.SLIResult{
						{
							Metric:  "my-test-metric-1",
							Value:   10.0,
							Success: true,
							Message: "",
						},
					},
				},
			},
			InSLOConfig: &apimodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &apimodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
				Objectives: []*apimodelsv2.SLO{
					{
						SLI: "my-test-metric-1",
						Pass: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=15.0"},
							},
							{
								Criteria: []string{"<=+10%"},
							},
						},
						Warning: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=20.0"},
							},
							{
								Criteria: []string{"<=+15%"},
							},
						},
						Weight: 1,
						KeySLI: false,
					},
				},
				TotalScore: &apimodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InPreviousEvaluationEvents: []*keptnv2.EvaluationFinishedEventData{
				{
					Evaluation: keptnv2.EvaluationDetails{
						TimeStart: "",
						TimeEnd:   "",
						Result:    "pass",
						Score:     2,
						IndicatorResults: []*keptnv2.SLIEvaluationResult{
							{
								Score: 2,
								Value: &keptnv2.SLIResult{
									Metric:  "my-test-metric-1",
									Value:   10.0,
									Success: true,
									Message: "",
								},
								PassTargets:    nil,
								WarningTargets: nil,
								KeySLI:         false,
								Status:         "pass",
							},
						},
					},
					EventData: keptnv2.EventData{
						Result:  "pass",
						Project: "sockshop",
						Service: "carts",
						Stage:   "dev",
					},
				},
			},
			ExpectedEvaluationResult: &keptnv2.EvaluationFinishedEventData{
				Evaluation: keptnv2.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "", // not set by the tested function
					Score:     0,  // not calculated by tested function
					IndicatorResults: []*keptnv2.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnv2.SLIResult{
								Metric:        "my-test-metric-1",
								Value:         10.0,
								ComparedValue: 10.0,
								Success:       true,
								Message:       "",
							},
							WarningTargets: []*keptnv2.SLITarget{
								{
									Criteria:    "<=20.0",
									TargetValue: 20,
									Violated:    false,
								},
								{
									Criteria:    "<=+15%",
									TargetValue: 11.5,
									Violated:    false,
								},
							},
							PassTargets: []*keptnv2.SLITarget{
								{
									Criteria:    "<=15.0",
									TargetValue: 15.0,
									Violated:    false,
								},
								{
									Criteria:    "<=+10%",
									TargetValue: 11,
									Violated:    false,
								},
							},
							KeySLI: false,
							Status: "pass",
						},
					},
				},
				EventData: keptnv2.EventData{
					Result:  "",
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
				},
			},
			ExpectedMaximumScore: 1,
			ExpectedKeySLIFailed: false,
		},
		{
			Name: "Expect Warning",
			InGetSLIDoneEvent: &keptnv2.GetSLIFinishedEventData{
				EventData: keptnv2.EventData{
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
				},
				GetSLI: keptnv2.GetSLIFinished{
					Start: "2019-10-20T07:57:27.152330783Z",
					End:   "2019-10-22T08:57:27.152330783Z",
					IndicatorValues: []*keptnv2.SLIResult{
						{
							Metric:  "my-test-metric-1",
							Value:   16.0,
							Success: true,
							Message: "",
						},
					},
				},
			},
			InSLOConfig: &apimodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &apimodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
				Objectives: []*apimodelsv2.SLO{
					{
						SLI: "my-test-metric-1",
						Pass: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=15.0"},
							},
							{
								Criteria: []string{"<=+10%"},
							},
						},
						Warning: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=20.0"},
							},
							{
								Criteria: []string{"<=+15%"},
							},
						},
						Weight: 1,
						KeySLI: false,
					},
				},
				TotalScore: &apimodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InPreviousEvaluationEvents: []*keptnv2.EvaluationFinishedEventData{
				{
					Evaluation: keptnv2.EvaluationDetails{
						TimeStart: "",
						TimeEnd:   "",
						Result:    "pass",
						Score:     2,
						IndicatorResults: []*keptnv2.SLIEvaluationResult{
							{
								Score: 2,
								Value: &keptnv2.SLIResult{
									Metric:  "my-test-metric-1",
									Value:   10.0,
									Success: true,
									Message: "",
								},
								PassTargets:    nil,
								WarningTargets: nil,
								KeySLI:         false,
								Status:         "pass",
							},
						},
					},
					EventData: keptnv2.EventData{
						Result:  "pass",
						Project: "sockshop",
						Service: "carts",
						Stage:   "dev",
					},
				},
			},
			ExpectedEvaluationResult: &keptnv2.EvaluationFinishedEventData{
				Evaluation: keptnv2.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "", // not set by the tested function
					Score:     0,  // not calculated by tested function
					IndicatorResults: []*keptnv2.SLIEvaluationResult{
						{
							Score: 0.5,
							Value: &keptnv2.SLIResult{
								Metric:        "my-test-metric-1",
								Value:         16.0,
								ComparedValue: 10.0,
								Success:       true,
								Message:       "",
							},
							WarningTargets: []*keptnv2.SLITarget{
								{
									Criteria:    "<=20.0",
									TargetValue: 20,
									Violated:    false,
								},
								{
									Criteria:    "<=+15%",
									TargetValue: 11.5,
									Violated:    true,
								},
							},
							PassTargets: []*keptnv2.SLITarget{
								{
									Criteria:    "<=15.0",
									TargetValue: 15.0,
									Violated:    true,
								},
								{
									Criteria:    "<=+10%",
									TargetValue: 11,
									Violated:    true,
								},
							},
							KeySLI: false,
							Status: "warning",
						},
					},
				},
				EventData: keptnv2.EventData{
					Result:  "",
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
				},
			},
			ExpectedMaximumScore: 1,
			ExpectedKeySLIFailed: false,
		},
		{
			Name: "Logging SLI with no pass criteria should not affect score",
			InGetSLIDoneEvent: &keptnv2.GetSLIFinishedEventData{
				EventData: keptnv2.EventData{
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
				},
				GetSLI: keptnv2.GetSLIFinished{
					Start: "2019-10-20T07:57:27.152330783Z",
					End:   "2019-10-22T08:57:27.152330783Z",
					IndicatorValues: []*keptnv2.SLIResult{
						{
							Metric:  "my-test-metric-1",
							Value:   10.0,
							Success: true,
							Message: "",
						},
						{
							Metric:  "my-log-metric",
							Value:   30.0,
							Success: true,
							Message: "",
						},
					},
				},
			},
			InSLOConfig: &apimodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &apimodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
				Objectives: []*apimodelsv2.SLO{
					{
						SLI: "my-test-metric-1",
						Pass: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=15.0"},
							},
							{
								Criteria: []string{"<=+10%"},
							},
						},
						Warning: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=20.0"},
							},
							{
								Criteria: []string{"<=+15%"},
							},
						},
						Weight: 1,
						KeySLI: false,
					},
					{
						SLI:    "my-log-metric",
						Weight: 1,
						KeySLI: false,
					},
				},
				TotalScore: &apimodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InPreviousEvaluationEvents: []*keptnv2.EvaluationFinishedEventData{
				{
					Evaluation: keptnv2.EvaluationDetails{
						TimeStart: "",
						TimeEnd:   "",
						Result:    "pass",
						Score:     2,
						IndicatorResults: []*keptnv2.SLIEvaluationResult{
							{
								Score: 2,
								Value: &keptnv2.SLIResult{
									Metric:  "my-test-metric-1",
									Value:   10.0,
									Success: true,
									Message: "",
								},
								PassTargets:    nil,
								WarningTargets: nil,
								KeySLI:         false,
								Status:         "pass",
							},
						},
					},
					EventData: keptnv2.EventData{
						Result:  "pass",
						Project: "sockshop",
						Service: "carts",
						Stage:   "dev",
					},
				},
			},
			ExpectedEvaluationResult: &keptnv2.EvaluationFinishedEventData{
				Evaluation: keptnv2.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "", // not set by the tested function
					Score:     0,  // not calculated by tested function
					IndicatorResults: []*keptnv2.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnv2.SLIResult{
								Metric:        "my-test-metric-1",
								Value:         10.0,
								ComparedValue: 10.0,
								Success:       true,
								Message:       "",
							},
							WarningTargets: []*keptnv2.SLITarget{
								{
									Criteria:    "<=20.0",
									TargetValue: 20,
									Violated:    false,
								},
								{
									Criteria:    "<=+15%",
									TargetValue: 11.5,
									Violated:    false,
								},
							},
							PassTargets: []*keptnv2.SLITarget{
								{
									Criteria:    "<=15.0",
									TargetValue: 15.0,
									Violated:    false,
								},
								{
									Criteria:    "<=+10%",
									TargetValue: 11,
									Violated:    false,
								},
							},
							KeySLI: false,
							Status: "pass",
						},
						{
							Score: 0,
							Value: &keptnv2.SLIResult{
								Metric:  "my-log-metric",
								Value:   30.0,
								Success: true,
								Message: "",
							},
							Status: "info",
						},
					},
				},
				EventData: keptnv2.EventData{
					Result:  "",
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
				},
			},
			ExpectedMaximumScore: 1,
			ExpectedKeySLIFailed: false,
		},
		{
			Name: "Logging SLI with empty pass criteria array should not affect score and have status 'info' - BUG 2231",
			InGetSLIDoneEvent: &keptnv2.GetSLIFinishedEventData{
				EventData: keptnv2.EventData{
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
				},
				GetSLI: keptnv2.GetSLIFinished{
					Start: "2019-10-20T07:57:27.152330783Z",
					End:   "2019-10-22T08:57:27.152330783Z",
					IndicatorValues: []*keptnv2.SLIResult{
						{
							Metric:  "my-test-metric-1",
							Value:   10.0,
							Success: true,
							Message: "",
						},
						{
							Metric:  "my-log-metric",
							Value:   30.0,
							Success: true,
							Message: "",
						},
					},
				},
			},
			InSLOConfig: &apimodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &apimodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
				Objectives: []*apimodelsv2.SLO{
					{
						SLI: "my-test-metric-1",
						Pass: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=15.0"},
							},
							{
								Criteria: []string{"<=+10%"},
							},
						},
						Warning: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=20.0"},
							},
							{
								Criteria: []string{"<=+15%"},
							},
						},
						Weight: 1,
						KeySLI: false,
					},
					{
						SLI:     "my-log-metric",
						Weight:  1,
						KeySLI:  false,
						Pass:    []*apimodelsv2.SLOCriteria{},
						Warning: []*apimodelsv2.SLOCriteria{},
					},
				},
				TotalScore: &apimodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InPreviousEvaluationEvents: []*keptnv2.EvaluationFinishedEventData{
				{
					Evaluation: keptnv2.EvaluationDetails{
						TimeStart: "",
						TimeEnd:   "",
						Result:    "pass",
						Score:     2,
						IndicatorResults: []*keptnv2.SLIEvaluationResult{
							{
								Score: 2,
								Value: &keptnv2.SLIResult{
									Metric:  "my-test-metric-1",
									Value:   10.0,
									Success: true,
									Message: "",
								},
								PassTargets:    nil,
								WarningTargets: nil,
								KeySLI:         false,
								Status:         "pass",
							},
						},
					},
					EventData: keptnv2.EventData{
						Result:  "pass",
						Project: "sockshop",
						Service: "carts",
						Stage:   "dev",
					},
				},
			},
			ExpectedEvaluationResult: &keptnv2.EvaluationFinishedEventData{
				Evaluation: keptnv2.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "", // not set by the tested function
					Score:     0,  // not calculated by tested function
					IndicatorResults: []*keptnv2.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnv2.SLIResult{
								Metric:        "my-test-metric-1",
								Value:         10.0,
								ComparedValue: 10.0,
								Success:       true,
								Message:       "",
							},
							WarningTargets: []*keptnv2.SLITarget{
								{
									Criteria:    "<=20.0",
									TargetValue: 20,
									Violated:    false,
								},
								{
									Criteria:    "<=+15%",
									TargetValue: 11.5,
									Violated:    false,
								},
							},
							PassTargets: []*keptnv2.SLITarget{
								{
									Criteria:    "<=15.0",
									TargetValue: 15.0,
									Violated:    false,
								},
								{
									Criteria:    "<=+10%",
									TargetValue: 11,
									Violated:    false,
								},
							},
							KeySLI: false,
							Status: "pass",
						},
						{
							Score: 0,
							Value: &keptnv2.SLIResult{
								Metric:        "my-log-metric",
								Value:         30.0,
								ComparedValue: 0.0,
								Success:       true,
								Message:       "",
							},
							Status: "info",
						},
					},
				},
				EventData: keptnv2.EventData{
					Result:  "",
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
				},
			},
			ExpectedMaximumScore: 1,
			ExpectedKeySLIFailed: false,
		},
		{
			Name: "BUG 1125",
			InGetSLIDoneEvent: &keptnv2.GetSLIFinishedEventData{
				EventData: keptnv2.EventData{
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
				},
				GetSLI: keptnv2.GetSLIFinished{
					Start: "2019-10-20T07:57:27.152330783Z",
					End:   "2019-10-22T08:57:27.152330783Z",
					IndicatorValues: []*keptnv2.SLIResult{
						{
							Metric:  "response_time_p50",
							Value:   1011.0745528937252,
							Success: true,
							Message: "",
						},
					},
				},
			},
			InSLOConfig: &apimodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &apimodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
				Objectives: []*apimodelsv2.SLO{
					{
						SLI: "response_time_p50",
						Pass: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=+20%", "<500"},
							},
						},
						Weight: 1,
						KeySLI: false,
					},
				},
				TotalScore: &apimodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InPreviousEvaluationEvents: nil,
			ExpectedEvaluationResult: &keptnv2.EvaluationFinishedEventData{
				Evaluation: keptnv2.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "", // not set by the tested function
					Score:     0,  // not calculated by tested function
					IndicatorResults: []*keptnv2.SLIEvaluationResult{
						{
							Score: 0,
							Value: &keptnv2.SLIResult{
								Metric:  "response_time_p50",
								Value:   1011.0745528937252,
								Success: true,
								Message: "",
							},
							WarningTargets: nil,
							PassTargets: []*keptnv2.SLITarget{
								{
									Criteria:    "<=+20%",
									TargetValue: 0,
									Violated:    false,
								},
								{
									Criteria:    "<500",
									TargetValue: 500,
									Violated:    true,
								},
							},
							KeySLI: false,
							Status: "fail",
						},
					},
				},
				EventData: keptnv2.EventData{
					Result:  "",
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
				},
			},
			ExpectedMaximumScore: 1,
			ExpectedKeySLIFailed: false,
		},
		{
			Name: "BUG 1263",
			InGetSLIDoneEvent: &keptnv2.GetSLIFinishedEventData{
				EventData: keptnv2.EventData{
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
				},
				GetSLI: keptnv2.GetSLIFinished{
					Start: "2019-10-20T07:57:27.152330783Z",
					End:   "2019-10-22T08:57:27.152330783Z",
					IndicatorValues: []*keptnv2.SLIResult{
						{
							Metric:  "response_time_p50",
							Value:   100,
							Success: true,
							Message: "",
						},
					},
				},
			},
			InSLOConfig: &apimodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &apimodelsv2.SLOComparison{
					CompareWith:               "single_result",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 1,
					AggregateFunction:         "avg",
				},
				Objectives: []*apimodelsv2.SLO{
					{
						SLI: "response_time_p50",
						Pass: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=+20%"},
							},
						},
						Weight: 1,
						KeySLI: false,
					},
				},
				TotalScore: &apimodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InPreviousEvaluationEvents: []*keptnv2.EvaluationFinishedEventData{
				{
					Evaluation: keptnv2.EvaluationDetails{
						TimeStart: "",
						TimeEnd:   "",
						Result:    "fail",
						Score:     0,
						IndicatorResults: []*keptnv2.SLIEvaluationResult{
							{
								Score: 0,
								Value: &keptnv2.SLIResult{
									Metric:  "response_time_p50",
									Value:   0.0,
									Success: false,
									Message: "",
								},
								PassTargets:    nil,
								WarningTargets: nil,
								KeySLI:         false,
								Status:         "fail",
							},
						},
					},
					EventData: keptnv2.EventData{
						Result:  "pass",
						Project: "sockshop",
						Service: "carts",
						Stage:   "dev",
					},
				},
			},
			ExpectedEvaluationResult: &keptnv2.EvaluationFinishedEventData{
				Evaluation: keptnv2.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "", // not set by the tested function
					Score:     0,  // not calculated by tested function
					IndicatorResults: []*keptnv2.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnv2.SLIResult{
								Metric:  "response_time_p50",
								Value:   100,
								Success: true,
								Message: "",
							},
							WarningTargets: nil,
							PassTargets: []*keptnv2.SLITarget{
								{
									Criteria:    "<=+20%",
									TargetValue: 0,
									Violated:    false,
								},
							},
							KeySLI: false,
							Status: "pass",
						},
					},
				},
				EventData: keptnv2.EventData{
					Result:  "",
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
				},
			},
			ExpectedMaximumScore: 1,
			ExpectedKeySLIFailed: false,
		},
		{
			Name: "6096 if SLI does not have objective have a message",
			InGetSLIDoneEvent: &keptnv2.GetSLIFinishedEventData{
				EventData: keptnv2.EventData{
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
					Result:  "fail",
				},
				GetSLI: keptnv2.GetSLIFinished{
					Start: "2019-10-20T07:57:27.152330783Z",
					End:   "2019-10-22T08:57:27.152330783Z",
					IndicatorValues: []*keptnv2.SLIResult{
						{
							Metric:  "my-test-metric-1",
							Value:   10.0,
							Success: true,
							Message: "",
						},
					},
				},
			},
			InSLOConfig: &apimodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &apimodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
				Objectives: []*apimodelsv2.SLO{
					{
						SLI: "a_different_metric",
						Pass: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=15.0"},
							},
							{
								Criteria: []string{"<=+10%"},
							},
						},
						Warning: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=20.0"},
							},
							{
								Criteria: []string{"<=+15%"},
							},
						},
						Weight: 1,
						KeySLI: false,
					},
				},
				TotalScore: &apimodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InPreviousEvaluationEvents: []*keptnv2.EvaluationFinishedEventData{
				{
					Evaluation: keptnv2.EvaluationDetails{
						TimeStart: "",
						TimeEnd:   "",
						Result:    "pass",
						Score:     2,
						IndicatorResults: []*keptnv2.SLIEvaluationResult{
							{
								Score: 2,
								Value: &keptnv2.SLIResult{
									Metric:  "my-test-metric-1",
									Value:   10.0,
									Success: true,
									Message: "",
								},
								PassTargets:    nil,
								WarningTargets: nil,
								KeySLI:         false,
								Status:         "pass",
							},
							{
								Score: 2,
								Value: &keptnv2.SLIResult{
									Metric:  "a_different_metric",
									Value:   5.0,
									Success: true,
									Message: "",
								},
								PassTargets:    nil,
								WarningTargets: nil,
								KeySLI:         false,
								Status:         "pass",
							},
						},
					},
					EventData: keptnv2.EventData{
						Result:  "pass",
						Project: "sockshop",
						Service: "carts",
						Stage:   "dev",
					},
				},
				{
					Evaluation: keptnv2.EvaluationDetails{
						TimeStart: "",
						TimeEnd:   "",
						Result:    "pass",
						Score:     2,
						IndicatorResults: []*keptnv2.SLIEvaluationResult{
							{
								Score: 2,
								Value: &keptnv2.SLIResult{
									Metric:  "my-test-metric-1",
									Value:   10.0,
									Success: true,
									Message: "",
								},
								PassTargets:    nil,
								WarningTargets: nil,
								KeySLI:         false,
								Status:         "pass",
							},
							{
								Score:          2,
								Value:          nil,
								PassTargets:    nil,
								WarningTargets: nil,
								KeySLI:         false,
								Status:         "pass",
							},
						},
					},
					EventData: keptnv2.EventData{
						Result:  "pass",
						Project: "sockshop",
						Service: "carts",
						Stage:   "dev",
					},
				},
			},
			ExpectedEvaluationResult: &keptnv2.EvaluationFinishedEventData{
				Evaluation: keptnv2.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "", // not set by the tested function
					Score:     0,  // not calculated by tested function
					IndicatorResults: []*keptnv2.SLIEvaluationResult{
						{
							Score: 0,
							Value: &keptnv2.SLIResult{
								Metric:  "a_different_metric",
								Value:   0,
								Success: false,
								Message: "no value received from SLI provider",
							},
							PassTargets: []*keptnv2.SLITarget{
								{
									Criteria:    "<=15.0",
									TargetValue: 15,
									Violated:    true,
								},
								{
									Criteria:    "<=+10%",
									TargetValue: 5,
									Violated:    true,
								},
							},
							WarningTargets: []*keptnv2.SLITarget{
								{
									Criteria:    "<=20.0",
									TargetValue: 20,
									Violated:    true,
								},
								{
									Criteria:    "<=+15%",
									TargetValue: 5,
									Violated:    true,
								},
							},
							KeySLI: false,
							Status: "fail",
						},
					},
				},
				EventData: keptnv2.EventData{
					Result:  "",
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
					Message: "Lighthouse received additional SLIs, which are not specified as SLO: my-test-metric-1 . Please consider using them as an SLO.",
				},
			},
			ExpectedMaximumScore: 1,
			ExpectedKeySLIFailed: false,
		},
		{
			Name: "9198 SLO file does not have objectives",
			InGetSLIDoneEvent: &keptnv2.GetSLIFinishedEventData{
				EventData: keptnv2.EventData{
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
					Result:  "fail",
					Status:  "succeeded",
				},
				GetSLI: keptnv2.GetSLIFinished{
					Start: "2019-10-20T07:57:27.152330783Z",
					End:   "2019-10-22T08:57:27.152330783Z",
					IndicatorValues: []*keptnv2.SLIResult{
						{
							Metric:  "no metric",
							Value:   0,
							Success: false,
							Message: "no SLIs were requested",
						},
					},
				},
			},
			InSLOConfig: &apimodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &apimodelsv2.SLOComparison{
					CompareWith:               "single_result",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 1,
					AggregateFunction:         "avg",
				},
				Objectives: []*apimodelsv2.SLO{},
				TotalScore: &apimodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InPreviousEvaluationEvents: []*keptnv2.EvaluationFinishedEventData{},
			ExpectedEvaluationResult: &keptnv2.EvaluationFinishedEventData{
				Evaluation: keptnv2.EvaluationDetails{
					TimeStart:        "2019-10-20T07:57:27.152330783Z",
					TimeEnd:          "2019-10-22T08:57:27.152330783Z",
					Result:           "", // not set by the tested function
					Score:            0,  // not calculated by tested function
					IndicatorResults: nil,
				},
				EventData: keptnv2.EventData{
					Result:  "fail",
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
					Message: "lighthouse failed because SLI failed with message no SLIs were requested",
				},
			},
			ExpectedMaximumScore: 1,
			ExpectedKeySLIFailed: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			evaluationDoneData, maximumScore, keySLIFailed := evaluateObjectives(test.InGetSLIDoneEvent, test.InSLOConfig, test.InPreviousEvaluationEvents)
			assert.EqualValues(t, test.ExpectedEvaluationResult, evaluationDoneData)
			assert.EqualValues(t, test.ExpectedMaximumScore, maximumScore)
			assert.EqualValues(t, test.ExpectedKeySLIFailed, keySLIFailed)
		})
	}
}

type calculateScoreTestObject struct {
	Name                     string
	InMaximumScore           float64
	InEvaluationResult       *keptnv2.EvaluationFinishedEventData
	InSLOConfig              *apimodelsv2.ServiceLevelObjectives
	InKeySLIFailed           bool
	ExpectedEvaluationResult *keptnv2.EvaluationFinishedEventData
	ExpectedError            error
}

func TestCalculateScore(t *testing.T) {
	tests := []*calculateScoreTestObject{
		{
			Name:           "Simple comparison",
			InMaximumScore: 1,
			InEvaluationResult: &keptnv2.EvaluationFinishedEventData{
				Evaluation: keptnv2.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "", // to be calculated
					Score:     0,  // to be calculated
					IndicatorResults: []*keptnv2.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnv2.SLIResult{
								Metric:  "my-test-metric-1",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							PassTargets:    nil,
							WarningTargets: nil,
							KeySLI:         false,
							Status:         "pass",
						},
					},
				},
				EventData: keptnv2.EventData{
					Result:  "",
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
				},
			},
			InSLOConfig: &apimodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &apimodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
				Objectives: []*apimodelsv2.SLO{
					{
						SLI: "my-test-metric-1",
						Pass: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=15.0"},
							},
							{
								Criteria: []string{"<=+10%"},
							},
						},
						Warning: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=20.0"},
							},
							{
								Criteria: []string{"<=+15%"},
							},
						},
						Weight: 1,
						KeySLI: false,
					},
				},
				TotalScore: &apimodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InKeySLIFailed: false,
			ExpectedEvaluationResult: &keptnv2.EvaluationFinishedEventData{
				Evaluation: keptnv2.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "pass",
					Score:     100.0,
					IndicatorResults: []*keptnv2.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnv2.SLIResult{
								Metric:  "my-test-metric-1",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							PassTargets:    nil,
							WarningTargets: nil,
							KeySLI:         false,
							Status:         "pass",
						},
					},
				},
				EventData: keptnv2.EventData{
					Result:  "pass",
					Status:  "succeeded",
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
				},
			},
			ExpectedError: nil,
		},
		{
			Name:           "Key SLI failed",
			InMaximumScore: 2,
			InEvaluationResult: &keptnv2.EvaluationFinishedEventData{
				Evaluation: keptnv2.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "", // to be calculated
					Score:     0,  // to be calculated
					IndicatorResults: []*keptnv2.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnv2.SLIResult{
								Metric:  "my-test-metric-1",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							PassTargets:    nil,
							WarningTargets: nil,
							KeySLI:         false,
							Status:         "pass",
						},
						{
							Score: 0,
							Value: &keptnv2.SLIResult{
								Metric:  "my-key-metric",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							PassTargets:    nil,
							WarningTargets: nil,
							KeySLI:         false,
							Status:         "fail",
						},
					},
				},
				EventData: keptnv2.EventData{
					Result:  "",
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
				},
			},
			InSLOConfig: &apimodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &apimodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
				Objectives: []*apimodelsv2.SLO{
					{
						SLI: "my-test-metric-1",
						Pass: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=15.0"},
							},
							{
								Criteria: []string{"<=+10%"},
							},
						},
						Warning: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=20.0"},
							},
							{
								Criteria: []string{"<=+15%"},
							},
						},
						Weight: 1,
						KeySLI: false,
					},
					{
						SLI: "my-key-metric",
						Pass: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=15.0"},
							},
							{
								Criteria: []string{"<=+10%"},
							},
						},
						Warning: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=20.0"},
							},
							{
								Criteria: []string{"<=+15%"},
							},
						},
						Weight: 1,
						KeySLI: true,
					},
				},
				TotalScore: &apimodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InKeySLIFailed: true,
			ExpectedEvaluationResult: &keptnv2.EvaluationFinishedEventData{
				Evaluation: keptnv2.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "fail",
					Score:     50.0,
					IndicatorResults: []*keptnv2.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnv2.SLIResult{
								Metric:  "my-test-metric-1",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							PassTargets:    nil,
							WarningTargets: nil,
							KeySLI:         false,
							Status:         "pass",
						},
						{
							Score: 0,
							Value: &keptnv2.SLIResult{
								Metric:  "my-key-metric",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							PassTargets:    nil,
							WarningTargets: nil,
							KeySLI:         false,
							Status:         "fail",
						},
					},
				},
				EventData: keptnv2.EventData{
					Result:  "fail",
					Status:  "succeeded",
					Labels:  nil,
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
					Message: "Evaluation failed since the calculated score of 50 is below the target value of 90",
				},
			},
			ExpectedError: nil,
		}, {
			Name:           "Non-Key SLI warning",
			InMaximumScore: 2,
			InEvaluationResult: &keptnv2.EvaluationFinishedEventData{
				Evaluation: keptnv2.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "", // to be calculated
					Score:     0,  // to be calculated
					IndicatorResults: []*keptnv2.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnv2.SLIResult{
								Metric:  "my-test-metric-1",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							PassTargets:    nil,
							WarningTargets: nil,
							KeySLI:         false,
							Status:         "pass",
						},
						{
							Score: 0.506,
							Value: &keptnv2.SLIResult{
								Metric:  "my-metric",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							PassTargets:    nil,
							WarningTargets: nil,
							KeySLI:         false,
							Status:         "pass",
						},
					},
				},
				EventData: keptnv2.EventData{
					Result:  "",
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
				},
			},
			InSLOConfig: &apimodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &apimodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
				Objectives: []*apimodelsv2.SLO{
					{
						SLI: "my-test-metric-1",
						Pass: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=8.0"},
							},
							{
								Criteria: []string{"<=+10%"},
							},
						},
						Warning: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=13.0"},
							},
							{
								Criteria: []string{"<=+15%"},
							},
						},
						Weight: 1,
						KeySLI: false,
					},
					{
						SLI: "my-metric",
						Pass: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=8.0"},
							},
							{
								Criteria: []string{"<=+10%"},
							},
						},
						Warning: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=12.0"},
							},
							{
								Criteria: []string{"<=+15%"},
							},
						},
						Weight: 1,
						KeySLI: false,
					},
				},
				TotalScore: &apimodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InKeySLIFailed: false,
			ExpectedEvaluationResult: &keptnv2.EvaluationFinishedEventData{
				Evaluation: keptnv2.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "warning",
					Score:     75.3,
					IndicatorResults: []*keptnv2.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnv2.SLIResult{
								Metric:  "my-test-metric-1",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							PassTargets:    nil,
							WarningTargets: nil,
							KeySLI:         false,
							Status:         "pass",
						},
						{
							Score: 0.506,
							Value: &keptnv2.SLIResult{
								Metric:  "my-metric",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							PassTargets:    nil,
							WarningTargets: nil,
							KeySLI:         false,
							Status:         "pass",
						},
					},
				},
				EventData: keptnv2.EventData{
					Result:  "warning",
					Status:  "succeeded",
					Labels:  nil,
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
					Message: "Evaluation returned a warning: the calculated score of 75.3 is close to the warning target value of 75",
				},
			},
			ExpectedError: nil,
		}, {
			Name:           "Non-Key SLI fail",
			InMaximumScore: 2,
			InEvaluationResult: &keptnv2.EvaluationFinishedEventData{
				Evaluation: keptnv2.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "", // to be calculated
					Score:     0,  // to be calculated
					IndicatorResults: []*keptnv2.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnv2.SLIResult{
								Metric:  "my-test-metric-1",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							PassTargets:    nil,
							WarningTargets: nil,
							KeySLI:         false,
							Status:         "pass",
						},
						{
							Score: 0.48,
							Value: &keptnv2.SLIResult{
								Metric:  "my-metric",
								Value:   10.0,
								Success: false,
								Message: "",
							},
							PassTargets:    nil,
							WarningTargets: nil,
							KeySLI:         false,
							Status:         "fail",
						},
					},
				},
				EventData: keptnv2.EventData{
					Result:  "",
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
				},
			},
			InSLOConfig: &apimodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &apimodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
				Objectives: []*apimodelsv2.SLO{
					{
						SLI: "my-test-metric-1",
						Pass: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=8.0"},
							},
							{
								Criteria: []string{"<=+10%"},
							},
						},
						Warning: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=10.0"},
							},
							{
								Criteria: []string{"<=+15%"},
							},
						},
						Weight: 1,
						KeySLI: false,
					},
					{
						SLI: "my-metric",
						Pass: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=8.0"},
							},
							{
								Criteria: []string{"<=+10%"},
							},
						},
						Warning: []*apimodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=10.0"},
							},
							{
								Criteria: []string{"<=+15%"},
							},
						},
						Weight: 1,
						KeySLI: false,
					},
				},
				TotalScore: &apimodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InKeySLIFailed: false,
			ExpectedEvaluationResult: &keptnv2.EvaluationFinishedEventData{
				Evaluation: keptnv2.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "fail",
					Score:     74,
					IndicatorResults: []*keptnv2.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnv2.SLIResult{
								Metric:  "my-test-metric-1",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							PassTargets:    nil,
							WarningTargets: nil,
							KeySLI:         false,
							Status:         "pass",
						},
						{
							Score: 0.48,
							Value: &keptnv2.SLIResult{
								Metric:  "my-metric",
								Value:   10.0,
								Success: false,
								Message: "",
							},
							PassTargets:    nil,
							WarningTargets: nil,
							KeySLI:         false,
							Status:         "fail",
						},
					},
				},
				EventData: keptnv2.EventData{
					Result:  "fail",
					Status:  "succeeded",
					Labels:  nil,
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
					Message: "Evaluation failed since the calculated score of 74 is below the warning value of 75",
				},
			},
			ExpectedError: nil,
		},
		{
			Name:           "Only Info SLIs",
			InMaximumScore: 0,
			InEvaluationResult: &keptnv2.EvaluationFinishedEventData{
				Evaluation: keptnv2.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "", // to be calculated
					Score:     0,  // to be calculated
					IndicatorResults: []*keptnv2.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnv2.SLIResult{
								Metric:  "my-test-metric-1",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							PassTargets:    nil,
							WarningTargets: nil,
							KeySLI:         false,
							Status:         "pass",
						},
						{
							Score: 1,
							Value: &keptnv2.SLIResult{
								Metric:  "my-key-metric",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							PassTargets:    nil,
							WarningTargets: nil,
							KeySLI:         false,
							Status:         "pass",
						},
					},
				},
				EventData: keptnv2.EventData{
					Result:  "",
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
				},
			},
			InSLOConfig: &apimodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &apimodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
				Objectives: []*apimodelsv2.SLO{
					{
						SLI:    "my-test-metric-1",
						Weight: 1,
						KeySLI: false,
					},
					{
						SLI:    "my-key-metric",
						Weight: 1,
						KeySLI: false,
					},
				},
				TotalScore: &apimodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InKeySLIFailed: true,
			ExpectedEvaluationResult: &keptnv2.EvaluationFinishedEventData{
				Evaluation: keptnv2.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "pass",
					Score:     100.0,
					IndicatorResults: []*keptnv2.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnv2.SLIResult{
								Metric:        "my-test-metric-1",
								Value:         10.0,
								ComparedValue: 0.0,
								Success:       true,
								Message:       "",
							},
							PassTargets:    nil,
							WarningTargets: nil,
							KeySLI:         false,
							Status:         "pass",
						},
						{
							Score: 1,
							Value: &keptnv2.SLIResult{
								Metric:        "my-key-metric",
								Value:         10.0,
								ComparedValue: 0.0,
								Success:       true,
								Message:       "",
							},
							PassTargets:    nil,
							WarningTargets: nil,
							KeySLI:         false,
							Status:         "pass",
						},
					},
				},
				EventData: keptnv2.EventData{
					Result:  "pass",
					Status:  "succeeded",
					Labels:  nil,
					Project: "sockshop",
					Service: "carts",
					Stage:   "dev",
				},
			},
			ExpectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := calculateScore(test.InMaximumScore, test.InEvaluationResult, test.InSLOConfig, test.InKeySLIFailed)

			assert.EqualValues(t, test.ExpectedError, err)
			assert.EqualValues(t, test.ExpectedEvaluationResult, test.InEvaluationResult)
		})
	}
}

func TestEvaluateSLIHandler_getPreviousEvaluations(t *testing.T) {

	var returnedResult datastoreResult

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)

			marshal, _ := json.Marshal(&returnedResult)
			w.Write(marshal)
		}),
	)
	defer ts.Close()

	t.Setenv("MONGODB_DATASTORE", strings.TrimPrefix(ts.URL, "http://"))

	type fields struct {
		Logger     *keptncommon.Logger
		Event      cloudevents.Event
		HTTPClient *http.Client
	}
	type args struct {
		e                       *keptnv2.GetSLIFinishedEventData
		numberOfPreviousResults int
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		resultFromDatastore datastoreResult
		want                []*keptnv2.EvaluationFinishedEventData
		want2               []string
		wantErr             bool
	}{
		{
			name: "get evaluation-done events",
			fields: fields{
				Logger:     nil,
				Event:      cloudevents.Event{},
				HTTPClient: &http.Client{},
			},
			args: args{
				e: &keptnv2.GetSLIFinishedEventData{
					EventData: keptnv2.EventData{
						Project: "sockshop",
						Stage:   "dev",
						Service: "carts",
					},
				},
				numberOfPreviousResults: 1,
			},
			resultFromDatastore: datastoreResult{
				NextPageKey: "",
				TotalCount:  1,
				PageSize:    1,
				Events: []struct {
					Data interface{} `json:"data"`
					ID   string      `json:"id"`
				}{
					{
						Data: &keptnv2.EvaluationFinishedEventData{
							EventData: keptnv2.EventData{
								Project: "sockshop",
								Service: "carts",
								Stage:   "dev",
								Labels:  nil,
								Result:  "",
							},
						},
						ID: "my-id",
					},
				},
			},
			want: []*keptnv2.EvaluationFinishedEventData{
				{
					EventData: keptnv2.EventData{
						Project: "sockshop",
						Service: "carts",
						Stage:   "dev",
						Labels:  nil,
						Result:  "",
					},
				},
			},
			want2:   []string{"my-id"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		returnedResult = tt.resultFromDatastore
		t.Run(tt.name, func(t *testing.T) {

			eh := &EvaluateSLIHandler{
				KeptnHandler: nil,
				Event:        tt.fields.Event,
				HTTPClient:   tt.fields.HTTPClient,
			}
			got, got2, err := eh.getPreviousEvaluations(tt.args.e, tt.args.numberOfPreviousResults, "all")
			if (err != nil) != tt.wantErr {
				t.Errorf("getPreviousEvaluations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPreviousEvaluations() got = %v, found %v", got, tt.want)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("getPreviousEvaluations() got = %v, found %v", got, tt.want)
			}
		})
	}
}

func TestEvaluateSLIHandler_HandleEvent(t *testing.T) {
	wg := &sync.WaitGroup{}
	ctx := cloudevents.WithEncodingStructured(context.WithValue(context.Background(), GracefulShutdownKey, wg))
	incomingEvent := cloudevents.NewEvent()
	incomingEvent.SetID("my-id")
	incomingEvent.SetSource("my-source")
	incomingEvent.SetExtension("shkeptncontext", "my-context")
	keptn, _ := keptnv2.NewKeptn(&incomingEvent, keptncommon.KeptnOpts{
		EventSender: &keptnfake.EventSender{},
	})
	var commitID string
	type fields struct {
		Event            cloudevents.Event
		HTTPClient       *http.Client
		KeptnHandler     *keptnv2.Keptn
		SLOFileRetriever SLOFileRetriever
		EventStore       EventStore
	}
	tests := []struct {
		name       string
		fields     fields
		wantID     string
		wantErr    bool
		wantEvents []keptnv2.EvaluationFinishedEventData
	}{
		{
			name: "no SLO file available",
			fields: fields{
				Event: incomingEvent,
				EventStore: &event_handler_mock.EventStoreMock{GetEventsFunc: func(filter *keptnapi.EventFilter) ([]*models.KeptnContextExtendedCE, *models.Error) {
					return []*models.KeptnContextExtendedCE{
						{
							Contenttype:        "",
							Data:               keptnv2.EvaluationTriggeredEventData{},
							ID:                 "my-id",
							Shkeptncontext:     "my-context",
							Shkeptnspecversion: "0.2.0",
							Source:             strutils.Stringp("my-source"),
							Specversion:        "1.0",
							Triggeredid:        "my-triggered-id",
							Type:               strutils.Stringp(keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName)),
						},
					}, nil
				},
				},
				KeptnHandler: keptn,
				SLOFileRetriever: SLOFileRetriever{
					ResourceHandler: &event_handler_mock.ResourceHandlerMock{
						GetResourceFunc: func(scope keptnapi.ResourceScope, options ...keptnapi.URIOption) (*models.Resource, error) {
							return nil, nil
						},
					},
				},
			},
			wantErr: false,
			wantEvents: []keptnv2.EvaluationFinishedEventData{
				{
					EventData: keptnv2.EventData{
						Status:  keptnv2.StatusSucceeded,
						Result:  keptnv2.ResultPass,
						Message: "no evaluation performed by lighthouse because no SLO file configured for project ",
					},
					Evaluation: keptnv2.EvaluationDetails{},
				},
			},
		},
		{
			name: "error reading SLO file",
			fields: fields{
				Event: incomingEvent,
				EventStore: &event_handler_mock.EventStoreMock{GetEventsFunc: func(filter *keptnapi.EventFilter) ([]*models.KeptnContextExtendedCE, *models.Error) {
					return []*models.KeptnContextExtendedCE{
						{
							Contenttype:        "",
							Data:               keptnv2.EvaluationTriggeredEventData{},
							ID:                 "my-id",
							Shkeptncontext:     "my-context",
							Shkeptnspecversion: "0.2.0",
							Source:             strutils.Stringp("my-source"),
							Specversion:        "1.0",
							Triggeredid:        "my-triggered-id",
							Type:               strutils.Stringp(keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName)),
						},
					}, nil
				},
				},
				KeptnHandler: keptn,
				SLOFileRetriever: SLOFileRetriever{
					ResourceHandler: &event_handler_mock.ResourceHandlerMock{
						GetResourceFunc: func(scope keptnapi.ResourceScope, options ...keptnapi.URIOption) (*models.Resource, error) {
							return nil, errors.New("Could not check out branch containing stage config")
						},
					},
					ServiceHandler: &event_handler_mock.ServiceHandlerMock{GetServiceFunc: func(project string, stage string, service string) (*models.Service, error) {
						return &models.Service{}, nil
					}},
				},
			},
			wantErr: false,
			wantEvents: []keptnv2.EvaluationFinishedEventData{
				{
					EventData: keptnv2.EventData{
						Status:  keptnv2.StatusErrored,
						Result:  keptnv2.ResultFailed,
						Message: "could not checkout the SLO",
					},
					Evaluation: keptnv2.EvaluationDetails{},
				},
			},
		},
		{
			name:   "reading SLO file with correct CommitId",
			wantID: "12345",
			fields: fields{
				Event: incomingEvent,
				EventStore: &event_handler_mock.EventStoreMock{GetEventsFunc: func(filter *keptnapi.EventFilter) ([]*models.KeptnContextExtendedCE, *models.Error) {
					return []*models.KeptnContextExtendedCE{
						{
							Contenttype:        "",
							Data:               keptnv2.EvaluationTriggeredEventData{},
							ID:                 "my-id",
							Shkeptncontext:     "my-context",
							Shkeptnspecversion: "0.2.0",
							Source:             strutils.Stringp("my-source"),
							Specversion:        "1.0",
							Triggeredid:        "my-triggered-id",
							Type:               strutils.Stringp(keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName)),
						},
					}, nil
				},
				},
				KeptnHandler: keptn,
				SLOFileRetriever: SLOFileRetriever{
					ResourceHandler: &event_handler_mock.ResourceHandlerMock{
						GetResourceFunc: func(scope keptnapi.ResourceScope, options ...keptnapi.URIOption) (*models.Resource, error) {
							commitID = strings.TrimPrefix(options[0](commitID), "?gitCommitID=")
							myres := models.Resource{Metadata: &models.Version{Version: commitID}}
							return &myres, nil
						},
					},
					ServiceHandler: &event_handler_mock.ServiceHandlerMock{GetServiceFunc: func(project string, stage string, service string) (*models.Service, error) {
						return &models.Service{}, nil
					}},
				},
			},
			wantErr: false,
			wantEvents: []keptnv2.EvaluationFinishedEventData{
				{
					EventData: keptnv2.EventData{
						Status:  keptnv2.StatusSucceeded,
						Result:  keptnv2.ResultPass,
						Message: "duno",
					},
					Evaluation: keptnv2.EvaluationDetails{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commitID = ""
			if tt.wantID != "" {
				tt.fields.Event.SetExtension("gitcommitid", tt.wantID)
			}
			eh := &EvaluateSLIHandler{
				Event:            tt.fields.Event,
				HTTPClient:       tt.fields.HTTPClient,
				KeptnHandler:     tt.fields.KeptnHandler,
				SLOFileRetriever: tt.fields.SLOFileRetriever,
				EventStore:       tt.fields.EventStore,
			}
			if err := eh.HandleEvent(ctx); (err != nil) != tt.wantErr {
				t.Errorf("HandleEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			sender := tt.fields.KeptnHandler.EventSender.(*keptnfake.EventSender)

			// wait for all events to be sent
			require.Eventually(t,
				func() bool {
					if len(sender.SentEvents) != len(tt.wantEvents) {
						return false
					}
					//recycling sender for the next test
					sender.SentEvents = nil
					return true
				},
				time.Second*3, time.Second*1,
			)

			// evaluate which events have been sent
			for index, event := range sender.SentEvents {
				evaluationFinishedEvent := &keptnv2.EvaluationFinishedEventData{}
				if err := event.DataAs(evaluationFinishedEvent); err != nil {
					t.Errorf("could not decode event: %s", err.Error())
				}
				assert.EqualValues(t, tt.wantEvents[index].EventData, (*evaluationFinishedEvent).EventData)
			}
			require.EqualValues(t, commitID, tt.wantID)
		})
	}
}

func Test_aggregateValues(t *testing.T) {
	type fields struct {
		InPreviousResults []*keptnv2.SLIEvaluationResult
		InComparison      *apimodelsv2.SLOComparison
	}
	tests := []struct {
		name        string
		fields      fields
		wantedValue float64
		shouldSkip  bool
	}{

		{name: "Aggregate 2 values with AVG",
			fields: fields{
				InPreviousResults: []*keptnv2.SLIEvaluationResult{
					{
						Score: 2,
						Value: &keptnv2.SLIResult{
							Metric:  "my-test-metric",
							Value:   5.0,
							Success: true,
							Message: "",
						},
						PassTargets:    nil,
						WarningTargets: nil,
						KeySLI:         false,
						Status:         "pass",
					},
					{
						Score: 2,
						Value: &keptnv2.SLIResult{
							Metric:  "my-test-metric",
							Value:   15.0,
							Success: true,
							Message: "",
						},
						PassTargets:    nil,
						WarningTargets: nil,
						KeySLI:         false,
						Status:         "pass",
					},
				},
				InComparison: &apimodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
			},
			wantedValue: 10.0,
			shouldSkip:  false,
		},
		{name: "Skip because of no previous results",
			fields: fields{
				InComparison: &apimodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
			},
			wantedValue: 0.0,
			shouldSkip:  true,
		},
		{name: "Skip because of no previous success",
			fields: fields{
				InPreviousResults: []*keptnv2.SLIEvaluationResult{
					{
						Score: 2,
						Value: &keptnv2.SLIResult{
							Metric:  "my-test-metric",
							Value:   5.0,
							Success: false,
							Message: "",
						},
						PassTargets:    nil,
						WarningTargets: nil,
						KeySLI:         false,
						Status:         "pass",
					},
					{
						Score: 2,
						Value: &keptnv2.SLIResult{
							Metric:  "my-test-metric",
							Value:   15.0,
							Success: false,
							Message: "",
						},
						PassTargets:    nil,
						WarningTargets: nil,
						KeySLI:         false,
						Status:         "pass",
					},
				},
				InComparison: &apimodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
			},
			wantedValue: 0.0,
			shouldSkip:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := aggregateValues(tt.fields.InPreviousResults, tt.fields.InComparison)
			if got != tt.wantedValue {
				t.Errorf("aggregateValues() got = %v, want %v", got, tt.wantedValue)
			}
			if got1 != tt.shouldSkip {
				t.Errorf("aggregateValues() got1 = %v, want %v", got1, tt.shouldSkip)
			}
		})
	}
}

func Test_getSLIResult(t *testing.T) {

	tests := []struct {
		name    string
		results *[]*keptnv2.SLIResult
		sli     string
		found   *keptnv2.SLIResult
		left    []*keptnv2.SLIResult
	}{
		{
			name: "none found",
			sli:  "not_this_metric",
			results: &[]*keptnv2.SLIResult{
				{
					Metric:  "response_time_p50",
					Value:   100,
					Success: true,
					Message: "",
				},
			},
			found: nil,
			left: []*keptnv2.SLIResult{
				{
					Metric:  "response_time_p50",
					Value:   100,
					Success: true,
					Message: "",
				},
			},
		},
		{
			name: "none left",
			sli:  "response_time_p50",
			results: &[]*keptnv2.SLIResult{
				{
					Metric:  "response_time_p50",
					Value:   100,
					Success: true,
					Message: "",
				},
			},
			found: &keptnv2.SLIResult{
				Metric:  "response_time_p50",
				Value:   100,
				Success: true,
				Message: "",
			},
			left: []*keptnv2.SLIResult{},
		},

		{
			name: "one sli left",
			sli:  "response_time_p50",
			results: &[]*keptnv2.SLIResult{
				{
					Metric:  "response_time_p50",
					Value:   100,
					Success: true,
					Message: "",
				},
				{
					Metric:  "wrong_metric",
					Value:   100,
					Success: true,
					Message: "",
				},
			},
			found: &keptnv2.SLIResult{
				Metric:  "response_time_p50",
				Value:   100,
				Success: true,
				Message: "",
			},
			left: []*keptnv2.SLIResult{
				{
					Metric:  "wrong_metric",
					Value:   100,
					Success: true,
					Message: "",
				},
			},
		},
		{
			name: "found last",
			sli:  "response_time_p50",
			results: &[]*keptnv2.SLIResult{

				{
					Metric:  "wrong_metric",
					Value:   100,
					Success: true,
					Message: "",
				},
				{
					Metric:  "response_time_p50",
					Value:   100,
					Success: true,
					Message: "",
				},
			},
			found: &keptnv2.SLIResult{
				Metric:  "response_time_p50",
				Value:   100,
				Success: true,
				Message: "",
			},
			left: []*keptnv2.SLIResult{
				{
					Metric:  "wrong_metric",
					Value:   100,
					Success: true,
					Message: "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getSLIResult(tt.results, tt.sli)
			if !reflect.DeepEqual(got, tt.found) {
				t.Errorf("getSLIResult() = %v, found %v", got, tt.found)
			}
			assert.Equal(t, len(tt.left), len(*tt.results))

		})
	}
}
