package websockethelper

import (
	"log"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
)

const logCEString = `{
    "cloudEventsVersion":"0.1",
    "contentType":"application/json",
    "data":{"message":"InfoMsg","terminate":true},
    "eventID":"af326f24-8705-4332-b32b-affcb62f3567",
    "eventTime":"2019-03-12T17:18:14.187682+01:00",
    "eventType":"sh.keptn.events.log",
    "extensions":null,
    "source":"https://github.com/keptn/keptn/cli#wstest"
}`

func TestClient(t *testing.T) {

	credentialmanager.MockCreds = true

	ws, _, err := OpenWS("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjaGFubmVsSWQiOiI2MDY0Njc5Zi0wNmJiLTRmMDEtOWIyNS1lMjM5Yjc4YmMzYTMiLCJpYXQiOjE1NTI0Nzk3NzgsImV4cCI6MTYzODg3OTc3OH0.ZDphJPxXrJtk4Qyk77t1nafNSzxBXBZmvcGTR7Vz064", "channelId")

	if err != nil {
		t.Errorf("An error occured %v", err)
	}

	err = ws.WriteMessage(websocket.TextMessage, ([]byte(logCEString)))
	if err != nil {
		log.Fatal(err)
	}

	err = PrintWSContent(ws, true)
	if err != nil {
		t.Errorf("An error occured %v", err)
	}
}
