package pkg

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

const DefaultHTTPEventEndpoint = "http://localhost:8081/event"

//go:generate moq  -pkg fake -out ./fake/event_sender_mock.go . EventSender
type EventSender interface {
	SendEvent(event cloudevents.Event) error
}

type EventReceiver interface {
	StartReceiver(ctx context.Context, fn interface{}) error
}
