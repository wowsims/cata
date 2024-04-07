package priest

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func (priest *Priest) ApplyTalents() {
	// TODO:
	// Reflective Shield
	// Improved Flash Heal
	// Renewed Hope
	// Rapture
	// Pain Suppression
	// Test of Faith
	// Guardian Spirit

	// priest.applyDivineAegis()
	// priest.applyGrace()
	// priest.applyBorrowedTime()
	// priest.applyInspiration()
	// priest.applyHolyConcentration()
	// priest.applySerendipity()
	// priest.applySurgeOfLight()
	// priest.registerInnerFocus()

	// priest.AddStat(stats.SpellCrit, 1*float64(priest.Talents.FocusedWill)*core.CritRatingPerCritChance)
	// priest.PseudoStats.SpiritRegenRateCasting = []float64{0.0, 0.17, 0.33, 0.5}[priest.Talents.Meditation]
	// priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= 1 - .02*float64(priest.Talents.SpellWarding)
	// priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] *= 1 - .02*float64(priest.Talents.SpellWarding)
	// priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= 1 - .02*float64(priest.Talents.SpellWarding)
	// priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= 1 - .02*float64(priest.Talents.SpellWarding)
	// priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= 1 - .02*float64(priest.Talents.SpellWarding)
	// priest.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= 1 - .02*float64(priest.Talents.SpellWarding)

	// if priest.Talents.SpiritualGuidance > 0 {
	// 	priest.AddStatDependency(stats.Spirit, stats.SpellPower, 0.05*float64(priest.Talents.SpiritualGuidance))
	// }

	// if priest.Talents.MentalStrength > 0 {
	// 	priest.MultiplyStat(stats.Intellect, 1.0+0.03*float64(priest.Talents.MentalStrength))
	// }

	// if priest.Talents.ImprovedPowerWordFortitude > 0 {
	// 	priest.MultiplyStat(stats.Stamina, 1.0+.02*float64(priest.Talents.ImprovedPowerWordFortitude))
	// }

	// if priest.Talents.Enlightenment > 0 {
	// 	priest.MultiplyStat(stats.Spirit, 1+.02*float64(priest.Talents.Enlightenment))
	// 	priest.PseudoStats.CastSpeedMultiplier *= 1 + .02*float64(priest.Talents.Enlightenment)
	// }

	// if priest.Talents.FocusedPower > 0 {
	// 	priest.PseudoStats.DamageDealtMultiplier *= 1 + .02*float64(priest.Talents.FocusedPower)
	// }

	// if priest.Talents.SpiritOfRedemption {
	// 	priest.MultiplyStat(stats.Spirit, 1.05)
	// }

	// Disciplin Talents
	// Improved Power Word: Shield - TBD
	// Twin Disciplines
	if priest.Talents.TwinDisciplines > 0 {
		priest.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1 + (0.02 * float64(priest.Talents.TwinDisciplines))
		priest.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexHoly] *= 1 + (0.02 * float64(priest.Talents.TwinDisciplines))
	}

	// Mental Agillity
	if priest.Talents.MentalAgility > 0 {
		AddOrReplaceMod(&priest.PowerCostPercentMods, &PriestAuraMod[float64]{
			ClassSpell: PriestSpellInstant,
			BaseValue:  -0.04 * float64(priest.Talents.MentalAgility),
			Stacks:     1,
			SpellID:    14781,
		})
	}

	// Evangelism
	priest.applyEvangelism()

	// Archangel
	priest.applyArchangel()

	// Shadow Talents
	// Darkness
	if priest.Talents.Darkness > 0 {
		priest.PseudoStats.CastSpeedMultiplier *= 1 + (0.01 * float64(priest.Talents.Darkness))
	}

	// Improved Shadow Word: Pain
	if priest.Talents.ImprovedShadowWordPain > 0 {
		AddOrReplaceMod(&priest.DamageDonePercentAddMods, &PriestAuraMod[float64]{
			ClassSpell: PriestSpellShadowWordPain,
			BaseValue:  0.03 * float64(priest.Talents.ImprovedShadowWordPain),
			Stacks:     1,
			SpellID:    15317,
		})
	}

	// Veiled Shadows
	if priest.Talents.VeiledShadows > 0 {
		AddOrReplaceMod(&priest.CooldownMods, &PriestAuraMod[time.Duration]{
			ClassSpell: PriestSpellFade,
			BaseValue:  time.Duration(3) * time.Second * time.Duration(priest.Talents.VeiledShadows),
			Stacks:     1,
			SpellID:    15311,
		})

		AddOrReplaceMod(&priest.CooldownMods, &PriestAuraMod[time.Duration]{
			ClassSpell: PriestSpellShadowFiend,
			BaseValue:  time.Duration(30) * time.Second * time.Duration(priest.Talents.VeiledShadows),
			Stacks:     1,
			SpellID:    15311,
		})
	}

	// Improved Psychic Scream
	if priest.Talents.ImprovedPsychicScream > 0 {
		AddOrReplaceMod(&priest.CooldownMods, &PriestAuraMod[time.Duration]{
			ClassSpell: PriestSpellPsychicScream,
			BaseValue:  time.Duration(priest.Talents.ImprovedPsychicScream*2) * time.Second,
			Stacks:     1,
			SpellID:    15448,
		})
	}

	// Improved Mind Blast
	priest.applyImprovedMindBlast()

	// Improved Devouring Plague
	priest.applyImprovedDevouringPlague()

	// Twisted Faith
	if priest.Talents.TwistedFaith > 0 {
		priest.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1 + (0.01 * float64(priest.Talents.TwistedFaith))
		priest.AddStatDependency(stats.Spirit, stats.SpellHit, 0.5*float64(priest.Talents.TwistedFaith))
	}

	// Shadowform
	if priest.Talents.Shadowform {
		priest.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.15
	}

	// Phantasm: Not implemented
	// Harnessed Shadows
	priest.applyHarnessedShadows()

	// Silence: Not implemented
	// Vampiric Embrace: vampiric_embrace.go
	// Masochism
	priest.applyMasochism()

	// Mind Melt - shadow_word_pain.go <25% part
	priest.applyMindMelt()

	// pain and suffering
	priest.applyPainAndSuffering()

	// vampiric touch - vampiric_touch.go
	// Paralysis - Not implemented
	// Psychic Horror - Not implemented
	// Sin and Punishment
	priest.applySinAndPunishment()

	// Shadowy Apparition
	priest.applyShadowyApparition()

	// Dispersion - TBD

	// Glphys
	if priest.HasGlyph(int32(proto.PriestPrimeGlyph_GlyphOfShadowWordPain)) {
		AddOrReplaceMod(&priest.DamageDonePercentAddMods, &PriestAuraMod[float64]{
			SpellID:    55681,
			BaseValue:  0.1,
			ClassSpell: PriestSpellShadowWordPain,
		})
	}

	if priest.HasGlyph(int32(proto.PriestPrimeGlyph_GlyphOfMindFlay)) {
		AddOrReplaceMod(&priest.DamageDonePercentAddMods, &PriestAuraMod[float64]{
			SpellID:    55687,
			BaseValue:  0.1,
			ClassSpell: PriestSpellMindFlay,
		})
	}

	if priest.HasGlyph(int32(proto.PriestPrimeGlyph_GlyphOfDispersion)) {
		AddOrReplaceMod(&priest.CooldownMods, &PriestAuraMod[time.Duration]{
			SpellID:    63229,
			BaseValue:  time.Second * -45,
			ClassSpell: PriestSpellDispersion,
		})
	}

	if priest.HasGlyph(int32(proto.PriestPrimeGlyph_GlyphOfShadowWordDeath)) {
		priest.RegisterAura(core.Aura{
			Label:    "Glyph of Shadow Word: Death",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},

			Icd: &core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 6,
			},

			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if priest.ShadowWordDeath.IsEqual(spell) && sim.IsExecutePhase25() && aura.Icd.IsReady(sim) {
					aura.Icd.Use(sim)
					priest.ShadowWordDeath.CD.Reset()
				}
			},
		})
	}
}

