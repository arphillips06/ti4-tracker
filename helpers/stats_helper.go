package helpers

import (
	"database/sql"
	"errors"
	"math"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

func CountTotalGames() (int64, error) {
	var count int64
	err := database.DB.Model(&models.Game{}).Count(&count).Error
	return count, err
}

func CalculateFactionStats() (map[string]int, map[string]int, map[string]float64, map[string]models.FactionPlayWinStat, error) {
	var factionPlays, factionWins []struct {
		Faction string
		Count   int
	}

	db := database.DB
	plays := make(map[string]int)
	wins := make(map[string]int)

	if err := db.Table("game_players").
		Select("faction, COUNT(*) as count").
		Group("faction").
		Scan(&factionPlays).Error; err != nil {
		return nil, nil, nil, nil, err
	}
	for _, f := range factionPlays {
		plays[f.Faction] = f.Count
	}

	if err := db.
		Table("game_players AS gp").
		Select("gp.faction, COUNT(*) AS count").
		Joins("JOIN games g ON g.id = gp.game_id AND g.winner_id = gp.player_id").
		Group("gp.faction").
		Scan(&factionWins).Error; err != nil {
		return nil, nil, nil, nil, err
	}

	for _, f := range factionWins {
		wins[f.Faction] = f.Count
	}

	winRates := make(map[string]float64)
	totalPlays, totalWins := 0, 0
	for f, count := range plays {
		totalPlays += count
		winCount := wins[f]
		if count > 0 {
			winRates[f] = float64(winCount) / float64(count) * 100
		} else {
			winRates[f] = 0
		}
	}
	for _, count := range wins {
		totalWins += count
	}

	distribution := make(map[string]models.FactionPlayWinStat)
	for f := range plays {
		distribution[f] = models.FactionPlayWinStat{
			PlayedCount: plays[f],
			WinCount:    wins[f],
			PlayRate:    float64(plays[f]) / float64(totalPlays) * 100,
			WinRate:     float64(wins[f]) / float64(totalWins) * 100,
		}
	}

	return plays, wins, winRates, distribution, nil
}

func CalculatePlayerWinRates() ([]models.PlayerWinRate, error) {
	var rows []struct {
		Name        string
		GamesPlayed int
		GamesWon    int
	}

	err := database.DB.
		Table("game_players AS gp").
		Select(`
		p.name,
		COUNT(DISTINCT gp.game_id) AS games_played,
		COUNT(DISTINCT CASE WHEN g.winner_id = gp.player_id THEN gp.game_id END) AS games_won`).
		Joins("JOIN players p ON p.id = gp.player_id").
		Joins("JOIN games g ON g.id = gp.game_id").
		Group("p.name").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	var rates []models.PlayerWinRate
	for _, r := range rows {
		rate := 0.0
		if r.GamesPlayed > 0 {
			rate = float64(r.GamesWon) / float64(r.GamesPlayed) * 100
		}
		rates = append(rates, models.PlayerWinRate{
			Player:      r.Name,
			GamesPlayed: r.GamesPlayed,
			GamesWon:    r.GamesWon,
			WinRate:     rate,
		})
	}

	return rates, nil
}

func CalculateObjectiveCounts() (map[string]int, error) {
	result := make(map[string]int)

	var publicCount, secretCount, stage1Count, stage2Count int64

	if err := database.DB.Model(&models.Score{}).Where("type = ?", "public").Count(&publicCount).Error; err != nil {
		return nil, err
	}
	result["publicScored"] = int(publicCount)

	if err := database.DB.Model(&models.Score{}).Where("type = ?", "secret").Count(&secretCount).Error; err != nil {
		return nil, err
	}
	result["secretScored"] = int(secretCount)

	if err := database.DB.Model(&models.Score{}).
		Joins("JOIN objectives ON scores.objective_id = objectives.id").
		Where("objectives.stage = ?", "I").
		Count(&stage1Count).Error; err != nil {
		return nil, err
	}
	result["stage1Scored"] = int(stage1Count)

	if err := database.DB.Model(&models.Score{}).
		Joins("JOIN objectives ON scores.objective_id = objectives.id").
		Where("objectives.stage = ?", "II").
		Count(&stage2Count).Error; err != nil {
		return nil, err
	}
	result["stage2Scored"] = int(stage2Count)

	return result, nil
}

func CalculateStandardDeviation(points []float64, mean float64) float64 {
	if len(points) == 0 {
		return 0
	}
	var sumSquares float64
	for _, x := range points {
		sumSquares += (x - mean) * (x - mean)
	}
	return math.Sqrt(sumSquares / float64(len(points)))
}

func CalculateObjectiveFrequencies() (map[string]int, error) {
	var rows []struct {
		Name  string
		Count int
	}

	err := database.DB.Table("scores").
		Select("objectives.name, COUNT(*) as count").
		Joins("JOIN objectives ON scores.objective_id = objectives.id").
		Group("objectives.name").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]int)
	for _, row := range rows {
		result[row.Name] = row.Count
	}
	return result, nil
}

