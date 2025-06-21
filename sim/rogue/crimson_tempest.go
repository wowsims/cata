package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (rogue *Rogue) registerCrimsonTempest() {
	hit_baseDamage := rogue.GetBaseDamageFromCoefficient(0.47600001097)
	hit_cpScaling := rogue.GetBaseDamageFromCoefficient(0.47600001097) // Same number...?
	hit_apScaling := 0.0275

	hit_minDamage := hit_baseDamage * 0.5

	dot_modifier := 2.4
	var lastCTDamage []float64

	// The DoT does not benefit from any external buff/debuff/passive
	rogue.CrimsonTempestDoT = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 121411, Tag: 7},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagIgnoreTargetModifiers | core.SpellFlagIgnoreAttackerModifiers | core.SpellFlagPassiveSpell,
		ClassSpellMask: RogueSpellCrimsonTempestDoT,

		DamageMultiplier: 1,
		CritMultiplier:   rogue.CritMultiplier(false),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Crimson Tempest",
				Tag:   RogueBleedTag,
			},
			NumberOfTicks: 6,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotPhysical(target, lastCTDamage[target.UnitIndex]/6)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHitNoHitCounter)
		},
	})

	rogue.CrimsonTempest = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 121411},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagFinisher | core.SpellFlagAPL,
		MetricSplits:   6,
		ClassSpellMask: RogueSpellCrimsonTempest,

		DamageMultiplier: 1,
		CritMultiplier:   rogue.CritMultiplier(false),
		ThreatMultiplier: 1,

		EnergyCost: core.EnergyCostOptions{
			Cost:          35,
			Refund:        0.8,
			RefundMetrics: rogue.EnergyRefundMetrics,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
				// Omitting the GCDMin - does not appear affected by either Shadow Blades or Adrenaline Rush
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(rogue.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			lastCTDamage = make([]float64, sim.GetNumTargets())
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				damage := hit_minDamage +
					sim.RandomFloat("Crimson Tempest")*hit_baseDamage +
					hit_cpScaling*float64(rogue.ComboPoints()) +
					hit_apScaling*float64(rogue.ComboPoints())*spell.MeleeAttackPower()

				result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeSpecialHitAndCrit)
				lastCTDamage[aoeTarget.UnitIndex] = result.Damage * dot_modifier

				if result.Landed() {
					rogue.CrimsonTempestDoT.Cast(sim, aoeTarget)
					// Currently, CT is triggering a Relentless Strikes refund for _every single target hit_
					// I'm assuming this to be a bug currently, but will model it should it stay for some time
				}
			}
			rogue.ApplyFinisher(sim, spell)
		},
	})
}
