package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"pictionary/matchmaking"
	"pictionary/models"
	"pictionary/store"
)

type joinRequest struct {
	Name string `json:"name"`
}

type joinResponse struct {
	PlayerID string `json:"playerId"`
	RoomID   string `json:"roomId"`
	WsURL    string `json:"wsUrl"`
}

type joinRateEntry struct {
	count       int
	windowStart time.Time
}

var (
	joinRateMu      sync.Mutex
	joinRateByIP    = make(map[string]*joinRateEntry)
	joinRateWindow  = time.Minute
	joinRateMaxHits = 10
)

// GameStarter starts a game for a room.
type GameStarter interface {
	StartGame(roomID string) error
}

// HandleJoin handles POST /api/join for matchmaking.
func HandleJoin(st store.Store, starter GameStarter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !allowJoinRequest(remoteIP(r.RemoteAddr)) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		defer r.Body.Close()

		var req joinRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		name := strings.TrimSpace(req.Name)
		if name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}
		if utf8.RuneCountInString(name) > 20 {
			http.Error(w, "name too long", http.StatusBadRequest)
			return
		}

		room, player, err := matchmaking.FindOrCreateRoom(st, name)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		if starter != nil && len(room.Players) >= 2 && room.Status == models.StatusWaiting {
			go func(roomID string) {
				_ = starter.StartGame(roomID)
			}(room.ID)
		}

		wsURL := fmt.Sprintf("ws://localhost:8080/ws?playerId=%s&roomId=%s", player.ID, room.ID)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(joinResponse{
			PlayerID: player.ID,
			RoomID:   room.ID,
			WsURL:    wsURL,
		}); err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}
}

func remoteIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return host
}

func allowJoinRequest(ip string) bool {
	if ip == "" {
		return true
	}
	now := time.Now()
	joinRateMu.Lock()
	defer joinRateMu.Unlock()
	entry, exists := joinRateByIP[ip]
	if !exists || now.Sub(entry.windowStart) >= joinRateWindow {
		joinRateByIP[ip] = &joinRateEntry{count: 1, windowStart: now}
		return true
	}
	if entry.count >= joinRateMaxHits {
		return false
	}
	entry.count++
	return true
}
