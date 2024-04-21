package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func (hunter *Hunter) ApplyTalents() {
	hunter.EnableArmorSpecialization(stats.Agility, proto.ArmorType_ArmorTypeMail)
	if hunter.Pet != nil {
		hunter.applyFrenzy()
		hunter.registerBestialWrathCD()
		// Todo: BM stuff
		hunter.Pet.ApplyTalents()
	}

	if hunter.Talents.Pathing > 0 {
		bonus := 0.01 * float64(hunter.Talents.Pathing)
		hunter.PseudoStats.RangedSpeedMultiplier *= 1 + bonus
	}

	if hunter.Talents.HunterVsWild > 0 {
		bonus := 0.05 * float64(hunter.Talents.HunterVsWild)
		hunter.MultiplyStat(stats.Stamina, 1+bonus)
	}

	if hunter.Talents.HuntingParty {
		agiBonus := 0.02
		hunter.MultiplyStat(stats.Agility, 1.0+agiBonus)
	}

	if hunter.Talents.KillingStreak > 0 {
		hunter.applyKillingStreak()
	}

	if hunter.Talents.Efficiency > 0 {
		hunter.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_PowerCost_Flat,
			ClassMask:  HunterSpellArcaneShot,
			FloatValue: -float64(hunter.Talents.Efficiency),
		})
		hunter.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_PowerCost_Flat,
			ClassMask:  HunterSpellExplosiveShot | HunterSpellChimeraShot,
			FloatValue: -(float64(hunter.Talents.Efficiency) * 2),
		})
	}
	hunter.registerSicEm()
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
	hunter.applyMasterMarksman()
	hunter.applyTermination()
}
func (hunter *Hunter) applyMasterMarksman() {
	if hunter.Talents.MasterMarksman == 0 {
		return
	}

	procChance := float64(hunter.Talents.MasterMarksman) * 0.2
	hunter.MasterMarksmanAura = hunter.RegisterAura(core.Aura{
		Label:    "Ready, Set, Aim...",
		ActionID: core.ActionID{SpellID: 82925},
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.AimedShot != nil {
				hunter.AimedShot.CostMultiplier = 0
				hunter.AimedShot.DefaultCast.CastTime = 0
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.AimedShot != nil {
				hunter.AimedShot.CostMultiplier = 1
				hunter.AimedShot.DefaultCast.CastTime = time.Second * 3
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == hunter.AimedShot {
				hunter.MasterMarksmanCounterAura.SetStacks(sim, 0)
				hunter.MasterMarksmanCounterAura.Activate(sim)
				aura.Deactivate(sim) // Consume effect
			}

		},
	})
	hunter.MasterMarksmanCounterAura = hunter.RegisterAura(core.Aura{
		Label:     "Master Marksman",
		Duration:  time.Second * 30,
		ActionID:  core.ActionID{SpellID: 34486},
		MaxStacks: 4,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell != hunter.SteadyShot {
				return
			}
			if procChance == 1 || sim.Proc(procChance, "Master Marksman Proc") {
				if aura.GetStacks() == 4 {
					hunter.MasterMarksmanAura.Activate(sim)
				} else {
					aura.AddStack(sim)
				}
			}
		},
	})
}

func (hunter *Hunter) applySpiritBond() {
	if hunter.Talents.SpiritBond == 0 || hunter.Pet == nil {
		return
	}

	hunter.PseudoStats.HealingTakenMultiplier *= 1 + 0.05*float64(hunter.Talents.SpiritBond)
	hunter.Pet.PseudoStats.HealingTakenMultiplier *= 1 + 0.05*float64(hunter.Talents.SpiritBond)

	actionID := core.ActionID{SpellID: 20895}
	healthMultiplier := 0.01 * float64(hunter.Talents.SpiritBond)
	healthMetrics := hunter.NewHealthMetrics(actionID)
	petHealthMetrics := hunter.Pet.NewHealthMetrics(actionID)

	hunter.RegisterResetEffect(func(sim *core.Simulation) {
		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
			Period: time.Second * 10,
			OnAction: func(sim *core.Simulation) {
				hunter.GainHealth(sim, hunter.MaxHealth()*healthMultiplier, healthMetrics)
				hunter.Pet.GainHealth(sim, hunter.Pet.MaxHealth()*healthMultiplier, petHealthMetrics)
			},
		})
	})
}

