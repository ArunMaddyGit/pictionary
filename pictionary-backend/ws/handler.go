package ws

import (
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"pictionary/models"
	"pictionary/store"
	"pictionary/words"
)

// GameStarter starts a room game lifecycle.
type GameStarter interface {
	StartGame(roomID string) error
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// HandleWebSocket upgrades GET /ws?playerId=...&roomId=... to WebSocket and registers the client.
func HandleWebSocket(hub *Hub, st store.Store, starter GameStarter) http.HandlerFunc {
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

		// Fallback: start the game when players are connected, not only on join.
		if starter != nil {
			if room, exists := st.GetRoom(roomID); exists && room.Status == models.StatusWaiting && len(room.Players) >= 2 {
				go func(id string) {
					_ = starter.StartGame(id)
				}(roomID)
			}
		}

		// Always send a fresh ROOM_STATE to the newly connected player so a late
		// websocket connection still gets the current game snapshot.
		if roomSnapshot, exists := st.GetRoom(roomID); exists {
			players := make([]PlayerInfo, 0, len(roomSnapshot.Players))
			for _, p := range roomSnapshot.Players {
				if p == nil {
					continue
				}
				players = append(players, PlayerInfo{
					ID:         p.ID,
					Name:       p.Name,
					Score:      p.Score,
					IsDrawer:   p.IsDrawer,
					HasGuessed: p.HasGuessed,
				})
			}
			drawerID := ""
			if roomSnapshot.CurrentDrawerIndex >= 0 && roomSnapshot.CurrentDrawerIndex < len(roomSnapshot.Players) {
				drawer := roomSnapshot.Players[roomSnapshot.CurrentDrawerIndex]
				if drawer != nil {
					drawerID = drawer.ID
				}
			}
			if msg, err := NewMessage(TypeRoomState, RoomStatePayload{
				Players:  players,
				Round:    roomSnapshot.Round,
				DrawerID: drawerID,
				Timer:    roomSnapshot.TurnDuration,
			}); err == nil {
				hub.SendToPlayer(playerID, msg)
			}

			// Also sync phase-specific data so a client that connected late does not
			// miss critical start-of-turn events.
			if roomSnapshot.CurrentDrawerIndex >= 0 && roomSnapshot.CurrentDrawerIndex < len(roomSnapshot.Players) {
				drawer := roomSnapshot.Players[roomSnapshot.CurrentDrawerIndex]
				if drawer != nil {
					switch roomSnapshot.Phase {
					case models.PhaseChoosingWord:
						if drawer.ID == playerID {
							if msg, err := NewMessage(TypeWordOptions, WordOptionsPayload{Words: words.GetRandomWords(3)}); err == nil {
								hub.SendToPlayer(playerID, msg)
							}
						}
					case models.PhaseDrawing:
						isDrawer := drawer.ID == playerID
						wordForPlayer := maskWord(roomSnapshot.CurrentWord)
						if isDrawer {
							wordForPlayer = roomSnapshot.CurrentWord
						}
						if msg, err := NewMessage(TypeTurnStart, TurnStartPayload{
							DrawerID:   drawer.ID,
							MaskedWord: wordForPlayer,
							IsDrawer:   isDrawer,
						}); err == nil {
							hub.SendToPlayer(playerID, msg)
						}
					case models.PhaseReveal:
						if msg, err := NewMessage(TypeRoundEnd, RoundEndPayload{Word: roomSnapshot.CurrentWord}); err == nil {
							hub.SendToPlayer(playerID, msg)
						}
					}
				}
			}
		}
		go client.WritePump()
		go client.ReadPump()
	}
}

func maskWord(word string) string {
	if strings.TrimSpace(word) == "" {
		return ""
	}
	runes := []rune(word)
	var out []string
	for _, r := range runes {
		if r == ' ' {
			out = append(out, " ")
		} else {
			out = append(out, "_")
		}
	}
	return strings.Join(out, " ")
}

func playerInRoom(room *models.Room, playerID string) bool {
	for _, p := range room.Players {
		if p != nil && p.ID == playerID {
			return true
		}
	}
	return false
}
