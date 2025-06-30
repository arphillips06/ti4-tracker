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
	r.GET("/games/:id/players", controllers.ListPlayersInGame)
	r.GET("/players/:id/games", controllers.GetPlayerGames)
	r.POST("/score", controllers.AddScore)
	r.POST("/games/:game_id/advance-round", controllers.AdvanceRound)

	// Start server on port 8080
	r.Run(":8080")
}
