package achievements

import (
	"fmt"

	ah "github.com/arphillips06/TI4-stats/helpers/achievements"
	"github.com/arphillips06/TI4-stats/models"
	"gorm.io/gorm"
)

func ComputeGlobalAchievements(db *gorm.DB) ([]Badge, error) {
	out := make([]Badge, 0, 2)

	if b, ok, err := globalFastestWin(db); err != nil {
		return nil, err
	} else if ok {
		out = append(out, b)
	}

	if b, ok, err := globalMostPointsInRound(db); err != nil {
		return nil, err
	} else if ok {
		out = append(out, b)
	}
	if b, ok, err := globalLargestWinMargin(db); err != nil {
		return nil, err
	} else if ok {
		out = append(out, b)
	}

	return out, nil
}

func globalFastestWin(db *gorm.DB) (Badge, bool, error) {
	roundsPerGame := db.Model(&models.Round{}).
		Select("rounds.game_id, COUNT(*) AS cnt").
		Joins("JOIN games g ON g.id = rounds.game_id").
		Where("g.partial = ? AND g.finished_at IS NOT NULL", false).
		Group("rounds.game_id")

	var min struct{ Value *int }
	if err := db.Table("(?) r", roundsPerGame).
		Select("MIN(cnt) AS value").
		Scan(&min).Error; err != nil {
		return Badge{}, false, err
	}
	if min.Value == nil {
		return Badge{}, false, nil
	}

	var gameIDs []uint
	if err := db.Table("(?) r", roundsPerGame).
		Where("r.cnt = ?", *min.Value).
		Pluck("r.game_id", &gameIDs).Error; err != nil {
		return Badge{}, false, err
	}
	if len(gameIDs) == 0 {
		return Badge{}, false, nil
	}

	holders := make([]ah.Holder, 0, len(gameIDs))
	for _, gid := range gameIDs {
		hs, err := ah.GetWinnerHolders(db, gid)
		if err != nil {
			return Badge{}, false, fmt.Errorf("getting winners for game %d: %w", gid, err)
		}
		holders = append(holders, hs...)
	}

	return Badge{
		Key:     "fastest_win",
		Label:   "Fastest Win",
		Value:   *min.Value,
		Status:  "record",
		Holders: holders,
	}, true, nil
}

func globalMostPointsInRound(db *gorm.DB) (Badge, bool, error) {
	perRoundTotals := db.Model(&models.Score{}).
		Select("scores.game_id, scores.player_id, scores.round_id, SUM(scores.points) AS total").
		Joins("JOIN games ON games.id = scores.game_id").
		Where("games.partial = FALSE AND games.finished_at IS NOT NULL").
		Group("scores.game_id, scores.player_id, scores.round_id")

	var max struct{ Value *int }
	if err := db.Table("(?) t", perRoundTotals).
		Select("MAX(total) AS value").
		Scan(&max).Error; err != nil {
		return Badge{}, false, err
	}
	if max.Value == nil {
		return Badge{}, false, nil
	}

	type row struct {
		GameID   uint
		PlayerID uint
		RoundID  uint
	}
	var rows []row
	if err := db.Table("(?) t", perRoundTotals).
		Select("t.game_id, t.player_id, t.round_id").
		Where("t.total = ?", *max.Value).
		Scan(&rows).Error; err != nil {
		return Badge{}, false, err
	}

	holders := make([]ah.Holder, 0, len(rows))
	for _, r := range rows {
		gid, rid := r.GameID, r.RoundID
		holders = append(holders, ah.Holder{
			PlayerID: r.PlayerID,
			GameID:   &gid,
			RoundID:  &rid,
		})
	}

	return Badge{
		Key:     "most_points_in_round",
		Label:   "Most Points In A Round",
		Value:   *max.Value,
		Status:  "record",
		Holders: holders,
	}, true, nil
}

func globalLargestWinMargin(db *gorm.DB) (Badge, bool, error) {
	type row struct {
		GameID   uint
		PlayerID uint
		Total    int
	}
	var rows []row
	if err := db.Model(&models.Score{}).
		Select("scores.game_id, scores.player_id, SUM(scores.points) AS total").
		Joins("JOIN games ON games.id = scores.game_id").
		Where("games.partial = FALSE AND games.finished_at IS NOT NULL").
		Group("scores.game_id, scores.player_id").
		Having("SUM(scores.points) IS NOT NULL").
		Order("scores.game_id ASC, total DESC").
		Scan(&rows).Error; err != nil {
		return Badge{}, false, err
	}
	if len(rows) == 0 {
		return Badge{}, false, nil
	}

	margins := make(map[uint]int)
	currentGame := uint(0)
	var first, second *int
	bump := func(gid uint) {
		if first == nil {
			return
		}
		margin := *first
		if second != nil {
			margin = *first - *second
		}
		margins[gid] = margin
		first, second = nil, nil
	}
	for _, r := range rows {
		if r.GameID != currentGame {
			if currentGame != 0 {
				bump(currentGame)
			}
			currentGame = r.GameID
		}
		if first == nil {
			v := r.Total
			first = &v
			continue
		}
		if second == nil {
			v := r.Total
			second = &v
		}
	}
	if currentGame != 0 {
		bump(currentGame)
	}

	var maxMargin *int
	for _, m := range margins {
		if maxMargin == nil || m > *maxMargin {
			v := m
			maxMargin = &v
		}
	}
	if maxMargin == nil {
		return Badge{}, false, nil
	}

	gameIDs := make([]uint, 0, 4)
	for gid, m := range margins {
		if m == *maxMargin {
			gameIDs = append(gameIDs, gid)
		}
	}

	holders := make([]ah.Holder, 0, len(gameIDs))
	for _, gid := range gameIDs {
		hs, err := ah.GetWinnerHolders(db, gid)
		if err != nil {
			return Badge{}, false, fmt.Errorf("getting winners for game %d: %w", gid, err)
		}
		holders = append(holders, hs...)
	}

	return Badge{
		Key:     "largest_win_margin",
		Label:   "Largest Win Margin",
		Value:   *maxMargin,
		Status:  "record",
		Holders: holders,
	}, true, nil
}
