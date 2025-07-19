package services

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
)

func GetAllPublicObjectivesForGame(gameID string) ([]models.GameObjective, error) {
	var gameObjectives []models.GameObjective

	err := database.DB.
		Preload("Objective").
		Preload("Round").
		Where("game_id = ? AND revealed = true", gameID).
		Find(&gameObjectives).Error
	if err != nil {
		return nil, err
	}

	var scores []models.Score
	err = database.DB.
		Where("game_id = ? AND type = ? AND agenda_title = ?", gameID, "agenda", models.AgendaCDL).
		Find(&scores).Error
	if err != nil {
		return nil, err
	}

	// Convert gameID to uint
	gameIDUint, err := strconv.ParseUint(gameID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid game ID: %v", err)
	}

	gameObjectives = helpers.InjectCDLObjectives(uint(gameIDUint), gameObjectives, scores)

	sort.Slice(gameObjectives, func(i, j int) bool {
		if gameObjectives[i].Stage != gameObjectives[j].Stage {
			return gameObjectives[i].Stage < gameObjectives[j].Stage
		}
		return gameObjectives[i].Position < gameObjectives[j].Position
	})

	return gameObjectives, nil
}
