package tbc

import (
	"time"

	"github.com/wowsims/cata/sim/common/shared"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func init() {

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

	core.AddEffectsToTest = true
}
