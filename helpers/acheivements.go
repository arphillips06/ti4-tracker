package helpers

import (
	"context"
	"time"

	"github.com/arphillips06/TI4-stats/models"
	"gorm.io/gorm"
)

type RecordHolder struct {
	PlayerID uint
	GameID   *uint
	RoundID  *uint
	Value    int
}

func GetAchievementByKey(db *gorm.DB, key string) (models.Achievement, error) {
	var ach models.Achievement
	err := db.Where("key = ?", key).First(&ach).Error
	return ach, err
}

func UpsertRecordHolders(ctx context.Context, db *gorm.DB, key string, holders []RecordHolder, newRecord, equalRecord bool) error {
	ach, err := GetAchievementByKey(db, key)
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if newRecord {
			if err := tx.Where("achievement_id = ?", ach.ID).Delete(&models.PlayerAchievement{}).Error; err != nil {
				return err
			}
			for _, h := range holders {
				if err := tx.Create(&models.PlayerAchievement{
					PlayerID:      h.PlayerID,
					AchievementID: ach.ID,
					GameID:        h.GameID,
					RoundID:       h.RoundID,
					NumericValue:  &h.Value,
					AwardedAt:     time.Now(),
				}).Error; err != nil {
					return err
				}
			}
			return nil
		}

		if equalRecord {
			for _, h := range holders {
				var count int64
				if err := tx.Model(&models.PlayerAchievement{}).
					Where("player_id = ? AND achievement_id = ?", h.PlayerID, ach.ID).
					Count(&count).Error; err != nil {
					return err
				}
				if count == 0 {
					if err := tx.Create(&models.PlayerAchievement{
						PlayerID:      h.PlayerID,
						AchievementID: ach.ID,
						GameID:        h.GameID,
						RoundID:       h.RoundID,
						NumericValue:  &h.Value,
						AwardedAt:     time.Now(),
					}).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}

func ReplaceRecordHolders(ctx context.Context, db *gorm.DB, key string, holders []RecordHolder) error {
	ach, err := GetAchievementByKey(db, key)
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("achievement_id = ?", ach.ID).Delete(&models.PlayerAchievement{}).Error; err != nil {
			return err
		}
		for _, h := range holders {
			if err := tx.Create(&models.PlayerAchievement{
				PlayerID:      h.PlayerID,
				AchievementID: ach.ID,
				NumericValue:  &h.Value,
				AwardedAt:     time.Now(),
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
