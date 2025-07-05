package database

import (
	"database/sql"
	"log"

	"github.com/arphillips06/TI4-stats/database/objectives"
	"github.com/arphillips06/TI4-stats/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "modernc.org/sqlite" // pure Go SQLite driver
)

var DB *gorm.DB

func InitDatabase() {
	// Open pure Go sqlite driver via database/sql
	sqlDB, err := sql.Open("sqlite", "ti4stats.db")
	if err != nil {
		log.Fatal("Failed to open database/sql DB:", err)
	}

	// Pass sql.DB to GORM sqlite dialector
	DB, err = gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database (gorm):", err)
	}

	// Automigrate your models
	err = DB.AutoMigrate(&models.Game{}, &models.Player{}, &models.Round{}, &models.Score{}, &models.GamePlayer{}, &models.GameObjective{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

func SeedObjectives() {
	//add "stage I" && "Stage II to objectives"
	for _, obj := range objectives.StageOne {
		obj.Stage = "I"
		insertObjective(obj)
	}
	for _, obj := range objectives.StageTwo {
		obj.Stage = "II"
		insertObjective(obj)
	}
	for _, obj := range objectives.Secret {
		obj.Stage = "Secret"
		insertObjective(obj)
	}
}

func insertObjective(obj models.Objective) {
	var existing models.Objective
	if err := DB.Where("name = ?", obj.Name).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := DB.Create(&obj).Error; err != nil {
				log.Printf("Failed to seed objective '%s': %v\n", obj.Name, err)
			} else {
				log.Printf("Seeded objective: %s\n", obj.Name)
			}
		} else {
			log.Printf("Error checking objective '%s': %v\n", obj.Name, err)
		}
	} else {
		// Update existing record with any missing fields
		existing.Type = obj.Type
		existing.Description = obj.Description
		existing.Points = obj.Points
		existing.Stage = obj.Stage
		existing.Phase = obj.Phase
		if err := DB.Save(&existing).Error; err != nil {
			log.Printf("Failed to update objective '%s': %v\n", obj.Name, err)
		}
	}
}

// allObjectives := append(append(objectives.StageOne, objectives.StageTwo...), objectives.Secret...)

// for _, obj := range allObjectives {
// 	var existing models.Objective
// 	if err := DB.Where("name = ?", obj.Name).First(&existing).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			if err := DB.Create(&obj).Error; err != nil {
// 				log.Printf("Failed to seed objective '%s': %v\n", obj.Name, err)
// 			} else {
// 				log.Printf("Seeded objective: %s\n", obj.Name)
// 			}
// 		} else {
// 			log.Printf("Error checking objective '%s': %v\n", obj.Name, err)
// 		}
// 	}
// }
//}
