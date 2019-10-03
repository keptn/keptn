package websockethelper

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/gorilla/websocket"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils"
)

// PrintWSContentCEResponse opens a websocket using the passed
// connection data (in form of a cloud event) and prints status data
func PrintWSContentCEResponse(responseCE *cloudevents.Event, apiEndPoint url.URL) error {

	connectionData := &keptnutils.ConnectionData{}
	err := responseCE.DataAs(connectionData)

	if err != nil {
		return err
	}
	return printWSContent(*connectionData, apiEndPoint)
}

// PrintWSContentByteResponse opens a websocket using the passed
// connection data (in form of a byte slice) and prints status data
func PrintWSContentByteResponse(response []byte, apiEndPoint url.URL) error {

	ceData := &keptnutils.IncompleteCE{}
	err := json.Unmarshal(response, ceData)
	if err != nil {
		return err
	}

	return printWSContent(ceData.ConnData, apiEndPoint)
}

func printWSContent(connData keptnutils.ConnectionData, apiEndPoint url.URL) error {

	err := validateConnectionData(connData)
	if err != nil {
		return err
	}

	ws, _, err := openWS(connData, apiEndPoint)
	if err != nil {
		fmt.Println("Opening websocket failed")
		return err
	}
	// PrintLogLevel(LogData{Message: "Websocket successfully opened", LogLevel: "DEBUG"}, loglevel)
	defer ws.Close()

	return readAndPrintCE(ws)
}

func validateConnectionData(connData keptnutils.ConnectionData) error {
	if connData.ChannelInfo.Token == "" && connData.ChannelInfo.ChannelID == "" {
		return errors.New("Could not open websocket because Token or Channel ID are missing")
	}
	return nil
}

// openWS opens a websocket
func openWS(connData keptnutils.ConnectionData, apiEndPoint url.URL) (*websocket.Conn, *http.Response, error) {

	wsEndPoint := apiEndPoint
	wsEndPoint.Scheme = "wss"

	header := http.Header{}
	header.Add("Token", connData.ChannelInfo.Token)
	header.Add("Keptn-Ws-Channel-Id", connData.ChannelInfo.ChannelID)
	dialer := websocket.DefaultDialer
	dialer.NetDial = utils.ResolveXipIo
	dialer.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, resp, err := dialer.Dial(wsEndPoint.String(), header)
	if err != nil {
		return nil, nil, err
	}
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })
	return conn, resp, err
}

// readAndPrintCE reads a cloud event from the websocket
func readAndPrintCE(ws *websocket.Conn) error {
	for {
		messageType, message, err := ws.ReadMessage()
		if messageType == 1 { // 1.. textmessage
			var messageCE keptnutils.MyCloudEvent

			dec := json.NewDecoder(strings.NewReader(string(message)))
			if err := dec.Decode(&messageCE); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			if printCE(messageCE) {
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

func printCE(ce keptnutils.MyCloudEvent) bool {
	var log keptnutils.LogData
	if err := json.Unmarshal(ce.Data, &log); err != nil {
		fmt.Println("JSON unmarshalling error. LogData format expected.")
		//return nil, err
	}
	switch ce.Type {
	case "sh.keptn.events.log":
		if strings.TrimSpace(log.Message) != "" {
			logging.PrintLogStringLevel(log.Message, log.LogLevel)
		}
		return log.Terminate
	default:
		fmt.Println("type of event could not be processed")
	}
	return true
}
