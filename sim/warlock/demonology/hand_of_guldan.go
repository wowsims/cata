package demonology

import (
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
			NumberOfTicks: 6,
			TickLength:    time.Second,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, demonology.CalcScalingSpellDmg(shadowFlameScale))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				demonology.DemonicFury.Gain(2, dot.Spell.ActionID, sim)
			},
		},
	})

	demonology.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 105174},
		SpellSchool:    core.SpellSchoolFire | core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
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
