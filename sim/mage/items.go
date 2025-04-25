package mage

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

// T11
var ItemSetFirelordsVestments = core.NewItemSet(core.ItemSet{
	Name: "Firelord's Vestments",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			// Increases the critical strike chance of your Death Coil and Frost Strike abilities by 5%.
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Percent,
				ClassMask:  MageSpellArcaneMissilesTick | MageSpellIceLance | MageSpellPyroblast | MageSpellPyroblastDot,
				FloatValue: 5,
			})
		},
		4: func(_ core.Agent, setBonusAura *core.Aura) {
			//Reduces cast time of Arcane Blast, Fireball, FFB, and Frostbolt by 10%
			setBonusAura.AttachSpellMod(core.SpellModConfig{
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
	Bonuses: map[int32]core.ApplySetBonus{
		// You have a chance to summon a Mirror Image to assist you in battle for 15 sec when you cast Frostbolt, Fireball, Frostfire Bolt, or Arcane Blast.
		// (Proc chance: 20%, 45s cooldown)
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			mage := agent.(MageAgent).GetMage()

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
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
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			mage := agent.(MageAgent).GetMage()

			setBonusAura.ApplyOnGain(func(_ *core.Aura, _ *core.Simulation) {
				mage.brainFreezeProcChance += .15
				mage.baseHotStreakProcChance += 0.30
			})

			setBonusAura.ApplyOnExpire(func(_ *core.Aura, _ *core.Simulation) {
				mage.brainFreezeProcChance -= .15
				mage.baseHotStreakProcChance -= .30
			})

			setBonusAura.ExposeToAPL(99064)

			// Arcane Power Cost reduction implemented in:
			// talents_arcane.go#278
			mage.T12_4pc = setBonusAura
		},
	},
})

// T13
var ItemSetTimeLordsRegalia = core.NewItemSet(core.ItemSet{
	Name: "Time Lord's Regalia",
	Bonuses: map[int32]core.ApplySetBonus{
		// Your Arcane Blast has a 100% chance and your Fireball, Pyroblast, Frostfire Bolt, and Frostbolt spells have a 50% chance to grant Stolen Time, increasing your haste rating by 50 for 30 sec and stacking up to 10 times.
		// When Arcane Power, Combustion, or Icy Veins expires, all stacks of Stolen Time are lost.
		2: func(agent core.Agent, setBonusAura *core.Aura) {
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

			newStolenTimeTrigger := func(procChance float64, spellMask int64) *core.Aura {
				return setBonusAura.MakeDependentProcTriggerAura(&character.Unit, core.ProcTrigger{
					Name:           fmt.Sprintf("Stolen Time Trigger %f", procChance),
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
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Cooldown reduction handlers can be found in:
			// combustion.go
			// talents_arcane.go
			// talents_frost.go

			mage := agent.(MageAgent).GetMage()
			mage.T13_4pc = setBonusAura
		},
	},
})
