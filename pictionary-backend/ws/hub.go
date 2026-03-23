package ws

import (
	"sync"
	"time"
)

// InboundMessage is a raw WebSocket payload from a client.
type InboundMessage struct {
	Client *Client
	Data   []byte
}

// DisconnectHandler receives callbacks when players disconnect.
type DisconnectHandler interface {
	HandlePlayerDisconnect(roomID, playerID string) error
}

// Hub tracks connected clients and routes messages.
type Hub struct {
	Clients    map[string]*Client
	Rooms      map[string]map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Inbound    chan *InboundMessage
	Router     *MessageRouter
	Engine     DisconnectHandler
	mutex      sync.RWMutex
}

// waitUntilRegistered blocks until Clients reflects the Register receive (avoids racing the Run loop).
func (h *Hub) waitUntilRegistered(c *Client) {
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		h.mutex.RLock()
		_, ok := h.Clients[c.PlayerID]
		h.mutex.RUnlock()
		if ok {
			return
		}
		time.Sleep(time.Millisecond)
	}
}

// NewHub creates a Hub with initialized maps and channels.
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Rooms:      make(map[string]map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Inbound:    make(chan *InboundMessage, 256),
	}
}

// Run processes registration, unregistration, and inbound messages.
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.mutex.Lock()
			h.Clients[c.PlayerID] = c
			if _, ok := h.Rooms[c.RoomID]; !ok {
				h.Rooms[c.RoomID] = make(map[string]*Client)
			}
			h.Rooms[c.RoomID][c.PlayerID] = c
			h.mutex.Unlock()

		case c := <-h.Unregister:
			removed := false
			h.mutex.Lock()
			if _, ok := h.Clients[c.PlayerID]; ok {
				delete(h.Clients, c.PlayerID)
				removed = true
			}
			if room, ok := h.Rooms[c.RoomID]; ok && removed {
				delete(room, c.PlayerID)
				if len(room) == 0 {
					delete(h.Rooms, c.RoomID)
				}
			}
			h.mutex.Unlock()
			if removed {
				close(c.Send)
				if h.Engine != nil {
					_ = h.Engine.HandlePlayerDisconnect(c.RoomID, c.PlayerID)
				}
			}

		case msg := <-h.Inbound:
			if h.Router != nil && msg != nil && msg.Client != nil {
				_ = h.Router.Route(msg.Client, msg.Data)
			}
		}
	}
}

// BroadcastToRoom sends a copy of message to every client in the room.
func (h *Hub) BroadcastToRoom(roomID string, message []byte) {
	h.mutex.RLock()
	room := h.Rooms[roomID]
	clients := make([]*Client, 0, len(room))
	for _, c := range room {
		clients = append(clients, c)
	}
	h.mutex.RUnlock()

	for _, c := range clients {
		msg := append([]byte(nil), message...)
		c.Send <- msg
	}
}

// SendToPlayer sends a copy of message to a single client.
func (h *Hub) SendToPlayer(playerID string, message []byte) {
	h.mutex.RLock()
	c := h.Clients[playerID]
	h.mutex.RUnlock()
	if c == nil {
		return
	}
	msg := append([]byte(nil), message...)
	c.Send <- msg
}
