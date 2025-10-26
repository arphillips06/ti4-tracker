package services

import (
	"fmt"
	"log"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/helpers/stats"
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
		Preload("Speaker").
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
	if game.CurrentRound == 0 {
		log.Printf("WARNING: Game %d has CurrentRound = 0", game.ID)
	}

	var vpSummary *models.VictoryPathSummary
	if game.WinnerID != nil {
		vp, err := stats.CalculateVictoryPath(game.ID, *game.WinnerID)
		if err == nil {
			key := stats.FormatVictoryPathKey(vp)
			if _, found := CachedVictoryPathCounts[key]; !found {
				log.Printf("[VictoryPath] New key '%s' not found in cache. Refreshing cache.", key)
				RefreshVictoryPathCache()
			}

			freq := CachedVictoryPathCounts[key]
			uniqueness := 100
			if freq > 1 {
				uniqueness = int(100.0 / float64(freq))
			}

			vpSummary = &models.VictoryPathSummary{
				Path:       vp,
				Frequency:  freq,
				Uniqueness: uniqueness,
			}
			for k := range CachedVictoryPathCounts {
				fmt.Println("  ", k)
			}
			log.Printf("[VictoryPath] Cache pointer address: %p", &CachedVictoryPathCounts)

		}
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
	var allScoreDTOs []models.ScoreDTO

	for _, s := range scores {
		allScoreDTOs = append(allScoreDTOs, models.ScoreDTO{
			ID:               s.ID,
			GameID:           s.GameID,
			RoundID:          s.RoundID,
			PlayerID:         s.PlayerID,
			ObjectiveID:      s.ObjectiveID,
			Points:           s.Points,
			Type:             s.Type,
			AgendaTitle:      s.AgendaTitle,
			RelicTitle:       s.RelicTitle,
			OriginallySecret: s.OriginallySecret,
			CreatedAt:        s.CreatedAt,
		})
	}
	scoreDTOsByObjective := make(map[uint][]models.ScoreDTO)
	for _, s := range allScoreDTOs {
		if s.ObjectiveID != 0 {
			scoreDTOsByObjective[s.ObjectiveID] = append(scoreDTOsByObjective[s.ObjectiveID], s)
		}
	}
	var speakerID *uint
	speakerName := ""

	if game.CurrentRound != 0 {
		var currentRound models.Round
		err := database.DB.
			Where("game_id = ? AND number = ?", game.ID, game.CurrentRound).
			First(&currentRound).Error
		if err == nil {
			var assignment models.SpeakerAssignment
			err := database.DB.
				Where("game_id = ? AND round_id = ?", game.ID, currentRound.ID).
				First(&assignment).Error
			if err == nil {
				speakerID = &assignment.PlayerID
				var gp models.GamePlayer
				if err := database.DB.Preload("Player").First(&gp, assignment.PlayerID).Error; err == nil {
					speakerName = gp.Player.Name
				} else {
					log.Printf("failed to load GamePlayer for speaker: %v", err)
				}
			} else {
				log.Printf("no speaker assignment found for game %d, round %d: %v", game.ID, game.CurrentRound, err)
			}
		} else {
			log.Printf("could not fetch round object for game %d, round %d: %v", game.ID, game.CurrentRound, err)
		}
	}

	if speakerID == nil && game.SpeakerID != nil {
		speakerID = game.SpeakerID
		var gp models.GamePlayer
		if err := database.DB.Preload("Player").First(&gp, *speakerID).Error; err == nil {
			speakerName = gp.Player.Name
		} else {
			log.Printf("failed to load GamePlayer for speaker: %v", err)
		}
	}
	var all []models.SpeakerAssignment
	database.DB.Find(&all)

	return models.GameDetailResponse{
		ID:                 game.ID,
		GameNumber:         game.GameNumber,
		WinningPoints:      game.WinningPoints,
		CurrentRound:       game.CurrentRound,
		FinishedAt:         game.FinishedAt,
		UseObjectiveDecks:  game.UseObjectiveDecks,
		Players:            game.GamePlayers,
		Rounds:             game.Rounds,
		Objectives:         game.GameObjectives,
		Scores:             summaryList,
		AllScores:          allScoreDTOs,
		ScoresByObjective:  scoreDTOsByObjective,
		Winner:             &game.Winner,
		CustodiansPlayerID: custodiansPlayerID,
		WinnerVictoryPath:  vpSummary,
		SpeakerID:          speakerID,
		SpeakerName:        speakerName,
	}, nil
}
