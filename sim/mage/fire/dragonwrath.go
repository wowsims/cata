package fire

import (
	"github.com/wowsims/cata/sim/common/cata"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	cata.CreateDTRClassConfig(proto.Spec_SpecFireMage, 0.08).
		AddSpell(88148, cata.NewDragonwrathSpellConfig().SupressImpact()). // Flamestrike (Blast Wave)
		AddSpell(11113, cata.NewDragonwrathSpellConfig().ProcPerCast()).
		AddSpell(83619, cata.NewDragonwrathSpellConfig().SupressSpell()) // Fire Power
}
