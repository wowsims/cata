package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (mage *Mage) ApplyTalents() {
	// Ordered row-by-row, left to right
	// mage.applyArcaneConcentration
	mage.applyMasterOfElements()
	// mage.applyEarlyWinter()
	mage.applyIgnite()
	// mage.applyImpact()
	// mage.applyIceFloes()
	// mage.applyPiercingChill()
	// mage.applyPermaFrost
	// mage.applyArcaneFlows()
	// mage.applyFingersOfFrost()
	// mage.applyImprovedFreeze()
	// mage.applyEnduringWinter()
	// mage.applyColdSnap()
	mage.applyBrainFreeze()
	// mage.applyArcanePotency()
	// mage.applyImprovedFlamestrike()
	// mage.apllyMoltenFury()
	// mage.applyFocusMagic()
	// mage.applyImprovedManaGem()
	// mage.applyPyromaniac()
	// mage.applyCriticalMass()
	// mage.applyFrostfireOrb()
	// mage.applyDeepFreeze()

	mage.applyArcaneMissileProc()
	mage.applyHotStreak()

	mage.registerArcanePowerCD()
	mage.registerPresenceOfMindCD()
	//mage.registerCombustionCD()
	mage.registerIcyVeinsCD()
	mage.registerColdSnapCD()

	// Stat Buffs
	// mage.PseudoStats.SpiritRegenRateCasting += float64(mage.Talents.ArcaneMeditation) / 6
	// if mage.Talents.StudentOfTheMind > 0 {
	// 	mage.MultiplyStat(stats.Spirit, 1.0+[]float64{0, .04, .07, .10}[mage.Talents.StudentOfTheMind])
	// }
	// if mage.Talents.ArcaneMind > 0 {
	// 	mage.MultiplyStat(stats.Intellect, 1.0+0.03*float64(mage.Talents.ArcaneMind))
	// }
	// if mage.Talents.MindMastery > 0 {
	// 	mage.AddStatDependency(stats.Intellect, stats.SpellPower, 0.03*float64(mage.Talents.MindMastery))
	// }
	// mage.AddStat(stats.SpellCrit, float64(mage.Talents.ArcaneInstability)*1*core.CritRatingPerCritChance)
	// mage.PseudoStats.DamageDealtMultiplier *= 1 + .01*float64(mage.Talents.ArcaneInstability)
	// mage.PseudoStats.DamageDealtMultiplier *= 1 + .01*float64(mage.Talents.PlayingWithFire)
	mage.PseudoStats.CastSpeedMultiplier *= 1 + .01*float64(mage.Talents.NetherwindPresence)

	mage.AddStat(stats.SpellCrit, 0.01*float64(mage.Talents.PiercingIce)*core.CritRatingPerCritChance)
	// mage.PseudoStats.SpiritRegenRateCasting += float64(mage.Talents.Pyromaniac) / 6

	// mage.AddStat(stats.SpellHit, float64(mage.Talents.Precision)*core.SpellHitRatingPerHitChance)
	// mage.PseudoStats.CostMultiplier *= 1 - .01*float64(mage.Talents.Precision)

	// mage.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] *= 1 + .01*float64(mage.Talents.ArcticWinds)
	// mage.PseudoStats.CostMultiplier *= 1 - .04*float64(mage.Talents.FrostChanneling)

	// magicAbsorptionBonus := 2 * float64(mage.Talents.MagicAbsorption)
	// mage.AddStat(stats.ArcaneResistance, magicAbsorptionBonus)
	// mage.AddStat(stats.FireResistance, magicAbsorptionBonus)
	// mage.AddStat(stats.FrostResistance, magicAbsorptionBonus)
	// mage.AddStat(stats.NatureResistance, magicAbsorptionBonus)
	// mage.AddStat(stats.ShadowResistance, magicAbsorptionBonus)
}
func (mage *Mage) applyHotStreak() {
	/* 	if mage.Talents.ImprovedHotStreak == 0 {
		return
	} */

	procChance := 1.0
	t10ProcAura := mage.BloodmagesRegalia2pcAura()

	mage.HotStreakAura = mage.RegisterAura(core.Aura{
		Label:    "HotStreak",
		ActionID: core.ActionID{SpellID: 48108},
		Duration: time.Second * 10,
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if t10ProcAura != nil {
				t10ProcAura.Activate(sim)
			}
		},
	})

	mage.hotStreakCritAura = mage.RegisterAura(core.Aura{
		Label:     "Hot Streak Proc Aura",
		ActionID:  core.ActionID{SpellID: 44448, Tag: 1},
		MaxStacks: 2,
		Duration:  time.Hour,
	})

	mage.RegisterAura(core.Aura{
		Label:    "Hot Streak Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(HotStreakSpells) {
				return
			}

			// Hot Streak Base Talent
			roll := sim.RandomFloat("Hot Streak")
			if result.DidCrit() && spell.Flags.Matches(HotStreakSpells) {
				if roll < procChance {
					if mage.HotStreakAura.IsActive() {
						mage.HotStreakAura.Refresh(sim)
					} else {
						mage.HotStreakAura.Activate(sim)
					}
				}
			}

			// Improved Hot Streak
			if !result.DidCrit() {
				mage.hotStreakCritAura.SetStacks(sim, 0)
				mage.hotStreakCritAura.Deactivate(sim)
				return
			}

			if mage.hotStreakCritAura.GetStacks() == 1 {
				if procChance == 1 || sim.Proc(procChance, "Hot Streak") {
					mage.hotStreakCritAura.SetStacks(sim, 0)
					mage.hotStreakCritAura.Deactivate(sim)

					mage.HotStreakAura.Activate(sim)
				}
			} else {
				mage.hotStreakCritAura.Activate(sim)
				mage.hotStreakCritAura.AddStack(sim)
			}
		},
	})

}

