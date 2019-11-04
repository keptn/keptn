// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import keptnutils "github.com/keptn/go-utils/pkg/utils"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	cliClients map[channelIDType]map[*cliClientType]bool

	// Registered clients.
	clients map[*clientType]bool

	// Inbound messages from the clients.
	broadcast chan *broadcastData

	// Register requests from the clients.
	register chan *clientType

	// Register requests from the clients.
	registerCLI chan *cliClientType

	// Unregister requests from clients.
	unregister chan *clientType

	unregisterCLI chan *cliClientType

	buffers map[channelIDType][][]byte
}

type broadcastData struct {
	channelID channelIDType
	data      []byte
}

func NewHub() *Hub {
	return &Hub{
		broadcast:     make(chan *broadcastData),
		register:      make(chan *clientType),
		registerCLI:   make(chan *cliClientType),
		unregister:    make(chan *clientType),
		unregisterCLI: make(chan *cliClientType),
		clients:       make(map[*clientType]bool),
		cliClients:    make(map[channelIDType]map[*cliClientType]bool),
		buffers:       make(map[channelIDType][][]byte),
	}
}

func (h *Hub) Run() {
	l := keptnutils.NewLogger("", "", "api")

	for {
		select {
		case client := <-h.register:
			if wsLogging {
				l.Debug("Registered service")
			}
			h.clients[client] = true
		case cliClient := <-h.registerCLI:
			if wsLogging {
				l.Debug("Registered CLI")
			}
			if _, available := h.cliClients[cliClient.channelID]; !available {
				h.cliClients[cliClient.channelID] = make(map[*cliClientType]bool)
			}
			h.cliClients[cliClient.channelID][cliClient] = true

			if len(h.buffers[cliClient.channelID]) > 0 {
				// Send buffered messages
				for _, data := range h.buffers[cliClient.channelID] {
					cliClient.send <- data
				}
				// Clear buffer
				delete(h.buffers, cliClient.channelID)
			}

		case client := <-h.unregister:
			if wsLogging {
				l.Debug("Unregistered service")
			}
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
			}
		case cliClient := <-h.unregisterCLI:
			if wsLogging {
				l.Debug("Unregistered CLI")
			}
			if _, ok := h.cliClients[cliClient.channelID][cliClient]; ok {
				delete(h.cliClients[cliClient.channelID], cliClient)
				close(cliClient.send)

				if len(h.cliClients[cliClient.channelID]) == 0 {
					delete(h.cliClients, cliClient.channelID)
				}
			}
		case message := <-h.broadcast:
			if wsLogging {
				l.Debug("Broadcast message")
			}
			if _, available := h.cliClients[message.channelID]; !available {
				// Buffer message
				if _, available := h.buffers[message.channelID]; !available {
					h.buffers[message.channelID] = [][]byte{}
				}
				h.buffers[message.channelID] = append(h.buffers[message.channelID], message.data)
			} else {
				cliClients := h.cliClients[message.channelID]
				for cliClient := range cliClients {
					select {
					case cliClient.send <- message.data:
					default:
						close(cliClient.send)
						delete(cliClients, cliClient)
					}
				}
			}
		}
	}
}
