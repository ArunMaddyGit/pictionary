package store

import (
	"errors"
	"sync"

	"pictionary/models"
)

// MemoryStore is an in-memory Store backed by a map.
type MemoryStore struct {
	mu    sync.RWMutex
	rooms map[string]*models.Room
}

// NewMemoryStore returns an empty MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		rooms: make(map[string]*models.Room),
	}
}

// GetRoom returns the room with the given id.
func (m *MemoryStore) GetRoom(id string) (*models.Room, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.rooms[id]
	return r, ok
}

// CreateRoom inserts a new room; id must be unique.
func (m *MemoryStore) CreateRoom(room *models.Room) error {
	if room == nil {
		return errors.New("room is nil")
	}
	if room.ID == "" {
		return errors.New("room id is empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.rooms[room.ID]; exists {
		return errors.New("room already exists")
	}
	m.rooms[room.ID] = room
	return nil
}

// UpdateRoom replaces an existing room by id.
func (m *MemoryStore) UpdateRoom(room *models.Room) error {
	if room == nil {
		return errors.New("room is nil")
	}
	if room.ID == "" {
		return errors.New("room id is empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.rooms[room.ID]; !exists {
		return errors.New("room not found")
	}
	m.rooms[room.ID] = room
	return nil
}

// DeleteRoom removes a room by id.
func (m *MemoryStore) DeleteRoom(id string) error {
	if id == "" {
		return errors.New("room id is empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.rooms[id]; !exists {
		return errors.New("room not found")
	}
	delete(m.rooms, id)
	return nil
}

// ListRooms returns a snapshot of all rooms.
func (m *MemoryStore) ListRooms() []*models.Room {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*models.Room, 0, len(m.rooms))
	for _, r := range m.rooms {
		out = append(out, r)
	}
	return out
}
