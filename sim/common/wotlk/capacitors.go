package wotlk

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

type CapacitorHandler func(*core.Simulation)

type CapacitorAura struct {
	Aura    core.Aura
	Handler CapacitorHandler
}

// Creates an aura which activates a handler function upon reaching a certain number of stacks.
func makeCapacitorAura(unit *core.Unit, config CapacitorAura) *core.Aura {
	handler := config.Handler
	config.Aura.OnStacksChange = func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
		if newStacks == aura.MaxStacks {
			handler(sim)
			aura.SetStacks(sim, 0)
		}
	}
	return unit.RegisterAura(config.Aura)
}

type CapacitorDamageEffect struct {
	Name      string
	ID        int32
	MaxStacks int32
	Trigger   core.ProcTrigger

	School core.SpellSchool
	MinDmg float64
	MaxDmg float64
}

func newCapacitorDamageEffect(config CapacitorDamageEffect) {
	core.NewItemEffect(config.ID, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()

		minDmg := config.MinDmg
		maxDmg := config.MaxDmg
		damageSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: config.ID},
			SpellSchool: config.School,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(minDmg, maxDmg), spell.OutcomeMagicHitAndCrit)
			},
		})

		capacitorAura := makeCapacitorAura(&character.Unit, CapacitorAura{
			Aura: core.Aura{
				Label:     config.Name,
				ActionID:  core.ActionID{ItemID: config.ID},
				Duration:  core.NeverExpires,
				MaxStacks: config.MaxStacks,
			},
			Handler: func(sim *core.Simulation) {
				damageSpell.Cast(sim, character.CurrentTarget)
			},
		})

		config.Trigger.Name = config.Name + " Trigger"
		config.Trigger.Handler = func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			capacitorAura.Activate(sim)
			capacitorAura.AddStack(sim)
		}
		triggerAura := core.MakeProcTriggerAura(&character.Unit, config.Trigger)

		character.ItemSwap.RegisterProc(config.ID, triggerAura)
	})
}

