package ws

import (
	"encoding/json"
	"errors"
	"strings"
)

// Message type constants (WebSocket envelope).
const (
	TypeDraw          = "DRAW"
	TypeGuess         = "GUESS"
	TypeSelectWord    = "SELECT_WORD"
	TypeClearCanvas   = "CLEAR_CANVAS"
	TypeRoomState     = "ROOM_STATE"
	TypeDrawBroadcast = "DRAW_BROADCAST"
	TypeTurnStart     = "TURN_START"
	TypeWordOptions   = "WORD_OPTIONS"
	TypeCorrectGuess  = "CORRECT_GUESS"
	TypeRoundEnd      = "ROUND_END"
	TypeGameEnd       = "GAME_END"
	TypeGuessMessage  = "GUESS_MESSAGE"
)

// Message is the JSON envelope for all WebSocket frames.
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// --- Client → Server payloads ---

// DrawPayload carries stroke data from the drawer.
type DrawPayload struct {
	Points     [][2]float64 `json:"points"`
	Color      string        `json:"color"`
	BrushSize  int           `json:"brushSize"`
}

// GuessPayload is a guess attempt from a player.
type GuessPayload struct {
	Text string `json:"text"`
}

// SelectWordPayload is the chosen word before drawing.
type SelectWordPayload struct {
	Word string `json:"word"`
}

// ClearCanvasPayload requests clearing the shared canvas.
type ClearCanvasPayload struct{}

// --- Server → Client payloads ---

// RoomStatePayload is a snapshot of room/game state for clients.
type RoomStatePayload struct {
	Players  []PlayerInfo `json:"players"`
	Round    int          `json:"round"`
	DrawerID string       `json:"drawerId"`
	Timer    int          `json:"timer"`
}

// DrawBroadcastPayload relays draw data to non-drawers.
type DrawBroadcastPayload struct {
	Points    [][2]float64 `json:"points"`
	Color     string        `json:"color"`
	BrushSize int           `json:"brushSize"`
}

// TurnStartPayload announces a new drawing turn.
type TurnStartPayload struct {
	DrawerID    string `json:"drawerId"`
	MaskedWord  string `json:"maskedWord"`
	IsDrawer    bool   `json:"isDrawer"`
}

// WordOptionsPayload lists words the drawer may pick.
type WordOptionsPayload struct {
	Words []string `json:"words"`
}

// CorrectGuessPayload notifies that someone guessed correctly.
type CorrectGuessPayload struct {
	PlayerID string `json:"playerId"`
	Score    int    `json:"score"`
}

// RoundEndPayload reveals the word at round end.
type RoundEndPayload struct {
	Word string `json:"word"`
}

// GameEndPayload carries final standings.
type GameEndPayload struct {
	Leaderboard []PlayerScore `json:"leaderboard"`
}

// GuessMessagePayload is a chat-style guess line for spectators.
type GuessMessagePayload struct {
	PlayerID   string `json:"playerId"`
	PlayerName string `json:"playerName"`
	Text       string `json:"text"`
}

// --- Helpers ---

// PlayerInfo is a player summary for ROOM_STATE.
type PlayerInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Score      int    `json:"score"`
	IsDrawer   bool   `json:"isDrawer"`
	HasGuessed bool   `json:"hasGuessed"`
}

// PlayerScore is a row on the end-game leaderboard.
type PlayerScore struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

// NewMessage marshals the envelope { type, payload } to JSON.
func NewMessage(msgType string, payload interface{}) ([]byte, error) {
	if payload == nil {
		return nil, errors.New("payload is nil")
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	env := Message{
		Type:    msgType,
		Payload: raw,
	}
	return json.Marshal(env)
}

// ParseMessage unmarshals a JSON envelope. Returns an error on invalid JSON.
func ParseMessage(data []byte) (*Message, error) {
	var m Message
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// ParseDrawPayload unmarshals a draw payload from raw JSON.
func ParseDrawPayload(raw json.RawMessage) (*DrawPayload, error) {
	var p DrawPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// ParseGuessPayload unmarshals a guess payload and requires non-empty text.
func ParseGuessPayload(raw json.RawMessage) (*GuessPayload, error) {
	var p GuessPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, err
	}
	if strings.TrimSpace(p.Text) == "" {
		return nil, errors.New("guess text is empty")
	}
	return &p, nil
}

// ParseSelectWordPayload unmarshals a select-word payload from raw JSON.
func ParseSelectWordPayload(raw json.RawMessage) (*SelectWordPayload, error) {
	var p SelectWordPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, err
	}
	return &p, nil
}
