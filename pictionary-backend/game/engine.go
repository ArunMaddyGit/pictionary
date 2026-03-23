package game

import (
	"errors"
	"strings"
	"time"

	"pictionary/models"
	"pictionary/scoring"
	"pictionary/store"
	"pictionary/words"
	"pictionary/ws"
)

// GameEngine coordinates turn lifecycle and game progression.
type GameEngine struct {
	Store store.Store
	Hub   *ws.Hub

	turnGapSleep func(time.Duration)
}

func (g *GameEngine) sleep(d time.Duration) {
	if g.turnGapSleep != nil {
		g.turnGapSleep(d)
		return
	}
	time.Sleep(d)
}

// StartGame validates room state and starts round 1.
func (g *GameEngine) StartGame(roomID string) error {
	room, ok := g.Store.GetRoom(roomID)
	if !ok {
		return errors.New("room not found")
	}
	if room.Status != models.StatusWaiting {
		return errors.New("room is not in waiting state")
	}
	if len(room.Players) < 2 {
		return errors.New("need at least 2 players")
	}
	room.Status = models.StatusPlaying
	room.Round = 1
	if room.MaxRounds == 0 {
		room.MaxRounds = 3
	}
	if room.TurnDuration == 0 {
		room.TurnDuration = 60
	}
	if err := g.Store.UpdateRoom(room); err != nil {
		return err
	}
	return g.StartTurn(roomID)
}

// StartTurn sets drawer/phase, emits options, and syncs room state.
func (g *GameEngine) StartTurn(roomID string) error {
	room, ok := g.Store.GetRoom(roomID)
	if !ok {
		return errors.New("room not found")
	}
	if len(room.Players) == 0 {
		return errors.New("room has no players")
	}
	if room.CurrentDrawerIndex < 0 || room.CurrentDrawerIndex >= len(room.Players) {
		room.CurrentDrawerIndex = 0
	}

	drawer := room.Players[room.CurrentDrawerIndex]
	for i, p := range room.Players {
		p.IsDrawer = i == room.CurrentDrawerIndex
		p.HasGuessed = false
	}

	wordOptions := g.pickWordsExcludingHistory(room.WordHistory, 3)
	room.Phase = models.PhaseChoosingWord
	room.TurnStartTime = time.Now()

	if payload, err := ws.NewMessage(ws.TypeWordOptions, ws.WordOptionsPayload{Words: wordOptions}); err == nil {
		g.Hub.SendToPlayer(drawer.ID, payload)
	}

	for _, p := range room.Players {
		if p.ID == drawer.ID {
			continue
		}
		if payload, err := ws.NewMessage(ws.TypeTurnStart, ws.TurnStartPayload{
			DrawerID:   drawer.ID,
			MaskedWord: g.MaskWord(""),
			IsDrawer:   false,
		}); err == nil {
			g.Hub.SendToPlayer(p.ID, payload)
		}
	}

	statePlayers := make([]ws.PlayerInfo, 0, len(room.Players))
	for _, p := range room.Players {
		statePlayers = append(statePlayers, ws.PlayerInfo{
			ID:         p.ID,
			Name:       p.Name,
			Score:      p.Score,
			IsDrawer:   p.IsDrawer,
			HasGuessed: p.HasGuessed,
		})
	}
	statePayload, err := ws.NewMessage(ws.TypeRoomState, ws.RoomStatePayload{
		Players:  statePlayers,
		Round:    room.Round,
		DrawerID: drawer.ID,
		Timer:    room.TurnDuration,
	})
	if err == nil {
		g.Hub.BroadcastToRoom(room.ID, statePayload)
	}
	return g.Store.UpdateRoom(room)
}

