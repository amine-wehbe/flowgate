package main

import "github.com/gorilla/websocket"

// Channels enforce that only run() touches clients — prevents concurrent map writes
type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}

// Channels must be initialized with make — nil channels block forever
func newHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// Single goroutine owns clients map — select serializes all mutations to avoid data races
func (h *Hub) run() {
	for {
		select {
		case conn := <-h.register:
			h.clients[conn] = true
		case conn := <-h.unregister:
			delete(h.clients, conn)
		case msg := <-h.broadcast:
			for conn := range h.clients {
				conn.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}
}
