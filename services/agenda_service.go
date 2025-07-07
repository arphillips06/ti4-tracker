package services

import (
	"errors"
	"fmt"
	"log"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
	"gorm.io/gorm"
)

func ApplyPoliticalCensure(input models.PoliticalCensureRequest) error {
	points := 1
	if !input.Gained {
		points = -1
	}

	score := models.Score{
		GameID:      input.GameID,
		RoundID:     input.RoundID,
		PlayerID:    input.PlayerID,
		Points:      points,
		Type:        "agenda", // lowercase to match existing entries
		AgendaTitle: "Political Censure",
	}

	if err := database.DB.Create(&score).Error; err != nil {
		return err
	}

	return nil
}

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

	// Step 3: Add up actual points from Score table
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

	if input.Result == "for" {
		targetPoints = -1
		for id, points := range totals {
			if points > targetPoints {
				targetPoints = points
				targetPlayerIDs = []uint{id}
			} else if points == targetPoints {
				targetPlayerIDs = append(targetPlayerIDs, id)
			}
		}
	} else if input.Result == "against" {
		targetPoints = 999
		for id, points := range totals {
			if points < targetPoints {
				targetPoints = points
				targetPlayerIDs = []uint{id}
			} else if points == targetPoints {
				targetPlayerIDs = append(targetPlayerIDs, id)
			}
		}
	}

	// Step 5: Apply the agenda score for each target player
	for _, id := range targetPlayerIDs {
		score := models.Score{
			GameID:      input.GameID,
			RoundID:     input.RoundID,
			PlayerID:    id,
			Points:      1,
			Type:        "agenda",
			AgendaTitle: "Seed of an Empire",
		}
		if err := database.DB.Create(&score).Error; err != nil {
			return err
		}
	}

	return nil
}

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

	if input.Result == "for" {
		log.Println("ApplyMutinyAgenda called with:", input)
		for _, playerID := range input.ForVotes {
			score := models.Score{
				GameID:      input.GameID,
				RoundID:     input.RoundID,
				PlayerID:    playerID,
				Points:      1,
				Type:        "agenda",
				AgendaTitle: "Mutiny",
			}

			if err := database.DB.Create(&score).Error; err != nil {
				return err
			}
		}
	}

	if input.Result == "against" {
		log.Println("ApplyMutinyAgenda called with:", input)
		for _, playerID := range input.ForVotes {
			var total int64
			err := database.DB.Model(&models.Score{}).
				Where("game_id = ? AND player_id = ?", input.GameID, playerID).
				Select("SUM(points)").Scan(&total).Error
			if err != nil {
				return err
			}

			if total > 0 {
				score := models.Score{
					GameID:      input.GameID,
					RoundID:     input.RoundID,
					PlayerID:    playerID,
					Points:      -1,
					Type:        "agenda",
					AgendaTitle: "Mutiny",
				}

				if err := database.DB.Create(&score).Error; err != nil {
					return err
				}
			}
		}
	}
	if input.Result != "for" && input.Result != "against" || len(input.ForVotes) == 0 {
		err := database.DB.Create(&models.Score{
			GameID:      input.GameID,
			RoundID:     input.RoundID,
			PlayerID:    0, // or nullable if supported
			ObjectiveID: 0,
			Points:      0,
			Type:        "agenda",
			AgendaTitle: "Mutiny",
		}).Error
		if err != nil {
			return fmt.Errorf("Failed to record Mutiny usage: %w", err)
		}
	}

	return nil
}

func ApplyClassifiedDocumentLeaks(input models.ClassifiedDocumentLeaksRequest) error {
	// Prevent duplicate resolution
	var count int64
	if err := database.DB.
		Model(&models.Score{}).
		Where("game_id = ? AND agenda_title = ?", input.GameID, "Classified Document Leaks").
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("Classified Document Leaks has already been resolved for this game")
	}

	log.Println("ApplyClassifiedDocumentLeaks called with:", input)

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
	score.Type = "public"
	if err := database.DB.Save(&score).Error; err != nil {
		return err
	}

	// Record the agenda use (0-point marker)
	agendaScore := models.Score{
		GameID:      input.GameID,
		PlayerID:    input.PlayerID,
		ObjectiveID: input.ObjectiveID,
		Points:      0,
		Type:        "agenda",
		AgendaTitle: "Classified Document Leaks",
	}
	if err := database.DB.Create(&agendaScore).Error; err != nil {
		return err
	}

	return nil
}

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
	_ = database.DB.Create(&models.Score{
		GameID:      gameID,
		PlayerID:    0,
		RoundID:     0,
		ObjectiveID: 0,
		Points:      0,
		Type:        "agenda",
		AgendaTitle: "Incentive Program",
	})

	return database.DB.Create(&gameObj).Error
}
