// This file is safe to edit. Once it exists it will not be overwritten

package utils

import (
	"context"

	keptn "github.com/keptn/go-utils/pkg/lib"

	cloudevents "github.com/cloudevents/sdk-go"
)

// PostToEventBroker makes a post request to the eventbroker
func PostToEventBroker(event cloudevents.Event) (*cloudevents.Event, error) {

	url, err := keptn.GetServiceEndpoint("EVENTBROKER_URI")
	if err != nil {
		return nil, err
	}
	t, err := cloudevents.NewHTTPTransport(
		cloudevents.WithTarget(url.String()),
		cloudevents.WithEncoding(cloudevents.HTTPStructuredV02),
	)

	if err != nil {
		return nil, err
	}

	c, err := cloudevents.NewClient(t)
	if err != nil {
		return nil, err
	}
	_, sent, err := c.Send(context.Background(), event)
	return sent, err
}
