package models

import "time"

//Game represents a single game
type Game struct {
	ID                uint `gorm:"primaryKey"`
	CreatedAt         time.Time
	FinishedAt        *time.Time
	WinnerID          uint
	Rounds            []Round      `gorm:"foreignKey:GameID"`
	GamePlayers       []GamePlayer `gorm:"foreignKey:GameID"`
	WinningPoints     int
	CurrentRound      int `gorm:"default:1"`
	GameObjectives    []GameObjective
	UseObjectiveDecks bool `json:"use_objective_decks"`
}

//Single player
type Player struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Games []GamePlayer `gorm:"foreignKey:PlayerID" json:"-"`
}

//Round counter
type Round struct {
	ID     uint `gorm:"primaryKey"`
	GameID uint
	Number int
	Scores []Score `gorm:"foreignKey:RoundID"`
}

//Scoring information
type Score struct {
	ID          uint `gorm:"primaryKey"`
	RoundID     uint
	PlayerID    uint
	GameID      uint
	ObjectiveID uint `gorm:"not null"`
	Objective   Objective
	Points      int
	Round       Round  `gorm:"foreignKey:RoundID"`
	Player      Player `gorm:"foreignKey:PlayerID"`
	Type        string `gorm:"type:VARCHAR(20)"` //e.g. objective, imperial, support
}

//Objective information
type Objective struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Type        string
	Description string
	Points      int
	Stage       string `gorm:"type:VARCHAR(5)"`
	Phase       string `gorm:"type:VARCHAR(10)"`
}

//links game and player together into one struct
type GamePlayer struct {
	ID       uint `gorm:"PrimaryKey"`
	GameID   uint
	PlayerID uint
	Faction  string
	Player   Player `gorm:"foreignKey:PlayerID"`
	Game     Game   `gorm:"foreignKey:GameID;references:ID" json:"-"`
}

//links game and ovjective into one struct
type GameObjective struct {
	ID          uint `gorm:"primaryKey"`
	GameID      uint
	ObjectiveID uint
	RoundID     uint   // null if revealed at game start
	Stage       string `gorm:"type:VARCHAR(10)"`
	Objective   Objective
	Round       Round
}
type PlayerInput struct {
	ID      string
	Name    string
	Faction string
}
