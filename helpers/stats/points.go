package stats

import (
	"sort"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

func CalculateVictoryPointSpreads() (map[int]int, error) {
	var games []models.Game
	err := database.DB.
		Preload("GamePlayers").
		Preload("Rounds.Scores").
		Find(&games).Error
	if err != nil {
		return nil, err
	}

	spreads := map[int]int{}
	for _, game := range games {
		scoreMap := map[uint]int{}
		for _, round := range game.Rounds {
			for _, score := range round.Scores {
				scoreMap[score.PlayerID] += score.Points
			}
		}

		var scores []int
		for _, s := range scoreMap {
			scores = append(scores, s)
		}
		if len(scores) < 2 {
			continue
		}
		sort.Sort(sort.Reverse(sort.IntSlice(scores)))
		diff := scores[0] - scores[1]
		spreads[diff]++
	}

	return spreads, nil
}
