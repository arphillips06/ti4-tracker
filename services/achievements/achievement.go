// achievements/types.go
package achievements

type Holder struct {
	PlayerID uint  `json:"player_id"`
	GameID   *uint `json:"game_id,omitempty"`
	RoundID  *uint `json:"round_id,omitempty"`
}

type Badge struct {
	Key     string   `json:"key"`     // e.g. "most_points_in_round"
	Label   string   `json:"label"`   // human label (optional)
	Value   int      `json:"value"`   // e.g. 10 points, 4 rounds
	Status  string   `json:"status"`  // "new" | "tied"
	Holders []Holder `json:"holders"` // who achieved it (in this game)
}
