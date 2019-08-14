package websocketutil

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/keptn/keptn/cli/utils/websockethelper"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/gorilla/websocket"
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
	ShKeptnContext     string          `json:"shkeptncontext"`
}

// OpenWS opens a websocket
func OpenWS(connData websockethelper.ConnectionData, apiEndPoint url.URL) (*websocket.Conn, *http.Response, error) {

	wsEndPoint := apiEndPoint
	wsEndPoint.Scheme = "ws"

	header := http.Header{}
	header.Add("Token", connData.ChannelInfo.Token)

	dialer := websocket.DefaultDialer

	return dialer.Dial(wsEndPoint.String(), header)
}

// WriteWSLog writes the log event to the websocket
func WriteWSLog(ws *websocket.Conn, logEvent cloudevents.Event, message string, terminate bool, logLevel string) error {
	logData := websockethelper.LogData{
		Message:   message,
		Terminate: terminate,
		LogLevel:  logLevel,
	}
	logDataRaw, _ := json.Marshal(logData)

	var shkeptncontext string
	logEvent.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	messageCE := MyCloudEvent{
		CloudEventsVersion: logEvent.SpecVersion(),
		ContentType:        logEvent.DataContentType(),
		Data:               logDataRaw,
		EventID:            logEvent.ID(),
		EventTime:          logEvent.Time().String(),
		EventType:          logEvent.Type(),
		Type:               "sh.keptn.events.log",
		Source:             logEvent.Source(),
		ShKeptnContext:     shkeptncontext,
	}

	data, _ := json.Marshal(messageCE)
	return ws.WriteMessage(1, data) // websocket.TextMessage = 1; ws.WriteJSON not supported because keptn CLI does a ReadMessage
}