// disciplin talents
func (priest *Priest) applyEvangelism() {
	if priest.Talents.Evangelism == 0 {
		return
	}

	priest.DarkEvangelismProcAura = priest.GetOrRegisterAura(core.Aura{
		Label:     "Dark EvangelismProc",
		ActionID:  core.ActionID{SpellID: 87118},
		Duration:  time.Second * 20,
		MaxStacks: 5,

		// dummy aura used to track stacks in spells
	})

	priest.HolyEvangelismProcAura = priest.GetOrRegisterAura(core.Aura{
		Label:     "EvangelismProc",
		ActionID:  core.ActionID{SpellID: 81661},
		Duration:  time.Second * 20,
		MaxStacks: 5,

		// dummy aura used to track stacks in spells
	})

	priest.GetOrRegisterAura(core.Aura{
		Label:    "Evangilism",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			switch spell.SpellID {
			case 585: // smite
			case 14914: // holy fire
				priest.AddHolyEvanglismStack(sim)
			case 15407: // mind flay
				priest.AddDarkEvangelismStack(sim)
			}
		},
	})
}

func (priest *Priest) applyArchangel() {

	archAngelMana := priest.NewManaMetrics(core.ActionID{SpellID: 87152})
	darkArchAngelMana := priest.NewManaMetrics(core.ActionID{SpellID: 87153})

	archAngelAura := priest.Unit.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 81700},
		Label:     "Archangel Aura",
		MaxStacks: 5,
		Duration:  time.Second * 18,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			if newStacks > oldStacks {
				priest.AddMana(sim, 0.01*priest.MaxMana()*float64((newStacks-oldStacks)), archAngelMana)
			}

			// place holder mod healing done not present yet
		},
	})

	darkArchAngelAura := priest.Unit.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 87153},
		Label:     "Dark Archangel Aura",
		MaxStacks: 5,
		Duration:  time.Second * 18,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			AddOrReplaceMod(&priest.DamageDonePercentAddMods, &PriestAuraMod[float64]{
				SpellID:    87153,
				BaseValue:  0.04,
				Stacks:     newStacks,
				ClassSpell: PriestSpellMindFlay | PriestSpellMindSpike | PriestSpellMindBlast | PriestSpellShadowWordDeath,
			})

			if newStacks > oldStacks {
				priest.AddMana(sim, 0.05*priest.MaxMana()*float64((newStacks-oldStacks)), darkArchAngelMana)
			}
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			RemoveMod(&priest.DamageDonePercentAddMods, 87153)
		},
	})

	priest.Archangel = priest.RegisterSpell(PriestSpellArchangel, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 87151},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,
		ManaCost: core.ManaCostOptions{
			BaseCost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 30,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if priest.HolyEvangelismProcAura.IsActive() {
				archAngelAura.Activate(sim)
				for i := 0; i < int(priest.HolyEvangelismProcAura.GetStacks()); i++ {
					archAngelAura.AddStack(sim)
				}

				priest.HolyEvangelismProcAura.Deactivate(sim)
			}
		},
	})
	priest.DarkArchangel = priest.RegisterSpell(PriestSpellDarkArchangel, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 87153},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,
		ManaCost: core.ManaCostOptions{
			BaseCost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 90,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if priest.DarkEvangelismProcAura.IsActive() {
				darkArchAngelAura.Activate(sim)
				for i := 0; i < int(priest.DarkEvangelismProcAura.GetStacks()); i++ {
					darkArchAngelAura.AddStack(sim)
				}

				priest.DarkEvangelismProcAura.Deactivate(sim)
			}
		},
	})
}

