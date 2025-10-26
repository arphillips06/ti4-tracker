package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/services"
	"github.com/gin-gonic/gin"
)

// GetStatsOverview godoc
// @Summary      Stats overview
// @Description  Returns headline stats plus Custodians (Mecatol) stats per player.
// @Tags         stats
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string  "error"
// @Router       /stats/overview [get]
func GetStatsOverview(c *gin.Context) (int, any, error) {
	overview, err := services.CalculateStatsOverview()
	if err != nil {
		return 0, nil, fmt.Errorf("failed to generate overview stats: %w", err)
	}

	custodians, err := services.GetPlayerCustodiansStats()
	if err != nil {
		return 0, nil, fmt.Errorf("failed to generate custodians stats: %w", err)
	}

	overview.CustodiansStats = custodians
	return http.StatusOK, overview, nil
}

// GetObjectiveDifficulty godoc
// @Summary      Get objective difficulty
// @Description  Calculates and returns difficulty metrics for TI4 objectives.
// @Tags         objectives, stats
// @Produce      json
// @Param        stage             query   string  false  "Filter by stage (I, II, secret, or all)"  default(all)
// @Param        minAppearances    query   int     false  "Minimum appearances required to include"  default(5)
// @Param        minOpportunities  query   int     false  "Minimum scoring opportunities required"   default(0)
// @Success      200  {object}  models.ObjectiveDifficultyResponse
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /stats/objectives/difficulty [get]
func GetObjectiveDifficulty(c *gin.Context) (int, any, error) {
	stage := c.DefaultQuery("stage", "all")
	minApp := parseIntDefault(c.Query("minAppearances"), 5)
	minOpp := parseIntDefault(c.Query("minOpportunities"), 0)

	res, err := services.CalculateObjectiveDifficulty(
		c.Request.Context(),
		database.DB,
		services.ObjectiveDifficultyOptions{
			Stage:            stage,
			MinAppearances:   minApp,
			MinOpportunities: minOpp,
		},
	)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, res, nil
}

func parseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
