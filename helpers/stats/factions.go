package stats

import (
	"sort"

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
	var results []models.FactionAggregateStats
	db := database.DB

	// Step 1: Get raw totals
	err := db.Raw(`
		SELECT gp.faction AS faction,
		       COUNT(DISTINCT s.game_id) AS total_plays,
		       SUM(s.points) AS total_points_scored
		FROM scores s
		JOIN game_players gp ON s.player_id = gp.player_id AND s.game_id = gp.game_id
		GROUP BY gp.faction
	`).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Step 2: Count how many times each faction won
	var wins []struct {
		Faction string
		Wins    int
	}
	err = db.Raw(`
		SELECT gp.faction AS faction,
		       COUNT(*) AS wins
		FROM games g
		JOIN game_players gp ON g.winner_id = gp.player_id AND g.id = gp.game_id
		GROUP BY gp.faction
	`).Scan(&wins).Error
	if err != nil {
		return nil, err
	}
	winMap := map[string]int{}
	for _, w := range wins {
		winMap[w.Faction] = w.Wins
	}

	// Step 3: Compute VP histogram
	type HistogramRow struct {
		Faction string
		VP      int
		Count   int
	}
	var histoRows []HistogramRow
	err = db.Raw(`
		SELECT gp.faction,
			final_scores.vp AS vp,
			COUNT(*) AS count
		FROM (
			SELECT s.game_id, s.player_id, SUM(s.points) AS vp
			FROM scores s
			JOIN rounds r ON s.round_id = r.id
			WHERE r.number = (
				SELECT MAX(number) FROM rounds WHERE game_id = s.game_id
			)
			GROUP BY s.game_id, s.player_id
		) AS final_scores
		JOIN game_players gp ON gp.player_id = final_scores.player_id AND gp.game_id = final_scores.game_id
		GROUP BY gp.faction, final_scores.vp
	`).Scan(&histoRows).Error
	if err != nil {
		return nil, err
	}

	// Group histogram data by faction
	intermediate := map[string]map[int]int{} // faction -> vp -> count

	for _, row := range histoRows {
		if _, ok := intermediate[row.Faction]; !ok {
			intermediate[row.Faction] = map[int]int{}
		}
		intermediate[row.Faction][row.VP] += row.Count
	}

	histograms := map[string][]models.VPBucket{}
	for faction, vpMap := range intermediate {
		var buckets []models.VPBucket
		for vp, count := range vpMap {
			buckets = append(buckets, models.VPBucket{VP: vp, Count: count})
		}
		sort.Slice(buckets, func(i, j int) bool {
			return buckets[i].VP < buckets[j].VP
		})
		histograms[faction] = buckets
	}
	for i := range results {
		results[i].WonCount = winMap[results[i].Faction]

		if buckets, ok := histograms[results[i].Faction]; ok {
			results[i].VPHistogram = buckets
		} else {
			results[i].VPHistogram = []models.VPBucket{}
		}
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
