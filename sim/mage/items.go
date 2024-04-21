package mage

import (
	"github.com/wowsims/cata/sim/core"
)

// T11
var ItemSetFirelordsVestments = core.NewItemSet(core.ItemSet{
	Name: "Firelord's Vestments",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the critical strike chance of your Death Coil and Frost Strike abilities by 5%.
			agent.GetCharacter().AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Rating,
				ClassMask:  MageSpellArcaneMissilesTick | MageSpellIceLance | MageSpellPyroblast,
				FloatValue: 5 * core.CritRatingPerCritChance,
			})
		},
		4: func(agent core.Agent) {
			//Reduces cast time of Arcane Blast, Fireball, FFB, and Frostbolt by 10%
			agent.GetCharacter().AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_CastTime_Pct,
				ClassMask:  MageSpellArcaneBlast | MageSpellFireball | MageSpellFrostfireBolt | MageSpellFrostbolt,
				FloatValue: -0.1,
			})
		},
	},
})

// T12
/* var ItemSetFirehawkRobesOfConflagration = core.NewItemSet(core.ItemSet{
	Name:            "Firehawk Robes of Conflagration",
	AlternativeName: "Firehawk",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:           "Firehawk Mirror Image",
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: MageSpellArcaneBlast | MageSpellFireball | MageSpellFrostfireBolt | MageSpellFrostbolt,
				ProcChance:     0.1, //Just a made up number for now
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					&character.EnableWithTimeout(character.NewMirrorImage())
				},
			})
		},
		4: func(agent core.Agent) {
			character := agent.GetCharacter()

			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_CastTime_Pct,
				ClassMask:  MageSpellArcaneBlast | MageSpellFireball | MageSpellFrostfireBolt | MageSpellFrostbolt,
				FloatValue: -0.1,
			})

			mage.arcanepowercostmod = &character.AddDynamicMod
		},
	},
}) */
