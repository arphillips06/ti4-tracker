package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
	"gorm.io/gorm"
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

func ValidateSecretScoringRules(playerID, roundID uint, phase string) error {
	var count int64
	err := database.DB.
		Model(&models.Score{}).
		Joins("JOIN objectives ON objectives.id = scores.objective_id").
		Where("scores.player_id = ? AND scores.round_id = ? AND objectives.type = ? AND objectives.phase = ?",
			playerID, roundID, "Secret", phase).
		Count(&count).Error
	if err != nil {
		return errors.New("failed to validate secret scoring rules")
	}
	if count > 0 {
		return errors.New("player has already scored a secret objective in this phase this round")
	}
	return nil
}

func AddScoreToGame(gameID, playerID uint, objectiveName string) (*models.Score, int, error) {
	var game models.Game
	if err := database.DB.Preload("Rounds").First(&game, gameID).Error; err != nil {
		return nil, 0, errors.New("game not found")
	}

	if game.FinishedAt != nil {
		return nil, 0, errors.New("game is already finished")
	}

	var obj models.Objective
	if err := database.DB.Where("LOWER(name) = ?", strings.ToLower(objectiveName)).First(&obj).Error; err != nil {
		return nil, 0, errors.New("objective not found")
	}

	var round models.Round
	if err := database.DB.Where("game_id = ? AND number = ?", game.ID, game.CurrentRound).First(&round).Error; err != nil {
		return nil, 0, errors.New("current round not found")
	}

	if obj.Type == "Secret" {
		if err := ValidateSecretScoringRules(playerID, round.ID, obj.Phase); err != nil {
			return nil, 0, err
		}
	}

	exists, err := CheckIfScoreExists(game.ID, playerID, obj.ID)
	if err != nil {
		return nil, 0, err
	}
	if exists {
		return nil, 0, errors.New("objective already scored")
	}

	score := models.Score{
		GameID:      game.ID,
		PlayerID:    playerID,
		ObjectiveID: obj.ID,
		Points:      obj.Points,
		RoundID:     round.ID,
	}

	if err := database.DB.Create(&score).Error; err != nil {
		return nil, 0, err
	}

	var total int
	database.DB.Model(&models.Score{}).
		Where("game_id = ? AND player_id = ?", game.ID, playerID).
		Select("SUM(points)").Scan(&total)

	if total >= game.WinningPoints {
		if err := MaybeFinishGameFromScore(&game, playerID); err != nil {
			return &score, total, err
		}
	}

	return &score, total, nil
}

func CheckIfScoreExists(gameID, playerID, objectiveID uint) (bool, error) {
	var existing models.Score
	err := database.DB.Where("game_id = ? AND player_id = ? AND objective_id = ?",
		gameID, playerID, objectiveID).First(&existing).Error

	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

func WinnerByScore(game *models.Game) error {
	var topScore struct {
		PlayerID uint
		Points   int
	}

	err := database.DB.
		Table("scores").
		Select("player_id, SUM(points) as points").
		Where("game_id = ?", game.ID).
		Group("player_id").
		Order("points DESC").
		Limit(1).
		Scan(&topScore).Error

	if err != nil {
		return err
	}

	if topScore.PlayerID != 0 {
		game.WinnerID = topScore.PlayerID
		return nil
	}

	return nil // No winner yet
}

func MaybeFinishGameFromExhaustion(game *models.Game) error {
	now := time.Now()
	game.FinishedAt = &now

	if err := WinnerByScore(game); err != nil {
		return err
	}

	return database.DB.Save(game).Error
}

func GetGameAndRounds(gameID uint) (*models.Game, error) {
	var game models.Game
	if err := database.DB.Preload("Rounds").First(&game, gameID).Error; err != nil {
		return nil, err
	}
	return &game, nil
}

func MaybeFinishGameFromScore(game *models.Game, scoringPlayerID uint) error {
	var totalPoints int
	err := database.DB.Model(&models.Score{}).
		Where("game_id = ? AND player_id = ?", game.ID, scoringPlayerID).
		Select("SUM(points)").Scan(&totalPoints).Error
	if err != nil {
		return err
	}

	if totalPoints >= game.WinningPoints {
		now := time.Now()
		game.FinishedAt = &now
		game.WinnerID = scoringPlayerID
		return database.DB.Save(game).Error
	}

	return nil
}
