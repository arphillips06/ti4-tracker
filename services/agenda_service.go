package services

import (
	"errors"
	"fmt"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
	"gorm.io/gorm"
)

const (
	AgendaMutiny = "Mutiny"
)

// ApplyPoliticalCensure adjusts agenda score based on whether the player was censured or not.
// If Gained is false, a point is removed.
func ApplyPoliticalCensure(input models.PoliticalCensureRequest) error {
	points := 1
	if !input.Gained {
		points = -1
	}

	return helpers.CreateAgendaScore(int(input.GameID), int(input.RoundID), int(input.PlayerID), points, models.AgendaCensure, 0)
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
	totals, err := helpers.GetPlayerScoresMap(input.GameID)
	if err != nil {
		return err
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
	fmt.Printf("SeedOfEmpire totals: %+v\n", totals)
	if len(targetPlayerIDs) == 0 {
		return fmt.Errorf("no valid target players found for Seed of an Empire (%s)", input.Result)
	}

	for _, id := range targetPlayerIDs {
		if err := helpers.CreateAgendaScore(int(input.GameID), int(input.RoundID), int(id), 1, models.AgendaSeed, 0); err != nil {
			return err
		}
	}

	return nil
}

// ApplyMutinyAgenda awards or removes points based on the Mutiny agenda result.
func ApplyMutinyAgenda(input models.AgendaResolution) error {
	exists, err := helpers.AgendaAlreadyResolved(input.GameID, models.AgendaMutiny)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("mutiny has already been resolved for this game")
	}

	switch input.Result {
	case "for":
		for _, playerID := range input.ForVotes {
			if err := helpers.CreateAgendaScore(int(input.GameID), int(input.RoundID), int(playerID), 1, models.AgendaMutiny, 0); err != nil {
				return err
			}
		}
	case "against":
		for _, playerID := range input.ForVotes {
			total, err := helpers.GetPlayerTotalPoints(input.GameID, playerID)
			if err != nil {
				return err
			}
			if total > 0 {
				if err := helpers.CreateAgendaScore(int(input.GameID), int(input.RoundID), int(playerID), -1, models.AgendaMutiny, 0); err != nil {
					return err
				}
			}
		}
	default:
		return helpers.CreateAgendaScore(int(input.GameID), int(input.RoundID), 0, 0, models.AgendaMutiny, 0)
	}

	return nil
}

// This converts the scored secret objective to a public one.
// It also marks that it was originally secret, and records that CDL was used.
func ApplyClassifiedDocumentLeaks(input models.ClassifiedDocumentLeaksRequest) error {
	exists, err := helpers.AgendaAlreadyResolved(input.GameID, models.AgendaCDL)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("classified Document Leaks has already been resolved for this game")
	}

	// Locate the secret score
	var score models.Score
	err = database.DB.
		Where("game_id = ? AND player_id = ? AND objective_id = ? AND type = ?", input.GameID, input.PlayerID, input.ObjectiveID, models.ScoreTypeSecret).
		First(&score).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("secret objective score not found for that player")
		}
		return err
	}

	// Update the score to public
	score.Type = models.ScoreTypePublic
	score.OriginallySecret = true
	if err := database.DB.Save(&score).Error; err != nil {
		return err
	}

	return helpers.CreateAgendaScore(
		int(input.GameID),
		int(input.RoundID),
		int(input.PlayerID),
		0,
		models.AgendaCDL,
		input.ObjectiveID,
	)
}

// Incentive Program reveals the next unrevealed Stage I/II objective
// depending on the vote outcome: "for" → Stage I, "against" → Stage II
func ApplyIncentiveProgramEffect(gameID uint, outcome string) error {
	game, err := helpers.GetUnfinishedGame(gameID)
	if err != nil {
		return err // handles both not found and already finished
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
		return fmt.Errorf("invalid outcome: must be 'for' or 'against'")
	}

	var existingObjectiveIDs []uint
	if err := database.DB.
		Model(&models.GameObjective{}).
		Where("game_id = ?", gameID).
		Pluck("objective_id", &existingObjectiveIDs).Error; err != nil {
		return err
	}

	var newObjective models.Objective
	err = database.DB.
		Where("stage = ? AND id NOT IN ?", stage, existingObjectiveIDs).
		Order("id").
		First(&newObjective).Error
	if err != nil {
		return fmt.Errorf("no additional objectives remain in Stage %s", stage)
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

	return helpers.CreateAgendaScore(int(gameID), 0, 0, 0, models.AgendaIncentive, 0)

}
