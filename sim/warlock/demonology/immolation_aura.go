package demonology

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

const immolationAuraScale = 0.17499999702 * 1.25 // 2025.06.13 Changes to Beta - Immolation Aura damage increased by 25%
const immolationAuraCoeff = 0.17499999702 * 1.25

func (demonology *DemonologyWarlock) registerImmolationAura() {
	var baseDamage = demonology.CalcScalingSpellDmg(immolationAuraScale)

	immolationAura := demonology.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 104025},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAoE | core.SpellFlagAPL | core.SpellFlagNoMetrics,
		ClassSpellMask: warlock.WarlockSpellImmolationAura,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           demonology.DefaultCritMultiplier(),
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Immolation Aura (DoT)",
			},

			TickLength:           time.Second,
			NumberOfTicks:        8,
			HasteReducesDuration: true,
			BonusCoefficient:     immolationAuraCoeff,
			IsAOE:                true,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if !demonology.DemonicFury.CanSpend(core.TernaryInt32(demonology.T15_2pc.IsActive(), 18, 25)) {
					dot.Deactivate(sim)
					return
				}

				demonology.DemonicFury.Spend(sim, core.TernaryInt32(demonology.T15_2pc.IsActive(), 18, 25), dot.Spell.ActionID)

				for _, unit := range sim.Encounter.TargetUnits {
					dot.Spell.CalcAndDealPeriodicDamage(sim, unit, baseDamage, dot.OutcomeTick)
				}
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return demonology.IsInMeta() && demonology.DemonicFury.CanSpend(core.TernaryInt32(demonology.T15_2pc.IsActive(), 18, 25))
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})

	demonology.Metamorphosis.RelatedSelfBuff.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
		immolationAura.AOEDot().Deactivate(sim)
	})
}
