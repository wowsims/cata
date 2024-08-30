package affliction

import (
	"github.com/wowsims/cata/sim/common/cata"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	cata.CreateDTRClassConfig(proto.Spec_SpecAfflictionWarlock, 0.085).
		AddSpell(348, cata.NewDragonwrathSpellConfig().SupressImpact()).  // Immolate
		AddSpell(47897, cata.NewDragonwrathSpellConfig().SupressImpact()) // Shadowflame
}
