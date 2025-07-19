package stats

import (
	"fmt"
	"sort"
	"time"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

func CountTotalGames() (int64, error) {
	var count int64
	err := database.DB.Model(&models.Game{}).Count(&count).Error
	return count, err
}
func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %02dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}
func computeStats(games []models.Game) models.GameLengthCategoryStats {
	var durations []models.GameDurationStat
	var totalRoundSeconds int64
	var totalRounds int

	for _, game := range games {
		duration := game.FinishedAt.Sub(game.CreatedAt)

		stat := models.GameDurationStat{
			GameID:    game.ID,
			Duration:  formatDuration(duration),
			Seconds:   int64(duration.Seconds()),
			StartedAt: game.CreatedAt,
		}

		if !game.Partial {
			roundCount := len(game.Rounds)
			if roundCount == 0 {
				continue
			}
			stat.RoundCount = roundCount
			totalRoundSeconds += int64(duration.Seconds())
			totalRounds += roundCount
		}

		durations = append(durations, stat)
	}

	if len(durations) == 0 {
		return models.GameLengthCategoryStats{}
	}

	var totalGameSeconds int64
	for _, d := range durations {
		totalGameSeconds += d.Seconds
	}

	avgGame := formatDuration(time.Duration(totalGameSeconds/int64(len(durations))) * time.Second)

	var sortByRounds, sortByTime []models.GameDurationStat
	for _, d := range durations {
		if d.RoundCount > 0 {
			sortByRounds = append(sortByRounds, d)
		}
	}
	sortByTime = append([]models.GameDurationStat(nil), durations...)

	sort.Slice(sortByTime, func(i, j int) bool {
		return sortByTime[i].Seconds < sortByTime[j].Seconds
	})

	var shortestByRounds, longestByRounds models.GameDurationStat
	if len(sortByRounds) > 0 {
		sort.Slice(sortByRounds, func(i, j int) bool {
			return sortByRounds[i].RoundCount < sortByRounds[j].RoundCount
		})
		shortestByRounds = sortByRounds[0]
		longestByRounds = sortByRounds[len(sortByRounds)-1]
	}

	averageRound := ""
	if totalRounds > 0 {
		avgRound := totalRoundSeconds / int64(totalRounds)
		averageRound = formatDuration(time.Duration(avgRound) * time.Second)
	}

	return models.GameLengthCategoryStats{
		LongestByRounds:  longestByRounds,
		ShortestByRounds: shortestByRounds,
		LongestByTime:    sortByTime[len(sortByTime)-1],
		ShortestByTime:   sortByTime[0],
		AverageRoundTime: averageRound,
		AverageGameTime:  avgGame,
	}
}

func GetGameLengthStats() (models.GameLengthStats, error) {
	var games []models.Game
	db := database.DB

	err := db.Preload("Rounds").Preload("GamePlayers").Find(&games).Error
	if err != nil {
		return models.GameLengthStats{}, err
	}

	var (
		allGames         []models.Game
		threePlayerGames []models.Game
		fourPlayerGames  []models.Game
	)

	for _, game := range games {
		if game.FinishedAt == nil || game.CreatedAt.IsZero() {
			continue
		}
		count := len(game.GamePlayers)
		switch count {
		case 3:
			threePlayerGames = append(threePlayerGames, game)

		case 4:
			fourPlayerGames = append(fourPlayerGames, game)
		}
		allGames = append(allGames, game)
	}

	return models.GameLengthStats{
		All:         computeStats(allGames),
		ThreePlayer: computeStats(threePlayerGames),
		FourPlayer:  computeStats(fourPlayerGames),
	}, nil
}
