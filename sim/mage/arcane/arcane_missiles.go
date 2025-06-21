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
			spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
		} else {
			result.Outcome = core.OutcomeHit
			spell.SpellMetrics[result.Target.UnitIndex].Ticks++
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
	actionID := core.ActionID{SpellID: 7268}
	arcane.arcaneMissilesTickSpell = arcane.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(1),
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: mage.MageSpellArcaneMissilesTick,
		MissileSpeed:   20,

		DamageMultiplier: 1,
		CritMultiplier:   arcane.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: arcaneMissilesCoefficient,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := arcane.CalcAndRollDamageRange(sim, arcaneMissilesScaling, 0)
			result := spell.CalcDamage(sim, target, baseDamage, arcane.OutcomeArcaneMissiles)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealPeriodicDamage(sim, result)
			})
			spell.SpellMetrics[result.Target.UnitIndex].Casts--
		},
	})

	arcane.arcaneMissiles = arcane.RegisterSpell(core.SpellConfig{
		ActionID:         actionID, // Real SpellID: 5143
		SpellSchool:      core.SpellSchoolArcane,
		ProcMask:         core.ProcMaskSpellDamage,
		Flags:            core.SpellFlagChanneled | core.SpellFlagAPL,
		ClassSpellMask:   mage.MageSpellArcaneMissilesCast,
		DamageMultiplier: 0,

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
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					arcane.ArcaneChargesAura.Activate(sim)
					arcane.ArcaneChargesAura.AddStack(sim)
				},
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					arcane.ArcaneChargesAura.Deactivate(sim)
				},
			},
			NumberOfTicks:        5,
			TickLength:           time.Millisecond * 400,
			HasteReducesDuration: true,
			AffectedByCastSpeed:  true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				arcane.arcaneMissilesTickSpell.Cast(sim, target)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if arcane.arcaneMissilesProcAura.IsActive() {
				arcane.arcaneMissilesProcAura.RemoveStack(sim)
			}
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				// Snapshot crit chance
				arcane.arcaneMissileCritSnapshot = arcane.arcaneMissilesTickSpell.SpellCritChance(target)
				spell.Dot(target).Apply(sim)
				arcane.arcaneMissilesTickSpell.SpellMetrics[target.UnitIndex].Hits++
				arcane.arcaneMissilesTickSpell.SpellMetrics[target.UnitIndex].Casts++
			}
		},
	})

	// Aura for when proc is successful
	arcane.arcaneMissilesProcAura = arcane.RegisterAura(core.Aura{
		Label:     "Arcane Missiles Proc",
		ActionID:  core.ActionID{SpellID: 79683},
		Duration:  time.Second * 20,
		MaxStacks: 2,
	})

	// Listener for procs
	core.MakeProcTriggerAura(&arcane.Unit, core.ProcTrigger{
		Name:              "Arcane Missiles Activation",
		ActionID:          core.ActionID{SpellID: 79684},
		ClassSpellMask:    mage.MageSpellsAll ^ mage.MageSpellArcaneMissilesTick,
		SpellFlagsExclude: core.SpellFlagHelpful,
		ProcChance:        0.4,
		Callback:          core.CallbackOnSpellHitDealt,
		Outcome:           core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			arcane.arcaneMissilesProcAura.Activate(sim)
			arcane.arcaneMissilesProcAura.AddStack(sim)
		},
	})
}
