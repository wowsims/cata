package balance

import (
	"github.com/wowsims/cata/sim/common/cata"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	cata.CreateDTRClassConfig(proto.Spec_SpecBalanceDruid, 0.072).
		AddSpell(8921, cata.NewDragonwrathSpellConfig().SupressImpact()). // Moonfire
		AddSpell(93402, cata.NewDragonwrathSpellConfig().SupressImpact()) // Sunfire
}
