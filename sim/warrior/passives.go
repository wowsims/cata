package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func (war *Warrior) registerEnrage() {
	actionID := core.ActionID{SpellID: 12880}
	rageMetrics := war.NewRageMetrics(actionID)
	duration := time.Second * 6
	if war.Spec == proto.Spec_SpecFuryWarrior {
		duration = time.Second * 8
	}

	war.EnrageAura = war.RegisterAura(core.Aura{
		Label:    "Enrage",
		Tag:      EnrageTag,
		ActionID: actionID,
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			war.AddRage(sim, 10, rageMetrics)
			war.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			war.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.1
		},
	})

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:           "Enrage - Trigger",
		ActionID:       actionID,
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: SpellMaskColossusSmash | SpellMaskShieldSlam | SpellMaskDevastate | SpellMaskBloodthirst | SpellMaskMortalStrike,
		Outcome:        core.OutcomeCrit,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			war.EnrageAura.Deactivate(sim)
			war.EnrageAura.Activate(sim)
		},
	})
}

func (warrior *Warrior) registerDeepWounds() {
	actionID := core.ActionID{SpellID: 115768}
	deepWoundsCoeff := 0.28900000453
	deepWoundsBonusCoeff := 0.39599999785

	damageMultiplier := 1.0
	if warrior.Spec == proto.Spec_SpecArmsWarrior {
		// Arms has a 200% damage bonus to Deep Wounds
		damageMultiplier *= 2
	} else if warrior.Spec == proto.Spec_SpecProtectionWarrior {
		// Protection has a 50% damage bonus to Deep Wounds
		damageMultiplier *= 0.5
	}

	warrior.DeepWounds = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreArmor | core.SpellFlagIgnoreAttackerModifiers | SpellFlagBleed | core.SpellFlagPassiveSpell,

		DamageMultiplier: damageMultiplier,
		CritMultiplier:   warrior.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Deep Wounds",
			},
			NumberOfTicks: 5,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseDamage := warrior.CalcScalingSpellDmg(deepWoundsCoeff)
				baseDamage += deepWoundsBonusCoeff * dot.Spell.MeleeAttackPower()
				dot.SnapshotPhysical(target, baseDamage)
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickPhysicalCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dot := spell.Dot(target)
			dot.Apply(sim)
		},
	})

	core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
		Name:           "Deep Wounds - Trigger",
		ActionID:       actionID,
		ClassSpellMask: SpellMaskMortalStrike | SpellMaskBloodthirst | SpellMaskDevastate,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.DeepWounds.Cast(sim, result.Target)
		},
	})
}

func (war *Warrior) registerBloodAndThunder() {
	if war.Spec == proto.Spec_SpecFuryWarrior {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  SpellMaskThunderClap,
		FloatValue: 0.5,
	})

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:           "Blood and Thunder",
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: SpellMaskThunderClap,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			for _, target := range sim.Encounter.TargetUnits {
				dot := war.DeepWounds.Dot(target)
				dot.Apply(sim)
			}
		},
	})
}
