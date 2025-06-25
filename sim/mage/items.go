package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

// T14
var ItemSetRegaliaOfTheBurningScroll = core.NewItemSet(core.ItemSet{
	Name: "Regalia of the Burning Scroll",
	Bonuses: map[int32]core.ApplySetBonus{
		// Increases the damage done by your Arcane Missiles spell by 7%,
		// increases the damage done by your Pyroblast spell by 8%, and increases the damage done by your Ice Lance spell by 12%.
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  MageSpellIceLance,
				FloatValue: 0.12,
			}).AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  MageSpellArcaneMissilesTick,
				FloatValue: 0.07,
			}).AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  MageSpellPyroblast | MageSpellPyroblastDot,
				FloatValue: 0.08,
			})
			setBonusAura.ExposeToAPL(123097)
		},
		// Increases the damage bonus of Arcane Power by an additional 10%,
		// reduces the cooldown of Icy Veins by 50%, and reduces the cooldown of Combustion by 20%.
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			mage := agent.(MageAgent).GetMage()
			mage.T14_4pc = setBonusAura

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				FloatValue: 0.5,
				Kind:       core.SpellMod_Cooldown_Multiplier,
				ClassMask:  MageSpellIcyVeins,
			}).AttachSpellMod(core.SpellModConfig{
				FloatValue: 1 - 0.2,
				Kind:       core.SpellMod_Cooldown_Multiplier,
				ClassMask:  MageSpellCombustion,
			})
			setBonusAura.ExposeToAPL(123101)
		},
	},
})

// T15
var ItemSetRegaliaOfTheChromaticHydra = core.NewItemSet(core.ItemSet{
	Name: "Regalia of the Chromatic Hydra",
	Bonuses: map[int32]core.ApplySetBonus{
		// When Alter Time expires, you gain 1800 Haste, Crit, and Mastery for 30 sec.
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			mage := agent.(MageAgent).GetMage()
			statValue := 1800.0
			aura := mage.NewTemporaryStatsAura(
				"Time Lord",
				core.ActionID{SpellID: 138317},
				stats.Stats{stats.HasteRating: statValue, stats.CritRating: statValue, stats.MasteryRating: statValue},
				time.Second*30,
			)

			mage.OnSpellRegistered(func(spell *core.Spell) {
				if !spell.Matches(MageSpellAlterTime) {
					return
				}

				spell.RelatedSelfBuff.ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
					if setBonusAura.IsActive() {
						aura.Activate(sim)
					}
				})
			})

			setBonusAura.ExposeToAPL(138316)
		},
		// Increases the effects of Arcane Charges by 5%,
		// increases the critical strike chance of Pyroblast by 5%,
		// and increases the chance for your Frostbolt to trigger Fingers of Frost by an additional 6%.
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			mage := agent.(MageAgent).GetMage()

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Percent,
				ClassMask:  MageSpellPyroblast | MageSpellPyroblastDot,
				FloatValue: 5,
			})
			setBonusAura.ApplyOnGain(func(_ *core.Aura, _ *core.Simulation) {
				mage.T15_4PC_ArcaneChargeEffect += 0.05
				mage.T15_4PC_FrostboltProcChance += 0.06
			}).ApplyOnExpire(func(_ *core.Aura, _ *core.Simulation) {
				mage.T15_4PC_ArcaneChargeEffect -= 0.05
				mage.T15_4PC_FrostboltProcChance -= 0.06
			})

			setBonusAura.ExposeToAPL(138376)
		},
	},
})

