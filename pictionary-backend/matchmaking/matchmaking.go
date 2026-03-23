package matchmaking

import (
	"errors"
	"strings"
	"sync"

	"github.com/google/uuid"

	"pictionary/models"
	"pictionary/store"
)

var findOrCreateMu sync.Mutex

// FindOrCreateRoom joins a waiting room with space or creates a new room.
func FindOrCreateRoom(s store.Store, playerName string) (*models.Room, *models.Player, error) {
	if strings.TrimSpace(playerName) == "" {
		return nil, nil, errors.New("empty player name")
	}

	findOrCreateMu.Lock()
	defer findOrCreateMu.Unlock()

	rooms := s.ListRooms()
	for _, room := range rooms {
		if room.Status == models.StatusWaiting && len(room.Players) < 8 {
			player := newPlayer(playerName)
			room.Players = append(room.Players, player)
			if err := s.UpdateRoom(room); err != nil {
				return nil, nil, err
			}
			return room, player, nil
		}
	}

	roomID := uuid.NewString()
	player := newPlayer(playerName)
	room := &models.Room{
		ID:            roomID,
		Players:       []*models.Player{player},
		Status:        models.StatusWaiting,
		MaxRounds:     3,
		TurnDuration:  60,
		Phase:         models.PhaseWaiting,
	}
	if err := s.CreateRoom(room); err != nil {
		return nil, nil, err
	}
	return room, player, nil
}

func newPlayer(name string) *models.Player {
	return &models.Player{
		ID:         uuid.NewString(),
		Name:       strings.TrimSpace(name),
		Score:      0,
		IsDrawer:   false,
		HasGuessed: false,
	}
}
