package utils

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
)

const timeout = 60

// Send creates a request including the X-Keptn-Signature and sends the data
// struct to the provided target. It returns the obtained http.Response.
func Send(url url.URL, event cloudevents.Event, apiToken string) (*cloudevents.Event, error) {
	ec := event.Context.AsV02()
	if ec.Time == nil || ec.Time.IsZero() {
		ec.Time = &types.Timestamp{Time: time.Now()}
		event.Context = ec
	}

	t, err := cloudeventshttp.New(
		cloudeventshttp.WithTarget(url.String()),
		cloudeventshttp.WithEncoding(cloudeventshttp.StructuredV02),
	)
	// Reset Client because we need TLS
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext:     ResolveXipIoWithContext,
	}
	t.Client = &http.Client{Timeout: timeout * time.Second, Transport: tr}

	if err != nil {
		return nil, err
	}

	c, err := client.New(t)

	if err != nil {
		return nil, err
	}

	// Add signature header
	usedContext := cloudeventshttp.ContextWithHeader(context.Background(), "x-token", apiToken)
	return c.Send(usedContext, event)
}

// status is a helper method to read the response of the target.
func status(resp *http.Response) string {
	status := resp.Status
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Status[%s] error reading response body: %v", status, err)
	}
	return fmt.Sprintf("Status[%s] %s", status, body)
}
