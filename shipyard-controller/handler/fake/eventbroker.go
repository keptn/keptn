package fake

import (
	"encoding/json"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type EventBroker struct {
	Server           *httptest.Server
	ReceivedEvents   []apimodels.KeptnContextExtendedCE
	Test             *testing.T
	HandleEventFunc  func(meb *EventBroker, event *apimodels.KeptnContextExtendedCE)
	VerificationFunc func(meb *EventBroker)
}

func NewEventBroker(test *testing.T, handleEventFunc func(meb *EventBroker, event *apimodels.KeptnContextExtendedCE), verificationFunc func(meb *EventBroker)) *EventBroker {
	meb := &EventBroker{
		Server:           nil,
		ReceivedEvents:   []apimodels.KeptnContextExtendedCE{},
		Test:             test,
		HandleEventFunc:  handleEventFunc,
		VerificationFunc: verificationFunc,
	}

	meb.Server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		b, _ := ioutil.ReadAll(request.Body)
		defer func() {
			_ = request.Body.Close()
		}()
		event := &apimodels.KeptnContextExtendedCE{}

		_ = json.Unmarshal(b, event)
		meb.HandleEventFunc(meb, event)

	}))

	return meb
}
