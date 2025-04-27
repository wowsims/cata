package affliction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

func (affliction *AfflictionWarlock) registerUnstableAffliction() {
	affliction.UnstableAffliction = affliction.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 30108},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagHauntSE | core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellUnstableAffliction,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 15},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 1500 * time.Millisecond,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           affliction.DefaultCritMultiplier(),
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura:                core.Aura{Label: "UnstableAffliction"},
			NumberOfTicks:       5,
			TickLength:          3 * time.Second,
			AffectedByCastSpeed: true,
			BonusCoefficient:    0.20000000298,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.Snapshot(target, affliction.CalcScalingSpellDmg(0.23199999332))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
				affliction.Immolate.Dot(target).Deactivate(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
