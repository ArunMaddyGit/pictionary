package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
	"pictionary/models"
	"pictionary/store"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// HandleWebSocket upgrades GET /ws?playerId=...&roomId=... to WebSocket and registers the client.
func HandleWebSocket(hub *Hub, st store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		playerID := r.URL.Query().Get("playerId")
		roomID := r.URL.Query().Get("roomId")
		if playerID == "" || roomID == "" {
			http.Error(w, "playerId and roomId are required", http.StatusBadRequest)
			return
		}
		room, ok := st.GetRoom(roomID)
		if !ok || !playerInRoom(room, playerID) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		client := &Client{
			Conn:     conn,
			PlayerID: playerID,
			RoomID:   roomID,
			Send:     make(chan []byte, 256),
			Hub:      hub,
		}

		hub.Register <- client
		hub.waitUntilRegistered(client)
		go client.WritePump()
		go client.ReadPump()
	}
}

func playerInRoom(room *models.Room, playerID string) bool {
	for _, p := range room.Players {
		if p != nil && p.ID == playerID {
			return true
		}
	}
	return false
}
