package event_handler

import (
	"errors"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/types"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	keptn "github.com/keptn/go-utils/pkg/lib"

	"testing"

	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type getSLOTestObject struct {
	Name           string
	SLOFileContent string
	ExpectedSLO    *keptn.ServiceLevelObjectives
	ExpectedError  error
}

func TestParseLO(t *testing.T) {
	tests := []*getSLOTestObject{
		{
			Name: "Simple SLO file",
			SLOFileContent: `---
spec_version: '1.0'
filter:
  id: "<prometheus_scrape_job_id>"
comparison:
  compare_with: "single_result"
  include_result_with_score: "pass"
  number_of_comparison_results: 3
  aggregate_function: avg
objectives:
  - sli: responseTime95
    pass:
      - criteria:
          - "<=+10%"
      - criteria:
          - "<200"
    warning:
      - criteria:
          - "<+15%"
          - ">-8%"
          - "<500"
  - null
  - sli: security_vulnerabilities
    weight: 2
    pass:
      - criteria:
          - "=0"
  - sli: sql_statements
    key_sli: true
    pass:
      - criteria:
          - "=0%"
      - criteria:
          - "<100"
    warning:
      - criteria:
          - "<+5%"
          - ">-5%"
total_score:
  pass: "90%"
  warning: 75%`,
			ExpectedSLO: &keptn.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter: map[string]string{
					"id": "<prometheus_scrape_job_id>",
				},
				Comparison: &keptn.SLOComparison{
					CompareWith:               "single_result",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 3,
					AggregateFunction:         "avg",
				},
				Objectives: []*keptn.SLO{
					{
						SLI: "responseTime95",
						Pass: []*keptn.SLOCriteria{
							{
								Criteria: []string{"<=+10%"},
							},
							{
								Criteria: []string{"<200"},
							},
						},
						Warning: []*keptn.SLOCriteria{
							{
								Criteria: []string{"<+15%", ">-8%", "<500"},
							},
						},
						Weight: 1,
						KeySLI: false,
					},
					{
						SLI: "security_vulnerabilities",
						Pass: []*keptn.SLOCriteria{
							{
								Criteria: []string{"=0"},
							},
						},
						Weight: 2,
						KeySLI: false,
					},
					{
						SLI: "sql_statements",
						Pass: []*keptn.SLOCriteria{
							{
								Criteria: []string{"=0%"},
							},
							{
								Criteria: []string{"<100"},
							},
						},
						Warning: []*keptn.SLOCriteria{
							{
								Criteria: []string{"<+5%", ">-5%"},
							},
						},
						Weight: 1,
						KeySLI: true,
					},
				},
				TotalScore: &keptn.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			ExpectedError: nil,
		},
		{
			Name: "Simple SLO file without comparison spec",
			SLOFileContent: `---
spec_version: '1.0'
filter:
  id: "<prometheus_scrape_job_id>"
objectives:
  - sli: responseTime95
    pass:
      - criteria:
          - "<=+10%"
      - criteria:
          - "<200"
    warning:
      - criteria:
          - "<+15%"
          - ">-8%"
          - "<500"
  - sli: security_vulnerabilities
    weight: 2
    pass:
      - criteria:
          - "=0"
  - sli: sql_statements
    key_sli: true
    pass:
      - criteria:
          - "=0%"
      - criteria:
          - "<100"
    warning:
      - criteria:
          - "<+5%"
          - ">-5%"
total_score:
  pass: "90%"
  warning: 75%`,
			ExpectedSLO: &keptn.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter: map[string]string{
					"id": "<prometheus_scrape_job_id>",
				},
				Comparison: &keptn.SLOComparison{
					CompareWith:               "single_result",
					IncludeResultWithScore:    "all",
					NumberOfComparisonResults: 1,
					AggregateFunction:         "avg",
				},
				Objectives: []*keptn.SLO{
					{
						SLI: "responseTime95",
						Pass: []*keptn.SLOCriteria{
							{
								Criteria: []string{"<=+10%"},
							},
							{
								Criteria: []string{"<200"},
							},
						},
						Warning: []*keptn.SLOCriteria{
							{
								Criteria: []string{"<+15%", ">-8%", "<500"},
							},
						},
						Weight: 1,
						KeySLI: false,
					},
					{
						SLI: "security_vulnerabilities",
						Pass: []*keptn.SLOCriteria{
							{
								Criteria: []string{"=0"},
							},
						},
						Weight: 2,
						KeySLI: false,
					},
					{
						SLI: "sql_statements",
						Pass: []*keptn.SLOCriteria{
							{
								Criteria: []string{"=0%"},
							},
							{
								Criteria: []string{"<100"},
							},
						},
						Warning: []*keptn.SLOCriteria{
							{
								Criteria: []string{"<+5%", ">-5%"},
							},
						},
						Weight: 1,
						KeySLI: true,
					},
				},
				TotalScore: &keptn.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			ExpectedError: nil,
		},
		{
			Name: "Issue 6096 SLO file",
			SLOFileContent: `---
spec_version: ""
filter: {}
comparison:
  compare_with: single_result
  include_result_with_score: pass
  number_of_comparison_results: 1
  aggregate_function: avg
objectives:
- null
- sli: srt
  displayName: ""
  pass: []
  warning: []
  weight: 1
  key_sli: false
total_score:
  pass: 90%
  warning: 75%`,
			ExpectedSLO: &keptn.ServiceLevelObjectives{
				SpecVersion: "",
				Filter:      map[string]string{},
				Comparison: &keptn.SLOComparison{
					CompareWith:               "single_result",
					IncludeResultWithScore:    "pass",
					NumberOfComparisonResults: 1,
					AggregateFunction:         "avg",
				},
				Objectives: []*keptn.SLO{
					{
						SLI:     "srt",
						Pass:    []*keptn.SLOCriteria{},
						Warning: []*keptn.SLOCriteria{},
						Weight:  1,
						KeySLI:  false,
					},
				},
				TotalScore: &keptn.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			ExpectedError: nil,
		},
		{
			Name: "Simple SLO file",
			SLOFileContent: `---
spec_version: '1.0'
filter:
  id: "<prometheus_scrape_job_id>"
comparison:
  compare_with: "single_result"
objectives:
  - sli: responseTime95
    pass:
      - criteria:
          - "<=+10%"
      - criteria:
          - "<200"
    warning:
      - criteria:
          - "<+15%"
          - ">-8%"
          - "<500"
  - sli: security_vulnerabilities
    weight: 2
    pass:
      - criteria:
          - "=0"
  - sli: sql_statements
    key_sli: true
    pass:
      - criteria:
          - "=0%"
      - criteria:
          - "<100"
    warning:
      - criteria:
          - "<+5%"
          - ">-5%"
total_score:
  pass: "90%"
  warning: 75%`,
			ExpectedSLO: &keptn.ServiceLevelObjectives{
				SpecVersion: "1.0",
				Filter: map[string]string{
					"id": "<prometheus_scrape_job_id>",
				},
				Comparison: &keptn.SLOComparison{
					CompareWith:               "single_result",
					IncludeResultWithScore:    "all",
					NumberOfComparisonResults: 3,
					AggregateFunction:         "avg",
				},
				Objectives: []*keptn.SLO{
					{
						SLI: "responseTime95",
						Pass: []*keptn.SLOCriteria{
							{
								Criteria: []string{"<=+10%"},
							},
							{
								Criteria: []string{"<200"},
							},
						},
						Warning: []*keptn.SLOCriteria{
							{
								Criteria: []string{"<+15%", ">-8%", "<500"},
							},
						},
						Weight: 1,
						KeySLI: false,
					},
					{
						SLI: "security_vulnerabilities",
						Pass: []*keptn.SLOCriteria{
							{
								Criteria: []string{"=0"},
							},
						},
						Weight: 2,
						KeySLI: false,
					},
					{
						SLI: "sql_statements",
						Pass: []*keptn.SLOCriteria{
							{
								Criteria: []string{"=0%"},
							},
							{
								Criteria: []string{"<100"},
							},
						},
						Warning: []*keptn.SLOCriteria{
							{
								Criteria: []string{"<+5%", ">-5%"},
							},
						},
						Weight: 1,
						KeySLI: true,
					},
				},
				TotalScore: &keptn.SLOScore{
					Pass:    "90%",
					Warning: "75%",
				},
			},
			ExpectedError: nil,
		},
		{
			Name:           "Invalid SLO file",
			SLOFileContent: `invalid`,
			ExpectedSLO:    nil,
			ExpectedError:  errors.New("yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `invalid` into keptn.ServiceLevelObjectives"),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			objectives, err := parseSLO([]byte(test.SLOFileContent))
			assert.EqualValues(t, test.ExpectedSLO, objectives)
			assert.EqualValues(t, test.ExpectedError, err)
		})
	}
}

func getStartEventWithCommitId(id string) cloudevents.Event {
	return cloudevents.Event{
		Context: &cloudevents.EventContextV1{
			Type:            keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName),
			Source:          types.URIRef{},
			ID:              "",
			Time:            nil,
			DataContentType: strutils.Stringp("application/json"),
			Extensions: map[string]interface{}{
				"shkeptncontext": "my-context",
				"gitcommitid":    id,
			},
		},
		DataEncoded: []byte(`{
    "project": "sockshop",
    "stage": "staging",
    "service": "carts",
    "testStrategy": "",
    "deploymentStrategy": "direct",
	"evaluation": {
		"timeframe": "5m"
    },
    "labels": {
      "testid": "12345",
      "buildnr": "build17",
      "runby": "JohnDoe"
    },
    "result": "pass"
  }`),
		DataBase64: false,
	}
}

