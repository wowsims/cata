package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func (hunter *Hunter) ApplyTalents() {
	if hunter.pet != nil {
		hunter.applyFrenzy()
		hunter.registerBestialWrathCD()
		// Todo: BM stuff
		// hunter.pet.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*2*float64(hunter.Talents.Ferocity))
		// hunter.pet.AddStat(stats.SpellCrit, core.CritRatingPerCritChance*2*float64(hunter.Talents.Ferocity))
		// hunter.pet.AddStat(stats.Dodge, 3*core.DodgeRatingPerDodgeChance*float64(hunter.Talents.CatlikeReflexes))
		// hunter.pet.PseudoStats.DamageDealtMultiplier *= 1 + 0.03*float64(hunter.Talents.UnleashedFury)
		// hunter.pet.PseudoStats.MeleeSpeedMultiplier *= 1 + 0.04*float64(hunter.Talents.SerpentsSwiftness)

		hunter.pet.ApplyTalents()
	}

	// hunter.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*1*float64(hunter.Talents.FocusedAim))
	// hunter.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*1*float64(hunter.Talents.KillerInstinct))
	// hunter.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*1*float64(hunter.Talents.MasterMarksman))
	// hunter.AddStat(stats.Parry, core.ParryRatingPerParryChance*1*float64(hunter.Talents.Deflection))
	// hunter.AddStat(stats.Dodge, 1*core.DodgeRatingPerDodgeChance*float64(hunter.Talents.CatlikeReflexes))
	// hunter.PseudoStats.RangedSpeedMultiplier *= 1 + 0.04*float64(hunter.Talents.SerpentsSwiftness)
	// hunter.PseudoStats.DamageTakenMultiplier *= 1 - 0.02*float64(hunter.Talents.SurvivalInstincts)
	// hunter.AutoAttacks.RangedConfig().DamageMultiplier *= hunter.markedForDeathMultiplier()

	// if hunter.Talents.LethalShots > 0 {
	// 	hunter.AddBonusRangedCritRating(1 * float64(hunter.Talents.LethalShots) * core.CritRatingPerCritChance)
	// }
	// if hunter.Talents.RangedWeaponSpecialization > 0 {
	// 	mult := 1 + []float64{0, .01, .03, .05}[hunter.Talents.RangedWeaponSpecialization]
	// 	hunter.OnSpellRegistered(func(spell *core.Spell) {
	// 		if spell.ProcMask.Matches(core.ProcMaskRanged) {
	// 			spell.DamageMultiplier *= mult
	// 		}
	// 	})
	// }

	if hunter.Talents.Pathing > 0 {
		bonus := 0.1*float64(hunter.Talents.Pathing)
		hunter.MultiplyCastSpeed(bonus) //Todo: Should this be attackspeed?
	}

	if hunter.Talents.HunterVsWild > 0 {
		bonus := 0.5 * float64(hunter.Talents.HunterVsWild)
		hunter.MultiplyStat(stats.Stamina, 1+bonus)
	}

	if hunter.Talents.HuntingParty {
		agiBonus := 0.01
		hunter.MultiplyStat(stats.Agility, 1.0+agiBonus)
	}

	if hunter.Talents.KillingStreak > 0 {
		hunter.applyKillingStreak()
	}

	hunter.applyCobraStrikes()
	hunter.applyPiercingShots()
	hunter.applySpiritBond()
	hunter.applyInvigoration()
	hunter.applyGoForTheThroat()
	hunter.applyThrillOfTheHunt()
	hunter.applyTNT()
	hunter.applySniperTraining()
	hunter.applyHuntingParty()
	hunter.applyImprovedSteadyShot()
	hunter.applyFocusFireCD()
	hunter.applyFervorCD()
	hunter.registerReadinessCD()
}


func (hunter *Hunter) applySpiritBond() {
	if hunter.Talents.SpiritBond == 0 || hunter.pet == nil {
		return
	}

	hunter.PseudoStats.HealingTakenMultiplier *= 1 + 0.05*float64(hunter.Talents.SpiritBond)
	hunter.pet.PseudoStats.HealingTakenMultiplier *= 1 + 0.05*float64(hunter.Talents.SpiritBond)

	actionID := core.ActionID{SpellID: 20895}
	healthMultiplier := 0.01 * float64(hunter.Talents.SpiritBond)
	healthMetrics := hunter.NewHealthMetrics(actionID)
	petHealthMetrics := hunter.pet.NewHealthMetrics(actionID)

	hunter.RegisterResetEffect(func(sim *core.Simulation) {
		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
			Period: time.Second * 10,
			OnAction: func(sim *core.Simulation) {
				hunter.GainHealth(sim, hunter.MaxHealth()*healthMultiplier, healthMetrics)
				hunter.pet.GainHealth(sim, hunter.pet.MaxHealth()*healthMultiplier, petHealthMetrics)
			},
		})
	})
}

