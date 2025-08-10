package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/database"
	handle "github.com/arphillips06/TI4-stats/errors"
	achievements "github.com/arphillips06/TI4-stats/services/achievements"
	"github.com/gin-gonic/gin"
)

// GetGameAchievements godoc
// @Summary      Get achievements for a game
// @Description  Computes and returns per-game achievements (only for finished, non-partial games).
// @Tags         games
// @Param        id   path      int  true  "Game ID"
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "value, Count"
// @Failure      400  {object}  map[string]string       "error"
// @Failure      500  {object}  map[string]string       "error"
// @Router       /games/{id}/achievements [get]
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

// GetGlobalAchievements godoc
// @Summary      Global achievements (records)
// @Description  Returns current records across all finished, non-partial games, including all holders for each record.
// @Tags         achievements
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "value, Count"
// @Failure      500  {object}  map[string]string       "error"
// @Router       /achievements [get]
func GetGlobalAchievements(c *gin.Context) (int, any, error) {
	badges, err := achievements.ComputeGlobalAchievements(database.DB)
	if err != nil {
		return http.StatusInternalServerError, gin.H{"error": err.Error()}, nil
	}
	return http.StatusOK, gin.H{"value": badges, "Count": len(badges)}, nil
}
