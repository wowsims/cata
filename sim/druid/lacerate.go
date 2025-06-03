package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (druid *Druid) registerLacerateSpell() {
	druid.Lacerate = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 33745},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},

			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 3,
			},

			IgnoreHaste: true,
		},

		BonusCritPercent: 0,
		DamageMultiplier: 1,
		CritMultiplier:   druid.DefaultCritMultiplier(),
		ThreatMultiplier: 1, // Changed in Cata
		MaxRange:         core.MaxMeleeRange,
		FlatThreatBonus:  0, // Removed in Cata

		Dot: core.DotConfig{
			Aura: druid.applyRendAndTear(core.Aura{
				Label:     "Lacerate",
				MaxStacks: 3,
				Duration:  time.Second * 15,
			}),
			NumberOfTicks: 5,
			TickLength:    time.Second * 3,

			OnSnapshot: func(_ *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				if isRollover {
					panic("Lacerate cannot roll over snapshots!")
				}

				dot.SnapshotPhysical(target, 0.0512*dot.Spell.MeleeAttackPower()*float64(dot.Aura.GetStacks()))
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.616 * spell.MeleeAttackPower()
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				dot := spell.Dot(target)
				if dot.IsActive() {
					dot.Refresh(sim)
					dot.AddStack(sim)
					dot.TakeSnapshot(sim, false)
				} else {
					dot.Apply(sim)
					dot.SetStacks(sim, 1)
					dot.TakeSnapshot(sim, false)
				}

				if sim.Proc(0.25, "Mangle CD Reset") {
					druid.MangleBear.CD.Reset()
				}
			}
		},
	})
}