func (hunter *Hunter) applyInvigoration() {
	if hunter.Talents.Invigoration == 0 || hunter.pet == nil {
		return
	}

	focusMetrics := hunter.NewFocusMetrics(core.ActionID{SpellID: 53253})

	hunter.pet.RegisterAura(core.Aura{
		Label:    "Invigoration",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
				return
			}

			if !result.DidCrit() {
				return
			}

			hunter.AddFocus(sim, 3*float64(hunter.Talents.Invigoration), focusMetrics)
		},
	})
}

func (hunter *Hunter) applyCobraStrikes() {
	if hunter.Talents.CobraStrikes == 0 || hunter.pet == nil {
		return
	}

	actionID := core.ActionID{SpellID: 53260}
	procChance := 0.05 * float64(hunter.Talents.CobraStrikes)

	hunter.pet.CobraStrikesAura = hunter.pet.RegisterAura(core.Aura{
		Label:     "Cobra Strikes",
		ActionID:  actionID,
		Duration:  time.Second * 10,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.pet.focusDump.BonusCritRating += 100 * core.CritRatingPerCritChance
			if hunter.pet.specialAbility != nil {
				hunter.pet.specialAbility.BonusCritRating += 100 * core.CritRatingPerCritChance
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.pet.focusDump.BonusCritRating -= 100 * core.CritRatingPerCritChance
			if hunter.pet.specialAbility != nil {
				hunter.pet.specialAbility.BonusCritRating -= 100 * core.CritRatingPerCritChance
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
				aura.RemoveStack(sim)
			}
		},
	})

	hunter.RegisterAura(core.Aura{
		Label:    "Cobra Strikes",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell != hunter.ArcaneShot { // Only arcane shot, but also can proc on non crits
				return
			}

			if sim.RandomFloat("Cobra Strikes") < procChance {
				hunter.pet.CobraStrikesAura.Activate(sim)
				hunter.pet.CobraStrikesAura.SetStacks(sim, 2)
			}
		},
	})
}

func (hunter *Hunter) applyPiercingShots() {
	if hunter.Talents.PiercingShots == 0 {
		return
	}

	psSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53238},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "PiercingShots",
				Duration: time.Second * 8,
			},
			NumberOfTicks: 8,
			TickLength:    time.Second * 1,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// Specifically account for bleed modifiers, since it still affects the spell, but we're ignoring all modifiers.
				dot.SnapshotAttackerMultiplier = target.PseudoStats.PeriodicPhysicalDamageTakenMultiplier
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).ApplyOrReset(sim)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
		},
	})

	hunter.RegisterAura(core.Aura{
		Label:    "Piercing Shots Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.DidCrit() {
				return
			}
			if spell != hunter.AimedShot && spell != hunter.SteadyShot && spell != hunter.ChimeraShot {
				return
			}

			dot := psSpell.Dot(result.Target)
			outstandingDamage := core.TernaryFloat64(dot.IsActive(), dot.SnapshotBaseDamage*float64(dot.NumberOfTicks-dot.TickCount), 0)
			newDamage := result.Damage * 0.1 * float64(hunter.Talents.PiercingShots)

			dot.SnapshotBaseDamage = (outstandingDamage + newDamage) / float64(dot.NumberOfTicks)
			psSpell.Cast(sim, result.Target)
		},
	})
}

