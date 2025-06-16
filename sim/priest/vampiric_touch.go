package priest

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

const VtScaleCoeff = 0.071
const VtSpellCoeff = 0.415

func (priest *Priest) registerVampiricTouchSpell() {
	manaMetric := priest.NewManaMetrics(core.ActionID{SpellID: 34914}.WithTag(1))
	priest.VampiricTouch = priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 34914},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: PriestSpellVampiricTouch,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 3,
		},

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           priest.DefaultCritMultiplier(),
		ThreatMultiplier:         1,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "VampiricTouch",
			},
			NumberOfTicks:       5,
			TickLength:          time.Second * 3,
			AffectedByCastSpeed: true,

			BonusCoefficient: VtSpellCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.Snapshot(target, priest.CalcScalingSpellDmg(VtScaleCoeff))
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				priest.AddMana(sim, priest.MaxMana()*0.02, manaMetric)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
				spell.DealOutcome(sim, result)
			}
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, dot.OutcomeExpectedMagicSnapshotCrit)
			} else {
				return spell.CalcPeriodicDamage(sim, target, priest.CalcScalingSpellDmg(VtScaleCoeff), spell.OutcomeExpectedMagicCrit)
			}
		},
	})
}
