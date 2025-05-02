package balance

import (
	cata "github.com/wowsims/mop/sim/common/cata"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	// https://docs.google.com/spreadsheets/d/e/2PACX-1vTaCACFb7dqXpF2qwAZIAgXX-p2VTuJWqmyWXqaJ3c49FNWm61E9-unEdN3cn7YHevoGWWPmqkqJv6h/pubhtml
	cata.CreateDTRClassConfig(proto.Spec_SpecBalanceDruid, 0.08).
		AddSpell(42231, cata.NewDragonwrathSpellConfig().IsAoESpell()).
		AddSpell(78777, cata.NewDragonwrathSpellConfig().IsAoESpell())
}
