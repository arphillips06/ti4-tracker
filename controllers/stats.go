package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/services"
	"github.com/gin-gonic/gin"
)

func GetStatsOverview(c *gin.Context) {
	overview, err := services.CalculateStatsOverview()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate overview stats"})
		return
	}

	custodians, err := services.GetPlayerCustodiansStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate custodians stats"})
		return
	}

	overview.CustodiansStats = custodians
	c.JSON(http.StatusOK, overview)
}