// HandleWordSelected validates drawer choice, enters drawing phase, and starts turn timer.
func (g *GameEngine) HandleWordSelected(roomID string, playerID string, word string) error {
	room, ok := g.Store.GetRoom(roomID)
	if !ok {
		return errors.New("room not found")
	}
	if len(room.Players) == 0 {
		return errors.New("room has no players")
	}
	if room.CurrentDrawerIndex < 0 || room.CurrentDrawerIndex >= len(room.Players) {
		return errors.New("invalid drawer index")
	}
	drawer := room.Players[room.CurrentDrawerIndex]
	if drawer.ID != playerID {
		return errors.New("only current drawer can select a word")
	}
	if strings.TrimSpace(word) == "" {
		return errors.New("word is empty")
	}

	room.CurrentWord = word
	room.WordHistory = append(room.WordHistory, word)
	room.Phase = models.PhaseDrawing
	room.TurnStartTime = time.Now()

	masked := g.MaskWord(word)
	for _, p := range room.Players {
		wordForPlayer := masked
		isDrawer := false
		if p.ID == drawer.ID {
			wordForPlayer = word
			isDrawer = true
		}
		if payload, err := ws.NewMessage(ws.TypeTurnStart, ws.TurnStartPayload{
			DrawerID:   drawer.ID,
			MaskedWord: wordForPlayer,
			IsDrawer:   isDrawer,
		}); err == nil {
			g.Hub.SendToPlayer(p.ID, payload)
		}
	}

	if err := g.Store.UpdateRoom(room); err != nil {
		return err
	}

	duration := room.TurnDuration
	if duration <= 0 {
		duration = 60
	}
	go func() {
		time.Sleep(time.Duration(duration) * time.Second)
		_ = g.EndTurn(roomID)
	}()
	return nil
}

// MaskWord replaces letters with underscores while preserving spaces.
func (g *GameEngine) MaskWord(word string) string {
	if word == "" {
		return ""
	}
	runes := []rune(word)
	var b strings.Builder
	for i, r := range runes {
		if i > 0 {
			b.WriteString(" ")
		}
		if r == ' ' {
			b.WriteString(" ")
		} else {
			b.WriteString("_")
		}
	}
	return b.String()
}

// EndTurn reveals the word, awards drawer bonus, and advances progression.
func (g *GameEngine) EndTurn(roomID string) error {
	room, ok := g.Store.GetRoom(roomID)
	if !ok {
		return errors.New("room not found")
	}
	if room.Phase != models.PhaseDrawing {
		return nil
	}
	if len(room.Players) == 0 {
		return errors.New("room has no players")
	}

	room.Phase = models.PhaseReveal
	if payload, err := ws.NewMessage(ws.TypeRoundEnd, ws.RoundEndPayload{Word: room.CurrentWord}); err == nil {
		g.Hub.BroadcastToRoom(room.ID, payload)
	}

	if scoring.CountCorrectGuessers(room.Players) > 0 &&
		room.CurrentDrawerIndex >= 0 &&
		room.CurrentDrawerIndex < len(room.Players) {
		room.Players[room.CurrentDrawerIndex].Score += 10
	}

	room.CurrentDrawerIndex++
	if room.CurrentDrawerIndex >= len(room.Players) {
		room.CurrentDrawerIndex = 0
		room.Round++
	}

	if room.Round > room.MaxRounds {
		if err := g.Store.UpdateRoom(room); err != nil {
			return err
		}
		return g.EndGame(roomID)
	}

	if err := g.Store.UpdateRoom(room); err != nil {
		return err
	}

	go func() {
		g.sleep(3 * time.Second)
		_ = g.StartTurn(roomID)
	}()
	return nil
}

// EndGame finalizes room status and broadcasts leaderboard.
func (g *GameEngine) EndGame(roomID string) error {
	room, ok := g.Store.GetRoom(roomID)
	if !ok {
		return errors.New("room not found")
	}
	room.Status = models.StatusFinished
	leaderboard := scoring.BuildLeaderboard(room.Players)
	if payload, err := ws.NewMessage(ws.TypeGameEnd, ws.GameEndPayload{Leaderboard: leaderboard}); err == nil {
		g.Hub.BroadcastToRoom(room.ID, payload)
	}
	return g.Store.UpdateRoom(room)
}

func (g *GameEngine) pickWordsExcludingHistory(history []string, count int) []string {
	out := make([]string, 0, count)
	seen := make(map[string]struct{}, count)
	for attempts := 0; attempts < 16 && len(out) < count; attempts++ {
		candidates := words.GetRandomWords(count * 2)
		for _, w := range candidates {
			if words.IsDuplicateInHistory(w, history) {
				continue
			}
			if _, ok := seen[w]; ok {
				continue
			}
			seen[w] = struct{}{}
			out = append(out, w)
			if len(out) == count {
				return out
			}
		}
	}
	if len(out) == 0 {
		return words.GetRandomWords(count)
	}
	return out
}
