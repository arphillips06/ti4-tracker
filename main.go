package main

import (
	"github.com/arphillips06/TI4-stats/controllers"
	"github.com/arphillips06/TI4-stats/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize DB and seed objectives
	database.InitDatabase()
	database.SeedObjectives()

	// Setup Gin router
	r := gin.Default()

	// Player routes
	r.POST("/players", controllers.CreatePlayer)
	r.POST("/games", controllers.CreateGame)
	r.POST("/gameplayers", controllers.AssignPlayerToGame)

	// Start server on port 8080
	r.Run(":8080")
}
