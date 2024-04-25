package death_knight

import (
	//"github.com/wowsims/cata/sim/core/proto"

	//"time"

	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (dk *DeathKnight) ApplyBloodTalents() {
	// Butchery
	dk.applyButchery()

	// Blade barrier
	if dk.Talents.BladeBarrier > 0 {
		dk.PseudoStats.DamageTakenMultiplier *= 1.0 - 0.02*float64(dk.Talents.BladeBarrier)
	}

	if dk.Talents.ImprovedBloodTap > 0 {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: time.Second * -15 * time.Duration(dk.Talents.ImprovedBloodTap),
			ClassMask: DeathKnightSpellBloodTap,
		})
	}

	// Bladed Armor
	if dk.Talents.BladedArmor > 0 {
		coeff := float64(dk.Talents.BladedArmor) * 2
		dk.AddStatDependency(stats.Armor, stats.AttackPower, coeff/180.0)
	}

	// Scent of Blood
	dk.applyScentOfBlood()

	// Scarlet Fever
	dk.applyScarletFever()

	// Blood-Caked Blade
	dk.applyBloodCakedBlade()

	//Toughness
	if dk.Talents.Toughness > 0 {
		dk.ApplyEquipScaling(stats.Armor, []float64{1.0, 1.03, 1.06, 1.1}[dk.Talents.Toughness])
	}

	// Abomination's Might
	if dk.Talents.AbominationsMight > 0 {
		strengthCoeff := 0.01 * float64(dk.Talents.AbominationsMight)
		dk.MultiplyStat(stats.Strength, 1.0+strengthCoeff)
	}

	// Sanguine Fortitude
	if dk.Talents.SanguineFortitude > 0 {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_RunicPowerCost_Flat,
			ClassMask:  DeathKnightSpellIceboundFortitude,
			FloatValue: -10.0 * float64(dk.Talents.SanguineFortitude),
		})
	}

	// Blood Parasite
	dk.applyBloodworms()

	// Will of the Necropolis
	dk.applyWillOfTheNecropolis()

	// Improved Death Strike
	if dk.Talents.ImprovedDeathStrike > 0 {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  DeathKnightSpellDeathStrike,
			FloatValue: 0.4 * float64(dk.Talents.ImprovedDeathStrike),
		})
		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  DeathKnightSpellDeathStrikeHeal,
			FloatValue: 0.15 * float64(dk.Talents.ImprovedDeathStrike),
		})

		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusCrit_Rating,
			ClassMask:  DeathKnightSpellDeathStrike,
			FloatValue: 10 * core.CritRatingPerCritChance * float64(dk.Talents.ImprovedDeathStrike),
		})
	}
}

func (dk *DeathKnight) applyScarletFever() {
	if dk.Talents.ScarletFever == 0 {
		return
	}

	dk.ScarletFeverAura = dk.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.ScarletFeverAura(dk.GetCharacter(), target, dk.Talents.ScarletFever, dk.Talents.Epidemic)
	})
	dk.Env.RegisterPreFinalizeEffect(func() {
		dk.BloodPlagueSpell.RelatedAuras = append(dk.BloodPlagueSpell.RelatedAuras, dk.ScarletFeverAura)
	})

	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:           "Scarlet Fever Activate",
		Callback:       core.CallbackOnApplyEffects,
		ClassSpellMask: DeathKnightSpellBloodPlague,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			dk.ScarletFeverAura.Get(result.Target).Activate(sim)
		},
	})
}

func (dk *DeathKnight) applyBloodCakedBlade() {
	if dk.Talents.BloodCakedBlade == 0 {
		return
	}

	procChance := float64(dk.Talents.BloodCakedBlade) * 0.10
	bloodCakedBladeHitMh := dk.bloodCakedBladeHit(true)

	var bloodCakedBladeHitOh *core.Spell
	if dk.HasOHWeapon() {
		bloodCakedBladeHitOh = dk.bloodCakedBladeHit(false)
	}

	core.MakePermanent(dk.RegisterAura(core.Aura{
		Label: "Blood-Caked Blade",
		// ActionID: core.ActionID{SpellID: 49628}, // Hide from metrics
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Damage == 0 || !spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				return
			}

			if sim.RandomFloat("Blood-Caked Blade Roll") < procChance {
				isMh := spell.ProcMask.Matches(core.ProcMaskMeleeMHAuto)
				if isMh {
					bloodCakedBladeHitMh.Cast(sim, result.Target)
				} else {
					bloodCakedBladeHitOh.Cast(sim, result.Target)
				}
			}
		},
	}))
}

