package subtlety

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

func (subRogue *SubtletyRogue) applyPassives() {
	// Sanguinary Vein - 50% Increase to Rupture
	subRogue.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  rogue.RogueSpellRupture,
		FloatValue: 0.5,
	})

	// Apply Mastery
	masteryMod := subRogue.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  rogue.RogueSpellDamagingFinisher,
		FloatValue: subRogue.GetMasteryBonus(),
	})

	subRogue.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		masteryMod.UpdateFloatValue(subRogue.GetMasteryBonus())
	})

	core.MakePermanent(subRogue.GetOrRegisterAura(core.Aura{
		Label:    "Executioner",
		ActionID: core.ActionID{SpellID: 76808},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			masteryMod.UpdateFloatValue(subRogue.GetMasteryBonus())
			masteryMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			masteryMod.Deactivate()
		},
	}))
}
