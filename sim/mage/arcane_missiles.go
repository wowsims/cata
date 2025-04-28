package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) OutcomeArcaneMissiles(sim *core.Simulation, result *core.SpellResult, attackTable *core.AttackTable) {
	spell := mage.arcaneMissilesTickSpell
	if spell.MagicHitCheck(sim, attackTable) {
		if sim.RandomFloat("Magical Crit Roll") < mage.arcaneMissileCritSnapshot {
			result.Outcome = core.OutcomeCrit
			result.Damage *= spell.CritMultiplier
			spell.SpellMetrics[result.Target.UnitIndex].Crits++
		} else {
			result.Outcome = core.OutcomeHit
			spell.SpellMetrics[result.Target.UnitIndex].Hits++
		}
	} else {
		result.Outcome = core.OutcomeMiss
		result.Damage = 0
		spell.SpellMetrics[result.Target.UnitIndex].Misses++
	}
}

func (mage *Mage) registerArcaneMissilesSpell() {
	mage.arcaneMissilesTickSpell = mage.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 7268},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: MageSpellArcaneMissilesTick,
		MissileSpeed:   20,

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 0.278,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := 0.432 * mage.ClassSpellScaling
			result := spell.CalcDamage(sim, target, damage, mage.OutcomeArcaneMissiles)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 7268},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagChanneled | core.SpellFlagAPL,
		ClassSpellMask: MageSpellArcaneMissilesCast,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return mage.arcaneMissilesProcAura.IsActive()
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "ArcaneMissiles" + mage.Label,
			},
			NumberOfTicks:        3,
			TickLength:           time.Millisecond * 700,
			HasteReducesDuration: true,
			AffectedByCastSpeed:  true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				mage.arcaneMissilesTickSpell.Cast(sim, target)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				// Snapshot crit chance
				mage.arcaneMissileCritSnapshot = mage.arcaneMissilesTickSpell.SpellCritChance(target)
				dot := spell.Dot(target)
				dot.Apply(sim)
			}
		},
	})
}
