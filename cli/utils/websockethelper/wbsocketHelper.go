package websockethelper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
)

type myCloudEvent struct {
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

// OpenWS opens a websocket
func OpenWS(token string, channelID string) (*websocket.Conn, *http.Response, error) {

	endPoint, _, err := credentialmanager.GetCreds()
	if err != nil {
		return nil, nil, err
	}

	wsEndPoint := endPoint
	wsEndPoint.Scheme = "ws"

	header := http.Header{}
	header.Add("Token", token)
	header.Add("x-keptn-ws-channel-id", channelID)
	return websocket.DefaultDialer.Dial(wsEndPoint.String(), header)
}

// readCE reads a cloud event
func readAndPrintCE(ws *websocket.Conn, verboseLogging bool) (interface{}, error) {
	ws.SetReadDeadline(time.Now().Add(time.Minute))
	for {
		messageType, message, err := ws.ReadMessage()
		if messageType == 1 { // 1.. textmessage
			var messageCE myCloudEvent
			dec := json.NewDecoder(strings.NewReader(string(message)))
			if err := dec.Decode(&messageCE); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			if printCE(messageCE, verboseLogging) {
				return nil, ws.Close()
			}

		}
		if err != nil {
			log.Println("read: ", err)
			return nil, err
		}

	}
	return nil, nil
}

func printCE(ce myCloudEvent, verboseLogging bool) bool {
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

// PrintWSContent prints received cloud events
func PrintWSContent(ws *websocket.Conn, verboseLogging bool) error {
	ceData, err := readAndPrintCE(ws, verboseLogging)
	if err != nil || ceData == nil {
		return err
	}
	return nil
}
