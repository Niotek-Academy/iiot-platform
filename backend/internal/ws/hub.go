package ws

import (
	"context"
	"sync"

	
)

type Client struct {
	conn          *websocketConn
	send          chan []byte
	machineFilter string // empty = receive updates for every machine
}

type Hub struct {
	mu         sync.RWMutex
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run must be started once, in a goroutine, before any client connects.
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = true
			h.mu.Unlock()
		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
			h.mu.Unlock()
		}
	}
}

// Broadcast sends payload to every client watching machineID (or watching
// everything, if it registered with no filter). A slow/stuck client is
// skipped rather than blocking everyone else.
func (h *Hub) Broadcast(machineID string, payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for c := range h.clients {
		if c.machineFilter != "" && c.machineFilter != machineID {
			continue
		}
		select {
		case c.send <- payload:
		default:
		}
	}
}