package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Client is a WebSocket peer registered with the Hub.
type Client struct {
	Conn     *websocket.Conn
	PlayerID string
	RoomID   string
	Send     chan []byte
	Hub      *Hub

	unregisterOnce sync.Once
}

func (c *Client) unregister() {
	c.unregisterOnce.Do(func() {
		c.Hub.Unregister <- c
	})
}

// ReadPump reads messages from the WebSocket and forwards them to the Hub.
func (c *Client) ReadPump() {
	defer func() {
		c.unregister()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		data := append([]byte(nil), message...)
		c.Hub.Inbound <- &InboundMessage{Client: c, Data: data}
	}
}

// WritePump drains Send and writes to the WebSocket until Send is closed.
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
		c.unregister()
	}()

	for {
		message, ok := <-c.Send
		if !ok {
			return
		}
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}
