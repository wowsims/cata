package destruction

import (
	"github.com/wowsims/mop/sim/core"
)

func (destruction *DestructionWarlock) ApplyMastery() {

	spenderMod := destruction.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: destruction.getSpenderMasteryBonus(),
		ClassMask:  SpellMaskCinderSpender,
	})

	generatorMod := destruction.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: destruction.getGeneratorMasteryBonus(),
		ClassMask:  SpellMaskCinderGenerator,
	})

	destruction.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMasteryRating, newMasteryRating float64) {
		generatorMod.UpdateFloatValue(destruction.getGeneratorMasteryBonus())
		spenderMod.UpdateFloatValue(destruction.getSpenderMasteryBonus())
	})

	core.MakePermanent(destruction.RegisterAura(core.Aura{
		Label: "Mastery: Emberstorm",
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			generatorMod.Activate()
			spenderMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			generatorMod.Deactivate()
			spenderMod.Deactivate()
		},
	}))
}
