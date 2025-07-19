package stats

import (
	"database/sql"

	"github.com/arphillips06/TI4-stats/database"
)

func CalculateAverageRounds() (float64, error) {
	var avg sql.NullFloat64

	subQuery := database.DB.
		Table("rounds").
		Select("game_id, MAX(number) as round_count").
		Group("game_id")

	err := database.DB.
		Table("(?) as game_rounds", subQuery).
		Select("AVG(round_count)").
		Joins("JOIN games ON games.id = game_rounds.game_id").
		Where("games.partial = false").
		Scan(&avg).Error

	if err != nil {
		return 0, err
	}
	if !avg.Valid {
		// No valid data to average from
		return 0, nil
	}

	return avg.Float64, nil
}
