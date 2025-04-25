package tbc

import (
	"time"

	"github.com/wowsims/mop/sim/common/shared"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func init() {
	core.AddEffectsToTest = false

	// Band of the Eternal Restorer
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Band of the Eternal Restorer",
		ItemID:     29309,
		AuraID:     35087,
		Bonus:      stats.Stats{stats.SpellPower: 93},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})

	// Shattered Sun Pendant of Restoration
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Light's Salvation",
		ItemID:     34677,
		AuraID:     45478,
		Bonus:      stats.Stats{stats.SpellPower: 117},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellHealing,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})

}
