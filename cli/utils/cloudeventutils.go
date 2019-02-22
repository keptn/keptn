package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/knative/pkg/cloudevents"
)

// Send creates a request including the X-Keptn-Signature and sends the data
// struct to the provided target. It returns error if there was an
// issue sending the event, otherwise nil means the event was accepted.
func Send(target string, apiToken string, builder cloudevents.Builder, data interface{}, overrides ...cloudevents.SendContext) error {

	req, err := builder.Build(target, data, overrides...)

	bodyBytes, err := ioutil.ReadAll(req.Body)
	// Restore the io.ReadCloser to its original state
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	if err != nil {
		fmt.Println("Unable to read body", err)
	}

	mac := hmac.New(sha1.New, []byte(apiToken))
	mac.Write(bodyBytes)
	signatureVal := mac.Sum(nil)
	sha1Hash := "sha1=" + fmt.Sprintf("%x", signatureVal)

	// Add signature header
	req.Header.Set("X-Keptn-Signature", sha1Hash)

	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println(resp.StatusCode)
		return nil
	}
	return fmt.Errorf("error sending cloudevent: %s", status(resp))
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