func (hunter *Hunter) applyInvigoration() {
	if hunter.Talents.Invigoration == 0 || hunter.Pet == nil {
		return
	}

	focusMetrics := hunter.NewFocusMetrics(core.ActionID{SpellID: 53253})

	hunter.Pet.RegisterAura(core.Aura{
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
	if hunter.Talents.CobraStrikes == 0 || hunter.Pet == nil {
		return
	}

	actionID := core.ActionID{SpellID: 53260}
	procChance := 0.05 * float64(hunter.Talents.CobraStrikes)

	hunter.Pet.CobraStrikesAura = hunter.Pet.RegisterAura(core.Aura{
		Label:     "Cobra Strikes",
		ActionID:  actionID,
		Duration:  time.Second * 10,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.Pet.focusDump.BonusCritRating += 100 * core.CritRatingPerCritChance
			if hunter.Pet.specialAbility != nil {
				hunter.Pet.specialAbility.BonusCritRating += 100 * core.CritRatingPerCritChance
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.Pet.focusDump.BonusCritRating -= 100 * core.CritRatingPerCritChance
			if hunter.Pet.specialAbility != nil {
				hunter.Pet.specialAbility.BonusCritRating -= 100 * core.CritRatingPerCritChance
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
				hunter.Pet.CobraStrikesAura.Activate(sim)
				hunter.Pet.CobraStrikesAura.SetStacks(sim, 2)
			}
		},
	})
}
func (hunter *Hunter) applyTermination() {
	if hunter.Talents.Termination == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 83490}

	focusMetrics := hunter.NewFocusMetrics(actionID)
	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label: "Termination",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.IsExecutePhase25() && spell == hunter.SteadyShot || spell == hunter.CobraShot {
				hunter.AddFocus(sim, float64(hunter.Talents.Termination)*3, focusMetrics)
			}
		},
	}))
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

	attackspeedMultiplier := 1 + (float64(hunter.Talents.ImprovedSteadyShot) * 0.05)
	hunter.ImprovedSteadyShotAura = hunter.RegisterAura(core.Aura{
		Label:     "Improved Steady Shot",
		ActionID:  core.ActionID{SpellID: 53221, Tag: 1},
		Duration:  time.Second * 8,
		MaxStacks: 1,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyRangedSpeed(sim, attackspeedMultiplier)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyRangedSpeed(sim, 1/attackspeedMultiplier)
		},
	})
	hunter.ImprovedSteadyShotAuraCounter = core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label:     "Imp SS Counter",
		ActionID:  core.ActionID{SpellID: 53221, Tag: 2},
		MaxStacks: 2,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskRangedAuto) || spell.ActionID.SpellID == 0 || !spell.Flags.Matches(core.SpellFlagAPL) {
				return
			}
			if spell != hunter.SteadyShot {
				aura.SetStacks(sim, 1)
			} else {
				if aura.GetStacks() == 2 {
					hunter.ImprovedSteadyShotAura.Activate(sim)
					aura.SetStacks(sim, 1)
				} else {
					aura.SetStacks(sim, 2)
				}
			}
		},
	}))
}
func (hunter *Hunter) applyKillingStreak() {
	if hunter.Talents.KillingStreak == 0 {
		return
	}
	damageMod := hunter.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  HunterSpellKillCommand,
		FloatValue: float64(hunter.Talents.KillingStreak) * 0.1,
	})
	costMod := hunter.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Flat,
		ClassMask:  HunterSpellKillCommand,
		FloatValue: -(float64(hunter.Talents.KillingStreak) * 5),
	})
	hunter.KillingStreakAura = hunter.RegisterAura(core.Aura{
		Label:    "Killing Streak",
		ActionID: core.ActionID{SpellID: 82748},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
			costMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Deactivate()
			costMod.Deactivate()
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == hunter.KillCommand {
				aura.Deactivate(sim)
			}
		},
	})
	hunter.KillingStreakCounterAura = hunter.RegisterAura(core.Aura{
		Label:     "Killing Streak (KC Crit)",
		Duration:  core.NeverExpires,
		MaxStacks: 2,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == hunter.KillCommand {
				if aura.GetStacks() == 2 && result.DidCrit() {
					hunter.KillingStreakAura.Activate(sim)
					aura.SetStacks(sim, 1)
					return
				}
				if result.DidCrit() {
					aura.AddStack(sim)
				}
			}
		},
	})
}
func (hunter *Hunter) applyFrenzy() {
	if hunter.Talents.Frenzy == 0 {
		return
	}
	actionID := core.ActionID{SpellID: 19622}
	hunter.Pet.FrenzyAura = hunter.Pet.RegisterAura(core.Aura{
		Label:     "Frenzy",
		Duration:  time.Second * 10,
		ActionID:  actionID,
		MaxStacks: 5,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyMeleeSpeed(sim, 1.02)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyMeleeSpeed(sim, 1/1.02)
		},
	})

	hunter.Pet.RegisterAura(core.Aura{
		Label:    "FrenzyHandler",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
				return
			}
			if hunter.Pet.FrenzyAura.IsActive() {
				if hunter.Pet.FrenzyAura.GetStacks() != 5 {
					hunter.Pet.FrenzyAura.AddStack(sim)
					hunter.Pet.FrenzyAura.Refresh(sim)
				}
			} else {
				hunter.Pet.FrenzyAura.Activate(sim)
				hunter.Pet.FrenzyAura.SetStacks(sim, 1)
			}
		},
	})
}

