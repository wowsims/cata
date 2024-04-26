package death_knight

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func init() {
	// Rune of the Nerubian Carapace
	core.NewEnchantEffect(3883, func(agent core.Agent) {
		character := agent.GetCharacter()

		character.MultiplyStat(stats.Armor, 1.02)
		character.MultiplyStat(stats.Stamina, 1.01)
	})

	// Rune of the Stoneskin Gargoyle
	core.NewEnchantEffect(3847, func(agent core.Agent) {
		character := agent.GetCharacter()

		character.MultiplyStat(stats.Armor, 1.04)
		character.MultiplyStat(stats.Stamina, 1.02)
	})

	// Rune of the Swordbreaking
	core.NewEnchantEffect(3594, func(agent core.Agent) {
		character := agent.GetCharacter()

		character.AddStat(stats.Parry, 2*core.ParryRatingPerParryChance)
	})

	// Rune of Swordshattering
	core.NewEnchantEffect(3365, func(agent core.Agent) {
		character := agent.GetCharacter()

		character.AddStat(stats.Parry, 4*core.ParryRatingPerParryChance)
	})

	// Rune of the Spellbreaking
	core.NewEnchantEffect(3595, func(agent core.Agent) {
		// TODO:
		// Add 2% magic deflection
	})

	// Rune of Spellshattering
	core.NewEnchantEffect(3367, func(agent core.Agent) {
		// TODO:
		// Add 4% magic deflection
	})

	// Rune of the Fallen Crusader
	core.NewEnchantEffect(3368, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.GetAura("Rune Of The Fallen Crusader") != nil {
			// Already registerd from one weapon
			return
		}

		procMask := character.GetProcMaskForEnchant(3368)

		rfcAura := character.NewTemporaryStatsAuraWrapped("Rune Of The Fallen Crusader Proc", core.ActionID{SpellID: 53365}, stats.Stats{}, time.Second*15, func(aura *core.Aura) {
			statDep := character.NewDynamicMultiplyStat(stats.Strength, 1.15)

			aura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.EnableDynamicStatDep(sim, statDep)
			})

			aura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.DisableDynamicStatDep(sim, statDep)
			})
		})

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Rune Of The Fallen Crusader",
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeLanded,
			PPM:      2.0,
			ProcMask: procMask,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				rfcAura.Activate(sim)
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(3368, 2.0, aura.Ppmm, aura)
	})

	// Rune of Cinderglacier
	core.NewEnchantEffect(3369, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.GetAura("Rune of Cinderglacier") != nil {
			// Already registerd from one weapon
			return
		}

		cinderMod := character.AddDynamicMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			FloatValue: 0.2,
			ClassMask:  DeathKnightSpellsAll,
			School:     core.SpellSchoolShadow | core.SpellSchoolFrost,
		})

		cinderAura := character.GetOrRegisterAura(core.Aura{
			ActionID:  core.ActionID{SpellID: 53386},
			Label:     "Cinderglacier",
			Duration:  time.Second * 30,
			MaxStacks: 2,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				cinderMod.Activate()
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				cinderMod.Deactivate()
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if spell.ClassSpellMask&DeathKnightSpellsAll == 0 {
					return
				}

				if !spell.SpellSchool.Matches(core.SpellSchoolShadow | core.SpellSchoolFrost) {
					return
				}

				if aura.IsActive() {
					aura.RemoveStack(sim)
				}
			},
		})

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Rune of Cinderglacier",
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeLanded,
			ProcMask: character.GetProcMaskForEnchant(3369),
			PPM:      1.0,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				cinderAura.Activate(sim)
				cinderAura.SetStacks(sim, cinderAura.MaxStacks)
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(3369, 1.0, aura.Ppmm, aura)
	})

	// Rune of Razorice
	core.NewEnchantEffect(3370, func(agent core.Agent) {
		character := agent.GetCharacter()

		if character.GetAura("Razor Frost") != nil {
			// Already registerd from one weapon
			return
		}

		actionID := core.ActionID{SpellID: 50401}

		// Rune of Razorice
		newRazoriceHitSpell := func(character *core.Character, isMH bool) *core.Spell {
			weapon := character.GetMHWeapon()
			if !isMH {
				weapon = character.GetOHWeapon()
			}
			if weapon == nil {
				return nil
			}

			return character.GetOrRegisterSpell(core.SpellConfig{
				ActionID:    actionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
				SpellSchool: core.SpellSchoolFrost,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagNoOnCastComplete,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					dmg := 0.0
					if isMH {
						dmg = spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) * 0.02
					} else {
						dmg = spell.Unit.OHWeaponDamage(sim, spell.MeleeAttackPower()) * 0.02
					}
					spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeAlwaysHit)
				},
			})
		}

		var vulnAuras core.AuraArray

		ddbcHandler := func(sim *core.Simulation, spell *core.Spell, attackTable *core.AttackTable) float64 {
			if spell.ClassSpellMask&DeathKnightSpellsAll == 0 || !spell.SpellSchool.Matches(core.SpellSchoolFrost) {
				return 1.0
			}
			stacks := vulnAuras.Get(attackTable.Defender).GetStacks()
			return 1.0 + 0.02*float64(stacks)
		}

		vulnAuras = character.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
			return target.GetOrRegisterAura(core.Aura{
				Label:     "RuneOfRazoriceVulnerability" + character.Label,
				ActionID:  core.ActionID{SpellID: 51714},
				Duration:  core.NeverExpires,
				MaxStacks: 5,

				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					core.EnableDamageDoneByCaster(DDBC_RuneOfRazorice, DDBC_Total, character.AttackTables[aura.Unit.UnitIndex], ddbcHandler)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					core.DisableDamageDoneByCaster(DDBC_RuneOfRazorice, character.AttackTables[aura.Unit.UnitIndex])
				},
			})
		})

		procMask := character.GetProcMaskForEnchant(3370)
		mhRazoriceSpell := newRazoriceHitSpell(character, true)
		ohRazoriceSpell := newRazoriceHitSpell(character, false)

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Razor Frost",
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeLanded,
			ProcMask: procMask,
			ICD:      time.Millisecond * 8,

			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				vulnAura := vulnAuras.Get(result.Target)
				if spell.IsMH() {
					mhRazoriceSpell.Cast(sim, result.Target)
				} else {
					ohRazoriceSpell.Cast(sim, result.Target)
				}
				vulnAura.Activate(sim)
				vulnAura.AddStack(sim)
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(3370, aura)
	})
}
