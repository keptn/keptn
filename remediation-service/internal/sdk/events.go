package sdk

import (
	"context"
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	httpprotocol "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	"time"
)

const DefaultHTTPEventEndpoint = "http://localhost:8081/event"
const MAX_SEND_RETRIES = 3

//go:generate moq  -pkg fake -out ./fake/event_sender_mock.go . EventSender
type EventSender interface {
	SendEvent(event cloudevents.Event) error
}

type EventReceiver interface {
	StartReceiver(ctx context.Context, fn interface{}) error
}

type HTTPEventSender struct {
	// EventsEndpoint is the http endpoint the events are sent to
	EventsEndpoint string
	// Client is an implementation of the cloudevents.Client interface
	Client cloudevents.Client
}

func NewHTTPEventSender(ceClient cloudevents.Client) *HTTPEventSender {
	c := &HTTPEventSender{
		EventsEndpoint: DefaultHTTPEventEndpoint,
		Client:         ceClient,
	}
	return c
}

// SendEvent sends a CloudEvent
func (httpSender HTTPEventSender) SendEvent(event cloudevents.Event) error {
	ctx := cloudevents.ContextWithTarget(context.Background(), httpSender.EventsEndpoint)
	ctx = cloudevents.WithEncodingStructured(ctx)

	var result protocol.Result
	for i := 0; i <= MAX_SEND_RETRIES; i++ {
		result = httpSender.Client.Send(ctx, event)
		httpResult, ok := result.(*httpprotocol.Result)
		if ok {
			if httpResult.StatusCode >= 200 && httpResult.StatusCode < 300 {
				return nil
			}
			<-time.After(keptn.GetExpBackoffTime(i + 1))
		} else if cloudevents.IsUndelivered(result) {
			<-time.After(keptn.GetExpBackoffTime(i + 1))
		} else {
			return nil
		}
	}
	return errors.New("Failed to send cloudevent: " + result.Error())
}
