package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
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
			// Arcane Power Cost reduction implemented in:
			// talents_arcane.go#278

			mage := agent.(MageAgent).GetMage()

			core.MakePermanent(mage.RegisterAura(core.Aura{
				Label:    "Item - Mage T12 4P Bonus",
				ActionID: core.ActionID{SpellID: 99064},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					mage.brainFreezeProcChance += .15
					mage.baseHotStreakProcChance += 0.30
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					mage.brainFreezeProcChance -= .15
					mage.baseHotStreakProcChance -= .30
				},
			}))

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
			character := agent.GetCharacter()
			mage := agent.(MageAgent).GetMage()

			// Stack reset handlers can be found in:
			// combustion.go
			// talents_arcane.go
			// talents_frost.go
			mage.t13ProcAura = core.MakeStackingAura(character, core.StackingStatAura{
				Aura: core.Aura{
					Label:     "Stolen Time",
					ActionID:  core.ActionID{SpellID: 105785},
					Duration:  time.Second * 30,
					MaxStacks: 10,
				},
				BonusPerStack: stats.Stats{stats.HasteRating: 50},
			})

			newStolenTimeTrigger := func(procChance float64, spellMask int64) {
				core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
					Name:           "Stolen Time Trigger",
					ActionID:       core.ActionID{ItemID: 105788},
					Callback:       core.CallbackOnSpellHitDealt,
					ClassSpellMask: spellMask,
					ProcChance:     procChance,
					Outcome:        core.OutcomeLanded,
					Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						mage.t13ProcAura.Activate(sim)
						mage.t13ProcAura.AddStack(sim)
					},
				})
			}

			newStolenTimeTrigger(1, MageSpellArcaneBlast)
			newStolenTimeTrigger(0.5, MageSpellFireball|MageSpellPyroblast|MageSpellFrostfireBolt|MageSpellFrostbolt)
		},
		// Each stack of Stolen Time also reduces the cooldown of Arcane Power by 7 sec, Combustion by 5 sec, and Icy Veins by 15 sec.
		4: func(agent core.Agent) {
			// Cooldown reduction handlers can be found in:
			// combustion.go
			// talents_arcane.go
			// talents_frost.go
		},
	},
})
