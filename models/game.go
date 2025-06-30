package models

import "time"

//Game represents a single game
type Game struct {
	ID            uint `gorm:"primaryKey"`
	CreatedAt     time.Time
	FinishedAt    *time.Time
	WinnerID      uint
	Rounds        []Round      `gorm:"foreignKey:GameID"`
	GamePlayers   []GamePlayer `gorm:"foreignKey:GameID"`
	WinningPoints int
	CurrentRound  int `gorm:"default:1"`
}

type Player struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Games []GamePlayer `gorm:"foreignKey:PlayerID" json:"-"`
}

type Round struct {
	ID     uint `gorm:"primaryKey"`
	GameID uint
	Number int
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

type GamePlayer struct {
	ID       uint `gorm:"PrimaryKey"`
	GameID   uint
	PlayerID uint
	Faction  string
	Player   Player `gorm:"foreignKey:PlayerID"`
	Game     Game   `gorm:"foreignKey:GameID;references:ID"`
}