func (priest *Priest) applyImprovedMindBlast() {
	if priest.Talents.ImprovedMindBlast == 0 {
		return
	}

	AddOrReplaceMod(&priest.CooldownMods, &PriestAuraMod[time.Duration]{
		SpellID:    48301,
		ClassSpell: PriestSpellMindBlast,
		BaseValue:  time.Duration(priest.Talents.ImprovedMindBlast) * time.Millisecond * 500,
	})

	mindTraumaSpell := priest.RegisterSpell(PriestSpellMindTrauma, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48301},
		ProcMask:    core.ProcMaskProc,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagNoMetrics,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			MindTraumaAura(target).Activate(sim)
		},
	})

	priest.RegisterAura(core.Aura{
		Label:    "Improved Mind Blast",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && priest.MindBlast.IsEqual(spell) {
				if sim.RandomFloat("Improved Mind Blast") < 0.33*float64(priest.Talents.ImprovedMindBlast) {
					mindTraumaSpell.Cast(sim, result.Target)
				}
			}
		},
	})
}

func MindTraumaAura(target *core.Unit) *core.Aura {
	return target.GetOrRegisterAura(core.Aura{
		Label:    "Mind Trauma",
		ActionID: core.ActionID{SpellID: 48301},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.HealingTakenMultiplier *= 0.9
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.HealingTakenMultiplier /= 0.9
		},
	})
}

