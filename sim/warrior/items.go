package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

// T14 - DPS
var ItemSetBattleplateOfResoundingRings = core.NewItemSet(core.ItemSet{
	Name: "Battleplate of Resounding Rings",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				ClassMask:  SpellMaskMortalStrike | SpellMaskBloodthirst,
				Kind:       core.SpellMod_DamageDone_Pct,
				FloatValue: 0.25,
			})

			setBonusAura.ExposeToAPL(123142)
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				ClassMask: SpellMaskRecklessness,
				Kind:      core.SpellMod_Cooldown_Flat,
				TimeValue: -90 * time.Second,
			})

			setBonusAura.ExposeToAPL(123144)
		},
	},
})

// T14 - Tank
var ItemSetPlateOfResoundingRings = core.NewItemSet(core.ItemSet{
	Name: "Plate of Resounding Rings",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				ClassMask: SpellMaskLastStand,
				Kind:      core.SpellMod_Cooldown_Flat,
				TimeValue: -60 * time.Second,
			})

			setBonusAura.ExposeToAPL(123146)
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			war := agent.(WarriorAgent).GetWarrior()

			war.T14Tank2P = setBonusAura

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				ClassMask: SpellMaskShieldBlock,
				Kind:      core.SpellMod_PowerCost_Flat,
				IntValue:  -5,
			})

			setBonusAura.ExposeToAPL(123147)
		},
	},
})

// T15 - DPS
var ItemSetBattleplateOfTheLastMogu = core.NewItemSet(core.ItemSet{
	Name: "Battleplate of the Last Mogu",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			war := agent.(WarriorAgent).GetWarrior()

			core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
				Name:     "Item - Warrior T15 DPS 2P Bonus",
				ActionID: core.ActionID{SpellID: 138120},
				ICD:      250 * time.Millisecond,
				DPM: war.NewSetBonusRPPMProcManager(138120, setBonusAura, core.ProcMaskMeleeWhiteHit, core.RPPMConfig{
					PPM: 1.1,
				}.WithSpecMod(-0.625, proto.Spec_SpecFuryWarrior)),
				Outcome:  core.OutcomeHit,
				Callback: core.CallbackOnSpellHitDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					war.EnrageAura.Deactivate(sim)
					war.EnrageAura.Activate(sim)
				},
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			war := agent.(WarriorAgent).GetWarrior()

			aura := war.RegisterAura(core.Aura{
				Label:    "Skull Banner - T15 4P Bonus",
				ActionID: core.ActionID{SpellID: 138127},
				Duration: 10 * time.Second,
			}).AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Percent,
				FloatValue: 35,
			})

			war.OnSpellRegistered(func(spell *core.Spell) {
				if !spell.Matches(SpellMaskSkullBanner) {
					return
				}

				war.SkullBannerAura.AttachDependentAura(aura)
			})

			setBonusAura.ExposeToAPL(138126)
		},
	},
})

// T15 - Tank
var ItemSetPlaceOfTheLastMogu = core.NewItemSet(core.ItemSet{
	Name: "Plate of the Last Mogu",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			war := agent.(WarriorAgent).GetWarrior()
			war.T15Tank2P = core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
				Name:           "Victorious -  T15 Protection 2P Bonus",
				ActionID:       core.ActionID{SpellID: 138279},
				ClassSpellMask: SpellMaskRevenge | SpellMaskShieldSlam,
				ProcChance:     0.1,
				Outcome:        core.OutcomeHit,
				Callback:       core.CallbackOnSpellHitDealt,
				Duration:       15 * time.Second,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					war.VictoryRushAura.Activate(sim)
				},
			})

			setBonusAura.ExposeToAPL(138280)
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			war := agent.(WarriorAgent).GetWarrior()
			war.T15Tank4P = setBonusAura

			setBonusAura.ExposeToAPL(138281)
		},
	},
})

