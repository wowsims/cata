package tbc

import (
	"github.com/wowsims/cata/sim/common/shared"
	"github.com/wowsims/cata/sim/core"
)

func init() {
	core.AddEffectsToTest = false

	// Band of the Eternal Restorer
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Band of the Eternal Restorer",
		ItemID:   29309,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
	})

	// Shattered Sun Pendant of Restoration
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "Light's Salvation",
		ItemID:   34677,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellHealing,
	})

}