func (priest *Priest) applyImprovedDevouringPlague() {
	if priest.Talents.ImprovedDevouringPlague == 0 {
		return
	}

	// simple spell here as it does not use any dmg mods or calculations
	impDPDamage := priest.Unit.RegisterSpell(core.SpellConfig{
		ActionID:                 core.ActionID{SpellID: 63675},
		SpellSchool:              core.SpellSchoolShadow,
		ProcMask:                 core.ProcMaskProc,
		Flags:                    core.SpellFlagIgnoreAttackerModifiers,
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,
		CritMultiplier:           priest.SpellCritMultiplier(1, priest.ShadowCritMultiplier),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dp := priest.DevouringPlague
			dot := dp.Dot(target)

			// No need to use snapshot dmg. It won't be initialized in time and we're in the same sim cycle
			// so all mods and buffs will be the same
			dmg := float64(dot.NumberOfTicks*int32(dp.ExpectedTickDamage(sim, target))*priest.Talents.ImprovedDevouringPlague) * 0.15
			spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMagicCrit)
		},
	})

	priest.RegisterAura(core.Aura{
		Label:    "Improved Devouring Plague Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !priest.DevouringPlague.IsEqual(spell) || !result.Landed() {
				return
			}

			impDPDamage.Cast(sim, result.Target)
		},
	})
}

func (priest *Priest) applyHarnessedShadows() {
	if priest.Talents.HarnessedShadows == 0 {
		return
	}

	AddOrReplaceMod(&priest.ProcChanceMods, &PriestAuraMod[float64]{
		SpellID:    78228,
		BaseValue:  0.04 * float64(priest.Talents.HarnessedShadows),
		ClassSpell: PriestSpellShadowOrbPassive,
	})

	// Proc on Dmg Taken not implemented yet
}

func (priest *Priest) applyPainAndSuffering() {
	if priest.Talents.PainAndSuffering == 0 {
		return
	}

	handler := func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if !result.Landed() || spell.SpellID != priest.MindFlayAPL.SpellID {
			return
		}

		procChance := float64(priest.Talents.PainAndSuffering) * 0.3
		if sim.RandomFloat("Pain and Suffering") < procChance {
			swp := priest.ShadowWordPain.Dot(result.Target)
			if swp.IsActive() {
				swp.Rollover(sim)
			}
		}
	}

	priest.RegisterAura(core.Aura{
		Label:    "Pain and Suffering",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: handler,
		OnSpellHitDealt:       handler,
	})
}

func (priest *Priest) applyMasochism() {
	if priest.Talents.Masochism == 0 {
		return
	}

	// Should we care for the different ranks here?
	manaMetrics := priest.NewManaMetrics(core.ActionID{
		SpellID: 88995,
	})

	damageTakenHandler := func(sim *core.Simulation, damage float64) {
		if damage < priest.MaxHealth()*0.1 {
			return
		}

		priest.AddMana(sim, priest.MaxMana()*0.05*float64(priest.Talents.Masochism), manaMetrics)
	}

	priest.RegisterAura(
		core.Aura{
			Label:    "Masochism",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitTaken: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				damageTakenHandler(sim, result.Damage)
			},

			OnPeriodicDamageTaken: func(_ *core.Aura, sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				damageTakenHandler(sim, result.Damage)
			},

			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				// SW:D
				if priest.ShadowWordDeath.IsEqual(spell) && result.Landed() {
					priest.AddMana(sim, priest.MaxMana()*0.05*float64(priest.Talents.Masochism), manaMetrics)
					return
				}
			},
		},
	)
}

func (priest *Priest) applyMindMelt() {
	if priest.Talents.MindMelt == 0 {
		return
	}

	priest.MindMeltProcAura = priest.RegisterAura(core.Aura{
		Label:     "Mind Melt Proc",
		ActionID:  core.ActionID{SpellID: 87160},
		Duration:  time.Second * 15,
		MaxStacks: 2,
		OnStacksChange: func(_ *core.Aura, _ *core.Simulation, _ int32, newStacks int32) {
			AddOrReplaceMod(&priest.CastTimePercentMods, &PriestAuraMod[float64]{
				SpellID:    87160,
				Stacks:     newStacks,
				BaseValue:  -0.5,
				ClassSpell: PriestSpellMindBlast,
			})
		},
		OnExpire: func(_ *core.Aura, _ *core.Simulation) {
			RemoveMod(&priest.CastTimePercentMods, 87160)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if priest.MindBlast.IsEqual(spell) {
				aura.Deactivate(sim)
			}
		},
	})

	priest.RegisterAura(core.Aura{
		Label:    "Mind Melt",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && priest.MindSpike.IsEqual(spell) {
				priest.MindMeltProcAura.Activate(sim)
				priest.MindMeltProcAura.AddStack(sim)
			}
		},
	})
}