func init() {
	core.AddEffectsToTest = false

	// newCapacitorDamageEffect(CapacitorDamageEffect{
	// 	Name:      "Thunder Capacitor",
	// 	ID:        38072,
	// 	MaxStacks: 4,
	// 	Trigger: core.ProcTrigger{
	// 		Callback: core.CallbackOnSpellHitDealt,
	// 		ProcMask: core.ProcMaskSpellOrProc | core.ProcMaskWeaponProc,
	// 		Outcome:  core.OutcomeCrit,
	// 		ICD:      time.Millisecond * 2500,
	// 		ActionID: core.ActionID{ItemID: 38072},
	// 	},
	// 	School: core.SpellSchoolNature,
	// 	MinDmg: 1181,
	// 	MaxDmg: 1371,
	// })
	// newCapacitorDamageEffect(CapacitorDamageEffect{
	// 	Name:      "Reign of the Unliving",
	// 	ID:        47182,
	// 	MaxStacks: 3,
	// 	Trigger: core.ProcTrigger{
	// 		Callback: core.CallbackOnSpellHitDealt,
	// 		ProcMask: core.ProcMaskSpellOrProc | core.ProcMaskWeaponProc | core.ProcMaskSuppressedProc,
	// 		Outcome:  core.OutcomeCrit,
	// 		ICD:      time.Millisecond * 2000,
	// 		ActionID: core.ActionID{ItemID: 47182},
	// 	},
	// 	School: core.SpellSchoolFire,
	// 	MinDmg: 1741,
	// 	MaxDmg: 2023,
	// })
	// newCapacitorDamageEffect(CapacitorDamageEffect{
	// 	Name:      "Reign of the Unliving H",
	// 	ID:        47188,
	// 	MaxStacks: 3,
	// 	Trigger: core.ProcTrigger{
	// 		Callback: core.CallbackOnSpellHitDealt,
	// 		ProcMask: core.ProcMaskSpellOrProc | core.ProcMaskWeaponProc | core.ProcMaskSuppressedProc,
	// 		Outcome:  core.OutcomeCrit,
	// 		ICD:      time.Millisecond * 2000,
	// 		ActionID: core.ActionID{ItemID: 47188},
	// 	},
	// 	School: core.SpellSchoolFire,
	// 	MinDmg: 1959,
	// 	MaxDmg: 2275,
	// })

	core.AddEffectsToTest = true

	// newCapacitorDamageEffect(CapacitorDamageEffect{
	// 	Name:      "Reign of the Dead",
	// 	ID:        47316,
	// 	MaxStacks: 3,
	// 	Trigger: core.ProcTrigger{
	// 		Callback: core.CallbackOnSpellHitDealt,
	// 		ProcMask: core.ProcMaskSpellOrProc | core.ProcMaskWeaponProc | core.ProcMaskSuppressedProc,
	// 		Outcome:  core.OutcomeCrit,
	// 		ICD:      time.Millisecond * 2000,
	// 		ActionID: core.ActionID{ItemID: 47316},
	// 	},
	// 	School: core.SpellSchoolFire,
	// 	MinDmg: 1741,
	// 	MaxDmg: 2023,
	// })
	// newCapacitorDamageEffect(CapacitorDamageEffect{
	// 	Name:      "Reign of the Dead H",
	// 	ID:        47477,
	// 	MaxStacks: 3,
	// 	Trigger: core.ProcTrigger{
	// 		Callback: core.CallbackOnSpellHitDealt,
	// 		ProcMask: core.ProcMaskSpellOrProc | core.ProcMaskWeaponProc | core.ProcMaskSuppressedProc,
	// 		Outcome:  core.OutcomeCrit,
	// 		ICD:      time.Millisecond * 2000,
	// 		ActionID: core.ActionID{ItemID: 47477},
	// 	},
	// 	School: core.SpellSchoolFire,
	// 	MinDmg: 1959,
	// 	MaxDmg: 2275,
	// })

	// see various posts around https://web.archive.org/web/20100530203708/http://elitistjerks.com/f78/t39136-combat_mutilate_spreadsheets_updated_3_3_a/p96/#post1518212
	NewItemEffectWithHeroic(func(isHeroic bool) {
		name := "Tiny Abomination in a Jar"
		itemID := int32(50351)
		maxStacks := int32(8)
		if isHeroic {
			name += " H"
			itemID = 50706
			maxStacks = 7
		}

		core.NewItemEffect(itemID, func(agent core.Agent, _ proto.ItemLevelState) {
			character := agent.GetCharacter()
			if !character.AutoAttacks.AutoSwingMelee {
				return
			}

			registerSpell := func(spellID int32, procMask core.ProcMask, autoAttackConfig *core.SpellConfig, weaponDamagefn func(sim *core.Simulation, ap float64) float64) *core.Spell {
				return character.GetOrRegisterSpell(core.SpellConfig{
					ActionID:         core.ActionID{SpellID: spellID}, // "Manifest Anger"
					SpellSchool:      core.SpellSchoolPhysical,
					ProcMask:         procMask,
					Flags:            core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete,
					DamageMultiplier: autoAttackConfig.DamageMultiplier * 0.5,
					CritMultiplier:   autoAttackConfig.CritMultiplier,
					ThreatMultiplier: autoAttackConfig.ThreatMultiplier,
					ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
						baseDamage := weaponDamagefn(sim, spell.MeleeAttackPower())
						spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
					},
				})
			}

			var mhSpell *core.Spell
			var ohSpell *core.Spell

			firstProc := core.MainHand

			capacitorAura := makeCapacitorAura(&character.Unit, CapacitorAura{
				Aura: core.Aura{
					Label:     name,
					ActionID:  core.ActionID{SpellID: 71432}, // "Motes of Anger", the aura is either 71406 or 71545 (H) ("Anger Capacitor")
					Duration:  core.NeverExpires,
					MaxStacks: maxStacks,
					OnInit: func(aura *core.Aura, sim *core.Simulation) {
						mhSpell = registerSpell(71433, core.ProcMaskMeleeMHSpecial, character.AutoAttacks.MHConfig(), character.MHWeaponDamage)
						ohSpell = registerSpell(71434, core.ProcMaskMeleeOHSpecial, character.AutoAttacks.OHConfig(), character.OHWeaponDamage)
					},
				},
				Handler: func(sim *core.Simulation) {
					if firstProc == core.MainHand {
						mhSpell.Cast(sim, character.CurrentTarget)
					} else {
						ohSpell.Cast(sim, character.CurrentTarget)
					}
				},
			})

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       name + " Trigger",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeOrMeleeProc,
				Outcome:    core.OutcomeLanded,
				ActionID:   core.ActionID{ItemID: itemID},
				ProcChance: 0.5,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell == mhSpell || spell == ohSpell { // can't proc itself
						return
					}
					hasOffHand := character.OffHand().ID != 0
					if !capacitorAura.IsActive() {
						if spell.ProcMask.Matches(core.ProcMaskMeleeMH | core.ProcMaskMeleeProc) {
							firstProc = core.MainHand
						} else if hasOffHand {
							firstProc = core.OffHand
						}
					}
					if firstProc == core.OffHand && !hasOffHand {
						return
					}

					capacitorAura.Activate(sim)
					capacitorAura.AddStack(sim)
				},
			})

			character.ItemSwap.RegisterProc(itemID, triggerAura)
		})
	})
}
