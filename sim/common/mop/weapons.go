package mop

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func init() {

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
