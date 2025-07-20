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
	ID                 uint                 `json:"id"`
	WinningPoints      int                  `json:"winning_points"`
	CurrentRound       int                  `json:"current_round"`
	FinishedAt         *time.Time           `json:"finished_at"`
	UseObjectiveDecks  bool                 `json:"use_objective_decks"`
	Players            []GamePlayer         `json:"players"`
	Rounds             []Round              `json:"rounds"`
	Objectives         []GameObjective      `json:"objectives"`
	Scores             []PlayerScoreSummary `json:"scores"`
	AllScores          []Score              `json:"all_scores"`
	Winner             *Player              `json:"winner"`
	CustodiansPlayerID *uint                `json:"custodiansPlayerId,omitempty"`
	ScoresByObjective  map[uint][]Score     `json:"ScoresByObjective"`
	WinnerVictoryPath  *VictoryPathSummary  `json:"victory_path,omitempty"`
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

type ObjectiveScoreSummary struct {
	ObjectiveID uint     `json:"objective_id"`
	Name        string   `json:"name"`
	Stage       string   `json:"stage"`
	ScoredBy    []string `json:"scored_by"` // player names
}

type AgendaResolution struct {
	GameID   uint   `json:"game_id"`
	RoundID  uint   `json:"round_id"`
	Result   string `json:"result"`
	ForVotes []uint `json:"for_votes"`
}

type PoliticalCensureRequest struct {
	GameID   uint `json:"game_id"`
	RoundID  uint `json:"round_id"`
	PlayerID uint `json:"player_id"`
	Gained   bool `json:"gained"`
}

type SeedOfEmpireResolution struct {
	GameID  uint   `json:"game_id"`
	RoundID uint   `json:"round_id"`
	Result  string `json:"result"` // "for" or "against"
}

type ClassifiedDocumentLeaksRequest struct {
	GameID      uint `json:"game_id"`
	RoundID     uint `json:"round_id"`
	PlayerID    uint `json:"player_id"`
	ObjectiveID uint `json:"objective_id"`
}

type ObjectiveWithMetadata struct {
	Objective Objective `json:"Objective"`
	IsCDL     bool      `json:"IsCDL"`
}

type IncentiveProgramRequest struct {
	GameID  uint   `json:"game_id"`
	Outcome string `json:"outcome"`
}

type ObjectiveDeck struct {
	ID          uint `gorm:"primaryKey"`
	GameID      uint
	Stage       string // "I" or "II"
	ObjectiveID uint
	Assigned    bool
	Position    int
	Objective   Objective `gorm:"foreignKey:ObjectiveID"`
}

type AssignPlayerInput struct {
	GameID   uint   `json:"game_id"`
	PlayerID uint   `json:"player_id"`
	Faction  string `json:"faction"`
}
