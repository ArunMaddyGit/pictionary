package ws

import (
	"log"
	"strings"

	"pictionary/models"
	"pictionary/scoring"
	"pictionary/store"
)

// Engine handles game transitions requested by routed messages.
type Engine interface {
	HandleWordSelected(roomID string, playerID string, word string) error
	EndTurn(roomID string) error
}

// MessageRouter routes inbound WebSocket messages to game actions.
type MessageRouter struct {
	Engine Engine
	Store  store.Store
}

// Route parses and processes one inbound client frame.
func (r *MessageRouter) Route(client *Client, data []byte) error {
	msg, err := ParseMessage(data)
	if err != nil {
		return err
	}

	switch msg.Type {
	case TypeDraw:
		return r.routeDraw(client, msg.Payload)
	case TypeGuess:
		return r.routeGuess(client, msg.Payload)
	case TypeSelectWord:
		return r.routeSelectWord(client, msg.Payload)
	case TypeClearCanvas:
		return r.routeClearCanvas(client)
	default:
		log.Printf("unknown ws message type: %s", msg.Type)
		return nil
	}
}

func (r *MessageRouter) routeDraw(client *Client, raw []byte) error {
	payload, err := ParseDrawPayload(raw)
	if err != nil {
		return err
	}
	room, ok := r.Store.GetRoom(client.RoomID)
	if !ok || len(room.Players) == 0 || room.CurrentDrawerIndex < 0 || room.CurrentDrawerIndex >= len(room.Players) {
		return nil
	}
	drawer := room.Players[room.CurrentDrawerIndex]
	if drawer.ID != client.PlayerID {
		return nil
	}
	out, err := NewMessage(TypeDrawBroadcast, DrawBroadcastPayload{
		Points:    payload.Points,
		Color:     payload.Color,
		BrushSize: payload.BrushSize,
	})
	if err != nil {
		return err
	}
	client.Hub.BroadcastToRoom(client.RoomID, out)
	return nil
}

func (r *MessageRouter) routeGuess(client *Client, raw []byte) error {
	payload, err := ParseGuessPayload(raw)
	if err != nil {
		return err
	}
	room, ok := r.Store.GetRoom(client.RoomID)
	if !ok || room.Phase != models.PhaseDrawing {
		return nil
	}

	player := findPlayer(room.Players, client.PlayerID)
	if player == nil || player.HasGuessed || player.IsDrawer {
		return nil
	}

	guessText := strings.TrimSpace(payload.Text)
	normalizedGuess := strings.ToLower(guessText)
	normalizedWord := strings.ToLower(strings.TrimSpace(room.CurrentWord))

	if normalizedGuess == normalizedWord {
		order := scoring.CountCorrectGuessers(room.Players) + 1
		scoring.ApplyScore(player, order)

		correctMsg, err := NewMessage(TypeCorrectGuess, CorrectGuessPayload{
			PlayerID: player.ID,
			Score:    player.Score,
		})
		if err != nil {
			return err
		}
		client.Hub.BroadcastToRoom(room.ID, correctMsg)

		guessMsg, err := NewMessage(TypeGuessMessage, GuessMessagePayload{
			PlayerID:   player.ID,
			PlayerName: player.Name,
			Text:       "✓ " + player.Name + " guessed correctly!",
		})
		if err != nil {
			return err
		}
		client.Hub.BroadcastToRoom(room.ID, guessMsg)

		if err := r.Store.UpdateRoom(room); err != nil {
			return err
		}
		if allNonDrawerGuessed(room.Players) {
			return r.Engine.EndTurn(room.ID)
		}
		return nil
	}

	guessMsg, err := NewMessage(TypeGuessMessage, GuessMessagePayload{
		PlayerID:   player.ID,
		PlayerName: player.Name,
		Text:       guessText,
	})
	if err != nil {
		return err
	}
	client.Hub.BroadcastToRoom(room.ID, guessMsg)
	return nil
}

func (r *MessageRouter) routeSelectWord(client *Client, raw []byte) error {
	payload, err := ParseSelectWordPayload(raw)
	if err != nil {
		return err
	}
	return r.Engine.HandleWordSelected(client.RoomID, client.PlayerID, payload.Word)
}

func (r *MessageRouter) routeClearCanvas(client *Client) error {
	room, ok := r.Store.GetRoom(client.RoomID)
	if !ok || len(room.Players) == 0 || room.CurrentDrawerIndex < 0 || room.CurrentDrawerIndex >= len(room.Players) {
		return nil
	}
	drawer := room.Players[room.CurrentDrawerIndex]
	if drawer.ID != client.PlayerID {
		return nil
	}
	out, err := NewMessage(TypeClearCanvas, ClearCanvasPayload{})
	if err != nil {
		return err
	}
	client.Hub.BroadcastToRoom(room.ID, out)
	return nil
}

func findPlayer(players []*models.Player, playerID string) *models.Player {
	for _, p := range players {
		if p != nil && p.ID == playerID {
			return p
		}
	}
	return nil
}

func allNonDrawerGuessed(players []*models.Player) bool {
	nonDrawerCount := 0
	for _, p := range players {
		if p != nil && !p.IsDrawer {
			nonDrawerCount++
		}
	}
	return scoring.CountCorrectGuessers(players) >= nonDrawerCount
}
