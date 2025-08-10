package achievements

import (
	"context"

	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
	"gorm.io/gorm"
)

const (
	AchFastestWin        = "fastest_win"
	AchMostPointsInRound = "most_points_in_round"
	AchMostCustodians    = "most_custodians_taken"
)

func EvaluateAchievementsAfterGame(db *gorm.DB, gameID uint) error {
	ctx := context.Background()
	if err := evalFastestWin(ctx, db, gameID); err != nil {
		return err
	}
	if err := evalMostPointsInRound(ctx, db, gameID); err != nil {
		return err
	}
	if err := evalMostCustodiansTaken(ctx, db); err != nil {
		return err
	}
	return nil
}

// ---------- Fastest Win (fewest rounds among winners) ----------
func evalFastestWin(ctx context.Context, db *gorm.DB, gameID uint) error {
	// 1) How many rounds did THIS game have?
	var roundCount int64
	if err := db.Model(&models.Round{}).
		Where("game_id = ?", gameID).
		Count(&roundCount).Error; err != nil {
		return err
	}
	// Fallbacks if needed (defensive)
	if roundCount == 0 {
		var cr struct{ CurrentRound int }
		if err := db.Raw(`SELECT current_round FROM games WHERE id = ?`, gameID).Scan(&cr).Error; err == nil && cr.CurrentRound > 0 {
			roundCount = int64(cr.CurrentRound)
		}
	}
	if roundCount == 0 {
		var rc struct{ C int64 }
		_ = db.Raw(`SELECT COUNT(DISTINCT round_id) AS c FROM scores WHERE game_id = ?`, gameID).Scan(&rc).Error
		roundCount = rc.C
	}
	if roundCount == 0 {
		return nil // nothing to evaluate
	}

	// 2) Winner(s) of THIS game
	var gr struct{ WinnerID *uint }
	if err := db.Raw(`SELECT winner_id FROM games WHERE id = ?`, gameID).Scan(&gr).Error; err != nil {
		return err
	}

	holders := make([]helpers.RecordHolder, 0)
	if gr.WinnerID != nil {
		holders = append(holders, helpers.RecordHolder{
			PlayerID: *gr.WinnerID, GameID: &gameID, Value: int(roundCount),
		})
	} else {
		type WinRow struct{ PlayerID uint }
		var rows []WinRow
		if err := db.Raw(`SELECT player_id FROM game_players WHERE game_id = ? AND won = TRUE`, gameID).Scan(&rows).Error; err != nil {
			return err
		}
		for _, r := range rows {
			holders = append(holders, helpers.RecordHolder{
				PlayerID: r.PlayerID, GameID: &gameID, Value: int(roundCount),
			})
		}
	}
	if len(holders) == 0 {
		return nil // no winner -> skip
	}

	currentValue := int(roundCount)

	// 3) Global record = MIN rounds among FINISHED games only
	var rec struct{ Value *int }
	if err := db.Raw(`
		SELECT MIN(r.cnt) AS value
		FROM (
			SELECT r.game_id, COUNT(*) AS cnt
			FROM rounds r
			JOIN games g ON g.id = r.game_id
			WHERE g.finished_at IS NOT NULL
			GROUP BY r.game_id
		) r
	`).Scan(&rec).Error; err != nil {
		return err
	}

	newRecord := false
	equalRecord := false
	if rec.Value == nil {
		newRecord = true
	} else if currentValue < *rec.Value {
		newRecord = true
	} else if currentValue == *rec.Value {
		equalRecord = true
	}

	// 4) Upsert holders (replace on new, add on tie)
	return helpers.UpsertRecordHolders(ctx, db, AchFastestWin, holders, newRecord, equalRecord)
}

// ---------- Most Points in a Round (global max of per-round deltas) ----------
func evalMostPointsInRound(ctx context.Context, db *gorm.DB, gameID uint) error {
	type Row struct {
		PlayerID uint
		RoundID  uint
		Total    int
	}
	var gameBest []Row
	if err := db.Raw(`
		SELECT s.player_id, s.round_id, SUM(s.points) AS total
		FROM scores s
		WHERE s.game_id = ?
		GROUP BY s.player_id, s.round_id
		HAVING SUM(s.points) IS NOT NULL
		ORDER BY total DESC`, gameID).Scan(&gameBest).Error; err != nil {
		return err
	}
	if len(gameBest) == 0 {
		return nil
	}
	currentMax := gameBest[0].Total

	type Rec struct{ Value int }
	var rec Rec
	if err := db.Raw(`
		SELECT MAX(total) AS value FROM (
			SELECT game_id, player_id, round_id, SUM(points) AS total
			FROM scores
			GROUP BY game_id, player_id, round_id
		) t`).Scan(&rec).Error; err != nil {
		return err
	}

	newRecord := rec.Value == 0 || currentMax > rec.Value
	equalRecord := rec.Value != 0 && currentMax == rec.Value

	holders := []helpers.RecordHolder{}
	for _, r := range gameBest {
		if r.Total != currentMax {
			break
		}
		rid := r.RoundID
		holders = append(holders, helpers.RecordHolder{
			PlayerID: r.PlayerID,
			GameID:   &gameID,
			RoundID:  &rid,
			Value:    currentMax,
		})
	}
	return helpers.UpsertRecordHolders(ctx, db, AchMostPointsInRound, holders, newRecord, equalRecord)
}

// ---------- Most Custodians Taken (lifetime, type="mecatol") ----------
func evalMostCustodiansTaken(ctx context.Context, db *gorm.DB) error {
	type Row struct {
		PlayerID uint
		Cnt      int
	}
	var rows []Row
	if err := db.Raw(`
		SELECT s.player_id, COUNT(*) AS cnt
		FROM scores s
		WHERE s.type = 'mecatol' AND s.points > 0
		GROUP BY s.player_id
		ORDER BY cnt DESC`).Scan(&rows).Error; err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}
	best := rows[0].Cnt

	holders := []helpers.RecordHolder{}
	for _, r := range rows {
		if r.Cnt != best {
			break
		}
		holders = append(holders, helpers.RecordHolder{
			PlayerID: r.PlayerID,
			Value:    r.Cnt,
		})
	}
	return helpers.ReplaceRecordHolders(ctx, db, AchMostCustodians, holders)
}
