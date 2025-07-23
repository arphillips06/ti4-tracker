package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/arphillips06/TI4-stats/controllers"
	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/helpers/stats"
	"github.com/arphillips06/TI4-stats/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize DB and seed objectives
	database.InitDatabase()
	database.SeedObjectives()

	// Setup Gin router
	r := gin.Default()
	pathCounts, err := stats.CalculateCommonVictoryPaths()
	if err != nil {
		log.Printf("Could not preload victory paths: %v", err)
		pathCounts = make(map[string]int)
	}
	services.CachedVictoryPathCounts = pathCounts

	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if strings.HasPrefix(origin, "http://192.168.1.") ||
			strings.HasPrefix(origin, "http://100.") ||
			origin == "http://localhost:3000" ||
			strings.HasSuffix(origin, ".ts.net:3000") ||
			strings.HasSuffix(origin, "ross-lab.org:3000") ||
			origin == "http://[2a01:4b00:bf28:bb00:be24:11ff:fe04:d1a5]" ||
			origin == "http://2a01:4b00:bf28:bb00:be24:11ff:fe04:d1a5" {

			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

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
	r.GET("/objectives/secrets/all", controllers.GetAllSecretObjectives)
	r.POST("/assign_objective", controllers.AssignObjective)
	r.GET("/objectives/public/all", controllers.GetAllPublicObjectives)
	r.GET("/api/games/:id/exists", controllers.GetGameExists)

	//scoring
	r.POST("/score", controllers.AddScore)
	r.POST("/score/imperial", controllers.ScoreImperialPoint)
	r.POST("/score/mecatol", controllers.ScoreMecatolPoint)
	r.POST("/score/imperial-rider", controllers.ScoreImperialRiderPoint)
	r.GET("/games/:id/objectives/scores", controllers.GetObjectiveScoreSummary)
	r.POST("/unscore", controllers.DeleteScore)

	//expose factions to API
	r.GET("/api/factions", controllers.GetFactions)

	//agendas
	r.POST("/agenda/mutiny", controllers.ResolveMutinyAgenda)
	r.POST("/agenda/political-censure", controllers.HandlePoliticalCensure)
	r.POST("/agenda/seed", controllers.HandleSeedOfEmpire)
	r.POST("/agenda/classified-document-leaks", controllers.HandleClassifiedDocumentLeaks)
	r.POST("/agenda/incentive-program", controllers.HandleIncentiveProgram)

	//relics
	r.POST("/relic/shard", controllers.HandleShardRelic)
	r.POST("/relic/crown", controllers.HandleCrownRelic)
	r.POST("/relic/obsidian", controllers.HandleObsidianRelic)

	r.POST("/games/:game_id/support/:player_id", controllers.SFTT)
	r.GET("/stats/overview", controllers.GetStatsOverview)

	// Serve static frontend files from /build
	r.Static("/static", "./build/static") // serve JS/CSS etc.

	// Serve index.html on root and fallback for SPA routing
	r.GET("/", func(c *gin.Context) {
		c.File("./build/index.html")
	})

	// For any unmatched route (client side routing), serve index.html
	r.NoRoute(func(c *gin.Context) {
		// Only serve index.html for paths that are not API routes
		if !strings.HasPrefix(c.Request.URL.Path, "/api") && !strings.HasPrefix(c.Request.URL.Path, "/games") && !strings.HasPrefix(c.Request.URL.Path, "/players") {
			c.File("./build/index.html")
		} else {
			c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		}
	})

	// Start server on port 8080
	bindAddr := os.Getenv("BIND_ADDRESS")
	if bindAddr == "" {
		bindAddr = "127.0.0.1:8080" // default for dev
	}

	r.Run(bindAddr)

}
