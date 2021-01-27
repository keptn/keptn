package fake

import (
	"encoding/json"
	"github.com/keptn/keptn/shipyard-controller/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockEventBroker struct {
	Server           *httptest.Server
	ReceivedEvents   []models.Event
	Test             *testing.T
	HandleEventFunc  func(meb *MockEventBroker, event *models.Event)
	VerificationFunc func(meb *MockEventBroker)
}

func (meb *MockEventBroker) handleEvent(event *models.Event) {
	meb.HandleEventFunc(meb, event)
}

func NewMockEventbroker(test *testing.T, handleEventFunc func(meb *MockEventBroker, event *models.Event), verificationFunc func(meb *MockEventBroker)) *MockEventBroker {
	meb := &MockEventBroker{
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