func CalculatePlayerAverages() ([]models.PlayerAveragePoints, error) {
	var rows []struct {
		Name        string
		GamesPlayed int
		TotalPoints float64
	}

	err := database.DB.
		Table("game_players AS gp").
		Select(`
		p.name,
		COUNT(DISTINCT gp.game_id) AS games_played,
		COALESCE(SUM(s.points), 0) AS total_points`).
		Joins("JOIN players p ON p.id = gp.player_id").
		Joins("LEFT JOIN scores s ON s.player_id = gp.player_id AND s.game_id = gp.game_id").
		Group("p.name").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	var result []models.PlayerAveragePoints
	for _, row := range rows {
		avg := 0.0
		if row.GamesPlayed > 0 {
			avg = row.TotalPoints / float64(row.GamesPlayed)
		}
		result = append(result, models.PlayerAveragePoints{
			Player:        row.Name,
			GamesPlayed:   row.GamesPlayed,
			TotalPoints:   row.TotalPoints,
			AveragePoints: avg,
		})
	}
	return result, nil
}

func CalculateTopFactionsPerPlayer() ([]models.PlayerFactionStats, error) {
	var rows []struct {
		Name    string
		Faction string
		Count   int
	}

	err := database.DB.
		Table("game_players AS gp").
		Select("p.name, gp.faction, COUNT(*) as count").
		Joins("JOIN players p ON p.id = gp.player_id").
		Group("p.name, gp.faction").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	mapping := make(map[string]map[string]int)
	for _, row := range rows {
		if _, ok := mapping[row.Name]; !ok {
			mapping[row.Name] = make(map[string]int)
		}
		mapping[row.Name][row.Faction] = row.Count
	}

	var result []models.PlayerFactionStats
	for player, factions := range mapping {
		result = append(result, models.PlayerFactionStats{
			Player:   player,
			Factions: factions,
		})
	}
	return result, nil
}

func CountUniquePlayers() (int, error) {
	var count int64
	err := database.DB.Model(&models.Player{}).Count(&count).Error
	return int(count), err
}

func DetermineMostPlayedAndVictoriousFactions(plays, wins map[string]int) (string, string) {
	mostPlayed := ""
	mostVictorious := ""
	maxPlayed := 0
	maxWins := 0

	for faction, count := range plays {
		if count > maxPlayed {
			mostPlayed = faction
			maxPlayed = count
		}
	}
	for faction, count := range wins {
		if count > maxWins {
			mostVictorious = faction
			maxWins = count
		}
	}
	return mostPlayed, mostVictorious
}

func CalculateAverageRounds() (float64, error) {
	var avg sql.NullFloat64

	subQuery := database.DB.
		Table("rounds").
		Select("game_id, MAX(number) as round_count").
		Group("game_id")

	err := database.DB.
		Table("(?) as game_rounds", subQuery).
		Select("AVG(round_count)").
		Joins("JOIN games ON games.id = game_rounds.game_id").
		Where("games.partial = false").
		Scan(&avg).Error

	if err != nil {
		return 0, err
	}
	if !avg.Valid {
		// No valid data to average from
		return 0, nil
	}

	return avg.Float64, nil
}

func CalculateAveragePlayerPoints() (float64, error) {
	var avg float64
	subQuery := database.DB.
		Model(&models.Score{}).
		Select("SUM(points) as total").
		Group("game_id, player_id")

	err := database.DB.
		Table("(?) as sub", subQuery).
		Select("AVG(total)").
		Scan(&avg).Error

	return avg, err
}

