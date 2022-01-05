package go_tests

import (
	"encoding/json"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"
	"time"
)

const evaluationDonePayload = `{
  "type": "sh.keptn.events.evaluation-done",
  "specversion": "0.2",
  "source": "lighthouse-service",
  "contenttype": "application/json",
  "data": {
    "project": "legacy-project",
    "stage": "hardening",
    "service": "legacy-service",
    "evaluationdetails":{
      "indicatorResults":[
        {
          "score":0,
          "status":"failed",
          "targets":[
            {
              "criteria":"<=800",
              "targetValue":800,
              "violated":true
            },
            {
              "criteria":"<=+10%",
              "targetValue":549.1967956487127,
              "violated":true
            },
            {
              "criteria":"<600",
              "targetValue":600,
              "violated":true
            }
          ],
          "value":{
            "metric":"response_time_p95",
            "success":true,
            "value":1002.6278552658177
          }
        }
      ],
      "result":"fail",
      "score":0,
      "sloFileContent":"LS0tDQpzcGVjX3ZlcnNpb246ICcxLjAnDQpjb21wYXJpc29uOg0KICBjb21wYXJlX3dpdGg6ICJzaW5nbGVfcmVzdWx0Ig0KICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyINCiAgYWdncmVnYXRlX2Z1bmN0aW9uOiBhdmcNCm9iamVjdGl2ZXM6DQogIC0gc2xpOiByZXNwb25zZV90aW1lX3A5NQ0KICAgIHBhc3M6ICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNTAwKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQ0KICAgICAgICAgIC0gIjw2MDAiICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcg0KICAgIHdhcm5pbmc6ICAgICAjIGlmIHRoZSByZXNwb25zZSB0aW1lIGlzIGJlbG93IDgwMG1zLCB0aGUgcmVzdWx0IHNob3VsZCBiZSBhIHdhcm5pbmcNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD04MDAiDQp0b3RhbF9zY29yZToNCiAgcGFzczogIjkwJSINCiAgd2FybmluZzogNzUl",
      "timeEnd":"2019-11-18T11:29:36Z",
      "timeStart":"2019-11-18T11:21:06Z"
    }
  }
}`

const evaluationInvalidatedEvent = `{
  "type": "sh.keptn.event.evaluation.invalidated",
  "specversion": "1.0",
  "source": "travis-ci",
  "contenttype": "application/json",
  "triggeredid": "$TRIGGERED_ID",
  "shkeptncontext": "$KEPTN_CONTEXT",
  "data": {
    "project": "legacy-project",
    "stage": "hardening",
    "service": "legacy-service"
  }
}`

func Test_QualityGates_BackwardsCompatibility(t *testing.T) {
	evaluationDoneEvent := &models.KeptnContextExtendedCE{}

	err := json.Unmarshal([]byte(evaluationDonePayload), evaluationDoneEvent)
	require.Nil(t, err)

	resp, err := ApiPOSTRequest("/v1/event", *evaluationDoneEvent, 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	context := &models.EventContext{}
	err = resp.ToJSON(context)
	require.Nil(t, err)
	require.NotNil(t, context.KeptnContext)

	eventsQuery := "/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?filter=shkeptncontext:" + *context.KeptnContext + "%20AND%20data.project:legacy-project&excludeInvalidated=true"

	var events *models.Events
	require.Eventually(t, func() bool {
		resp, err := ApiGETRequest(eventsQuery, 3)
		if err != nil {
			return false
		}

		events = &models.Events{}
		err = resp.ToJSON(events)
		if err != nil {
			return false
		}
		if events == nil {
			return false
		}
		if len(events.Events) == 0 {
			return false
		}
		return true
	}, 10*time.Second, 2*time.Second)

	require.Equal(t, "lighthouse-service", *events.Events[0].Source)
	require.Equal(t, keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName), *events.Events[0].Type)
	require.Equal(t, "1.0", events.Events[0].Specversion)

	evaluationFinishedData := &keptnv2.EvaluationFinishedEventData{}
	err = keptnv2.EventDataAs(*events.Events[0], evaluationFinishedData)

	require.Nil(t, err)
	require.NotEmpty(t, evaluationFinishedData.Evaluation)
	require.Equal(t, string(keptnv2.ResultFailed), evaluationFinishedData.Evaluation.Result)
	require.NotEmpty(t, evaluationFinishedData.Evaluation.IndicatorResults)
	require.NotEmpty(t, evaluationFinishedData.Evaluation.SLOFileContent)

	evaluationInvalidatedPayload := strings.ReplaceAll(evaluationInvalidatedEvent, "$TRIGGERED_ID", events.Events[0].ID)
	evaluationInvalidatedPayload = strings.ReplaceAll(evaluationInvalidatedPayload, "$KEPTN_CONTEXT", events.Events[0].Shkeptncontext)

	evaluationInvalidatedEvent := &models.KeptnContextExtendedCE{}
	err = json.Unmarshal([]byte(evaluationInvalidatedPayload), evaluationInvalidatedEvent)

	require.Nil(t, err)

	resp, err = ApiPOSTRequest("/v1/event", *evaluationInvalidatedEvent, 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	require.Eventually(t, func() bool {
		resp, err = ApiGETRequest(eventsQuery, 3)
		if err != nil {
			return false
		}

		events = &models.Events{}
		err = resp.ToJSON(events)
		if err != nil {
			return false
		}
		if events == nil {
			return false
		}
		if len(events.Events) != 0 {
			return false
		}
		return true
	}, 10*time.Second, 2*time.Second)
}
