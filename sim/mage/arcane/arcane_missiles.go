package arcane

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/mage"
)

func (arcane *ArcaneMage) OutcomeArcaneMissiles(sim *core.Simulation, result *core.SpellResult, attackTable *core.AttackTable) {
	spell := arcane.arcaneMissilesTickSpell
	if spell.MagicHitCheck(sim, attackTable) {
		if sim.RandomFloat("Magical Crit Roll") < arcane.arcaneMissileCritSnapshot {
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

func (arcane *ArcaneMage) registerArcaneMissilesSpell() {

	// Values found at https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=exact%253A7268
	arcaneMissilesScaling := 0.22
	arcaneMissilesCoefficient := 0.22
	arcaneMissilesTickSpell := arcane.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 7268},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: mage.MageSpellArcaneMissilesTick,
		MissileSpeed:   20,

		DamageMultiplier: 1,
		CritMultiplier:   arcane.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: arcaneMissilesCoefficient,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := arcane.CalcScalingSpellDmg(arcaneMissilesScaling)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, arcane.OutcomeArcaneMissiles)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})

	arcane.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 7268},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagChanneled | core.SpellFlagAPL,
		ClassSpellMask: mage.MageSpellArcaneMissilesCast,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return arcane.arcaneMissilesProcAura.IsActive()
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "ArcaneMissiles",
			},
			NumberOfTicks:        5,
			TickLength:           time.Millisecond * 400,
			HasteReducesDuration: true,
			AffectedByCastSpeed:  true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				arcaneMissilesTickSpell.Cast(sim, target)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				// Snapshot crit chance
				arcane.arcaneMissileCritSnapshot = arcaneMissilesTickSpell.SpellCritChance(target)
				spell.Dot(target).Apply(sim)
			}
		},
	})

	// Aura for when proc is successful
	arcane.arcaneMissilesProcAura = arcane.RegisterAura(core.Aura{
		Label:    "Arcane Missiles Proc",
		ActionID: core.ActionID{SpellID: 79683},
		Duration: time.Second * 20,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(mage.MageSpellArcaneMissilesCast) {
				aura.Deactivate(sim)
			}
		},
	})

	// Listener for procs
	core.MakePermanent(arcane.RegisterAura(core.Aura{
		Label: "Arcane Missiles Activation",
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(mage.MageSpellsAllDamaging) {
				if sim.Proc(0.3, "Arcane Missiles") {
					arcane.arcaneMissilesProcAura.Activate(sim)
				}
			}
		},
	}))
}
