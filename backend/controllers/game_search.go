// controllers/games_search.go
package controllers

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/arphillips06/TI4-stats/database"
	"github.com/arphillips06/TI4-stats/models"
	"gorm.io/gorm"
)

var tokRe = regexp.MustCompile(`"([^"]+)"|(\S+)`)

type searchFilters struct {
	Winner     string
	Player     string
	Faction    string
	Agenda     string
	Relic      string
	Custodians *bool
	RoundsOp   string // "=", ">=", "<="
	RoundsVal  *int
	After      *time.Time
	Before     *time.Time
	FreeText   []string
}

func parseSearchQuery(q string) (searchFilters, error) {
	q = strings.TrimSpace(q)
	var f searchFilters
	if q == "" {
		return f, nil
	}

	matches := tokRe.FindAllStringSubmatch(q, -1)
	for _, m := range matches {
		token := m[1]
		if token == "" {
			token = m[2]
		}
		if token == "" {
			continue
		}

		lc := strings.ToLower(token)
		switch {
		case strings.HasPrefix(lc, "winner:"):
			f.Winner = strings.Trim(strings.TrimPrefix(token, "winner:"), `"`)
		case strings.HasPrefix(lc, "player:"):
			f.Player = strings.Trim(strings.TrimPrefix(token, "player:"), `"`)
		case strings.HasPrefix(lc, "faction:"):
			f.Faction = strings.Trim(strings.TrimPrefix(token, "faction:"), `"`)
		case strings.HasPrefix(lc, "agenda:"):
			f.Agenda = strings.Trim(strings.TrimPrefix(token, "agenda:"), `"`)
		case strings.HasPrefix(lc, "relic:"):
			f.Relic = strings.Trim(strings.TrimPrefix(token, "relic:"), `"`)
		case strings.HasPrefix(lc, "custodians:"):
			val := strings.Trim(strings.TrimPrefix(lc, "custodians:"), `"`)
			switch val {
			case "true", "1", "yes":
				t := true
				f.Custodians = &t
			case "false", "0", "no":
				t := false
				f.Custodians = &t
			}
		case strings.HasPrefix(lc, "rounds>="), strings.HasPrefix(lc, "rounds<="), strings.HasPrefix(lc, "rounds="):
			op := ">="
			if strings.Contains(lc, "<=") {
				op = "<="
			}
			if strings.Contains(lc, "rounds=") && !strings.Contains(lc, "<=") && !strings.Contains(lc, ">=") {
				op = "="
			}
			parts := strings.SplitN(token, op, 2)
			if len(parts) == 2 {
				if n, err := strconv.Atoi(parts[1]); err == nil {
					f.RoundsOp = op
					f.RoundsVal = &n
				}
			}
		case strings.HasPrefix(lc, "after:"):
			if t, err := time.Parse("2006-01-02", token[len("after:"):]); err == nil {
				f.After = &t
			}
		case strings.HasPrefix(lc, "before:"):
			if t, err := time.Parse("2006-01-02", token[len("before:"):]); err == nil {
				f.Before = &t
			}
		default:
			f.FreeText = append(f.FreeText, strings.Trim(token, `"`))
		}
	}
	return f, nil
}

func applyGameSearch(db *gorm.DB, f searchFilters) *gorm.DB {
	db = db.Table("games").Where("games.partial = FALSE")

	if f.After != nil {
		db = db.Where("games.finished_at IS NOT NULL AND games.finished_at >= ?", *f.After)
	}
	if f.Before != nil {
		db = db.Where("games.finished_at IS NOT NULL AND games.finished_at < ?", *f.Before)
	}
	if f.Winner != "" {
		name := "%" + strings.ToLower(f.Winner) + "%"
		db = db.Where(`EXISTS (SELECT 1 FROM players w WHERE w.id = games.winner_id AND LOWER(w.name) LIKE ?)`, name)
	}
	if f.Player != "" {
		name := "%" + strings.ToLower(f.Player) + "%"
		db = db.Where(`
			EXISTS (
			  SELECT 1
			  FROM game_players gp
			  JOIN players p ON p.id = gp.player_id
			  WHERE gp.game_id = games.id AND LOWER(p.name) LIKE ?
			)`, name)
	}
	if f.Faction != "" {
		name := "%" + strings.ToLower(f.Faction) + "%"
		db = db.Where(`
			EXISTS (
			  SELECT 1
			  FROM game_players gp
			  JOIN factions f ON f.id = gp.faction_id
			  WHERE gp.game_id = games.id AND LOWER(f.name) LIKE ?
			)`, name)
	}
	if f.Agenda != "" {
		title := "%" + strings.ToLower(f.Agenda) + "%"
		db = db.Where(`
			EXISTS (
			  SELECT 1 FROM scores s
			  WHERE s.game_id = games.id AND s.type='agenda' AND LOWER(s.agenda_title) LIKE ?
			)`, title)
	}
	if f.Relic != "" {
		title := "%" + strings.ToLower(f.Relic) + "%"
		db = db.Where(`
			EXISTS (
			  SELECT 1 FROM scores s
			  WHERE s.game_id = games.id AND s.type='relic' AND LOWER(s.relic_title) LIKE ?
			)`, title)
	}
	if f.Custodians != nil {
		if *f.Custodians {
			db = db.Where(`EXISTS (SELECT 1 FROM scores s WHERE s.game_id=games.id AND s.type='mecatol' AND s.points>0)`)
		} else {
			db = db.Where(`NOT EXISTS (SELECT 1 FROM scores s WHERE s.game_id=games.id AND s.type='mecatol' AND s.points>0)`)
		}
	}
	if f.RoundsVal != nil && f.RoundsOp != "" {
		switch f.RoundsOp {
		case "=":
			db = db.Where(`(SELECT COALESCE(MAX(number),0) FROM rounds r WHERE r.game_id=games.id) = ?`, *f.RoundsVal)
		case ">=":
			db = db.Where(`(SELECT COALESCE(MAX(number),0) FROM rounds r WHERE r.game_id=games.id) >= ?`, *f.RoundsVal)
		case "<=":
			db = db.Where(`(SELECT COALESCE(MAX(number),0) FROM rounds r WHERE r.game_id=games.id) <= ?`, *f.RoundsVal)
		}
	}
	if len(f.FreeText) > 0 {
		for _, term := range f.FreeText {
			like := "%" + strings.ToLower(term) + "%"
			db = db.Where("(LOWER(games.title) LIKE ? OR LOWER(games.notes) LIKE ?)", like, like)
		}
	}
	return db
}

func listGamesWithSearch(q string) *gorm.DB {
	f, _ := parseSearchQuery(q)
	return applyGameSearch(database.DB.Model(&models.Game{}), f).
		Order("COALESCE(games.finished_at, games.created_at) DESC")
}
