package ws

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	id     string
	in     chan Message
	out    chan Message
	hub    *Hub
	conn   *websocket.Conn
	cancel func()
	ctx    context.Context
}

type Hub struct {
	in         chan Message
	pending    map[string]chan<- any
	cancel     context.CancelFunc
	ctx        context.Context
	clients    map[string]*Client
	register   chan *Client
	unregister chan string
	broadcast  chan Message
}

func (h *Hub) RegisterClient(id string, conn *websocket.Conn) {
	if _, ok := h.clients[id]; ok {
		return
	}
	h.clients[id] = newClient(id, h, h.ctx, conn)
}

func (c *Client) read() {
	c.conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	c.conn.SetReadLimit(1024)
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(5 * time.Second)); return nil })
	defer c.cancel()
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				panic(err)
			}
			var wsMessage Message
			if err := json.Unmarshal(msg, &wsMessage); err != nil {
				panic(err)
			}
			c.hub.broadcast <- wsMessage
		}

	}
}
func (c *Client) write() {

	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		ticker.Stop()
		c.cancel()
	}()
	for {
		select {
		case <-c.ctx.Done():
			return
		case message, ok := <-c.in:
			if !ok {
				panic("channel closed")
			}
			c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			rawMessage := message.toJSON()
			if rawMessage == nil {
				panic("Raw message failed")
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, rawMessage); err != nil {
				panic(err)
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				panic(err)
			}
		}
	}
}

type Message struct {
	Dest             string          `json:"dest"`
	ResolvePendingId *string         `json:"pending_id"`
	Content          json.RawMessage `json:"content"`
}

func (m Message) toJSON() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return data

}

func (h *Hub) RegisterPending(clientId string) {
	h.pending[clientId] = make(chan<- any)
}

func (h *Hub) removeClient(id string) {
	delete(h.clients, id)
}

func (h *Hub) run() {
	for {
		select {
		case message, ok := <-h.in:
			if !ok {
				return
			}
			if message.ResolvePendingId != nil {
				h.pending[*message.ResolvePendingId] <- message.Content
				close(h.pending[*message.ResolvePendingId])
				delete(h.pending, *message.ResolvePendingId)
			}

		case client, ok := <-h.register:
			if !ok {
				return
			}
			h.RegisterClient(client.id, client.conn)
		case clientId, ok := <-h.unregister:
			if !ok {
				return
			}
			h.removeClient(clientId)
		case message, ok := <-h.broadcast:
			if !ok {
				return
			}
			for _, client := range h.clients {
				timeout, cancel := context.WithTimeout(h.ctx, time.Duration(3*time.Second))
				defer cancel()
				select {
				case client.in <- message:
				case <-timeout.Done():
					client.cancel()
				}
			}

		}
	}
}
