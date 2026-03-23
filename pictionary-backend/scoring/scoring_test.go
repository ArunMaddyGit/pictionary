package scoring

import (
	"testing"

	"pictionary/models"
)

func TestCalculateScore_FirstGuesser(t *testing.T) {
	if got := CalculateScore(1); got != 120 {
		t.Fatalf("CalculateScore(1) = %d, want 120", got)
	}
}

func TestCalculateScore_LastGuesser(t *testing.T) {
	if got := CalculateScore(8); got != 10 {
		t.Fatalf("CalculateScore(8) = %d, want 10", got)
	}
}

func TestCalculateScore_OutOfRange(t *testing.T) {
	if got := CalculateScore(0); got != 0 {
		t.Fatalf("CalculateScore(0) = %d, want 0", got)
	}
	if got := CalculateScore(9); got != 0 {
		t.Fatalf("CalculateScore(9) = %d, want 0", got)
	}
}

func TestApplyScore(t *testing.T) {
	p := &models.Player{ID: "p1", Name: "A", Score: 0, HasGuessed: false}
	ApplyScore(p, 1)
	if p.Score != 120 {
		t.Fatalf("player.Score = %d, want 120", p.Score)
	}
	if !p.HasGuessed {
		t.Fatal("player.HasGuessed = false, want true")
	}
}

func TestCountCorrectGuessers(t *testing.T) {
	players := []*models.Player{
		{ID: "p1", Name: "A", HasGuessed: true, IsDrawer: false},
		{ID: "p2", Name: "B", HasGuessed: true, IsDrawer: false},
		{ID: "p3", Name: "C", HasGuessed: false, IsDrawer: true},
	}
	if got := CountCorrectGuessers(players); got != 2 {
		t.Fatalf("CountCorrectGuessers = %d, want 2", got)
	}
}

func TestBuildLeaderboard_Sorted(t *testing.T) {
	players := []*models.Player{
		{ID: "p1", Name: "Alice", Score: 50},
		{ID: "p2", Name: "Bob", Score: 120},
		{ID: "p3", Name: "Cara", Score: 80},
	}
	got := BuildLeaderboard(players)
	if len(got) != 3 {
		t.Fatalf("len(leaderboard) = %d, want 3", len(got))
	}
	if got[0].ID != "p2" || got[1].ID != "p3" || got[2].ID != "p1" {
		t.Fatalf("unexpected order: %#v", got)
	}
}
