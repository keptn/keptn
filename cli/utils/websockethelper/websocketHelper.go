package websockethelper

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils"
)

// PrintWSContentEventContext opens a websocket using the passed
// connection data and prints status data
func PrintWSContentEventContext(eventContext *apimodels.EventContext, apiEndPoint url.URL, useTSL bool) error {
	connectionData := &keptnutils.ConnectionData{EventContext: *eventContext}
	return printWSContent(*connectionData, apiEndPoint, useTSL)
}

func printWSContent(connData keptnutils.ConnectionData, apiEndPoint url.URL, useTSL bool) error {

	err := validateConnectionData(connData)
	if err != nil {
		return err
	}

	ws, _, err := openWS(connData, apiEndPoint, useTSL)
	if err != nil {
		fmt.Println("Opening websocket failed")
		return err
	}
	// PrintLogLevel(LogData{Message: "Websocket successfully opened", LogLevel: "DEBUG"}, loglevel)
	defer ws.Close()

	return readAndPrintCE(ws)
}

func validateConnectionData(connData keptnutils.ConnectionData) error {
	if *connData.EventContext.Token == "" && *connData.EventContext.KeptnContext == "" {
		return errors.New("Could not open websocket because Token or KeptnContext are missing")
	}
	return nil
}

// openWS opens a websocket
func openWS(connData keptnutils.ConnectionData, apiEndPoint url.URL, useTSL bool) (*websocket.Conn, *http.Response, error) {

	wsEndPoint := apiEndPoint
	if useTSL {
		wsEndPoint.Scheme = "wss"
	} else {
		wsEndPoint.Scheme = "ws"
	}
	header := http.Header{}
	header.Add("Token", *connData.EventContext.Token)
	header.Add("Keptn-Ws-Channel-Id", *connData.EventContext.KeptnContext)
	header.Add("Host", "api.keptn")

	dialer := websocket.DefaultDialer
	dialer.NetDial = utils.ResolveXipIo
	dialer.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, resp, err := dialer.Dial(wsEndPoint.String(), header)
	if err != nil {
		return nil, nil, err
	}
	conn.SetReadDeadline(time.Now().Add(readDeadline))
	return conn, resp, err
}

const readDeadline = 90 * time.Second

// readAndPrintCE reads a cloud event from the websocket
func readAndPrintCE(ws *websocket.Conn) error {
	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("Warning: Websocket connection timed out")
				return nil
			}
			return err
		}

		ws.SetReadDeadline(time.Now().Add(readDeadline))
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
