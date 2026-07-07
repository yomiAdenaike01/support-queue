package sse

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	id       string
	in       chan json.RawMessage
	conn     *gin.Context
	stopCh   chan struct{}
	stopOnce sync.Once
}

func newClient(conn *gin.Context, id string) *Client {
	return &Client{
		id:       id,
		in:       make(chan json.RawMessage, 100),
		conn:     conn,
		stopCh:   make(chan struct{}),
		stopOnce: sync.Once{},
	}
}

func (c *Client) stop() {
	c.stopOnce.Do(func() { close(c.stopCh) })
}

func (c *Client) start(h *Hub) {
	if c.conn == nil {
		return
	}
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	reqDone := c.conn.Request.Context().Done()

	c.conn.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-c.in:
			if !ok {
				return false
			}
			c.conn.SSEvent("message", msg)
			return true
		case <-ticker.C:
			c.conn.SSEvent("", "ping")
			return true
		case <-reqDone:
			return false
		case <-c.stopCh:
			return false
		}
	})

	h.unregisterClient(c)
}

type Hub struct {
	mu      sync.RWMutex
	ctx     context.Context
	clients map[string]*Client
}

func newHub(ctx context.Context) *Hub {
	h := &Hub{
		ctx:     ctx,
		clients: map[string]*Client{},
	}
	go func() {
		<-ctx.Done()
		h.shutdown()
	}()
	return h
}

func (h *Hub) send(payload json.RawMessage) error {
	clients := h.getClients()
	if len(clients) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error
	timeout, cancel := context.WithTimeout(h.ctx, 5*time.Second)
	defer cancel()

	for _, client := range clients {
		wg.Add(1)
		go func(c *Client) {
			defer wg.Done()
			select {
			case c.in <- payload:
			case <-timeout.Done():
				mu.Lock()
				errs = append(errs, fmt.Errorf("client %s: send timed out", c.id))
				mu.Unlock()
			}
		}(client)
	}
	wg.Wait()
	return errors.Join(errs...)
}

func (h *Hub) getClients() map[string]*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	clients := make(map[string]*Client, len(h.clients))
	maps.Copy(clients, h.clients)
	return clients
}

func (h *Hub) registerClient(c *Client) {
	h.mu.Lock()
	old, ok := h.clients[c.id]
	h.clients[c.id] = c
	h.mu.Unlock()
	if ok {
		old.stop()
	}
}

func (h *Hub) unregisterClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if existing, ok := h.clients[c.id]; ok && existing == c {
		delete(h.clients, c.id)
	}
}

func (h *Hub) shutdown() {
	for _, client := range h.getClients() {
		client.stop()
	}
}

func (s *Server) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientId := c.Query("id")
		if clientId == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")

		client := newClient(c, clientId)
		s.hub.registerClient(client)
		client.start(s.hub)
	}
}

func (s *Server) Emit(payload json.RawMessage) {
	if err := s.hub.send(payload); err != nil {
		log.Printf("Failed to send message to client err=%s", err.Error())
	}
}

func New(ctx context.Context) *Server {
	return &Server{
		ctx: ctx,
		hub: newHub(ctx),
	}
}