func (priest *Priest) applySinAndPunishment() {
	if priest.Talents.SinAndPunishment == 0 {
		return
	}

	priest.RegisterAura(core.Aura{
		Label:    "SinAndPunishment",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// check for critical mind flay tick
			if result.Outcome.Matches(core.OutcomeCrit) && priest.MindFlayAPL.IsEqual(spell) {

				// reduce cooldown
				remaining := max(0, time.Duration(*priest.Shadowfiend.CD.Timer)-time.Second*5)
				priest.Shadowfiend.CD.Set(remaining)
			}
		},
	})
}

func (priest *Priest) applyShadowyApparition() {
	if priest.Talents.ShadowyApparition == 0 {
		return
	}

	const spellScaling = 0.515
	const levelScaling = 0.514

	priest.ShadowyApparition = priest.RegisterSpell(PriestSpellShadowyApparation, core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 87532},
		MissileSpeed: 3.5,
		ProcMask:     core.ProcMaskEmpty, // summoned guardian, should not be able to proc stuff - verify
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: 0,
				GCD:  0,
			},
		},

		SpellSchool:      core.SpellSchoolShadow,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := priest.ScalingBaseDamage*levelScaling + spellScaling*spell.SpellPower()

			// snapshot values on spawn
			dmgMulti := priest.GetClassSpellDamageDonePercent(PriestSpellShadowyApparation, core.SpellSchoolShadow)
			dmgMultiAdd := priest.GetClassSpellDamageDonePercent(PriestSpellShadowyApparation, core.SpellSchoolShadow)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {

				// calculate dmg on hit, as the apparations profit from the debuffs on the target
				// when they reach them
				// spell and other modifiers are snapshotted when the apparations spawn
				spell.DamageMultiplier = dmgMulti
				spell.DamageMultiplierAdditive = dmgMultiAdd

				result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
				spell.DealDamage(sim, result)
			})
		},
	})

	priest.RegisterAura(core.Aura{
		Label:    "Shadowy Apparition Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},

		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell != priest.ShadowWordPain.Spell || !result.Landed() {
				return
			}

			procChance := priest.GetClassSpellProcChance(0.04*float64(priest.Talents.ShadowyApparition), PriestSpellShadowyApparation, core.SpellSchoolShadow)
			if sim.RandomFloat("Shadowy Apparition") < procChance {
				priest.ShadowyApparition.Cast(sim, result.Target)
			}
		},
	})
}

// func (priest *Priest) applyDivineAegis() {
// 	if priest.Talents.DivineAegis == 0 {
// 		return
// 	}

// 	divineAegis := priest.RegisterSpell(core.SpellConfig{
// 		ActionID:    core.ActionID{SpellID: 47515},
// 		SpellSchool: core.SpellSchoolHoly,
// 		ProcMask:    core.ProcMaskSpellHealing,
// 		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,

// 		DamageMultiplier: 1 *
// 			(0.1 * float64(priest.Talents.DivineAegis)) *
// 			core.TernaryFloat64(priest.HasSetBonus(ItemSetZabrasRaiment, 4), 1.1, 1),
// 		ThreatMultiplier: 1,

// 		Shield: core.ShieldConfig{
// 			Aura: core.Aura{
// 				Label:    "Divine Aegis",
// 				Duration: time.Second * 12,
// 			},
// 		},
// 	})

// 	priest.RegisterAura(core.Aura{
// 		Label:    "Divine Aegis Talent",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if result.Outcome.Matches(core.OutcomeCrit) {
// 				divineAegis.Shield(result.Target).Apply(sim, result.Damage)
// 			}
// 		},
// 	})
// }

// func (priest *Priest) applyGrace() {
// 	if priest.Talents.Grace == 0 {
// 		return
// 	}

// 	procChance := .5 * float64(priest.Talents.Grace)