// T16 - DPS
var ItemSetBattleplateOfThePrehistoricMarauder = core.NewItemSet(core.ItemSet{
	Name: "Battleplate of the Prehistoric Marauder",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			war := agent.(WarriorAgent).GetWarrior()
			actionID := core.ActionID{SpellID: 144438}
			rageMetrics := war.NewRageMetrics(actionID)

			core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
				Name:     "Colossal Rage",
				ActionID: actionID,
				ProcMask: core.ProcMaskMeleeSpecial,
				Outcome:  core.OutcomeHit,
				Callback: core.CallbackOnSpellHitDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if war.ColossusSmashAuras.Get(result.Target).IsActive() {
						war.AddRage(sim, 5, rageMetrics)
					}
				},
			})

			setBonusAura.ExposeToAPL(144436)
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			war := agent.(WarriorAgent).GetWarrior()

			costMod := war.AddDynamicMod(core.SpellModConfig{
				ClassMask: SpellMaskExecute,
				Kind:      core.SpellMod_PowerCost_Flat,
				IntValue:  -30,
			})

			war.T16Dps4P = war.RegisterAura(core.Aura{
				Label:    "Death Sentence",
				ActionID: core.ActionID{SpellID: 144442},
				Duration: 12 * time.Second,
			}).ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
				if sim.IsExecutePhase20() {
					costMod.Activate()
				}
			}).ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
				costMod.Deactivate()
			})

			core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
				Name:           "Death Sentence - Consume",
				ClassSpellMask: SpellMaskExecute,
				Callback:       core.CallbackOnCastComplete,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					war.T16Dps4P.Deactivate(sim)
				},
			})

			core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
				Name:           "Death Sentence - Trigger",
				ActionID:       core.ActionID{SpellID: 144442},
				ClassSpellMask: SpellMaskMortalStrike | SpellMaskBloodthirst,
				Outcome:        core.OutcomeHit,
				ProcChance:     0.1,
				Callback:       core.CallbackOnSpellHitDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					war.T16Dps4P.Activate(sim)
				},
			})

			setBonusAura.ExposeToAPL(144441)
		},
	},
})

// T16 - Tank
var ItemSetPlateOfThePrehistoricMarauder = core.NewItemSet(core.ItemSet{
	Name: "Plate of the Prehistoric Marauder",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// TODO: You heal for 30% of all damage blocked with a shield
			war := agent.(WarriorAgent).GetWarrior()
			healthMetrics := war.NewHealthMetrics(core.ActionID{SpellID: 144503})

			war.OnSpellRegistered(func(spell *core.Spell) {
				if !spell.Matches(SpellMaskShieldBarrier) {
					return
				}

				war.ShieldBarrierAura.Aura.ApplyOnStacksChange(func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
					if setBonusAura.IsActive() {
						absorbLoss := float64(oldStacks - newStacks)
						war.GainHealth(sim, absorbLoss*0.3, healthMetrics)
					}
				})
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:     "Item - Warrior T16 Tank 2P Bonus",
				Callback: core.CallbackOnSpellHitTaken,
				Outcome:  core.OutcomeBlock,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					blockDamageReduction := result.Target.BlockDamageReduction()
					preBlockDamage := result.Damage / (1 - blockDamageReduction) * blockDamageReduction
					blockedDamage := result.Damage - preBlockDamage
					war.GainHealth(sim, blockedDamage*0.3, healthMetrics)
				},
			})

			setBonusAura.ExposeToAPL(144503)
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			war := agent.(WarriorAgent).GetWarrior()
			actionID := core.ActionID{SpellID: 144500}
			rageMetrics := war.NewRageMetrics(actionID)

			aura := war.RegisterAura(core.Aura{
				Label:    "Reckless Defense",
				ActionID: actionID,
				Duration: 10 * time.Second,
			})

			for _, aura := range war.DemoralizingShoutAuras {
				aura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
					if setBonusAura.IsActive() {
						aura.Activate(sim)
					}
				})
			}

			war.AddDynamicDamageTakenModifier(func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult, isPeriodic bool) {
				if aura.IsActive() {
					war.AddRage(sim, result.Damage/war.MaxHealth()*100, rageMetrics)
				}
			})

			setBonusAura.ExposeToAPL(144502)
		},
	},
})

func init() {
}