// func (mage *Mage) applyArcaneConcentration() {
// 	if mage.Talents.ArcaneConcentration == 0 {
// 		return
// 	}

// 	bonusCrit := float64(mage.Talents.ArcanePotency) * 15 * core.CritRatingPerCritChance

// 	// The result that caused the proc. Used to check we don't deactivate from the same proc.
// 	var proccedAt time.Duration
// 	var proccedSpell *core.Spell

// 	if mage.Talents.ArcanePotency > 0 {
// 		mage.ArcanePotencyAura = mage.RegisterAura(core.Aura{
// 			Label:    "Arcane Potency",
// 			ActionID: core.ActionID{SpellID: 31572},
// 			Duration: time.Second * 15,
// 			OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 				aura.Unit.AddStatDynamic(sim, stats.SpellCrit, bonusCrit)
// 			},
// 			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 				aura.Unit.AddStatDynamic(sim, stats.SpellCrit, -bonusCrit)
// 			},
// 			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
// 				if !spell.Flags.Matches(SpellFlagMage) {
// 					return
// 				}
// 				if proccedAt == sim.CurrentTime && proccedSpell == spell {
// 					// Means this is another hit from the same cast that procced CC.
// 					return
// 				}
// 				aura.Deactivate(sim)
// 			},
// 		})
// 	}

// 	mage.ClearcastingAura = mage.RegisterAura(core.Aura{
// 		Label:    "Clearcasting",
// 		ActionID: core.ActionID{SpellID: 12536},
// 		Duration: time.Second * 15,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Unit.PseudoStats.CostMultiplier -= 1
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Unit.PseudoStats.CostMultiplier += 1
// 		},
// 		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
// 			if !spell.Flags.Matches(SpellFlagMage) {
// 				return
// 			}
// 			if spell.DefaultCast.Cost == 0 {
// 				return
// 			}
// 			if spell == mage.ArcaneMissiles && mage.ArcaneMissilesAura.IsActive() {
// 				return
// 			}
// 			if proccedAt == sim.CurrentTime && proccedSpell == spell {
// 				// Means this is another hit from the same cast that procced CC.
// 				return
// 			}
// 			aura.Deactivate(sim)
// 		},
// 	})

