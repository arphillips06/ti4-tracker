package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/database/factions"
	"github.com/gin-gonic/gin"
)

func GetFactions(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, factions.AllFactions)
}