// 	auras := make([]*core.Aura, len(priest.Env.AllUnits))
// 	for _, unit := range priest.Env.AllUnits {
// 		if !priest.IsOpponent(unit) {
// 			aura := unit.RegisterAura(core.Aura{
// 				Label:     "Grace" + strconv.Itoa(int(priest.Index)),
// 				ActionID:  core.ActionID{SpellID: 47517},
// 				Duration:  time.Second * 15,
// 				MaxStacks: 3,
// 				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
// 					priest.AttackTables[aura.Unit.UnitIndex].HealingDealtMultiplier /= 1 + .03*float64(oldStacks)
// 					priest.AttackTables[aura.Unit.UnitIndex].HealingDealtMultiplier *= 1 + .03*float64(newStacks)
// 				},
// 			})
// 			auras[unit.UnitIndex] = aura
// 		}
// 	}

// 	priest.RegisterAura(core.Aura{
// 		Label:    "Grace Talent",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if spell == priest.FlashHeal || spell == priest.GreaterHeal || spell == priest.PenanceHeal {
// 				if sim.Proc(procChance, "Grace") {
// 					aura := auras[result.Target.UnitIndex]
// 					aura.Activate(sim)
// 					aura.AddStack(sim)
// 				}
// 			}
// 		},
// 	})
// }

// // This one is called from healing priest sim initialization because it needs an input.
// func (priest *Priest) ApplyRapture(ppm float64) {
// 	if priest.Talents.Rapture == 0 {
// 		return
// 	}

// 	if ppm <= 0 {
// 		return
// 	}

// 	raptureManaCoeff := []float64{0, .015, .020, .025}[priest.Talents.Rapture]
// 	raptureMetrics := priest.NewManaMetrics(core.ActionID{SpellID: 47537})

// 	priest.RegisterResetEffect(func(sim *core.Simulation) {
// 		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
// 			Period: time.Minute / time.Duration(ppm),
// 			OnAction: func(sim *core.Simulation) {
// 				priest.AddMana(sim, raptureManaCoeff*priest.MaxMana(), raptureMetrics)
// 			},
// 		})
// 	})
// }

// func (priest *Priest) applyBorrowedTime() {
// 	if priest.Talents.BorrowedTime == 0 {
// 		return
// 	}

// 	multiplier := 1 + .05*float64(priest.Talents.BorrowedTime)

// 	procAura := priest.RegisterAura(core.Aura{
// 		Label:    "Borrowed Time",
// 		ActionID: core.ActionID{SpellID: 52800},
// 		Duration: time.Second * 6,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			priest.MultiplyCastSpeed(multiplier)
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			priest.MultiplyCastSpeed(1 / multiplier)
// 		},
// 	})

// 	priest.RegisterAura(core.Aura{
// 		Label:    "Borrwed Time Talent",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
// 			if spell == priest.PowerWordShield {
// 				procAura.Activate(sim)
// 			} else if spell.CurCast.CastTime > 0 {
// 				procAura.Deactivate(sim)
// 			}
// 		},
// 	})
// }

// func (priest *Priest) applyInspiration() {
// 	if priest.Talents.Inspiration == 0 {
// 		return
// 	}

// 	auras := make([]*core.Aura, len(priest.Env.AllUnits))
// 	for _, unit := range priest.Env.AllUnits {
// 		if !priest.IsOpponent(unit) {
// 			aura := core.InspirationAura(unit, priest.Talents.Inspiration)
// 			auras[unit.UnitIndex] = aura
// 		}
// 	}

// 	priest.RegisterAura(core.Aura{
// 		Label:    "Inspiration Talent",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if spell == priest.FlashHeal ||
// 				spell == priest.GreaterHeal ||
// 				spell == priest.BindingHeal ||
// 				spell == priest.PrayerOfMending ||
// 				spell == priest.PrayerOfHealing ||
// 				spell == priest.CircleOfHealing ||
// 				spell == priest.PenanceHeal {
// 				auras[result.Target.UnitIndex].Activate(sim)
// 			}
// 		},
// 	})
// }

// func (priest *Priest) applyHolyConcentration() {
// 	if priest.Talents.HolyConcentration == 0 {
// 		return
// 	}

// 	multiplier := 1 + []float64{0, .16, .32, .50}[priest.Talents.HolyConcentration]

// 	procAura := priest.RegisterAura(core.Aura{
// 		Label:    "Holy Concentration",
// 		ActionID: core.ActionID{SpellID: 34860},
// 		Duration: time.Second * 8,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			priest.PseudoStats.SpiritRegenMultiplier *= multiplier
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			priest.PseudoStats.SpiritRegenMultiplier /= multiplier
// 		},
// 	})

