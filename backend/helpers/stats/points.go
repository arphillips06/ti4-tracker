package stats

import (
	"fmt"
	"sort"
	"strings"

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
		sort.Ints(scores)                         // ascending
		diff := scores[len(scores)-1] - scores[0] // max - min
		spreads[diff]++
	}

	return spreads, nil
}
func CalculateCommonVictoryPaths() (map[string]int, error) {
	var games []models.Game
	err := database.DB.Find(&games).Error
	if err != nil {
		return nil, err
	}

	pathCounts := make(map[string]int)

	for _, game := range games {

		var winner models.GamePlayer
		err := database.DB.
			Where("game_id = ? AND won = ?", game.ID, true).
			First(&winner).Error
		if err != nil {
			continue // skip if no winner or error
		}

		path, err := CalculateVictoryPath(game.ID, winner.PlayerID)
		if err != nil {
			continue
		}

		key := FormatVictoryPathKey(path)

		pathCounts[key] += 1
		fmt.Println("Generated Key:", key)

	}

	return pathCounts, nil
}

func CalculateVictoryPath(gameID uint, playerID uint) (models.VictoryPath, error) {
	var scores []models.Score
	err := database.DB.
		Preload("Objective").
		Where("game_id = ? AND player_id = ?", gameID, playerID).
		Find(&scores).Error
	if err != nil {
		return models.VictoryPath{}, err
	}

	vp := models.VictoryPath{}
	for _, score := range scores {
		switch strings.ToLower(score.Type) {
		case "public":
			if score.OriginallySecret {
				vp.SecretPoints += score.Points
			} else if score.Objective.ID != 0 && score.Objective.Stage != "" {
				switch score.Objective.Stage {
				case "I":
					vp.Stage1Points += score.Points
				case "II":
					vp.Stage2Scored += 1
				}
			}
		case "secret":
			vp.SecretPoints += score.Points
		case "custodians", "mecatol":
			vp.Custodians += score.Points
		case "imperial":
			vp.Imperial += score.Points
		case "relic":
			vp.Relics += score.Points
		case "agenda":
			vp.Agenda += score.Points
		case "imperial_rider":
			vp.ActionCard += score.Points
		case "support":
			vp.Support += score.Points
		}
	}
	return vp, nil
}
