package fake

import (
	"encoding/json"
	"github.com/keptn/keptn/shipyard-controller/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type EventBroker struct {
	Server           *httptest.Server
	ReceivedEvents   []models.Event
	Test             *testing.T
	HandleEventFunc  func(meb *EventBroker, event *models.Event)
	VerificationFunc func(meb *EventBroker)
}

func NewEventBroker(test *testing.T, handleEventFunc func(meb *EventBroker, event *models.Event), verificationFunc func(meb *EventBroker)) *EventBroker {
	meb := &EventBroker{
		Server:           nil,
		ReceivedEvents:   []models.Event{},
		Test:             test,
		HandleEventFunc:  handleEventFunc,
		VerificationFunc: verificationFunc,
	}

	meb.Server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		b, _ := ioutil.ReadAll(request.Body)
		defer func() {
			_ = request.Body.Close()
		}()
		event := &models.Event{}

		_ = json.Unmarshal(b, event)
		meb.HandleEventFunc(meb, event)

	}))

	return meb
}
