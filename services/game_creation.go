package services

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/database/factions"
	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
	"gorm.io/gorm"
)

// Validates player input and returns matched players with faction info.
func ParseAndValidatePlayers(inputPlayers []models.PlayerInput) ([]models.SelectedPlayersWithFaction, error) {
	var allPlayers []models.Player
	if err := database.DB.Find(&allPlayers).Error; err != nil {
		return nil, err
	}
	// Map for quick lookup by ID or lowercase name
	playerMap := make(map[string]models.Player)
	for _, p := range allPlayers {
		playerMap[strconv.Itoa(int(p.ID))] = p
		playerMap[strings.ToLower(p.Name)] = p
	}

	var selected []models.SelectedPlayersWithFaction
	for _, p := range inputPlayers {
		if strings.TrimSpace(p.Name) == "" {
			return nil, fmt.Errorf("player name cannot be blank")
		}
		lookup := strings.ToLower(p.Name)
		if p.ID != "" {
			lookup = p.ID
		}

		player, exists := playerMap[lookup]
		if !exists {
			newplayer, err := CreatePlayer(p.Name)
			if err != nil {
				return nil, fmt.Errorf("failed to create player: %s", p.Name)
			}
			player = newplayer
		}

		if !factions.IsValidFaction(p.Faction) {
			return nil, fmt.Errorf("invalid faction: %s", p.Faction)
		}

		selected = append(selected, models.SelectedPlayersWithFaction{
			Player:  player,
			Faction: p.Faction,
		})
	}

	return selected, nil
}

// Creates a new game and initial round
func CreateGameAndRound(winningPoints int, useDecks bool) (models.Game, models.Round, error) {
	game := models.Game{
		WinningPoints:     winningPoints,
		UseObjectiveDecks: useDecks,
	}
	if err := database.DB.Create(&game).Error; err != nil {
		return game, models.Round{}, err
	}

	round1 := models.Round{GameID: game.ID, Number: 1}
	if err := database.DB.Create(&round1).Error; err != nil {
		return game, models.Round{}, err
	}

	game.CurrentRound = 1
	if err := database.DB.Save(&game).Error; err != nil {
		return game, round1, err
	}

	return game, round1, nil
}

