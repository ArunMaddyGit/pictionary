package ws

import (
	"testing"
	"time"
)

func TestHubRegisterClient(t *testing.T) {
	h := NewHub()
	go h.Run()

	c := &Client{
		PlayerID: "player-1",
		RoomID:   "room-a",
		Send:     make(chan []byte, 256),
		Hub:      h,
	}

	h.Register <- c
	h.waitUntilRegistered(c)

	h.mutex.RLock()
	got, ok := h.Clients["player-1"]
	h.mutex.RUnlock()
	if !ok || got != c {
		t.Fatal("expected client in Clients")
	}

	h.mutex.RLock()
	room := h.Rooms["room-a"]
	h.mutex.RUnlock()
	if room == nil {
		t.Fatal("expected room map entry")
	}
	if room["player-1"] != c {
		t.Fatal("expected client in Rooms[room-a][player-1]")
	}
}

func TestHubUnregisterClient(t *testing.T) {
	h := NewHub()
	go h.Run()

	c := &Client{
		PlayerID: "player-1",
		RoomID:   "room-a",
		Send:     make(chan []byte, 256),
		Hub:      h,
	}

	h.Register <- c
	h.waitUntilRegistered(c)
	h.Unregister <- c

	deadline := time.After(500 * time.Millisecond)
	for {
		h.mutex.RLock()
		_, inClients := h.Clients["player-1"]
		room := h.Rooms["room-a"]
		h.mutex.RUnlock()
		if !inClients && (room == nil || room["player-1"] == nil) {
			break
		}
		select {
		case <-deadline:
			t.Fatal("timeout waiting for unregister")
		default:
			time.Sleep(1 * time.Millisecond)
		}
	}

	select {
	case _, ok := <-c.Send:
		if ok {
			t.Fatal("Send channel should be closed after unregister")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timeout waiting for Send to close")
	}
}

func TestBroadcastToRoom(t *testing.T) {
	h := NewHub()
	go h.Run()

	c1 := &Client{
		PlayerID: "p1",
		RoomID:   "r1",
		Send:     make(chan []byte, 256),
		Hub:      h,
	}
	c2 := &Client{
		PlayerID: "p2",
		RoomID:   "r1",
		Send:     make(chan []byte, 256),
		Hub:      h,
	}

	h.Register <- c1
	h.waitUntilRegistered(c1)
	h.Register <- c2
	h.waitUntilRegistered(c2)

	msg := []byte("hello-room")
	h.BroadcastToRoom("r1", msg)

	select {
	case got := <-c1.Send:
		if string(got) != string(msg) {
			t.Fatalf("c1 got %q, want %q", got, msg)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timeout waiting for c1")
	}

	select {
	case got := <-c2.Send:
		if string(got) != string(msg) {
			t.Fatalf("c2 got %q, want %q", got, msg)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timeout waiting for c2")
	}
}

func TestSendToPlayer(t *testing.T) {
	h := NewHub()
	go h.Run()

	c1 := &Client{
		PlayerID: "p1",
		RoomID:   "r1",
		Send:     make(chan []byte, 256),
		Hub:      h,
	}
	c2 := &Client{
		PlayerID: "p2",
		RoomID:   "r1",
		Send:     make(chan []byte, 256),
		Hub:      h,
	}

	h.Register <- c1
	h.waitUntilRegistered(c1)
	h.Register <- c2
	h.waitUntilRegistered(c2)

	msg := []byte("direct")
	h.SendToPlayer("p1", msg)

	select {
	case got := <-c1.Send:
		if string(got) != string(msg) {
			t.Fatalf("c1 got %q, want %q", got, msg)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timeout waiting for c1")
	}

	select {
	case got := <-c2.Send:
		t.Fatalf("c2 should not receive, got %q", got)
	default:
	}
}