func (hunter *Hunter) applyImprovedSteadyShot() {
	if hunter.Talents.ImprovedSteadyShot == 0 {
		return
	}
		hunter.ImprovedSteadyShotAura = hunter.RegisterAura(core.Aura{
		Label:     "Improved Steady Shot",
		ActionID:  core.ActionID{SpellID: 53221},
		Duration:  time.Second * 8,
		MaxStacks: 1,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			attackspeedMultiplier := 1 + (float64(hunter.Talents.ImprovedSteadyShot) * 0.05)
			aura.Unit.MultiplyRangedSpeed(sim, attackspeedMultiplier)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			attackspeedMultiplier := 1 + (float64(hunter.Talents.ImprovedSteadyShot) * 0.05)
			aura.Unit.MultiplyRangedSpeed(sim, 1 / attackspeedMultiplier)
		},
	})
	hunter.ImprovedSteadyShotAuraCounter = hunter.RegisterAura(core.Aura{
		Label:     "Imp SS Counter",
		Duration:  core.NeverExpires,
		MaxStacks: 1,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ActionID.SpellID == 0 { // Todo: Better way to stop auto attacks from counting?
				return
			}
			if spell != hunter.SteadyShot {
				aura.SetStacks(sim, 0)
			} else {
				if aura.GetStacks() == 1 {
					hunter.ImprovedSteadyShotAura.Activate(sim)
					aura.SetStacks(sim, 0)
				} else {
					aura.SetStacks(sim, 1)
				}
			}
		},
	})
}
func (hunter *Hunter) applyKillingStreak() {
	if hunter.Talents.KillingStreak == 0 {
		return
	}
	hunter.KillingStreakAura = hunter.RegisterAura(core.Aura{
		Label:     "Killing Streak",
		ActionID:  core.ActionID{SpellID: 82748},
		Duration:  core.NeverExpires,
		MaxStacks: 1,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.KillCommand.DamageMultiplier /= 1 + (float64(hunter.Talents.KillingStreak) * 0.1)
			if hunter.Talents.KillingStreak == 1 {
				hunter.KillCommand.CostMultiplier *= 35.0 / 40.0
			} else if hunter.Talents.KillingStreak == 2 {
				hunter.KillCommand.CostMultiplier *= 30.0 / 40.0
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.KillCommand.DamageMultiplier /= 1 + (float64(hunter.Talents.KillingStreak) * 0.1)
			if hunter.Talents.KillingStreak == 1 { //Todo: ??? Static cost reduction?
				hunter.KillCommand.CostMultiplier *= 40.0 / 35.0
			} else if hunter.Talents.KillingStreak == 2 {
				hunter.KillCommand.CostMultiplier *= 40.0 / 30.0
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if hunter.Talents.KillingStreak > 0 {
				if spell == hunter.KillCommand {
					aura.SetStacks(sim, 0)
				}
			}
		},
	})
	hunter.KillingStreakCounterAura = hunter.RegisterAura(core.Aura{
		Label:     "Killing Streak (KC Crit)",
		Duration:  core.NeverExpires,
		MaxStacks: 2,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if hunter.Talents.KillingStreak > 0 {
				if spell == hunter.KillCommand {
					if aura.GetStacks() == 1 && result.DidCrit() {
						hunter.KillingStreakAura.Activate(sim)
						aura.SetStacks(sim, 0)
						return
					}
					if result.DidCrit() {
						aura.Activate(sim)
						aura.AddStack(sim)
					} else if aura.IsActive() {
						aura.SetStacks(sim, 0)
					}
				}
			}
		},
	})
}
func (hunter *Hunter) applyFrenzy() {
	if hunter.Talents.Frenzy == 0 {
		return
	}

	hunter.pet.FrenzyAura = hunter.pet.RegisterAura(core.Aura{
		Label:    "Frenzy",
		Duration: time.Second * 10,
		MaxStacks: 5,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyMeleeSpeed(sim, 1.02)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyMeleeSpeed(sim, 1/1.02)
		},
	})
	hunter.pet.RegisterAura(core.Aura{
		Label: "FrenzyHandler",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
				return
			}
			if hunter.pet.FrenzyAura.IsActive() {
				if hunter.pet.FrenzyAura.GetStacks() != 5 {
					hunter.pet.FrenzyAura.AddStack(sim)
					hunter.pet.FrenzyAura.Refresh(sim)
				}
			} else {
				hunter.pet.FrenzyAura.Activate(sim)
				hunter.pet.FrenzyAura.SetStacks(sim, 1)
			}
		},
	})
}

 func (hunter *Hunter) applyLongevity(dur time.Duration) time.Duration {
 	return time.Duration(float64(dur) * (1.0 - 0.1*float64(hunter.Talents.Longevity)))
 }