// Assigns 10 public objectives (5 stage I, 5 stage II) to a game.
// Two stage I objectives are revealed in round 1.
func AssignObjectivesToGame(game models.Game, round1 models.Round) error {
	var stage1 []models.Objective
	var stage2 []models.Objective

	database.DB.Where("stage = ?", "I").Find(&stage1)
	database.DB.Where("stage = ?", "II").Find(&stage2)

	rand.Shuffle(len(stage1), func(i, j int) { stage1[i], stage1[j] = stage1[j], stage1[i] })
	rand.Shuffle(len(stage2), func(i, j int) { stage2[i], stage2[j] = stage2[j], stage2[i] })

	selectedStage1 := stage1[:5]
	selectedStage2 := stage2[:5]

	for i, obj := range selectedStage1 {
		roundID := uint(0)
		if i < 2 {
			roundID = round1.ID
		}

		gameObj := models.GameObjective{
			GameID:      game.ID,
			ObjectiveID: obj.ID,
			RoundID:     roundID,
			Stage:       obj.Stage,
			Revealed:    i < 2,
			Position:    i,
		}
		if err := database.DB.Create(&gameObj).Error; err != nil {
			return err
		}
	}
	for j, obj := range selectedStage2 {
		gameObj := models.GameObjective{
			GameID:      game.ID,
			ObjectiveID: obj.ID,
			RoundID:     0,
			Stage:       obj.Stage,
			Revealed:    false,
			Position:    j,
		}
		if err := database.DB.Create(&gameObj).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateNewGameWithPlayers(input models.CreateGameInput) (models.Game, []models.GameObjective, error) {
	const (
		DefaultWinningPoints   = 10
		AlternateWinningPoints = 14
	)

	useDecks := true
	if input.UseObjectiveDecks != nil {
		useDecks = *input.UseObjectiveDecks
	}
	if input.WinningPoints != DefaultWinningPoints && input.WinningPoints != AlternateWinningPoints {
		input.WinningPoints = DefaultWinningPoints
	}

	selected, err := ParseAndValidatePlayers(input.Players)
	if err != nil {
		return models.Game{}, nil, err
	}

	game, round1, err := CreateGameAndRound(input.WinningPoints, useDecks)
	if err != nil {
		return models.Game{}, nil, err
	}

	for _, entry := range selected {
		if err := helpers.CreateGamePlayer(game.ID, entry.Player.ID, entry.Faction); err != nil {
			return models.Game{}, nil, err
		}
	}

	if err := database.DB.First(&game, game.ID).Error; err != nil {
		return models.Game{}, nil, errors.New("failed to reload game")
	}
	var gamePlayers []models.GamePlayer
	if err := database.DB.Preload("Player").Where("game_id = ?", game.ID).Find(&gamePlayers).Error; err != nil {
		return models.Game{}, nil, errors.New("failed to load game players for speaker assignment")
	}

	log.Printf("💬 UseRandomSpeaker: %v", input.UseRandomSpeaker)
	log.Printf("🎲 Players loaded: %v", gamePlayers)

	if input.UseRandomSpeaker != nil && *input.UseRandomSpeaker {
		if len(gamePlayers) > 0 {
			chosen := gamePlayers[rand.Intn(len(gamePlayers))]
			log.Printf("🎙️  Chosen speaker: %v", chosen)
			game.SpeakerID = &chosen.ID
		}
	}

	log.Printf("💾 SpeakerID set to: %v", game.SpeakerID)

	if err := database.DB.Save(&game).Error; err != nil {
		return models.Game{}, nil, errors.New("failed to save speaker assignment")
	}

	var revealed []models.GameObjective
	if game.UseObjectiveDecks {
		if err := AssignObjectivesToGame(game, round1); err != nil {
			return models.Game{}, nil, err
		}
		_ = database.DB.
			Preload("Objective").
			Joins("JOIN rounds ON rounds.id = game_objectives.round_id").
			Where("game_objectives.game_id = ?", game.ID).
			Find(&revealed)
	}

	return game, revealed, nil
}

func ManuallyAssignObjective(gameID uint, roundNumber uint, objectiveID uint) error {
	var round models.Round
	if err := database.DB.
		Where("game_id = ? AND number = ?", gameID, roundNumber).
		First(&round).Error; err != nil {
		return errors.New("round not found")
	}

	var obj models.Objective
	if err := database.DB.
		First(&obj, objectiveID).Error; err != nil {
		return errors.New("objective not found")
	}

	var existing models.GameObjective
	err := database.DB.
		Where("game_id = ? AND objective_id = ?", gameID, obj.ID).
		First(&existing).Error

	if err == nil {
		return errors.New("objective already assigned to this game")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	var position int64
	_ = database.DB.Model(&models.GameObjective{}).
		Where("game_id = ? AND stage = ?", gameID, obj.Stage).
		Count(&position)

	reveal := models.GameObjective{
		GameID:      gameID,
		ObjectiveID: obj.ID,
		RoundID:     round.ID,
		Stage:       obj.Stage,     // e.g., "I" or "II"
		Revealed:    true,          // so it shows up in /objectives
		Position:    int(position), // controls display order
	}
	log.Printf("Assigned objective %s (ID %d) to game %d round %d", obj.Name, obj.ID, gameID, roundNumber)

	return database.DB.Create(&reveal).Error
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

	// Update the Game's SpeakerID
	if err := database.DB.Model(&models.Game{}).Where("id = ?", gameID).
		Update("speaker_id", chosen.ID).Error; err != nil {
		return nil, errors.New("failed to update speaker")
	}

	return &chosen, nil
}
