package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/database/factions"
	"github.com/gin-gonic/gin"
)

// GetFactions godoc
// @Summary      List factions
// @Description  Returns all available factions.
// @Tags         factions
// @Produce      json
// @Success      200  {array}   map[string]interface{}
// @Failure      500  {object}  map[string]string  "error"
// @Router       /factions [get]
func GetFactions(c *gin.Context) (int, any, error) {
	return http.StatusOK, factions.AllFactions, nil
}