// 	mage.RegisterAura(core.Aura{
// 		Label:    "Arcane Concentration",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if !spell.Flags.Matches(SpellFlagMage) || spell == mage.ArcaneMissiles {
// 				return
// 			}

// 			if !result.Landed() {
// 				return
// 			}

// 			procChance := 0.02 * float64(mage.Talents.ArcaneConcentration)

// 			// Arcane Missile ticks can proc CC, just at a low rate of about 1.5% with 5/5 Arcane Concentration
// 			if spell == mage.ArcaneMissilesTickSpell {
// 				procChance *= 0.15
// 			}

// 			if sim.RandomFloat("Arcane Concentration") > procChance {
// 				return
// 			}

// 			proccedAt = sim.CurrentTime
// 			proccedSpell = spell
// 			mage.ClearcastingAura.Activate(sim)
// 			if mage.ArcanePotencyAura != nil {
// 				mage.ArcanePotencyAura.Activate(sim)
// 			}
// 		},
// 	})
// }

func (mage *Mage) applyArcaneMissileProc() {

	t10ProcAura := mage.BloodmagesRegalia2pcAura()

	mage.ArcaneMissilesAura = mage.RegisterAura(core.Aura{
		Label:    "Missile Barrage Proc",
		ActionID: core.ActionID{SpellID: 44401},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mage.ArcaneMissiles.CostMultiplier -= 100
			mage.ArcaneMissiles.CastTimeMultiplier /= 2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.ArcaneMissiles.CostMultiplier += 100
			mage.ArcaneMissiles.CastTimeMultiplier *= 2
			if t10ProcAura != nil {
				t10ProcAura.Activate(sim)
			}
		},
	})

	var procChance float64
	if mage.Talents.HotStreak {
		procChance = .05 // Chance for hot streak crit
	} else if mage.Talents.BrainFreeze == 0 {
		procChance = .05 * float64(mage.Talents.BrainFreeze)
	} else {
		procChance = 0.4
	}

	mage.RegisterAura(core.Aura{
		Label:    "Missile Barrage Activation",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) { // Arcane Missiles and Brain Freeze proc on cast complete
			if !spell.Flags.Matches(BarrageSpells) {
				return
			}
			roll := sim.RandomFloat("Missile Barrage")
			if roll < procChance {
				mage.ArcaneMissilesAura.Activate(sim)
			}
		},
	})
}

func (mage *Mage) registerPresenceOfMindCD() {
	if !mage.Talents.PresenceOfMind {
		return
	}

	actionID := core.ActionID{SpellID: 12043}
	var spellToUse *core.Spell
	mage.Env.RegisterPostFinalizeEffect(func() {
		if mage.Pyroblast != nil {
			spellToUse = mage.Pyroblast
		} else if mage.PrimaryTalentTree == 1 {
			spellToUse = mage.Fireball
		} else if mage.PrimaryTalentTree == 2 {
			spellToUse = mage.Frostbolt
		} else {
			spellToUse = mage.ArcaneBlast
		}
	})

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * time.Duration(120*(1-mage.GetArcaneFlowsCDReduction())),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			if !mage.GCD.IsReady(sim) {
				return false
			}
			if mage.ArcanePowerAura.IsActive() {
				return false
			}

			manaCost := spellToUse.DefaultCast.Cost * mage.PseudoStats.CostMultiplier
			if spellToUse == mage.ArcaneBlast {
				manaCost *= float64(mage.ArcaneBlastAura.GetStacks()) * 1.75
			}
			return mage.CurrentMana() >= manaCost
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if mage.ArcanePotencyAura != nil {
				mage.ArcanePotencyAura.Activate(sim)
			}
			normalCastTime := spellToUse.DefaultCast.CastTime
			spellToUse.DefaultCast.CastTime = 0
			spellToUse.Cast(sim, mage.CurrentTarget)
			spellToUse.DefaultCast.CastTime = normalCastTime
		},
	})
	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) registerArcanePowerCD() {
	if !mage.Talents.ArcanePower {
		return
	}
	actionID := core.ActionID{SpellID: 12042}

	var affectedSpells []*core.Spell
	mage.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagMage) {
			affectedSpells = append(affectedSpells, spell)
		}
	})

	mage.ArcanePowerAura = mage.RegisterAura(core.Aura{
		Label:    "Arcane Power",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplierAdditive += 0.2
				spell.CostMultiplier += 0.1
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplierAdditive -= 0.2
				spell.CostMultiplier -= 0.2
			}
		},
	})
	core.RegisterPercentDamageModifierEffect(mage.ArcanePowerAura, 1.2)

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * time.Duration(120*(1-mage.GetArcaneFlowsCDReduction())),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.ArcanePowerAura.Activate(sim)
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !mage.ArcanePotencyAura.IsActive()
		},
	})
	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) GetArcaneFlowsCDReduction() float64 {
	switch float64(mage.Talents.ArcaneFlows) {
	case 2:
		return 0.25
	case 1:
		return 0.12
	case 0:
		return 0
	}
	return 0
}

