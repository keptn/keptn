package event_handler

import (
	"testing"

	"github.com/stretchr/testify/assert"

	keptn "github.com/keptn/go-utils/pkg/lib"
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
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			objectives, err := parseSLO([]byte(test.SLOFileContent))
			assert.EqualValues(t, test.ExpectedSLO, objectives)
			assert.EqualValues(t, test.ExpectedError, err)
		})
	}
}