// 	priest.RegisterAura(core.Aura{
// 		Label:    "Holy Concentration Talent",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if result.DidCrit() &&
// 				(spell == priest.FlashHeal || spell == priest.GreaterHeal || spell == priest.EmpoweredRenew) {
// 				procAura.Activate(sim)
// 			}
// 		},
// 	})
// }

// func (priest *Priest) applySerendipity() {
// 	if priest.Talents.Serendipity == 0 {
// 		return
// 	}

// 	reductionPerStack := .04 * float64(priest.Talents.Serendipity)

// 	procAura := priest.RegisterAura(core.Aura{
// 		Label:     "Serendipity",
// 		ActionID:  core.ActionID{SpellID: 63737},
// 		Duration:  time.Second * 20,
// 		MaxStacks: 3,
// 		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
// 			priest.PrayerOfHealing.CastTimeMultiplier += reductionPerStack * float64(oldStacks)
// 			priest.PrayerOfHealing.CastTimeMultiplier -= reductionPerStack * float64(newStacks)
// 			priest.GreaterHeal.CastTimeMultiplier += reductionPerStack * float64(oldStacks)
// 			priest.GreaterHeal.CastTimeMultiplier -= reductionPerStack * float64(newStacks)
// 		},
// 	})

// 	priest.RegisterAura(core.Aura{
// 		Label:    "Serendipity Talent",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
// 			if spell == priest.FlashHeal || spell == priest.BindingHeal {
// 				procAura.Activate(sim)
// 				procAura.AddStack(sim)
// 			} else if spell == priest.GreaterHeal || spell == priest.PrayerOfHealing {
// 				procAura.Deactivate(sim)
// 			}
// 		},
// 	})
// }

// func (priest *Priest) applySurgeOfLight() {
// 	if priest.Talents.SurgeOfLight == 0 {
// 		return
// 	}

// 	procHandler := func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 		if spell == priest.Smite || spell == priest.FlashHeal {
// 			aura.Deactivate(sim)
// 		}
// 	}