func (mage *Mage) applyMasterOfElements() {
	if mage.Talents.MasterOfElements == 0 {
		return
	}

	refundCoeff := 0.15 * float64(mage.Talents.MasterOfElements)
	manaMetrics := mage.NewManaMetrics(core.ActionID{SpellID: 29077})

	mage.RegisterAura(core.Aura{
		Label:    "Master of Elements",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
				return
			}
			if spell.CurCast.Cost == 0 {
				return
			}
			if result.DidCrit() {
				if refundCoeff < 0 {
					mage.SpendMana(sim, -1*spell.DefaultCast.Cost*refundCoeff, manaMetrics)
				} else {
					mage.AddMana(sim, spell.DefaultCast.Cost*refundCoeff, manaMetrics)
				}
			}
		},
	})
}

/* func (mage *Mage) registerCombustionCD() {

	if !mage.Talents.Combustion {
		return
	}

	var combustionDotDamage float64
	var spellsToDuplicate []*core.Spell
	if mage.LivingBomb.Dot(target) {
		spellsToDuplicate = append(dotsToDuplicate, mage.LivingBomb)
	}
	if mage.Pyroblast.Get(target) {
		spellsToDuplicate = append(dotsToDuplicate, mage.Pyroblast)
	}
	if mage.Ignite.Get(target) {
		spellsToDuplicate = append(dotsToDuplicate, mage.Ignite)
	}
	for _, activeDots := range spellsToDuplicate {
		combustionDotDamage += float64(activeDots.CalcPeriodicDamage().Outcome) * float64(activeDots.Dot(activeDots.Unit).NumberOfTicks)
	}

	combustionAura := mage.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Combustion",
			ActionID: actionID,
			Duration: 15 * time.Second,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
			},
		})
	})
	mage.Combustion = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 11129},
		SpellSchool: core.SpellSchoolFire,
		//ProcMask:    core.SpellFlagNoOnCastComplete,
		Flags: SpellFlagMage,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: mage.GetFireMasteryBonusMultiplier(),
		DamageMultiplierAdditive: 1 +
			.01*float64(mage.Talents.FirePower),
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.429*mage.ScalingBaseDamage + 1.113*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
			spell.Dot(target).Apply(sim)
			spell.DealDamage(sim, result)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.Combustion,
		Type:  core.CooldownTypeDPS,
	})
}
*/

