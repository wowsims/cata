package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func (shaman *Shaman) ApplyTalents() {
	// TODO: Confirm this is best way to apply acuity?
	shaman.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*1*float64(shaman.Talents.Acuity))
	shaman.AddStat(stats.SpellCrit, core.CritRatingPerCritChance*1*float64(shaman.Talents.Acuity))
	shaman.AddStat(stats.Expertise, 4*core.ExpertisePerQuarterPercentReduction*float64(shaman.Talents.UnleashedRage))

	if shaman.Spec == proto.Spec_SpecEnhancementShaman {
		shaman.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*6)

		//TODO: Get shaman mastery points
		masteryBonusPoints := 0.0
		shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1.2 + (masteryBonusPoints * 0.025)
		shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] *= 1.2 + (masteryBonusPoints * 0.025)
		shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] *= 1.2 + (masteryBonusPoints * 0.025)

		shaman.PseudoStats.CanParry = true
		//TODO: Enhancement bonuses for mental quickness and primal wisdom
	}

	if shaman.Talents.Toughness > 0 {
		shaman.MultiplyStat(stats.Stamina, []float64{1.0, 1.03, 1.07, 1.1}[shaman.Talents.Toughness])
	}

	shaman.applyElementalFocus()
	shaman.applyRollingThunder()
	shaman.applyFulmination()
	shaman.applyElementalDevastation()
	shaman.applyFlurry()
	shaman.applyMaelstromWeapon()
	shaman.registerElementalMasteryCD()
	shaman.registerNaturesSwiftnessCD()
	shaman.registerShamanisticRageCD()
	shaman.registerManaTideTotemCD()
}

func (shaman *Shaman) applyElementalFocus() {
	if !shaman.Talents.ElementalFocus {
		return
	}

	oathBonus := 1 + 0.05*float64(shaman.Talents.ElementalOath)
	var affectedSpells []*core.Spell

	// TODO: fix this.
	// Right now: Set to 3 so that the spell that cast it consumes a charge down to expected 2.
	// Correct fix would be to figure out how to make 'onCastComplete' fire before 'onspellhitdealt' without breaking all the other things.
	maxStacks := int32(3)

	// TODO: need to check for additional spells that benefit from the cost reduction
	// TODO: the damage bonus may apply to more spells now need to check how this works
	// I tested on beta and totems like magma benefit from the spell bonus
	clearcastingAura := shaman.RegisterAura(core.Aura{
		Label:     "Clearcasting",
		ActionID:  core.ActionID{SpellID: 16246},
		Duration:  time.Second * 15,
		MaxStacks: maxStacks,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice([]*core.Spell{
				shaman.LightningBolt,
				shaman.ChainLightning,
				shaman.LavaBurst,
				shaman.FireNova,
				shaman.EarthShock,
				shaman.FlameShock,
				shaman.FrostShock,
			}, func(spell *core.Spell) bool { return spell != nil })
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.CostMultiplier -= 0.4
			}
			if oathBonus > 1 {
				shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] *= oathBonus
				shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= oathBonus
				shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] *= oathBonus
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.CostMultiplier += 0.4
			}
			if oathBonus > 1 {
				shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] /= oathBonus
				shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= oathBonus
				shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] /= oathBonus
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Flags.Matches(SpellFlagShock | SpellFlagFocusable) {
				return
			}
			if spell.ActionID.Tag == 6 { // Filter LO casts
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
			// TODO: need to check for additional spells that cause the proc
			if !spell.Flags.Matches(SpellFlagShock | SpellFlagFocusable) {
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

	shaman.RegisterAura(core.Aura{
		Label:    "Rolling Thunder",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if (spell == shaman.LightningBolt || spell == shaman.ChainLightning) && shaman.SelfBuffs.Shield == proto.ShamanShield_LightningShield {
				if sim.RandomFloat("Rolling Thunder") < 0.3*float64(shaman.Talents.RollingThunder) {
					shaman.AddMana(sim, 0.02*shaman.MaxMana(), manaMetrics)
					shaman.LightningShieldAura.AddStack(sim)
				}
			}
		},
	})
}

func (shaman *Shaman) applyFulmination() {
	if !shaman.Talents.Fulmination {
		return
	}

	shaman.Fulmination = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 88767},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskProc,
		Flags:       SpellFlagElectric | SpellFlagFocusable,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: 0,
				GCD:      0,
			},
		},

		BonusHitRating:   float64(shaman.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance,
		BonusCritRating:  0,
		DamageMultiplier: 1 + 0.02*float64(shaman.Talents.Concussion) + 0.05*float64(shaman.Talents.ImprovedShields),
		CritMultiplier:   shaman.ElementalFuryCritMultiplier(0),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damagePerOrb := 350 + 0.267*spell.SpellPower()
			totalDamage := damagePerOrb * float64(shaman.LightningShieldAura.GetStacks())
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

	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfElementalMastery) {
		// TODO: apply damage reduction if necessary
	}

	// TODO: Share CD with Natures Swiftness

	buffAura := shaman.RegisterAura(core.Aura{
		Label:    "Elemental Mastery Buff",
		ActionID: core.ActionID{SpellID: 64701},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyCastSpeed(1.20)
			// This is updated for the new elemental mastery this says fire/frost/nature damage increased by 15%.
			// In beta this looks like it is applying to magma totem even if you use it after searing totem is dropped.
			// It is not doing the same for fire elemental totem
			// need to look into how this multiplier works and if it is affecting totems
			shaman.PseudoStats.DamageDealtMultiplier *= 1.15
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyCastSpeed(1 / 1.20)
			shaman.PseudoStats.DamageDealtMultiplier /= 1.18
		},
	})

	emAura := shaman.RegisterAura(core.Aura{
		Label:    "Elemental Mastery",
		ActionID: eleMasterActionID,
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
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell != shaman.LightningBolt && spell != shaman.ChainLightning && spell != shaman.LavaBurst {
				return
			}
			// Remove the buff and put skill on CD
			aura.Deactivate(sim)
			cdTimer.Set(sim.CurrentTime + cd)
			shaman.UpdateMajorCooldowns()
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
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
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

	// TODO: This also reduces mana cost in cata (20% per stack?) did I do this right?
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
	ppmm := shaman.AutoAttacks.NewPPMManager([]float64{3.0, 7.0, 10.0}[shaman.Talents.MaelstromWeapon], core.ProcMaskMelee)
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
