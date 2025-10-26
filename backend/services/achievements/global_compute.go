package achievements

import (
	"fmt"
	"math"
	"sort"
	"time"

	ah "github.com/arphillips06/TI4-stats/helpers/achievements"
	"github.com/arphillips06/TI4-stats/models"
	"gorm.io/gorm"
)

var RecordComebackKid = models.Achievement{
	Key:  "record_comeback_kid",
	Name: "Record: Comeback Kid",
	Type: "record",
}

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
	if b, ok, err := globalComebackKid(db); err != nil {
		return nil, err
	} else if ok {
		out = append(out, b)
	}
	if b, ok, err := globalCurrentWinningStreak(db); err != nil {
		return nil, err
	} else if ok {
		out = append(out, b)
	}
	if b, ok, err := globalLongestWinningStreak(db); err != nil {
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
		Select("scores.game_id, scores.player_id, r.number AS round_id, SUM(scores.points) AS total").
		Joins("JOIN games ON games.id = scores.game_id").
		Joins("JOIN rounds r ON r.id = scores.round_id").
		Where("games.partial = FALSE AND games.finished_at IS NOT NULL").
		Group("scores.game_id, scores.player_id, r.number")

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

func globalComebackKid(db *gorm.DB) (Badge, bool, error) {
	var games []models.Game
	if err := db.
		Where("partial = FALSE AND finished_at IS NOT NULL").
		Preload("GamePlayers.Player").
		Preload("Rounds.Scores").
		Find(&games).Error; err != nil {
		return Badge{}, false, err
	}

	bestDeficit := -1
	holders := make([]ah.Holder, 0, 4)

	for _, g := range games {
		sort.SliceStable(g.Rounds, func(i, j int) bool {
			if g.Rounds[i].Number != g.Rounds[j].Number {
				return g.Rounds[i].Number < g.Rounds[j].Number
			}
			return g.Rounds[i].ID < g.Rounds[j].ID
		})

		if len(g.Rounds) == 0 || len(g.GamePlayers) < 2 {
			continue
		}
		finalTotals := map[uint]int{}
		for _, r := range g.Rounds {
			for _, s := range r.Scores {
				finalTotals[s.PlayerID] += s.Points
			}
		}
		if len(finalTotals) < 2 {
			continue
		}

		var winnerID uint
		maxFinal := math.MinInt
		for pid, vp := range finalTotals {
			if vp > maxFinal {
				maxFinal = vp
				winnerID = pid
			}
		}
		running := map[uint]int{}
		for pid := range finalTotals {
			running[pid] = 0
		}

		for i, r := range g.Rounds {
			for _, s := range r.Scores {
				running[s.PlayerID] += s.Points
			}
			if i == len(g.Rounds)-1 {
				break
			}
			minVP := math.MaxInt
			maxVP := math.MinInt
			for _, vp := range running {
				if vp < minVP {
					minVP = vp
				}
				if vp > maxVP {
					maxVP = vp
				}
			}

			winnerVP := running[winnerID]
			if winnerVP == minVP {
				deficit := maxVP - winnerVP

				switch {
				case deficit > bestDeficit:
					bestDeficit = deficit
					holders = holders[:0]
					gid := g.ID
					rnum := uint(r.Number)
					holders = append(holders, ah.Holder{
						PlayerID: winnerID,
						GameID:   &gid,
						RoundID:  &rnum,
					})
				case deficit == bestDeficit:
					gid := g.ID
					rnum := uint(r.Number)
					holders = append(holders, ah.Holder{
						PlayerID: winnerID,
						GameID:   &gid,
						RoundID:  &rnum,
					})
				}
			}
		}
	}

	if bestDeficit < 0 {
		return Badge{}, false, nil
	}

	return Badge{
		Key:     "record_comeback_kid",
		Label:   "Biggest Comeback",
		Value:   bestDeficit,
		Status:  "record",
		Holders: holders,
	}, true, nil
}

func globalCurrentWinningStreak(db *gorm.DB) (Badge, bool, error) {
	type lastPlay struct {
		PlayerID   uint
		GameID     uint
		FinishedAt time.Time
		Won        bool
	}
	var lastPlays []lastPlay
	err := db.Table("game_players gp").
		Select("gp.player_id, g.id as game_id, g.finished_at, gp.won").
		Joins("JOIN games g ON g.id = gp.game_id").
		Where("g.partial = FALSE AND g.finished_at IS NOT NULL").
		Order("g.finished_at DESC").
		Scan(&lastPlays).Error
	if err != nil {
		return Badge{}, false, err
	}
	streaks := make(map[uint]int)

	for _, lp := range lastPlays {
		if _, ok := streaks[lp.PlayerID]; ok {
			continue
		}
		if !lp.Won {
			streaks[lp.PlayerID] = 0
			continue
		}
		var wins []struct {
			Won bool
		}
		if err := db.Table("game_players gp").
			Select("gp.won").
			Joins("JOIN games g ON g.id = gp.game_id").
			Where("gp.player_id = ? AND g.partial = FALSE AND g.finished_at IS NOT NULL", lp.PlayerID).
			Order("g.finished_at DESC").
			Scan(&wins).Error; err != nil {
			return Badge{}, false, err
		}

		streak := 0
		for _, w := range wins {
			if w.Won {
				streak++
			} else {
				break
			}
		}
		streaks[lp.PlayerID] = streak
	}

	// Find the max streak
	maxStreak := -1
	holders := []ah.Holder{}
	for pid, s := range streaks {
		switch {
		case s > maxStreak:
			maxStreak = s
			holders = holders[:0]
			holders = append(holders, ah.Holder{PlayerID: pid})
		case s == maxStreak:
			holders = append(holders, ah.Holder{PlayerID: pid})
		}
	}

	if maxStreak <= 0 {
		return Badge{}, false, nil
	}

	return Badge{
		Key:     "current_winning_streak",
		Label:   "Current Winning Streak",
		Value:   maxStreak,
		Status:  "record",
		Holders: holders,
	}, true, nil
}

func globalLongestWinningStreak(db *gorm.DB) (Badge, bool, error) {
	type row struct {
		PlayerID   uint
		Won        bool
		FinishedAt time.Time
	}
	var rows []row
	if err := db.Table("game_players gp").
		Select("gp.player_id, gp.won, g.finished_at").
		Joins("JOIN games g ON g.id = gp.game_id").
		Where("g.partial = FALSE AND g.finished_at IS NOT NULL").
		Order("gp.player_id, g.finished_at").
		Scan(&rows).Error; err != nil {
		return Badge{}, false, err
	}

	streaks := make(map[uint]int)
	bestStreaks := make(map[uint]int)

	for _, r := range rows {
		if r.Won {
			streaks[r.PlayerID]++
			if streaks[r.PlayerID] > bestStreaks[r.PlayerID] {
				bestStreaks[r.PlayerID] = streaks[r.PlayerID]
			}
		} else {
			streaks[r.PlayerID] = 0
		}
	}

	maxStreak := -1
	holders := []ah.Holder{}
	for pid, s := range bestStreaks {
		switch {
		case s > maxStreak:
			maxStreak = s
			holders = holders[:0]
			holders = append(holders, ah.Holder{PlayerID: pid})
		case s == maxStreak:
			holders = append(holders, ah.Holder{PlayerID: pid})
		}
	}

	if maxStreak <= 0 {
		return Badge{}, false, nil
	}

	return Badge{
		Key:     "longest_winning_streak",
		Label:   "Longest Winning Streak",
		Value:   maxStreak,
		Status:  "record",
		Holders: holders,
	}, true, nil
}
