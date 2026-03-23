package store

import (
	"sort"
	"testing"
	"time"

	"pictionary/models"
)

func TestCreateAndGetRoom(t *testing.T) {
	s := NewMemoryStore()
	r := &models.Room{
		ID:            "room-1",
		Status:        models.StatusWaiting,
		MaxRounds:     3,
		TurnDuration:  60,
		Phase:         models.PhaseWaiting,
		TurnStartTime: time.Now().UTC(),
	}
	if err := s.CreateRoom(r); err != nil {
		t.Fatalf("CreateRoom: %v", err)
	}
	got, ok := s.GetRoom("room-1")
	if !ok {
		t.Fatal("expected room to exist")
	}
	if got.ID != "room-1" {
		t.Errorf("ID = %q, want room-1", got.ID)
	}
	if got.Status != models.StatusWaiting {
		t.Errorf("Status = %v, want WAITING", got.Status)
	}
}

func TestUpdateRoom(t *testing.T) {
	s := NewMemoryStore()
	r := &models.Room{ID: "r1", Status: models.StatusWaiting, MaxRounds: 3, TurnDuration: 60}
	if err := s.CreateRoom(r); err != nil {
		t.Fatalf("CreateRoom: %v", err)
	}
	r.Status = models.StatusPlaying
	r.Round = 1
	if err := s.UpdateRoom(r); err != nil {
		t.Fatalf("UpdateRoom: %v", err)
	}
	got, _ := s.GetRoom("r1")
	if got.Status != models.StatusPlaying {
		t.Errorf("Status = %v, want PLAYING", got.Status)
	}
	if got.Round != 1 {
		t.Errorf("Round = %d, want 1", got.Round)
	}
}

func TestDeleteRoom(t *testing.T) {
	s := NewMemoryStore()
	r := &models.Room{ID: "gone", MaxRounds: 3, TurnDuration: 60}
	if err := s.CreateRoom(r); err != nil {
		t.Fatalf("CreateRoom: %v", err)
	}
	if err := s.DeleteRoom("gone"); err != nil {
		t.Fatalf("DeleteRoom: %v", err)
	}
	if _, ok := s.GetRoom("gone"); ok {
		t.Fatal("expected room to be deleted")
	}
}

func TestListRooms(t *testing.T) {
	s := NewMemoryStore()
	for _, id := range []string{"a", "b", "c"} {
		r := &models.Room{ID: id, MaxRounds: 3, TurnDuration: 60}
		if err := s.CreateRoom(r); err != nil {
			t.Fatalf("CreateRoom %s: %v", id, err)
		}
	}
	list := s.ListRooms()
	if len(list) != 3 {
		t.Fatalf("len(ListRooms) = %d, want 3", len(list))
	}
	ids := make([]string, len(list))
	for i, room := range list {
		ids[i] = room.ID
	}
	sort.Strings(ids)
	want := []string{"a", "b", "c"}
	for i := range want {
		if ids[i] != want[i] {
			t.Errorf("ids[%d] = %q, want %q", i, ids[i], want[i])
		}
	}
}
