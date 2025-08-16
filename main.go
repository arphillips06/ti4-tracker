package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/arphillips06/TI4-stats/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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
	docs.SwaggerInfo.Title = "TI4 Stats API"
	docs.SwaggerInfo.Version = "0.1"
	docs.SwaggerInfo.Description = "Endpoints for TI4-stats backend."
	docs.SwaggerInfo.BasePath = "/"

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
	r.POST("/players", controllers.Wrap(controllers.CreatePlayer))
	r.GET("/players", controllers.Wrap(controllers.ListPlayers))
	r.GET("/players/:id/games", controllers.Wrap(controllers.GetPlayerGames))

	// game routes
	r.POST("/games", controllers.Wrap(controllers.CreateGame))
	r.GET("/games/:id/players", controllers.Wrap(controllers.ListPlayersInGame))
	r.GET("/games", controllers.Wrap(controllers.ListGames))
	r.POST("/gameplayers", controllers.Wrap(controllers.AssignPlayerToGame))
	r.POST("/games/:game_id/advance-round", controllers.Wrap(controllers.AdvanceRound))
	r.GET("/games/:id/score-summary", controllers.Wrap(controllers.GetScoreSummary))
	r.GET("/games/:id/scores-by-round", controllers.Wrap(controllers.GetScoresByRound))
	r.GET("/games/:id", controllers.Wrap(controllers.GetGameByID))
	r.GET("/games/:id/objectives", controllers.Wrap(controllers.GetGameObjectives))
	r.GET("/objectives/secrets/all", controllers.Wrap(controllers.GetAllSecretObjectives))
	r.POST("/assign_objective", controllers.Wrap(controllers.AssignObjective))
	r.GET("/objectives/public/all", controllers.Wrap(controllers.GetAllPublicObjectives))
	r.GET("/api/games/:id/exists", controllers.Wrap(controllers.GetGameExists))
	r.POST("/game/:id/randomise-speaker", controllers.Wrap(controllers.RandomiseSpeaker))
	r.POST("/games/:game_id/speaker", controllers.Wrap(controllers.PostAssignSpeaker))

	//scoring
	r.POST("/score", controllers.Wrap(controllers.AddScore))
	r.POST("/score/imperial", controllers.Wrap(controllers.ScoreImperialPoint))
	r.POST("/score/mecatol", controllers.Wrap(controllers.ScoreMecatolPoint))
	r.POST("/score/imperial-rider", controllers.Wrap(controllers.ScoreImperialRiderPoint))
	r.GET("/games/:id/objectives/scores", controllers.Wrap(controllers.GetObjectiveScoreSummary))
	r.POST("/unscore", controllers.Wrap(controllers.DeleteScore))

	//expose factions to API
	r.GET("/api/factions", controllers.Wrap(controllers.GetFactions))

	//agendas
	r.POST("/agenda/mutiny", controllers.ResolveMutinyAgenda)
	r.POST("/agenda/political-censure", controllers.HandlePoliticalCensure)
	r.POST("/agenda/seed", controllers.HandleSeedOfEmpire)
	r.POST("/agenda/classified-document-leaks", controllers.HandleClassifiedDocumentLeaks)
	r.POST("/agenda/incentive-program", controllers.HandleIncentiveProgram)

	//relics
	r.POST("/relic/shard", controllers.Wrap(controllers.HandleShardRelic))
	r.POST("/relic/crown", controllers.Wrap(controllers.HandleCrownRelic))
	r.POST("/relic/obsidian", controllers.Wrap(controllers.HandleObsidianRelic))

	r.POST("/games/:game_id/support/:player_id", controllers.Wrap(controllers.SFTT))
	r.GET("/stats/overview", controllers.Wrap(controllers.GetStatsOverview))

	// Serve static frontend files from /build
	r.Static("/static", "./build/static") // serve JS/CSS etc.

	// Serve index.html on root and fallback for SPA routing
	r.GET("/", func(c *gin.Context) {
		c.File("./build/index.html")
	})

	r.GET("/games/:id/achievements", controllers.Wrap(controllers.GetGameAchievements))
	r.GET("/achievements", controllers.Wrap(controllers.GetGlobalAchievements))

	//swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
