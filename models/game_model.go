package models

import "time"

//these structs are to be used with the SQL database
//Game represents a single game
type Game struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	CreatedAt         time.Time       `json:"created_at"`
	FinishedAt        *time.Time      `json:"finished_at"`
	WinnerID          *uint           `json:"winner_id"`
	Winner            Player          `gorm:"foreignKey:WinnerID" json:"winner"`
	Rounds            []Round         `gorm:"foreignKey:GameID" json:"rounds"`
	GamePlayers       []GamePlayer    `gorm:"foreignKey:GameID" json:"players"`
	WinningPoints     int             `json:"winning_points"`
	CurrentRound      int             `gorm:"default:1" json:"current_round"`
	GameObjectives    []GameObjective `json:"game_objectives"`
	UseObjectiveDecks bool            `json:"use_objective_decks"`
	Partial           bool            `gorm:"default:false"`
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
	ID               uint `gorm:"primaryKey"`
	RoundID          uint
	PlayerID         uint
	GameID           uint
	ObjectiveID      uint `gorm:"not null"`
	Objective        Objective
	Points           int
	Round            Round     `gorm:"foreignKey:RoundID"`
	Player           Player    `gorm:"foreignKey:PlayerID"`
	Type             string    `gorm:"type:VARCHAR(20)"` //e.g. objective, imperial, support
	AgendaTitle      string    `gorm:"type:VARCHAR(100)"`
	RelicTitle       string    `gorm:"type:VARCHAR(20)"`
	CreatedAt        time.Time `json:"created_at"`
	OriginallySecret bool      `gorm:"default:false"`
}

//Objective information
type Objective struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Points      int    `json:"points"`
	Stage       string `gorm:"type:VARCHAR(5)" json:"stage"`
	Phase       string `gorm:"type:VARCHAR(10)" json:"phase"`
}

//links game and player together into one struct
type GamePlayer struct {
	ID       uint `gorm:"PrimaryKey"`
	GameID   uint
	PlayerID uint
	Faction  string
	Player   Player `gorm:"foreignKey:PlayerID"`
	Game     Game   `gorm:"foreignKey:GameID;references:ID" json:"-"`
	Won      bool
}

//links game and ovjective into one struct
type GameObjective struct {
	ID          uint      `json:"ID"`
	GameID      uint      `json:"GameID"`
	ObjectiveID uint      `json:"ObjectiveID"`
	RoundID     uint      `json:"RoundID"`
	Stage       string    `json:"Stage"`
	Revealed    bool      `gorm:"default:false"`
	Objective   Objective `gorm:"foreignKey:ObjectiveID;references:ID" json:"Objective"`
	Round       Round     `gorm:"foreignKey:RoundID;references:ID" json:"Round"`
	IsCDL       bool      `json:"IsCDL" gorm:"-"`
	Position    int
}

type PlayerInput struct {
	ID      string
	Name    string
	Faction string
}

type AssignObjectiveRequest struct {
	GameID      uint `json:"game_id"`
	RoundID     int  `json:"round_id"`
	ObjectiveID uint `json:"objective_id"`
}
