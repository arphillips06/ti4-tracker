package services

import (
	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

func CreatePlayer(name string) (models.Player, error) {
	player := models.Player{Name: name}
	err := database.DB.Create(&player).Error
	return player, err
}

func GetPlayersInGame(gameID string) ([]models.GamePlayer, error) {
	var gamePlayers []models.GamePlayer
	err := database.DB.Where("game_id = ?", gameID).Preload("Player").Find(&gamePlayers).Error
	return gamePlayers, err
}

func GetGamesForPlayer(playerID string) (models.Player, error) {
	var player models.Player
	err := database.DB.
		Preload("Games.Game").
		Preload("Games.Game.GamePlayers.Player").
		First(&player, playerID).Error
	return player, err
}

func ListAllPlayers() ([]models.Player, error) {
	var players []models.Player
	err := database.DB.Find(&players).Error
	return players, err
}

func AssignPlayerToGame(gameID, playerID uint, faction string) (models.GamePlayer, error) {
	gp := models.GamePlayer{
		GameID:   gameID,
		PlayerID: playerID,
		Faction:  faction,
	}
	err := database.DB.Create(&gp).Error
	return gp, err
}
