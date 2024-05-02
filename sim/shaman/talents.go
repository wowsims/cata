package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func (shaman *Shaman) ApplyTalents() {
	shaman.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*1*float64(shaman.Talents.Acuity))
	shaman.AddStat(stats.SpellCrit, core.CritRatingPerCritChance*1*float64(shaman.Talents.Acuity))
	shaman.AddStat(stats.Expertise, 4*core.ExpertisePerQuarterPercentReduction*float64(shaman.Talents.UnleashedRage))

	if shaman.Talents.Concussion > 0 {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask: SpellMaskLightningBolt | SpellMaskLightningBoltOverload | SpellMaskChainLightning | SpellMaskChainLightningOverload |
				SpellMaskThunderstorm | SpellMaskLavaBurst | SpellMaskLavaBurstOverload | SpellMaskEarthShock | SpellMaskFlameShock | SpellMaskFrost,
			Kind:       core.SpellMod_DamageDone_Flat,
			FloatValue: 0.02 * float64(shaman.Talents.Concussion),
		})
	}

	if shaman.Talents.Toughness > 0 {
		shaman.MultiplyStat(stats.Stamina, []float64{1.0, 1.03, 1.07, 1.1}[shaman.Talents.Toughness])
	}

	if shaman.Talents.ElementalPrecision > 0 {
		shaman.AddStats(stats.Stats{stats.SpellHit: []float64{0.0, -0.33, -0.66, -1.0}[shaman.Talents.ElementalPrecision] * shaman.GetBaseStats()[stats.Spirit]})
		shaman.AddStatDependency(stats.Spirit, stats.SpellHit, []float64{0.0, 0.33, 0.66, 1.0}[shaman.Talents.ElementalPrecision])
		shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1 + 0.01*float64(shaman.Talents.ElementalPrecision)
		shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] *= 1 + 0.01*float64(shaman.Talents.ElementalPrecision)
		shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] *= 1 + 0.01*float64(shaman.Talents.ElementalPrecision)
	}

	if shaman.Talents.CallOfFlame > 0 {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask:  SpellMaskLavaBurst,
			Kind:       core.SpellMod_DamageDone_Flat,
			FloatValue: 0.05 * float64(shaman.Talents.CallOfFlame),
		})

		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask:  SpellMaskSearingTotem | SpellMaskMagmaTotem | SpellMaskFireNova,
			Kind:       core.SpellMod_DamageDone_Flat,
			FloatValue: 0.10 * float64(shaman.Talents.CallOfFlame),
		})
	}

	shaman.applyElementalFocus()
	shaman.applyRollingThunder()

	if shaman.Talents.LavaFlows > 0 {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask:  SpellMaskFlameShockDot,
			Kind:       core.SpellMod_DamageDone_Flat,
			FloatValue: 0.20 * float64(shaman.Talents.LavaFlows),
		})
		/* TODO: Use if CritMultiplier_Flat gets added
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask:  SpellMaskLavaBurst | SpellMaskLavaBurstOverload,
			Kind:       core.SpellMod_CritMultiplier_Flat,
			FloatValue: 0.08 * float64(shaman.Talents.LavaFlows),
		})
		*/
	}

	shaman.applyLavaSurge()

	shaman.applyFulmination()

	if shaman.Talents.Earthquake {
		shaman.registerEarthquakeSpell()
	}

	if shaman.Talents.FocusedStrikes > 0 {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask:  SpellMaskPrimalStrike,
			Kind:       core.SpellMod_DamageDone_Pct,
			FloatValue: 0.15 * float64(shaman.Talents.FocusedStrikes),
		})

		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask:  SpellMaskStormstrike,
			Kind:       core.SpellMod_DamageDone_Pct,
			FloatValue: 0.15 * float64(shaman.Talents.FocusedStrikes),
		})
	}

	if shaman.Talents.ImprovedShields > 0 {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask:  SpellMaskLightningShield | SpellMaskFulmination | SpellMaskEarthShield,
			Kind:       core.SpellMod_DamageDone_Flat,
			FloatValue: 0.05 * float64(shaman.Talents.ImprovedShields),
		})
	}

	shaman.applyElementalDevastation()

	if shaman.Talents.Stormstrike {
		shaman.registerStormstrikeSpell()
	}

	shaman.applyFlurry()
	shaman.applyMaelstromWeapon()
	shaman.applySearingFlames()
	shaman.applyTotemicFocus()

	if shaman.Talents.FeralSpirit {
		shaman.registerFeralSpirit()
	}

	shaman.registerElementalMasteryCD()
	shaman.registerNaturesSwiftnessCD()
	shaman.registerShamanisticRageCD()
	shaman.registerManaTideTotemCD()

	shaman.ApplyGlyphs()
}

