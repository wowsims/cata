package fire

import (
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	mop.CreateDTRClassConfig(proto.Spec_SpecFireMage, 0.08).
		AddSpell(83619, mop.NewDragonwrathSpellConfig().SupressSpell()).              // Fire Power
		AddSpell(2120, mop.NewDragonwrathSpellConfig().IsAoESpell()).                 // Flame Strike
		AddSpell(11113, mop.NewDragonwrathSpellConfig().IsAoESpell()).                // Blast Wave
		AddSpell(88148, mop.NewDragonwrathSpellConfig().IsAoESpell().SupressImpact()) // Improved Flame Strike
}
