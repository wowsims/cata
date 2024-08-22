package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

// T11
var ItemSetFirelordsVestments = core.NewItemSet(core.ItemSet{
	Name: "Firelord's Vestments",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the critical strike chance of your Death Coil and Frost Strike abilities by 5%.
			agent.GetCharacter().AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Percent,
				ClassMask:  MageSpellArcaneMissilesTick | MageSpellIceLance | MageSpellPyroblast | MageSpellPyroblastDot,
				FloatValue: 5,
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
var ItemSetFirehawkRobesOfConflagration = core.NewItemSet(core.ItemSet{
	Name: "Firehawk Robes of Conflagration",
	Bonuses: map[int32]core.ApplyEffect{
		// You have a chance to summon a Mirror Image to assist you in battle for 15 sec when you cast Frostbolt, Fireball, Frostfire Bolt, or Arcane Blast.
		// (Proc chance: 20%, 45s cooldown)
		2: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()

			core.MakeProcTriggerAura(&mage.Unit, core.ProcTrigger{
				Name:           "Item - Mage T12 2P Bonus",
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: MageSpellArcaneBlast | MageSpellFireball | MageSpellFrostfireBolt | MageSpellFrostbolt,
				ProcChance:     0.20,
				ICD:            time.Second * 45,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					mage.t12MirrorImage.EnableWithTimeout(sim, mage.t12MirrorImage, time.Second*15)
				},
			})
		},
		// Your spells have an increased chance to trigger Brain Freeze or Hot Streak.
		// In addition, Arcane Power decreases the cost of your damaging spells by 10% instead of increasing their cost.
		4: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()

			core.MakePermanent(mage.RegisterAura(core.Aura{
				Label: "Item - Mage T12 4P Bonus",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					mage.brainFreezeProcChance += .15
					mage.hotStreakProcChance += 0.30
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					mage.brainFreezeProcChance -= .15
					mage.hotStreakProcChance -= .30
				},
			}))

			mage.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ClassSpellMask == MageSpellArcanePower {
					mage.arcanePowerCostMod.UpdateFloatValue(-0.1)
				}
			})
		},
	},
})

// T13
var ItemSetTimeLordsRegalia = core.NewItemSet(core.ItemSet{
	Name: "Time Lord's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		// Your Arcane Blast has a 100% chance and your Fireball, Pyroblast, Frostfire Bolt, and Frostbolt spells have a 50% chance to grant Stolen Time, increasing your haste rating by 50 for 30 sec and stacking up to 10 times.
		// When Arcane Power, Combustion, or Icy Veins expires, all stacks of Stolen Time are lost.
		2: func(agent core.Agent) {
			// mage := agent.(MageAgent).GetMage()
		},
		// Each stack of Stolen Time also reduces the cooldown of Arcane Power by 7 sec, Combustion by 5 sec, and Icy Veins by 15 sec.
		4: func(agent core.Agent) {
			// mage := agent.(MageAgent).GetMage()
		},
	},
})
