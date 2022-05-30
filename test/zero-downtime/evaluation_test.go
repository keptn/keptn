package zero_downtime

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	testutils "github.com/keptn/keptn/test/go-tests"
	keptnkubeutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math/rand"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const zdShipyard = `apiVersion: "spec.keptn.sh/0.2.0"
kind: "Shipyard"
metadata:
  name: "shipyard-sockshop"
spec:
  stages:
    - name: "hardening"
      sequences:
        - name: "remediation"
          tasks:
            - name: "action"
            - name: "approval"
              properties:
                pass: "automatic"
                warning: "automatic"
            - name: "evaluation"
              properties:
                timeframe: "5m"
        - name: "evaluation"
          tasks:
            - name: "evaluation"
            - name: "approval"
              properties:
                pass: "automatic"
                warning: "automatic"`

const webhookYaml = `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: sh.keptn.event.action.triggered
      requests:
        - >-
          curl --request GET http://shipyard-controller:8080/v1/project
      subscriptionID: ${action-sub-id}
      sendFinished: true
      sendStarted: true`

const remediationYaml = `apiVersion: spec.keptn.sh/0.1.4
kind: Remediation
metadata:
  name: remediation-configuration
spec:
  remediations: 
  - problemType: "default"
    actionsOnOpen:
    - name: Execute webhook
      action: webhook
      description: Execute a nice webhook`

const sloYaml = `---
spec_version: '0.1.0'
comparison:
  compare_with: "single_result"
  include_result_with_score: "pass"
  aggregate_function: avg
objectives:
  - sli: test-metric
    pass:
      - criteria:
          - "<=4"
    warning:
      - criteria:
          - "<=5"
