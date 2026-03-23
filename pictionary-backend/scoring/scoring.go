package scoring

import (
	"sort"

	"pictionary/models"
	"pictionary/ws"
)

var scoreTable = []int{120, 100, 80, 60, 40, 30, 20, 10}

// CalculateScore returns points for a 1-based guess order.
func CalculateScore(guessOrder int) int {
	if guessOrder < 1 || guessOrder > len(scoreTable) {
		return 0
	}
	return scoreTable[guessOrder-1]
}

// CountCorrectGuessers counts non-drawer players that already guessed.
func CountCorrectGuessers(players []*models.Player) int {
	count := 0
	for _, p := range players {
		if p != nil && p.HasGuessed && !p.IsDrawer {
			count++
		}
	}
	return count
}

// ApplyScore updates player score based on guess order and marks guessed.
func ApplyScore(player *models.Player, guessOrder int) {
	if player == nil {
		return
	}
	player.Score += CalculateScore(guessOrder)
	player.HasGuessed = true
}

// BuildLeaderboard returns all players sorted by score descending.
func BuildLeaderboard(players []*models.Player) []ws.PlayerScore {
	leaderboard := make([]ws.PlayerScore, 0, len(players))
	for _, p := range players {
		if p == nil {
			continue
		}
		leaderboard = append(leaderboard, ws.PlayerScore{
			ID:    p.ID,
			Name:  p.Name,
			Score: p.Score,
		})
	}
	sort.SliceStable(leaderboard, func(i, j int) bool {
		return leaderboard[i].Score > leaderboard[j].Score
	})
	return leaderboard
}
