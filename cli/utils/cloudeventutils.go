package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type WebsocketDescription struct {
	ChannelID string `json:"channelId"`
	Token     string `json:"token"`
}

type RespData struct {
	Desc WebsocketDescription `json:"websocketChannel"`
}

// Send creates a request including the X-Keptn-Signature and sends the data
// struct to the provided target. It returns the obtained http.Response.
func Send(req *http.Request, apiToken string, desc *WebsocketDescription) (*http.Response, error) {

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
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 60 * time.Second, Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if desc != nil {
		// TODO: Check if there is a more elegant way to get the content-type
		if strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {

			gResp := RespData{}
			err = json.NewDecoder(resp.Body).Decode(&gResp)
			*desc = gResp.Desc
		}
	}

	return resp, err
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
