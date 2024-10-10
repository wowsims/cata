package demonology

import (
	"github.com/wowsims/cata/sim/common/cata"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	cata.CreateDTRClassConfig(proto.Spec_SpecDemonologyWarlock, 0.12).
		AddSpell(348, cata.NewDragonwrathSpellConfig().SupressImpact()).  // Immolate
		AddSpell(47897, cata.NewDragonwrathSpellConfig().SupressImpact()) // Shadowflame
}