func (hunter *Hunter) applyLongevity(dur time.Duration) time.Duration {
	return time.Duration(float64(dur) * (1.0 - 0.1*float64(hunter.Talents.Longevity)))
}

func (hunter *Hunter) applyFocusFireCD() {
	if !hunter.Talents.FocusFire {
		return
	}

	actionID := core.ActionID{SpellID: 82692}
	petFocusMetrics := hunter.Pet.NewFocusMetrics(actionID)
	focusFireAura := hunter.RegisterAura(core.Aura{
		Label:    "Focus Fire",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.Pet.FrenzyStacksSnapshot = float64(hunter.Pet.FrenzyAura.GetStacks())
			if hunter.Pet.FrenzyStacksSnapshot >= 1 {
				hunter.Pet.FrenzyAura.Deactivate(sim)
				hunter.Pet.AddFocus(sim, 4, petFocusMetrics)
				aura.Unit.MultiplyRangedSpeed(sim, 1+(float64(hunter.Pet.FrenzyStacksSnapshot)*0.03))
				if sim.Log != nil {
					hunter.Pet.Log(sim, "Consumed %0f stacks of Frenzy for Focus Fire.", hunter.Pet.FrenzyStacksSnapshot)
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if hunter.Pet.FrenzyStacksSnapshot > 0 {
				aura.Unit.MultiplyRangedSpeed(sim, 1/(1+(float64(hunter.Pet.FrenzyStacksSnapshot)*0.03)))
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
			hunter.Pet.AddFocus(sim, 50, focusMetrics)
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

	bestialWrathPetAura := hunter.Pet.RegisterAura(core.Aura{
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
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.2
		},
	})
	core.RegisterPercentDamageModifierEffect(bestialWrathAura, 1.2)

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
	if hunter.Pet == nil {
		return
	}

	focusMetrics := hunter.NewFocusMetrics(core.ActionID{SpellID: 34950})

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
			if !hunter.Pet.IsEnabled() {
				return
			}
			hunter.Pet.AddFocus(sim, amount, focusMetrics)
		},
	})
}

// Todo: Should we support precasting freezing/ice trap?
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

				hunter.ExplosiveShot.CD.Reset()

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
				if hunter.ExplosiveShot != nil {
					hunter.ExplosiveShot.CD.Reset()
				}
			}
		},
	})
}

func (hunter *Hunter) applyThrillOfTheHunt() {
	if hunter.Talents.ThrillOfTheHunt == 0 {
		return
	}

	procChance := float64(hunter.Talents.ThrillOfTheHunt) * 0.05
	focusMetrics := hunter.NewFocusMetrics(core.ActionID{SpellID: 34499})

	hunter.RegisterAura(core.Aura{
		Label:    "Thrill of the Hunt",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// mask 256
			if spell == hunter.ArcaneShot || spell == hunter.ExplosiveShot || spell == hunter.BlackArrow {
				if sim.Proc(procChance, "ThrillOfTheHunt") {
					hunter.AddFocus(sim, spell.CurCast.Cost*0.4, focusMetrics)
				}
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

	dmgMod := hunter.AddDynamicMod(core.SpellModConfig{
		ClassMask:  HunterSpellCobraShot | HunterSpellSteadyShot,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: .02 * float64(hunter.Talents.SniperTraining),
	})

	stAura := hunter.RegisterAura(core.Aura{
		Label:    "Sniper Training",
		ActionID: core.ActionID{SpellID: 53304},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dmgMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dmgMod.Deactivate()
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
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, 1/1.10)
		},
	})
}
func (hunter *Hunter) registerSicEm() {
	if hunter.Talents.SicEm == 0 {
		return
	}

	actionId := core.ActionID{SpellID: 83356}
	sicEmMod := hunter.Pet.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -(float64(hunter.Talents.SicEm) * 0.5),
		ProcMask:   core.ProcMaskMeleeMHSpecial,
	})

	core.MakePermanent(hunter.RegisterAura(core.Aura{
		Label:    "Sic'Em Mod",
		ActionID: actionId,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == hunter.ArcaneShot || spell == hunter.AimedShot || spell == hunter.ExplosiveShot {
				if result.DidCrit() {
					sicEmMod.Activate()
				}
			}
		},
	}))
	core.MakePermanent(hunter.Pet.RegisterAura(core.Aura{
		ActionID: actionId,
		Label:    "Sic'Em",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask == core.ProcMaskMeleeMHSpecial {
				if sicEmMod.IsActive {
					sicEmMod.Deactivate()
				}
			}
		},
	}))

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
			hunter.KillShot.CD.Reset()
			hunter.RaptorStrike.CD.Reset()
			hunter.ExplosiveTrap.CD.Reset()
			if hunter.KillCommand != nil {
				hunter.KillCommand.CD.Reset()
			}
			if hunter.ChimeraShot != nil {
				hunter.ChimeraShot.CD.Reset()
			}
			if hunter.BlackArrow != nil {
				hunter.BlackArrow.CD.Reset()
			}
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
