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

	//player management
	r.POST("/players", controllers.CreatePlayer)
	r.GET("/players", controllers.ListPlayers)
	r.GET("/players/:id/games", controllers.GetPlayerGames)

	// game routes
	r.POST("/games", controllers.CreateGame)
	r.GET("/games/:id/players", controllers.ListPlayersInGame)
	r.GET("/games", controllers.ListGames)
	r.POST("/gameplayers", controllers.AssignPlayerToGame)
	r.POST("/games/:game_id/advance-round", controllers.AdvanceRound)
	r.GET("/games/:id/score-summary", controllers.GetScoreSummary)
	r.GET("/games/:id/scores-by-round", controllers.GetScoresByRound)
	r.GET("/games/:id", controllers.GetGameByID)
	r.GET("/games/:id/objectives", controllers.GetGameObjectives)
	//scoring
	r.POST("/score", controllers.AddScore)

	// Start server on port 8080
	r.Run("127.0.0.1:8080")
}
