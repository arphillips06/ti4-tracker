// achievements/types.go
package achievements

import achievements_helper "github.com/arphillips06/TI4-stats/helpers/achievements"

type Holder struct {
	PlayerID uint  `json:"player_id"`
	GameID   *uint `json:"game_id,omitempty"`
	RoundID  *uint `json:"round_id,omitempty"`
}

type Badge struct {
	Key     string                       `json:"key"`
	Label   string                       `json:"label"`
	Value   int                          `json:"value"`
	Status  string                       `json:"status,omitempty"`
	Holders []achievements_helper.Holder `json:"holders"`
}

type intVal struct{ Value *int }

type roundTotal struct {
	PlayerID uint
	RoundID  uint
	Total    int
}
