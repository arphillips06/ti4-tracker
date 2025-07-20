package stats

import (
	"fmt"

	"github.com/arphillips06/TI4-stats/models"
)

func FormatVictoryPathKey(vp models.VictoryPath) string {
	return fmt.Sprintf("S1:%d S2:%d Sec:%d Cust:%d Imp:%d Rel:%d Ag:%d AC:%d Sup:%d",
		vp.Stage1Points, vp.Stage2Scored, vp.SecretPoints,
		vp.Custodians, vp.Imperial,
		vp.Relics, vp.Agenda, vp.ActionCard, vp.Support)
}
