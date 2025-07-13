package services

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
	"gorm.io/gorm"
)

func ScoreImperialPoint(gameID, roundID, playerID uint) error {
	roundID, err := helpers.GetCurrentRoundID(gameID)
	if err != nil {
		return err
	}

	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		return err
	}
	if game.FinishedAt != nil {
		return fmt.Errorf("game is already finished")
	}

	score := models.Score{
		GameID:   gameID,
		RoundID:  roundID,
		PlayerID: playerID,
		Points:   1,
		Type:     "imperial",
	}
	if err := database.DB.Create(&score).Error; err != nil {
		return err
	}
	return MaybeFinishGameFromScore(&game, playerID)
}
func ScoreAgendaPoint(gameID, playerID uint, points int, agendaTitle string) error {
	roundID, err := helpers.GetCurrentRoundID(gameID)
	if err != nil {
		return fmt.Errorf("failed to get current round: %w", err)
	}

	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		return fmt.Errorf("failed to load game: %w", err)
	}

	var totalPoints int64
	err = database.DB.
		Model(&models.Score{}).
		Where("game_id = ? AND player_id = ?", gameID, playerID).
		Select("SUM(points)").
		Row().
		Scan(&totalPoints)
	if err != nil {
		return fmt.Errorf("failed to calculate current score: %w", err)
	}

	if int(totalPoints)+points < 0 {
		return fmt.Errorf("agenda scoring would reduce points below zero")
	}

	score := models.Score{
		GameID:      gameID,
		RoundID:     roundID,
		PlayerID:    playerID,
		Points:      points,
		Type:        "agenda",
		AgendaTitle: agendaTitle,
	}
	if err := database.DB.Create(&score).Error; err != nil {
		return err
	}

	return MaybeFinishGameFromScore(&game, playerID)
}

func ScoreMecatolPoint(gameID, playerID uint) error {
	roundID, err := helpers.GetCurrentRoundID(gameID)
	if err != nil {
		return err
	}

	var existing models.Score
	err = database.DB.
		Where("game_id = ? AND type = ?", gameID, "mecatol").
		First(&existing).Error
	if err == nil {
		return fmt.Errorf("Mecatol Rex point already awarded")
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		return err
	}

	score := models.Score{
		GameID:   gameID,
		RoundID:  roundID,
		PlayerID: playerID,
		Points:   1,
		Type:     "mecatol",
	}
	if err := database.DB.Create(&score).Error; err != nil {
		return err
	}

	return MaybeFinishGameFromScore(&game, playerID)
}

func ScoreSupportPoint(gameID, playerID uint) error {
	roundID, err := helpers.GetCurrentRoundID(gameID)
	if err != nil {
		return err
	}

	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		return err
	}
	if game.FinishedAt != nil {
		return fmt.Errorf("game is already finished")
	}

	var playerCount int64
	err = database.DB.
		Model(&models.GamePlayer{}).
		Where("game_id = ?", gameID).
		Count(&playerCount).Error
	if err != nil {
		return err
	}

	var supportCount int64
	err = database.DB.
		Model(&models.Score{}).
		Where("game_id = ? AND type = ?", gameID, "Support").
		Select("COALESCE(SUM(points), 0)").
		Row().
		Scan(&supportCount)

	if err != nil {
		return err
	}

	if supportCount >= playerCount-1 {
		return fmt.Errorf("Support for the Throne can only be scored by %d players in a %d-player game", playerCount-1, playerCount)
	}

	score := models.Score{
		GameID:   gameID,
		RoundID:  roundID,
		PlayerID: playerID,
		Points:   1,
		Type:     "Support",
	}
	if err := database.DB.Create(&score).Error; err != nil {
		return err
	}

	return MaybeFinishGameFromScore(&game, playerID)
}

func LoseOneSupportPoint(gameID, playerID uint) error {
	roundID, err := helpers.GetCurrentRoundID(gameID)
	if err != nil {
		return err
	}

	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		return err
	}
	if game.FinishedAt != nil {
		return fmt.Errorf("game is already finished")
	}

	// Create a negative support score record
	score := models.Score{
		GameID:   gameID,
		RoundID:  roundID,
		PlayerID: playerID,
		Points:   -1,
		Type:     "Support",
	}

	if err := database.DB.Create(&score).Error; err != nil {
		return err
	}

	return nil
}

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

