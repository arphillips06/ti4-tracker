package stats

import (
	"errors"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
)

func CalculateObjectiveCounts() (map[string]int, error) {
	result := make(map[string]int)

	var publicCount, secretCount, stage1Count, stage2Count, cdlCount int64

	// Count unique (game_id, objective_id) pairs for each category

	err := database.DB.
		Table("scores").
		Select("COUNT(DISTINCT game_id || '-' || objective_id)").
		Where("type = ?", "public").
		Scan(&publicCount).Error
	if err != nil {
		return nil, err
	}
	result["publicScored"] = int(publicCount)

	err = database.DB.
		Table("scores").
		Select("COUNT(DISTINCT game_id || '-' || objective_id)").
		Where("type = ?", "secret").
		Scan(&secretCount).Error
	if err != nil {
		return nil, err
	}
	result["secretScored"] = int(secretCount)

	err = database.DB.
		Table("scores").
		Joins("JOIN objectives ON scores.objective_id = objectives.id").
		Select("COUNT(DISTINCT scores.game_id || '-' || scores.objective_id)").
		Where("objectives.stage = ?", "I").
		Scan(&stage1Count).Error
	if err != nil {
		return nil, err
	}
	result["stage1Scored"] = int(stage1Count)

	err = database.DB.
		Table("scores").
		Joins("JOIN objectives ON scores.objective_id = objectives.id").
		Select("COUNT(DISTINCT scores.game_id || '-' || scores.objective_id)").
		Where("objectives.stage = ?", "II").
		Scan(&stage2Count).Error
	if err != nil {
		return nil, err
	}
	result["stage2Scored"] = int(stage2Count)

	// CDL tracking (if applicable)
	err = database.DB.
		Table("scores").
		Where("agenda_title = ?", "Classified Document Leaks").
		Count(&cdlCount).Error
	if err != nil {
		return nil, err
	}
	result["cdlPromoted"] = int(cdlCount)

	return result, nil
}

func CalculateObjectiveFrequencies() (map[string]int, map[string]int, error) {
	publicMap := make(map[string]int)
	secretMap := make(map[string]int)

	var publicRows, secretRows []struct {
		Name  string
		Count int
	}

	// Public objectives
	err := database.DB.Table("scores").
		Joins("JOIN objectives ON scores.objective_id = objectives.id").
		Select("objectives.name, COUNT(DISTINCT scores.game_id) as count").
		Where("objectives.stage IN ('I', 'II')").
		Group("objectives.name").
		Scan(&publicRows).Error
	if err != nil {
		return nil, nil, err
	}
	for _, row := range publicRows {
		publicMap[row.Name] = row.Count
	}

	// Secret objectives
	err = database.DB.Table("scores").
		Joins("JOIN objectives ON scores.objective_id = objectives.id").
		Select("objectives.name, COUNT(DISTINCT scores.game_id) as count").
		Where("objectives.stage = 'S'").
		Group("objectives.name").
		Scan(&secretRows).Error
	if err != nil {
		return nil, nil, err
	}
	for _, row := range secretRows {
		secretMap[row.Name] = row.Count
	}

	return publicMap, secretMap, nil
}

func CalculateObjectiveAppearanceStats(totalGames int64) (map[string]models.ObjectiveStats, error) {
	if totalGames == 0 {
		return nil, errors.New("cannot calculate stats with 0 total games")
	}

	type ObjectiveRow struct {
		Name      string
		Type      string
		GameCount int
	}

	var appearances []ObjectiveRow
	var scored []ObjectiveRow

	// Only count revealed objectives that are not secret
	err := database.DB.
		Table("game_objectives").
		Select("objectives.name, objectives.type, COUNT(DISTINCT game_objectives.game_id) as game_count").
		Joins("JOIN objectives ON game_objectives.objective_id = objectives.id").
		Joins("JOIN games ON game_objectives.game_id = games.id").
		Where("games.partial = false AND game_objectives.revealed = true AND objectives.type != ?", "secret").
		Group("objectives.name, objectives.type").
		Scan(&appearances).Error
	if err != nil {
		return nil, err
	}

	// Only count scores for non-secret objectives
	err = database.DB.
		Model(&models.Score{}).
		Select("objectives.name, objectives.type, COUNT(DISTINCT scores.game_id) as game_count").
		Joins("JOIN objectives ON scores.objective_id = objectives.id").
		Joins("JOIN games ON scores.game_id = games.id").
		Where("games.partial = false AND objectives.type != ?", "secret").
		Group("objectives.name, objectives.type").
		Scan(&scored).Error
	if err != nil {
		return nil, err
	}

	// Build lookup maps
	appearanceMap := make(map[string]ObjectiveRow)
	scoredMap := make(map[string]ObjectiveRow)

	for _, row := range appearances {
		appearanceMap[row.Name] = row
	}
	for _, row := range scored {
		scoredMap[row.Name] = row
	}

	// Merge into result
	result := make(map[string]models.ObjectiveStats)
	allNames := make(map[string]bool)
	for name := range appearanceMap {
		allNames[name] = true
	}
	for name := range scoredMap {
		allNames[name] = true
	}

	for name := range allNames {
		appear := appearanceMap[name]
		score := scoredMap[name]

		appeared := appear.GameCount
		scored := score.GameCount

		// fallback in case type is missing from one of the maps
		typ := appear.Type
		if typ == "" {
			typ = score.Type
		}

		appearanceRate := float64(appeared) / float64(totalGames) * 100
		scoreRate := 0.0
		if appeared > 0 {
			scoreRate = float64(scored) / float64(appeared) * 100
		}

		result[name] = models.ObjectiveStats{
			Type:                   typ,
			AppearanceRate:         appearanceRate,
			ScoredWhenAppearedRate: scoreRate,
			AppearedCount:          appeared,
			ScoredCount:            scored,
		}
	}

	return result, nil
}

