package matchmaking

import (
	"fmt"
	"testing"

	"pictionary/models"
	"pictionary/store"
)

func TestFindOrCreateRoom_NewRoom(t *testing.T) {
	s := store.NewMemoryStore()
	room, player, err := FindOrCreateRoom(s, "Alice")
	if err != nil {
		t.Fatalf("FindOrCreateRoom: %v", err)
	}
	if room == nil || player == nil {
		t.Fatal("expected non-nil room and player")
	}
	if len(room.Players) != 1 {
		t.Fatalf("len(Players) = %d, want 1", len(room.Players))
	}
	if room.Players[0].ID != player.ID {
		t.Fatal("returned player should be the one in the room")
	}
	if player.Name != "Alice" {
		t.Errorf("Name = %q, want Alice", player.Name)
	}
	if room.ID == "" || player.ID == "" {
		t.Error("room and player must have UUID ids")
	}
}

func TestFindOrCreateRoom_JoinExisting(t *testing.T) {
	s := store.NewMemoryStore()
	r1, p1, err := FindOrCreateRoom(s, "Alice")
	if err != nil {
		t.Fatalf("first FindOrCreateRoom: %v", err)
	}
	r2, p2, err := FindOrCreateRoom(s, "Bob")
	if err != nil {
		t.Fatalf("second FindOrCreateRoom: %v", err)
	}
	if r2.ID != r1.ID {
		t.Fatalf("expected same room, got %q and %q", r1.ID, r2.ID)
	}
	if p1.ID == p2.ID {
		t.Fatal("players must have distinct ids")
	}
	if len(r2.Players) != 2 {
		t.Fatalf("len(Players) = %d, want 2", len(r2.Players))
	}
}

func TestFindOrCreateRoom_RoomFull(t *testing.T) {
	s := store.NewMemoryStore()
	full := &models.Room{
		ID:            "full-room",
		Status:        models.StatusWaiting,
		MaxRounds:     3,
		TurnDuration:  60,
		Phase:         models.PhaseWaiting,
		Players:       make([]*models.Player, 8),
	}
	for i := range full.Players {
		full.Players[i] = &models.Player{ID: fmt.Sprintf("player-%d", i), Name: "P", Score: 0}
	}
	if err := s.CreateRoom(full); err != nil {
		t.Fatalf("CreateRoom: %v", err)
	}

	room, player, err := FindOrCreateRoom(s, "Ninth")
	if err != nil {
		t.Fatalf("FindOrCreateRoom: %v", err)
	}
	if room.ID == full.ID {
		t.Fatal("expected a new room when the waiting room is full")
	}
	if len(room.Players) != 1 || room.Players[0].ID != player.ID {
		t.Fatal("new room should have exactly the new player")
	}
}

func TestFindOrCreateRoom_EmptyName(t *testing.T) {
	s := store.NewMemoryStore()
	_, _, err := FindOrCreateRoom(s, "")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
	_, _, err = FindOrCreateRoom(s, "   ")
	if err == nil {
		t.Fatal("expected error for whitespace-only name")
	}
}
