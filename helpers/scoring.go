package helpers

import (
	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

func AggregatePlayerScores(scores []models.Score) []models.PlayerScoreSummary {
	summaryMap := make(map[uint]models.PlayerScoreSummary)

	for _, s := range scores {
		summary := summaryMap[s.PlayerID]
		summary.PlayerID = s.PlayerID
		summary.PlayerName = s.Player.Name
		summary.Points += s.Points
		summaryMap[s.PlayerID] = summary
	}

	var summaries []models.PlayerScoreSummary
	for _, s := range summaryMap {
		summaries = append(summaries, s)
	}

	return summaries
}

func CreateRelicScore(gameID, playerID uint, points int, relicTitle string) error {
	score := models.Score{
		GameID:     gameID,
		PlayerID:   playerID,
		Points:     points,
		Type:       "relic",
		RelicTitle: relicTitle,
	}
	return CreateGenericScore(score)
}

func CreateBasicScore(gameID, roundID, playerID uint, points int, scoreType string) error {
	score := models.Score{
		GameID:   gameID,
		RoundID:  roundID,
		PlayerID: playerID,
		Points:   points,
		Type:     scoreType,
	}
	return CreateGenericScore(score)
}

func GetPlayerTotalPoints(gameID, playerID uint) (int, error) {
	var total int
	err := database.DB.Model(&models.Score{}).
		Where("game_id = ? AND player_id = ?", gameID, playerID).
		Select("SUM(points)").Scan(&total).Error
	return total, err
}

func CreateGenericScore(score models.Score) error {
	return database.DB.Create(&score).Error
}

func CreateAgendaScore(gameID, roundID, playerID, points int, agendaTitle string, objectiveID uint) error {
	score := models.Score{
		GameID:      uint(gameID),
		RoundID:     uint(roundID),
		PlayerID:    uint(playerID),
		Points:      points,
		Type:        "agenda",
		AgendaTitle: agendaTitle,
		ObjectiveID: objectiveID,
	}
	return CreateGenericScore(score)
}

func GetPlayerScoresMap(gameID uint) (map[uint]int, error) {
	// Get all players in the game
	var gamePlayers []models.GamePlayer
	if err := database.DB.Where("game_id = ?", gameID).Find(&gamePlayers).Error; err != nil {
		return nil, err
	}

	// Initialize all player totals to 0
	playerTotals := make(map[uint]int)
	for _, gp := range gamePlayers {
		playerTotals[gp.PlayerID] = 0
	}

	// Get all scores for the game
	var scores []models.Score
	if err := database.DB.Where("game_id = ?", gameID).Find(&scores).Error; err != nil {
		return nil, err
	}

	for _, s := range scores {
		playerTotals[s.PlayerID] += s.Points
	}

	return playerTotals, nil
}

func CreateObjectiveScore(gameID, roundID, playerID, objectiveID uint, points int) error {
	var obj models.Objective
	if err := database.DB.First(&obj, objectiveID).Error; err != nil {
		return err
	}

	scoreType := "public"
	if obj.Stage == "Secret" {
		scoreType = "secret"
	}

	score := models.Score{
		GameID:      gameID,
		RoundID:     roundID,
		PlayerID:    playerID,
		ObjectiveID: objectiveID,
		Points:      points,
		Type:        scoreType,
	}
	return database.DB.Create(&score).Error
}
