package websockethelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/gorilla/websocket"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
)

// MyCloudEvent represents a keptn cloud event
type MyCloudEvent struct {
	CloudEventsVersion string          `json:"cloudEventsVersion"`
	ContentType        string          `json:"contentType"`
	Data               json.RawMessage `json:"data"`
	EventID            string          `json:"eventID"`
	EventTime          string          `json:"eventTime"`
	EventType          string          `json:"eventType"`
	Type               string          `json:"type"`
	Source             string          `json:"source"`
}

// LogData represents log data
type LogData struct {
	Message   string `json:"message"`
	Terminate bool   `json:"terminate"`
	LogLevel  string `json:"loglevel"`
}

// incompleteCE is a helper type for unmarshalling the CE data
type incompleteCE struct {
	ConnData ConnectionData `json:"data"`
}

// ConnectionData stores ChannelInfo and Success data
type ConnectionData struct {
	ChannelInfo ChannelInfo `json:"channelInfo"`
	Success     bool        `json:"success"`
}

// ChannelInfo stores a token and a channelID used for opening the websocket
type ChannelInfo struct {
	Token     string `json:"token"`
	ChannelID string `json:"channelId"`
}

// PrintWSContent opens a websocket using the passed connection data and
// prints status data
func PrintWSContent(responseCE *cloudevents.Event, verbose bool) error {

	ceData := &incompleteCE{}
	err := responseCE.DataAs(ceData)
	if err != nil {
		return err
	}
	connData := ceData.ConnData

	err = validateConnectionData(connData)
	if err != nil {
		return err
	}

	ws, _, err := openWS(connData)
	if err != nil {
		fmt.Println("Opening websocket failed")
		return err
	}
	defer ws.Close()

	return readAndPrintCE(ws, verbose)
}

func validateConnectionData(connData ConnectionData) error {
	if !connData.Success || connData.ChannelInfo.Token == "" && connData.ChannelInfo.ChannelID == "" {
		return errors.New("Could not open websocket because Token or Channel ID might be missing or request received unsuccessful status")
	}
	return nil
}

// openWS opens a websocket
func openWS(connData ConnectionData) (*websocket.Conn, *http.Response, error) {

	endPoint, _, err := credentialmanager.GetCreds()
	if err != nil {
		return nil, nil, err
	}

	wsEndPoint := endPoint
	wsEndPoint.Scheme = "ws"

	header := http.Header{}
	header.Add("Token", connData.ChannelInfo.Token)
	header.Add("x-keptn-ws-channel-id", connData.ChannelInfo.ChannelID)
	return websocket.DefaultDialer.Dial(wsEndPoint.String(), header)
}

// readAndPrintCE reads a cloud event from the websocket
func readAndPrintCE(ws *websocket.Conn, verboseLogging bool) error {
	ws.SetReadDeadline(time.Now().Add(time.Minute))
	for {
		messageType, message, err := ws.ReadMessage()
		if messageType == 1 { // 1.. textmessage
			var messageCE MyCloudEvent
			dec := json.NewDecoder(strings.NewReader(string(message)))
			if err := dec.Decode(&messageCE); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			if printCE(messageCE, verboseLogging) {
				return nil
			}
		}
		if err != nil {
			log.Println("read: ", err)
			return err
		}

	}
	return nil
}

func printCE(ce MyCloudEvent, verboseLogging bool) bool {
	var log LogData
	if err := json.Unmarshal(ce.Data, &log); err != nil {
		fmt.Println("JSON unmarshalling error. LogData format expected.")
		//return nil, err
	}
	switch ce.Type {
	case "sh.keptn.events.log":
		printLogLevel(log, verboseLogging)
		return log.Terminate
	default:
		fmt.Println("type of event could not be processed")
	}
	return true
}

func printLogLevel(log LogData, verboseLogging bool) {
	if log.LogLevel == "DEBUG" && verboseLogging {
		fmt.Println(log.Message)
	} else if log.LogLevel != "DEBUG" {
		fmt.Println(log.Message)
	}
}
