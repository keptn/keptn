package fake

import (
	"encoding/json"
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type EventBroker struct {
	Server           *httptest.Server
	ReceivedEvents   []keptnmodels.KeptnContextExtendedCE
	Test             *testing.T
	HandleEventFunc  func(meb *EventBroker, event *keptnmodels.KeptnContextExtendedCE)
	VerificationFunc func(meb *EventBroker)
}

func NewEventBroker(test *testing.T, handleEventFunc func(meb *EventBroker, event *keptnmodels.KeptnContextExtendedCE), verificationFunc func(meb *EventBroker)) *EventBroker {
	meb := &EventBroker{
		Server:           nil,
		ReceivedEvents:   []keptnmodels.KeptnContextExtendedCE{},
		Test:             test,
		HandleEventFunc:  handleEventFunc,
		VerificationFunc: verificationFunc,
	}

	meb.Server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		b, _ := ioutil.ReadAll(request.Body)
		defer func() {
			_ = request.Body.Close()
		}()
		event := &keptnmodels.KeptnContextExtendedCE{}

		_ = json.Unmarshal(b, event)
		meb.HandleEventFunc(meb, event)

	}))

	return meb
}
