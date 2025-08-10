// achievements/compute.go
package achievements

import (
	"time"

	"github.com/arphillips06/TI4-stats/models"
	"gorm.io/gorm"
)

func ComputeGameAchievements(db *gorm.DB, gameID uint) ([]Badge, error) {
	out := make([]Badge, 0)

	var st struct {
		Partial    bool
		FinishedAt *time.Time
	}
	if err := db.Model(&models.Game{}).
		Select("partial, finished_at").
		Where("id = ?", gameID).
		Scan(&st).Error; err != nil {
		return out, err
	}
	if st.Partial || st.FinishedAt == nil {
		return out, nil
	}

	var roundCount int64
	if err := db.Model(&models.Round{}).
		Where("game_id = ?", gameID).
		Count(&roundCount).Error; err != nil {
		return out, err
	}
	if roundCount == 0 {
		var cr struct{ CurrentRound int }
		_ = db.Model(&models.Game{}).Select("current_round").
			Where("id = ?", gameID).Scan(&cr).Error
		if cr.CurrentRound > 0 {
			roundCount = int64(cr.CurrentRound)
		}
		if roundCount == 0 {
			var rc struct{ C int64 }
			_ = db.Model(&models.Score{}).
				Select("COUNT(DISTINCT round_id) AS c").
				Where("game_id = ?", gameID).
				Scan(&rc).Error
			roundCount = rc.C
		}
	}
	if roundCount > 0 {
		roundsPerGame := db.Model(&models.Round{}).
			Select("rounds.game_id, COUNT(*) AS cnt").
			Joins("JOIN games g ON g.id = rounds.game_id").
			Where("g.partial = ? AND g.finished_at IS NOT NULL", false).
			Group("rounds.game_id")

		var rec struct{ Value *int }
		if err := db.Table("(?) r", roundsPerGame).
			Select("MIN(cnt) AS value").Scan(&rec).Error; err != nil {
			return out, err
		}

		current := int(roundCount)
		status := ""
		if rec.Value == nil || current < *rec.Value {
			status = "new"
		} else if rec.Value != nil && current == *rec.Value {
			status = "tied"
		}

		holders := []Holder{}
		var w struct{ WinnerID *uint }
		_ = db.Model(&models.Game{}).
			Select("winner_id").Where("id = ?", gameID).Scan(&w).Error
		if w.WinnerID != nil {
			gid := gameID
			holders = append(holders, Holder{PlayerID: *w.WinnerID, GameID: &gid})
		} else {
			var rows []struct{ PlayerID uint }
			_ = db.Model(&models.GamePlayer{}).
				Select("player_id").Where("game_id = ? AND won = ?", gameID, true).
				Scan(&rows).Error
			for _, r := range rows {
				gid := gameID
				holders = append(holders, Holder{PlayerID: r.PlayerID, GameID: &gid})
			}
		}

		out = append(out, Badge{
			Key:     "fastest_win",
			Label:   "Fastest Win",
			Value:   current,
			Status:  status,
			Holders: holders,
		})
	}

	type Row struct {
		PlayerID, RoundID uint
		Total             int
	}
	var gameBest []Row
	if err := db.Model(&models.Score{}).
		Select("scores.player_id, scores.round_id, SUM(scores.points) AS total").
		Joins("JOIN games ON games.id = scores.game_id").
		Where("scores.game_id = ? AND games.partial = FALSE", gameID).
		Group("scores.player_id, scores.round_id").
		Having("SUM(scores.points) IS NOT NULL").
		Order("total DESC").
		Scan(&gameBest).Error; err != nil {
		return out, err
	}
	if len(gameBest) > 0 {
		currentMax := gameBest[0].Total

		perRoundTotals := db.Model(&models.Score{}).
			Select("scores.game_id, scores.player_id, scores.round_id, SUM(scores.points) AS total").
			Joins("JOIN games ON games.id = scores.game_id").
			Where("games.partial = FALSE AND games.finished_at IS NOT NULL").
			Group("scores.game_id, scores.player_id, scores.round_id")

		var rec2 struct{ Value *int }
		if err := db.Table("(?) t", perRoundTotals).
			Select("MAX(total) AS value").Scan(&rec2).Error; err != nil {
			return out, err
		}

		status := ""
		if rec2.Value == nil || currentMax > *rec2.Value {
			status = "new"
		} else if rec2.Value != nil && currentMax == *rec2.Value {
			status = "tied"
		}

		holders := []Holder{}
		for _, r := range gameBest {
			if r.Total != currentMax {
				break
			}
			gid, rid := gameID, r.RoundID
			holders = append(holders, Holder{PlayerID: r.PlayerID, GameID: &gid, RoundID: &rid})
		}

		out = append(out, Badge{
			Key:     "most_points_in_round",
			Label:   "Most Points In A Round",
			Value:   currentMax,
			Status:  status,
			Holders: holders,
		})
	}

	return out, nil
}
