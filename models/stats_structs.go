package models

import "time"

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
	Type                   string  `json:"type"` // "public" or "secret"
	AppearanceRate         float64 `json:"appearanceRate"`
	ScoredWhenAppearedRate float64 `json:"scoredWhenAppearedRate"`
	AppearedCount          int     `json:"appearedCount"`
	ScoredCount            int     `json:"scoredCount"`
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

type FactionPlayerStats struct {
	Faction           string `json:"faction"`
	Player            string `json:"player"`
	PlayedCount       int    `json:"playedCount"`
	WonCount          int    `json:"wonCount"`
	TotalPointsScored int    `json:"totalPointsScored"`
}

type FactionAggregateStats struct {
	Faction           string `json:"faction"`
	TotalPlays        int    `json:"totalPlays"`
	TotalPointsScored int    `json:"totalPointsScored"`
	WonCount          int    `json:"wonCount"`
}

type GameDurationStat struct {
	GameID     uint      `json:"game_id"`
	RoundCount int       `json:"round_count"`
	Duration   string    `json:"duration"`
	Seconds    int64     `json:"seconds"`
	StartedAt  time.Time `json:"started_at"`
}

type GameLengthStats struct {
	All         GameLengthCategoryStats `json:"all"`
	ThreePlayer GameLengthCategoryStats `json:"three_player"`
	FourPlayer  GameLengthCategoryStats `json:"four_player"`
}

type GameLengthCategoryStats struct {
	LongestByRounds  GameDurationStat `json:"longest_by_rounds"`
	ShortestByRounds GameDurationStat `json:"shortest_by_rounds"`
	LongestByTime    GameDurationStat `json:"longest_by_time"`
	ShortestByTime   GameDurationStat `json:"shortest_by_time"`
	AverageRoundTime string           `json:"average_round_time"`
	AverageGameTime  string           `json:"average_game_time"`
}

type ObjectiveMeta struct {
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	TimesAppeared int     `json:"timesAppeared"`
	TimesScored   int     `json:"timesScored"`
	ScoredPercent float64 `json:"scoredPercent"`
	AverageRound  float64 `json:"averageRound"`
}