func (hunter *Hunter) applyFocusFireCD(){
	if !hunter.Talents.FocusFire {
		return
	}

	actionID := core.ActionID{SpellID: 82692}
	focusFireAura := hunter.RegisterAura(core.Aura{
		Label:    "Focus Fire",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.pet.FrenzyStacksSnapshot = float64(hunter.pet.FrenzyAura.GetStacks())
			if hunter.pet.FrenzyStacksSnapshot >= 1 {
				hunter.pet.FrenzyAura.Deactivate(sim)
				hunter.pet.AddFocus(sim, 4, nil)
				aura.Unit.MultiplyRangedSpeed(sim, 1 + (float64(hunter.pet.FrenzyStacksSnapshot) * 0.03))
				if sim.Log != nil {
					hunter.pet.Log(sim, "Consumed %0f stacks of Frenzy for Focus Fire.", hunter.pet.FrenzyStacksSnapshot)
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.pet.FrenzyStacksSnapshot > 0 {
				aura.Unit.MultiplyRangedSpeed(sim, 1 / (1 + (float64(hunter.pet.FrenzyStacksSnapshot) * 0.03)))
			}
		},
	})

	focusFireSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 1,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 15,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if focusFireAura.IsActive() {
				focusFireAura.Deactivate(sim) // Want to apply new one
			}
			focusFireAura.Activate(sim)
			//focusFireAura.OnGain(focusFireAura, sim)
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: focusFireSpell,
		Type:  core.CooldownTypeDPS,
	})

}
func (hunter *Hunter) applyFervorCD() {
	if !hunter.Talents.Fervor {
		return
	}

	actionID := core.ActionID{SpellID: 82726}
	focusMetrics := hunter.NewFocusMetrics(actionID)
	fervorSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 1,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.AddFocus(sim, 50, focusMetrics)
			hunter.pet.AddFocus(sim, 50, focusMetrics)
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: fervorSpell,
		Type:  core.CooldownTypeDPS,
	})
}
func (hunter *Hunter) registerBestialWrathCD() {
	if !hunter.Talents.BestialWrath {
		return
	}
	if hunter.Talents.TheBeastWithin {
		hunter.PseudoStats.DamageDealtMultiplier *= 1.1
	}

	actionID := core.ActionID{SpellID: 19574}

	bestialWrathPetAura := hunter.pet.RegisterAura(core.Aura{
		Label:    "Bestial Wrath Pet",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.2
		},
	})

	bestialWrathAura := hunter.RegisterAura(core.Aura{
		Label:    "Bestial Wrath",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.1
			aura.Unit.PseudoStats.CostMultiplier -= 0.5
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.1
			aura.Unit.PseudoStats.CostMultiplier += 0.5
		},
	})
	core.RegisterPercentDamageModifierEffect(bestialWrathAura, 1.1)

	bwSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 1,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: hunter.applyLongevity(time.Minute*2 - core.TernaryDuration(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfBestialWrath), time.Second*20, 0)),
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			bestialWrathPetAura.Activate(sim)

			if hunter.Talents.TheBeastWithin {
				bestialWrathAura.Activate(sim)
			}
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: bwSpell,
		Type:  core.CooldownTypeDPS,
	})
}

func (hunter *Hunter) applyGoForTheThroat() {
	if hunter.Talents.GoForTheThroat == 0 {
		return
	}
	if hunter.pet == nil {
		return
	}

	spellID := []int32{0, 34950, 34950}[hunter.Talents.GoForTheThroat]
	focusMetrics := hunter.NewFocusMetrics(core.ActionID{SpellID: spellID})

	amount := 5 * float64(hunter.Talents.GoForTheThroat)

	hunter.RegisterAura(core.Aura{
		Label:    "Go for the Throat",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskRangedAuto) || !result.DidCrit() {
				return
			}
			if !hunter.pet.IsEnabled() {
				return
			}
			hunter.pet.AddFocus(sim, amount, focusMetrics)
		},
	})
}

//Todo: Should we support precasting freezing/ice trap?
func (hunter *Hunter) applyTNT() {
	if hunter.Talents.TNT == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 56343}
	procChance := []float64{0, 0.10, 0.20}[hunter.Talents.TNT]

	icd := core.Cooldown{
		Timer:    hunter.NewTimer(),
		Duration: time.Second * 10,
	}

	hunter.LockAndLoadAura = hunter.RegisterAura(core.Aura{
		Icd:       &icd,
		Label:     "Lock and Load Proc",
		ActionID:  actionID,
		Duration:  time.Second * 12,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.ExplosiveShot != nil {
				hunter.ExplosiveShot.CostMultiplier -= 1
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.ExplosiveShot != nil {
				hunter.ExplosiveShot.CostMultiplier += 1
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == hunter.ExplosiveShot {
				aura.RemoveStack(sim)
			}
		},
	})

	hunter.RegisterAura(core.Aura{
		Label:    "TNT Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell != hunter.BlackArrow && spell != hunter.ExplosiveTrap {
				return
			}

			if !icd.IsReady(sim) {
				return
			}

			if sim.RandomFloat("Lock and Load") < procChance {
				icd.Use(sim)
				hunter.LockAndLoadAura.Activate(sim)
				hunter.LockAndLoadAura.SetStacks(sim, 2)
			}
		},
	})
}

