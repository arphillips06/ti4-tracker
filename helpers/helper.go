package helpers

import (
	"errors"
	"net/http"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/gin-gonic/gin"
)

const (
	ScoreTypeAgenda = "agenda"
)

func GetCurrentRoundID(gameID uint) (uint, error) {
	var game models.Game
	if err := database.DB.Select("id, current_round").First(&game, gameID).Error; err != nil {
		return 0, err
	}

	var round models.Round
	if err := database.DB.
		Where("game_id = ? AND number = ?", game.ID, game.CurrentRound).
		First(&round).Error; err != nil {
		return 0, errors.New("current round not found")
	}

	return round.ID, nil
}

func HandleRequest[T any](c *gin.Context, handler func(input T) error) {
	input, ok := BindJSON[T](c)
	if !ok {
		return
	}
	err := handler(*input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func GetTotalPoints(gameID, playerID uint) (int, error) {
	var total int
	err := database.DB.Model(&models.Score{}).
		Where("game_id = ? AND player_id = ?", gameID, playerID).
		Select("SUM(points)").Scan(&total).Error
	return total, err
}

func GetUnfinishedGame(gameID uint) (*models.Game, error) {
	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		return nil, err
	}
	if game.FinishedAt != nil {
		return nil, errors.New("game is already finished")
	}
	return &game, nil
}

func CreateGamePlayer(gameID, playerID uint, faction string) error {
	return database.DB.Create(&models.GamePlayer{
		GameID:   gameID,
		PlayerID: playerID,
		Faction:  faction,
	}).Error
}
