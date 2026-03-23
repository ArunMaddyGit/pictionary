package ws

import (
	"encoding/json"
	"testing"
)

func TestNewMessage_Draw(t *testing.T) {
	dp := DrawPayload{
		Points:    [][2]float64{{1, 2}, {3, 4}},
		Color:     "#000000",
		BrushSize: 4,
	}
	data, err := NewMessage(TypeDraw, dp)
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}

	var env Message
	if err := json.Unmarshal(data, &env); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}
	if env.Type != TypeDraw {
		t.Fatalf("type = %q, want %q", env.Type, TypeDraw)
	}

	got, err := ParseDrawPayload(env.Payload)
	if err != nil {
		t.Fatalf("ParseDrawPayload: %v", err)
	}
	if len(got.Points) != 2 || got.Points[0][0] != 1 || got.Points[0][1] != 2 {
		t.Fatalf("points[0] = %v", got.Points)
	}
	if got.Color != "#000000" || got.BrushSize != 4 {
		t.Fatalf("got %#v", got)
	}
}

func TestNewMessage_RoomState(t *testing.T) {
	p := RoomStatePayload{
		Players: []PlayerInfo{
			{ID: "a", Name: "Alice", Score: 10, IsDrawer: true, HasGuessed: false},
		},
		Round:    2,
		DrawerID: "a",
		Timer:    45,
	}
	data, err := NewMessage(TypeRoomState, p)
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}

	var env Message
	if err := json.Unmarshal(data, &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if env.Type != TypeRoomState {
		t.Fatalf("type = %q", env.Type)
	}
	var out RoomStatePayload
	if err := json.Unmarshal(env.Payload, &out); err != nil {
		t.Fatalf("payload: %v", err)
	}
	if len(out.Players) != 1 || out.Players[0].Name != "Alice" {
		t.Fatalf("players: %#v", out.Players)
	}
	if out.Round != 2 || out.DrawerID != "a" || out.Timer != 45 {
		t.Fatalf("round/drawer/timer: %#v", out)
	}
}

func TestParseMessage_Valid(t *testing.T) {
	raw := []byte(`{"type":"GUESS","payload":{"text":"cat"}}`)
	m, err := ParseMessage(raw)
	if err != nil {
		t.Fatalf("ParseMessage: %v", err)
	}
	if m.Type != TypeGuess {
		t.Fatalf("type = %q", m.Type)
	}
	g, err := ParseGuessPayload(m.Payload)
	if err != nil {
		t.Fatalf("ParseGuessPayload: %v", err)
	}
	if g.Text != "cat" {
		t.Fatalf("text = %q", g.Text)
	}
}

func TestParseMessage_Invalid(t *testing.T) {
	_, err := ParseMessage([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
	_, err = ParseMessage([]byte(`{"type":`))
	if err == nil {
		t.Fatal("expected error for truncated JSON")
	}
}

func TestParseDrawPayload(t *testing.T) {
	raw := json.RawMessage(`{"points":[[0,0],[1,1]],"color":"#fff","brushSize":2}`)
	p, err := ParseDrawPayload(raw)
	if err != nil {
		t.Fatalf("ParseDrawPayload: %v", err)
	}
	if len(p.Points) != 2 {
		t.Fatalf("len(points) = %d", len(p.Points))
	}
	if p.Color != "#fff" || p.BrushSize != 2 {
		t.Fatalf("got %#v", p)
	}
}

func TestParseGuessPayload_Empty(t *testing.T) {
	cases := []string{
		`{"text":""}`,
		`{"text":"   "}`,
	}
	for _, s := range cases {
		raw := json.RawMessage(s)
		_, err := ParseGuessPayload(raw)
		if err == nil {
			t.Fatalf("expected error for %s", s)
		}
	}
}
