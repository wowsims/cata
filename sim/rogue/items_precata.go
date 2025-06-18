package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func init() {
	// Felstriker
	core.NewItemEffect(12590, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()

		effectAura := character.NewTemporaryStatsAura("Felstriker", core.ActionID{SpellID: 16551}, stats.Stats{stats.PhysicalCritPercent: 100, stats.PhysicalHitPercent: 100}, time.Second*3)
		procMask := character.GetDynamicProcMaskForWeaponEffect(12590)
		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Felstriker Trigger",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          *procMask,
			SpellFlagsExclude: core.SpellFlagPassiveSpell,
			DPM:               character.NewDynamicLegacyProcForWeapon(12590, 1, 0),
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				effectAura.Activate(sim)
			},
		})
	})

	// Rod of the Sun King
	core.NewItemEffect(29996, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()

		actionID := core.ActionID{ItemID: 29996}
		resourceMetricsEnergy := character.NewEnergyMetrics(actionID)
		procMask := character.GetDynamicProcMaskForWeaponEffect(29996)

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Rod of the Sun King",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          *procMask,
			SpellFlagsExclude: core.SpellFlagPassiveSpell,
			DPM:               character.NewDynamicLegacyProcForWeapon(29996, 1, 0),
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				spell.Unit.AddEnergy(sim, 10, resourceMetricsEnergy)
			},
		})
	})

	// Heartpierce
	// Normal
	core.NewItemEffect(49982, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 49982}
		resourceMetricsEnergy := character.NewEnergyMetrics(actionID)

		hpSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Hot: core.DotConfig{
				TickLength:    time.Second * 2,
				NumberOfTicks: 6,
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					character.AddEnergy(sim, 4, resourceMetricsEnergy)
				},
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Heartpierce Trigger",
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeLanded,
			ProcMask: core.ProcMaskMelee,
			DPM:      character.NewDynamicLegacyProcForWeapon(49982, 1, 0),
			ActionID: core.ActionID{ItemID: 49982},

			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				hpSpell.Hot(&character.Unit).Activate(sim)
			},
		})
	})

	// Heroic
	core.NewItemEffect(50641, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 50641}
		resourceMetricsEnergy := character.NewEnergyMetrics(actionID)

		hpSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Hot: core.DotConfig{
				TickLength:    time.Second * 2,
				NumberOfTicks: 6,
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					character.AddEnergy(sim, 4, resourceMetricsEnergy)
				},
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Heartpierce Trigger",
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeLanded,
			ProcMask: core.ProcMaskMelee,
			DPM:      character.NewDynamicLegacyProcForWeapon(50641, 1, 0),
			ActionID: core.ActionID{ItemID: 50641},

			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				hpSpell.Hot(&character.Unit).Activate(sim)
			},
		})
	})
}
