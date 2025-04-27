package priest

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (priest *Priest) registerDevouringPlagueSpell() {
	actionID := core.ActionID{SpellID: 2944, Tag: 0}
	priest.DevouringPlague = priest.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagDisease | core.SpellFlagAPL,
		ClassSpellMask: PriestSpellDevouringPlague,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 25,
			PercentModifier: 100,
		},
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,
		CritMultiplier:           priest.DefaultCritMultiplier(),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "DevouringPlague",
			},

			NumberOfTicks:       8,
			TickLength:          time.Second * 3,
			AffectedByCastSpeed: true,

			BonusCoefficient: 0.163,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.Snapshot(target, 0.144*priest.ClassSpellScaling)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}

			spell.DealOutcome(sim, result)
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, dot.OutcomeExpectedMagicSnapshotCrit)
			} else {
				baseDamage := 0.144 * priest.ClassSpellScaling
				return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicCrit)
			}
		},
	})
}