func CalculateSecretObjectiveRates() ([]models.SecretObjectiveRate, error) {
	type Result struct {
		Name         string
		GamesPlayed  int64
		SecretScored int64
	}

	var rows []Result

	// Subquery: games played per player
	subGamesPlayed := database.DB.
		Table("game_players").
		Joins("JOIN games ON games.id = game_players.game_id").
		Where("games.partial = false").
		Select("player_id, COUNT(DISTINCT game_id) AS games_played").
		Group("player_id")

		// Subquery: secret scored per player
	subSecrets := database.DB.
		Table("scores").
		Joins("JOIN games ON games.id = scores.game_id").
		Where("games.partial = false AND type = ?", "secret").
		Select("player_id, COUNT(DISTINCT scores.id) AS secret_scored").
		Group("player_id")

	// Join both subqueries on player_id
	err := database.DB.
		Table("players AS p").
		Select("p.name, COALESCE(gp.games_played, 0) AS games_played, COALESCE(ss.secret_scored, 0) AS secret_scored").
		Joins("LEFT JOIN (?) AS gp ON p.id = gp.player_id", subGamesPlayed).
		Joins("LEFT JOIN (?) AS ss ON p.id = ss.player_id", subSecrets).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	// Compute final score rates
	var result []models.SecretObjectiveRate
	for _, row := range rows {
		appeared := int(row.GamesPlayed) * 3 // Max allowed draws per player
		rate := 0.0
		if appeared > 0 {
			rate = float64(row.SecretScored) / float64(appeared) * 100
		}
		result = append(result, models.SecretObjectiveRate{
			Player:          row.Name,
			SecretAppeared:  int(appeared),
			SecretScored:    int(row.SecretScored),
			SecretScoreRate: rate,
		})
	}

	return result, nil
}

func CalculateObjectiveMetaStats() ([]models.ObjectiveMeta, error) {
	var metas []models.ObjectiveMeta

	// Step 1: Get scored data (distinct games where it was scored)
	type ScoreStats struct {
		Name         string
		Type         string
		GamesScored  int
		AverageRound float64
	}
	var scoreStats []ScoreStats

	err := database.DB.
		Table("scores").
		Select(`
			objectives.name AS name,
			objectives.type AS type,
			COUNT(DISTINCT scores.game_id) AS games_scored,
			AVG(rounds.number) AS average_round
		`).
		Joins("JOIN objectives ON scores.objective_id = objectives.id").
		Joins("JOIN games ON scores.game_id = games.id").
		Joins("JOIN rounds ON scores.round_id = rounds.id").
		Where("games.partial = ?", false).
		Group("objectives.name, objectives.type").
		Scan(&scoreStats).Error
	if err != nil {
		return nil, err
	}

	metaMap := map[string]models.ObjectiveMeta{}
	for _, stat := range scoreStats {
		metaMap[stat.Name] = models.ObjectiveMeta{
			Name:         stat.Name,
			Type:         stat.Type,
			TimesScored:  stat.GamesScored, // ✅ now distinct games
			AverageRound: stat.AverageRound,
		}
	}

	// Step 2: Get appearance counts
	type Appearance struct {
		Name  string
		Type  string
		Count int
	}

	var appearances []Appearance

	err = database.DB.
		Table("game_objectives").
		Select("objectives.name, objectives.type, COUNT(DISTINCT game_objectives.game_id) as count").
		Joins("JOIN objectives ON game_objectives.objective_id = objectives.id").
		Joins("JOIN games ON game_objectives.game_id = games.id").
		Where("games.partial = ? AND game_objectives.revealed = ?", false, true).
		Group("objectives.name").
		Scan(&appearances).Error
	if err != nil {
		return nil, err
	}

	for _, a := range appearances {
		meta := metaMap[a.Name]
		meta.Name = a.Name
		meta.Type = a.Type
		meta.TimesAppeared = a.Count
		if a.Count > 0 {
			meta.ScoredPercent = float64(meta.TimesScored) / float64(a.Count) * 100
		}
		metaMap[a.Name] = meta
	}

	// Convert to slice
	for _, meta := range metaMap {
		metas = append(metas, meta)
	}

	return metas, nil
}