func (dk *DeathKnight) bloodCakedBladeHit(isMh bool) *core.Spell {
	return dk.Unit.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 50463}.WithTag(core.TernaryInt32(isMh, 1, 2)),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskProc,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagMeleeMetrics,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage float64
			if isMh {
				baseDamage = spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
			} else {
				baseDamage = spell.Unit.OHWeaponDamage(sim, spell.MeleeAttackPower())
			}
			baseDamage *= dk.GetDiseaseMulti(target, 0.25, 0.125)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialNoCrit)
		},
	})
}

func (dk *DeathKnight) applyWillOfTheNecropolis() {
	if dk.Talents.WillOfTheNecropolis == 0 {
		return
	}

	damageMit := 1.0 - []float64{0.0, 0.06, 0.16, 0.25}[dk.Talents.WillOfTheNecropolis]

	actionID := core.ActionID{SpellID: 96171}
	wotnAura := dk.RegisterAura(core.Aura{
		Label:    "Will of The Necropolis Proc",
		ActionID: actionID,
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier *= damageMit
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier /= damageMit
		},
	})

	runeTapMod := dk.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		ClassMask:  DeathKnightSpellRuneTap,
		FloatValue: -1,
	})

	core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
		Label: "Will of The Necropolis",
		Icd: &core.Cooldown{
			Timer:    dk.NewTimer(),
			Duration: time.Second * 45,
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !aura.Icd.IsReady(sim) {
				return
			}

			if dk.CurrentHealthPercent() <= 0.3 {
				aura.Icd.Use(sim)
				wotnAura.Activate(sim)
				runeTapMod.Activate()
				dk.GetSpell(RuneTapActionID).CD.Reset()
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask != DeathKnightSpellRuneTap {
				return
			}

			if spell.CurCast.Cost > 0 {
				return
			}

			runeTapMod.Deactivate()
		},
	}))
}

func (dk *DeathKnight) applyScentOfBlood() {
	if dk.Talents.ScentOfBlood == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 49509}
	procChance := 0.15

	rpMetrics := dk.NewRunicPowerMetrics(actionID)

	procAura := dk.RegisterAura(core.Aura{
		Label:     "Scent of Blood Proc",
		ActionID:  actionID,
		Duration:  core.NeverExpires,
		MaxStacks: dk.Talents.ScentOfBlood,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			dk.AddRunicPower(sim, 10.0, rpMetrics)
			aura.RemoveStack(sim)
		},
	})

	core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
		Label: "Scent of Blood",
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.Proc(procChance, "Scent Of Blood Proc Chance") {
				procAura.Activate(sim)
				procAura.SetStacks(sim, procAura.MaxStacks)
			}
		},
	}))
}

func (dk *DeathKnight) applyButchery() {
	if dk.Talents.Butchery == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 49483}
	amountOfRunicPower := 1.0 * float64(dk.Talents.Butchery)
	rpMetrics := dk.NewRunicPowerMetrics(actionID)

	var pa *core.PendingAction
	core.MakePermanent(dk.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Butchery",
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			pa = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period: time.Second * 5,
				OnAction: func(sim *core.Simulation) {
					dk.AddRunicPower(sim, amountOfRunicPower, rpMetrics)
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			pa.Cancel(sim)
		},
	}))
}

func (dk *DeathKnight) applyBloodworms() {
	// if dk.Talents.BloodParasite == 0 {
	// 	return
	// }

	// procChance := 0.03 * float64(dk.Talents.BloodParasite)
	// icd := core.Cooldown{
	// 	Timer:    dk.NewTimer(),
	// 	Duration: time.Second * 20,
	// }

	// // For tracking purposes
	// procSpell := dk.RegisterSpell(core.SpellConfig{
	// 	ActionID: core.ActionID{SpellID: 49543},
	// 	ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
	// 		// Summon Bloodworms
	// 		random := int(math.Round(sim.RandomFloat("Bloodworms count")*2.0)) + 2
	// 		for i := 0; i < random; i++ {
	// 			dk.Bloodworm[i].EnableWithTimeout(sim, dk.Bloodworm[i], time.Second*20)
	// 			dk.Bloodworm[i].CancelGCDTimer(sim)
	// 		}
	// 	},
	// })

	// core.MakePermanent(dk.RegisterAura(core.Aura{
	// 	Label: "Bloodworms Proc",
	// 	OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// 		if !spell.ProcMask.Matches(core.ProcMaskMelee) {
	// 			return
	// 		}

	// 		if !icd.IsReady(sim) {
	// 			return
	// 		}

	// 		if sim.RandomFloat("Bloodworms proc") < procChance {
	// 			icd.Use(sim)
	// 			procSpell.Cast(sim, result.Target)
	// 		}
	// 	},
	// }))
}
