package models

type StatsOverview struct {
	TotalGames                 int                           `json:"totalGames"`
	GamesWonByFaction          map[string]int                `json:"gamesWonByFaction"`
	GamesPlayedByFaction       map[string]int                `json:"gamesPlayedByFaction"`
	WinRateByFaction           map[string]float64            `json:"winRateByFaction"`
	ObjectiveStats             map[string]int                `json:"objectiveStats"`
	ObjectiveFrequency         map[string]int                `json:"objectiveFrequency"`
	PlayerWinRates             []PlayerWinRate               `json:"playerWinRates"`
	ObjectiveAppearanceStats   map[string]ObjectiveStats     `json:"objectiveAppearanceStats"`
	FactionPlayWinDistribution map[string]FactionPlayWinStat `json:"factionPlayWinDistribution"`
	PlayerAveragePoints        []PlayerAveragePoints         `json:"playerAveragePoints"`
	TopFactionsPerPlayer       []PlayerFactionStats          `json:"topFactionsPerPlayer"`
	PlayerMostCommonFinishes   []PlayerMostCommonFinish      `json:"playerMostCommonFinishes"`
	SecretObjectiveRates       []SecretObjectiveRate         `json:"secretObjectiveRates"`
	HeadlineStats              HeadlineStats                 `json:"headlineStats"`
}

type FactionPlayWinStat struct {
	PlayedCount int     `json:"playedCount"`
	WinCount    int     `json:"winCount"`
	PlayRate    float64 `json:"playRate"`
	WinRate     float64 `json:"winRate"`
}
type PlayerPointStdev struct {
	Player string  `json:"player"`
	Stdev  float64 `json:"stdev"`
}

type PlayerAveragePoints struct {
	Player        string  `json:"player"`
	GamesPlayed   int     `json:"gamesPlayed"`
	TotalPoints   float64 `json:"totalPoints"`
	AveragePoints float64 `json:"averagePoints"`
	Stdev         float64 `json:"stdev"`
}

type PlayerMostCommonFinish struct {
	Player     string `json:"player"`
	Position   int    `json:"position"`
	Count      int    `json:"count"`
	TotalGames int    `json:"totalGames"`
}

type PlayerWinRate struct {
	Player      string  `json:"player"`
	GamesPlayed int     `json:"gamesPlayed"`
	GamesWon    int     `json:"gamesWon"`
	WinRate     float64 `json:"winRate"`
}

type ObjectiveStats struct {
	AppearanceRate         float64 `json:"appearanceRate"`
	ScoredWhenAppearedRate float64 `json:"scoredWhenAppearedRate"`
}

type PlayerFactionStats struct {
	Player   string         `json:"player"`
	Factions map[string]int `json:"factions"`
}

type SecretObjectiveRate struct {
	Player          string  `json:"player"`
	SecretAppeared  int     `json:"secretAppeared"`
	SecretScored    int     `json:"secretScored"`
	SecretScoreRate float64 `json:"secretScoreRate"`
}

type HeadlineStats struct {
	AverageGameRounds      float64 `json:"averageGameRounds"`
	AveragePointsPerPlayer float64 `json:"averagePointsPerPlayer"`
	TotalUniquePlayers     int     `json:"totalUniquePlayers"`
	MostPlayedFaction      string  `json:"mostPlayedFaction"`
	MostVictoriousFaction  string  `json:"mostVictoriousFaction"`
}