// 	priest.SurgeOfLightProcAura = priest.RegisterAura(core.Aura{
// 		Label:    "Surge of Light Proc",
// 		ActionID: core.ActionID{SpellID: 33154},
// 		Duration: time.Second * 10,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			if priest.Smite != nil {
// 				priest.Smite.CastTimeMultiplier -= 1
// 				priest.Smite.CostMultiplier -= 1
// 				priest.Smite.BonusCritRating -= 100 * core.CritRatingPerCritChance
// 			}
// 			if priest.FlashHeal != nil {
// 				priest.FlashHeal.CastTimeMultiplier -= 1
// 				priest.FlashHeal.CostMultiplier -= 1
// 				priest.FlashHeal.BonusCritRating -= 100 * core.CritRatingPerCritChance
// 			}
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			if priest.Smite != nil {
// 				priest.Smite.CastTimeMultiplier += 1
// 				priest.Smite.CostMultiplier += 1
// 				priest.Smite.BonusCritRating += 100 * core.CritRatingPerCritChance
// 			}
// 			if priest.FlashHeal != nil {
// 				priest.FlashHeal.CastTimeMultiplier += 1
// 				priest.FlashHeal.CostMultiplier += 1
// 				priest.FlashHeal.BonusCritRating += 100 * core.CritRatingPerCritChance
// 			}
// 		},
// 		OnSpellHitDealt: procHandler,
// 		OnHealDealt:     procHandler,
// 	})

// 	procChance := 0.25 * float64(priest.Talents.SurgeOfLight)
// 	icd := core.Cooldown{
// 		Timer:    priest.NewTimer(),
// 		Duration: time.Second * 6,
// 	}
// 	priest.SurgeOfLightProcAura.Icd = &icd

// 	handler := func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 		if icd.IsReady(sim) && result.Outcome.Matches(core.OutcomeCrit) && sim.RandomFloat("SurgeOfLight") < procChance {
// 			icd.Use(sim)
// 			priest.SurgeOfLightProcAura.Activate(sim)
// 		}
// 	}

// 	priest.RegisterAura(core.Aura{
// 		Label:    "Surge of Light",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: handler,
// 		OnHealDealt:     handler,
// 	})
// }

// func (priest *Priest) applyMisery() {
// 	if priest.Talents.Misery == 0 {
// 		return
// 	}

// 	miseryAuras := priest.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
// 		return core.MiseryAura(target, priest.Talents.Misery)
// 	})
// 	priest.Env.RegisterPreFinalizeEffect(func() {
// 		priest.ShadowWordPain.RelatedAuras = append(priest.ShadowWordPain.RelatedAuras, miseryAuras)
// 		if priest.VampiricTouch != nil {
// 			priest.VampiricTouch.RelatedAuras = append(priest.VampiricTouch.RelatedAuras, miseryAuras)
// 		}
// 		if priest.MindFlay[1] != nil {
// 			priest.MindFlayAPL.RelatedAuras = append(priest.MindFlayAPL.RelatedAuras, miseryAuras)
// 			priest.MindFlay[1].RelatedAuras = append(priest.MindFlay[1].RelatedAuras, miseryAuras)
// 			priest.MindFlay[2].RelatedAuras = append(priest.MindFlay[2].RelatedAuras, miseryAuras)
// 			priest.MindFlay[3].RelatedAuras = append(priest.MindFlay[3].RelatedAuras, miseryAuras)
// 		}
// 	})

// 	priest.RegisterAura(core.Aura{
// 		Label:    "Priest Shadow Effects",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if !result.Landed() {
// 				return
// 			}

// 			if spell == priest.ShadowWordPain || spell == priest.VampiricTouch || spell.ActionID.SpellID == priest.MindFlay[1].ActionID.SpellID {
// 				miseryAuras.Get(result.Target).Activate(sim)
// 			}
// 		},
// 	})
// }

// func (priest *Priest) applyShadowWeaving() {
// 	if priest.Talents.ShadowWeaving == 0 {
// 		return
// 	}

// 	priest.ShadowWeavingAura = priest.GetOrRegisterAura(core.Aura{
// 		Label:     "Shadow Weaving",
// 		ActionID:  core.ActionID{SpellID: 15258},
// 		Duration:  time.Second * 15,
// 		MaxStacks: 5,
// 		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
// 			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= 1.0 + 0.02*float64(oldStacks)
// 			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.0 + 0.02*float64(newStacks)
// 		},
// 	})
// }

// func (priest *Priest) applyImprovedSpiritTap() {
// 	if priest.Talents.ImprovedSpiritTap == 0 {
// 		return
// 	}

// 	increase := 1 + 0.05*float64(priest.Talents.ImprovedSpiritTap)
// 	statDep := priest.NewDynamicMultiplyStat(stats.Spirit, increase)
// 	regen := []float64{0, 0.17, 0.33}[priest.Talents.ImprovedSpiritTap]

// 	priest.ImprovedSpiritTap = priest.GetOrRegisterAura(core.Aura{
// 		Label:    "Improved Spirit Tap",
// 		ActionID: core.ActionID{SpellID: 59000},
// 		Duration: time.Second * 8,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			priest.EnableDynamicStatDep(sim, statDep)
// 			priest.PseudoStats.SpiritRegenRateCasting += regen
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			priest.DisableDynamicStatDep(sim, statDep)
// 			priest.PseudoStats.SpiritRegenRateCasting -= regen
// 		},
// 	})
// }

// func (priest *Priest) registerInnerFocus() {
// 	if !priest.Talents.InnerFocus {
// 		return
// 	}

// 	actionID := core.ActionID{SpellID: 14751}

// 	priest.InnerFocusAura = priest.RegisterAura(core.Aura{
// 		Label:    "Inner Focus",
// 		ActionID: actionID,
// 		Duration: core.NeverExpires,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Unit.AddStatDynamic(sim, stats.SpellCrit, 25*core.CritRatingPerCritChance)
// 			aura.Unit.PseudoStats.CostMultiplier -= 1
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Unit.AddStatDynamic(sim, stats.SpellCrit, -25*core.CritRatingPerCritChance)
// 			aura.Unit.PseudoStats.CostMultiplier += 1
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			// Remove the buff and put skill on CD
// 			aura.Deactivate(sim)
// 			priest.InnerFocus.CD.Use(sim)
// 			priest.UpdateMajorCooldowns()
// 		},
// 	})

// 	priest.InnerFocus = priest.RegisterSpell(core.SpellConfig{
// 		ActionID: actionID,
// 		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

// 		Cast: core.CastConfig{
// 			CD: core.Cooldown{
// 				Timer:    priest.NewTimer(),
// 				Duration: time.Duration(float64(time.Minute*3) * (1 - .1*float64(priest.Talents.Aspiration))),
// 			},
// 		},

// 		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
// 			priest.InnerFocusAura.Activate(sim)
// 		},
// 	})
// }
