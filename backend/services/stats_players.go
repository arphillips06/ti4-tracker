package services

import (
	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

type PlayerCustodiansStats struct {
	PlayerID                uint   `json:"player_id"`
	PlayerName              string `json:"player_name"`
	GamesPlayed             int    `json:"games_played"`
	GamesWon                int    `json:"games_won"`
	CustodiansTaken         int    `json:"custodians_taken"`
	CustodiansWins          int    `json:"custodians_wins"`
	CustodiansWinPercentage int    `json:"custodians_win_percentage"`
}

func GetPlayerCustodiansStats() ([]PlayerCustodiansStats, error) {
	var players []models.Player
	if err := database.DB.Find(&players).Error; err != nil {
		return nil, err
	}

	var results []PlayerCustodiansStats

	for _, player := range players {
		var gamesPlayed int64
		var gamesWon int64
		var custodiansTaken int64
		var custodiansWins int64

		// Games Played
		if err := database.DB.Model(&models.GamePlayer{}).
			Where("player_id = ?", player.ID).
			Count(&gamesPlayed).Error; err != nil {
			return nil, err
		}

		if err := database.DB.Model(&models.Game{}).
			Where("winner_id = ?", player.ID).
			Count(&gamesWon).Error; err != nil {
			return nil, err
		}

		if err := database.DB.Model(&models.Score{}).
			Where("player_id = ? AND type = 'mecatol'", player.ID).
			Count(&custodiansTaken).Error; err != nil {
			return nil, err
		}

		if err := database.DB.Raw(`
			SELECT COUNT(DISTINCT s.game_id)
			FROM scores s
			JOIN games g ON g.id = s.game_id
			WHERE s.type = 'mecatol' AND s.player_id = ? AND g.winner_id = ?
		`, player.ID, player.ID).Scan(&custodiansWins).Error; err != nil {
			return nil, err
		}

		custodiansWinPercent := 0
		if gamesWon > 0 {
			custodiansWinPercent = int((float64(custodiansWins) / float64(gamesWon)) * 100)
		}

		results = append(results, PlayerCustodiansStats{
			PlayerID:                player.ID,
			PlayerName:              player.Name,
			GamesPlayed:             int(gamesPlayed),
			GamesWon:                int(gamesWon),
			CustodiansTaken:         int(custodiansTaken),
			CustodiansWins:          int(custodiansWins),
			CustodiansWinPercentage: custodiansWinPercent,
		})
	}

	return results, nil
}