func ValidateSecretScoringRules(gameID, playerID, roundID, objectiveID uint) error {
	var objective models.Objective
	if err := database.DB.First(&objective, objectiveID).Error; err != nil {
		log.Printf("[ERROR] Could not find objective %d: %v", objectiveID, err)
		return errors.New("Objective not found")
	}

	if strings.ToLower(objective.Type) != "secret" {
		return nil
	}

	// Phase-specific limit check
	var countThisPhase int64
	if err := database.DB.
		Model(&models.Score{}).
		Where(`
			player_id = ? AND 
			round_id = ? AND 
			LOWER(type) = 'secret' AND 
			objective_id IN (
				SELECT id FROM objectives WHERE LOWER(phase) = ?
			)`,
			playerID, roundID, strings.ToLower(objective.Phase)).
		Count(&countThisPhase).Error; err != nil {
		return errors.New("Failed to validate secret scoring rules")
	}

	if countThisPhase > 0 {
		return errors.New("Player has already scored a secret objective in this phase this round")
	}

	// Total secret scoring cap
	var totalSecrets int64
	if err := database.DB.
		Model(&models.Score{}).
		Joins("JOIN objectives ON objectives.id = scores.objective_id").
		Where("scores.player_id = ? AND scores.game_id = ? AND LOWER(scores.type) = 'secret'", playerID, gameID).
		Count(&totalSecrets).Error; err != nil {
		return errors.New("Failed to count total secret objectives")
	}

	// Obsidian check
	var obsidianUsed int64
	if err := database.DB.
		Model(&models.Score{}).
		Where("game_id = ? AND player_id = ? AND LOWER(type) = 'relic' AND LOWER(relic_title) = 'the obsidian'", gameID, playerID).
		Count(&obsidianUsed).Error; err != nil {
		return errors.New("Failed to check Obsidian use")
	}

	maxSecrets := int64(3)
	if obsidianUsed > 0 {
		maxSecrets = 4
	}

	if totalSecrets >= maxSecrets {
		return fmt.Errorf("Player has already scored the maximum of %d secret objectives", maxSecrets)
	}

	log.Printf("[DEBUG] Player has scored %d/%d secret objectives", totalSecrets, maxSecrets)
	return nil
}

func AddScoreToGame(gameID, playerID uint, objectiveName string) (*models.Score, int, error) {
	var game models.Game
	if err := database.DB.Preload("Rounds").First(&game, gameID).Error; err != nil {
		return nil, 0, errors.New("game not found")
	}

	if game.FinishedAt != nil {
		return nil, 0, errors.New("game is already finished")
	}

	var obj models.Objective
	if err := database.DB.Where("LOWER(name) = ?", strings.ToLower(objectiveName)).First(&obj).Error; err != nil {
		return nil, 0, errors.New("objective not found")
	}

	var round models.Round
	if err := database.DB.Where("game_id = ? AND number = ?", game.ID, game.CurrentRound).First(&round).Error; err != nil {
		return nil, 0, errors.New("current round not found")
	}

	if obj.Type == "Secret" {
		if err := ValidateSecretScoringRules(gameID, playerID, round.ID, obj.ID); err != nil {
			return nil, 0, err
		}
	}

	exists, err := CheckIfScoreExists(game.ID, playerID, obj.ID)
	if err != nil {
		return nil, 0, err
	}
	if exists {
		return nil, 0, errors.New("objective already scored")
	}

	score := models.Score{
		GameID:      game.ID,
		PlayerID:    playerID,
		ObjectiveID: obj.ID,
		Points:      obj.Points,
		RoundID:     round.ID,
	}

	if err := database.DB.Create(&score).Error; err != nil {
		return nil, 0, err
	}

	var total int
	database.DB.Model(&models.Score{}).
		Where("game_id = ? AND player_id = ?", game.ID, playerID).
		Select("SUM(points)").Scan(&total)

	if total >= game.WinningPoints {
		if err := MaybeFinishGameFromScore(&game, playerID); err != nil {
			return &score, total, err
		}
	}

	return &score, total, nil
}

func CheckIfScoreExists(gameID, playerID, objectiveID uint) (bool, error) {
	var existing models.Score
	err := database.DB.
		Where("game_id = ? AND player_id = ? AND objective_id = ?", gameID, playerID, objectiveID).
		First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check for existing score: %w", err)
	}

	return true, nil
}

func WinnerByScore(game *models.Game) error {
	var topScore struct {
		PlayerID uint
		Points   int
	}

	err := database.DB.
		Table("scores").
		Select("player_id, SUM(points) as points").
		Where("game_id = ?", game.ID).
		Group("player_id").
		Order("points DESC").
		Limit(1).
		Scan(&topScore).Error

	if err != nil {
		return err
	}

	if topScore.PlayerID != 0 {
		game.WinnerID = &topScore.PlayerID
		return nil
	}

	return nil // No winner yet
}

func MaybeFinishGameFromExhaustion(game *models.Game) error {
	now := time.Now()
	game.FinishedAt = &now

	if err := WinnerByScore(game); err != nil {
		return err
	}

	return database.DB.Save(game).Error
}

func GetGameAndRounds(gameID uint) (*models.Game, error) {
	var game models.Game
	if err := database.DB.Preload("Rounds").First(&game, gameID).Error; err != nil {
		return nil, err
	}
	return &game, nil
}

func MaybeFinishGameFromScore(game *models.Game, scoringPlayerID uint) error {
	var totalPoints int
	err := database.DB.Model(&models.Score{}).
		Where("game_id = ? AND player_id = ?", game.ID, scoringPlayerID).
		Select("SUM(points)").Scan(&totalPoints).Error
	if err != nil {
		return err
	}

	if totalPoints >= game.WinningPoints {
		now := time.Now()
		game.FinishedAt = &now
		game.WinnerID = &scoringPlayerID
		return database.DB.Save(game).Error
	}

	return nil
}

func RemoveScore(gameID, playerID, objectiveID int) error {
	return database.DB.
		Table("scores").
		Where("game_id = ? AND player_id = ? AND objective_id = ?", gameID, playerID, objectiveID).
		Delete(nil).Error
}