func (shaman *Shaman) applyElementalFocus() {
	if !shaman.Talents.ElementalFocus {
		return
	}

	affectedSpells := SpellMaskLightningBolt | SpellMaskChainLightning | SpellMaskLavaBurst | SpellMaskFireNova | SpellMaskEarthShock | SpellMaskFlameShock | SpellMaskFrostShock
	costReductionMod := shaman.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		ClassMask:  affectedSpells,
		FloatValue: -0.4,
	})

	oathBonus := 0.05 * float64(shaman.Talents.ElementalOath)
	oathMod := shaman.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		School:     core.SpellSchoolFire | core.SpellSchoolFrost | core.SpellSchoolNature,
		FloatValue: oathBonus,
	})
	oathModEarthquake := shaman.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  SpellMaskEarthquake,
		FloatValue: oathBonus,
	})

	maxStacks := int32(2)

	// TODO: need to check for additional spells that benefit from the cost reduction
	clearcastingAura := shaman.RegisterAura(core.Aura{
		Label:     "Clearcasting",
		ActionID:  core.ActionID{SpellID: 16246},
		Duration:  time.Second * 15,
		MaxStacks: maxStacks,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			costReductionMod.Activate()
			oathMod.Activate()
			oathModEarthquake.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			costReductionMod.Deactivate()
			oathMod.Deactivate()
			oathModEarthquake.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Flags.Matches(SpellFlagShock|SpellFlagFocusable) || (spell.ClassSpellMask&(SpellMaskOverload|SpellMaskThunderstorm) != 0) {
				return
			}
			aura.RemoveStack(sim)
		},
	})

	shaman.RegisterAura(core.Aura{
		Label:    "Elemental Focus",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(SpellFlagShock|SpellFlagFocusable) || (spell.ClassSpellMask&(SpellMaskOverload|SpellMaskUnleashFlame) != 0) || spell == shaman.Earthquake {
				return
			}
			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			clearcastingAura.Activate(sim)
			clearcastingAura.SetStacks(sim, maxStacks)
		},
	})
}

func (shaman *Shaman) applyRollingThunder() {
	if shaman.Talents.RollingThunder == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 88765}
	manaMetrics := shaman.NewManaMetrics(actionID)

	// allowedSpells := make([]*core.Spell, 0)
	// allowedSpells = append(allowedSpells, shaman.LightningBolt, shaman.LightningBoltOverload, shaman.ChainLightning)
	// allowedSpells = append(allowedSpells, shaman.ChainLightningOverloads...)

	shaman.RegisterAura(core.Aura{
		Label:    "Rolling Thunder",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if (spell == shaman.LightningBolt || spell == shaman.ChainLightning || spell == shaman.LightningBoltOverload) && shaman.SelfBuffs.Shield == proto.ShamanShield_LightningShield {
				// for _, allowedSpell := range allowedSpells {
				// 	if spell == allowedSpell {
				if sim.RandomFloat("Rolling Thunder") < 0.3*float64(shaman.Talents.RollingThunder) {
					shaman.AddMana(sim, 0.02*shaman.MaxMana(), manaMetrics)
					shaman.LightningShieldAura.Refresh(sim)
					shaman.LightningShieldAura.AddStack(sim)
				}
				//  }
			}
		},
	})
}

