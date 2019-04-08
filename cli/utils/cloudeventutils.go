package utils

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
)

type addHeader func(context.Context, transport.Message, string) context.Context

// AddXKeptnSignatureHeader adds the X-Keptn-Signature header in the provided context
func AddXKeptnSignatureHeader(c context.Context, msg transport.Message, apiToken string) context.Context {
	if m, ok := msg.(*cloudeventshttp.Message); ok {
		mac := hmac.New(sha1.New, []byte(apiToken))
		mac.Write(m.Body)
		signatureVal := mac.Sum(nil)
		sha1Hash := "sha1=" + fmt.Sprintf("%x", signatureVal)

		// Add signature header
		return cloudeventshttp.ContextWithHeader(c, "X-Keptn-Signature", sha1Hash)
	}

	Warning.Printf("Cannto add header")
	return c
}

// AddAuthorizationHeader adds the Autorization header in the provided context
func AddAuthorizationHeader(c context.Context, msg transport.Message, apiToken string) context.Context {

	return cloudeventshttp.ContextWithHeader(c, "Authorization", "Bearer "+apiToken)
}

// Send creates a request including the X-Keptn-Signature and sends the data
// struct to the provided target. It returns the obtained http.Response.
func Send(url url.URL, event cloudevents.Event, apiToken string, fn addHeader) (*cloudevents.Event, error) {

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
	}
	t.Client = &http.Client{Timeout: 60 * time.Second, Transport: tr}

	if err != nil {
		return nil, err
	}

	c, err := client.New(t)

	if err != nil {
		return nil, err
	}

	myCodec := &cloudeventshttp.Codec{
		Encoding:                   t.Encoding,
		DefaultEncodingSelectionFn: t.DefaultEncodingSelectionFn,
	}

	msg, err := myCodec.Encode(event)
	if err != nil {
		return nil, err
	}

	usedContext := context.Background()
	if fn != nil {
		usedContext = fn(usedContext, msg, apiToken)
	}

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
