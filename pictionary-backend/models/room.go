package models

import "time"

// RoomStatus is the lifecycle state of a room.
type RoomStatus string

const (
	StatusWaiting  RoomStatus = "WAITING"
	StatusPlaying  RoomStatus = "PLAYING"
	StatusFinished RoomStatus = "FINISHED"
)

// GamePhase is the phase within a round.
type GamePhase string

const (
	PhaseWaiting       GamePhase = "WAITING"
	PhaseChoosingWord  GamePhase = "CHOOSING_WORD"
	PhaseDrawing       GamePhase = "DRAWING"
	PhaseReveal        GamePhase = "REVEAL"
)

// Room holds shared game state for a Pictionary session.
type Room struct {
	ID                 string
	Players            []*Player
	Status             RoomStatus
	CurrentDrawerIndex int
	Round              int
	MaxRounds          int // default 3
	CurrentWord        string
	WordHistory        []string
	TurnStartTime      time.Time
	TurnDuration       int // default 60 seconds
	Phase              GamePhase
}
