package factions

import "strings"

var AllFactions = []string{
	"Arborec",
	"Argent Flight",
	"Barony of Letnev",
	"Clan of Saar",
	"Embers of Muaat",
	"Empyrean",
	"Federation of Sol",
	"Ghosts of Creuss",
	"Hakann",
	"L1Z1X Mindnet",
	"Mahact Gene-Sorcerers",
	"Mentak Coalition",
	"Naalu Collective",
	"Naaz-Rokha Alliance",
	"Nekro Virus",
	"Nomad",
	"Sardakk N'orr",
	"Titans of Ul",
	"Universities of Jol-Nar",
	"Vuil'raith Cabal",
	"Winnu",
	"Xxcha Kingdom",
	"Yin Brotherhood",
	"Yssaril Tribes",
	"Council Keleres",
}

func IsValidFaction(name string) bool {
	for _, f := range AllFactions {
		if strings.EqualFold(f, name) {
			return true
		}
	}
	return false
}
