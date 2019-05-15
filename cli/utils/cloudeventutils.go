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
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
)

type WebsocketDescription struct {
	ChannelID string `json:"channelId"`
	Token     string `json:"token"`
}

type RespData struct {
	Desc WebsocketDescription `json:"websocketChannel"`
}

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
	}
	t.Client = &http.Client{Timeout: timeout * time.Second, Transport: tr}

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
	if m, ok := msg.(*cloudeventshttp.Message); ok {
		mac := hmac.New(sha1.New, []byte(apiToken))
		mac.Write(m.Body)
		signatureVal := mac.Sum(nil)
		sha1Hash := "sha1=" + fmt.Sprintf("%x", signatureVal)

		// Add signature header
		usedContext = cloudeventshttp.ContextWithHeader(usedContext, "X-Keptn-Signature", sha1Hash)
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
