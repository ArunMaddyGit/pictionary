package models

// Player represents a participant in a room.
type Player struct {
	ID         string
	Name       string
	Score      int
	IsDrawer   bool
	HasGuessed bool
}
