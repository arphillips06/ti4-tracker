package services

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
	"gorm.io/gorm"
)

func ValidateSecretScoringRules(gameID, playerID, roundID, objectiveID uint) error {
	var objective models.Objective
	if err := database.DB.First(&objective, objectiveID).Error; err != nil {
		log.Printf("[ERROR] Could not find objective %d: %v", objectiveID, err)
		return errors.New("Objective not found")
	}

	if strings.ToLower(objective.Type) != models.ScoreTypeSecret {
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

func CheckIfScoreExists(gameID, playerID, objectiveID uint) (bool, error) {
	var existing models.Score
	err := database.DB.
		Where("game_id = ? AND player_id = ? AND objective_id = ?", gameID, playerID, objectiveID).
		First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check for existing score: %w", err)
	}

	return true, nil
}
