package mop

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func init() {

	// Yaungol Fire Carrier
	core.NewItemEffect(86518, func(agent core.Agent, state proto.ItemLevelState) {
		character := agent.GetCharacter()

		statValue := core.GetItemEffectScaling(86518, 0.58200001717, state)

		dot := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 126211},
			SpellSchool: core.SpellSchoolFire,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagIgnoreArmor | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultCritMultiplier(),
			ThreatMultiplier: 1,

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "Yaungol Fire",
				},
				NumberOfTicks: 5,
				TickLength:    time.Second * 2,

				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.SnapshotPhysical(target, statValue)
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.Dot(target).Apply(sim)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Yaungol Fire Carrier",
			ActionID:   core.ActionID{SpellID: 126212},
			Harmful:    true,
			ProcMask:   core.ProcMaskMeleeOrMeleeProc,
			ProcChance: 0.1,
			Outcome:    core.OutcomeLanded,
			Callback:   core.CallbackOnSpellHitDealt,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				dot.Cast(sim, result.Target)
			},
		})
	})

	// The Gloaming Blade
	core.NewItemEffect(88149, func(agent core.Agent, state proto.ItemLevelState) {
		character := agent.GetCharacter()

		statValue := core.GetItemEffectScaling(88149, 0.19400000572, state)

		aura, _ := character.NewTemporaryStatBuffWithStacks(core.TemporaryStatBuffWithStacksConfig{
			AuraLabel:     "The Deepest Night",
			ActionID:      core.ActionID{SpellID: 127890},
			Duration:      time.Second * 10,
			MaxStacks:     3,
			BonusPerStack: stats.Stats{stats.CritRating: statValue},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "The Gloaming Blade",
			Harmful:  true,
			DPM:      character.NewDynamicLegacyProcForWeapon(88149, 2, 0),
			Outcome:  core.OutcomeLanded,
			Callback: core.CallbackOnSpellHitDealt,
			Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				aura.Activate(sim)
				aura.AddStack(sim)
			},
		})
	})

}