func (shaman *Shaman) applyLavaSurge() {
	if shaman.Talents.LavaSurge == 0 {
		return
	}

	has4PT12 := shaman.HasSetBonus(ItemSetVolcanicRegalia, 4)

	var instantLavaSurgeMod *core.SpellMod

	if has4PT12 {
		instantLavaSurgeMod = shaman.AddDynamicMod(core.SpellModConfig{
			Kind:       core.SpellMod_CastTime_Pct,
			FloatValue: -1,
			ClassMask:  SpellMaskLavaBurst,
		})
	}

	shaman.RegisterAura(core.Aura{
		Label:    "Lava Surge",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ClassSpellMask != SpellMaskFlameShockDot || !sim.Proc(0.1*float64(shaman.Talents.LavaSurge), "LavaSurge") {
				return
			}
			shaman.LavaBurst.CD.Reset()
			if has4PT12 {
				instantLavaSurgeMod.Activate()
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask != SpellMaskLavaBurst || !has4PT12 {
				return
			}
			instantLavaSurgeMod.Deactivate()
		},
	})
}

func (shaman *Shaman) applyFulmination() {
	if !shaman.Talents.Fulmination {
		return
	}

	shaman.Fulmination = shaman.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 88767},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskProc,
		Flags:          SpellFlagFocusable,
		ClassSpellMask: SpellMaskFulmination,
		ManaCost: core.ManaCostOptions{
			BaseCost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: 0,
				GCD:      0,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   shaman.DefaultSpellCritMultiplier(),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			totalDamage := (shaman.ClassSpellScaling*0.38899999857 + 0.267*spell.SpellPower()) * (float64(shaman.LightningShieldAura.GetStacks()) - 3)
			result := spell.CalcDamage(sim, target, totalDamage, spell.OutcomeMagicHitAndCrit)
			spell.DealDamage(sim, result)
		},
	})
}

func (shaman *Shaman) applyElementalDevastation() {
	if shaman.Talents.ElementalDevastation == 0 {
		return
	}

	critBonus := 3.0 * float64(shaman.Talents.ElementalDevastation) * core.CritRatingPerCritChance
	procAura := shaman.NewTemporaryStatsAura("Elemental Devastation Proc", core.ActionID{SpellID: 30160}, stats.Stats{stats.MeleeCrit: critBonus}, time.Second*10)

	shaman.RegisterAura(core.Aura{
		Label:    "Elemental Devastation",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
				return
			}
			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			procAura.Activate(sim)
		},
	})
}

var eleMasterActionID = core.ActionID{SpellID: 16166}

func (shaman *Shaman) registerElementalMasteryCD() {
	if !shaman.Talents.ElementalMastery {
		return
	}

	cdTimer := shaman.NewTimer()
	cd := time.Minute * 3

	damageMod := shaman.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		School:     core.SpellSchoolFire | core.SpellSchoolFrost | core.SpellSchoolNature,
		FloatValue: 0.15,
	})

	has2PT13 := shaman.HasSetBonus(ItemSetSpiritwalkersRegalia, 2)

	buffAura := shaman.RegisterAura(core.Aura{
		Label:    "Elemental Mastery Buff",
		ActionID: core.ActionID{SpellID: 64701},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyCastSpeed(1.20)
			damageMod.Activate()
			if has2PT13 {
				shaman.AddStatDynamic(sim, stats.Mastery, 2000)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyCastSpeed(1 / 1.20)
			damageMod.Deactivate()
			if has2PT13 {
				shaman.AddStatDynamic(sim, stats.Mastery, -2000)
			}
		},
	})

	affectedSpells := SpellMaskChainLightning | SpellMaskLavaBurst | SpellMaskLightningBolt
	castTimeMod := shaman.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		ClassMask:  affectedSpells,
		FloatValue: -1,
	})

	emAura := shaman.RegisterAura(core.Aura{
		Label:    "Elemental Mastery",
		ActionID: eleMasterActionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Deactivate()
		},
		//TODO: there are timeline graphical anomalies due to elemental mastery's logic.
		//It doesn't know which lighting bolt cast correspond to which hit so put a grey bar under one and two hit to the second one
		//Travel time visual is also missing
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask&affectedSpells > 0 {
				// Remove the buff and put skill on CD
				aura.Deactivate(sim)
				cdTimer.Set(sim.CurrentTime + cd)
				shaman.UpdateMajorCooldowns()
			}
		},
	})

	eleMastSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID: eleMasterActionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			buffAura.Activate(sim)
			emAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: eleMastSpell,
		Type:  core.CooldownTypeDPS,
	})

	if shaman.Talents.Feedback > 0 {
		shaman.RegisterAura(core.Aura{
			Label:    "Feedback",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if (spell == shaman.LightningBolt || spell == shaman.ChainLightning) && !eleMastSpell.CD.IsReady(sim) {
					*eleMastSpell.CD.Timer = core.Timer(time.Duration(*eleMastSpell.CD.Timer) - time.Second*time.Duration(shaman.Talents.Feedback))
					shaman.UpdateMajorCooldowns() // this could get expensive because it will be called all the time.
				}
			},
		})
	}
}