total_score:
  pass: "51"
  warning: "20"`

func TestMessWithResourceService(t *testing.T) {
	//images := []string{"0.15.1-dev.202205240824", "0.15.1-dev.202205240902"}
	images := []string{"nonff", "nonff2"}
	services := []string{"resource-service"}

	project := "a-resource-service-test"
	shipyardFile, err := testutils.CreateTmpShipyardFile(zdShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFile)

	t.Logf("Creating project %s", project)
	project, err = testutils.CreateProject(project, shipyardFile)
	require.Nil(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	for _, svc := range services {
		go func(service string) {
			err := updateImageOfService(ctx, t, service, images)
			if err != nil {
				t.Logf("%v", err)
			}
		}(svc)
	}

	resourceContent := "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPGptZXRlclRlc3RQbGFuIHZlcnNpb249IjEuMiIgcHJvcGVydGllcz0iNS4wIiBqbWV0ZXI9IjUuNCI+CiAgPGhhc2hUcmVlPgogICAgPFRlc3RQbGFuIGd1aWNsYXNzPSJUZXN0UGxhbkd1aSIgdGVzdGNsYXNzPSJUZXN0UGxhbiIgdGVzdG5hbWU9IlRlc3QgUGxhbiIgZW5hYmxlZD0idHJ1ZSI+CiAgICAgIDxzdHJpbmdQcm9wIG5hbWU9IlRlc3RQbGFuLmNvbW1lbnRzIj48L3N0cmluZ1Byb3A+CiAgICAgIDxib29sUHJvcCBuYW1lPSJUZXN0UGxhbi5mdW5jdGlvbmFsX21vZGUiPmZhbHNlPC9ib29sUHJvcD4KICAgICAgPGJvb2xQcm9wIG5hbWU9IlRlc3RQbGFuLnNlcmlhbGl6ZV90aHJlYWRncm91cHMiPmZhbHNlPC9ib29sUHJvcD4KICAgICAgPGVsZW1lbnRQcm9wIG5hbWU9IlRlc3RQbGFuLnVzZXJfZGVmaW5lZF92YXJpYWJsZXMiIGVsZW1lbnRUeXBlPSJBcmd1bWVudHMiIGd1aWNsYXNzPSJBcmd1bWVudHNQYW5lbCIgdGVzdGNsYXNzPSJBcmd1bWVudHMiIHRlc3RuYW1lPSJVc2VyIERlZmluZWQgVmFyaWFibGVzIiBlbmFibGVkPSJ0cnVlIj4KICAgICAgICA8Y29sbGVjdGlvblByb3AgbmFtZT0iQXJndW1lbnRzLmFyZ3VtZW50cyI+CiAgICAgICAgICA8ZWxlbWVudFByb3AgbmFtZT0iU0VSVkVSX1VSTCIgZWxlbWVudFR5cGU9IkFyZ3VtZW50Ij4KICAgICAgICAgICAgPHN0cmluZ1Byb3AgbmFtZT0iQXJndW1lbnQubmFtZSI+U0VSVkVSX1VSTDwvc3RyaW5nUHJvcD4KICAgICAgICAgICAgPHN0cmluZ1Byb3AgbmFtZT0iQXJndW1lbnQudmFsdWUiPmVjMi01NC0xNjQtMTY0LTEyMS5jb21wdXRlLTEuYW1hem9uYXdzLmNvbTwvc3RyaW5nUHJvcD4KICAgICAgICAgICAgPHN0cmluZ1Byb3AgbmFtZT0iQXJndW1lbnQubWV0YWRhdGEiPj08L3N0cmluZ1Byb3A+CiAgICAgICAgICA8L2VsZW1lbnRQcm9wPgogICAgICAgICAgPGVsZW1lbnRQcm9wIG5hbWU9IkRlZmF1bHRUaGlua1RpbWUiIGVsZW1lbnRUeXBlPSJBcmd1bWVudCI+CiAgICAgICAgICAgIDxzdHJpbmdQcm9wIG5hbWU9IkFyZ3VtZW50Lm5hbWUiPkRlZmF1bHRUaGlua1RpbWU8L3N0cmluZ1Byb3A+CiAgICAgICAgICAgIDxzdHJpbmdQcm9wIG5hbWU9IkFyZ3VtZW50LnZhbHVlIj4yNTA8L3N0cmluZ1Byb3A+CiAgICAgICAgICAgIDxzdHJpbmdQcm9wIG5hbWU9IkFyZ3VtZW50Lm1ldGFkYXRhIj49PC9zdHJpbmdQcm9wPgogICAgICAgICAgPC9lbGVtZW50UHJvcD4KICAgICAgICAgIDxlbGVtZW50UHJvcCBuYW1lPSJEVF9MVE4iIGVsZW1lbnRUeXBlPSJBcmd1bWVudCI+CiAgICAgICAgICAgIDxzdHJpbmdQcm9wIG5hbWU9IkFyZ3VtZW50Lm5hbWUiPkRUX0xUTjwvc3RyaW5nUHJvcD4KICAgICAgICAgICAgPHN0cmluZ1Byb3AgbmFtZT0iQXJndW1lbnQudmFsdWUiPlRlc3RKdW5lMDM8L3N0cmluZ1Byb3A+CiAgICAgICAgICAgIDxzdHJpbmdQcm9wIG5hbWU9IkFyZ3VtZW50Lm1ldGFkYXRhIj49PC9zdHJpbmdQcm9wPgogICAgICAgICAgPC9lbGVtZW50UHJvcD4KICAgICAgICAgIDxlbGVtZW50UHJvcCBuYW1lPSJTRVJWRVJfUE9SVCIgZWxlbWVudFR5cGU9IkFyZ3VtZW50Ij4KICAgICAgICAgICAgPHN0cmluZ1Byb3AgbmFtZT0iQXJndW1lbnQubmFtZSI+U0VSVkVSX1BPUlQ8L3N0cmluZ1Byb3A+CiAgICAgICAgICAgIDxzdHJpbmdQcm9wIG5hbWU9IkFyZ3VtZW50LnZhbHVlIj44MDwvc3RyaW5nUHJvcD4KICAgICAgICAgICAgPHN0cmluZ1Byb3AgbmFtZT0iQXJndW1lbnQubWV0YWRhdGEiPj08L3N0cmluZ1Byb3A+CiAgICAgICAgICA8L2VsZW1lbnRQcm9wPgogICAgICAgICAgPGVsZW1lbnRQcm9wIG5hbWU9Ikxvb3BDb3VudCIgZWxlbWVudFR5cGU9IkFyZ3VtZW50Ij4KICAgICAgICAgICAgPHN0cmluZ1Byb3AgbmFtZT0iQXJndW1lbnQubmFtZSI+TG9vcENvdW50PC9zdHJpbmdQcm9wPgogICAgICAgICAgICA8c3RyaW5nUHJvcCBuYW1lPSJBcmd1bWVudC52YWx1ZSI+MTAwMDwvc3RyaW5nUHJvcD4KICAgICAgICAgICAgPHN0cmluZ1Byb3AgbmFtZT0iQXJndW1lbnQubWV0YWRhdGEiPj08L3N0cmluZ1Byb3A+CiAgICAgICAgICA8L2VsZW1lbnRQcm9wPgogICAgICAgICAgPGVsZW1lbnRQcm9wIG5hbWU9IkNIRUNLX1BBVEgiIGVsZW1lbnRUeXBlPSJBcmd1bWVudCI+CiAgICAgICAgICAgIDxzdHJpbmdQcm9wIG5hbWU9IkFyZ3VtZW50Lm5hbWUiPkNIRUNLX1BBVEg8L3N0cmluZ1Byb3A+CiAgICAgICAgICAgIDxzdHJpbmdQcm9wIG5hbWU9IkFyZ3VtZW50LnZhbHVlIj4vPC9zdHJpbmdQcm9wPgogICAgICAgICAgICA8c3RyaW5nUHJvcCBuYW1lPSJBcmd1bWVudC5tZXRhZGF0YSI+PTwvc3RyaW5nUHJvcD4KICAgICAgICAgIDwvZWxlbWVudFByb3A+CiAgICAgICAgICA8ZWxlbWVudFByb3AgbmFtZT0iUFJPVE9DT0wiIGVsZW1lbnRUeXBlPSJBcmd1bWVudCI+CiAgICAgICAgICAgIDxzdHJpbmdQcm9wIG5hbWU9IkFyZ3VtZW50Lm5hbWUiPlBST1RPQ09MPC9zdHJpbmdQcm9wPgogICAgICAgICAgICA8c3RyaW5nUHJvcCBuYW1lPSJBcmd1bWVudC52YWx1ZSI+aHR0cDwvc3RyaW5nUHJvcD4KICAgICAgICAgICAgPHN0cmluZ1Byb3AgbmFtZT0iQXJndW1lbnQubWV0YWRhdGEiPj08L3N0cmluZ1Byb3A+CiAgICAgICAgICA8L2VsZW1lbnRQcm9wPgogICAgICAgICAgPGVsZW1lbnRQcm9wIG5hbWU9IlZVQ291bnQiIGVsZW1lbnRUeXBlPSJBcmd1bWVudCI+CiAgICAgICAgICAgIDxzdHJpbmdQcm9wIG5hbWU9IkFyZ3VtZW50Lm5hbWUiPlZVQ291bnQ8L3N0cmluZ1Byb3A+CiAgICAgICAgICAgIDxzdHJpbmdQcm9wIG5hbWU9IkFyZ3VtZW50LnZhbHVlIj4xMDwvc3RyaW5nUHJvcD4KICAgICAgICAgICAgPHN0cmluZ1Byb3AgbmFtZT0iQXJndW1lbnQubWV0YWRhdGEiPj08L3N0cmluZ1Byb3A+CiAgICAgICAgICA8L2VsZW1lbnRQcm9wPgogICAgICAgIDwvY29sbGVjdGlvblByb3A+CiAgICAgIDwvZWxlbWVudFByb3A+CiAgICAgIDxzdHJpbmdQcm9wIG5hbWU9IlRlc3RQbGFuLnVzZXJfZGVmaW5lX2NsYXNzcGF0aCI+PC9zdHJpbmdQcm9wPgogICAgPC9UZXN0UGxhbj4KICAgIDxoYXNoVHJlZT4KICAgICAgPFRocmVhZEdyb3VwIGd1aWNsYXNzPSJUaHJlYWRHcm91cEd1aSIgdGVzdGNsYXNzPSJUaHJlYWRHcm91cCIgdGVzdG5hbWU9IlRocmVhZCBHcm91cCIgZW5hYmxlZD0idHJ1ZSI+CiAgICAgICAgPHN0cmluZ1Byb3AgbmFtZT0iVGhyZWFkR3JvdXAub25fc2FtcGxlX2Vycm9yIj5jb250aW51ZTwvc3RyaW5nUHJvcD4KICAgICAgICA8ZWxlbWVudFByb3AgbmFtZT0iVGhyZWFkR3JvdXAubWFpbl9jb250cm9sbGVyIiBlbGVtZW50VHlwZT0iTG9vcENvbnRyb2xsZXIiIGd1aWNsYXNzPSJMb29wQ29udHJvbFBhbmVsIiB0ZXN0Y2xhc3M9Ikxvb3BDb250cm9sbGVyIiB0ZXN0bmFtZT0iTG9vcCBDb250cm9sbGVyIiBlbmFibGVkPSJ0cnVlIj4KICAgICAgICAgIDxib29sUHJvcCBuYW1lPSJMb29wQ29udHJvbGxlci5jb250aW51ZV9mb3JldmVyIj5mYWxzZTwvYm9vbFByb3A+CiAgICAgICAgICA8c3RyaW5nUHJvcCBuYW1lPSJMb29wQ29udHJvbGxlci5sb29wcyI+JHtfX1AoTG9vcENvdW50LCR7TG9vcENvdW50fSl9PC9zdHJpbmdQcm9wPgogICAgICAgIDwvZWxlbWVudFByb3A+CiAgICAgICAgPHN0cmluZ1Byb3AgbmFtZT0iVGhyZWFkR3JvdXAubnVtX3RocmVhZHMiPiR7X19QKFZVQ291bnQsJHtWVUNvdW50fSl9PC9zdHJpbmdQcm9wPgogICAgICAgIDxzdHJpbmdQcm9wIG5hbWU9IlRocmVhZEdyb3VwLnJhbXBfdGltZSI+MTwvc3RyaW5nUHJvcD4KICAgICAgICA8bG9uZ1Byb3AgbmFtZT0iVGhyZWFkR3JvdXAuc3RhcnRfdGltZSI+MTUzNjA2NDUxNzAwMDwvbG9uZ1Byb3A+CiAgICAgICAgPGxvbmdQcm9wIG5hbWU9IlRocmVhZEdyb3VwLmVuZF90aW1lIj4xNTM2MDY0NTE3MDAwPC9sb25nUHJvcD4KICAgICAgICA8Ym9vbFByb3AgbmFtZT0iVGhyZWFkR3JvdXAuc2NoZWR1bGVyIj5mYWxzZTwvYm9vbFByb3A+CiAgICAgICAgPHN0cmluZ1Byb3AgbmFtZT0iVGhyZWFkR3JvdXAuZHVyYXRpb24iPjwvc3RyaW5nUHJvcD4KICAgICAgICA8c3RyaW5nUHJvcCBuYW1lPSJUaHJlYWRHcm91cC5kZWxheSI+PC9zdHJpbmdQcm9wPgogICAgICAgIDxib29sUHJvcCBuYW1lPSJUaHJlYWRHcm91cC5zYW1lX3VzZXJfb25fbmV4dF9pdGVyYXRpb24iPnRydWU8L2Jvb2xQcm9wPgogICAgICA8L1RocmVhZEdyb3VwPgogICAgICA8aGFzaFRyZWU+CiAgICAgICAgPENvb2tpZU1hbmFnZXIgZ3VpY2xhc3M9IkNvb2tpZVBhbmVsIiB0ZXN0Y2xhc3M9IkNvb2tpZU1hbmFnZXIiIHRlc3RuYW1lPSJIVFRQIENvb2tpZSBNYW5hZ2VyIiBlbmFibGVkPSJ0cnVlIj4KICAgICAgICAgIDxjb2xsZWN0aW9uUHJvcCBuYW1lPSJDb29raWVNYW5hZ2VyLmNvb2tpZXMiLz4KICAgICAgICAgIDxib29sUHJvcCBuYW1lPSJDb29raWVNYW5hZ2VyLmNsZWFyRWFjaEl0ZXJhdGlvbiI+ZmFsc2U8L2Jvb2xQcm9wPgogICAgICAgICAgPGJvb2xQcm9wIG5hbWU9IkNvb2tpZU1hbmFnZXIuY29udHJvbGxlZEJ5VGhyZWFkR3JvdXAiPmZhbHNlPC9ib29sUHJvcD4KICAgICAgICA8L0Nvb2tpZU1hbmFnZXI+CiAgICAgICAgPGhhc2hUcmVlLz4KICAgICAgICA8SGVhZGVyTWFuYWdlciBndWljbGFzcz0iSGVhZGVyUGFuZWwiIHRlc3RjbGFzcz0iSGVhZGVyTWFuYWdlciIgdGVzdG5hbWU9IkhUVFAgSGVhZGVyIE1hbmFnZXIiIGVuYWJsZWQ9InRydWUiPgogICAgICAgICAgPGNvbGxlY3Rpb25Qcm9wIG5hbWU9IkhlYWRlck1hbmFnZXIuaGVhZGVycyI+CiAgICAgICAgICAgIDxlbGVtZW50UHJvcCBuYW1lPSIiIGVsZW1lbnRUeXBlPSJIZWFkZXIiPgogICAgICAgICAgICAgIDxzdHJpbmdQcm9wIG5hbWU9IkhlYWRlci5uYW1lIj5DYWNoZS1Db250cm9sPC9zdHJpbmdQcm9wPgogICAgICAgICAgICAgIDxzdHJpbmdQcm9wIG5hbWU9IkhlYWRlci52YWx1ZSI+bm8tY2FjaGU8L3N0cmluZ1Byb3A+CiAgICAgICAgICAgIDwvZWxlbWVudFByb3A+CiAgICAgICAgICAgIDxlbGVtZW50UHJvcCBuYW1lPSIiIGVsZW1lbnRUeXBlPSJIZWFkZXIiPgogICAgICAgICAgICAgIDxzdHJpbmdQcm9wIG5hbWU9IkhlYWRlci5uYW1lIj5Db250ZW50LVR5cGU8L3N0cmluZ1Byb3A+CiAgICAgICAgICAgICAgPHN0cmluZ1Byb3AgbmFtZT0iSGVhZGVyLnZhbHVlIj5hcHBsaWNhdGlvbi9qc29uPC9zdHJpbmdQcm9wPgogICAgICAgICAgICA8L2VsZW1lbnRQcm9wPgogICAgICAgICAgICA8ZWxlbWVudFByb3AgbmFtZT0iIiBlbGVtZW50VHlwZT0iSGVhZGVyIj4KICAgICAgICAgICAgICA8c3RyaW5nUHJvcCBuYW1lPSJIZWFkZXIubmFtZSI+anNvbjwvc3RyaW5nUHJvcD4KICAgICAgICAgICAgICA8c3RyaW5nUHJvcCBuYW1lPSJIZWFkZXIudmFsdWUiPnRydWU8L3N0cmluZ1Byb3A+CiAgICAgICAgICAgIDwvZWxlbWVudFByb3A+CiAgICAgICAgICA8L2NvbGxlY3Rpb25Qcm9wPgogICAgICAgIDwvSGVhZGVyTWFuYWdlcj4KICAgICAgICA8aGFzaFRyZWUvPgogICAgICAgIDxCZWFuU2hlbGxQcmVQcm9jZXNzb3IgZ3VpY2xhc3M9IlRlc3RCZWFuR1VJIiB0ZXN0Y2xhc3M9IkJlYW5TaGVsbFByZVByb2Nlc3NvciIgdGVzdG5hbWU9IlNldCBEeW5hdHJhY2UgSGVhZGVycyIgZW5hYmxlZD0idHJ1ZSI+CiAgICAgICAgICA8c3RyaW5nUHJvcCBuYW1lPSJmaWxlbmFtZSI+PC9zdHJpbmdQcm9wPgogICAgICAgICAgPHN0cmluZ1Byb3AgbmFtZT0icGFyYW1ldGVycyI+bG9hZC5qbXg8L3N0cmluZ1Byb3A+CiAgICAgICAgICA8Ym9vbFByb3AgbmFtZT0icmVzZXRJbnRlcnByZXRlciI+ZmFsc2U8L2Jvb2xQcm9wPgogICAgICAgICAgPHN0cmluZ1Byb3AgbmFtZT0ic2NyaXB0Ij4KCmltcG9ydCBvcmcuYXBhY2hlLmptZXRlci51dGlsLkpNZXRlclV0aWxzOwppbXBvcnQgb3JnLmFwYWNoZS5qbWV0ZXIucHJvdG9jb2wuaHR0cC5jb250cm9sLkhlYWRlck1hbmFnZXI7CmltcG9ydCBqYXZhLmlvOwppbXBvcnQgamF2YS51dGlsOwoKLy8gLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLQovLyBHZW5lcmF0ZSB0aGUgeC1keW5hdHJhY2UtdGVzdCBoZWFkZXIKLy8gLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLQpTdHJpbmcgTFROPUpNZXRlclV0aWxzLmdldFByb3BlcnR5KCZxdW90O0RUX0xUTiZxdW90Oyk7CmlmKChMVE4gPT0gbnVsbCkgfHwgKExUTi5sZW5ndGgoKSA9PSAwKSkgewogICAgaWYodmFycyAhPSBudWxsKSB7CiAgICAgICAgTFROID0gdmFycy5nZXQoJnF1b3Q7RFRfTFROJnF1b3Q7KTsKICAgIH0KfQppZihMVE4gPT0gbnVsbCkgTFROID0gJnF1b3Q7Tm9UZXN0TmFtZSZxdW90OzsKClN0cmluZyBMU04gPSAoYnNoLmFyZ3MubGVuZ3RoICZndDsgMCkgPyBic2guYXJnc1swXSA6ICZxdW90O1Rlc3QgU2NlbmFyaW8mcXVvdDs7ClN0cmluZyBUU04gPSBzYW1wbGVyLmdldE5hbWUoKTsKU3RyaW5nIFZVID0gY3R4LmdldFRocmVhZEdyb3VwKCkuZ2V0TmFtZSgpICsgY3R4LmdldFRocmVhZE51bSgpOwpTdHJpbmcgaGVhZGVyVmFsdWUgPSAmcXVvdDtMU049JnF1b3Q7KyBMU04gKyAmcXVvdDs7VFNOPSZxdW90OyArIFRTTiArICZxdW90OztMVE49JnF1b3Q7ICsgTFROICsgJnF1b3Q7O1ZVPSZxdW90OyArIFZVICsgJnF1b3Q7OyZxdW90OzsKCi8vIC0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0KLy8gU2V0IGhlYWRlcgovLyAtLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tCkhlYWRlck1hbmFnZXIgaG0gPSBzYW1wbGVyLmdldEhlYWRlck1hbmFnZXIoKTsKaG0ucmVtb3ZlSGVhZGVyTmFtZWQoJnF1b3Q7eC1keW5hdHJhY2UtdGVzdCZxdW90Oyk7CmhtLmFkZChuZXcgb3JnLmFwYWNoZS5qbWV0ZXIucHJvdG9jb2wuaHR0cC5jb250cm9sLkhlYWRlcigmcXVvdDt4LWR5bmF0cmFjZS10ZXN0JnF1b3Q7LCBoZWFkZXJWYWx1ZSkpOwoKICAgICAgICAgIDwvc3RyaW5nUHJvcD4KICAgICAgICA8L0JlYW5TaGVsbFByZVByb2Nlc3Nvcj4KICAgICAgICA8aGFzaFRyZWUvPgogICAgICAgIDxIVFRQU2FtcGxlclByb3h5IGd1aWNsYXNzPSJIdHRwVGVzdFNhbXBsZUd1aSIgdGVzdGNsYXNzPSJIVFRQU2FtcGxlclByb3h5IiB0ZXN0bmFtZT0iaG9tZXBhZ2UiIGVuYWJsZWQ9InRydWUiPgogICAgICAgICAgPGVsZW1lbnRQcm9wIG5hbWU9IkhUVFBzYW1wbGVyLkFyZ3VtZW50cyIgZWxlbWVudFR5cGU9IkFyZ3VtZW50cyIgZ3VpY2xhc3M9IkhUVFBBcmd1bWVudHNQYW5lbCIgdGVzdGNsYXNzPSJBcmd1bWVudHMiIGVuYWJsZWQ9InRydWUiPgogICAgICAgICAgICA8Y29sbGVjdGlvblByb3AgbmFtZT0iQXJndW1lbnRzLmFyZ3VtZW50cyIvPgogICAgICAgICAgPC9lbGVtZW50UHJvcD4KICAgICAgICAgIDxzdHJpbmdQcm9wIG5hbWU9IkhUVFBTYW1wbGVyLmRvbWFpbiI+JHtfX1AoU0VSVkVSX1VSTCwke1NFUlZFUl9VUkx9KX08L3N0cmluZ1Byb3A+CiAgICAgICAgICA8c3RyaW5nUHJvcCBuYW1lPSJIVFRQU2FtcGxlci5wb3J0Ij4ke19fUChTRVJWRVJfUE9SVCwke1NFUlZFUl9QT1JUfSl9PC9zdHJpbmdQcm9wPgogICAgICAgICAgPHN0cmluZ1Byb3AgbmFtZT0iSFRUUFNhbXBsZXIucHJvdG9jb2wiPiR7X19QKFBST1RPQ09MLCR7UFJPVE9DT0x9KX08L3N0cmluZ1Byb3A+CiAgICAgICAgICA8c3RyaW5nUHJvcCBuYW1lPSJIVFRQU2FtcGxlci5jb250ZW50RW5jb2RpbmciPjwvc3RyaW5nUHJvcD4KICAgICAgICAgIDxzdHJpbmdQcm9wIG5hbWU9IkhUVFBTYW1wbGVyLnBhdGgiPi88L3N0cmluZ1Byb3A+CiAgICAgICAgICA8c3RyaW5nUHJvcCBuYW1lPSJIVFRQU2FtcGxlci5tZXRob2QiPkdFVDwvc3RyaW5nUHJvcD4KICAgICAgICAgIDxib29sUHJvcCBuYW1lPSJIVFRQU2FtcGxlci5mb2xsb3dfcmVkaXJlY3RzIj50cnVlPC9ib29sUHJvcD4KICAgICAgICAgIDxib29sUHJvcCBuYW1lPSJIVFRQU2FtcGxlci5hdXRvX3JlZGlyZWN0cyI+ZmFsc2U8L2Jvb2xQcm9wPgogICAgICAgICAgPGJvb2xQcm9wIG5hbWU9IkhUVFBTYW1wbGVyLnVzZV9rZWVwYWxpdmUiPnRydWU8L2Jvb2xQcm9wPgogICAgICAgICAgPGJvb2xQcm9wIG5hbWU9IkhUVFBTYW1wbGVyLkRPX01VTFRJUEFSVF9QT1NUIj5mYWxzZTwvYm9vbFByb3A+CiAgICAgICAgICA8Ym9vbFByb3AgbmFtZT0iSFRUUFNhbXBsZXIuQlJPV1NFUl9DT01QQVRJQkxFX01VTFRJUEFSVCI+dHJ1ZTwvYm9vbFByb3A+CiAgICAgICAgICA8c3RyaW5nUHJvcCBuYW1lPSJIVFRQU2FtcGxlci5lbWJlZGRlZF91cmxfcmUiPjwvc3RyaW5nUHJvcD4KICAgICAgICAgIDxzdHJpbmdQcm9wIG5hbWU9IkhUVFBTYW1wbGVyLmNvbm5lY3RfdGltZW91dCI+PC9zdHJpbmdQcm9wPgogICAgICAgICAgPHN0cmluZ1Byb3AgbmFtZT0iSFRUUFNhbXBsZXIucmVzcG9uc2VfdGltZW91dCI+PC9zdHJpbmdQcm9wPgogICAgICAgIDwvSFRUUFNhbXBsZXJQcm94eT4KICAgICAgICA8aGFzaFRyZWUvPgogICAgICAgIDxDb25zdGFudFRpbWVyIGd1aWNsYXNzPSJDb25zdGFudFRpbWVyR3VpIiB0ZXN0Y2xhc3M9IkNvbnN0YW50VGltZXIiIHRlc3RuYW1lPSJEZWZhdWx0IFRoaW5rIFRpbWUiIGVuYWJsZWQ9InRydWUiPgogICAgICAgICAgPHN0cmluZ1Byb3AgbmFtZT0iQ29uc3RhbnRUaW1lci5kZWxheSI+e19fUChUaGlua1RpbWUsJHtEZWZhdWx0VGhpbmtUaW1lfSl9PC9zdHJpbmdQcm9wPgogICAgICAgIDwvQ29uc3RhbnRUaW1lcj4KICAgICAgICA8aGFzaFRyZWUvPgogICAgICA8L2hhc2hUcmVlPgogICAgICA8UmVzdWx0Q29sbGVjdG9yIGd1aWNsYXNzPSJWaWV3UmVzdWx0c0Z1bGxWaXN1YWxpemVyIiB0ZXN0Y2xhc3M9IlJlc3VsdENvbGxlY3RvciIgdGVzdG5hbWU9IlZpZXcgUmVzdWx0cyBUcmVlIiBlbmFibGVkPSJ0cnVlIj4KICAgICAgICA8Ym9vbFByb3AgbmFtZT0iUmVzdWx0Q29sbGVjdG9yLmVycm9yX2xvZ2dpbmciPmZhbHNlPC9ib29sUHJvcD4KICAgICAgICA8b2JqUHJvcD4KICAgICAgICAgIDxuYW1lPnNhdmVDb25maWc8L25hbWU+CiAgICAgICAgICA8dmFsdWUgY2xhc3M9IlNhbXBsZVNhdmVDb25maWd1cmF0aW9uIj4KICAgICAgICAgICAgPHRpbWU+dHJ1ZTwvdGltZT4KICAgICAgICAgICAgPGxhdGVuY3k+dHJ1ZTwvbGF0ZW5jeT4KICAgICAgICAgICAgPHRpbWVzdGFtcD50cnVlPC90aW1lc3RhbXA+CiAgICAgICAgICAgIDxzdWNjZXNzPnRydWU8L3N1Y2Nlc3M+CiAgICAgICAgICAgIDxsYWJlbD50cnVlPC9sYWJlbD4KICAgICAgICAgICAgPGNvZGU+dHJ1ZTwvY29kZT4KICAgICAgICAgICAgPG1lc3NhZ2U+dHJ1ZTwvbWVzc2FnZT4KICAgICAgICAgICAgPHRocmVhZE5hbWU+dHJ1ZTwvdGhyZWFkTmFtZT4KICAgICAgICAgICAgPGRhdGFUeXBlPnRydWU8L2RhdGFUeXBlPgogICAgICAgICAgICA8ZW5jb2Rpbmc+ZmFsc2U8L2VuY29kaW5nPgogICAgICAgICAgICA8YXNzZXJ0aW9ucz50cnVlPC9hc3NlcnRpb25zPgogICAgICAgICAgICA8c3VicmVzdWx0cz50cnVlPC9zdWJyZXN1bHRzPgogICAgICAgICAgICA8cmVzcG9uc2VEYXRhPmZhbHNlPC9yZXNwb25zZURhdGE+CiAgICAgICAgICAgIDxzYW1wbGVyRGF0YT5mYWxzZTwvc2FtcGxlckRhdGE+CiAgICAgICAgICAgIDx4bWw+ZmFsc2U8L3htbD4KICAgICAgICAgICAgPGZpZWxkTmFtZXM+ZmFsc2U8L2ZpZWxkTmFtZXM+CiAgICAgICAgICAgIDxyZXNwb25zZUhlYWRlcnM+ZmFsc2U8L3Jlc3BvbnNlSGVhZGVycz4KICAgICAgICAgICAgPHJlcXVlc3RIZWFkZXJzPmZhbHNlPC9yZXF1ZXN0SGVhZGVycz4KICAgICAgICAgICAgPHJlc3BvbnNlRGF0YU9uRXJyb3I+ZmFsc2U8L3Jlc3BvbnNlRGF0YU9uRXJyb3I+CiAgICAgICAgICAgIDxzYXZlQXNzZXJ0aW9uUmVzdWx0c0ZhaWx1cmVNZXNzYWdlPmZhbHNlPC9zYXZlQXNzZXJ0aW9uUmVzdWx0c0ZhaWx1cmVNZXNzYWdlPgogICAgICAgICAgICA8YXNzZXJ0aW9uc1Jlc3VsdHNUb1NhdmU+MDwvYXNzZXJ0aW9uc1Jlc3VsdHNUb1NhdmU+CiAgICAgICAgICAgIDxieXRlcz50cnVlPC9ieXRlcz4KICAgICAgICAgICAgPHRocmVhZENvdW50cz50cnVlPC90aHJlYWRDb3VudHM+CiAgICAgICAgICA8L3ZhbHVlPgogICAgICAgIDwvb2JqUHJvcD4KICAgICAgICA8c3RyaW5nUHJvcCBuYW1lPSJmaWxlbmFtZSI+PC9zdHJpbmdQcm9wPgogICAgICA8L1Jlc3VsdENvbGxlY3Rvcj4KICAgICAgPGhhc2hUcmVlLz4KICAgICAgPFJlc3VsdENvbGxlY3RvciBndWljbGFzcz0iU3VtbWFyeVJlcG9ydCIgdGVzdGNsYXNzPSJSZXN1bHRDb2xsZWN0b3IiIHRlc3RuYW1lPSJTdW1tYXJ5IFJlcG9ydCIgZW5hYmxlZD0idHJ1ZSI+CiAgICAgICAgPGJvb2xQcm9wIG5hbWU9IlJlc3VsdENvbGxlY3Rvci5lcnJvcl9sb2dnaW5nIj5mYWxzZTwvYm9vbFByb3A+CiAgICAgICAgPG9ialByb3A+CiAgICAgICAgICA8bmFtZT5zYXZlQ29uZmlnPC9uYW1lPgogICAgICAgICAgPHZhbHVlIGNsYXNzPSJTYW1wbGVTYXZlQ29uZmlndXJhdGlvbiI+CiAgICAgICAgICAgIDx0aW1lPnRydWU8L3RpbWU+CiAgICAgICAgICAgIDxsYXRlbmN5PnRydWU8L2xhdGVuY3k+CiAgICAgICAgICAgIDx0aW1lc3RhbXA+dHJ1ZTwvdGltZXN0YW1wPgogICAgICAgICAgICA8c3VjY2Vzcz50cnVlPC9zdWNjZXNzPgogICAgICAgICAgICA8bGFiZWw+dHJ1ZTwvbGFiZWw+CiAgICAgICAgICAgIDxjb2RlPnRydWU8L2NvZGU+CiAgICAgICAgICAgIDxtZXNzYWdlPnRydWU8L21lc3NhZ2U+CiAgICAgICAgICAgIDx0aHJlYWROYW1lPnRydWU8L3RocmVhZE5hbWU+CiAgICAgICAgICAgIDxkYXRhVHlwZT50cnVlPC9kYXRhVHlwZT4KICAgICAgICAgICAgPGVuY29kaW5nPmZhbHNlPC9lbmNvZGluZz4KICAgICAgICAgICAgPGFzc2VydGlvbnM+dHJ1ZTwvYXNzZXJ0aW9ucz4KICAgICAgICAgICAgPHN1YnJlc3VsdHM+dHJ1ZTwvc3VicmVzdWx0cz4KICAgICAgICAgICAgPHJlc3BvbnNlRGF0YT5mYWxzZTwvcmVzcG9uc2VEYXRhPgogICAgICAgICAgICA8c2FtcGxlckRhdGE+ZmFsc2U8L3NhbXBsZXJEYXRhPgogICAgICAgICAgICA8eG1sPmZhbHNlPC94bWw+CiAgICAgICAgICAgIDxmaWVsZE5hbWVzPnRydWU8L2ZpZWxkTmFtZXM+CiAgICAgICAgICAgIDxyZXNwb25zZUhlYWRlcnM+ZmFsc2U8L3Jlc3BvbnNlSGVhZGVycz4KICAgICAgICAgICAgPHJlcXVlc3RIZWFkZXJzPmZhbHNlPC9yZXF1ZXN0SGVhZGVycz4KICAgICAgICAgICAgPHJlc3BvbnNlRGF0YU9uRXJyb3I+ZmFsc2U8L3Jlc3BvbnNlRGF0YU9uRXJyb3I+CiAgICAgICAgICAgIDxzYXZlQXNzZXJ0aW9uUmVzdWx0c0ZhaWx1cmVNZXNzYWdlPnRydWU8L3NhdmVBc3NlcnRpb25SZXN1bHRzRmFpbHVyZU1lc3NhZ2U+CiAgICAgICAgICAgIDxhc3NlcnRpb25zUmVzdWx0c1RvU2F2ZT4wPC9hc3NlcnRpb25zUmVzdWx0c1RvU2F2ZT4KICAgICAgICAgICAgPGJ5dGVzPnRydWU8L2J5dGVzPgogICAgICAgICAgICA8c2VudEJ5dGVzPnRydWU8L3NlbnRCeXRlcz4KICAgICAgICAgICAgPHVybD50cnVlPC91cmw+CiAgICAgICAgICAgIDx0aHJlYWRDb3VudHM+dHJ1ZTwvdGhyZWFkQ291bnRzPgogICAgICAgICAgICA8aWRsZVRpbWU+dHJ1ZTwvaWRsZVRpbWU+CiAgICAgICAgICAgIDxjb25uZWN0VGltZT50cnVlPC9jb25uZWN0VGltZT4KICAgICAgICAgIDwvdmFsdWU+CiAgICAgICAgPC9vYmpQcm9wPgogICAgICAgIDxzdHJpbmdQcm9wIG5hbWU9ImZpbGVuYW1lIj48L3N0cmluZ1Byb3A+CiAgICAgIDwvUmVzdWx0Q29sbGVjdG9yPgogICAgICA8aGFzaFRyZWUvPgogICAgPC9oYXNoVHJlZT4KICA8L2hhc2hUcmVlPgo8L2ptZXRlclRlc3RQbGFuPgo="
	newResourceContent := "LS0tCmFwaVZlcnNpb246IHYxCmtpbmQ6IFBlcnNpc3RlbnRWb2x1bWVDbGFpbQptZXRhZGF0YToKICBjcmVhdGlvblRpbWVzdGFtcDogbnVsbAogIG5hbWU6IGNvbmZpZ3VyYXRpb24tdm9sdW1lCiAgbmFtZXNwYWNlOiBrZXB0bgogIGxhYmVsczoKICAgIGFwcC5rdWJlcm5ldGVzLmlvL25hbWU6IGNvbmZpZ3VyYXRpb24tdm9sdW1lCiAgICBhcHAua3ViZXJuZXRlcy5pby9pbnN0YW5jZToga2VwdG4KICAgIGFwcC5rdWJlcm5ldGVzLmlvL3BhcnQtb2Y6IGtlcHRuLWtlcHRuCiAgICBhcHAua3ViZXJuZXRlcy5pby9jb21wb25lbnQ6IGNvbnRyb2wtcGxhbmUKc3BlYzoKICBhY2Nlc3NNb2RlczoKICAtIFJlYWRXcml0ZU9uY2UKICByZXNvdXJjZXM6CiAgICByZXF1ZXN0czoKICAgICAgc3RvcmFnZTogMTAwTWkKc3RhdHVzOiB7fQ=="
	nrUploadResources := 50
	nrThreadsPerUpdate := 10

	var nrSuccessfulRequests int32
	nrTotalRequests := nrUploadResources * nrThreadsPerUpdate

	defer func() {
		t.Logf("Succesful requests: %d/%d", nrSuccessfulRequests, nrTotalRequests)
	}()

	resourceUri := "my-resource.txt"

	createResourceRequest := models.Resources{
		Resources: []*models.Resource{
			{
				ResourceContent: resourceContent,
				ResourceURI:     &resourceUri,
			},
		},
	}

	updateResourceRequest := models.Resources{
		Resources: []*models.Resource{
			{
				ResourceContent: newResourceContent,
				ResourceURI:     &resourceUri,
			},
		},
	}

	for i := 0; i < nrUploadResources; i++ {
		for j := 0; j < nrThreadsPerUpdate; j++ {
			var resourceReq models.Resources
			if j%2 == 0 {
				resourceReq = createResourceRequest
			} else {
				resourceReq = updateResourceRequest
			}
			go func(r models.Resources) {
				<-time.After(time.Duration(rand.Intn(5000)) * time.Millisecond)
				t.Logf("Creating a new resource for project %s", project)
				resp, err := testutils.ApiPOSTRequest("/configuration-service/v1/project/"+project+"/resource", r, 0)
				assert.Nil(t, err)
				if err != nil {
					t.Logf("error: %s", err.Error())
				}
				if resp != nil {
					t.Logf("Completed request with status code %s", resp.Response().Status)
					assert.Equal(t, 201, resp.Response().StatusCode)
					if resp.Response().StatusCode == 201 {
						atomic.AddInt32(&nrSuccessfulRequests, 1)
					}
				}
			}(resourceReq)
		}
		<-time.After(5 * time.Second)
	}

	cancel()
}

func TestSequences(t *testing.T) {
	images := []string{"0.15.1-dev.202205240824", "0.15.1-dev.202205240902"}
	services := []string{"api-service", "shipyard-controller", "resource-service", "lighthouse-service", "approval-service", "webhook-service", "remediation-service", "mongodb-datastore"}
	//services := []string{"mongodb-datastore"}

	project := "a-zd-test"
	stage := "hardening"
	service := "myservice"

	nrSequences := 30

	shipyardFile, err := testutils.CreateTmpShipyardFile(zdShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFile)

	t.Logf("Creating project %s", project)
	project, err = testutils.CreateProject(project, shipyardFile)
	require.Nil(t, err)

	t.Logf("creating service %s", service)
	_, err = testutils.ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", service, project))

	taskTypes := []string{"action"}

	t.Logf("Setting up subscription")
	webhookYamlWithSubscriptionIDs := webhookYaml
	webhookYamlWithSubscriptionIDs = getWebhookYamlWithSubscriptionIDs(t, taskTypes, project, webhookYamlWithSubscriptionIDs)

	// wait some time to make sure the webhook service has pulled the updated subscription
	<-time.After(20 * time.Second) // sorry :(

	t.Logf("Adding webhook")
	// now, let's add an webhook.yaml file to our service
	webhookFilePath, err := testutils.CreateTmpFile("webhook.yaml", webhookYamlWithSubscriptionIDs)
	require.Nil(t, err)
	defer func() {
		err := os.Remove(webhookFilePath)
		if err != nil {
			t.Logf("Could not delete tmp file: %s", err.Error())
		}

	}()

	t.Log("Adding webhook.yaml to our service")
	_, err = testutils.ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --service=%s --resource=%s --resourceUri=webhook/webhook.yaml --all-stages", project, service, webhookFilePath))

	remediationFilePath, err := testutils.CreateTmpFile("remediation.yaml", remediationYaml)
	t.Log("Adding remediation.yaml to our service")
	_, err = testutils.ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --service=%s --resource=%s --resourceUri=remediation.yaml --all-stages", project, service, remediationFilePath))

	t.Log("deleting lighthouse configmap from previous test run")
	testutils.ExecuteCommandf("kubectl delete configmap -n %s lighthouse-config-%s", testutils.GetKeptnNameSpaceFromEnv(), project)

	t.Log("adding SLI provider")
	_, err = testutils.ExecuteCommand(fmt.Sprintf("kubectl create configmap -n %s lighthouse-config-%s --from-literal=sli-provider=my-sli-provider", testutils.GetKeptnNameSpaceFromEnv(), project))
	require.Nil(t, err)

	sloFilePath, err := testutils.CreateTmpFile("slo.yaml", sloYaml)
	t.Log("Adding slo.yaml to our service")
	_, err = testutils.ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --service=%s --resource=%s --resourceUri=slo.yaml --all-stages", project, service, sloFilePath))

	ctx, cancel := context.WithCancel(context.Background())

	for _, svc := range services {
		go func(service string) {
			err := updateImageOfService(ctx, t, service, images)
			if err != nil {
				t.Logf("%v", err)
			}
		}(svc)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go startSLIRetrieval(t, project, stage, service)
	executeSequences(nrSequences, project, stage, service)

	t.Log("Checking if all sequences are finished")
	require.Eventually(t, func() bool {
		states := &models.SequenceStates{}
		resp, err := testutils.ApiGETRequest("/controlPlane/v1/sequence/"+project+"?state=finished", 3)
		if err != nil {
			return false
		}
		err = resp.ToJSON(states)
		if err != nil {
			return false
		}

		if states.TotalCount != int64(nrSequences) {
			t.Logf("Finished %d/%d triggered sequences", states.TotalCount, nrSequences)
			return false
		}
		t.Logf("All sequences completed!")
		return true
	}, 10*time.Minute, 10*time.Second)

	cancel()
}

func executeSequences(nrSequences int, project, stage, service string) {
	for i := 0; i < nrSequences; i++ {
		nrTriggered := 0
		go func() {
			//_, err := triggerEvaluation("podtatohead", "hardening", "helloservice")
			_, err := triggerRemediation(project, stage, service)
			if err != nil {
				nrTriggered++
			}
		}()

		<-time.After(3 * time.Second)
	}
}

func triggerEvaluation(projectName, stageName, serviceName string) (string, error) {
	cliResp, err := testutils.ExecuteCommand(fmt.Sprintf("keptn trigger evaluation --project=%s --stage=%s --service=%s --timeframe=5m", projectName, stageName, serviceName))

	if err != nil {
		return "", err
	}
	var keptnContext string
	split := strings.Split(cliResp, "\n")
	for _, line := range split {
		if strings.Contains(line, "ID of") {
			splitLine := strings.Split(line, ":")
			if len(splitLine) == 2 {
				keptnContext = strings.TrimSpace(splitLine[1])
			}
		}
	}
	return keptnContext, err
}

func triggerRemediation(projectName, stageName, serviceName string) (string, error) {
	source := "golang-test"
	eventData := keptnv2.EventData{}
	eventType := keptnv2.GetTriggeredEventType(stageName + ".remediation")
	eventData.SetProject(projectName)
	eventData.SetService(serviceName)
	eventData.SetStage(stageName)

	resp, err := testutils.ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
		Contenttype:        "application/json",
		Data:               eventData,
		ID:                 uuid.NewString(),
		Shkeptnspecversion: "0.2.0",
		Source:             &source,
		Specversion:        "1.0",
		Type:               &eventType,
	}, 0)

	if err != nil {
		return "", err
	}

	eventContext := &models.EventContext{}
	err = resp.ToJSON(eventContext)
	if err != nil {
		return "", err
	}
	return *eventContext.KeptnContext, nil
}

func updateImageOfService(ctx context.Context, t *testing.T, service string, images []string) error {
	clientset, err := keptnkubeutils.GetClientset(false)

	if err != nil {
		return err
	}

	i := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			nextImage := images[i%len(images)]
			get, err := clientset.AppsV1().Deployments(testutils.GetKeptnNameSpaceFromEnv()).Get(context.TODO(), service, v1.GetOptions{})
			if err != nil {
				break
			}

			imageWithTag := get.Spec.Template.Spec.Containers[0].Image
			split := strings.Split(imageWithTag, ":")
			updatedImage := fmt.Sprintf("%s:%s", split[0], nextImage)

			get.Spec.Template.Spec.Containers[0].Image = updatedImage

			t.Logf("upgrading %s to %s", service, updatedImage)
			_, err = clientset.AppsV1().Deployments(testutils.GetKeptnNameSpaceFromEnv()).Update(context.TODO(), get, v1.UpdateOptions{})
			if err != nil {
				break
			}

			require.Eventually(t, func() bool {
				pods, err := clientset.CoreV1().Pods(testutils.GetKeptnNameSpaceFromEnv()).List(context.TODO(), v1.ListOptions{LabelSelector: "app.kubernetes.io/name=" + service})
				if err != nil {
					return false
				}

				if int32(len(pods.Items)) != 1 {
					// make sure only one pod is running
					return false
				}

				for _, item := range pods.Items {
					if len(item.Spec.Containers) == 0 {
						continue
					}
					if item.Spec.Containers[0].Image == updatedImage {
						return true
					}
				}
				return false
			}, 3*time.Minute, 10*time.Second)
			<-time.After(5 * time.Second)
			i++
		}
	}
}

func startSLIRetrieval(t *testing.T, project, stage, service string) {
	retrievedSLIs := map[string]bool{}
	var err error
	for {
		<-time.After(2 * time.Second)
		retrievedSLIs, err = reportSLIValues(retrievedSLIs, project, stage, service)
		if err != nil {
			t.Logf("Error while SLI retrieval: %v", err)
		}

	}
}

func reportSLIValues(retrievedSLIs map[string]bool, project string, stage string, service string) (map[string]bool, error) {
	resp, err := testutils.ApiGETRequest(fmt.Sprintf("/mongodb-datastore/event?project=%s&stage=%s&service=%s&type=sh.keptn.event.get-sli.triggered", project, stage, service), 3)
	if err != nil {
		return retrievedSLIs, err
	}
	events := &models.Events{}
	if err := resp.ToJSON(events); err != nil {
		return retrievedSLIs, err
	}

	if len(events.Events) == 0 {
		return retrievedSLIs, nil
	}

	sliFinishedEventType := keptnv2.GetFinishedEventType(keptnv2.GetSLITaskName)
	source := "golang-test"
	for _, sliTriggeredEvent := range events.Events {
		if retrievedSLIs[sliTriggeredEvent.Shkeptncontext] {
			continue
		}
		retrievedSLIs[sliTriggeredEvent.Shkeptncontext] = true
		_, err := testutils.ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
			Contenttype: "application/json",
			Data: keptnv2.GetSLIFinishedEventData{
				EventData: keptnv2.EventData{
					Project: project,
					Stage:   stage,
					Service: service,
					Status:  keptnv2.StatusSucceeded,
					Result:  keptnv2.ResultPass,
				},
				GetSLI: keptnv2.GetSLIFinished{
					IndicatorValues: []*keptnv2.SLIResult{
						{
							Metric:  "test-metric",
							Value:   1,
							Success: true,
						},
					},
				},
			},
			ID:                 uuid.NewString(),
			Shkeptnspecversion: "0.2.0",
			Source:             &source,
			Specversion:        "1.0",
			Shkeptncontext:     sliTriggeredEvent.Shkeptncontext,
			Triggeredid:        sliTriggeredEvent.ID,
			Type:               &sliFinishedEventType,
		}, 0)

		if err != nil {
			continue
		}
	}
	return retrievedSLIs, nil
}

func getWebhookYamlWithSubscriptionIDs(t *testing.T, taskTypes []string, projectName string, webhookYamlWithSubscriptionIDs string) string {
	for _, taskType := range taskTypes {
		eventType := keptnv2.GetTriggeredEventType(taskType)
		if strings.HasSuffix(taskType, "-finished") {
			eventType = keptnv2.GetFinishedEventType(strings.TrimSuffix(taskType, "-finished"))
		}
		subscriptionID, err := testutils.CreateSubscription(t, "webhook-service", models.EventSubscription{
			Event: eventType,
			Filter: models.EventSubscriptionFilter{
				Projects: []string{projectName},
			},
		})
		require.Nil(t, err)

		subscriptionPlaceholder := fmt.Sprintf("${%s-sub-id}", taskType)
		webhookYamlWithSubscriptionIDs = strings.Replace(webhookYamlWithSubscriptionIDs, subscriptionPlaceholder, subscriptionID, -1)
	}
	return webhookYamlWithSubscriptionIDs
}
