package destruction

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

func (destruction *DestructionWarlock) registerFireAndBrimstone() {
	destruction.FABAura = destruction.RegisterAura(core.Aura{
		Label:    "Fire and Brimstone",
		ActionID: core.ActionID{SpellID: 108683},
		Duration: core.NeverExpires,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !destruction.BurningEmbers.CanSpend(10) && aura.IsActive() {
				aura.Deactivate(sim)
			}
		},
	})

	destruction.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 108683},
		SpellSchool:      core.SpellSchoolFire,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		ProcMask:         core.ProcMaskEmpty,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{NonEmpty: true},
		},

		ClassSpellMask: warlock.WarlockSpellFireAndBrimstone,
		Flags:          core.SpellFlagAPL,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			destruction.FABAura.Activate(sim)
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return destruction.BurningEmbers.CanSpend(10)
		},
	})

	destruction.registerFireAndBrimstoneConflagrate()
	destruction.registerFireAndBrimstoneImmolate()
	destruction.registerFireAndBrimstoneIncinerate()
}

func (destruction *DestructionWarlock) getFABReduction() float64 {
	return 0.35 * (1 + destruction.getSpenderMasteryBonus())
}
