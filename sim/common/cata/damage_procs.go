package cata

import (
	"time"

	"github.com/wowsims/cata/sim/common/shared"
	"github.com/wowsims/cata/sim/core"
)

type SpellCritProvider interface {
	DefaultSpellCritMultiplier() float64
}

func init() {
	shared.NewProcDamageEffect(shared.ProcDamageEffect{
		ItemID:  62049,
		SpellID: 89087,
		School:  core.SpellSchoolNature,
		MinDmg:  5250,
		MaxDmg:  8750,
		Flags:   core.SpellFlagNoSpellMods | core.SpellFlagIgnoreModifiers | core.SpellFlagNoOnDamageDealt,
		Outcome: shared.OutcomeMeleeNoBlockDodgeParryCrit,
		Trigger: core.ProcTrigger{
			Name:     "Darkmoon Card: Hurricane",
			ProcMask: core.ProcMaskMeleeOrRanged,
			PPM:      1,
			Outcome:  core.OutcomeLanded,
			Callback: core.CallbackOnSpellHitDealt,
		},
	})

	shared.NewProcDamageEffect(shared.ProcDamageEffect{
		ItemID:  62051,
		SpellID: 89087,
		School:  core.SpellSchoolNature,
		MinDmg:  5250,
		MaxDmg:  8750,
		Outcome: shared.OutcomeMeleeNoBlockDodgeParryCrit,
		Flags:   core.SpellFlagNoSpellMods | core.SpellFlagIgnoreModifiers | core.SpellFlagNoOnDamageDealt,
		Trigger: core.ProcTrigger{
			Name:     "Darkmoon Card: Hurricane",
			ProcMask: core.ProcMaskMeleeOrRanged,
			PPM:      1,
			Outcome:  core.OutcomeLanded,
			Callback: core.CallbackOnSpellHitDealt,
		},
	})

	core.NewItemEffect(68925, func(agent core.Agent) {
		character := agent.GetCharacter()
		// Proc chance determined to be p=.48 by video research - Researched by InDebt & Frostbitten
		// Research: https://github.com/wowsims/cata/pull/1009#issuecomment-2348700653
		procChance := 0.5
		dummyAura := character.GetOrRegisterAura(core.Aura{
			Label:     "Electrical Charge",
			ActionID:  core.ActionID{SpellID: 96890},
			Duration:  core.NeverExpires,
			MaxStacks: 10,
		})

		lightningSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 96891},
			SpellSchool:      core.SpellSchoolNature,
			ProcMask:         core.ProcMaskEmpty,
			Flags:            core.SpellFlagNoOnDamageDealt,
			DamageMultiplier: 1,
			CritMultiplier:   agent.GetDefaultSpellValueProvider().DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(985, 1266) * float64(dummyAura.GetStacks())
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			},
		})

		core.MakePermanent(core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Electrical Charge Aura",
			ActionID:   core.ActionID{ItemID: 68925},
			Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
			ProcMask:   core.ProcMaskSpellOrProc,
			ProcChance: 1,
			Outcome:    core.OutcomeCrit,
			ICD:        time.Millisecond * 2500,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				dummyAura.Activate(sim)
				dummyAura.AddStack(sim)

				if sim.Proc(procChance, "Variable Pulse Lightning Capacitor") {
					lightningSpell.Cast(sim, result.Target)
					dummyAura.SetStacks(sim, 0)
				}
			},
		}))
	})

	core.NewItemEffect(69110, func(agent core.Agent) {
		character := agent.GetCharacter()
		// Proc chance determined to be p=.48 by video research - Researched by InDebt & Frostbitten
		// Research: https://github.com/wowsims/cata/pull/1009#issuecomment-2348700653
		procChance := 0.5
		dummyAura := character.GetOrRegisterAura(core.Aura{
			Label:     "Electrical Charge",
			ActionID:  core.ActionID{SpellID: 96890},
			Duration:  core.NeverExpires,
			MaxStacks: 10,
		})

		lightningSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 96891},
			SpellSchool:      core.SpellSchoolNature,
			ProcMask:         core.ProcMaskEmpty,
			Flags:            core.SpellFlagNoOnDamageDealt,
			DamageMultiplier: 1,
			CritMultiplier:   agent.GetDefaultSpellValueProvider().DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(2889, 3713) * float64(dummyAura.GetStacks())
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			},
		})

		core.MakePermanent(core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Electrical Charge Aura",
			ActionID:   core.ActionID{ItemID: 69110},
			Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
			ProcMask:   core.ProcMaskSpellOrProc,
			ProcChance: 1,
			Outcome:    core.OutcomeCrit,
			ICD:        time.Millisecond * 2500,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				dummyAura.Activate(sim)
				dummyAura.AddStack(sim)

				if sim.Proc(procChance, "Variable Pulse Lightning Capacitor Heroic") {
					lightningSpell.Cast(sim, result.Target)
					dummyAura.SetStacks(sim, 0)
				}
			},
		}))
	})

}
