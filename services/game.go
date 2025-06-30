package services

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/database/factions"
	"github.com/arphillips06/TI4-stats/models"
)

type CreateGameInput struct {
	WinningPoints int
	Players       []models.PlayerInput
}

type selectedPlayersWithFaction struct {
	Player  models.Player
	Faction string
}

func ParseAndValidatePlayers(inputPlayers []models.PlayerInput) ([]selectedPlayersWithFaction, error) {
	var allPlayers []models.Player
	if err := database.DB.Find(&allPlayers).Error; err != nil {
		return nil, err
	}

	playerMap := make(map[string]models.Player)
	for _, p := range allPlayers {
		playerMap[strconv.Itoa(int(p.ID))] = p
		playerMap[strings.ToLower(p.Name)] = p
	}

	var selected []selectedPlayersWithFaction
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

		selected = append(selected, selectedPlayersWithFaction{
			Player:  player,
			Faction: p.Faction,
		})
	}

	return selected, nil
}

func CreateGameAndRound(winningPoints int) (models.Game, models.Round, error) {
	game := models.Game{WinningPoints: winningPoints}
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

func AssignObjectivesToGame(game models.Game, round1 models.Round) error {
	var stage1 []models.Objective
	var stage2 []models.Objective

	database.DB.Where("stage = ?", "I").Find(&stage1)
	database.DB.Where("stage = ?", "II").Find(&stage2)

	if len(stage1) < 5 || len(stage2) < 5 {
		return errors.New("not enough public objectives in database")
	}

	rand.Shuffle(len(stage1), func(i, j int) { stage1[i], stage1[j] = stage1[j], stage1[i] })
	rand.Shuffle(len(stage2), func(i, j int) { stage2[i], stage2[j] = stage2[j], stage2[i] })

	for i, obj := range stage1[:5] {
		gameObj := models.GameObjective{
			GameID:      game.ID,
			ObjectiveID: obj.ID,
			Stage:       "I",
		}
		if i < 2 {
			gameObj.RoundID = round1.ID
		}
		database.DB.Create(&gameObj)
	}

	for _, obj := range stage2[:5] {
		gameObj := models.GameObjective{
			GameID:      game.ID,
			ObjectiveID: obj.ID,
			Stage:       "II",
		}
		database.DB.Create(&gameObj)
	}

	return nil
}
