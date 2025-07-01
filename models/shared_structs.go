package models

import "time"

type CreateGameInput struct {
	WinningPoints     int           `json:"winning_points"`
	UseObjectiveDecks *bool         `json:"use_objective_decks"`
	Players           []PlayerInput `json:"players"`
}

type PlayerScoreSummary struct {
	PlayerID   uint   `json:"player_id"`
	PlayerName string `json:"player_name"`
	Points     int    `json:"points"`
}

type RoundScore struct {
	Player    string `json:"player"`
	Objective string `json:"objective"`
	Points    int    `json:"points"`
}

type RoundScoresGroup struct {
	Round  int          `json:"round"`
	Scores []RoundScore `json:"scores"`
}

type GameDetailResponse struct {
	ID            uint                 `json:"id"`
	WinningPoints int                  `json:"winning_points"`
	CurrentRound  int                  `json:"current_round"`
	FinishedAt    *time.Time           `json:"finished_at"`
	Players       []GamePlayer         `json:"players"`
	Scores        []PlayerScoreSummary `json:"scores"`
}

type SelectedPlayersWithFaction struct {
	Player  Player
	Faction string
}

type ScoredObjective struct {
	ObjectiveID string
	PlayerID    string
	Round       int
}
