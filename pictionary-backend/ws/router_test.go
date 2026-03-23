package ws

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"pictionary/models"
	"pictionary/store"
)

type routerMockStore struct {
	rooms map[string]*models.Room
}

var _ store.Store = (*routerMockStore)(nil)

func newRouterMockStore() *routerMockStore {
	return &routerMockStore{rooms: map[string]*models.Room{}}
}

func (m *routerMockStore) GetRoom(id string) (*models.Room, bool) {
	r, ok := m.rooms[id]
	return r, ok
}
func (m *routerMockStore) CreateRoom(room *models.Room) error {
	if room == nil || room.ID == "" {
		return errors.New("invalid room")
	}
	m.rooms[room.ID] = room
	return nil
}
func (m *routerMockStore) UpdateRoom(room *models.Room) error {
	if room == nil || room.ID == "" {
		return errors.New("invalid room")
	}
	m.rooms[room.ID] = room
	return nil
}
func (m *routerMockStore) DeleteRoom(id string) error {
	delete(m.rooms, id)
	return nil
}
func (m *routerMockStore) ListRooms() []*models.Room {
	out := make([]*models.Room, 0, len(m.rooms))
	for _, r := range m.rooms {
		out = append(out, r)
	}
	return out
}

type routerMockEngine struct {
	endTurnCalls       int
	handleSelectCalls  int
	lastRoomID         string
	lastPlayerID       string
	lastSelectedWord   string
}

func (m *routerMockEngine) HandleWordSelected(roomID string, playerID string, word string) error {
	m.handleSelectCalls++
	m.lastRoomID = roomID
	m.lastPlayerID = playerID
	m.lastSelectedWord = word
	return nil
}

func (m *routerMockEngine) EndTurn(roomID string) error {
	m.endTurnCalls++
	m.lastRoomID = roomID
	return nil
}

func setupRouterHarness() (*MessageRouter, *Hub, *routerMockStore, *routerMockEngine, *Client, *Client, *models.Room) {
	s := newRouterMockStore()
	e := &routerMockEngine{}
	h := NewHub()

	room := &models.Room{
		ID:                 "room-1",
		Status:             models.StatusPlaying,
		Phase:              models.PhaseDrawing,
		CurrentWord:        "apple",
		CurrentDrawerIndex: 0,
		Players: []*models.Player{
			{ID: "p1", Name: "Drawer", IsDrawer: true},
			{ID: "p2", Name: "Guesser", IsDrawer: false},
		},
	}
	_ = s.CreateRoom(room)

	drawer := &Client{PlayerID: "p1", RoomID: room.ID, Send: make(chan []byte, 256), Hub: h}
	guesser := &Client{PlayerID: "p2", RoomID: room.ID, Send: make(chan []byte, 256), Hub: h}
	h.Clients[drawer.PlayerID] = drawer
	h.Clients[guesser.PlayerID] = guesser
	h.Rooms[room.ID] = map[string]*Client{
		drawer.PlayerID:  drawer,
		guesser.PlayerID: guesser,
	}

	r := &MessageRouter{Engine: e, Store: s}
	h.Router = r
	return r, h, s, e, drawer, guesser, room
}

func decodeEnvelope(t *testing.T, data []byte) *Message {
	t.Helper()
	m, err := ParseMessage(data)
	if err != nil {
		t.Fatalf("ParseMessage: %v", err)
	}
	return m
}

func drain(ch chan []byte) {
	for {
		select {
		case <-ch:
		default:
			return
		}
	}
}

