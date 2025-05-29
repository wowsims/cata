package demonology

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

const shadowFlameScale = 0.137
const shadowFlameCoeff = 0.137
const hogScale = 0.575
const hogCoeff = 0.575

func (demonology *DemonologyWarlock) registerHandOfGuldan() {
	shadowFlame := demonology.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 47960},
		SpellSchool:    core.SpellSchoolFire | core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagNoOnCastComplete,
		ClassSpellMask: warlock.WarlockSpellHandOfGuldan,

		ThreatMultiplier: 1,
		DamageMultiplier: 1,
		CritMultiplier:   demonology.DefaultCritMultiplier(),
		BonusCoefficient: shadowFlameCoeff,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Shadowflame",
				MaxStacks: 2,
			},
			NumberOfTicks:    6,
			TickLength:       time.Second,
			BonusCoefficient: shadowFlameCoeff,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, 0)
				stacks := math.Min(float64(dot.Aura.GetStacks())+1, 2)
				dot.SnapshotBaseDamage = demonology.CalcScalingSpellDmg(shadowFlameScale) + stacks*dot.BonusCoefficient*dot.Spell.BonusDamage()
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				demonology.DemonicFury.Gain(2, dot.Spell.ActionID, sim)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Activate(sim)
			spell.Dot(target).Aura.AddStack(sim)
		},
	})

	demonology.HandOfGuldan = demonology.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 105174},
		SpellSchool:    core.SpellSchoolFire | core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAoE | core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellHandOfGuldan,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 5},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		Charges:      2,
		RechargeTime: time.Second * 15,

		DamageMultiplier: 1,
		CritMultiplier:   demonology.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: hogCoeff,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !demonology.IsInMeta()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// keep stacks in sync as they're shared
			demonology.ChaosWave.ConsumeCharge(sim)

			// Shadowflame is snapshotted at cast time
			for _, enemy := range sim.Encounter.TargetUnits {
				shadowFlame.Dot(enemy).TakeSnapshot(sim, false)
			}

			sim.AddPendingAction(&core.PendingAction{
				NextActionAt: sim.CurrentTime + time.Millisecond*1300, // Fixed delay of 1.3 seconds
				Priority:     core.ActionPriorityAuto,
				OnAction: func(sim *core.Simulation) {
					for _, enemy := range sim.Encounter.TargetUnits {
						result := spell.CalcAndDealDamage(
							sim,
							enemy,
							demonology.CalcScalingSpellDmg(hogScale),
							spell.OutcomeMagicHitAndCrit,
						)

						if result.Landed() {
							shadowFlame.Cast(sim, enemy)
						}
					}
				},
			})
		},
	})
}
