package services

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

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

func ValidateSecretScoringRules(gameID, playerID, roundID, objectiveID uint) error {
	var objective models.Objective
	if err := database.DB.First(&objective, objectiveID).Error; err != nil {
		log.Printf("[ERROR] Could not find objective %d: %v", objectiveID, err)
		return errors.New("Objective not found")
	}

	if strings.ToLower(objective.Type) != "secret" {
		return nil
	}

	// Phase-specific limit check
	var countThisPhase int64
	if err := database.DB.
		Model(&models.Score{}).
		Where(`
			player_id = ? AND 
			round_id = ? AND 
			LOWER(type) = 'secret' AND 
			objective_id IN (
				SELECT id FROM objectives WHERE LOWER(phase) = ?
			)`,
			playerID, roundID, strings.ToLower(objective.Phase)).
		Count(&countThisPhase).Error; err != nil {
		return errors.New("Failed to validate secret scoring rules")
	}

	if countThisPhase > 0 {
		return errors.New("Player has already scored a secret objective in this phase this round")
	}

	// Total secret scoring cap
	var totalSecrets int64
	if err := database.DB.
		Model(&models.Score{}).
		Joins("JOIN objectives ON objectives.id = scores.objective_id").
		Where("scores.player_id = ? AND scores.game_id = ? AND LOWER(scores.type) = 'secret'", playerID, gameID).
		Count(&totalSecrets).Error; err != nil {
		return errors.New("Failed to count total secret objectives")
	}

	// Obsidian check
	var obsidianUsed int64
	if err := database.DB.
		Model(&models.Score{}).
		Where("game_id = ? AND player_id = ? AND LOWER(type) = 'relic' AND LOWER(relic_title) = 'the obsidian'", gameID, playerID).
		Count(&obsidianUsed).Error; err != nil {
		return errors.New("Failed to check Obsidian use")
	}

	maxSecrets := int64(3)
	if obsidianUsed > 0 {
		maxSecrets = 4
	}

	if totalSecrets >= maxSecrets {
		return fmt.Errorf("Player has already scored the maximum of %d secret objectives", maxSecrets)
	}

	return nil
}
