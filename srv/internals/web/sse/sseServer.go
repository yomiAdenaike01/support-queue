package sse

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"maps"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	ctx context.Context
	hub *Hub
}

type Client struct {
	id   string
	in   chan json.RawMessage
	conn *gin.Context
}

type Hub struct {
	mu               sync.RWMutex
	ctx              context.Context
	in               chan json.RawMessage
	addClientChan    chan *Client
	removeClientChan chan string
	clients          map[string]*Client
}

type sendClientError struct {
	error
	clientId string
}

func (h *Hub) send(payload json.RawMessage) error {
	clients := h.getClients()
	if clients == nil {
		return nil
	}
	var wg sync.WaitGroup
	errorList := make([]error, 0, len(clients))
	timeout, cancel := context.WithTimeout(h.ctx, 5*time.Second)
	defer cancel()
	i := 0
	for _, client := range clients {
		wg.Add(1)
		go func(c *Client, i int) {
			defer wg.Done()
			select {
			case c.in <- payload:
			case <-timeout.Done():
				errMsg := "Failed to send to client timeout"
				log.Println(errMsg)
				errorList[i] = errors.New(errMsg)
			}
		}(client, i)
		i++
	}
	wg.Wait()
	return nil

}

func (c *Client) send(msg json.RawMessage, attempt ...int) {
	currentAttempt := 0
	if len(attempt) > 0 {
		currentAttempt = attempt[0]
	}
	maxAttempts := 3
	timeoutctx, stop := context.WithTimeout(c.conn.Request.Context(), time.Duration(5*currentAttempt)*time.Second)
	defer stop()
	select {
	case <-timeoutctx.Done():
		if currentAttempt+1 == maxAttempts {
			c.conn.Status(http.StatusOK)
			return
		}
		c.send(msg, currentAttempt+1)
	}
}

func (c *Client) start(removeClient chan<- string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	c.conn.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-c.in:
			if !ok {
				return false
			}
			c.conn.SSEvent("message", msg)
		case <-ticker.C:
			c.conn.SSEvent("", "ping")
			return true
		case <-c.conn.Done():
			return false
		}
		return true
	})

	go func() {
		removeTimeout, cancel := context.WithTimeout(c.conn.Request.Context(), 5*time.Second)
		defer cancel()
		<-c.conn.Done()
		select {
		case removeClient <- c.id:
		case <-removeTimeout.Done():
			return
		default:
			return
		}
	}()

}

func (h *Hub) registerClient(c *Client) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if _, ok := h.clients[c.id]; !ok {
		h.mu.Lock()
		defer h.mu.Unlock()
		h.clients[c.id] = c
		go c.start(h.removeClientChan)
		return
	}
}

func (h *Hub) getClients() map[string]*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	clients := make(map[string]*Client, len(h.clients))
	maps.Copy(clients, h.clients)
	return clients
}

func (h *Hub) shutdown() {
	var wg sync.WaitGroup
	c, stop := context.WithTimeout(h.ctx, 30*time.Second)
	defer stop()
	clients := h.getClients()
	go func() {
		<-c.Done()
	}()
	go func() {
		wg.Wait()
	}()
	go func() {
		for _, client := range clients {
			wg.Add(1)
			defer wg.Done()
			h.removeClientChan <- client.id

		}
	}()

}

func (h *Hub) getClientById(id string) *Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if client, ok := h.clients[id]; ok {
		return client
	}
	return nil
}

func (h *Hub) unregisterClient(clientId string) {
	if client := h.getClientById(clientId); client != nil {
		h.mu.Lock()
		defer h.mu.Unlock()
		delete(h.clients, clientId)
	}
}

func (h *Hub) run() {
	for {
		select {
		case clientId, ok := <-h.removeClientChan:
			if !ok {
				return
			}
			h.unregisterClient(clientId)

		case client, ok := <-h.addClientChan:
			if !ok {
				return
			}
			h.registerClient(client)
		case <-h.ctx.Done():
			h.shutdown()
		case msg, ok := <-h.in:
			if !ok {
				return
			}
			h.send(msg)

		}
	}
}

func newHub(ctx context.Context) *Hub {
	hub := &Hub{
		in:  make(chan json.RawMessage, 100),
		ctx: ctx,
	}
	go hub.run()
	return hub

}

func newClient(conn *gin.Context) *Client {
	rawClientId, ok := conn.Get("client_id")
	if !ok {
		panic("Failed to find client_id")
	}
	clientId, ok := rawClientId.(string)
	if !ok {
		panic("Failed to cast to string")
	}
	return &Client{
		id:   clientId,
		in:   make(chan json.RawMessage, 100),
		conn: conn,
	}
}

func (s *Server) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		if c.Query("id") == "" {
			c.Status(http.StatusBadRequest)
		}
		c.Set("client_id", c.Query("id"))
		s.hub.registerClient(newClient(c))
		c.Next()
	}
}

func (s *Server) Emit(payload json.RawMessage) {
	if err := s.hub.send(payload); err != nil {
		log.Printf("Failed to send message to client err=%s", err.Error())
		return
	}
	return

}

func New(ctx context.Context) *Server {
	return &Server{
		ctx: ctx,
		hub: newHub(ctx),
	}
}
