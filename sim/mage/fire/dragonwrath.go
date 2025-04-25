package fire

import (
	cata "github.com/wowsims/mop/sim/common/cata"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	cata.CreateDTRClassConfig(proto.Spec_SpecFireMage, 0.08).
		AddSpell(83619, cata.NewDragonwrathSpellConfig().SupressSpell()).              // Fire Power
		AddSpell(2120, cata.NewDragonwrathSpellConfig().IsAoESpell()).                 // Flame Strike
		AddSpell(11113, cata.NewDragonwrathSpellConfig().IsAoESpell()).                // Blast Wave
		AddSpell(88148, cata.NewDragonwrathSpellConfig().IsAoESpell().SupressImpact()) // Improved Flame Strike
}
