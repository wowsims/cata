package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

// T14 DPS
var ItemSetBattlegearOfTheLostCatacomb = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of the Lost Catacomb",
	ID:   1123,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Your Obliterate, Frost Strike, and Scourge Strike deal 4% increased damage.
			dk := agent.(DeathKnightAgent).GetDeathKnight()
			if dk.Spec == proto.Spec_SpecBloodDeathKnight {
				return
			}

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  DeathKnightSpellFrostStrike | DeathKnightSpellObliterate | DeathKnightSpellScourgeStrike,
				FloatValue: 0.04,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Your Pillar of Frost ability grants 5% additional Strength, and your Unholy Frenzy ability grants 10% additional haste.
			dk := agent.(DeathKnightAgent).GetDeathKnight()
			if dk.Spec == proto.Spec_SpecBloodDeathKnight {
				return
			}

			// Handled in sim/core/buffs.go and sim/death_knight/frost/pillar_of_frost.go
			dk.T14Dps4pc = setBonusAura
		},
	},
})

// T14 Tank
var ItemSetPlateOfTheLostCatacomb = core.NewItemSet(core.ItemSet{
	Name: "Plate of the Lost Catacomb",
	ID:   1124,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Reduces the cooldown of your Vampiric Blood ability by 20 sec.
			dk := agent.(DeathKnightAgent).GetDeathKnight()
			if dk.Spec != proto.Spec_SpecBloodDeathKnight {
				return
			}

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_Cooldown_Flat,
				ClassMask: DeathKnightSpellVampiricBlood,
				TimeValue: time.Second * -20,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Increases the healing received from your Death Strike by 10%.
			dk := agent.(DeathKnightAgent).GetDeathKnight()

			setBonusAura.AttachMultiplicativePseudoStatBuff(
				&dk.deathStrikeHealingMultiplier, 1.1,
			)
		},
	},
})

// T15 DPS
var ItemSetBattleplateOfTheAllConsumingMaw = core.NewItemSet(core.ItemSet{
	Name: "Battleplate of the All-Consuming Maw",
	ID:   1152,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Your attacks have a chance to raise the spirit of a fallen Zandalari as your Death Knight minion for 15 sec.
			// (Approximately 1.15 procs per minute)
			dk := agent.(DeathKnightAgent).GetDeathKnight()

			risenZandalariSpell := dk.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 138342},
				SpellSchool: core.SpellSchoolPhysical,
				Flags:       core.SpellFlagPassiveSpell,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					for _, troll := range dk.FallenZandalari {
						if troll.IsActive() {
							continue
						}

						troll.EnableWithTimeout(sim, troll, time.Second*15)

						return
					}

					if sim.Log != nil {
						dk.Log(sim, "No Fallen Zandalari available for the T15 4pc to proc, this is unreasonable.")
					}
				},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Callback: core.CallbackOnSpellHitDealt,
				Outcome:  core.OutcomeLanded,
				DPM: dk.NewRPPMProcManager(138343, false, core.ProcMaskDirect, core.RPPMConfig{
					PPM: 1.14999997616,
				}),
				ICD: time.Millisecond * 250,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					risenZandalariSpell.Cast(sim, result.Target)
				},
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Your Soul Reaper ability now deals additional Shadow damage to targets below 45% instead of below 35%.
			// Additionally, Killing Machine now also increases the critical strike chance of Soul Reaper.
			dk := agent.(DeathKnightAgent).GetDeathKnight()

			// KM effect handled in sim/death_knight/frost/killing_machine.go
			dk.T15Dps4pc = setBonusAura.AttachAdditivePseudoStatBuff(
				&dk.soulReaperHealthThreshold, 0.1,
			).ExposeToAPL(138347)
		},
	},
})

// T15 Tank
var ItemSetPlateOfTheAllConsumingMaw = core.NewItemSet(core.ItemSet{
	Name: "Plate of the All-Consuming Maw",
	ID:   1151,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Reduces the cooldown of your Rune Tap ability by 10 sec and removes its Rune cost.
			dk := agent.(DeathKnightAgent).GetDeathKnight()
			if dk.Spec != proto.Spec_SpecBloodDeathKnight {
				return
			}

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_Cooldown_Flat,
				ClassMask: DeathKnightSpellRuneTap,
				TimeValue: time.Second * -10,
			}).AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_PowerCost_Pct,
				ClassMask:  DeathKnightSpellRuneTap,
				FloatValue: -2.0,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Your Bone Shield ability grants you 15 Runic Power each time one of its charges is consumed.
			dk := agent.(DeathKnightAgent).GetDeathKnight()
			if dk.Spec != proto.Spec_SpecBloodDeathKnight {
				return
			}

			rpMetrics := dk.NewRunicPowerMetrics(core.ActionID{SpellID: 138214})

			dk.OnSpellRegistered(func(spell *core.Spell) {
				if !spell.Matches(DeathKnightSpellBoneShield) {
					return
				}

				dk.BoneShieldAura.ApplyOnStacksChange(func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
					if !setBonusAura.IsActive() {
						return
					}

					if newStacks < oldStacks {
						dk.AddRunicPower(sim, 15, rpMetrics)
					}
				})
			})
		},
	},
})
