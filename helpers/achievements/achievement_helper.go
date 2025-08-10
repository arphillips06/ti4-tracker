package achievements_helper

import (
	"time"

	"github.com/arphillips06/TI4-stats/models"
	"gorm.io/gorm"
)

type roundCountRow struct{ CurrentRound int }
type countRow struct{ C int64 }
type Holder struct {
	PlayerID uint  `json:"player_id"`
	GameID   *uint `json:"game_id,omitempty"`
	RoundID  *uint `json:"round_id,omitempty"`
}
type winnerRow struct{ WinnerID *uint }

const (
	StatusNew  = "new"
	StatusTied = "tied"
)

func IsFinishedNonPartial(db *gorm.DB, gameID uint) (bool, error) {
	var st struct {
		Partial    bool
		FinishedAt *time.Time
	}
	if err := db.Model(&models.Game{}).
		Select("partial, finished_at").
		Where("id = ?", gameID).
		Scan(&st).Error; err != nil {
		return false, err
	}
	return !st.Partial && st.FinishedAt != nil, nil
}

func GetRoundCountForGame(db *gorm.DB, gameID uint) (int, error) {
	var roundCount int64
	if err := db.Model(&models.Round{}).
		Where("game_id = ?", gameID).
		Count(&roundCount).Error; err != nil {
		return 0, err
	}
	if roundCount > 0 {
		return int(roundCount), nil
	}

	var cr roundCountRow
	_ = db.Model(&models.Game{}).
		Select("current_round").
		Where("id = ?", gameID).
		Scan(&cr).Error
	if cr.CurrentRound > 0 {
		return cr.CurrentRound, nil
	}

	var rc countRow
	_ = db.Model(&models.Score{}).
		Select("COUNT(DISTINCT round_id) AS c").
		Where("game_id = ?", gameID).
		Scan(&rc).Error

	return int(rc.C), nil
}

func GetWinnerHolders(db *gorm.DB, gameID uint) ([]Holder, error) {
	holders := make([]Holder, 0, 2)
	gid := gameID

	var w winnerRow
	if err := db.Model(&models.Game{}).
		Select("winner_id").
		Where("id = ?", gameID).
		Scan(&w).Error; err != nil {
		return nil, err
	}
	if w.WinnerID != nil {
		holders = append(holders, Holder{PlayerID: *w.WinnerID, GameID: &gid})
		return holders, nil
	}

	var rows []struct{ PlayerID uint }
	if err := db.Model(&models.GamePlayer{}).
		Select("player_id").
		Where("game_id = ? AND won = ?", gameID, true).
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	for _, r := range rows {
		holders = append(holders, Holder{PlayerID: r.PlayerID, GameID: &gid})
	}
	return holders, nil
}

func CompareMinRecord(current int, record *int) string {
	if record == nil || current < *record {
		return StatusNew
	}
	if current == *record {
		return StatusTied
	}
	return ""
}

func CompareMaxRecord(current int, record *int) string {
	if record == nil || current > *record {
		return StatusNew
	}
	if current == *record {
		return StatusTied
	}
	return ""
}
