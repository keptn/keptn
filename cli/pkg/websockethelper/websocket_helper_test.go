package websockethelper

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	keptnutils "github.com/keptn/go-utils/pkg/lib"
)

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func TestSingleCERead(t *testing.T) {
	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	msg := "Test"

	sendCE(t, ws, msg, true, "DEBUG")

	r, w, old := beginRedirectStdOut()
	readAndPrintCE(ws)
	out := endRedirectStdOut(r, w, old)

	if strings.TrimSpace(out) != msg {
		t.Fatalf("Actual and expected output do not match")
	}
}

func TestDoubleCERead(t *testing.T) {
	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	msg := "Test"
	sendCE(t, ws, msg, false, "DEBUG")
	sendCE(t, ws, msg, true, "DEBUG")

	r, w, old := beginRedirectStdOut()
	readAndPrintCE(ws)
	out := endRedirectStdOut(r, w, old)

	if strings.TrimSpace(out) != msg+"\n"+msg {
		t.Fatalf("Actual and expected output do not match")
	}
}

func sendCE(t *testing.T, ws *websocket.Conn, msg string, terminate bool, logLevel string) {
	testCloudEvent1 := struct {
		Type string
		Data keptnutils.LogData
	}{
		Type: "sh.keptn.events.log",
		Data: keptnutils.LogData{
			Message:   msg,
			Terminate: terminate,
			LogLevel:  logLevel,
		},
	}

	data, _ := json.Marshal(testCloudEvent1)
	if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
		t.Fatalf("%v", err)
	}
}

func beginRedirectStdOut() (*os.File, *os.File, *os.File) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	return r, w, old
}

func endRedirectStdOut(r *os.File, w *os.File, old *os.File) string {
	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	return out
}
