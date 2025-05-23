package shadow

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/priest"
)

func (shadow *ShadowPriest) registerSurgeOfDarkness() {
	shadow.SurgeOfDarkness = shadow.RegisterAura(core.Aura{
		Label:     "Surge of Darkness",
		ActionID:  core.ActionID{SpellID: 87160},
		MaxStacks: 2,
		Duration:  core.NeverExpires,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell != shadow.MindSpike {
				return
			}

			aura.RemoveStack(sim)
		},
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask:  priest.PriestSpellMindSpike,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.5,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask:  priest.PriestSpellMindSpike,
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -100,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask:  priest.PriestSpellMindSpike,
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -1,
	})

	// Always register the auras above, or else APL might behave funky
	// if the aura is not registered and we write a condition for it, it might evaluate to nil
	// causing us to not have a condition to begin with
	if !shadow.Talents.FromDarknessComesLight {
		return
	}

	core.MakePermanent(shadow.RegisterAura(core.Aura{
		Label: "Surge of Darkness (Talent)",
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// use class mask here due to mastery duplication
			if spell.ClassSpellMask&priest.PriestSpellVampiricTouch == 0 {
				return
			}

			if sim.Proc(0.2, "Roll Surge of Darkness") {
				shadow.SurgeOfDarkness.Activate(sim)
				shadow.SurgeOfDarkness.AddStack(sim)
			}
		},
	}))
}

func (shadow *ShadowPriest) registerSolaceAndInstanity() {
	if !shadow.Talents.SolaceAndInsanity {
		return
	}

	shadow.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 129197},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagChanneled | core.SpellFlagAPL,
		ClassSpellMask: priest.PriestSpellMindFlay,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 1,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,
		CritMultiplier:           shadow.DefaultCritMultiplier(),
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "MindFlay-Insanity",
			},
			NumberOfTicks:        3,
			TickLength:           time.Second * 1,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: true,
			BonusCoefficient:     MfCoeff,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, shadow.CalcScalingSpellDmg(MfScale))
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
			return spell.CalcPeriodicDamage(sim, target, shadow.CalcScalingSpellDmg(MfScale), spell.OutcomeExpectedMagicCrit)
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return shadow.DevouringPlague.Dot(target).IsActive()
		},
	})

	dmgMod := shadow.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  priest.PriestSpellMindFlay,
		FloatValue: 0.33,
	})

	shadow.OnSpellRegistered(func(spell *core.Spell) {
		if spell.ClassSpellMask == priest.PriestSpellDevouringPlagueDoT {
			for _, target := range shadow.Env.Encounter.TargetUnits {
				dot := spell.Dot(target)
				if dot != nil {
					dot.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
						dmgMod.UpdateFloatValue(float64(shadow.orbsConsumed) * 1 / 3)
						dmgMod.Activate()
					})
					dot.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
						dmgMod.Deactivate()
					})
				}
			}
		}
	})
}

func (shadow *ShadowPriest) registerTwistOfFate() {
	dmgMod := shadow.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		School:     core.SpellSchoolShadow | core.SpellSchoolHoly,
		FloatValue: 0.15,
	})

	tofAura := shadow.RegisterAura(core.Aura{
		Label:    "Twist of Fate",
		ActionID: core.ActionID{SpellID: 123254},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dmgMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dmgMod.Deactivate()
		},
	})

	if !shadow.Talents.TwistOfFate {
		return
	}
	core.MakePermanent(shadow.RegisterAura(core.Aura{
		Label: "Twist of Fate (Talent)",
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.IsExecutePhase35() {
				tofAura.Activate(sim)
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.IsExecutePhase35() {
				tofAura.Activate(sim)
			}
		},
	}))
}

func (shadow *ShadowPriest) registerDivineInsight() {
	castTimeMod := shadow.AddDynamicMod(core.SpellModConfig{
		ClassMask:  priest.PriestSpellMindBlast,
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -1,
	})

	costMod := shadow.AddDynamicMod(core.SpellModConfig{
		ClassMask:  priest.PriestSpellMindBlast,
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -100,
	})

	procAura := shadow.RegisterAura(core.Aura{
		Label:    "Divine Insight",
		Duration: time.Second * 12,
		ActionID: core.ActionID{SpellID: 124430},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Activate()
			costMod.Activate()
			shadow.MindBlast.CD.Reset()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			costMod.Deactivate()
			castTimeMod.Deactivate()
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == shadow.MindBlast {
				aura.Deactivate(sim)
			}
		},
	})

	if !shadow.Talents.DivineInsight {
		return
	}

	core.MakePermanent(shadow.RegisterAura(core.Aura{
		Label: "Divine Insight (Talent)",
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ClassSpellMask&priest.PriestSpellShadowWordPain == 0 {
				return
			}

			if sim.Proc(0.05, "Divine Insight (Proc)") {
				procAura.Activate(sim)
			}
		},
	}))
}
