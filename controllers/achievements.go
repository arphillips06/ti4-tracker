package controllers

import (
	"net/http"
	"strconv"

	"github.com/arphillips06/TI4-stats/database"
	achievements "github.com/arphillips06/TI4-stats/services/acheivements"
	"github.com/gin-gonic/gin"
)

func GetPlayerAchievements(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player id"})
		return
	}
	badges, err := achievements.GetPlayerAchievements(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// always return [] not null
	if badges == nil {
		badges = []achievements.AchievementBadge{}
	}
	c.JSON(http.StatusOK, badges)
}

func GetGameAchievements(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game id"})
		return
	}
	badges, err := achievements.GetGameAchievements(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if badges == nil {
		badges = []achievements.AchievementBadge{}
	}
	c.JSON(http.StatusOK, badges)
}

// Optional admin/debug endpoints (nice to keep):
func RecomputeAchievements(c *gin.Context) {
	n, err := achievements.RecomputeAllFinishedGames(database.DB, achievements.EvaluateAchievementsAfterGame)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "recomputed": n})
}