func (shaman *Shaman) registerNaturesSwiftnessCD() {
	if !shaman.Talents.NaturesSwiftness {
		return
	}
	actionID := core.ActionID{SpellID: 16188}
	cdTimer := shaman.NewTimer()
	cd := time.Minute * 2

	nsAura := shaman.RegisterAura(core.Aura{
		Label:    "Natures Swiftness",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.ChainLightning.CastTimeMultiplier -= 1
			shaman.LavaBurst.CastTimeMultiplier -= 1
			shaman.LightningBolt.CastTimeMultiplier -= 1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.ChainLightning.CastTimeMultiplier += 1
			shaman.LavaBurst.CastTimeMultiplier += 1
			shaman.LightningBolt.CastTimeMultiplier += 1
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell != shaman.LightningBolt && spell != shaman.ChainLightning && spell != shaman.LavaBurst {
				return
			}

			// Remove the buff and put skill on CD
			aura.Deactivate(sim)
			cdTimer.Set(sim.CurrentTime + cd)
			shaman.UpdateMajorCooldowns()
		},
	})

	nsSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Don't use NS unless we're casting a full-length lightning bolt, which is
			// the only spell shamans have with a cast longer than GCD.
			return !shaman.HasTemporarySpellCastSpeedIncrease()
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			nsAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: nsSpell,
		Type:  core.CooldownTypeDPS,
	})
}

// TODO: Updated talent and id for cata but this might be working differently in cata, ie is there an ICD anymore? wowhead shows no icd vs wrath 500ms
func (shaman *Shaman) applyFlurry() {
	if shaman.Talents.Flurry == 0 {
		return
	}

	bonus := 1.0 + 0.10*float64(shaman.Talents.Flurry)

	inverseBonus := 1 / bonus

	procAura := shaman.RegisterAura(core.Aura{
		Label:     "Flurry Proc",
		ActionID:  core.ActionID{SpellID: 16282},
		Duration:  core.NeverExpires,
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyMeleeSpeed(sim, bonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyMeleeSpeed(sim, inverseBonus)
		},
	})

	icd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Millisecond * 500,
	}

	shaman.RegisterAura(core.Aura{
		Label:    "Flurry",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			if result.Outcome.Matches(core.OutcomeCrit) {
				procAura.Activate(sim)
				procAura.SetStacks(sim, 3)
				icd.Reset() // the "charge protection" ICD isn't up yet
				return
			}

			// Remove a stack.
			if procAura.IsActive() && spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) && icd.IsReady(sim) {
				icd.Use(sim)
				procAura.RemoveStack(sim)
			}
		},
	})
}

