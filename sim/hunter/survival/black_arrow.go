package survival

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (hunter *SurvivalHunter) registerBlackArrowSpell(timer *core.Timer) {
	if !hunter.Talents.BlackArrow {
		return
	}

	actionID := core.ActionID{SpellID: 3674}

	hunter.Hunter.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskRangedSpecial,
		
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		FocusCost: core.FocusCostOptions{
			Cost: 35,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second*30 - time.Second*2*time.Duration(hunter.Talents.Resourcefulness),
			},
		},

		DamageMultiplierAdditive: 1 +
			.10*float64(hunter.Talents.TrapMastery),
		DamageMultiplier: 1 *
			(1.0 / 1.06), // Black Arrow is not affected by its own 1.06 aura.
		ThreatMultiplier: 1,
		CritMultiplier: (0.5) * ( 1 + float64(hunter.Talents.Toxicology) * 0.5), //Todo: SimC, is this crit damage multiplier?
		

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "BlackArrow-3674",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					hunter.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier *= 1.06
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					hunter.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier /= 1.06
				},
			},
			NumberOfTicks: 10,
			TickLength:    time.Second * 2,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				// scales slightly better (11.5%) than the tooltip implies (10%), but isn't affected by Hunter's Mark
				dot.SnapshotBaseDamage = 2849 + 0.665 * (dot.Spell.Unit.GetStat(stats.RangedAttackPower)+dot.Spell.Unit.PseudoStats.MobTypeAttackPower)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeRangedHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
