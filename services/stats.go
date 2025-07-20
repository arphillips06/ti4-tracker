package services

import (
	"github.com/arphillips06/TI4-stats/helpers/stats"
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
	CustodiansStats            []PlayerCustodiansStats              `json:"custodiansStats"`
	FactionPlayerStats         []models.FactionPlayerStats          `json:"factionPlayerStats"`
	GameLengthStats            models.GameLengthStats               `json:"gameLengthStats"`
	FactionAggregateStats      []models.FactionAggregateStats       `json:"factionAggregateStats"`
	SecretObjectiveFrequency   map[string]int                       `json:"publicSecretFrequency"`
	PublicObjectiveFrequency   map[string]int                       `json:"publicObjectiveFrequency"`
	ObjectiveMetaStats         []models.ObjectiveMeta               `json:"objectiveMetaStats"`
	PointSpreadDistribution    map[int]int                          `json:"pointSpreadDistribution"`
	GameLengthDistribution     map[int]int                          `json:"gameLengthDistribution"`
	CommonVictoryPaths         map[string]int                       `json:"commonVictoryPaths"`
}

var CachedVictoryPathCounts = map[string]int{}

func CalculateStatsOverview() (*StatsOverview, error) {
	totalGames, err := stats.CountTotalGames()
	if err != nil {
		return nil, err
	}

	factionPlays, factionWins, winRates, playWinDist, err := stats.CalculateFactionStats()
	if err != nil {
		return nil, err
	}

	objectiveStats, err := stats.CalculateObjectiveCounts()
	if err != nil {
		return nil, err
	}

	publicFreq, secretFreq, err := stats.CalculateObjectiveFrequencies()
	if err != nil {
		return nil, err
	}

	playerWinRates, err := stats.CalculatePlayerWinRates()
	if err != nil {
		return nil, err
	}

	playerAverages, err := stats.CalculatePlayerAverages()
	if err != nil {
		return nil, err
	}

	topFactionsPerPlayer, err := stats.CalculateTopFactionsPerPlayer()
	if err != nil {
		return nil, err
	}

	playerFinishes, err := stats.CalculateMostCommonFinishes()
	if err != nil {
		return nil, err
	}

	secretRates, err := stats.CalculateSecretObjectiveRates()
	if err != nil {
		return nil, err
	}

	pointStdevs, err := stats.CalculatePointStandardDeviations()
	if err != nil {
		return nil, err
	}

	avgRounds, err := stats.CalculateAverageRounds()
	if err != nil {
		return nil, err
	}

	avgPoints, err := stats.CalculateAveragePlayerPoints()
	if err != nil {
		return nil, err
	}

	totalPlayers, err := stats.CountUniquePlayers()
	if err != nil {
		return nil, err
	}

	mostPlayed, mostVictorious := stats.DetermineMostPlayedAndVictoriousFactions(factionPlays, factionWins)

	objectiveAppearanceStats, err := stats.CalculateObjectiveAppearanceStats(totalGames)
	if err != nil {
		return nil, err
	}

	gameLengthStats, err := stats.GetGameLengthStats()
	if err != nil {
		return nil, err
	}
	factionPlayerStats, err := stats.GetFactionPlayerStats()
	if err != nil {
		return nil, err
	}
	factionAggStats, err := stats.GetFactionAggregateStats()
	if err != nil {
		return nil, err
	}
	objectiveMetaStats, err := stats.CalculateObjectiveMetaStats()
	if err != nil {
		return nil, err
	}
	pointSpreads, err := stats.CalculateVictoryPointSpreads()
	if err != nil {
		return nil, err
	}

	lengths, err := stats.CalculateGameLengthDistribution()
	if err != nil {
		return nil, err
	}

	victoryPaths, err := stats.CalculateCommonVictoryPaths()
	if err != nil {
		return nil, err
	}

	return &StatsOverview{
		TotalGames:                 int(totalGames),
		GamesPlayedByFaction:       factionPlays,
		GamesWonByFaction:          factionWins,
		WinRateByFaction:           winRates,
		ObjectiveStats:             objectiveStats,
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
		FactionPlayerStats:         factionPlayerStats,
		GameLengthStats:            gameLengthStats,
		FactionAggregateStats:      factionAggStats,
		PublicObjectiveFrequency:   publicFreq,
		SecretObjectiveFrequency:   secretFreq,
		ObjectiveMetaStats:         objectiveMetaStats,
		PointSpreadDistribution:    pointSpreads,
		GameLengthDistribution:     lengths,
		CommonVictoryPaths:         victoryPaths,
	}, nil
}
