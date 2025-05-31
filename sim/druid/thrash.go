package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (druid *Druid) registerThrashBearSpell() {
	flatHitDamage := 1.125 * druid.ClassSpellScaling          // ~1232
	flatTickDamage := 0.62699997425 * druid.ClassSpellScaling // ~686

	druid.ThrashBear = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 77758},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreArmor | core.SpellFlagAPL | core.SpellFlagAoE,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   druid.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: druid.applyRendAndTear(core.Aura{
				Label:    "Thrash (Bear)",
				Duration: time.Second * 16,
			}),
			NumberOfTicks: 8,
			TickLength:    time.Second * 2,

			OnSnapshot: func(_ *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				if isRollover {
					panic("Thrash cannot roll-over snapshots!")
				}

				dot.SnapshotPhysical(target, flatTickDamage+0.141*dot.Spell.MeleeAttackPower())
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			baseDamage := flatHitDamage + 0.191*spell.MeleeAttackPower()

			for _, aoeTarget := range sim.Encounter.TargetUnits {
				result := spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

				if result.Landed() {
					spell.Dot(aoeTarget).Apply(sim)
					druid.WeakenedBlowsAuras.Get(aoeTarget).Activate(sim)

					if sim.Proc(0.25, "Mangle CD Reset") {
						druid.MangleBear.CD.Reset()
					}
				}
			}
		},

		RelatedAuraArrays: druid.WeakenedBlowsAuras.ToMap(),
	})
}

func (druid *Druid) registerThrashCatSpell() {
	flatHitDamage := 1.125 * druid.ClassSpellScaling          // ~1232
	flatTickDamage := 0.62699997425 * druid.ClassSpellScaling // ~686

	druid.ThrashCat = druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 106830},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreArmor | core.SpellFlagAPL | core.SpellFlagAoE,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		EnergyCost: core.EnergyCostOptions{
			Cost: 50,
		},

		DamageMultiplier: 1,
		CritMultiplier:   druid.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: druid.applyRendAndTear(core.Aura{
				Label:    "Thrash (Cat)",
				Duration: time.Second * 15,
			}),
			NumberOfTicks: 5,
			TickLength:    time.Second * 3,

			OnSnapshot: func(_ *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				if isRollover {
					panic("Thrash cannot roll-over snapshots!")
				}

				dot.SnapshotPhysical(target, flatTickDamage+0.141*dot.Spell.MeleeAttackPower())
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			baseDamage := flatHitDamage + 0.191*spell.MeleeAttackPower()

			for _, aoeTarget := range sim.Encounter.TargetUnits {
				result := spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

				if result.Landed() {
					spell.Dot(aoeTarget).Apply(sim)
					druid.WeakenedBlowsAuras.Get(aoeTarget).Activate(sim)
				}
			}
		},

		RelatedAuraArrays: druid.WeakenedBlowsAuras.ToMap(),
	})
}
