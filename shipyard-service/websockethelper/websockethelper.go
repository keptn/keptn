package websockethelper

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/url"

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
}

// LogData represents log data
type LogData struct {
	Message   string `json:"message"`
	Terminate bool   `json:"terminate"`
	LogLevel  string `json:"loglevel"`
}

// ConnectionData stores ChannelInfo and Success data
type ConnectionData struct {
	ChannelInfo ChannelInfo `json:"channelInfo"`
}

// ChannelInfo stores a token and a channelID used for opening the websocket
type ChannelInfo struct {
	Token     string `json:"token"`
	ChannelID string `json:"channelID"`
}

// OpenWS opens a websocket
func OpenWS(connData ConnectionData, apiEndPoint url.URL) (*websocket.Conn, *http.Response, error) {

	wsEndPoint := apiEndPoint
	wsEndPoint.Scheme = "wss"

	header := http.Header{}
	header.Add("Token", connData.ChannelInfo.Token)
	header.Add("Keptn-Ws-Channel-Id", connData.ChannelInfo.ChannelID)

	dialer := websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	return dialer.Dial(wsEndPoint.String(), header)
}

// WriteWSLog writes the log event to the websocket
func WriteWSLog(ws *websocket.Conn, logEvent cloudevents.Event, message string, terminate bool, logLevel string) error {

	logData := LogData{
		Message:   message,
		Terminate: terminate,
		LogLevel:  logLevel,
	}

	logEvent.Data = logData
	return ws.WriteJSON(logEvent)
}
