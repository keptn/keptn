package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gbrlsnchs/jwt/v2"
	"github.com/gorilla/websocket"
	keptnutils "github.com/keptn/go-utils/pkg/lib"
)

const wsLogging = false

const (
	// Time allowed to write a message to the peer.
	writeWait = 3 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 3 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type clientType struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn
}

type channelIDType string

type cliClientType struct {
	channelID channelIDType

	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

type receivedData struct {
	Shkeptncontext string `json:"shkeptncontext"`
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *clientType) readPump(l *keptnutils.Logger) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				l.Error(fmt.Sprintf("Received error while reading: %s", err.Error()))
			}
			break
		}
		if wsLogging {
			l.Debug(fmt.Sprintf("Received message from service: %s", message))
		}
		var data receivedData
		err = json.Unmarshal(message, &data)
		if err != nil {
			l.Error(fmt.Sprintf("Unmarshaling error in websocket communication: %s", err.Error()))
		}
		bData := broadcastData{channelIDType(data.Shkeptncontext), message}
		c.hub.broadcast <- &bData
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *cliClientType) writePump(l *keptnutils.Logger) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.hub.unregisterCLI <- c
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if wsLogging {
				l.Debug(fmt.Sprintf("Received message to CLI: %s", message))
			}
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				l.Debug("Cannot write message because hub closed the channel")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				return
			}

		}
	}
}

// ServeWs handles websocket requests from the services.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) error {
	l := keptnutils.NewLogger("", "", "api")
	if wsLogging {
		l.Debug("Serve internal service")
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	client := &clientType{hub: hub, conn: conn}
	client.hub.register <- client

	go client.readPump(l)
	return nil
}

// ServeWsCLI handles websocket requests from the CLI.
func ServeWsCLI(hub *Hub, w http.ResponseWriter, r *http.Request, channelID string) error {
	l := keptnutils.NewLogger("", "", "api")
	if wsLogging {
		l.Debug("Serve CLI")
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	client := &cliClientType{hub: hub, conn: conn, send: make(chan []byte, 256), channelID: channelIDType(channelID)}
	client.hub.registerCLI <- client

	go client.writePump(l)
	return nil
}

// VerifyToken verifies the Token containted in the HTTP Header
func VerifyToken(header http.Header) error {

	if val, ok := header["Token"]; ok && len(val) == 1 {
		token := val[0]

		// Define a signer.
		hs256 := jwt.NewHS256(os.Getenv("SECRET_TOKEN"))

		payload, sig, err := jwt.Parse(token)
		if err != nil {
			return err
		}
		if err = hs256.Verify(payload, sig); err != nil {
			return err
		}
		return nil
	}

	return errors.New("No Token in Header")
}

// CreateChannelInfo creates a new channel info for websockets
func CreateChannelInfo(keptnContext string) (string, error) {

	hs256 := jwt.NewHS256(os.Getenv("SECRET_TOKEN"))
	jot := &jwt.JWT{
		ExpirationTime: time.Now().Add(24 * 30 * 12 * time.Hour).Unix(),
	}
	jot.SetAlgorithm(hs256)
	payload, err := jwt.Marshal(jot)
	if err != nil {
		return "", err
	}
	token, err := hs256.Sign(payload)
	if err != nil {
		return "", err
	}

	return string(token), nil
}
