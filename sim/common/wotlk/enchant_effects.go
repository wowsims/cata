package wotlk

import (
	"time"

	"github.com/wowsims/cata/sim/common/shared"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func init() {
	// Keep these in order by item ID.

	// Enchant: 3251, Spell: 44622 - Giant Slayer
	core.NewEnchantEffect(3251, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()

		dpm := character.AutoAttacks.NewDynamicProcManagerForEnchant(3251, 4.0, 0)

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 44622},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 237, spell.OutcomeMagicHitAndCrit)
			},
		})

		aura := core.MakePermanent(character.GetOrRegisterAura(core.Aura{
			Label: "Giant Slayer",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if result.Target.MobType != proto.MobType_MobTypeGiant {
					return
				}

				if dpm.Proc(sim, spell.ProcMask, "Giant Slayer") {
					procSpell.Cast(sim, result.Target)
				}
			},
		}))

		character.ItemSwap.RegisterEnchantProc(3251, aura)
	})

	// Enchant: 3239, Spell: 44525 - Icebreaker
	core.NewEnchantEffect(3239, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()

		dpm := character.AutoAttacks.NewDynamicProcManagerForEnchant(3239, 4.0, 0)

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 44525},
			SpellSchool: core.SpellSchoolFire,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(185, 215), spell.OutcomeMagicHitAndCrit)
			},
		})

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Icebreaker",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if dpm.Proc(sim, spell.ProcMask, "Icebreaker") {
					procSpell.Cast(sim, result.Target)
				}
			},
		})

		character.ItemSwap.RegisterEnchantProc(3239, aura)
	})

	// Enchant: 3607, Spell: 55076, Item: 41146 - Sun Scope
	core.NewEnchantEffect(3607, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()
		// TODO: This should be ranged-only haste. For now just make it hunter-only.
		if character.Class == proto.Class_ClassHunter {
			character.AddStat(stats.HasteRating, 40)
		}
	})

	// Enchant: 3608, Spell: 55135, Item: 41167 - Heartseeker Scope
	core.NewEnchantEffect(3608, func(agent core.Agent, _ proto.ItemLevelState) {
		agent.GetCharacter().AddBonusRangedCritPercent(40 / core.CritRatingPerCritPercent)
	})

	// Enchant: 3748, Spell: 56353, Item: 42500 - Titanium Shield Spike
	shared.NewProcDamageEffect(shared.ProcDamageEffect{
		EnchantID: 3748,
		SpellID:   56353,
		Trigger: core.ProcTrigger{
			Name:     "Titanium Shield Spike",
			Callback: core.CallbackOnSpellHitTaken,
			ProcMask: core.ProcMaskMelee,
			Outcome:  core.OutcomeBlock,
		},
		School:  core.SpellSchoolPhysical,
		MinDmg:  45,
		MaxDmg:  67,
		Outcome: shared.OutcomeMeleeCanCrit,
		IsMelee: true,
	})

	// Enchant: 3247, Spell: 44595 - Scourgebane
	core.NewEnchantEffect(3247, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeUndead {
			character.PseudoStats.MobTypeAttackPower += 140
		}
	})

	// Enchant: 3253, Spell: 44625 - Armsman
	core.NewEnchantEffect(3253, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()
		character.PseudoStats.ThreatMultiplier *= 1.02
	})

	// Enchant: 3296, Spell: 47899 - Wisdom
	core.NewEnchantEffect(3296, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()
		character.PseudoStats.ThreatMultiplier *= 0.98
	})

	// Enchant: 3789, Spell: 59620 - Berserking
	core.NewEnchantEffect(3789, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()

		dpm := character.AutoAttacks.NewDynamicProcManagerForEnchant(3789, 1.0, 0)

		// Modify only gear armor, including from agility
		fivePercentOfArmor := (character.EquipStats()[stats.Armor] + 2.0*character.EquipStats()[stats.Agility]) * 0.05
		procAuraMH := character.NewTemporaryStatsAura("Berserking MH Proc", core.ActionID{SpellID: 59620, Tag: 1}, stats.Stats{stats.AttackPower: 400, stats.RangedAttackPower: 400, stats.Armor: -fivePercentOfArmor}, time.Second*15)
		procAuraOH := character.NewTemporaryStatsAura("Berserking OH Proc", core.ActionID{SpellID: 59620, Tag: 2}, stats.Stats{stats.AttackPower: 400, stats.RangedAttackPower: 400, stats.Armor: -fivePercentOfArmor}, time.Second*15)

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Berserking (Enchant)",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if dpm.Proc(sim, spell.ProcMask, "Berserking") {
					if spell.IsMH() {
						procAuraMH.Activate(sim)
					} else {
						procAuraOH.Activate(sim)
					}
				}
			},
		})

		character.ItemSwap.RegisterEnchantProc(3789, aura)
	})

	// TODO: These are stand-in values without any real reference.
	core.NewEnchantEffect(3241, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()

		dpm := character.AutoAttacks.NewDynamicProcManagerForEnchant(3241, 3.0, 0)

		healthMetrics := character.NewHealthMetrics(core.ActionID{ItemID: 44494})

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Lifeward",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if dpm.Proc(sim, spell.ProcMask, "Lifeward") {
					character.GainHealth(sim, 300*character.PseudoStats.HealingTakenMultiplier, healthMetrics)
				}
			},
		})

		character.ItemSwap.RegisterEnchantProc(3241, aura)
	})

	// Enchant: 3790, Spell: 59630 - Black Magic
	core.NewEnchantEffect(3790, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()

		procAura := character.NewTemporaryStatsAura("Black Magic Proc", core.ActionID{SpellID: 59626}, stats.Stats{stats.HasteRating: 250}, time.Second*10)
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 35,
		}
		procAura.Icd = &icd

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Black Magic",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				// Special case for spells that aren't spells that can proc black magic.
				specialCaseSpell := spell.ActionID.SpellID == 47465 || spell.ActionID.SpellID == 12867

				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskSpellOrSpellProc) && !specialCaseSpell {
					return
				}

				if icd.IsReady(sim) && sim.RandomFloat("Black Magic") < 0.35 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})

		character.ItemSwap.RegisterEnchantProc(3790, aura)
	})

	// Enchant: 3843, Spell: 61471 - Diamond-cut Refractor Scope
	core.NewEnchantEffect(3843, func(agent core.Agent, _ proto.ItemLevelState) {
		w := agent.GetCharacter().AutoAttacks.Ranged()
		w.BaseDamageMin += 15
		w.BaseDamageMax += 15
	})

	//core.NewEnchantEffect(3603, func(agent core.Agent, _ proto.ItemLevelState) {
	//	character := agent.GetCharacter()
	//	actionID := core.ActionID{SpellID: 54757}

	//	spell := character.GetOrRegisterSpell(core.SpellConfig{
	//		ActionID:    actionID,
	//		SpellSchool: core.SpellSchoolFire,
	//		ProcMask:    core.ProcMaskEmpty,
	//		Flags:       core.SpellFlagNoOnCastComplete,

	//		Cast: core.CastConfig{
	//			CD: core.Cooldown{
	//				Timer:    character.NewTimer(),
	//				Duration: time.Second * 45,
	//			},
	//			SharedCD: core.Cooldown{
	//				Timer:    character.GetOffensiveTrinketCD(),
	//				Duration: time.Second * 10,
	//			},
	//		},

	//		DamageMultiplier: 1,
	//		CritMultiplier:   character.GetSpellCritMultiplier(),
	//		ThreatMultiplier: 1,

	//		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	//			spell.CalcAndDealDamage(sim, target, sim.Roll(1654, 2020), spell.OutcomeMagicCrit)
	//		},
	//	})

	//	character.AddMajorCooldown(core.MajorCooldown{
	//		Spell:    spell,
	//		Priority: core.CooldownPriorityLow, // Use low prio so other actives get used first.
	//		Type:     core.CooldownTypeDPS,
	//	})
	//})

	//core.NewEnchantEffect(3604, func(agent core.Agent, _ proto.ItemLevelState) {
	//	character := agent.GetCharacter()
	//	actionID := core.ActionID{SpellID: 54758}

	//	procAura := character.NewTemporaryStatsAura("Hyperspeed Acceleration", actionID, stats.Stats{stats.HasteRating: 340}, time.Second*12)

	//	spell := character.GetOrRegisterSpell(core.SpellConfig{
	//		ActionID:    actionID,
	//		SpellSchool: core.SpellSchoolPhysical,
	//		Flags:       core.SpellFlagNoOnCastComplete,

	//		Cast: core.CastConfig{
	//			CD: core.Cooldown{
	//				Timer:    character.NewTimer(),
	//				Duration: time.Second * 60,
	//			},
	//			// Shared CD with Offensive trinkets has been removed.
	//			// https://twitter.com/AggrendWoW/status/1579664462843633664
	//			// Change possibly temporary, but developers have confirmed it was intended.

	//			// SharedCD: core.Cooldown{
	//			// 	Timer:    character.GetOffensiveTrinketCD(),
	//			// 	Duration: time.Second * 12,
	//			// },
	//		},

	//		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
	//			procAura.Activate(sim)
	//		},
	//	})

	//	character.AddMajorCooldown(core.MajorCooldown{
	//		Spell:    spell,
	//		Priority: core.CooldownPriorityLow, // Use low prio so other actives get used first.
	//		Type:     core.CooldownTypeDPS,
	//	})
	//})

	// Enchant: 3722, Spell: 55642 - Lightweave Embroidery
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Lightweave Embroidery",
		EnchantID: 3722,
		Callback:  core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt,
		ProcMask:  core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:   core.OutcomeLanded,
	})

	// Enchant: 3728, Spell: 55769 - Darkglow Embroidery
	core.NewEnchantEffect(3728, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()
		if !character.HasManaBar() {
			return
		}

		manaMetrics := character.NewManaMetrics(core.ActionID{SpellID: 55767})
		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 45,
		}

		aura := character.GetOrRegisterAura(core.Aura{
			Icd:      &icd,
			ActionID: core.ActionID{SpellID: 55769},
			Label:    "Darkglow Embroidery",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if icd.IsReady(sim) && sim.RandomFloat("Darkglow") < 0.35 {
					icd.Use(sim)
					character.AddMana(sim, 400, manaMetrics)
				}
			},
		})

		character.ItemSwap.RegisterEnchantProc(3728, aura)
	})

	// Enchant: 3730, Spell: 55777 - Swordguard Embroidery
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Swordguard Embroidery",
		EnchantID: 3730,
		Callback:  core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt,
		ProcMask:  core.ProcMaskMeleeOrRanged,
		Outcome:   core.OutcomeLanded,
	})

	// Enchant: 3870, Spell: 64568 - Blood Draining
	core.NewEnchantEffect(3870, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()
		healthMetrics := character.NewHealthMetrics(core.ActionID{SpellID: 64569})

		bloodReserveAura := character.GetOrRegisterAura(core.Aura{
			Label:     "Blood Reserve",
			ActionID:  core.ActionID{SpellID: 64568},
			Duration:  time.Second * 20,
			MaxStacks: 5,
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if character.CurrentHealth()/character.MaxHealth() < 0.35 {
					amountHealed := float64(aura.GetStacks()) * (360. + sim.RandomFloat("Blood Reserve")*80.)
					character.GainHealth(sim, amountHealed*character.PseudoStats.HealingTakenMultiplier, healthMetrics)
					aura.Deactivate(sim)
				}
			},
		})

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Blood Draining",
			Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
			ProcMask:   core.ProcMaskMelee,
			Harmful:    true,
			ProcChance: 0.5,
			ICD:        time.Second * 10,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				if bloodReserveAura.IsActive() {
					bloodReserveAura.Refresh(sim)
					bloodReserveAura.AddStack(sim)
				} else {
					bloodReserveAura.Activate(sim)
					bloodReserveAura.SetStacks(sim, 1)
				}
			},
		})

		character.ItemSwap.RegisterEnchantProc(3870, aura)
	})
}
