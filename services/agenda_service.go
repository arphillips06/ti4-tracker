package services

import (
	"errors"
	"fmt"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
	"gorm.io/gorm"
)

// ApplyPoliticalCensure adjusts agenda score based on whether the player was censured or not.
// If Gained is false, a point is removed.
func ApplyPoliticalCensure(input models.PoliticalCensureRequest) error {
	points := 1
	if !input.Gained {
		points = -1
	}

	return helpers.CreateAgendaScore(int(input.GameID), int(input.RoundID), int(input.PlayerID), points, models.AgendaCensure)
}

// ApplySeedOfEmpire awards 1 point to the player with most (or fewest) points depending on the vote result.
// Ties are handled by awarding all tied players.
func ApplySeedOfEmpire(input models.SeedOfEmpireResolution) error {
	// Step 1: Get all players in the game
	var gamePlayers []models.GamePlayer
	if err := database.DB.Where("game_id = ?", input.GameID).Find(&gamePlayers).Error; err != nil {
		return err
	}

	// Step 2: Initialize totals to 0
	totals := make(map[uint]int)
	for _, gp := range gamePlayers {
		totals[gp.PlayerID] = 0
	}

	// Step 3: Sum all current player scores from Score table
	var scores []models.Score
	if err := database.DB.Where("game_id = ?", input.GameID).Find(&scores).Error; err != nil {
		return err
	}
	for _, score := range scores {
		totals[score.PlayerID] += score.Points
	}

	// Step 4: Determine the target player(s)
	var targetPlayerIDs []uint
	var targetPoints int

	switch input.Result {
	case "for":
		targetPoints = -1
		for id, points := range totals {
			if points > targetPoints {
				targetPoints = points
				targetPlayerIDs = []uint{id}
			} else if points == targetPoints {
				targetPlayerIDs = append(targetPlayerIDs, id)
			}
		}
	case "against":
		targetPoints = 999
		for id, points := range totals {
			if points < targetPoints {
				targetPoints = points
				targetPlayerIDs = []uint{id}
			} else if points == targetPoints {
				targetPlayerIDs = append(targetPlayerIDs, id)
			}
		}
	default:
		return fmt.Errorf("invalid vote result: %s", input.Result)
	}

	for _, id := range targetPlayerIDs {
		if err := helpers.CreateAgendaScore(int(input.GameID), int(input.RoundID), int(id), 1, models.AgendaSeed); err != nil {
			return err
		}
	}

	return nil
}

// ApplyMutinyAgenda awards or removes points based on the Mutiny agenda result.
func ApplyMutinyAgenda(input models.AgendaResolution) error {
	var count int64
	if err := database.DB.
		Model(&models.Score{}).
		Where("game_id = ? AND agenda_title = ?", input.GameID, "Mutiny").
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("Mutiny has already been resolved for this game")
	}

	switch input.Result {
	case "for":
		for _, playerID := range input.ForVotes {
			if err := helpers.CreateAgendaScore(int(input.GameID), int(input.RoundID), int(playerID), 1, models.AgendaMutiny); err != nil {
				return err
			}
		}
	case "against":
		for _, playerID := range input.ForVotes {
			var total int64
			err := database.DB.Model(&models.Score{}).
				Where("game_id = ? AND player_id = ?", input.GameID, playerID).
				Select("SUM(points)").Scan(&total).Error
			if err != nil {
				return err
			}
			if total > 0 {
				if err := helpers.CreateAgendaScore(int(input.GameID), int(input.RoundID), int(playerID), -1, models.AgendaMutiny); err != nil {
					return err
				}
			}
		}
	default:
		return helpers.CreateAgendaScore(int(input.GameID), int(input.RoundID), 0, 0, models.AgendaMutiny)
	}

	return nil
}

// This converts the scored secret objective to a public one.
// It also marks that it was originally secret, and records that CDL was used.
func ApplyClassifiedDocumentLeaks(input models.ClassifiedDocumentLeaksRequest) error {
	// Prevent duplicate resolution
	var count int64
	if err := database.DB.
		Model(&models.Score{}).
		Where("game_id = ? AND agenda_title = ?", input.GameID, models.AgendaCDL).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("Classified Document Leaks has already been resolved for this game")
	}

	// Locate the secret score
	var score models.Score
	err := database.DB.Where("game_id = ? AND player_id = ? AND objective_id = ? AND type = ?", input.GameID, input.PlayerID, input.ObjectiveID, "secret").
		First(&score).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("Secret objective score not found for that player")
		}
		return err
	}

	// Update the score to public
	score.Type = models.ScoreTypePublic
	score.OriginallySecret = true
	if err := database.DB.Save(&score).Error; err != nil {
		return err
	}

	return helpers.CreateAgendaScore(int(input.GameID), int(input.RoundID), int(input.PlayerID), 0, models.AgendaCDL)
}

// Incentive Program reveals the next unrevealed Stage I/II objective
// depending on the vote outcome: "for" → Stage I, "against" → Stage II
func ApplyIncentiveProgramEffect(gameID uint, outcome string) error {
	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		return fmt.Errorf("Game not found")
	}

	if !game.UseObjectiveDecks {
		return nil // Manual mode: do nothing, admin will assign
	}

	var stage string
	switch outcome {
	case "for":
		stage = "I"
	case "against":
		stage = "II"
	default:
		return fmt.Errorf("Invalid outcome: must be 'for' or 'against'")
	}

	var existingObjectiveIDs []uint
	if err := database.DB.
		Model(&models.GameObjective{}).
		Where("game_id = ?", gameID).
		Pluck("objective_id", &existingObjectiveIDs).Error; err != nil {
		return err
	}

	var newObjective models.Objective
	err := database.DB.
		Where("stage = ? AND id NOT IN ?", stage, existingObjectiveIDs).
		Order("id").
		First(&newObjective).Error
	if err != nil {
		return fmt.Errorf("No additional objectives remain in Stage %s", stage)
	}

	// Find the next unrevealed objective in this stage
	gameObj := models.GameObjective{
		GameID:      gameID,
		ObjectiveID: newObjective.ID,
		Stage:       newObjective.Stage,
		RoundID:     0,
		Revealed:    true,
	}
	if err := database.DB.Create(&gameObj).Error; err != nil {
		return err
	}

	return helpers.CreateAgendaScore(int(gameID), 0, 0, 0, models.AgendaIncentive)

}
