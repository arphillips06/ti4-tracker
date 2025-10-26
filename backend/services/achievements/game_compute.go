package achievements

import (
	achievements_helper "github.com/arphillips06/TI4-stats/helpers/achievements"
	"github.com/arphillips06/TI4-stats/models"
	"gorm.io/gorm"
)

func ComputeGameAchievements(db *gorm.DB, gameID uint) ([]Badge, error) {
	ok, err := achievements_helper.IsFinishedNonPartial(db, gameID)
	if err != nil || !ok {
		return []Badge{}, err
	}

	out := make([]Badge, 0, 2)

	if b, yes, err := computeFastestWinBadge(db, gameID); err != nil {
		return nil, err
	} else if yes {
		out = append(out, b)
	}

	if b, yes, err := computeMostPointsInRoundBadge(db, gameID); err != nil {
		return nil, err
	} else if yes {
		out = append(out, b)
	}
	if b, yes, err := computeLargestWinMarginBadge(db, gameID); err != nil {
		return nil, err
	} else if yes {
		out = append(out, b)
	}

	return out, nil
}

func computeFastestWinBadge(db *gorm.DB, gameID uint) (Badge, bool, error) {
	rounds, err := achievements_helper.GetRoundCountForGame(db, gameID)
	if err != nil || rounds == 0 {
		return Badge{}, false, err
	}

	minRounds, err := getAllTimeMinRounds(db)
	if err != nil {
		return Badge{}, false, err
	}
	status := achievements_helper.CompareMinRecord(rounds, minRounds)

	holders, err := achievements_helper.GetWinnerHolders(db, gameID)
	if err != nil {
		return Badge{}, false, err
	}

	return Badge{
		Key:     "fastest_win",
		Label:   "Fastest Win",
		Value:   rounds,
		Status:  status,
		Holders: holders,
	}, true, nil
}

func computeMostPointsInRoundBadge(db *gorm.DB, gameID uint) (Badge, bool, error) {
	rows, err := getGameBestRoundTotals(db, gameID)
	if err != nil || len(rows) == 0 {
		return Badge{}, false, err
	}
	currentMax := rows[0].Total

	recordMax, err := getAllTimeMaxRoundPoints(db)
	if err != nil {
		return Badge{}, false, err
	}
	status := achievements_helper.CompareMaxRecord(currentMax, recordMax)

	holders := make([]achievements_helper.Holder, 0, 2)
	for _, r := range rows {
		if r.Total != currentMax {
			break
		}
		gid, rid := gameID, r.RoundID
		holders = append(holders, achievements_helper.Holder{
			PlayerID: r.PlayerID, GameID: &gid, RoundID: &rid,
		})
	}

	return Badge{
		Key:     "most_points_in_round",
		Label:   "Most Points In A Round",
		Value:   currentMax,
		Status:  status,
		Holders: holders,
	}, true, nil
}

func getAllTimeMinRounds(db *gorm.DB) (*int, error) {
	roundsPerGame := db.Model(&models.Round{}).
		Select("rounds.game_id, COUNT(*) AS cnt").
		Joins("JOIN games g ON g.id = rounds.game_id").
		Where("g.partial = ? AND g.finished_at IS NOT NULL", false).
		Group("rounds.game_id")

	var rec intVal
	if err := db.Table("(?) r", roundsPerGame).
		Select("MIN(cnt) AS value").
		Scan(&rec).Error; err != nil {
		return nil, err
	}
	return rec.Value, nil
}

func getGameBestRoundTotals(db *gorm.DB, gameID uint) ([]roundTotal, error) {
	var out []roundTotal
	if err := db.Model(&models.Score{}).
		Select("scores.player_id, r.number AS round_id, SUM(scores.points) AS total").
		Joins("JOIN games ON games.id = scores.game_id").
		Joins("JOIN rounds r ON r.id = scores.round_id").
		Where("scores.game_id = ? AND games.partial = FALSE", gameID).
		Group("scores.player_id, r.number").
		Having("SUM(scores.points) IS NOT NULL").
		Order("total DESC").
		Scan(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func getAllTimeMaxRoundPoints(db *gorm.DB) (*int, error) {
	perRoundTotals := db.Model(&models.Score{}).
		Select("scores.game_id, scores.player_id, r.number AS round_id, SUM(scores.points) AS total").
		Joins("JOIN games ON games.id = scores.game_id").
		Joins("JOIN rounds r ON r.id = scores.round_id").
		Where("games.partial = FALSE AND games.finished_at IS NOT NULL").
		Group("scores.game_id, scores.player_id, r.number")

	var rec intVal
	if err := db.Table("(?) t", perRoundTotals).
		Select("MAX(total) AS value").
		Scan(&rec).Error; err != nil {
		return nil, err
	}
	return rec.Value, nil
}

func computeLargestWinMarginBadge(db *gorm.DB, gameID uint) (Badge, bool, error) {
	rows, err := getGameFinalTotals(db, gameID)
	if err != nil || len(rows) == 0 {
		return Badge{}, false, err
	}
	currentMargin := 0
	if len(rows) >= 2 {
		currentMargin = rows[0].Total - rows[1].Total
	} else {
		currentMargin = rows[0].Total
	}

	recordMax, err := getAllTimeMaxWinningMargin(db)
	if err != nil {
		return Badge{}, false, err
	}
	status := achievements_helper.CompareMaxRecord(currentMargin, recordMax)

	holders, err := achievements_helper.GetWinnerHolders(db, gameID)
	if err != nil {
		return Badge{}, false, err
	}

	return Badge{
		Key:     "largest_win_margin",
		Label:   "Largest Win Margin",
		Value:   currentMargin,
		Status:  status,
		Holders: holders,
	}, true, nil
}

type playerTotal struct {
	PlayerID uint
	Total    int
}

func getGameFinalTotals(db *gorm.DB, gameID uint) ([]playerTotal, error) {
	var out []playerTotal
	if err := db.Model(&models.Score{}).
		Select("scores.player_id, SUM(scores.points) AS total").
		Joins("JOIN games ON games.id = scores.game_id").
		Where("scores.game_id = ? AND games.partial = FALSE", gameID).
		Group("scores.player_id").
		Having("SUM(scores.points) IS NOT NULL").
		Order("total DESC").
		Scan(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func getAllTimeMaxWinningMargin(db *gorm.DB) (*int, error) {
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
		return nil, err
	}

	var maxMargin *int
	var currentGame uint
	var first, second *int
	bump := func() {
		if first == nil {
			return
		}
		margin := *first
		if second != nil {
			margin = *first - *second
		}
		if maxMargin == nil || margin > *maxMargin {
			maxMargin = &margin
		}
		first, second = nil, nil
	}

	for _, r := range rows {
		if r.GameID != currentGame {
			if currentGame != 0 {
				bump()
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
		bump()
	}
	return maxMargin, nil
}
