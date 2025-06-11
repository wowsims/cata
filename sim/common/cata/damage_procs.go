package cata

import (
	"time"

	"github.com/wowsims/mop/sim/common/shared"
	"github.com/wowsims/mop/sim/core"
)

func init() {
	shared.NewProcDamageEffect(shared.ProcDamageEffect{
		ItemID:  62049,
		SpellID: 89087,
		School:  core.SpellSchoolNature,
		MinDmg:  5250,
		MaxDmg:  8750,
		Flags:   core.SpellFlagNoSpellMods | core.SpellFlagIgnoreModifiers | core.SpellFlagNoOnDamageDealt,
		Outcome: shared.OutcomeMeleeNoBlockDodgeParryCrit,
		TriggerDPM: func(c *core.Character) *core.DynamicProcManager {
			return c.NewLegacyPPMManager(1, core.ProcMaskMeleeOrRanged)
		},
		Trigger: core.ProcTrigger{
			Name:     "Darkmoon Card: Hurricane",
			ProcMask: core.ProcMaskMeleeOrRanged,
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
		Flags:   core.SpellFlagNoSpellMods | core.SpellFlagIgnoreModifiers | core.SpellFlagNoOnDamageDealt,
		Outcome: shared.OutcomeMeleeNoBlockDodgeParryCrit,
		TriggerDPM: func(c *core.Character) *core.DynamicProcManager {
			return c.NewLegacyPPMManager(1, core.ProcMaskMeleeOrRanged)
		},
		Trigger: core.ProcTrigger{
			Name:     "Darkmoon Card: Hurricane",
			ProcMask: core.ProcMaskMeleeOrRanged,
			Outcome:  core.OutcomeLanded,
			Callback: core.CallbackOnSpellHitDealt,
		},
	})

	core.NewItemEffect(68925, func(agent core.Agent) {
		character := agent.GetCharacter()
		// Proc chance determined to be p=.48 by video research - Researched by InDebt & Frostbitten
		// Research: https://github.com/wowsims/mop/pull/1009#issuecomment-2348700653
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
			ProcMask:         core.ProcMaskProc,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultCritMultiplier(),
			ThreatMultiplier: 1,
			MissileSpeed:     20,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				// Tooltip is wrong:
				// https://wago.tools/db2/SpellEffect?build=4.4.1.56574&filter[SpellID]=96887%7C96891%7C97119&page=1&sort[SpellID]=asc
				baseDamage := sim.Roll(2561, 3292) * float64(dummyAura.GetStacks())
				result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
				spell.WaitTravelTime(sim, func(sim *core.Simulation) {
					spell.DealDamage(sim, result)
				})
			},
		})

		aura := core.MakePermanent(core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Electrical Charge Aura",
			ActionID:   core.ActionID{ItemID: 68925},
			Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
			ProcMask:   core.ProcMaskSpellOrSpellProc,
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

		character.ItemSwap.RegisterProc(68925, aura)
	})

	core.NewItemEffect(69110, func(agent core.Agent) {
		character := agent.GetCharacter()
		// Proc chance determined to be p=.48 by video research - Researched by InDebt & Frostbitten
		// Research: https://github.com/wowsims/mop/pull/1009#issuecomment-2348700653
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
			ProcMask:         core.ProcMaskProc,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultCritMultiplier(),
			ThreatMultiplier: 1,
			MissileSpeed:     20,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(2889, 3713) * float64(dummyAura.GetStacks())
				result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
				spell.WaitTravelTime(sim, func(sim *core.Simulation) {
					spell.DealDamage(sim, result)
				})
			},
		})

		aura := core.MakePermanent(core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Electrical Charge Aura",
			ActionID:   core.ActionID{ItemID: 69110},
			Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
			ProcMask:   core.ProcMaskSpellOrSpellProc,
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

		character.ItemSwap.RegisterProc(69110, aura)
	})

}
