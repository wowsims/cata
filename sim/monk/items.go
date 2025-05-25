package monk

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

// T14 - Windwalker
var ItemSetBattlegearOfTheRedCrane = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of the Red Crane",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_Cooldown_Flat,
				ClassMask: MonkSpellFistsOfFury,
				TimeValue: 5 * time.Second,
			}).ExposeToAPL(123149)
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_BuffDuration_Flat,
				ClassMask: MonkSpellEnergizingBrew,
				TimeValue: 5 * time.Second,
			}).ExposeToAPL(123150)
		},
	},
})

// T14 - Brewmaster
var ItemSetArmorOfTheRedCrane = core.NewItemSet(core.ItemSet{
	Name: "Armor of the Red Crane",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			monk := agent.(MonkAgent).GetMonk()

			monk.OnSpellRegistered(func(spell *core.Spell) {
				if !spell.Matches(MonkSpellElusiveBrew) {
					return
				}

				hasDodgeBonus := false
				spell.RelatedSelfBuff.ApplyOnGain(func(_ *core.Aura, sim *core.Simulation) {
					if setBonusAura.IsActive() {
						monk.PseudoStats.BaseDodgeChance += 0.05
						hasDodgeBonus = true
					}
				}).ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
					if hasDodgeBonus {
						monk.PseudoStats.BaseDodgeChance -= 0.05
						hasDodgeBonus = false
					}
				})
			})

			setBonusAura.ExposeToAPL(123157)
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Implemented in guard.go
			monk := agent.(MonkAgent).GetMonk()
			monk.T14Brewmaster4P = setBonusAura

			setBonusAura.ExposeToAPL(123159)
		},
	},
})

// T15 - Windwalker
var ItemSetFireCharmBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Fire-Charm Battlegear",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			monk := agent.(MonkAgent).GetMonk()

			actionId := core.ActionID{SpellID: 138177}
			energyMetrics := monk.NewEnergyMetrics(actionId)

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:       "Item - Monk T15 Windwalker 2P Bonus",
				ActionID:   actionId,
				ProcChance: 0.15,
				ICD:        100 * time.Millisecond,
				SpellFlags: SpellFlagBuilder,
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					monk.AddEnergy(sim, 10, energyMetrics)
				},
			}).ExposeToAPL(actionId.SpellID)
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Implemented in windwalker/tigereye_brew.go
			monk := agent.(MonkAgent).GetMonk()
			monk.T15Windwalker4P = setBonusAura

			setBonusAura.ExposeToAPL(138315)
		},
	},
})

// T15 - Brewmaster
var ItemSetFireCharmArmor = core.NewItemSet(core.ItemSet{
	Name: "Fire-Charm Armor",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			monk := agent.(MonkAgent).GetMonk()

			monk.T15Brewmaster2P = monk.RegisterAura(core.Aura{
				Label:    "Item - Monk T15 Brewmaster 2P Bonus" + monk.Label,
				ActionID: core.ActionID{SpellID: 138233},
				Duration: 0,
			})

			monk.OnSpellRegistered(func(spell *core.Spell) {
				if !spell.Matches(MonkSpellElusiveBrew) {
					return
				}

				spell.RelatedSelfBuff.ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
					if setBonusAura.IsActive() {
						monk.T15Brewmaster2P.Duration = time.Duration(monk.ElusiveBrewStacks) * time.Second
						monk.T15Brewmaster2P.Activate(sim)
					}
				})
			})

			setBonusAura.ExposeToAPL(138231)
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			monk := agent.(MonkAgent).GetMonk()

			monk.T15Brewmaster4P = monk.RegisterAura(core.Aura{
				Label:    "Purifier" + monk.Label,
				ActionID: core.ActionID{SpellID: 138237},
				Duration: 15 * time.Second,
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Item - Monk T15 Brewmaster 4P Bonus",
				ActionID:       core.ActionID{SpellID: 138236},
				ClassSpellMask: MonkSpellStagger,
				ProcChance:     0.1,
				Callback:       core.CallbackOnPeriodicHealDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					monk.T15Brewmaster4P.Activate(sim)
				},
			})

			setBonusAura.ExposeToAPL(138236)
		},
	},
})

// T16 - Windwalker
var ItemSetBattlegearOfSevenSacredSeals = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of Seven Sacred Seals",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			monk := agent.(MonkAgent).GetMonk()

			registerComboBreakerDamageMod := func(spellID int32, spellMask int64) {
				monk.OnSpellRegistered(func(spell *core.Spell) {
					if !spell.Matches(spellMask) {
						return
					}

					aura := monk.GetAuraByID(core.ActionID{SpellID: spellID})
					if aura != nil {
						damageMod := monk.AddDynamicMod(core.SpellModConfig{
							Kind:       core.SpellMod_DamageDone_Pct,
							ClassMask:  spellMask,
							FloatValue: 0.4,
						})

						aura.ApplyOnGain(func(_ *core.Aura, sim *core.Simulation) {
							if setBonusAura.IsActive() {
								damageMod.Activate()
							}
						}).ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
							damageMod.Deactivate()
						})
					}
				})

			}

			registerComboBreakerDamageMod(118864, MonkSpellTigerPalm)
			registerComboBreakerDamageMod(116768, MonkSpellBlackoutKick)

			setBonusAura.ExposeToAPL(145004)
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Implemented in windwalker/tigereye_brew.go
			monk := agent.(MonkAgent).GetMonk()
			monk.T16Windwalker4P = monk.RegisterAura(core.Aura{
				Label:    "Focus of Xuen" + monk.Label,
				ActionID: core.ActionID{SpellID: 145024},
				Duration: 10 * time.Second,
			})

			setBonusAura.ExposeToAPL(145022)
		},
	},
})

// T16 - Brewmaster
var ItemSetArmorOfSevenSacredSeals = core.NewItemSet(core.ItemSet{
	Name: "Armor of Seven Sacred Seals",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Not implemented as not having Black Ox statue
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			monk := agent.(MonkAgent).GetMonk()

			monk.T16Brewmaster4P = setBonusAura

			setBonusAura.ExposeToAPL(145056)
		},
	},
})

func init() {
}
