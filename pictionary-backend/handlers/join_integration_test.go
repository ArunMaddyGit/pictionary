package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pictionary/game"
	"pictionary/store"
	"pictionary/ws"
)

func TestJoinAndStartGame(t *testing.T) {
	st := store.NewMemoryStore()
	hub := ws.NewHub()
	engine := &game.GameEngine{Store: st, Hub: hub}
	router := &ws.MessageRouter{Engine: engine, Store: st}
	hub.Router = router
	hub.Engine = engine
	go hub.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/join", HandleJoin(st, engine))
	mux.HandleFunc("/ws", ws.HandleWebSocket(hub, st, engine))
	srv := httptest.NewServer(mux)
	defer srv.Close()

	join := func(name string) joinResponse {
		t.Helper()
		body := []byte(`{"name":"` + name + `"}`)
		req, err := http.NewRequest(http.MethodPost, srv.URL+"/api/join", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("NewRequest: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("join request failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want 200", resp.StatusCode)
		}
		var out joinResponse
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			t.Fatalf("decode join response: %v", err)
		}
		return out
	}

	first := join("Alice")
	second := join("Bob")

	if first.RoomID == "" || second.RoomID == "" {
		t.Fatal("room IDs should be non-empty")
	}
	if first.PlayerID == "" || second.PlayerID == "" {
		t.Fatal("player IDs should be non-empty")
	}
	if first.RoomID != second.RoomID {
		t.Fatalf("expected same roomId, got %q and %q", first.RoomID, second.RoomID)
	}
	if !strings.HasPrefix(first.WsURL, "ws://localhost:8080/ws?playerId=") {
		t.Fatalf("unexpected first wsUrl: %q", first.WsURL)
	}
	if !strings.Contains(first.WsURL, "&roomId="+first.RoomID) {
		t.Fatalf("first wsUrl missing roomId: %q", first.WsURL)
	}
	if !strings.HasPrefix(second.WsURL, "ws://localhost:8080/ws?playerId=") {
		t.Fatalf("unexpected second wsUrl: %q", second.WsURL)
	}
	if !strings.Contains(second.WsURL, "&roomId="+second.RoomID) {
		t.Fatalf("second wsUrl missing roomId: %q", second.WsURL)
	}
}
