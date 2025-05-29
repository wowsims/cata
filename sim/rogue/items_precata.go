package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func init() {
	// Felstriker
	core.NewItemEffect(12590, func(agent core.Agent) {
		character := agent.GetCharacter()

		effectAura := character.NewTemporaryStatsAura("Felstriker", core.ActionID{SpellID: 16551}, stats.Stats{stats.PhysicalCritPercent: 100, stats.PhysicalHitPercent: 100}, time.Second*3)
		procMask := character.GetDynamicProcMaskForWeaponEffect(12590)
		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:              "Felstriker Trigger",
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          *procMask,
			SpellFlagsExclude: core.SpellFlagPassiveSpell,
			PPM:               1,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				effectAura.Activate(sim)
			},
		})
	})

	// Rod of the Sun King
	core.NewItemEffect(29996, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetDynamicProcMaskForWeaponEffect(29996)
		pppm := character.AutoAttacks.NewPPMManager(1.0, *procMask)

		actionID := core.ActionID{ItemID: 29996}

		var resourceMetricsRage *core.ResourceMetrics
		var resourceMetricsEnergy *core.ResourceMetrics
		if character.HasRageBar() {
			resourceMetricsRage = character.NewRageMetrics(actionID)
		}
		if character.HasEnergyBar() {
			resourceMetricsEnergy = character.NewEnergyMetrics(actionID)
		}

		character.GetOrRegisterAura(core.Aura{
			Label:    "Rod of the Sun King",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if pppm.Proc(sim, spell.ProcMask, "Rod of the Sun King") {
					switch spell.Unit.GetCurrentPowerBar() {
					case core.RageBar:
						spell.Unit.AddRage(sim, 5, resourceMetricsRage)
					case core.EnergyBar:
						spell.Unit.AddEnergy(sim, 10, resourceMetricsEnergy)
					}
				}
			},
		})
	})
}
