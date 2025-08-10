package models

import "time"

type Achievement struct {
	ID        uint   `gorm:"primaryKey"`
	Key       string `gorm:"uniqueIndex;size:64"` // e.g. "fastest_win"
	Name      string // e.g. "Fastest Win"
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PlayerAchievement struct {
	ID            uint      `gorm:"primaryKey"`
	PlayerID      uint      `gorm:"index"`
	AchievementID uint      `gorm:"index"`
	GameID        *uint     // optional: which game established/earned it
	RoundID       *uint     // optional
	NumericValue  *int      // record value (e.g., 5 rounds, 4 points-in-round, 7 custodians)
	TextValue     *string   // optional extra
	AwardedAt     time.Time `gorm:"index"`
}

type AchievementBadge struct {
	PlayerID   uint   `json:"player_id"`
	PlayerName string `json:"player_name"`
	Key        string `json:"key"`   // e.g. "fastest_win"
	Name       string `json:"name"`  // e.g. "Fastest Win"
	Value      *int   `json:"value"` // record value (rounds, points, count)
	GameID     *uint  `json:"game_id,omitempty"`
	RoundID    *uint  `json:"round_id,omitempty"`
}
