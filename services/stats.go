package services

import (
	"fmt"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

type StatsOverview struct {
	TotalGames               int                       `json:"totalGames"`
	GamesWonByFaction        map[string]int            `json:"gamesWonByFaction"`
	GamesPlayedByFaction     map[string]int            `json:"gamesPlayedByFaction"`
	WinRateByFaction         map[string]float64        `json:"winRateByFaction"`
	ObjectiveStats           map[string]int            `json:"objectiveStats"`
	ObjectiveFrequency       map[string]int            `json:"objectiveFrequency"`
	PlayerWinRates           []PlayerWinRate           `json:"playerWinRates"`
	ObjectiveAppearanceStats map[string]ObjectiveStats `json:"objectiveAppearanceStats"`
}

type PlayerWinRate struct {
	Player      string  `json:"player"`
	GamesPlayed int     `json:"gamesPlayed"`
	GamesWon    int     `json:"gamesWon"`
	WinRate     float64 `json:"winRate"`
}
type ObjectiveStats struct {
	AppearanceRate         float64 `json:"appearanceRate"`
	ScoredWhenAppearedRate float64 `json:"scoredWhenAppearedRate"`
}

func CalculateStatsOverview() (*StatsOverview, error) {
	// Step 1: Get total games
	var totalGames int64
	if err := database.DB.Model(&models.Game{}).Count(&totalGames).Error; err != nil {
		return nil, err
	}

	// Step 2: Faction stats (played + won)
	factionPlays := make(map[string]int)
	factionWins := make(map[string]int)

	var factionData []struct {
		Faction string
		Count   int
	}

	if err := database.DB.
		Raw(`SELECT faction, COUNT(*) as count FROM game_players GROUP BY faction`).
		Scan(&factionData).Error; err != nil {
		return nil, err
	}
	for _, row := range factionData {
		factionPlays[row.Faction] = row.Count
	}

	if err := database.DB.
		Raw(`
            SELECT gp.faction, COUNT(*) as count
            FROM games g
            JOIN game_players gp ON g.id = gp.game_id AND g.winner_id = gp.player_id
            GROUP BY gp.faction
        `).
		Scan(&factionData).Error; err != nil {
		return nil, err
	}
	for _, row := range factionData {
		factionWins[row.Faction] = row.Count
	}

	winRates := make(map[string]float64)
	for faction, played := range factionPlays {
		if played > 0 {
			winRates[faction] = float64(factionWins[faction]) / float64(played) * 100
		} else {
			winRates[faction] = 0
		}
	}

	// Step 3: Objective stats
	var publicCount, secretCount, stage1Count, stage2Count int64
	database.DB.Model(&models.Score{}).Where("type = 'public'").Count(&publicCount)
	database.DB.Model(&models.Score{}).Where("type = 'secret'").Count(&secretCount)
	database.DB.Model(&models.Score{}).Joins("JOIN objectives ON scores.objective_id = objectives.id").
		Where("objectives.stage = ?", "I").Count(&stage1Count)
	database.DB.Model(&models.Score{}).Joins("JOIN objectives ON scores.objective_id = objectives.id").
		Where("objectives.stage = ?", "II").Count(&stage2Count)

	// Step 4: Frequency of objectives scored
	var freqRows []struct {
		Name  string
		Count int
	}
	database.DB.
		Table("scores").
		Select("objectives.name, COUNT(*) as count").
		Joins("JOIN objectives ON scores.objective_id = objectives.id").
		Group("objectives.name").
		Scan(&freqRows)

	freqMap := make(map[string]int)
	for _, row := range freqRows {
		freqMap[row.Name] = row.Count
	}

	// Step 5: Player win rates
	var playerRows []struct {
		Name        string
		GamesPlayed int
		GamesWon    int
	}
	err := database.DB.
		Raw(`
        SELECT p.name, 
               COUNT(DISTINCT gp.game_id) AS games_played,
               COUNT(DISTINCT CASE WHEN g.winner_id = gp.player_id THEN gp.game_id END) AS games_won
        FROM game_players gp
        JOIN players p ON p.id = gp.player_id
        JOIN games g ON g.id = gp.game_id
        GROUP BY p.name
    `).Scan(&playerRows).Error
	if err != nil {
		return nil, err
	}

	var playerWinRates []PlayerWinRate
	for _, p := range playerRows {
		winRate := 0.0
		if p.GamesPlayed > 0 {
			winRate = float64(p.GamesWon) / float64(p.GamesPlayed) * 100
		}
		playerWinRates = append(playerWinRates, PlayerWinRate{
			Player:      p.Name,
			GamesPlayed: p.GamesPlayed,
			GamesWon:    p.GamesWon,
			WinRate:     winRate,
		})
	}

	// Step 6: Objective appearance vs scored rates
	var appearances []struct {
		Name      string
		GameCount int
	}
	var scored []struct {
		Name      string
		GameCount int
	}

	database.DB.Raw(`
    SELECT o.name, COUNT(DISTINCT go.game_id) as game_count
    FROM game_objectives go
    JOIN objectives o ON go.objective_id = o.id
    GROUP BY o.name
	`).Scan(&appearances)
	database.DB.Raw(`
	SELECT o.name, COUNT(DISTINCT s.game_id) as game_count
	FROM scores s
	JOIN objectives o ON s.objective_id = o.id
	GROUP BY o.name
	`).Scan(&scored)

	objectiveAppearanceMap := make(map[string]int)
	objectiveScoredMap := make(map[string]int)
	fmt.Println("Appearances:")
	for _, row := range appearances {
		fmt.Printf("%s appeared in %d games\n", row.Name, row.GameCount)
	}

	fmt.Println("Scored:")
	for _, row := range scored {
		fmt.Printf("%s was scored in %d games\n", row.Name, row.GameCount)
	}

	for _, row := range appearances {
		objectiveAppearanceMap[row.Name] = row.GameCount
	}
	for _, row := range scored {
		objectiveScoredMap[row.Name] = row.GameCount
	}

	objectiveAppearanceStats := make(map[string]ObjectiveStats)
	allObjectiveNames := make(map[string]bool)

	// Combine keys from both maps
	for name := range objectiveAppearanceMap {
		allObjectiveNames[name] = true
	}
	for name := range objectiveScoredMap {
		allObjectiveNames[name] = true
	}

	for name := range allObjectiveNames {
		appeared := objectiveAppearanceMap[name]
		scored := objectiveScoredMap[name]

		var appearanceRate, scoreRate float64

		if totalGames > 0 {
			appearanceRate = float64(appeared) / float64(totalGames) * 100
		}
		if appeared > 0 {
			scoreRate = float64(scored) / float64(appeared) * 100
		}

		objectiveAppearanceStats[name] = ObjectiveStats{
			AppearanceRate:         appearanceRate,
			ScoredWhenAppearedRate: scoreRate,
		}
	}
	fmt.Printf("Final ObjectiveAppearanceStats (%d items):\n", len(objectiveAppearanceStats))
	for k, v := range objectiveAppearanceStats {
		fmt.Printf("  %s: %+v\n", k, v)
	}

	return &StatsOverview{
		TotalGames:           int(totalGames),
		GamesPlayedByFaction: factionPlays,
		GamesWonByFaction:    factionWins,
		WinRateByFaction:     winRates,
		ObjectiveStats: map[string]int{
			"publicScored": int(publicCount),
			"secretScored": int(secretCount),
			"stage1Scored": int(stage1Count),
			"stage2Scored": int(stage2Count),
		},
		ObjectiveFrequency:       freqMap,
		PlayerWinRates:           playerWinRates,
		ObjectiveAppearanceStats: objectiveAppearanceStats,
	}, nil
}
