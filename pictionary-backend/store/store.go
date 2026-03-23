package store

import "pictionary/models"

// Store persists rooms; implementations may be swapped (e.g. Redis later).
type Store interface {
	GetRoom(id string) (*models.Room, bool)
	CreateRoom(room *models.Room) error
	UpdateRoom(room *models.Room) error
	DeleteRoom(id string) error
	ListRooms() []*models.Room
}
