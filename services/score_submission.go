package services

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
	"gorm.io/gorm"
)

func SubmitScore(gameID, playerID, objectiveID uint) (map[string]interface{}, error) {
	game, err := helpers.GetUnfinishedGame(gameID)
	if err != nil {
		return nil, err
	}

	var objective models.Objective
	if err := database.DB.First(&objective, objectiveID).Error; err != nil {
		return nil, errors.New("objective not found")
	}

	var round models.Round
	if err := database.DB.Where("game_id = ? AND number = ?", gameID, game.CurrentRound).First(&round).Error; err != nil {
		return nil, errors.New("current round not found")
	}

	if err := ValidateSecretScoringRules(gameID, playerID, round.ID, objectiveID); err != nil {
		return nil, err
	}

	exists, err := CheckIfScoreExists(gameID, playerID, objectiveID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("objective already scored by this player")
	}

	if err := helpers.CreateObjectiveScore(gameID, round.ID, playerID, objectiveID, objective.Points); err != nil {
		return nil, fmt.Errorf("failed to add score: %v", err)
	}

	totalPoints, err := helpers.GetTotalPoints(gameID, playerID)
	if err != nil {
		return nil, err
	}

	if err := MaybeFinishGameFromScore(game, playerID); err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"message":      "Score added",
		"objective":    objective.Name,
		"points":       objective.Points,
		"round":        game.CurrentRound,
		"total_points": totalPoints,
	}
	if game.FinishedAt != nil {
		resp["message"] = "Game finished"
		resp["winner"] = game.WinnerID
	}

	return resp, nil
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
		if err := ValidateSecretScoringRules(gameID, playerID, round.ID, obj.ID); err != nil {
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

func ScoreMecatolPoint(gameID, playerID uint) error {
	roundID, err := helpers.GetCurrentRoundID(gameID)
	if err != nil {
		log.Printf("[ScoreMecatolPoint] Failed to get round ID for game %d: %v", gameID, err)
		return err
	}

	var existing models.Score
	err = database.DB.
		Where("game_id = ? AND type = ?", gameID, models.ScoreTypeMecatol).
		First(&existing).Error
	if err == nil {
		log.Printf("[ScoreMecatolPoint] Mecatol already scored for game %d", gameID)
		return fmt.Errorf("mecatol Rex point already awarded")
	}
	if err != gorm.ErrRecordNotFound {
		log.Printf("[ScoreMecatolPoint] DB error when checking existing Mecatol: %v", err)
		return err
	}

	log.Printf("[ScoreMecatolPoint] No existing Mecatol score found for game %d. Creating one for player %d", gameID, playerID)

	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		log.Printf("[ScoreMecatolPoint] Failed to load game %d: %v", gameID, err)
		return err
	}

	log.Printf("[ScoreMecatolPoint] Creating Mecatol score: Game %d, Player %d, Round %d", gameID, playerID, roundID)
	if err := helpers.CreateGenericScore(models.Score{
		GameID:   gameID,
		RoundID:  roundID,
		PlayerID: playerID,
		Points:   1,
		Type:     models.ScoreTypeMecatol,
	}); err != nil {
		log.Printf("[ScoreMecatolPoint] Failed to create Mecatol score: %v", err)
		return err
	}

	log.Printf("[ScoreMecatolPoint] Mecatol score created. Checking if game is finished.")
	return MaybeFinishGameFromScore(&game, playerID)
}

func ScoreImperialPoint(gameID, playerID uint) error {
	roundID, err := helpers.GetCurrentRoundID(gameID)
	if err != nil {
		log.Printf("[ScoreImperialPoint] Failed to get round ID for game %d: %v", gameID, err)
		return err
	}

	game, err := helpers.GetUnfinishedGame(gameID)
	if err != nil {
		log.Printf("[ScoreImperialPoint] Could not get unfinished game %d: %v", gameID, err)
		return err
	}

	log.Printf("[ScoreImperialPoint] Creating Imperial score: Game %d, Player %d, Round %d", gameID, playerID, roundID)
	if err := helpers.CreateGenericScore(models.Score{
		GameID:   gameID,
		RoundID:  roundID,
		PlayerID: playerID,
		Points:   1,
		Type:     "imperial",
	}); err != nil {
		log.Printf("[ScoreImperialPoint] Failed to create Imperial score: %v", err)
		return err
	}

	log.Printf("[ScoreImperialPoint] Imperial score created. Checking if game is finished.")
	return MaybeFinishGameFromScore(game, playerID)
}

func ScoreSupportPoint(gameID, playerID uint) error {
	roundID, err := helpers.GetCurrentRoundID(gameID)
	if err != nil {
		return err
	}

	if _, err := helpers.GetUnfinishedGame(gameID); err != nil {
		return err
	}

	var playerCount int64
	if err := database.DB.
		Model(&models.GamePlayer{}).
		Where("game_id = ?", gameID).
		Count(&playerCount).Error; err != nil {
		return err
	}

	// Get total number of scored Support points (sum of all positive and negative values)
	var totalSupportPoints int64
	if err := database.DB.
		Model(&models.Score{}).
		Where("game_id = ? AND type = ?", gameID, "Support").
		Select("COALESCE(SUM(points), 0)"). // this handles negative reversals
		Scan(&totalSupportPoints).Error; err != nil {
		return err
	}

	if totalSupportPoints >= playerCount-1 {
		return fmt.Errorf("support for the Throne can only be scored a total of %d times in a %d-player game", playerCount-1, playerCount)
	}

	return helpers.CreateGenericScore(models.Score{
		GameID:   gameID,
		RoundID:  roundID,
		PlayerID: playerID,
		Points:   1,
		Type:     models.SFTT,
	})
}

func LoseOneSupportPoint(gameID, playerID uint) error {
	roundID, err := helpers.GetCurrentRoundID(gameID)
	if err != nil {
		return err
	}

	if _, err := helpers.GetUnfinishedGame(gameID); err != nil {
		return err
	}

	// Create a negative support score record
	helpers.CreateGenericScore(models.Score{
		GameID:   gameID,
		RoundID:  roundID,
		PlayerID: playerID,
		Points:   -1,
		Type:     models.SFTT,
	})

	return nil
}

func HandleSupportForTheThrone(gameID, playerID uint, action string) error {
	switch action {
	case "score":
		return ScoreSupportPoint(gameID, playerID)
	case "unscore":
		return LoseOneSupportPoint(gameID, playerID)
	default:
		return errors.New("invalid action: must be 'score' or 'unscore'")
	}
}

func ScoreImperialRiderPoint(gameID, roundID, playerID uint) error {
	if roundID == 0 {
		var err error
		roundID, err = helpers.GetCurrentRoundID(gameID)
		if err != nil {
			return err
		}
	}

	game, err := helpers.GetUnfinishedGame(gameID)
	if err != nil {
		return err
	}

	if err := helpers.CreateGenericScore(models.Score{
		GameID:   gameID,
		RoundID:  roundID,
		PlayerID: playerID,
		Points:   1,
		Type:     "imperial_rider", // Distinct type
	}); err != nil {
		return err
	}

	return MaybeFinishGameFromScore(game, playerID)
}