func (hunter *Hunter) applyThrillOfTheHunt() {
	if hunter.Talents.ThrillOfTheHunt == 0 {
		return
	}

	procChance := float64(hunter.Talents.ThrillOfTheHunt) * 5
	focusMetrics := hunter.NewFocusMetrics(core.ActionID{SpellID: 34499})

	hunter.RegisterAura(core.Aura{
		Label:    "Thrill of the Hunt",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// mask 256
			if !spell.ProcMask.Matches(core.ProcMaskRangedSpecial) {
				return
			}

			if sim.Proc(procChance, "ThrillOfTheHunt") {
				hunter.AddFocus(sim, spell.CurCast.Cost * 0.4, focusMetrics)
			}
		},
	})
}


func (hunter *Hunter) applySniperTraining() {
	if hunter.Talents.SniperTraining == 0 {
		return
	}

	uptime := hunter.SurvivalOptions.SniperTrainingUptime
	if uptime <= 0 {
		return
	}
	uptime = min(1, uptime)

	dmgMod := .02 * float64(hunter.Talents.SniperTraining)

	stAura := hunter.RegisterAura(core.Aura{
		Label:    "Sniper Training",
		ActionID: core.ActionID{SpellID: 53304},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.SteadyShot != nil {
				hunter.SteadyShot.DamageMultiplierAdditive += dmgMod
			}
			if hunter.CobraShot != nil {
				hunter.CobraShot.DamageMultiplierAdditive += dmgMod
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.SteadyShot != nil {
				hunter.SteadyShot.DamageMultiplierAdditive -= dmgMod
			}
			if hunter.CobraShot != nil {
				hunter.CobraShot.DamageMultiplierAdditive -= dmgMod
			}
		},
	})

	core.ApplyFixedUptimeAura(stAura, uptime, time.Second*15, 1)
}

func (hunter *Hunter) applyHuntingParty() {
	if !hunter.Talents.HuntingParty {
		return
	}


	hunter.RegisterAura(core.Aura{
		Label:    "Hunting Party",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, 1.10)
		},
	})
}

func (hunter *Hunter) registerReadinessCD() {
	if !hunter.Talents.Readiness {
		return
	}

	actionID := core.ActionID{SpellID: 23989}

	readinessSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second * 1,
			},
			IgnoreHaste: true, // Hunter GCD is locked
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Minute * 3,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Don't use if there are no cooldowns to reset.
			return !hunter.RapidFire.IsReady(sim)
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.RapidFire.CD.Reset()
			hunter.MultiShot.CD.Reset()
			hunter.KillShot.CD.Reset()
			hunter.RaptorStrike.CD.Reset()
			hunter.ExplosiveTrap.CD.Reset()
			if hunter.KillCommand != nil {
				hunter.KillCommand.CD.Reset()
			}
			if hunter.AimedShot != nil {
				hunter.AimedShot.CD.Reset()
			}
			if hunter.SilencingShot != nil {
				hunter.SilencingShot.CD.Reset()
			}
			if hunter.ChimeraShot != nil {
				hunter.ChimeraShot.CD.Reset()
			}
			if hunter.BlackArrow != nil {
				hunter.BlackArrow.CD.Reset()
			}

			// TODO: This is needed because there are edge cases where core doesn't re-use Rapid Fire.
			// Fix core so this isn't necessary.
			core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt: sim.CurrentTime + 1,
				OnAction: func(_ *core.Simulation) {
					hunter.UpdateMajorCooldowns()
				},
			})
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: readinessSpell,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			// If RF is about to become ready naturally, wait so we can get 2x usages.
			if !hunter.RapidFire.IsReady(sim) && hunter.RapidFire.TimeToReady(sim) < time.Second*10 {
				return false
			}
			return !hunter.RapidFireAura.IsActive() || hunter.RapidFireAura.RemainingDuration(sim) < time.Second*10
		},
	})
}
