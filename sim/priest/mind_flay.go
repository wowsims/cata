package priest

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (priest *Priest) newMindFlaySpell() *core.Spell {
	mindFlayCoefficient := 0.19799999893
	return priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 15407},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagChanneled | core.SpellFlagAPL,
		ClassSpellMask: PriestSpellMindFlay,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 8,
			PercentModifier: 100,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,
		CritMultiplier:           priest.DefaultCritMultiplier(),
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "MindFlay-" + priest.Label,
			},
			NumberOfTicks:        3,
			TickLength:           time.Second * 1,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: true,
			BonusCoefficient:     0.2879999876,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, priest.CalcScalingSpellDmg(mindFlayCoefficient))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
				spell.DealOutcome(sim, result)
			}
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			return spell.CalcPeriodicDamage(sim, target, priest.CalcScalingSpellDmg(mindFlayCoefficient), spell.OutcomeExpectedMagicCrit)
		},
	})
}