// T16
var ItemSetChronomancerRegalia = core.NewItemSet(core.ItemSet{
	Name: "Chronomancer Regalia",
	Bonuses: map[int32]core.ApplySetBonus{
		// Arcane Missiles causes your next Arcane Blast within 10 sec to cost 25% less mana, stacking up to 4 times.
		// Consuming Brain Freeze increases the damage of your next Ice Lance, Frostbolt, Frostfire Bolt, or Cone of Cold by 20%.
		// Consuming Pyroblast! increases your haste by 750 for 5 sec, stacking up to 5 times.
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			mage := agent.(MageAgent).GetMage()

			arcaneBlastMod := mage.AddDynamicMod(core.SpellModConfig{
				Kind:       core.SpellMod_PowerCost_Pct,
				ClassMask:  MageSpellArcaneBlast,
				FloatValue: 0,
			})

			arcaneAura := mage.GetOrRegisterAura(core.Aura{
				Label:     "Profound Magic",
				ActionID:  core.ActionID{SpellID: 145252},
				Duration:  time.Second * 10,
				MaxStacks: 4,
				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
					arcaneBlastMod.UpdateFloatValue(0.25 * float64(newStacks))
				},
			})

			setBonusAura.MakeDependentProcTriggerAura(&mage.Unit, core.ProcTrigger{
				Name:           "Profound Magic - Consume",
				ClassSpellMask: MageSpellArcaneBlast,
				Callback:       core.CallbackOnSpellHitDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					arcaneAura.Deactivate(sim)
				},
			})

			setBonusAura.MakeDependentProcTriggerAura(&mage.Unit, core.ProcTrigger{
				Name:           "Item - Mage T16 2P Bonus",
				ClassSpellMask: MageSpellArcaneMissilesCast,
				Callback:       core.CallbackOnSpellHitDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					arcaneAura.Activate(sim)
					arcaneAura.AddStack(sim)
				},
			})

			fireAura := core.MakeStackingAura(&mage.Character, core.StackingStatAura{
				Aura: core.Aura{
					Label:     "Potent Flames",
					ActionID:  core.ActionID{SpellID: 145254},
					Duration:  time.Second * 5,
					MaxStacks: 5,
				},
				BonusPerStack: stats.Stats{stats.HasteRating: 750},
			})

			mage.OnSpellRegistered(func(spell *core.Spell) {
				if !spell.Matches(MageSpellPyroblast) {
					return
				}

				mage.InstantPyroblastAura.ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
					if setBonusAura.IsActive() {
						fireAura.Activate(sim)
						fireAura.AddStack(sim)
					}
				})
			})

			frostClassMask := MageSpellIceLance | MageSpellFrostbolt | MageSpellFrostfireBolt | MageSpellConeOfCold
			frostAura := mage.GetOrRegisterAura(core.Aura{
				Label:    "Frozen Thoughts",
				ActionID: core.ActionID{SpellID: 146557},
				Duration: time.Second * 15,
			}).AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  frostClassMask,
				FloatValue: 0.25,
			})

			setBonusAura.MakeDependentProcTriggerAura(&mage.Unit, core.ProcTrigger{
				Name:           "Frozen Thoughts - Consume",
				ClassSpellMask: frostClassMask,
				Callback:       core.CallbackOnSpellHitDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					frostAura.Deactivate(sim)
				},
			})

			mage.OnSpellRegistered(func(spell *core.Spell) {
				if !spell.Matches(MageSpellFrostbolt) {
					return
				}

				mage.BrainFreezeAura.ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
					if setBonusAura.IsActive() {
						frostAura.Activate(sim)
					}
				})
			})

			setBonusAura.ExposeToAPL(145251)
		},
		// Arcane Missiles has a 15% chance to not consume Arcane Missiles!.
		// Consuming Brain Freeze has a 30% chance to drop an icy boulder on your target.
		// Inferno Blast also causes your next Pyroblast to be a critical strike.
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			mage := agent.(MageAgent).GetMage()

			mage.T16_4pc = setBonusAura

			frigidBlast := mage.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 145264},
				SpellSchool: core.SpellSchoolFrost,
				ProcMask:    core.ProcMaskSpellProc,
				Flags:       core.SpellFlagPassiveSpell,

				DamageMultiplier: 1,
				CritMultiplier:   mage.DefaultCritMultiplier(),
				ThreatMultiplier: 1,

				BonusCoefficient: 1.5,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					baseDamage := mage.CalcAndRollDamageRange(sim, 1.5, 0.15000000596)
					spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicCrit)
				},
			})

			mage.OnSpellRegistered(func(spell *core.Spell) {
				if !spell.Matches(MageSpellFrostbolt) {
					return
				}

				mage.BrainFreezeAura.ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
					if setBonusAura.IsActive() {
						frigidBlast.Cast(sim, mage.CurrentTarget)
					}
				})
			})

			fireAura := mage.GetOrRegisterAura(core.Aura{
				Label:    "Fiery Adept",
				ActionID: core.ActionID{SpellID: 145261},
				Duration: time.Second * 15,
			}).AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Percent,
				ClassMask:  MageSpellPyroblast | MageSpellPyroblastDot,
				FloatValue: 100,
			})

			setBonusAura.MakeDependentProcTriggerAura(&mage.Unit, core.ProcTrigger{
				Name:           "Fiery Adept - Consume",
				ClassSpellMask: MageSpellPyroblast,
				Harmful:        true,
				Callback:       core.CallbackOnSpellHitDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					fireAura.Deactivate(sim)
				},
			})

			setBonusAura.MakeDependentProcTriggerAura(&mage.Unit, core.ProcTrigger{
				Name:           "Item - Mage T16 4P Bonus",
				ClassSpellMask: MageSpellInfernoBlast,
				Callback:       core.CallbackOnSpellHitDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					fireAura.Activate(sim)
				},
			})

			setBonusAura.ExposeToAPL(145257)
		},
	},
})
