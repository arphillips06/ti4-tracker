package models

import "time"

//Game represents a single game
type Game struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	WinnerID  uint     //playerID of winner
	Rounds    []Round  `gorm:"foreignKey:GameID"`
	Players   []Player `gorm:"many2many:game_players;"`
}

type Player struct {
	ID      uint `gorm:"primaryKey"`
	Name    string
	Faction string
	Games   []Game `gorm:"many2many:game_players;"`
}

type Round struct {
	ID     uint `gorm:"primaryKey"`
	GameID uint
	Number int     //round number
	Scores []Score `gorm:"foreignKey:RoundID"`
}

type Score struct {
	ID          uint `gorm:"primaryKey"`
	RoundID     uint
	PlayerID    uint
	GameID      uint
	ObjectiveID uint `gorm:"not null"`
	Objective   Objective
	Points      int
}

type Objective struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Type        string
	Description string
	Points      int
}
