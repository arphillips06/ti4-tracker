package helpers

import (
	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

func InjectCDLObjectives(gameID uint, existing []models.GameObjective, scores []models.Score) []models.GameObjective {
	cdlObjectiveIDs := make(map[uint]bool)
	for _, score := range scores {
		if score.AgendaTitle == models.AgendaCDL {
			cdlObjectiveIDs[score.ObjectiveID] = true
		}
	}

	existingIDs := make(map[uint]bool)
	for _, obj := range existing {
		existingIDs[obj.ObjectiveID] = true
	}

	for objID := range cdlObjectiveIDs {
		if !existingIDs[objID] {
			var objective models.Objective
			_ = database.DB.First(&objective, objID)

			existing = append(existing, models.GameObjective{
				GameID:      gameID,
				ObjectiveID: objID,
				Objective:   objective,
				IsCDL:       true,
			})
		}
	}

	for i := range existing {
		if cdlObjectiveIDs[existing[i].ObjectiveID] {
			existing[i].IsCDL = true
		}
	}

	return existing
}