func TestRoute_Draw_ValidDrawer(t *testing.T) {
	r, _, _, _, drawer, guesser, _ := setupRouterHarness()
	drain(drawer.Send)
	drain(guesser.Send)

	in, _ := NewMessage(TypeDraw, DrawPayload{
		Points:    [][2]float64{{1, 2}},
		Color:     "#000",
		BrushSize: 3,
	})
	if err := r.Route(drawer, in); err != nil {
		t.Fatalf("Route: %v", err)
	}

	select {
	case msg := <-drawer.Send:
		env := decodeEnvelope(t, msg)
		if env.Type != TypeDrawBroadcast {
			t.Fatalf("type = %s, want %s", env.Type, TypeDrawBroadcast)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("drawer did not receive broadcast")
	}
	select {
	case msg := <-guesser.Send:
		env := decodeEnvelope(t, msg)
		if env.Type != TypeDrawBroadcast {
			t.Fatalf("type = %s, want %s", env.Type, TypeDrawBroadcast)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("guesser did not receive broadcast")
	}
}

func TestRoute_Draw_InvalidDrawer(t *testing.T) {
	r, _, _, _, drawer, guesser, _ := setupRouterHarness()
	drain(drawer.Send)
	drain(guesser.Send)

	in, _ := NewMessage(TypeDraw, DrawPayload{
		Points:    [][2]float64{{1, 2}},
		Color:     "#000",
		BrushSize: 3,
	})
	if err := r.Route(guesser, in); err != nil {
		t.Fatalf("Route: %v", err)
	}

	select {
	case <-drawer.Send:
		t.Fatal("drawer should not receive draw broadcast from invalid sender")
	default:
	}
	select {
	case <-guesser.Send:
		t.Fatal("guesser should not receive draw broadcast from invalid sender")
	default:
	}
}

func TestRoute_Guess_Correct(t *testing.T) {
	r, _, s, _, drawer, guesser, room := setupRouterHarness()
	drain(drawer.Send)
	drain(guesser.Send)

	in, _ := NewMessage(TypeGuess, GuessPayload{Text: " apple "})
	if err := r.Route(guesser, in); err != nil {
		t.Fatalf("Route: %v", err)
	}

	foundCorrect := false
	for i := 0; i < 2; i++ {
		select {
		case msg := <-drawer.Send:
			env := decodeEnvelope(t, msg)
			if env.Type == TypeCorrectGuess {
				foundCorrect = true
			}
		case <-time.After(200 * time.Millisecond):
			t.Fatal("missing expected broadcast")
		}
	}
	if !foundCorrect {
		t.Fatal("expected CORRECT_GUESS broadcast")
	}
	updated, _ := s.GetRoom(room.ID)
	p := findPlayer(updated.Players, guesser.PlayerID)
	if p == nil || !p.HasGuessed || p.Score == 0 {
		t.Fatal("player should be marked guessed and scored")
	}
}

func TestRoute_Guess_AlreadyGuessed(t *testing.T) {
	r, _, _, _, drawer, guesser, room := setupRouterHarness()
	drain(drawer.Send)
	drain(guesser.Send)

	in, _ := NewMessage(TypeGuess, GuessPayload{Text: "apple"})
	if err := r.Route(guesser, in); err != nil {
		t.Fatalf("first Route: %v", err)
	}
	drain(drawer.Send)
	drain(guesser.Send)

	if err := r.Route(guesser, in); err != nil {
		t.Fatalf("second Route: %v", err)
	}
	select {
	case <-drawer.Send:
		t.Fatal("second guess should be ignored")
	default:
	}
	select {
	case <-guesser.Send:
		t.Fatal("second guess should be ignored")
	default:
	}
	_ = room
}

func TestRoute_Guess_Incorrect(t *testing.T) {
	r, _, _, _, drawer, guesser, _ := setupRouterHarness()
	drain(drawer.Send)
	drain(guesser.Send)

	in, _ := NewMessage(TypeGuess, GuessPayload{Text: "banana"})
	if err := r.Route(guesser, in); err != nil {
		t.Fatalf("Route: %v", err)
	}

	select {
	case msg := <-drawer.Send:
		env := decodeEnvelope(t, msg)
		if env.Type != TypeGuessMessage {
			t.Fatalf("type = %s, want %s", env.Type, TypeGuessMessage)
		}
		var payload GuessMessagePayload
		if err := json.Unmarshal(env.Payload, &payload); err != nil {
			t.Fatalf("payload unmarshal: %v", err)
		}
		if payload.Text != "banana" {
			t.Fatalf("text = %q, want banana", payload.Text)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("missing guess message")
	}
}

func TestRoute_ClearCanvas_ValidDrawer(t *testing.T) {
	r, _, _, _, drawer, guesser, _ := setupRouterHarness()
	drain(drawer.Send)
	drain(guesser.Send)

	in, _ := NewMessage(TypeClearCanvas, ClearCanvasPayload{})
	if err := r.Route(drawer, in); err != nil {
		t.Fatalf("Route: %v", err)
	}

	select {
	case msg := <-drawer.Send:
		env := decodeEnvelope(t, msg)
		if env.Type != TypeClearCanvas {
			t.Fatalf("type = %s, want %s", env.Type, TypeClearCanvas)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("drawer did not receive clear canvas broadcast")
	}
	select {
	case msg := <-guesser.Send:
		env := decodeEnvelope(t, msg)
		if env.Type != TypeClearCanvas {
			t.Fatalf("type = %s, want %s", env.Type, TypeClearCanvas)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("guesser did not receive clear canvas broadcast")
	}
}
