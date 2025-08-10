package services

import (
	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

func GetObjectivesByType(objectiveType string) ([]models.Objective, error) {
	var objs []models.Objective
	err := database.DB.Where("type = ?", objectiveType).Find(&objs).Error
	return objs, err
}
