package game

import (
	"errors"
	"testing"
	"time"

	"pictionary/models"
	"pictionary/store"
	"pictionary/ws"
)

type fakeStore struct {
	rooms map[string]*models.Room
}

var _ store.Store = (*fakeStore)(nil)

func newFakeStore() *fakeStore {
	return &fakeStore{rooms: map[string]*models.Room{}}
}

func (f *fakeStore) GetRoom(id string) (*models.Room, bool) {
	r, ok := f.rooms[id]
	return r, ok
}

func (f *fakeStore) CreateRoom(room *models.Room) error {
	if room == nil || room.ID == "" {
		return errors.New("invalid room")
	}
	f.rooms[room.ID] = room
	return nil
}

func (f *fakeStore) UpdateRoom(room *models.Room) error {
	if room == nil || room.ID == "" {
		return errors.New("invalid room")
	}
	f.rooms[room.ID] = room
	return nil
}

func (f *fakeStore) DeleteRoom(id string) error {
	delete(f.rooms, id)
	return nil
}

func (f *fakeStore) ListRooms() []*models.Room {
	out := make([]*models.Room, 0, len(f.rooms))
	for _, r := range f.rooms {
		out = append(out, r)
	}
	return out
}

func newTestHubForRoom(roomID string, players []*models.Player) *ws.Hub {
	h := ws.NewHub()
	h.Rooms[roomID] = map[string]*ws.Client{}
	for _, p := range players {
		c := &ws.Client{
			PlayerID: p.ID,
			RoomID:   roomID,
			Send:     make(chan []byte, 256),
			Hub:      h,
		}
		h.Clients[p.ID] = c
		h.Rooms[roomID][p.ID] = c
	}
	return h
}

func TestMaskWord(t *testing.T) {
	e := &GameEngine{}
	if got := e.MaskWord("apple"); got != "_ _ _ _ _" {
		t.Fatalf("MaskWord(apple) = %q", got)
	}
	if got := e.MaskWord("hot dog"); got != "_ _ _   _ _ _" {
		t.Fatalf("MaskWord(hot dog) = %q", got)
	}
}

func TestStartGame_TooFewPlayers(t *testing.T) {
	s := newFakeStore()
	room := &models.Room{
		ID:        "r1",
		Status:    models.StatusWaiting,
		MaxRounds: 3,
		Players: []*models.Player{
			{ID: "p1", Name: "A"},
		},
	}
	_ = s.CreateRoom(room)
	e := &GameEngine{Store: s, Hub: newTestHubForRoom(room.ID, room.Players)}

	if err := e.StartGame(room.ID); err == nil {
		t.Fatal("expected error for too few players")
	}
}

func TestStartGame_SetsDrawer(t *testing.T) {
	s := newFakeStore()
	room := &models.Room{
		ID:                "r1",
		Status:            models.StatusWaiting,
		MaxRounds:         3,
		TurnDuration:      60,
		CurrentDrawerIndex: 0,
		Players: []*models.Player{
			{ID: "p1", Name: "A"},
			{ID: "p2", Name: "B"},
		},
	}
	_ = s.CreateRoom(room)
	e := &GameEngine{Store: s, Hub: newTestHubForRoom(room.ID, room.Players)}

	if err := e.StartGame(room.ID); err != nil {
		t.Fatalf("StartGame: %v", err)
	}
	got, _ := s.GetRoom(room.ID)
	if got.Status != models.StatusPlaying {
		t.Fatalf("Status = %v, want PLAYING", got.Status)
	}
	if got.Round != 1 {
		t.Fatalf("Round = %d, want 1", got.Round)
	}
	if !got.Players[0].IsDrawer || got.Players[1].IsDrawer {
		t.Fatal("first player should be drawer")
	}
}

func TestEndTurn_AdvancesDrawer(t *testing.T) {
	s := newFakeStore()
	room := &models.Room{
		ID:                 "r1",
		Status:             models.StatusPlaying,
		Phase:              models.PhaseDrawing,
		CurrentWord:        "apple",
		Round:              1,
		MaxRounds:          3,
		CurrentDrawerIndex: 0,
		Players: []*models.Player{
			{ID: "p1", Name: "A", IsDrawer: true},
			{ID: "p2", Name: "B"},
		},
	}
	_ = s.CreateRoom(room)
	e := &GameEngine{
		Store:       s,
		Hub:         newTestHubForRoom(room.ID, room.Players),
		turnGapSleep: func(_ time.Duration) {},
	}

	if err := e.EndTurn(room.ID); err != nil {
		t.Fatalf("EndTurn: %v", err)
	}
	got, _ := s.GetRoom(room.ID)
	if got.CurrentDrawerIndex != 1 {
		t.Fatalf("CurrentDrawerIndex = %d, want 1", got.CurrentDrawerIndex)
	}
}

func TestEndTurn_AdvancesRound(t *testing.T) {
	s := newFakeStore()
	room := &models.Room{
		ID:                 "r2",
		Status:             models.StatusPlaying,
		Phase:              models.PhaseDrawing,
		CurrentWord:        "apple",
		Round:              1,
		MaxRounds:          3,
		CurrentDrawerIndex: 1,
		Players: []*models.Player{
			{ID: "p1", Name: "A"},
			{ID: "p2", Name: "B", IsDrawer: true},
		},
	}
	_ = s.CreateRoom(room)
	e := &GameEngine{
		Store:       s,
		Hub:         newTestHubForRoom(room.ID, room.Players),
		turnGapSleep: func(_ time.Duration) {},
	}

	if err := e.EndTurn(room.ID); err != nil {
		t.Fatalf("EndTurn: %v", err)
	}
	got, _ := s.GetRoom(room.ID)
	if got.Round != 2 {
		t.Fatalf("Round = %d, want 2", got.Round)
	}
	if got.CurrentDrawerIndex != 0 {
		t.Fatalf("CurrentDrawerIndex = %d, want 0", got.CurrentDrawerIndex)
	}
}

func TestEndGame_SetsStatusFinished(t *testing.T) {
	s := newFakeStore()
	room := &models.Room{
		ID:     "r3",
		Status: models.StatusPlaying,
		Players: []*models.Player{
			{ID: "p1", Name: "A", Score: 10},
			{ID: "p2", Name: "B", Score: 20},
		},
	}
	_ = s.CreateRoom(room)
	e := &GameEngine{Store: s, Hub: newTestHubForRoom(room.ID, room.Players)}

	if err := e.EndGame(room.ID); err != nil {
		t.Fatalf("EndGame: %v", err)
	}
	got, _ := s.GetRoom(room.ID)
	if got.Status != models.StatusFinished {
		t.Fatalf("Status = %v, want FINISHED", got.Status)
	}
}