func TestBuildSubscriptions(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []models.EventSubscription
	}{
		{
			name:  "valid subscriptions",
			input: "sh.keptn.event.evaluation.triggered,sh.keptn.event.get-sli.finished,sh.keptn.event.monitoring.configure",
			want: []models.EventSubscription{
				{
					Event:  "sh.keptn.event.evaluation.triggered",
					Filter: models.EventSubscriptionFilter{},
				},
				{
					Event:  "sh.keptn.event.get-sli.finished",
					Filter: models.EventSubscriptionFilter{},
				},
				{
					Event:  "sh.keptn.event.monitoring.configure",
					Filter: models.EventSubscriptionFilter{},
				},
			},
		},
		{
			name:  "invalid subscriptions #1",
			input: "sh.keptn.event.evaluation.triggered.sh.keptn.event.get-sli.finished.sh.keptn.event.monitoring.configure",
			want: []models.EventSubscription{
				{
					Event:  "sh.keptn.event.evaluation.triggered.sh.keptn.event.get-sli.finished.sh.keptn.event.monitoring.configure",
					Filter: models.EventSubscriptionFilter{},
				},
			},
		},
		{
			name:  "valid subscriptions #2",
			input: "sh.keptn.event.evaluation.triggered;sh.keptn.event.get-sli.finished;sh.keptn.event.monitoring.configure",
			want: []models.EventSubscription{
				{
					Event:  "sh.keptn.event.evaluation.triggered;sh.keptn.event.get-sli.finished;sh.keptn.event.monitoring.configure",
					Filter: models.EventSubscriptionFilter{},
				},
			},
		},
		{
			name:  "no subscriptions #1",
			input: "",
			want:  []models.EventSubscription{},
		},
		{
			name:  "no subscriptions #2",
			input: ",",
			want:  []models.EventSubscription{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildSubscriptions(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}