func (mage *Mage) registerIcyVeinsCD() {
	if !mage.Talents.IcyVeins {
		return
	}

	actionID := core.ActionID{SpellID: 12472}
	icyVeinsAura := mage.RegisterAura(core.Aura{
		Label:    "Icy Veins",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyCastSpeed(1.2)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyCastSpeed(1 / 1.2)
		},
	})

	mage.IcyVeins = mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		ManaCost: core.ManaCostOptions{
			BaseCost: 0.03,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * time.Duration(180*[]float64{1, .93, .86, .80}[mage.Talents.IceFloes]),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Need to check for icy veins already active in case Cold Snap is used right after.
			return !icyVeinsAura.IsActive()
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			icyVeinsAura.Activate(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.IcyVeins,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) registerColdSnapCD() {
	if !mage.Talents.ColdSnap {
		return
	}

	actionID := core.ActionID{SpellID: 11958}
	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 8,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Don't use if there are no cooldowns to reset.
			return (mage.IcyVeins != nil && !mage.IcyVeins.IsReady(sim)) ||
				(mage.SummonWaterElemental != nil && !mage.SummonWaterElemental.IsReady(sim))
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if mage.IcyVeins != nil {
				mage.IcyVeins.CD.Reset()
			}
			if mage.SummonWaterElemental != nil {
				mage.SummonWaterElemental.CD.Reset()
			}
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			// Ideally wait for both water ele and icy veins so we can reset both.
			if mage.IcyVeins != nil && mage.IcyVeins.IsReady(sim) {
				return false
			}
			if mage.SummonWaterElemental != nil && mage.SummonWaterElemental.IsReady(sim) {
				return false
			}
			return true
		},
	})
}

// func (mage *Mage) applyMoltenFury() {
// 	if mage.Talents.MoltenFury == 0 {
// 		return
// 	}

// 	multiplier := 1.0 + 0.06*float64(mage.Talents.MoltenFury)

// 	mage.RegisterResetEffect(func(sim *core.Simulation) {
// 		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
// 			if isExecute == 35 {
// 				mage.PseudoStats.DamageDealtMultiplier *= multiplier
// 				// For some reason Molten Fury doesn't apply to living bomb DoT, so cancel it out.
// 				if mage.LivingBomb != nil {
// 					mage.LivingBomb.DamageMultiplier /= multiplier
// 				}
// 			}
// 		})
// 	})
// }

func (mage *Mage) hasChillEffect(spell *core.Spell) bool {
	return spell == mage.Frostbolt || spell == mage.FrostfireBolt || (spell == mage.Blizzard && mage.Talents.IceShards > 0)
}

// func (mage *Mage) applyFingersOfFrost() {
// 	if mage.Talents.FingersOfFrost == 0 {
// 		return
// 	}

// 	bonusCrit := []float64{0, 17, 34, 50}[mage.Talents.Shatter] * core.CritRatingPerCritChance
// 	iceLanceMultiplier := core.TernaryFloat64(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfIceLance), 4, 3)

// 	var proccedAt time.Duration

// 	mage.FingersOfFrostAura = mage.RegisterAura(core.Aura{
// 		Label:     "Fingers of Frost Proc",
// 		ActionID:  core.ActionID{SpellID: 44545},
// 		Duration:  time.Second * 15,
// 		MaxStacks: 2,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			mage.AddStatDynamic(sim, stats.SpellCrit, bonusCrit)
// 			mage.IceLance.DamageMultiplier *= iceLanceMultiplier
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			mage.AddStatDynamic(sim, stats.SpellCrit, -bonusCrit)
// 			mage.IceLance.DamageMultiplier /= iceLanceMultiplier
// 		},
// 		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
// 			if proccedAt != sim.CurrentTime {
// 				aura.RemoveStack(sim)
// 			}
// 		},
// 	})

// 	procChance := []float64{0, .07, .15}[mage.Talents.FingersOfFrost]
// 	mage.RegisterAura(core.Aura{
// 		Label:    "Fingers of Frost Talent",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if mage.hasChillEffect(spell) && sim.RandomFloat("Fingers of Frost") < procChance {
// 				mage.FingersOfFrostAura.Activate(sim)
// 				mage.FingersOfFrostAura.SetStacks(sim, 2)
// 				proccedAt = sim.CurrentTime
// 			}
// 		},
// 	})
// }

