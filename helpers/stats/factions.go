package stats

import (
	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

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

func GetFactionPlayerStats() ([]models.FactionPlayerStats, error) {
	db := database.DB
	var results []models.FactionPlayerStats

	err := db.
		Model(&models.GamePlayer{}).
		Select("faction, players.name as player, COUNT(*) as played_count, SUM(CASE WHEN game_players.won THEN 1 ELSE 0 END) as won_count").
		Joins("JOIN players ON players.id = game_players.player_id").
		Group("faction, players.name").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func GetFactionAggregateStats() ([]models.FactionAggregateStats, error) {
	db := database.DB
	var results []models.FactionAggregateStats

	err := db.
		Model(&models.GamePlayer{}).
		Select(`
			faction,
			COUNT(DISTINCT game_players.game_id) AS total_plays,
			SUM(CASE WHEN game_players.won THEN 1 ELSE 0 END) AS won_count,
			COALESCE(SUM(scores.points), 0) AS total_points_scored
		`).
		Joins("LEFT JOIN scores ON scores.player_id = game_players.player_id AND scores.game_id = game_players.game_id").
		Group("game_players.faction").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
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
			PlayRate: func() float64 {
				if totalPlays == 0 {
					return 0
				}
				return float64(plays[f]) / float64(totalPlays) * 100
			}(),
			WinRate: func() float64 {
				if totalWins == 0 {
					return 0
				}
				return float64(wins[f]) / float64(totalWins) * 100
			}(),
		}
	}

	return plays, wins, winRates, distribution, nil
}
