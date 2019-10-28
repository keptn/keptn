package event_handler

import (
	"encoding/json"
	"errors"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"io/ioutil"
	"net/http"
	"net/url"
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

type EvaluateSLIHandler struct {
	Logger     *keptnutils.Logger
	Event      cloudevents.Event
	HTTPClient *http.Client
}

func (eh *EvaluateSLIHandler) HandleEvent() error {
	e := &keptnevents.InternalGetSLIDoneEventData{}

	err := eh.Event.DataAs(e)

	if err != nil {
		eh.Logger.Error("Could not parse event payload: " + err.Error())
		return err
	}
	// get results of previous evaluations from data store (mongodb-datastore.keptn-datastore.svc.cluster.local)
	queryString := "mongodb-datastore.keptn-datastore.svc.cluster.local/event?type=" + keptnevents.EvaluationDoneEventType + "&project=" + e.Project + "&stage=" + e.Stage + "&service=" + e.Service + "&pageSize=20"
	req, err := http.NewRequest("GET", queryString, nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := eh.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return errors.New("could not retrieve previous evaluation-done events")
	}

	previousEvents := &datastoreResult{}
	err = json.Unmarshal(body, previousEvents)

	if err != nil {
		return err
	}

	var evaluationDoneEvents []keptnevents.EvaluationDoneEventData

	for _, event := range previousEvents.Events {
		evaluationDoneEvents = append(evaluationDoneEvents, event.Data.(keptnevents.EvaluationDoneEventData))
	}

	// compare the results based on the evaluation strategy

	// send the evaluation-done-event
	return nil
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
