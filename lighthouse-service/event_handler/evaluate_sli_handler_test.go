package event_handler

import (
	"errors"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnmodelsv2 "github.com/keptn/go-utils/pkg/models/v2"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
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
	InViolation      *keptnevents.SLIViolation
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
			InViolation: &keptnevents.SLIViolation{
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
			InViolation: &keptnevents.SLIViolation{
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
			InViolation: &keptnevents.SLIViolation{
				Criteria: "?9.0",
			},
			ExpectedResult: false,
			ExpectedError:  errors.New("no operator set"),
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := evaluateFixedThreshold(test.InSLIResult, test.InCriteriaObject, test.InViolation)

			assert.EqualValues(t, test.ExpectedResult, result)
			assert.EqualValues(t, test.ExpectedError, err)
			assert.EqualValues(t, test.InViolation.TargetValue, test.InCriteriaObject.Value)
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
	InViolation       *keptnevents.SLIViolation
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
					Violations: nil,
					Status:     "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Violations: nil,
					Status:     "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InViolation: &keptnevents.SLIViolation{
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
					Violations: nil,
					Status:     "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Violations: nil,
					Status:     "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InViolation: &keptnevents.SLIViolation{
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
					Violations: nil,
					Status:     "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Violations: nil,
					Status:     "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InViolation: &keptnevents.SLIViolation{
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
					Violations: nil,
					Status:     "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Violations: nil,
					Status:     "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InViolation: &keptnevents.SLIViolation{
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
					Violations: nil,
					Status:     "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Violations: nil,
					Status:     "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InViolation: &keptnevents.SLIViolation{
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
					Violations: nil,
					Status:     "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Violations: nil,
					Status:     "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InViolation: &keptnevents.SLIViolation{
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
					Violations: nil,
					Status:     "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Violations: nil,
					Status:     "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "p50",
			},
			InViolation: &keptnevents.SLIViolation{
				Criteria: "<=+10%",
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := evaluateComparison(test.InSLIResult, test.InCriteriaObject, test.InPreviousResults, test.InComparison, test.InViolation)
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
	InViolation       *keptnevents.SLIViolation
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
					Violations: nil,
					Status:     "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Violations: nil,
					Status:     "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InViolation: &keptnevents.SLIViolation{
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
					Violations: nil,
					Status:     "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Violations: nil,
					Status:     "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InViolation: &keptnevents.SLIViolation{
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
					Violations: nil,
					Status:     "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Violations: nil,
					Status:     "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InViolation: &keptnevents.SLIViolation{
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
					Violations: nil,
					Status:     "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Violations: nil,
					Status:     "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InViolation: &keptnevents.SLIViolation{
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
					Violations: nil,
					Status:     "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Violations: nil,
					Status:     "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InViolation: &keptnevents.SLIViolation{
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
					Violations: nil,
					Status:     "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Violations: nil,
					Status:     "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			InViolation: &keptnevents.SLIViolation{
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
					Violations: nil,
					Status:     "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Violations: nil,
					Status:     "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "p50",
			},
			InViolation: &keptnevents.SLIViolation{
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
			InViolation: &keptnevents.SLIViolation{
				Criteria: "<=+10%",
			},
			ExpectedResult: true,
			ExpectedError:  nil,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := evaluateSingleCriteria(test.InSLIResult, test.InCriteria, test.InPreviousResults, test.InComparison, test.InViolation)
			assert.EqualValues(t, test.ExpectedResult, result)
			assert.EqualValues(t, test.ExpectedError, err)
		})
	}
}

type evaluateCriteriaSetTestObject struct {
	Name               string
	InSLIResult        *keptnevents.SLIResult
	InCriteriaSet      *keptnmodelsv2.SLOCriteria
	InPreviousResults  []*keptnevents.SLIEvaluationResult
	InComparison       *keptnmodelsv2.SLOComparison
	ExpectedViolations []*keptnevents.SLIViolation
	ExpectedResult     bool
	ExpectedError      error
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
					Violations: nil,
					Status:     "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Violations: nil,
					Status:     "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			ExpectedViolations: nil,
			ExpectedResult:     true,
			ExpectedError:      nil,
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
					Violations: nil,
					Status:     "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Violations: nil,
					Status:     "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			ExpectedViolations: []*keptnevents.SLIViolation{
				{
					Criteria:    "<=10.0",
					TargetValue: 10.0,
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
			assert.EqualValues(t, test.ExpectedViolations, violations)
			assert.EqualValues(t, test.ExpectedError, err)
		})
	}
}

type evaluateOrCombinedCriteriaTestObject struct {
	Name               string
	InSLIResult        *keptnevents.SLIResult
	InCriteriaSets     []*keptnmodelsv2.SLOCriteria
	InPreviousResults  []*keptnevents.SLIEvaluationResult
	InComparison       *keptnmodelsv2.SLOComparison
	ExpectedViolations []*keptnevents.SLIViolation
	ExpectedResult     bool
	ExpectedError      error
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
					Violations: nil,
					Status:     "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Violations: nil,
					Status:     "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			ExpectedViolations: nil,
			ExpectedResult:     true,
			ExpectedError:      nil,
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
					Violations: nil,
					Status:     "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Violations: nil,
					Status:     "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			ExpectedViolations: []*keptnevents.SLIViolation{
				{
					Criteria:    "<=10.0",
					TargetValue: 10.0,
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
					Violations: nil,
					Status:     "pass",
				},
				{
					Score: 2,
					Value: &keptnevents.SLIResult{
						Metric:  "my-test-metric",
						Value:   10.0,
						Success: true,
						Message: "",
					},
					Violations: nil,
					Status:     "pass",
				},
			},
			InComparison: &keptnmodelsv2.SLOComparison{
				CompareWith:               "several_results",
				IncludeResultWithScore:    "pass",
				NumberOfComparisonResults: 2,
				AggregateFunction:         "avg",
			},
			ExpectedViolations: []*keptnevents.SLIViolation{
				{
					Criteria:    "<=10.0",
					TargetValue: 10.0,
				},
				{
					Criteria:    "<=+10%",
					TargetValue: 11.0,
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
				assert.EqualValues(t, test.ExpectedViolations, violations)
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
								Violations: nil,
								Status:     "pass",
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
							Violations: nil,
							Status:     "pass",
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
