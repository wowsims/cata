package demonology

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

func (demonology *DemonologyWarlock) registerImmolationAura() {
	demonology.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 104025},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellHellfire,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           demonology.DefaultCritMultiplier(),
		ThreatMultiplier:         1,

		Hot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Immolation Aura (DoT)",
				ActionID: core.ActionID{SpellID: 104025}.WithTag(1),
			},

			TickLength:           time.Second,
			NumberOfTicks:        8,
			HasteReducesDuration: true,
			IsAOE:                true,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if !demonology.DemonicFury.CanSpend(core.TernaryInt32(demonology.T15_2pc.IsActive(), 18, 25)) {
					dot.Deactivate(sim)
					return
				}

				demonology.DemonicFury.Spend(core.TernaryInt32(demonology.T15_2pc.IsActive(), 18, 25), dot.Spell.ActionID, sim)
				demonology.Hellfire.RelatedDotSpell.Cast(sim, target)
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return demonology.IsInMeta() && demonology.DemonicFury.CanSpend(core.TernaryInt32(demonology.T15_2pc.IsActive(), 18, 25))
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Hot(&demonology.Unit).Apply(sim)
		},
	})
}
