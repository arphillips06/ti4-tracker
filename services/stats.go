package services

import (
	"fmt"
	"math"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

type StatsOverview struct {
	TotalGames                 int                           `json:"totalGames"`
	GamesWonByFaction          map[string]int                `json:"gamesWonByFaction"`
	GamesPlayedByFaction       map[string]int                `json:"gamesPlayedByFaction"`
	WinRateByFaction           map[string]float64            `json:"winRateByFaction"`
	ObjectiveStats             map[string]int                `json:"objectiveStats"`
	ObjectiveFrequency         map[string]int                `json:"objectiveFrequency"`
	PlayerWinRates             []PlayerWinRate               `json:"playerWinRates"`
	ObjectiveAppearanceStats   map[string]ObjectiveStats     `json:"objectiveAppearanceStats"`
	FactionPlayWinDistribution map[string]FactionPlayWinStat `json:"factionPlayWinDistribution"`
	PlayerAveragePoints        []PlayerAveragePoints         `json:"playerAveragePoints"`
	TopFactionsPerPlayer       []PlayerFactionStats          `json:"topFactionsPerPlayer"`
	PlayerMostCommonFinishes   []PlayerMostCommonFinish      `json:"playerMostCommonFinishes"`
	SecretObjectiveRates       []SecretObjectiveRate         `json:"secretObjectiveRates"`
	PlayerPointStdevs          []PlayerPointStdev            `json:"playerPointStdevs"`
	AverageRounds              float64                       `json:"averageRounds"`
	AveragePlayerPoints        float64                       `json:"averagePlayerPoints"`
	TotalUniquePlayers         int                           `json:"totalUniquePlayers"`
	MostPlayedFaction          string                        `json:"mostPlayedFaction"`
	MostVictoriousFaction      string                        `json:"mostVictoriousFaction"`
	AverageGameRounds          float64                       `json:"averageGameRounds"`
}
type FactionPlayWinStat struct {
	PlayedCount int     `json:"playedCount"`
	WinCount    int     `json:"winCount"`
	PlayRate    float64 `json:"playRate"`
	WinRate     float64 `json:"winRate"`
}
type PlayerPointStdev struct {
	Player string  `json:"player"`
	Stdev  float64 `json:"stdev"`
}

type SecretObjectiveRate struct {
	Player          string  `json:"player"`
	SecretAppeared  int     `json:"secretAppeared"`
	SecretScored    int     `json:"secretScored"`
	SecretScoreRate float64 `json:"secretScoreRate"`
}

type PlayerAveragePoints struct {
	Player        string  `json:"player"`
	GamesPlayed   int     `json:"gamesPlayed"`
	TotalPoints   float64 `json:"totalPoints"`
	AveragePoints float64 `json:"averagePoints"`
}
type PlayerMostCommonFinish struct {
	Player     string `json:"player"`
	Position   int    `json:"position"`
	Count      int    `json:"count"`
	TotalGames int    `json:"totalGames"`
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
type PlayerFactionStats struct {
	Player   string         `json:"player"`
	Factions map[string]int `json:"factions"`
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
	// Total plays and wins across all factions
	totalFactionPlays := 0
	totalFactionWins := 0
	for _, count := range factionPlays {
		totalFactionPlays += count
	}
	for _, count := range factionWins {
		totalFactionWins += count
	}

	factionPlayWinDistribution := make(map[string]FactionPlayWinStat)
	allFactions := make(map[string]bool)

	// Ensure union of all faction keys
	for f := range factionPlays {
		allFactions[f] = true
	}
	for f := range factionWins {
		allFactions[f] = true
	}

	for faction := range allFactions {
		played := factionPlays[faction]
		won := factionWins[faction]

		var playRate, winRate float64
		if totalFactionPlays > 0 {
			playRate = float64(played) / float64(totalFactionPlays) * 100
		}
		if totalFactionWins > 0 {
			winRate = float64(won) / float64(totalFactionWins) * 100
		}

		factionPlayWinDistribution[faction] = FactionPlayWinStat{
			PlayedCount: played,
			WinCount:    won,
			PlayRate:    playRate,
			WinRate:     winRate,
		}
	}
	var avgPointsRows []struct {
		Name        string
		GamesPlayed int
		TotalPoints float64
	}

	err = database.DB.Raw(`
    SELECT p.name,
           COUNT(DISTINCT gp.game_id) AS games_played,
           COALESCE(SUM(s.points), 0) AS total_points
    FROM game_players gp
    JOIN players p ON p.id = gp.player_id
    LEFT JOIN scores s ON s.player_id = gp.player_id AND s.game_id = gp.game_id
    GROUP BY p.name
`).Scan(&avgPointsRows).Error
	if err != nil {
		return nil, err
	}

	var playerAverages []PlayerAveragePoints
	for _, row := range avgPointsRows {
		avg := 0.0
		if row.GamesPlayed > 0 {
			avg = row.TotalPoints / float64(row.GamesPlayed)
		}
		playerAverages = append(playerAverages, PlayerAveragePoints{
			Player:        row.Name,
			GamesPlayed:   row.GamesPlayed,
			TotalPoints:   row.TotalPoints,
			AveragePoints: avg,
		})
	}
	var factionRows []struct {
		Name    string
		Faction string
		Count   int
	}

	err = database.DB.Raw(`
    SELECT p.name, gp.faction, COUNT(*) as count
    FROM game_players gp
    JOIN players p ON p.id = gp.player_id
    GROUP BY p.name, gp.faction
`).Scan(&factionRows).Error
	if err != nil {
		return nil, err
	}

	playerFactionMap := make(map[string]map[string]int)
	for _, row := range factionRows {
		if _, exists := playerFactionMap[row.Name]; !exists {
			playerFactionMap[row.Name] = make(map[string]int)
		}
		playerFactionMap[row.Name][row.Faction] = row.Count
	}

	var topFactionsPerPlayer []PlayerFactionStats
	for player, factions := range playerFactionMap {
		topFactionsPerPlayer = append(topFactionsPerPlayer, PlayerFactionStats{
			Player:   player,
			Factions: factions,
		})
	}
	var playerPositions []PlayerMostCommonFinish

	// For each player, get their ranks across all games
	var positionData []struct {
		Player     string
		Position   int
		Count      int
		TotalGames int
	}

	err = database.DB.Raw(`
	WITH ranked_players AS (
		SELECT
			gp.game_id,
			p.name AS player,
			gp.player_id,
			SUM(
				CASE 
					WHEN s.type = 'public' OR s.type = 'secret' OR s.type = 'agenda' OR s.type = 'mecatol' OR s.type = 'support' OR s.type = 'imperial'
					THEN 1 ELSE 0
				END
			) AS score
		FROM game_players gp
		JOIN players p ON gp.player_id = p.id
		LEFT JOIN scores s ON gp.player_id = s.player_id AND gp.game_id = s.game_id
		GROUP BY gp.game_id, gp.player_id, p.name
	),
	ranked_with_position AS (
		SELECT
			player,
			game_id,
			DENSE_RANK() OVER (PARTITION BY game_id ORDER BY score DESC) AS position
		FROM ranked_players
	)
	SELECT
		player,
		position,
		COUNT(*) as count,
		(SELECT COUNT(*) FROM ranked_with_position WHERE player = r.player) as total_games
	FROM ranked_with_position r
	GROUP BY player, position
	ORDER BY player, count DESC
`).Scan(&positionData).Error
	if err != nil {
		return nil, err
	}

	// Build final slice with highest-frequency position per player
	seen := make(map[string]bool)
	for _, row := range positionData {
		if !seen[row.Player] {
			playerPositions = append(playerPositions, PlayerMostCommonFinish{
				Player:     row.Player,
				Position:   row.Position,
				Count:      row.Count,
				TotalGames: row.TotalGames,
			})
			seen[row.Player] = true
		}
	}
	// Step 7: Secret Objective Scoring Rate per Player (assume 3 per game)
	var secretRows []struct {
		Name         string
		GamesPlayed  int
		SecretScored int
	}

	err = database.DB.
		Raw(`
		SELECT p.name, 
		       COUNT(DISTINCT gp.game_id) AS games_played,
		       COUNT(DISTINCT s.id) AS secret_scored
		FROM players p
		LEFT JOIN game_players gp ON p.id = gp.player_id
		LEFT JOIN scores s ON s.player_id = gp.player_id AND s.game_id = gp.game_id AND s.type = 'secret'
		GROUP BY p.name
	`).Scan(&secretRows).Error
	if err != nil {
		return nil, err
	}

	var secretObjectiveRates []SecretObjectiveRate
	for _, row := range secretRows {
		// Each game gives 3 potential secret scoring opportunities
		secretAppeared := row.GamesPlayed * 3

		rate := 0.0
		if secretAppeared > 0 {
			rate = float64(row.SecretScored) / float64(secretAppeared) * 100
		}

		secretObjectiveRates = append(secretObjectiveRates, SecretObjectiveRate{
			Player:          row.Name,
			SecretAppeared:  secretAppeared,
			SecretScored:    row.SecretScored,
			SecretScoreRate: rate,
		})
	}
	// Step X: Standard deviation of points per player
	var pointRows []struct {
		Name  string
		Game  int
		Total float64
	}

	err = database.DB.Raw(`
	SELECT p.name, gp.game_id, COALESCE(SUM(s.points), 0) as total
	FROM game_players gp
	JOIN players p ON gp.player_id = p.id
	LEFT JOIN scores s ON s.player_id = gp.player_id AND s.game_id = gp.game_id
	GROUP BY p.name, gp.game_id
`).Scan(&pointRows).Error
	if err != nil {
		return nil, err
	}

	// Organize scores by player
	playerPoints := make(map[string][]float64)
	for _, row := range pointRows {
		playerPoints[row.Name] = append(playerPoints[row.Name], row.Total)
	}

	// Calculate standard deviation
	var playerPointStdevs []PlayerPointStdev
	for player, scores := range playerPoints {
		n := float64(len(scores))
		if n == 0 {
			playerPointStdevs = append(playerPointStdevs, PlayerPointStdev{Player: player, Stdev: 0})
			continue
		}
		var sum, mean, variance float64
		for _, s := range scores {
			sum += s
		}
		mean = sum / n
		for _, s := range scores {
			variance += (s - mean) * (s - mean)
		}
		stdev := math.Sqrt(variance / n)
		playerPointStdevs = append(playerPointStdevs, PlayerPointStdev{Player: player, Stdev: stdev})
	}
	var avgRounds float64
	err = database.DB.Raw(`
  SELECT AVG(round_count) 
  FROM (
    SELECT game_id, MAX(number) AS round_count 
    FROM rounds 
    GROUP BY game_id
  ) AS game_rounds
`).Scan(&avgRounds).Error
	if err != nil {
		return nil, err
	}
	var avgPlayerPoints float64
	err = database.DB.Raw(`
  SELECT AVG(points) FROM scores
`).Scan(&avgPlayerPoints).Error
	if err != nil {
		return nil, err
	}
	var uniquePlayers int64
	err = database.DB.Model(&models.Player{}).Count(&uniquePlayers).Error
	if err != nil {
		return nil, err
	}
	mostPlayedFaction := ""
	mostPlayedCount := 0
	for faction, count := range factionPlays {
		if count > mostPlayedCount {
			mostPlayedFaction = faction
			mostPlayedCount = count
		}
	}

	mostVictoriousFaction := ""
	mostVictoriousCount := 0
	for faction, count := range factionWins {
		if count > mostVictoriousCount {
			mostVictoriousFaction = faction
			mostVictoriousCount = count
		}
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
		ObjectiveFrequency:         freqMap,
		PlayerWinRates:             playerWinRates,
		ObjectiveAppearanceStats:   objectiveAppearanceStats,
		FactionPlayWinDistribution: factionPlayWinDistribution,
		PlayerAveragePoints:        playerAverages,
		TopFactionsPerPlayer:       topFactionsPerPlayer,
		PlayerMostCommonFinishes:   playerPositions,
		SecretObjectiveRates:       secretObjectiveRates,
		PlayerPointStdevs:          playerPointStdevs,
		AverageRounds:              avgRounds,
		AveragePlayerPoints:        avgPlayerPoints,
		TotalUniquePlayers:         int(uniquePlayers),
		MostPlayedFaction:          mostPlayedFaction,
		MostVictoriousFaction:      mostVictoriousFaction,
		AverageGameRounds:          avgRounds,
	}, nil
}
