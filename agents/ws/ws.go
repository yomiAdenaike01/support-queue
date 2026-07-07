package ws

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func newClient(id string, hub *Hub, parent context.Context, conn *websocket.Conn) *Client {
	ctx, cancel := context.WithCancel(parent)
	client := &Client{
		in:   make(chan Message),
		out:  make(chan Message),
		conn: conn,
		hub:  hub,
		cancel: sync.OnceFunc(func() {
			cancel()
			hub.unregister <- id
		}),
		ctx: ctx,
		id:  id,
	}

	go client.write()
	go client.read()
	return client
}

func New() *Hub {
	ctx, cancel := context.WithCancel(context.Background())
	hub := &Hub{
		pending: map[string]chan<- any{},
		ctx:     ctx,
		cancel:  cancel,
		clients: map[string]*Client{},
	}
	go hub.run()
	return hub
}
