package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/magiconair/properties/assert"
)

var hub *Hub

func handler(w http.ResponseWriter, r *http.Request) {
	if val, ok := r.Header["Keptn-Ws-Channel-Id"]; ok {
		ServeWsCLI(hub, w, r, val[0])
	} else {
		ServeWs(hub, w, r)
	}
}

func TestServiceRegistering(t *testing.T) {

	hub = NewHub()
	go hub.Run()

	http.DefaultServeMux = http.NewServeMux()
	srv := &http.Server{Addr: ":80"}
	http.HandleFunc("/", handler)

	go func() {
		srv.ListenAndServe()
	}()

	u := url.URL{Scheme: "ws", Host: "localhost", Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	waitCount := 0
	for len(hub.clients) < 1 && waitCount < 2 {
		time.Sleep(1 * time.Second)
		waitCount++
	}
	assert.Equal(t, len(hub.clients), 1, "Client not registered")

	// Close the client
	c.Close()
	waitCount = 0
	for len(hub.clients) > 0 && waitCount < 2 {
		time.Sleep(1 * time.Second)
		waitCount++
	}
	assert.Equal(t, len(hub.clients), 0, "Client not unregistered")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
}

func TestCLIRegistering(t *testing.T) {

	hub = NewHub()
	go hub.Run()

	http.DefaultServeMux = http.NewServeMux()
	srv := &http.Server{Addr: ":80"}
	http.HandleFunc("/", handler)

	go func() {
		srv.ListenAndServe()
	}()

	u := url.URL{Scheme: "ws", Host: "localhost", Path: "/"}
	log.Printf("connecting to %s", u.String())

	header := http.Header{}
	header.Add("Token", "adf")
	header.Add("Keptn-Ws-Channel-Id", "asdf")

	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("dial:", err)
	}

	waitCount := 0
	for len(hub.cliClients) < 1 && waitCount < 2 {
		time.Sleep(1 * time.Second)
		waitCount++
	}
	assert.Equal(t, len(hub.cliClients), 1, "CLI Client not registered")

	// Close the client
	c.Close()
	waitCount = 0
	for len(hub.cliClients) > 0 && waitCount < 20 {
		time.Sleep(1 * time.Second)
		waitCount++
	}
	assert.Equal(t, len(hub.cliClients), 0, "CLI Client not unregistered")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
}

func TestSendMessage(t *testing.T) {

	hub = NewHub()
	go hub.Run()

	http.DefaultServeMux = http.NewServeMux()
	srv := &http.Server{Addr: ":80"}
	http.HandleFunc("/", handler)

	go func() {
		srv.ListenAndServe()
	}()

	u := url.URL{Scheme: "ws", Host: "localhost", Path: "/"}
	log.Printf("connecting to %s", u.String())

	header := http.Header{}
	header.Add("Token", "adf")
	header.Add("Keptn-Ws-Channel-Id", "asdf")

	cliClient, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("dial:", err)
	}

	serviceClient, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan bool)

	message := struct {
		Shkeptncontext string `json:"shkeptncontext"`
	}{
		"asdf",
	}

	messageData, _ := json.Marshal(message)

	go func() {
		_, received, err := cliClient.ReadMessage()
		if err != nil {
			log.Fatal()
		}
		assert.Equal(t, received, messageData)
		fmt.Println("Received data match")
		done <- true
	}()

	writeMessage(serviceClient, messageData)

	// Close clients
	serviceClient.Close()
	<-done

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
}

func TestBuffering(t *testing.T) {

	hub = NewHub()
	go hub.Run()

	http.DefaultServeMux = http.NewServeMux()
	srv := &http.Server{Addr: ":80"}
	http.HandleFunc("/", handler)

	go func() {
		srv.ListenAndServe()
	}()

	u := url.URL{Scheme: "ws", Host: "localhost", Path: "/"}
	log.Printf("connecting to %s", u.String())

	header := http.Header{}
	header.Add("Token", "adf")
	header.Add("Keptn-Ws-Channel-Id", "asdf")

	serviceClient, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan bool)

	message := struct {
		Shkeptncontext string `json:"shkeptncontext"`
	}{
		"asdf",
	}

	messageData, _ := json.Marshal(message)
	writeMessage(serviceClient, messageData)

	// Close service client
	serviceClient.Close()

	cliClient, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	go func() {
		_, received, err := cliClient.ReadMessage()
		if err != nil {
			log.Fatal()
		}
		assert.Equal(t, received, messageData)
		fmt.Println("Received data match")
		done <- true
	}()

	<-done

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
}

func TestPositiveVerification(t *testing.T) {

	os.Setenv("SECRET_TOKEN", "test-token")

	keptnContext := uuid.New().String()
	token, err := CreateChannelInfo(keptnContext)

	assert.Equal(t, err, nil)

	var header http.Header
	header = make(http.Header)
	header.Add("Token", token)
	err = VerifyToken(header)
	assert.Equal(t, err, nil)
}

func TestNegativeVerification(t *testing.T) {

	os.Setenv("SECRET_TOKEN", "test-token")

	var header http.Header
	header = make(http.Header)
	header.Add("Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTQyODAxOTh9.1yZsr6r9F4Ftpj9AsN3AeE6N_Tjr2oGDjHMkdO1Z0P3")
	err := VerifyToken(header)
	assert.Equal(t, err, errors.New("jwt: HMAC verification failed"))
}

func writeMessage(client *websocket.Conn, message []byte) {
	w, err := client.NextWriter(websocket.TextMessage)
	if err != nil {
		log.Fatal()
	}
	w.Write(message)
	if err := w.Close(); err != nil {
		return
	}
}
