package event_handler

import (
	"encoding/json"
	"errors"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptnutils "github.com/keptn/go-utils/pkg/lib"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnmodelsv2 "github.com/keptn/go-utils/pkg/lib"
	"github.com/stretchr/testify/assert"
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
	InSLIResult      *keptnevents.SLIResult
	InCriteriaObject *criteriaObject
	InTarget         *keptnevents.SLITarget
	ExpectedResult   bool
	ExpectedError    error
}

func TestEvaluateFixedThreshold(t *testing.T) {
	tests := []*evaluateFixedThresholdTestObject{
		{
			Name: "10.0 > 9.0 should return true",
			InSLIResult: &keptnevents.SLIResult{
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
			InTarget: &keptnevents.SLITarget{
				Criteria: ">9.0",
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name: "10.0 = 9.0 should return false",
			InSLIResult: &keptnevents.SLIResult{
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
			InTarget: &keptnevents.SLITarget{
				Criteria: "=9.0",
			},
			ExpectedResult: false,
			ExpectedError:  nil,
		},
		{
			Name: "10.0 ? 9.0 should return an error",
			InSLIResult: &keptnevents.SLIResult{
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
			InTarget: &keptnevents.SLITarget{
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
	InSLIResult       *keptnevents.SLIResult
	InCriteriaObject  *criteriaObject
	InPreviousResults []*keptnevents.SLIEvaluationResult
	InComparison      *keptnmodelsv2.SLOComparison
	InTarget          *keptnevents.SLITarget
	ExpectedResult    bool
	ExpectedError     error
}

func TestEvaluateComparison(t *testing.T) {
	tests := []*evaluateComparisonTestObject{
		{
			Name: "Expect true for 10.0 <= avg([10.0, 10.0]) + 10%",
			InSLIResult: &keptnevents.SLIResult{
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
			InPreviousResults: []*keptnevents.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnevents.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name: "Expect false for 11.01 <= avg([10.0, 10.0]) + 10%",
			InSLIResult: &keptnevents.SLIResult{
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
			InPreviousResults: []*keptnevents.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnevents.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: false,
			ExpectedError:  nil,
		},
		{
			Name: "Expect false for 10.0 < avg([10.0, 10.0]) + 0%",
			InSLIResult: &keptnevents.SLIResult{
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
			InPreviousResults: []*keptnevents.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnevents.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: false,
			ExpectedError:  nil,
		},
		{
			Name: "Expect false for 10.0 > avg([10.0, 10.0]) + 0%",
			InSLIResult: &keptnevents.SLIResult{
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
			InPreviousResults: []*keptnevents.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnevents.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: false,
			ExpectedError:  nil,
		},
		{
			Name: "Expect true for 11.0 = avg([10.0, 10.0]) + 1.0",
			InSLIResult: &keptnevents.SLIResult{
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
			InPreviousResults: []*keptnevents.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnevents.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name: "Expect true for 9.0 = avg([10.0, 10.0]) - 1.0",
			InSLIResult: &keptnevents.SLIResult{
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
			InPreviousResults: []*keptnevents.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnevents.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name: "Expect true for 10.0 <= p50([10.0, 5.0]) + 10.0%",
			InSLIResult: &keptnevents.SLIResult{
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
			InPreviousResults: []*keptnevents.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "p50",
			},
			InTarget: &keptnevents.SLITarget{
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
	InSLIResult       *keptnevents.SLIResult
	InCriteria        string
	InPreviousResults []*keptnevents.SLIEvaluationResult
	InComparison      *keptnmodelsv2.SLOComparison
	InTarget          *keptnevents.SLITarget
	ExpectedResult    bool
	ExpectedError     error
}

func TestEvaluateSingleCriteria(t *testing.T) {
	tests := []*evaluateSingleCriteriaTestObject{
		{
			Name: "Expect true for 10.0 <= avg([10.0, 10.0]) + 10%",
			InSLIResult: &keptnevents.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteria: "<=+10%",
			InPreviousResults: []*keptnevents.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnevents.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name: "Expect false for 11.01 <= avg([10.0, 10.0]) + 10%",
			InSLIResult: &keptnevents.SLIResult{
				Metric:  "my-test-metric",
				Value:   11.01,
				Success: true,
				Message: "",
			},
			InCriteria: "<=+10%",
			InPreviousResults: []*keptnevents.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnevents.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: false,
			ExpectedError:  nil,
		},
		{
			Name: "Expect false for 10.0 < avg([10.0, 10.0]) + 0%",
			InSLIResult: &keptnevents.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteria: "<+0%",
			InPreviousResults: []*keptnevents.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnevents.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: false,
			ExpectedError:  nil,
		},
		{
			Name: "Expect false for 10.0 > avg([10.0, 10.0]) + 0%",
			InSLIResult: &keptnevents.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteria: ">+0%",
			InPreviousResults: []*keptnevents.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnevents.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: false,
			ExpectedError:  nil,
		},
		{
			Name: "Expect true for 11.0 = avg([10.0, 10.0]) + 1.0",
			InSLIResult: &keptnevents.SLIResult{
				Metric:  "my-test-metric",
				Value:   11.0,
				Success: true,
				Message: "",
			},
			InCriteria: "=+1",
			InPreviousResults: []*keptnevents.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnevents.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name: "Expect true for 9.0 = avg([10.0, 10.0]) - 1.0",
			InSLIResult: &keptnevents.SLIResult{
				Metric:  "my-test-metric",
				Value:   9.0,
				Success: true,
				Message: "",
			},
			InCriteria: "=-1",
			InPreviousResults: []*keptnevents.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InTarget: &keptnevents.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name: "Expect true for 10.0 <= p50([10.0, 5.0]) + 10.0%",
			InSLIResult: &keptnevents.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteria: "<=+10%",
			InPreviousResults: []*keptnevents.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "p50",
			},
			InTarget: &keptnevents.SLITarget{
				Criteria: "<=+10%",
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Name: "Expect true for 10.0 <= 10.0",
			InSLIResult: &keptnevents.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteria: "<=10",
			InTarget: &keptnevents.SLITarget{
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
	InSLIResult       *keptnevents.SLIResult
	InCriteriaSet     *keptnmodelsv2.SLOCriteria
	InPreviousResults []*keptnevents.SLIEvaluationResult
	InComparison      *keptnmodelsv2.SLOComparison
	ExpectedTargets   []*keptnevents.SLITarget
	ExpectedResult    bool
	ExpectedError     error
}

func TestEvaluateCriteriaSet(t *testing.T) {
	tests := []*evaluateCriteriaSetTestObject{
		{
			Name: "Expect true for (10.0 <= avg([10.0, 10.0]) + 10%) && (10.0 <= 10.0)",
			InSLIResult: &keptnevents.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteriaSet: &keptnmodelsv2.SLOCriteria{
				Criteria: []string{"<=+10%", "<=10.0"},
			},
			InPreviousResults: []*keptnevents.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			ExpectedTargets: []*keptnevents.SLITarget{
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
			InSLIResult: &keptnevents.SLIResult{
				Metric:  "my-test-metric",
				Value:   11.0,
				Success: true,
				Message: "",
			},
			InCriteriaSet: &keptnmodelsv2.SLOCriteria{
				Criteria: []string{"<=+10%", "<=10.0"},
			},
			InPreviousResults: []*keptnevents.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			ExpectedTargets: []*keptnevents.SLITarget{
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
	InSLIResult       *keptnevents.SLIResult
	InCriteriaSets    []*keptnmodelsv2.SLOCriteria
	InPreviousResults []*keptnevents.SLIEvaluationResult
	InComparison      *keptnmodelsv2.SLOComparison
	ExpectedTargets   []*keptnevents.SLITarget
	ExpectedResult    bool
	ExpectedError     error
}

func TestEvaluateOrCombinedCriteria(t *testing.T) {
	tests := []*evaluateOrCombinedCriteriaTestObject{
		{
			Name: "Expect true for (10.0 <= avg([10.0, 10.0]) + 10%) || (10.0 <= 10.0)",
			InSLIResult: &keptnevents.SLIResult{
				Metric:  "my-test-metric",
				Value:   10.0,
				Success: true,
				Message: "",
			},
			InCriteriaSets: []*keptnmodelsv2.SLOCriteria{
				{
					Criteria: []string{"<=10.0"},
				},
				{
					Criteria: []string{"<=+10%"},
				},
			},
			InPreviousResults: []*keptnevents.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			ExpectedTargets: []*keptnevents.SLITarget{
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
			InSLIResult: &keptnevents.SLIResult{
				Metric:  "my-test-metric",
				Value:   11.0,
				Success: true,
				Message: "",
			},
			InCriteriaSets: []*keptnmodelsv2.SLOCriteria{
				{
					Criteria: []string{"<=10.0"},
				},
				{
					Criteria: []string{"<=+10%"},
				},
			},
			InPreviousResults: []*keptnevents.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			ExpectedTargets: []*keptnevents.SLITarget{
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
			InSLIResult: &keptnevents.SLIResult{
				Metric:  "my-test-metric",
				Value:   20.0,
				Success: true,
				Message: "",
			},
			InCriteriaSets: []*keptnmodelsv2.SLOCriteria{
				{
					Criteria: []string{"<=10.0"},
				},
				{
					Criteria: []string{"<=+10%"},
				},
			},
			InPreviousResults: []*keptnevents.SLIEvaluationResult{
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Targets: nil,
					Status:  "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			ExpectedTargets: []*keptnevents.SLITarget{
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
	InGetSLIDoneEvent          *keptnevents.InternalGetSLIDoneEventData
	InSLOConfig                *keptnmodelsv2.ServiceLevelObjectives
	InPreviousEvaluationEvents []*keptnevents.EvaluationDoneEventData
	ExpectedEvaluationResult   *keptnevents.EvaluationDoneEventData
	ExpectedMaximumScore       float64
	ExpectedKeySLIFailed       bool
}

func TestEvaluateObjectives(t *testing.T) {
	tests := []*evaluateObjectivesTestObject{
		{
			Name: "Simple comparison evaluation",
			InGetSLIDoneEvent: &keptnevents.InternalGetSLIDoneEventData{
				Project: "sockshop",
				Service: "carts",
				Stage:   "dev",
				Start:   "2019-10-20T07:57:27.152330783Z",
				End:     "2019-10-22T08:57:27.152330783Z",
				IndicatorValues: []*keptnevents.SLIResult{
					{
						Metric:  "my-test-metric-1",
						Value:   10.0,
						Success: true,
						Message: "",
					},
				},
			},
			InSLOConfig: &keptnmodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &keptnmodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
				Objectives: []*keptnmodelsv2.SLO{
					{
						SLI: "my-test-metric-1",
						Pass: []*keptnmodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=15.0"},
							},
							{
								Criteria: []string{"<=+10%"},
							},
						},
						Warning: []*keptnmodelsv2.SLOCriteria{
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
				TotalScore: &keptnmodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InPreviousEvaluationEvents: []*keptnevents.EvaluationDoneEventData{
				{
					EvaluationDetails: &keptnevents.EvaluationDetails{
						TimeStart: "",
						TimeEnd:   "",
						Result:    "pass",
						Score:     2,
						IndicatorResults: []*keptnevents.SLIEvaluationResult{
							{
								Score: 2,
								Value: &keptnevents.SLIResult{
									Metric:  "my-test-metric-1",
									Value:   10.0,
									Success: true,
									Message: "",
								},
								Targets: nil,
								Status:  "pass",
							},
						},
					},
					Result:       "pass",
					Project:      "sockshop",
					Service:      "carts",
					Stage:        "dev",
					TestStrategy: "performance",
				},
			},
			ExpectedEvaluationResult: &keptnevents.EvaluationDoneEventData{
				EvaluationDetails: &keptnevents.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "", // not set by the tested function
					Score:     0,  // not calculated by tested function
					IndicatorResults: []*keptnevents.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnevents.SLIResult{
								Metric:  "my-test-metric-1",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							Targets: []*keptnevents.SLITarget{
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
							Status: "pass",
						},
					},
				},
				Result:       "", // not set by the tested function
				Project:      "sockshop",
				Service:      "carts",
				Stage:        "dev",
				TestStrategy: "",
			},
			ExpectedMaximumScore: 1,
			ExpectedKeySLIFailed: false,
		},
		{
			Name: "Expect Warning",
			InGetSLIDoneEvent: &keptnevents.InternalGetSLIDoneEventData{
				Project: "sockshop",
				Service: "carts",
				Stage:   "dev",
				Start:   "2019-10-20T07:57:27.152330783Z",
				End:     "2019-10-22T08:57:27.152330783Z",
				IndicatorValues: []*keptnevents.SLIResult{
					{
						Metric:  "my-test-metric-1",
						Value:   16.0,
						Success: true,
						Message: "",
					},
				},
			},
			InSLOConfig: &keptnmodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &keptnmodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
				Objectives: []*keptnmodelsv2.SLO{
					{
						SLI: "my-test-metric-1",
						Pass: []*keptnmodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=15.0"},
							},
							{
								Criteria: []string{"<=+10%"},
							},
						},
						Warning: []*keptnmodelsv2.SLOCriteria{
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
				TotalScore: &keptnmodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InPreviousEvaluationEvents: []*keptnevents.EvaluationDoneEventData{
				{
					EvaluationDetails: &keptnevents.EvaluationDetails{
						TimeStart: "",
						TimeEnd:   "",
						Result:    "pass",
						Score:     2,
						IndicatorResults: []*keptnevents.SLIEvaluationResult{
							{
								Score: 2,
								Value: &keptnevents.SLIResult{
									Metric:  "my-test-metric-1",
									Value:   10.0,
									Success: true,
									Message: "",
								},
								Targets: nil,
								Status:  "pass",
							},
						},
					},
					Result:       "pass",
					Project:      "sockshop",
					Service:      "carts",
					Stage:        "dev",
					TestStrategy: "performance",
				},
			},
			ExpectedEvaluationResult: &keptnevents.EvaluationDoneEventData{
				EvaluationDetails: &keptnevents.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "", // not set by the tested function
					Score:     0,  // not calculated by tested function
					IndicatorResults: []*keptnevents.SLIEvaluationResult{
						{
							Score: 0.5,
							Value: &keptnevents.SLIResult{
								Metric:  "my-test-metric-1",
								Value:   16.0,
								Success: true,
								Message: "",
							},
							Targets: []*keptnevents.SLITarget{
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
							Status: "warning",
						},
					},
				},
				Result:       "", // not set by the tested function
				Project:      "sockshop",
				Service:      "carts",
				Stage:        "dev",
				TestStrategy: "",
			},
			ExpectedMaximumScore: 1,
			ExpectedKeySLIFailed: false,
		},
		{
			Name: "Logging SLI with no pass criteria should not affect score",
			InGetSLIDoneEvent: &keptnevents.InternalGetSLIDoneEventData{
				Project: "sockshop",
				Service: "carts",
				Stage:   "dev",
				Start:   "2019-10-20T07:57:27.152330783Z",
				End:     "2019-10-22T08:57:27.152330783Z",
				IndicatorValues: []*keptnevents.SLIResult{
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
			InSLOConfig: &keptnmodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &keptnmodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
				Objectives: []*keptnmodelsv2.SLO{
					{
						SLI: "my-test-metric-1",
						Pass: []*keptnmodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=15.0"},
							},
							{
								Criteria: []string{"<=+10%"},
							},
						},
						Warning: []*keptnmodelsv2.SLOCriteria{
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
				TotalScore: &keptnmodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InPreviousEvaluationEvents: []*keptnevents.EvaluationDoneEventData{
				{
					EvaluationDetails: &keptnevents.EvaluationDetails{
						TimeStart: "",
						TimeEnd:   "",
						Result:    "pass",
						Score:     2,
						IndicatorResults: []*keptnevents.SLIEvaluationResult{
							{
								Score: 2,
								Value: &keptnevents.SLIResult{
									Metric:  "my-test-metric-1",
									Value:   10.0,
									Success: true,
									Message: "",
								},
								Targets: nil,
								Status:  "pass",
							},
						},
					},
					Result:       "pass",
					Project:      "sockshop",
					Service:      "carts",
					Stage:        "dev",
					TestStrategy: "performance",
				},
			},
			ExpectedEvaluationResult: &keptnevents.EvaluationDoneEventData{
				EvaluationDetails: &keptnevents.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "", // not set by the tested function
					Score:     0,  // not calculated by tested function
					IndicatorResults: []*keptnevents.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnevents.SLIResult{
								Metric:  "my-test-metric-1",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							Targets: []*keptnevents.SLITarget{
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
							Status: "pass",
						},
						{
							Score: 0,
							Value: &keptnevents.SLIResult{
								Metric:  "my-log-metric",
								Value:   30.0,
								Success: true,
								Message: "",
							},
							Status: "info",
						},
					},
				},
				Result:       "", // not set by the tested function
				Project:      "sockshop",
				Service:      "carts",
				Stage:        "dev",
				TestStrategy: "",
			},
			ExpectedMaximumScore: 1,
			ExpectedKeySLIFailed: false,
		},
		{
			Name: "BUG 1125",
			InGetSLIDoneEvent: &keptnevents.InternalGetSLIDoneEventData{
				Project: "sockshop",
				Service: "carts",
				Stage:   "dev",
				Start:   "2019-10-20T07:57:27.152330783Z",
				End:     "2019-10-22T08:57:27.152330783Z",
				IndicatorValues: []*keptnevents.SLIResult{
					{
						Metric:  "response_time_p50",
						Value:   1011.0745528937252,
						Success: true,
						Message: "",
					},
				},
			},
			InSLOConfig: &keptnmodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &keptnmodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
				Objectives: []*keptnmodelsv2.SLO{
					{
						SLI: "response_time_p50",
						Pass: []*keptnmodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=+20%", "<500"},
							},
						},
						Weight: 1,
						KeySLI: false,
					},
				},
				TotalScore: &keptnmodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InPreviousEvaluationEvents: nil,
			ExpectedEvaluationResult: &keptnevents.EvaluationDoneEventData{
				EvaluationDetails: &keptnevents.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "", // not set by the tested function
					Score:     0,  // not calculated by tested function
					IndicatorResults: []*keptnevents.SLIEvaluationResult{
						{
							Score: 0,
							Value: &keptnevents.SLIResult{
								Metric:  "response_time_p50",
								Value:   1011.0745528937252,
								Success: true,
								Message: "",
							},
							Targets: []*keptnevents.SLITarget{
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
							Status: "fail",
						},
					},
				},
				Result:       "", // not set by the tested function
				Project:      "sockshop",
				Service:      "carts",
				Stage:        "dev",
				TestStrategy: "",
			},
			ExpectedMaximumScore: 1,
			ExpectedKeySLIFailed: false,
		},
		{
			Name: "BUG 1263",
			InGetSLIDoneEvent: &keptnevents.InternalGetSLIDoneEventData{
				Project: "sockshop",
				Service: "carts",
				Stage:   "dev",
				Start:   "2019-10-20T07:57:27.152330783Z",
				End:     "2019-10-22T08:57:27.152330783Z",
				IndicatorValues: []*keptnevents.SLIResult{
					{
						Metric:  "response_time_p50",
						Value:   100,
						Success: true,
						Message: "",
					},
				},
			},
			InSLOConfig: &keptnmodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &keptnmodelsv2.SLOComparison{
					CompareWith:               "single_result",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 1,
					AggregateFunction:         "avg",
				},
				Objectives: []*keptnmodelsv2.SLO{
					{
						SLI: "response_time_p50",
						Pass: []*keptnmodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=+20%"},
							},
						},
						Weight: 1,
						KeySLI: false,
					},
				},
				TotalScore: &keptnmodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InPreviousEvaluationEvents: []*keptnevents.EvaluationDoneEventData{
				{
					EvaluationDetails: &keptnevents.EvaluationDetails{
						TimeStart: "",
						TimeEnd:   "",
						Result:    "fail",
						Score:     0,
						IndicatorResults: []*keptnevents.SLIEvaluationResult{
							{
								Score: 0,
								Value: &keptnevents.SLIResult{
									Metric:  "response_time_p50",
									Value:   0.0,
									Success: false,
									Message: "",
								},
								Targets: nil,
								Status:  "fail",
							},
						},
					},
					Result:       "pass",
					Project:      "sockshop",
					Service:      "carts",
					Stage:        "dev",
					TestStrategy: "performance",
				},
			},
			ExpectedEvaluationResult: &keptnevents.EvaluationDoneEventData{
				EvaluationDetails: &keptnevents.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "", // not set by the tested function
					Score:     0,  // not calculated by tested function
					IndicatorResults: []*keptnevents.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnevents.SLIResult{
								Metric:  "response_time_p50",
								Value:   100,
								Success: true,
								Message: "",
							},
							Targets: []*keptnevents.SLITarget{
								{
									Criteria:    "<=+20%",
									TargetValue: 0,
									Violated:    false,
								},
							},
							Status: "pass",
						},
					},
				},
				Result:       "", // not set by the tested function
				Project:      "sockshop",
				Service:      "carts",
				Stage:        "dev",
				TestStrategy: "",
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
	InEvaluationResult       *keptnevents.EvaluationDoneEventData
	InSLOConfig              *keptnmodelsv2.ServiceLevelObjectives
	InKeySLIFailed           bool
	ExpectedEvaluationResult *keptnevents.EvaluationDoneEventData
	ExpectedError            error
}

func TestCalculateScore(t *testing.T) {
	tests := []*calculateScoreTestObject{
		{
			Name:           "Simple comparison",
			InMaximumScore: 1,
			InEvaluationResult: &keptnevents.EvaluationDoneEventData{
				EvaluationDetails: &keptnevents.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "", // to be calculated
					Score:     0,  // to be calculated
					IndicatorResults: []*keptnevents.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnevents.SLIResult{
								Metric:  "my-test-metric-1",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							Targets: nil,
							Status:  "pass",
						},
					},
				},
				Result:       "", // to be set
				Project:      "sockshop",
				Service:      "carts",
				Stage:        "dev",
				TestStrategy: "",
			},
			InSLOConfig: &keptnmodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &keptnmodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
				Objectives: []*keptnmodelsv2.SLO{
					{
						SLI: "my-test-metric-1",
						Pass: []*keptnmodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=15.0"},
							},
							{
								Criteria: []string{"<=+10%"},
							},
						},
						Warning: []*keptnmodelsv2.SLOCriteria{
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
				TotalScore: &keptnmodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InKeySLIFailed: false,
			ExpectedEvaluationResult: &keptnevents.EvaluationDoneEventData{
				EvaluationDetails: &keptnevents.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "pass",
					Score:     100.0,
					IndicatorResults: []*keptnevents.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnevents.SLIResult{
								Metric:  "my-test-metric-1",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							Targets: nil,
							Status:  "pass",
						},
					},
				},
				Result:       "pass",
				Project:      "sockshop",
				Service:      "carts",
				Stage:        "dev",
				TestStrategy: "",
			},
			ExpectedError: nil,
		},
		{
			Name:           "Key SLI failed",
			InMaximumScore: 2,
			InEvaluationResult: &keptnevents.EvaluationDoneEventData{
				EvaluationDetails: &keptnevents.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "", // to be calculated
					Score:     0,  // to be calculated
					IndicatorResults: []*keptnevents.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnevents.SLIResult{
								Metric:  "my-test-metric-1",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							Targets: nil,
							Status:  "pass",
						},
						{
							Score: 0,
							Value: &keptnevents.SLIResult{
								Metric:  "my-key-metric",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							Targets: nil,
							Status:  "fail",
						},
					},
				},
				Result:       "", // to be set
				Project:      "sockshop",
				Service:      "carts",
				Stage:        "dev",
				TestStrategy: "",
			},
			InSLOConfig: &keptnmodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &keptnmodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
				Objectives: []*keptnmodelsv2.SLO{
					{
						SLI: "my-test-metric-1",
						Pass: []*keptnmodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=15.0"},
							},
							{
								Criteria: []string{"<=+10%"},
							},
						},
						Warning: []*keptnmodelsv2.SLOCriteria{
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
						Pass: []*keptnmodelsv2.SLOCriteria{
							{
								Criteria: []string{"<=15.0"},
							},
							{
								Criteria: []string{"<=+10%"},
							},
						},
						Warning: []*keptnmodelsv2.SLOCriteria{
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
				TotalScore: &keptnmodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InKeySLIFailed: true,
			ExpectedEvaluationResult: &keptnevents.EvaluationDoneEventData{
				EvaluationDetails: &keptnevents.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "fail",
					Score:     50.0,
					IndicatorResults: []*keptnevents.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnevents.SLIResult{
								Metric:  "my-test-metric-1",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							Targets: nil,
							Status:  "pass",
						},
						{
							Score: 0,
							Value: &keptnevents.SLIResult{
								Metric:  "my-key-metric",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							Targets: nil,
							Status:  "fail",
						},
					},
				},
				Result:       "fail",
				Project:      "sockshop",
				Service:      "carts",
				Stage:        "dev",
				TestStrategy: "",
			},
			ExpectedError: nil,
		},
		{
			Name:           "Only Info SLIs",
			InMaximumScore: 0,
			InEvaluationResult: &keptnevents.EvaluationDoneEventData{
				EvaluationDetails: &keptnevents.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "", // to be calculated
					Score:     0,  // to be calculated
					IndicatorResults: []*keptnevents.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnevents.SLIResult{
								Metric:  "my-test-metric-1",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							Targets: nil,
							Status:  "pass",
						},
						{
							Score: 1,
							Value: &keptnevents.SLIResult{
								Metric:  "my-key-metric",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							Targets: nil,
							Status:  "pass",
						},
					},
				},
				Result:       "", // to be set
				Project:      "sockshop",
				Service:      "carts",
				Stage:        "dev",
				TestStrategy: "",
			},
			InSLOConfig: &keptnmodelsv2.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter:      nil,
				Comparison: &keptnmodelsv2.SLOComparison{
					CompareWith:               "several_results",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 2,
					AggregateFunction:         "avg",
				},
				Objectives: []*keptnmodelsv2.SLO{
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
				TotalScore: &keptnmodelsv2.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			InKeySLIFailed: true,
			ExpectedEvaluationResult: &keptnevents.EvaluationDoneEventData{
				EvaluationDetails: &keptnevents.EvaluationDetails{
					TimeStart: "2019-10-20T07:57:27.152330783Z",
					TimeEnd:   "2019-10-22T08:57:27.152330783Z",
					Result:    "pass",
					Score:     100.0,
					IndicatorResults: []*keptnevents.SLIEvaluationResult{
						{
							Score: 1,
							Value: &keptnevents.SLIResult{
								Metric:  "my-test-metric-1",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							Targets: nil,
							Status:  "pass",
						},
						{
							Score: 1,
							Value: &keptnevents.SLIResult{
								Metric:  "my-key-metric",
								Value:   10.0,
								Success: true,
								Message: "",
							},
							Targets: nil,
							Status:  "pass",
						},
					},
				},
				Result:       "pass",
				Project:      "sockshop",
				Service:      "carts",
				Stage:        "dev",
				TestStrategy: "",
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

func TestEvaluateSLIHandler_getPreviousTestExecutionResult(t *testing.T) {

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

	_ = os.Setenv("MONGODB_DATASTORE", strings.TrimPrefix(ts.URL, "http://"))

	type fields struct {
		Logger     *keptnutils.Logger
		Event      cloudevents.Event
		HTTPClient *http.Client
	}
	type args struct {
		e            *keptnevents.InternalGetSLIDoneEventData
		keptnContext string
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		resultFromDatastore datastoreResult
		want                *keptnevents.TestsFinishedEventData
		wantErr             bool
	}{
		{
			name: "return a tests-finished event",
			fields: fields{
				Logger:     nil,
				Event:      cloudevents.Event{},
				HTTPClient: &http.Client{},
			},
			args: args{
				e: &keptnevents.InternalGetSLIDoneEventData{
					Project:            "sockshop",
					Stage:              "dev",
					Service:            "carts",
					Start:              "",
					End:                "",
					TestStrategy:       "",
					IndicatorValues:    nil,
					DeploymentStrategy: "",
					Deployment:         "",
					Labels:             nil,
				},
				keptnContext: "",
			},
			resultFromDatastore: datastoreResult{
				NextPageKey: "",
				TotalCount:  1,
				PageSize:    1,
				Events: []struct {
					Data interface{} `json:"data"`
				}{
					{
						Data: &keptnevents.TestsFinishedEventData{
							Project:            "sockshop",
							Service:            "carts",
							Stage:              "dev",
							TestStrategy:       "",
							DeploymentStrategy: "",
							Start:              "",
							End:                "",
							Labels:             nil,
							Result:             "",
						},
					},
				},
			},
			want: &keptnevents.TestsFinishedEventData{
				Project:            "sockshop",
				Service:            "carts",
				Stage:              "dev",
				TestStrategy:       "",
				DeploymentStrategy: "",
				Start:              "",
				End:                "",
				Labels:             nil,
				Result:             "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			returnedResult = tt.resultFromDatastore

			//keptnHandler, _ := keptnevents.NewKeptn(&tt.fields.Event, keptnevents.KeptnOpts{})
			eh := &EvaluateSLIHandler{
				KeptnHandler: nil,
				Event:        tt.fields.Event,
				HTTPClient:   tt.fields.HTTPClient,
			}
			got, err := eh.getPreviousTestExecutionResult(tt.args.e, tt.args.keptnContext)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPreviousTestExecutionResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPreviousTestExecutionResult() got = %v, want %v", got, tt.want)
			}
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

	_ = os.Setenv("MONGODB_DATASTORE", strings.TrimPrefix(ts.URL, "http://"))

	type fields struct {
		Logger     *keptnutils.Logger
		Event      cloudevents.Event
		HTTPClient *http.Client
	}
	type args struct {
		e                       *keptnevents.InternalGetSLIDoneEventData
		numberOfPreviousResults int
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		resultFromDatastore datastoreResult
		want                []*keptnevents.EvaluationDoneEventData
		wantErr             bool
	}{
		{
			name: "get eveluation-done events",
			fields: fields{
				Logger:     nil,
				Event:      cloudevents.Event{},
				HTTPClient: &http.Client{},
			},
			args: args{
				e: &keptnevents.InternalGetSLIDoneEventData{
					Project:            "sockshop",
					Stage:              "dev",
					Service:            "carts",
					Start:              "",
					End:                "",
					TestStrategy:       "",
					IndicatorValues:    nil,
					DeploymentStrategy: "",
					Deployment:         "",
					Labels:             nil,
				},
				numberOfPreviousResults: 1,
			},
			resultFromDatastore: datastoreResult{
				NextPageKey: "",
				TotalCount:  1,
				PageSize:    1,
				Events: []struct {
					Data interface{} `json:"data"`
				}{
					{
						Data: &keptnevents.EvaluationDoneEventData{
							Project:            "sockshop",
							Service:            "carts",
							Stage:              "dev",
							TestStrategy:       "",
							DeploymentStrategy: "",
							Labels:             nil,
							Result:             "",
						},
					},
				},
			},
			want: []*keptnevents.EvaluationDoneEventData{
				{
					Project:            "sockshop",
					Service:            "carts",
					Stage:              "dev",
					TestStrategy:       "",
					DeploymentStrategy: "",
					Labels:             nil,
					Result:             "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		returnedResult = tt.resultFromDatastore
		t.Run(tt.name, func(t *testing.T) {

			//keptnHandler, _ := keptnevents.NewKeptn(&tt.fields.Event, keptnevents.KeptnOpts{})

			eh := &EvaluateSLIHandler{
				KeptnHandler: nil,
				Event:        tt.fields.Event,
				HTTPClient:   tt.fields.HTTPClient,
			}
			got, err := eh.getPreviousEvaluations(tt.args.e, tt.args.numberOfPreviousResults)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPreviousEvaluations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPreviousEvaluations() got = %v, want %v", got, tt.want)
			}
		})
	}
}
