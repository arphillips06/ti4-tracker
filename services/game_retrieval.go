package services

import (
	"fmt"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
)

// Gets a game by its string ID
func GetGameByID(gameID uint) (*models.Game, error) {
	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		return nil, err
	}
	return &game, nil
}

func GetGameAndScores(gameID string) (models.Game, []models.Score, error) {
	var game models.Game
	if err := database.DB.
		Preload("GamePlayers.Player").
		Preload("Rounds").
		Preload("Winner").
		Preload("GameObjectives.Objective").
		Preload("GameObjectives.Round").
		First(&game, gameID).Error; err != nil {
		return game, nil, fmt.Errorf("game not found")
	}

	var scores []models.Score
	if err := database.DB.
		Preload("Player").
		Preload("Objective").
		Where("game_id = ?", game.ID).
		Find(&scores).Error; err != nil {
		return game, nil, fmt.Errorf("could not load scores")
	}

	// Inject CDL Objectives if needed
	game.GameObjectives = helpers.InjectCDLObjectives(game.ID, game.GameObjectives, scores)

	return game, scores, nil
}

func BuildGameDetailResponse(gameID string) (models.GameDetailResponse, error) {
	game, scores, err := GetGameAndScores(gameID)
	if err != nil {
		return models.GameDetailResponse{}, err
	}

	var custodiansPlayerID *uint
	scoreSummaryMap := make(map[uint]models.PlayerScoreSummary)

	for _, s := range scores {
		if s.Type == "mecatol" {
			custodiansPlayerID = &s.PlayerID
		}
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
	scoresByObjective := make(map[uint][]models.Score)
	for _, s := range scores {
		if s.ObjectiveID != 0 {
			scoresByObjective[s.ObjectiveID] = append(scoresByObjective[s.ObjectiveID], s)
		}
	}

	return models.GameDetailResponse{
		ID:                 game.ID,
		WinningPoints:      game.WinningPoints,
		CurrentRound:       game.CurrentRound,
		FinishedAt:         game.FinishedAt,
		UseObjectiveDecks:  game.UseObjectiveDecks,
		Players:            game.GamePlayers,
		Rounds:             game.Rounds,
		Objectives:         game.GameObjectives,
		Scores:             summaryList,
		AllScores:          scores,
		ScoresByObjective:  scoresByObjective,
		Winner:             &game.Winner,
		CustodiansPlayerID: custodiansPlayerID,
	}, nil
}
