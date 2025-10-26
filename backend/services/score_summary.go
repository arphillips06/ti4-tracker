package services

import (
	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

func GetObjectiveScoreSummary(gameID uint) ([]models.ObjectiveScoreSummary, error) {
	var objectives []models.Objective
	var summaries []models.ObjectiveScoreSummary

	// Get all objectives scored in this game
	err := database.DB.
		Raw(`
            SELECT DISTINCT o.id, o.name, o.stage
            FROM scores s
            JOIN objectives o ON o.id = s.objective_id
            WHERE s.game_id = ?
        `, gameID).Scan(&objectives).Error
	if err != nil {
		return nil, err
	}
	for _, obj := range objectives {
		var playerNames []string

		err := database.DB.
			Table("scores").
			Select("players.name").
			Joins("JOIN players ON players.id = scores.player_id").
			Where("scores.game_id = ? AND scores.objective_id = ?", gameID, obj.ID).
			Pluck("players.name", &playerNames).Error

		if err != nil {
			return nil, err
		}

		summaries = append(summaries, models.ObjectiveScoreSummary{
			ObjectiveID: obj.ID,
			Name:        obj.Name,
			Stage:       obj.Stage,
			ScoredBy:    playerNames,
		})
	}

	return summaries, nil
}

// GetScoreSummaryByPlayer returns a summary of total points scored by each player
func GetScoreSummaryByPlayer(gameID string) ([]models.PlayerScoreSummary, error) {
	_, scores, err := GetGameAndScores(gameID)
	if err != nil {
		return nil, err
	}

	scoreSummaryMap := make(map[uint]models.PlayerScoreSummary)
	for _, s := range scores {
		summary := scoreSummaryMap[s.PlayerID]
		summary.PlayerID = s.PlayerID
		summary.PlayerName = s.Player.Name
		summary.Points += s.Points
		scoreSummaryMap[s.PlayerID] = summary
	}

	var summaryList []models.PlayerScoreSummary
	for _, s := range scoreSummaryMap {
		summaryList = append(summaryList, s)
	}
	return summaryList, nil
}

// GetScoresGroupedByRound returns scores grouped by round for a specific game
func GetScoresGroupedByRound(gameID string) ([]map[string]any, error) {
	type rawScore struct {
		Round  int    `json:"round"`
		Player string `json:"player"`
		Source string `json:"source"`
		Points int    `json:"points"`
	}

	var results []rawScore

	err := database.DB.
		Table("scores").
		Select(`
			COALESCE(rounds.number, 0) AS round,
			players.name AS player,
			COALESCE(objectives.name, scores.agenda_title, scores.relic_title,
				CASE
					WHEN scores.type = 'imperial' THEN 'Imperial Point'
					WHEN scores.type = 'mecatol' THEN 'Custodians'
					WHEN scores.type = 'Support' THEN 'Support for the Throne'
					ELSE 'Unknown'
				END
			) AS source,
			scores.points
		`).
		Joins("JOIN players ON scores.player_id = players.id").
		Joins("LEFT JOIN rounds ON scores.round_id = rounds.id").
		Joins("LEFT JOIN objectives ON scores.objective_id = objectives.id").
		Where("scores.game_id = ?", gameID).
		Order("round, players.name").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	grouped := make(map[int][]rawScore)
	for _, r := range results {
		grouped[r.Round] = append(grouped[r.Round], r)
	}

	var response []map[string]any
	for round, scores := range grouped {
		response = append(response, map[string]any{
			"round":  round,
			"scores": scores,
		})
	}

	return response, nil
}
