package services

import (
	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
)

type StatsOverview struct {
	TotalGames                 int                                  `json:"totalGames"`
	GamesWonByFaction          map[string]int                       `json:"gamesWonByFaction"`
	GamesPlayedByFaction       map[string]int                       `json:"gamesPlayedByFaction"`
	WinRateByFaction           map[string]float64                   `json:"winRateByFaction"`
	ObjectiveStats             map[string]int                       `json:"objectiveStats"`
	ObjectiveFrequency         map[string]int                       `json:"objectiveFrequency"`
	PlayerWinRates             []models.PlayerWinRate               `json:"playerWinRates"`
	ObjectiveAppearanceStats   map[string]models.ObjectiveStats     `json:"objectiveAppearanceStats"`
	FactionPlayWinDistribution map[string]models.FactionPlayWinStat `json:"factionPlayWinDistribution"`
	PlayerAveragePoints        []models.PlayerAveragePoints         `json:"playerAveragePoints"`
	TopFactionsPerPlayer       []models.PlayerFactionStats          `json:"topFactionsPerPlayer"`
	PlayerMostCommonFinishes   []models.PlayerMostCommonFinish      `json:"playerMostCommonFinishes"`
	SecretObjectiveRates       []models.SecretObjectiveRate         `json:"secretObjectiveRates"`
	PlayerPointStdevs          []models.PlayerPointStdev            `json:"playerPointStdevs"`
	AveragePlayerPoints        float64                              `json:"averagePlayerPoints"`
	TotalUniquePlayers         int                                  `json:"totalUniquePlayers"`
	MostPlayedFaction          string                               `json:"mostPlayedFaction"`
	MostVictoriousFaction      string                               `json:"mostVictoriousFaction"`
	AverageGameRounds          float64                              `json:"averageGameRounds"`
}

func CalculateStatsOverview() (*StatsOverview, error) {
	totalGames, err := helpers.CountTotalGames()
	if err != nil {
		return nil, err
	}

	factionPlays, factionWins, winRates, playWinDist, err := helpers.CalculateFactionStats()
	if err != nil {
		return nil, err
	}

	objectiveStats, err := helpers.CalculateObjectiveCounts()
	if err != nil {
		return nil, err
	}

	objectiveFreq, err := helpers.CalculateObjectiveFrequencies()
	if err != nil {
		return nil, err
	}

	playerWinRates, err := helpers.CalculatePlayerWinRates()
	if err != nil {
		return nil, err
	}

	playerAverages, err := helpers.CalculatePlayerAverages()
	if err != nil {
		return nil, err
	}

	topFactionsPerPlayer, err := helpers.CalculateTopFactionsPerPlayer()
	if err != nil {
		return nil, err
	}

	playerFinishes, err := helpers.CalculateMostCommonFinishes()
	if err != nil {
		return nil, err
	}

	secretRates, err := helpers.CalculateSecretObjectiveRates()
	if err != nil {
		return nil, err
	}

	pointStdevs, err := helpers.CalculatePointStandardDeviations()
	if err != nil {
		return nil, err
	}

	avgRounds, err := helpers.CalculateAverageRounds()
	if err != nil {
		return nil, err
	}

	avgPoints, err := helpers.CalculateAveragePlayerPoints()
	if err != nil {
		return nil, err
	}

	totalPlayers, err := helpers.CountUniquePlayers()
	if err != nil {
		return nil, err
	}

	mostPlayed, mostVictorious := helpers.DetermineMostPlayedAndVictoriousFactions(factionPlays, factionWins)

	objectiveAppearanceStats, err := helpers.CalculateObjectiveAppearanceStats(totalGames)
	if err != nil {
		return nil, err
	}

	return &StatsOverview{
		TotalGames:                 int(totalGames),
		GamesPlayedByFaction:       factionPlays,
		GamesWonByFaction:          factionWins,
		WinRateByFaction:           winRates,
		ObjectiveStats:             objectiveStats,
		ObjectiveFrequency:         objectiveFreq,
		PlayerWinRates:             playerWinRates,
		ObjectiveAppearanceStats:   objectiveAppearanceStats,
		FactionPlayWinDistribution: playWinDist,
		PlayerAveragePoints:        playerAverages,
		TopFactionsPerPlayer:       topFactionsPerPlayer,
		PlayerMostCommonFinishes:   playerFinishes,
		SecretObjectiveRates:       secretRates,
		PlayerPointStdevs:          pointStdevs,
		AveragePlayerPoints:        avgPoints,
		TotalUniquePlayers:         int(totalPlayers),
		MostPlayedFaction:          mostPlayed,
		MostVictoriousFaction:      mostVictorious,
		AverageGameRounds:          avgRounds,
	}, nil
}