func (mage *Mage) applyBrainFreeze() {
	if mage.Talents.BrainFreeze == 0 {
		return
	}

	hasT8_4pc := mage.HasSetBonus(ItemSetKirinTorGarb, 4)
	t10ProcAura := mage.BloodmagesRegalia2pcAura()

	mage.BrainFreezeAura = mage.GetOrRegisterAura(core.Aura{
		Label:    "Brain Freeze Proc",
		ActionID: core.ActionID{SpellID: 57761},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mage.Fireball.CostMultiplier -= 100
			mage.Fireball.CastTimeMultiplier -= 1
			mage.FrostfireBolt.CostMultiplier -= 100
			mage.FrostfireBolt.CastTimeMultiplier -= 1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.Fireball.CostMultiplier += 100
			mage.Fireball.CastTimeMultiplier += 1
			mage.FrostfireBolt.CostMultiplier += 100
			mage.FrostfireBolt.CastTimeMultiplier += 1
			if t10ProcAura != nil {
				t10ProcAura.Activate(sim)
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == mage.FrostfireBolt || spell == mage.Fireball {
				if !hasT8_4pc || sim.RandomFloat("MageT84PC") > T84PcProcChance {
					aura.Deactivate(sim)
				}
			}
		},
	})

	procChance := .05 * float64(mage.Talents.BrainFreeze)
	mage.RegisterAura(core.Aura{
		Label:    "Brain Freeze Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if mage.hasChillEffect(spell) && sim.RandomFloat("Brain Freeze") < procChance {
				mage.BrainFreezeAura.Activate(sim)
			}
		},
	})
}

// func (mage *Mage) applyWintersChill() {
// 	if mage.Talents.WintersChill == 0 {
// 		return
// 	}

// 	procChance := []float64{0, 0.33, 0.66, 1}[mage.Talents.WintersChill]

// 	wcAuras := mage.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
// 		return core.WintersChillAura(target, 0)
// 	})
// 	mage.Env.RegisterPreFinalizeEffect(func() {
// 		for _, spell := range mage.GetSpellsMatchingSchool(core.SpellSchoolFrost) {
// 			spell.RelatedAuras = append(spell.RelatedAuras, wcAuras)
// 		}
// 	})

// 	mage.RegisterAura(core.Aura{
// 		Label:    "Winters Chill Talent",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if !result.Landed() || !spell.SpellSchool.Matches(core.SpellSchoolFrost) {
// 				return
// 			}

// 			if sim.Proc(procChance, "Winters Chill") {
// 				aura := wcAuras.Get(result.Target)
// 				aura.Activate(sim)
// 				if aura.IsActive() {
// 					aura.AddStack(sim)
// 				}
// 			}
// 		},
// 	})
// }

// func (mage *Mage) applyFireStarter() {
// 	if mage.Talents.Firestarter == 0 {
// 		return
// 	}

// 	firestarterAura := mage.RegisterAura(core.Aura{
// 		Label:    "Firestarter",
// 		ActionID: core.ActionID{SpellID: 54741},
// 		Duration: 10 * time.Second,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			mage.Flamestrike.CostMultiplier -= 100
// 			mage.Flamestrike.CastTimeMultiplier -= 1
// 			mage.FlamestrikeRank8.CostMultiplier -= 100
// 			mage.FlamestrikeRank8.CastTimeMultiplier -= 1
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			mage.Flamestrike.CostMultiplier += 100
// 			mage.Flamestrike.CastTimeMultiplier += 1
// 			mage.FlamestrikeRank8.CostMultiplier += 100
// 			mage.FlamestrikeRank8.CastTimeMultiplier += 1
// 		},
// 		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
// 			if spell == mage.Flamestrike || spell == mage.FlamestrikeRank8 {
// 				aura.Deactivate(sim)
// 			}
// 		},
// 	})

// 	mage.RegisterAura(core.Aura{
// 		Label:    "Firestarter talent",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if !result.Landed() {
// 				return
// 			}

// 			if spell == mage.BlastWave || spell == mage.DragonsBreath {
// 				firestarterAura.Activate(sim)
// 			}
// 		},
// 	})
// }
