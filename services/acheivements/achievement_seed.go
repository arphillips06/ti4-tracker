package achievements

import (
	"time"

	"github.com/arphillips06/TI4-stats/models"
	"gorm.io/gorm"
)

func EnsureAchievementsSeeded(db *gorm.DB) error {
	// Migrate tables (idempotent)
	if err := db.AutoMigrate(&models.Achievement{}, &models.PlayerAchievement{}); err != nil {
		return err
	}

	// Seed core records (idempotent)
	rows := []models.Achievement{
		{Key: "fastest_win", Name: "Fastest Win"},
		{Key: "most_points_in_round", Name: "Most Points in a Round"},
		{Key: "most_custodians_taken", Name: "Most Custodians Taken"},
	}
	for _, a := range rows {
		if err := db.
			Where("key = ?", a.Key).
			FirstOrCreate(&models.Achievement{}, a).Error; err != nil {
			return err
		}
	}
	return nil
}

// (Already shown earlier, but included for completeness)
type AchievementBadge struct {
	PlayerID   uint   `json:"player_id"`
	PlayerName string `json:"player_name"`
	Key        string `json:"key"`
	Name       string `json:"name"`
	Value      *int   `json:"value"`
	GameID     *uint  `json:"game_id,omitempty"`
	RoundID    *uint  `json:"round_id,omitempty"`
}

func GetPlayerAchievements(db *gorm.DB, playerID uint) ([]AchievementBadge, error) {
	out := make([]AchievementBadge, 0)
	err := db.
		Table("player_achievements pa").
		Select(`
			pa.player_id,
			COALESCE(p.name, CAST(pa.player_id AS TEXT)) AS player_name,
			a.key, a.name,
			pa.numeric_value AS value,
			pa.game_id, pa.round_id
		`).
		Joins("JOIN achievements a ON a.id = pa.achievement_id").
		Joins("LEFT JOIN players p ON p.id = pa.player_id").
		Where("pa.player_id = ?", playerID).
		Order("a.key ASC, pa.awarded_at DESC").
		Scan(&out).Error
	return out, err
}

func GetGameAchievements(db *gorm.DB, gameID uint) ([]AchievementBadge, error) {
	out := make([]AchievementBadge, 0)
	err := db.
		Table("player_achievements pa").
		Select(`
			pa.player_id,
			COALESCE(p.name, CAST(pa.player_id AS TEXT)) AS player_name,
			a.key, a.name,
			pa.numeric_value AS value,
			pa.game_id, pa.round_id
		`).
		Joins("JOIN achievements a ON a.id = pa.achievement_id").
		Joins("LEFT JOIN players p ON p.id = pa.player_id").
		Where("pa.game_id = ?", gameID).
		Order("a.key ASC, pa.awarded_at DESC").
		Scan(&out).Error
	return out, err
}

// Optional admin/debug helper
func RecomputeAllFinishedGames(db *gorm.DB, eval func(*gorm.DB, uint) error) (int, error) {
	type Row struct{ ID uint }
	var games []Row
	if err := db.Raw(`SELECT id FROM games WHERE finished_at IS NOT NULL`).Scan(&games).Error; err != nil {
		return 0, err
	}
	for _, g := range games {
		_ = eval(db, g.ID) // log inside eval on error; don’t hard-fail
	}
	return len(games), nil
}

// For completeness if PlayerAchievement missing
type PlayerAchievement struct {
	ID            uint `gorm:"primaryKey"`
	PlayerID      uint `gorm:"index"`
	AchievementID uint `gorm:"index"`
	GameID        *uint
	RoundID       *uint
	NumericValue  *int
	AwardedAt     time.Time `gorm:"index"`
}
