package controllers

import (
	"fmt"
	"net/http"

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
