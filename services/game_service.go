package services

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/database/factions"
	"github.com/arphillips06/TI4-stats/models"
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
		lookup := strings.ToLower(p.Name)
		if p.ID != "" {
			lookup = p.ID
		}

		player, exists := playerMap[lookup]
		if !exists {
			return nil, fmt.Errorf("player not found: %s", p.Name)
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
		}
		if err := database.DB.Create(&gameObj).Error; err != nil {
			return err
		}
	}
	for _, obj := range selectedStage2 {
		gameObj := models.GameObjective{
			GameID:      game.ID,
			ObjectiveID: obj.ID,
			RoundID:     0,
			Stage:       obj.Stage,
		}
		if err := database.DB.Create(&gameObj).Error; err != nil {
			return err
		}
	}
	return nil
}

// Gets a game by its string ID
func GetGameByID(id string) (*models.Game, error) {
	var game models.Game
	if err := database.DB.First(&game, id).Error; err != nil {
		return nil, err
	}
	return &game, nil
}

// Creates and advances to a new round
func CreateNewRound(game *models.Game) (*models.Round, error) {
	newRound := models.Round{
		GameID: game.ID,
		Number: game.CurrentRound + 1,
	}
	if err := database.DB.Create(&newRound).Error; err != nil {
		return nil, err
	}
	game.CurrentRound = newRound.Number
	if err := database.DB.Save(&game).Error; err != nil {
		return nil, err
	}
	return &newRound, nil
}

// Determines if we should reveal a Stage I or Stage II objective this round
func DetermineStageToReveal(gameID uint) string {
	var count int64
	database.DB.Model(&models.GameObjective{}).
		Where("game_id = ? AND stage = ? AND round_id > 0", gameID, "I").
		Count(&count)
	if count >= 5 {
		return "II"
	}
	return "I"
}

// Marks the next unrevealed objective of the given stage as revealed in the current round
func RevealNextObjective(gameID, roundID uint, stage string) error {
	var obj models.GameObjective
	err := database.DB.
		Where("game_id = ? AND round_id = 0 AND stage = ?", gameID, stage).
		First(&obj).Error
	if err != nil {
		return err
	}
	obj.RoundID = roundID
	return database.DB.Save(&obj).Error
}

// Counts total number of revealed public objectives for a game
func CountRevealedObjectives(gameID uint) int64 {
	var count int64
	database.DB.Model(&models.GameObjective{}).
		Where("game_id = ? AND round_id > 0", gameID).
		Count(&count)
	return count
}
