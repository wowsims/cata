package shadow

import (
	cata "github.com/wowsims/mop/sim/common/cata"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	cata.CreateDTRClassConfig(proto.Spec_SpecShadowPriest, 0.09).
		AddSpell(2944, cata.NewDragonwrathSpellConfig().SupressImpact()).                 // Improved Devouring Plague
		AddSpell(48045, cata.NewDragonwrathSpellConfig().IsAoESpell().TreatTickAsCast()). // Mind sear
		AddSpell(87532, cata.NewDragonwrathSpellConfig().SupressSpell())                  // Shadowy Apparition
}
