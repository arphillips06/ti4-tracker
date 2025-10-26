package services

import (
	"context"
	"math"
	"sort"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/arphillips06/TI4-stats/models"
)

type ObjectiveDifficultyOptions struct {
	Stage            string
	MinAppearances   int
	MinOpportunities int
}

func CalculateObjectiveDifficulty(ctx context.Context, db *gorm.DB, opts ObjectiveDifficultyOptions) (models.ObjectiveDifficultyResponse, error) {
	objMeta, err := loadObjectives(ctx, db, opts)
	if err != nil {
		return models.ObjectiveDifficultyResponse{}, err
	}

	appOpp, err := loadAppearancesAndOpportunities(ctx, db, opts)
	if err != nil {
		return models.ObjectiveDifficultyResponse{}, err
	}

	scoreCounts, firstRounds, err := loadScoresAndFirstRounds(ctx, db, opts)
	if err != nil {
		return models.ObjectiveDifficultyResponse{}, err
	}

	rows := make([]models.ObjectiveDifficultyRow, 0, len(objMeta))
	for _, om := range objMeta {
		A := appOpp[om.ID].Appearances
		O := appOpp[om.ID].Opportunities
		S := scoreCounts[om.ID]

		var raw, adj, lo, hi float64
		if O > 0 {
			raw = float64(S) / float64(O)
			adj = bayesAdj(int(S), int(O), 2, 5)
			lo, hi = wilsonInterval(int(S), int(O), 0.95)
		}

		avg, med := avgAndMedian(firstRounds[om.ID])

		row := models.ObjectiveDifficultyRow{
			ObjectiveID:   om.ID,
			Name:          om.Name,
			Stage:         om.Stage,
			Phase:         om.Phase,
			Appearances:   A,
			Opportunities: O,
			Scores:        S,
			RawRate:       raw,
			AdjRate:       adj,
			Difficulty:    1 - raw,
			WilsonLo:      lo,
			WilsonHi:      hi,
			AvgRound:      avg,
			MedianRound:   med,
		}

		if opts.MinAppearances > 0 && int(row.Appearances) < opts.MinAppearances {
			continue
		}
		if opts.MinOpportunities > 0 && int(row.Opportunities) < opts.MinOpportunities {
			continue
		}

		rows = append(rows, row)
	}

	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].Difficulty == rows[j].Difficulty {
			if rows[i].AdjRate == rows[j].AdjRate {
				return rows[i].Appearances > rows[j].Appearances
			}
			return rows[i].AdjRate < rows[j].AdjRate
		}
		return rows[i].Difficulty > rows[j].Difficulty
	})

	return models.ObjectiveDifficultyResponse{
		Rows:        rows,
		GeneratedAt: time.Now(),
		Filters: map[string]string{
			"stage":            opts.Stage,
			"minAppearances":   strconv.Itoa(opts.MinAppearances),
			"minOpportunities": strconv.Itoa(opts.MinOpportunities),
		},
	}, nil
}

type objectiveMeta struct {
	ID    uint
	Name  string
	Stage string
	Phase string
}

func loadObjectives(ctx context.Context, db *gorm.DB, opts ObjectiveDifficultyOptions) ([]objectiveMeta, error) {
	q := db.WithContext(ctx).Table("objectives").Select("id, name, stage, phase")
	if opts.Stage != "" && opts.Stage != "all" {
		q = q.Where("stage = ?", opts.Stage)
	} else {
		q = q.Where("stage IN ('I','II')")
	}
	var rows []objectiveMeta
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

type appOppRow struct{ Appearances, Opportunities int64 }

func loadAppearancesAndOpportunities(ctx context.Context, db *gorm.DB, opts ObjectiveDifficultyOptions) (map[uint]appOppRow, error) {
	q := db.WithContext(ctx).
		Table("game_objectives AS go").
		Select(`
            go.objective_id AS objective_id,
            COUNT(DISTINCT go.game_id) AS appearances,
            COUNT(gp.id) AS opportunities
        `).
		Joins("JOIN objectives o ON o.id = go.objective_id").
		Joins("JOIN game_players gp ON gp.game_id = go.game_id").
		Where("go.revealed = ?", true)

	if opts.Stage != "" && opts.Stage != "all" {
		q = q.Where("o.stage = ?", opts.Stage)
	} else {
		q = q.Where("o.stage IN ('I','II')")
	}

	type row struct {
		ObjectiveID   uint
		Appearances   int64
		Opportunities int64
	}
	var out []row
	if err := q.Group("go.objective_id").Scan(&out).Error; err != nil {
		return nil, err
	}

	m := make(map[uint]appOppRow, len(out))
	for _, r := range out {
		m[r.ObjectiveID] = appOppRow{Appearances: r.Appearances, Opportunities: r.Opportunities}
	}
	return m, nil
}

func loadScoresAndFirstRounds(ctx context.Context, db *gorm.DB, opts ObjectiveDifficultyOptions) (map[uint]int64, map[uint][]int, error) {
	base := db.WithContext(ctx).Table("scores AS s").
		Select("s.objective_id, s.game_id, s.player_id, MIN(r.number) AS first_round").
		Joins("JOIN rounds r ON r.id = s.round_id").
		Joins("JOIN objectives o ON o.id = s.objective_id").
		Where("s.objective_id <> 0").
		Where("LOWER(TRIM(s.type)) IN ?", []string{"public", "objective"})

	if opts.Stage != "" && opts.Stage != "all" {
		base = base.Where("o.stage = ?", opts.Stage)
	} else {
		base = base.Where("o.stage IN ('I','II')")
	}

	type row struct {
		ObjectiveID uint
		GameID      uint
		PlayerID    uint
		FirstRound  int
	}
	var rows []row
	if err := base.Group("s.objective_id, s.game_id, s.player_id").Scan(&rows).Error; err != nil {
		return nil, nil, err
	}

	scores := make(map[uint]int64)
	first := make(map[uint][]int)
	for _, r := range rows {
		scores[r.ObjectiveID]++
		first[r.ObjectiveID] = append(first[r.ObjectiveID], r.FirstRound)
	}
	return scores, first, nil
}

func bayesAdj(S, O, alpha, beta int) float64 {
	return float64(S+alpha) / float64(O+alpha+beta)
}

func wilsonInterval(S, O int, conf float64) (float64, float64) {
	if O == 0 {
		return 0, 0
	}
	z := 1.959963984540054
	phat := float64(S) / float64(O)
	denom := 1 + (z*z)/float64(O)
	center := phat + (z*z)/(2*float64(O))
	margin := z * math.Sqrt((phat*(1-phat)+(z*z)/(4*float64(O)))/float64(O))
	lo := (center - margin) / denom
	hi := (center + margin) / denom
	if lo < 0 {
		lo = 0
	}
	if hi > 1 {
		hi = 1
	}
	return lo, hi
}

func avgAndMedian(xs []int) (avg float64, median float64) {
	if len(xs) == 0 {
		return 0, 0
	}
	sum := 0
	for _, v := range xs {
		sum += v
	}
	avg = float64(sum) / float64(len(xs))
	ys := append([]int(nil), xs...)
	sort.Ints(ys)
	n := len(ys)
	if n%2 == 1 {
		median = float64(ys[n/2])
	} else {
		median = float64(ys[n/2-1]+ys[n/2]) / 2.0
	}
	return
}
