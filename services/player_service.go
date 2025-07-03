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
