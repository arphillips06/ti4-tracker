package services

import (
	"errors"
	"fmt"
	"log"
	"math/rand"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
	"gorm.io/gorm"
)

func AssignSpeaker(gameID, roundNumber, playerID uint) error {
	var round models.Round
	if err := database.DB.
		Where("game_id = ? AND number = ?", gameID, roundNumber).
		First(&round).Error; err != nil {
		return fmt.Errorf("could not find round %d for game %d: %w", roundNumber, gameID, err)
	}

	var player models.GamePlayer
	if err := database.DB.First(&player, playerID).Error; err != nil {
		return err
	}
	if player.GameID != gameID {
		return fmt.Errorf("player does not belong to this game")
	}

	var existing models.SpeakerAssignment
	err := database.DB.Where("game_id = ? AND round_id = ?", gameID, round.ID).First(&existing).Error
	if err == nil {
		existing.PlayerID = playerID
		return database.DB.Save(&existing).Error
	}

	sa := models.SpeakerAssignment{
		GameID:   gameID,
		RoundID:  round.ID,
		PlayerID: playerID,
	}
	return database.DB.Create(&sa).Error
}

func RandomiseSpeaker(gameID uint) (*models.Player, error) {
	var players []models.Player
	if err := database.DB.Where("game_id = ?", gameID).Find(&players).Error; err != nil {
		return nil, errors.New("failed to fetch players")
	}

	if len(players) == 0 {
		return nil, errors.New("no players found for this game")
	}

	chosen := players[rand.Intn(len(players))]

	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		return nil, errors.New("failed to fetch game")
	}

	var round models.Round
	err := database.DB.
		Where("game_id = ?", gameID).
		Order("number ASC").
		First(&round).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		round = models.Round{GameID: gameID, Number: 1}
		if err := database.DB.Create(&round).Error; err != nil {
			return nil, errors.New("failed to create round 1")
		}
		log.Println("🆕 Created round 1 for game", gameID)
	} else if err != nil {
		return nil, errors.New("failed to fetch round 1")
	}

	assignment := models.SpeakerAssignment{
		GameID:   gameID,
		RoundID:  round.ID,
		PlayerID: chosen.ID,
	}
	if err := database.DB.Create(&assignment).Error; err != nil {
		return nil, errors.New("failed to create speaker assignment")
	}

	if err := database.DB.Model(&models.Game{}).Where("id = ?", gameID).
		Update("speaker_id", chosen.ID).Error; err != nil {
		return nil, errors.New("failed to update game speaker")
	}

	return &chosen, nil
}
