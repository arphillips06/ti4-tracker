package services

import (
	"fmt"
	"log"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
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

	return nil
}
