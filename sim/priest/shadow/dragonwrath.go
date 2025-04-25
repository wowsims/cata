package shadow

import (
	"github.com/wowsims/mop/sim/common/mop"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	mop.CreateDTRClassConfig(proto.Spec_SpecShadowPriest, 0.09).
		AddSpell(2944, mop.NewDragonwrathSpellConfig().SupressImpact()).                 // Improved Devouring Plague
		AddSpell(48045, mop.NewDragonwrathSpellConfig().IsAoESpell().TreatTickAsCast()). // Mind sear
		AddSpell(87532, mop.NewDragonwrathSpellConfig().SupressSpell())                  // Shadowy Apparition
}
