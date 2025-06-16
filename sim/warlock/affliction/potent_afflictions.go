package affliction

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

func (affliction *AfflictionWarlock) registerPotentAffliction() {
	dmgMod := affliction.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: affliction.getMasteryBonus() / 100,
		ClassMask:  warlock.WarlockSpellAgony | warlock.WarlockSpellUnstableAffliction | warlock.WarlockSpellCorruption,
	})

	affliction.AddOnMasteryStatChanged(func(_ *core.Simulation, _, _ float64) {
		dmgMod.UpdateFloatValue(affliction.getMasteryBonus() / 100)
	})

	dmgMod.Activate()
}
