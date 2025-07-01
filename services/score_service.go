package services

import (
	"fmt"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

func ScoreImperialPoint(gameID, roundID, playerID uint) error {
	score := models.Score{
		GameID:   gameID,
		RoundID:  roundID,
		PlayerID: playerID,
		Points:   1,
		Type:     "imperial",
	}
	return database.DB.Create(&score).Error
}

func ScoreMecatolPoint(gameID, roundID, playerID uint) error {
	var existing models.Score
	err := database.DB.
		Where("game_id = ? AND type = ?", gameID, "mecatol").
		First(&existing).Error
	if err == nil {
		// Mecatol point already awarded
		return fmt.Errorf("Mecatol Rex point already awarded")
	}

	score := models.Score{
		GameID:   gameID,
		RoundID:  roundID,
		PlayerID: playerID,
		Points:   1,
		Type:     "mecatol",
	}
	return database.DB.Create(&score).Error
}

func GetObjectiveScoreSummary(gameID uint) ([]models.ObjectiveScoreSummary, error) {
	var objectives []models.Objective
	var summaries []models.ObjectiveScoreSummary

	// Get all objectives scored in this game
	err := database.DB.
		Raw(`
            SELECT DISTINCT o.id, o.name, o.stage
            FROM scores s
            JOIN objectives o ON o.id = s.objective_id
            WHERE s.game_id = ?
        `, gameID).Scan(&objectives).Error
	if err != nil {
		return nil, err
	}

	for _, obj := range objectives {
		var playerNames []string

		err := database.DB.
			Table("scores").
			Select("players.name").
			Joins("JOIN players ON players.id = scores.player_id").
			Where("scores.game_id = ? AND scores.objective_id = ?", gameID, obj.ID).
			Pluck("players.name", &playerNames).Error

		if err != nil {
			return nil, err
		}

		summaries = append(summaries, models.ObjectiveScoreSummary{
			ObjectiveID: obj.ID,
			Name:        obj.Name,
			Stage:       obj.Stage,
			ScoredBy:    playerNames,
		})
	}

	return summaries, nil
}
