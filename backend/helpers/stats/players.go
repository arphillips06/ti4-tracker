package stats

import (
	"math"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

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

func CountUniquePlayers() (int, error) {
	var count int64
	err := database.DB.Model(&models.Player{}).Count(&count).Error
	return int(count), err
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