func (shaman *Shaman) applyMaelstromWeapon() {
	if shaman.Talents.MaelstromWeapon == 0 {
		return
	}

	// TODO: Don't forget to make it so that AA don't reset when casting when MW is active
	// for LB / CL / LvB
	// They can't actually hit while casting, but the AA timer doesnt reset if you cast during the AA timer.

	// For sim purposes maelstrom weapon only impacts CL / LB
	shaman.MaelstromWeaponAura = shaman.RegisterAura(core.Aura{
		Label:     "MaelstromWeapon Proc",
		ActionID:  core.ActionID{SpellID: 51530},
		Duration:  time.Second * 30,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			multDiff := 0.2 * float64(newStacks-oldStacks)
			shaman.LightningBolt.CastTimeMultiplier -= multDiff
			shaman.LightningBolt.CostMultiplier -= multDiff
			shaman.ChainLightning.CastTimeMultiplier -= multDiff
			shaman.ChainLightning.CostMultiplier -= multDiff
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(SpellFlagElectric) {
				return
			}
			shaman.MaelstromWeaponAura.Deactivate(sim)
		},
	})

	// TODO: This was 2% per talent point and max of 10% proc in wotlk. Can't find data on proc chance in cata but the talent was reduced to 3 pts. Guessing it is 3/7/10 like other talents
	ppmm := shaman.AutoAttacks.NewPPMManager([]float64{0.0, 3.0, 6.0, 10.0}[shaman.Talents.MaelstromWeapon], core.ProcMaskMelee)
	// This aura is hidden, just applies stacks of the proc aura.
	shaman.RegisterAura(core.Aura{
		Label:    "MaelstromWeapon",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if ppmm.Proc(sim, spell.ProcMask, "Maelstrom Weapon") {
				shaman.MaelstromWeaponAura.Activate(sim)
				shaman.MaelstromWeaponAura.AddStack(sim)
			}
		},
	})
}

func (shaman *Shaman) applySearingFlames() {
	if shaman.Talents.SearingFlames == 0 {
		return
	}

	shaman.SearingFlames = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 77657},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers | core.SpellFlagNoOnDamageDealt,

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Searing Flames",
				MaxStacks: 5,
			},
			TickLength:    time.Second * 3,
			NumberOfTicks: 5,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
		},
	})

	core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
		Name:           "Searing Flames",
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: SpellMaskSearingTotem,
		Outcome:        core.OutcomeLanded,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			dot := shaman.SearingFlames.Dot(result.Target)

			if shaman.Talents.SearingFlames == 3 || sim.RandomFloat("Searing Flames") < 0.33*float64(shaman.Talents.SearingFlames) {
				dot.Aura.Activate(sim)
				dot.Aura.AddStack(sim)

				// recalc damage based on stacks, testing with searing totem seems to indicate the damage is updated dynamically on refesh
				// instantly taking the bonus of any procs or buffs and applying it times the number of stacks
				dot.SnapshotBaseDamage = float64(dot.GetStacks()) * result.Damage / float64(dot.NumberOfTicks)
				dot.SnapshotAttackerMultiplier = shaman.SearingFlames.DamageMultiplier
				shaman.SearingFlames.Cast(sim, result.Target)
			}
		},
	})

}

func (shaman *Shaman) applyTotemicFocus() {
	if shaman.Talents.TotemicFocus == 0 {
		return
	}

	shaman.FireElemental.maxFireBlastCasts = int32(float64(shaman.FireElemental.maxFireBlastCasts) * (1.0 + 0.20*float64(shaman.Talents.TotemicFocus)))
	shaman.FireElemental.maxFireNovaCasts = int32(float64(shaman.FireElemental.maxFireNovaCasts) * (1.0 + 0.20*float64(shaman.Talents.TotemicFocus)))
}

func (shaman *Shaman) registerManaTideTotemCD() {
	if !shaman.Talents.ManaTideTotem {
		return
	}

	mttAura := core.ManaTideTotemAura(shaman.GetCharacter(), shaman.Index)
	mttSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID: core.ManaTideTotemActionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Minute * 3,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mttAura.Activate(sim)

			// If healing stream is active, cancel it while mana tide is up.
			if shaman.HealingStreamTotem.Hot(&shaman.Unit).IsActive() {
				for _, agent := range shaman.Party.Players {
					shaman.HealingStreamTotem.Hot(&agent.GetCharacter().Unit).Cancel(sim)
				}
			}

			// TODO: Current water totem buff needs to be removed from party/raid.
			if shaman.Totems.Water != proto.WaterTotem_NoWaterTotem {
				shaman.TotemExpirations[WaterTotem] = sim.CurrentTime + time.Second*12
			}
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: mttSpell,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return sim.CurrentTime > time.Second*30
		},
	})
}