func CalculateObjectiveAppearanceStats(totalGames int64) (map[string]models.ObjectiveStats, error) {
	if totalGames == 0 {
		return nil, errors.New("cannot calculate stats with 0 total games")
	}

	var appearances []struct {
		Name      string
		GameCount int
	}
	var scored []struct {
		Name      string
		GameCount int
	}
	err := database.DB.
		Table("game_objectives").
		Select("objectives.name, COUNT(DISTINCT game_objectives.game_id) as game_count").
		Joins("JOIN objectives ON game_objectives.objective_id = objectives.id").
		Joins("JOIN games ON game_objectives.game_id = games.id").
		Where("games.partial = false").
		Group("objectives.name").
		Scan(&appearances).Error
	if err != nil {
		return nil, err
	}
	err = database.DB.
		Model(&models.Score{}).
		Select("objectives.name, COUNT(DISTINCT scores.game_id) as game_count").
		Joins("JOIN objectives ON scores.objective_id = objectives.id").
		Joins("JOIN games ON scores.game_id = games.id").
		Where("games.partial = false").
		Group("objectives.name").
		Scan(&scored).Error
	if err != nil {
		return nil, err
	}
	appearanceMap := make(map[string]int)
	scoredMap := make(map[string]int)
	for _, row := range appearances {
		appearanceMap[row.Name] = row.GameCount
	}
	for _, row := range scored {
		scoredMap[row.Name] = row.GameCount
	}

	result := make(map[string]models.ObjectiveStats)
	allNames := make(map[string]bool)
	for name := range appearanceMap {
		allNames[name] = true
	}
	for name := range scoredMap {
		allNames[name] = true
	}

	for name := range allNames {
		appeared := appearanceMap[name]
		scored := scoredMap[name]

		appearanceRate := float64(appeared) / float64(totalGames) * 100
		scoreRate := 0.0
		if appeared > 0 {
			scoreRate = float64(scored) / float64(appeared) * 100
		}

		result[name] = models.ObjectiveStats{
			AppearanceRate:         appearanceRate,
			ScoredWhenAppearedRate: scoreRate,
		}
	}
	return result, nil
}
func CalculateMostCommonFinishes() ([]models.PlayerMostCommonFinish, error) {
	var positionData []struct {
		Player     string
		Position   int
		Count      int
		TotalGames int
	}

	err := database.DB.Raw(`
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

	seen := make(map[string]bool)
	var results []models.PlayerMostCommonFinish
	for _, row := range positionData {
		if !seen[row.Player] {
			results = append(results, models.PlayerMostCommonFinish{
				Player:     row.Player,
				Position:   row.Position,
				Count:      row.Count,
				TotalGames: row.TotalGames,
			})
			seen[row.Player] = true
		}
	}
	return results, nil
}

func CalculateSecretObjectiveRates() ([]models.SecretObjectiveRate, error) {
	type Result struct {
		Name         string
		GamesPlayed  int64
		SecretScored int64
	}

	var rows []Result

	// Subquery: games played per player
	subGamesPlayed := database.DB.
		Table("game_players").
		Joins("JOIN games ON games.id = game_players.game_id").
		Where("games.partial = false").
		Select("player_id, COUNT(DISTINCT game_id) AS games_played").
		Group("player_id")

		// Subquery: secret scored per player
	subSecrets := database.DB.
		Table("scores").
		Joins("JOIN games ON games.id = scores.game_id").
		Where("games.partial = false AND type = ?", "secret").
		Select("player_id, COUNT(DISTINCT scores.id) AS secret_scored").
		Group("player_id")

	// Join both subqueries on player_id
	err := database.DB.
		Table("players AS p").
		Select("p.name, COALESCE(gp.games_played, 0) AS games_played, COALESCE(ss.secret_scored, 0) AS secret_scored").
		Joins("LEFT JOIN (?) AS gp ON p.id = gp.player_id", subGamesPlayed).
		Joins("LEFT JOIN (?) AS ss ON p.id = ss.player_id", subSecrets).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	// Compute final score rates
	var result []models.SecretObjectiveRate
	for _, row := range rows {
		appeared := int(row.GamesPlayed) * 3 // Max allowed draws per player
		rate := 0.0
		if appeared > 0 {
			rate = float64(row.SecretScored) / float64(appeared) * 100
		}
		result = append(result, models.SecretObjectiveRate{
			Player:          row.Name,
			SecretAppeared:  int(appeared),
			SecretScored:    int(row.SecretScored),
			SecretScoreRate: rate,
		})
	}

	return result, nil
}

func CalculatePointStandardDeviations() ([]models.PlayerPointStdev, error) {
	var rows []struct {
		Name  string
		Game  int
		Total float64
	}

	err := database.DB.
		Table("game_players AS gp").
		Select("p.name, gp.game_id AS game, COALESCE(SUM(s.points), 0) AS total").
		Joins("JOIN players p ON p.id = gp.player_id").
		Joins("LEFT JOIN scores s ON s.player_id = gp.player_id AND s.game_id = gp.game_id").
		Group("p.name, gp.game_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	playerPoints := make(map[string][]float64)
	for _, row := range rows {
		playerPoints[row.Name] = append(playerPoints[row.Name], row.Total)
	}

	var result []models.PlayerPointStdev
	for player, points := range playerPoints {
		n := float64(len(points))
		if n == 0 {
			result = append(result, models.PlayerPointStdev{Player: player, Stdev: 0})
			continue
		}
		var sum float64
		for _, p := range points {
			sum += p
		}
		mean := sum / n
		var variance float64
		for _, p := range points {
			variance += (p - mean) * (p - mean)
		}
		result = append(result, models.PlayerPointStdev{
			Player: player,
			Stdev:  math.Sqrt(variance / n),
		})
	}
	return result, nil
}
