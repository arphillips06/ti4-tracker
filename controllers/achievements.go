package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/database"
	handle "github.com/arphillips06/TI4-stats/errors"
	achievements "github.com/arphillips06/TI4-stats/services/achievements"
	"github.com/gin-gonic/gin"
)

func GetGameAchievements(c *gin.Context) (int, any, error) {
	id, err := handle.ParseID(c, "id")
	if err != nil {
		return 0, nil, err
	}
	badges, err := achievements.ComputeGameAchievements(database.DB, id)
	if err != nil {
		return 0, nil, err
	}
	if badges == nil {
		badges = []achievements.Badge{}
	}
	return http.StatusOK, gin.H{"value": badges, "Count": len(badges)}, nil
}
