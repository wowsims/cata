package priest

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func (priest *Priest) ApplyTalents() {
	priest.ApplyArmorSpecializationEffect(stats.Intellect, proto.ArmorType_ArmorTypeCloth, 89745)
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
		priest.AddStaticMod(core.SpellModConfig{
			School:     core.SpellSchoolHoly | core.SpellSchoolShadow,
			FloatValue: (0.02 * float64(priest.Talents.TwinDisciplines)),
			Kind:       core.SpellMod_DamageDone_Pct,
		})
	}

	// Mental Agillity
	if priest.Talents.MentalAgility > 0 {
		priest.AddStaticMod(core.SpellModConfig{
			ClassMask: PriestSpellInstant,
			IntValue:  -4 * priest.Talents.MentalAgility,
			Kind:      core.SpellMod_PowerCost_Pct,
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
		priest.AddStaticMod(core.SpellModConfig{
			ClassMask:  PriestSpellShadowWordPain,
			FloatValue: 0.03 * float64(priest.Talents.ImprovedShadowWordPain),
			Kind:       core.SpellMod_DamageDone_Flat,
		})
	}

	// Veiled Shadows
	if priest.Talents.VeiledShadows > 0 {
		priest.AddStaticMod(core.SpellModConfig{
			ClassMask: PriestSpellFade,
			TimeValue: time.Second * -3 * time.Duration(priest.Talents.VeiledShadows),
			Kind:      core.SpellMod_Cooldown_Flat,
		})

		priest.AddStaticMod(core.SpellModConfig{
			ClassMask: PriestSpellShadowFiend,
			TimeValue: time.Second * -30 * time.Duration(priest.Talents.VeiledShadows),
			Kind:      core.SpellMod_Cooldown_Flat,
		})
	}

	// Improved Psychic Scream
	if priest.Talents.ImprovedPsychicScream > 0 {
		priest.AddStaticMod(core.SpellModConfig{
			ClassMask: PriestSpellPsychicScream,
			TimeValue: time.Second * -2 * time.Duration(priest.Talents.ImprovedPsychicScream),
			Kind:      core.SpellMod_Cooldown_Flat,
		})
	}

	// Improved Mind Blast
	priest.applyImprovedMindBlast()

	// Improved Devouring Plague
	priest.applyImprovedDevouringPlague()

	// Twisted Faith
	if priest.Talents.TwistedFaith > 0 {
		priest.AddStaticMod(core.SpellModConfig{
			School:     core.SpellSchoolShadow,
			FloatValue: 0.01 * float64(priest.Talents.TwistedFaith),
			Kind:       core.SpellMod_DamageDone_Pct,
		})

		// Twisted Faith is not applied to base spirit
		priest.AddStat(stats.SpellHitPercent, -0.5*float64(priest.Talents.TwistedFaith)*priest.GetBaseStats()[stats.Spirit]/core.SpellHitRatingPerHitPercent)
		priest.AddStatDependency(stats.Spirit, stats.SpellHitPercent, 0.5*float64(priest.Talents.TwistedFaith)/core.SpellHitRatingPerHitPercent)
	}

	// Shadowform
	if priest.Talents.Shadowform {
		// no class restrictions
		priest.AddStaticMod(core.SpellModConfig{
			School:     core.SpellSchoolShadow,
			FloatValue: 0.15,
			Kind:       core.SpellMod_DamageDone_Pct,
		})
	}

	// Phantasm: Not implemented
	// Harnessed Shadows - shadow.go

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

	priest.ApplyGlyphs()
}

// disciplin talents
func (priest *Priest) applyEvangelism() {
	if priest.Talents.Evangelism == 0 {
		return
	}

	darkEvangelismMod := priest.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.0,
		ClassMask:  PriestSpellDoT,
	})

	priest.DarkEvangelismProcAura = priest.GetOrRegisterAura(core.Aura{
		Label:     "Dark EvangelismProc",
		ActionID:  core.ActionID{SpellID: 87118},
		Duration:  time.Second * 20,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			darkEvangelismMod.UpdateFloatValue(0.02 * float64(newStacks))
			darkEvangelismMod.Activate()
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			darkEvangelismMod.Deactivate()
		},
	})

	evangelismDmgMod := priest.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.0,
		ClassMask:  PriestSpellSmite | PriestSpellHolyFire | PriestSpellPenance,
	})

	evangelismManaMod := priest.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		IntValue:  0,
		ClassMask: PriestSpellSmite | PriestSpellHolyFire | PriestSpellPenance,
	})

	priest.HolyEvangelismProcAura = priest.GetOrRegisterAura(core.Aura{
		Label:     "EvangelismProc",
		ActionID:  core.ActionID{SpellID: 81661},
		Duration:  time.Second * 20,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, _ int32, newStacks int32) {
			evangelismDmgMod.UpdateFloatValue(0.04 * float64(newStacks))
			evangelismDmgMod.Activate()

			evangelismManaMod.UpdateIntValue(-6 * newStacks)
			evangelismManaMod.Activate()
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			evangelismDmgMod.Deactivate()
			evangelismManaMod.Deactivate()
		},
	})

	core.MakeProcTriggerAura(&priest.Unit, core.ProcTrigger{
		Name:           "Evangilism Hit",
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: PriestSpellSmite | PriestSpellHolyFire | PriestSpellMindFlay,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ClassSpellMask == PriestSpellMindFlay {
				priest.AddDarkEvangelismStack(sim)
				return
			}
			priest.AddHolyEvanglismStack(sim)
		},
	})

	core.MakeProcTriggerAura(&priest.Unit, core.ProcTrigger{
		Name:           "Evangilism Tick",
		Callback:       core.CallbackOnPeriodicDamageDealt,
		ClassSpellMask: PriestSpellMindFlay,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			priest.AddDarkEvangelismStack(sim)
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

	darkArchAngelMod := priest.AddDynamicMod(core.SpellModConfig{
		ClassMask:  PriestSpellMindFlay | PriestSpellMindSpike | PriestSpellMindBlast | PriestSpellShadowWordDeath,
		FloatValue: 0.04,
		Kind:       core.SpellMod_DamageDone_Flat,
	})

	darkArchAngelAura := priest.Unit.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 87153},
		Label:     "Dark Archangel Aura",
		MaxStacks: 5,
		Duration:  time.Second * 18,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			if newStacks > oldStacks {
				priest.AddMana(sim, 0.05*priest.MaxMana()*float64((newStacks-oldStacks)), darkArchAngelMana)
			}

			darkArchAngelMod.UpdateFloatValue(0.04 * float64(newStacks))
			darkArchAngelMod.Activate()
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			darkArchAngelMod.Deactivate()
		},
	})

	priest.RegisterSpell(core.SpellConfig{
		ActionID:                 core.ActionID{SpellID: 87151},
		SpellSchool:              core.SpellSchoolHoly,
		ProcMask:                 core.ProcMaskEmpty,
		Flags:                    core.SpellFlagAPL,
		ClassSpellMask:           PriestSpellArchangel,
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           priest.DefaultCritMultiplier(),
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 0,
		},
		Cast: core.CastConfig{
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
	priest.RegisterSpell(core.SpellConfig{
		ActionID:                 core.ActionID{SpellID: 87153},
		SpellSchool:              core.SpellSchoolHoly,
		ProcMask:                 core.ProcMaskEmpty,
		Flags:                    core.SpellFlagAPL,
		ClassSpellMask:           PriestSpellDarkArchangel,
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           priest.DefaultCritMultiplier(),
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 0,
		},
		Cast: core.CastConfig{
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

	priest.AddStaticMod(core.SpellModConfig{
		ClassMask: PriestSpellMindBlast,
		TimeValue: time.Duration(priest.Talents.ImprovedMindBlast) * time.Millisecond * -500,
		Kind:      core.SpellMod_Cooldown_Flat,
	})

	mindTraumaAura := priest.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return MindTraumaAura(target)
	})

	mindTraumaSpell := priest.RegisterSpell(core.SpellConfig{
		ActionID:                 core.ActionID{SpellID: 48301},
		ProcMask:                 core.ProcMaskSpellProc,
		SpellSchool:              core.SpellSchoolShadow,
		ClassSpellMask:           PriestSpellMindTrauma,
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           priest.DefaultCritMultiplier(),
		Flags:                    core.SpellFlagNoMetrics,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			mindTraumaAura.Get(target).Activate(sim)
		},
	})

	procChance := []float64{0.0, 0.33, 0.66, 1.0}[priest.Talents.ImprovedMindBlast]
	core.MakePermanent(priest.RegisterAura(core.Aura{
		Label: "Improved Mind Blast",
		OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.ClassSpellMask == PriestSpellMindBlast {
				if sim.Proc(procChance, "Improved Mind Blast") {
					mindTraumaSpell.Cast(sim, result.Target)
				}
			}
		},
	}))
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
	impDPDamage := priest.RegisterSpell(core.SpellConfig{

		// TODO: improve metric aggregation to show correct DPC
		ActionID:                 core.ActionID{SpellID: 2944, Tag: 1},
		SpellSchool:              core.SpellSchoolShadow,
		ProcMask:                 core.ProcMaskProc,
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,
		ClassSpellMask:           PriestSpellImprovedDevouringPlague,
		Flags:                    core.SpellFlagPassiveSpell,
		CritMultiplier:           priest.DefaultCritMultiplier(),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dot := priest.DevouringPlague.Dot(target)
			dpTickDamage := dot.SnapshotBaseDamage

			// Improved Devouring Plague only considers haste on gear nothing else for dot tick frequency
			// https://github.com/JamminL/cata-classic-bugs/issues/971
			tickPeriod := float64(dot.BaseTickLength) / (1 + (priest.GetStat(stats.HasteRating) / (core.HasteRatingPerHastePercent * 100)))
			ticks := math.Ceil(float64(dot.BaseDuration()) / tickPeriod)
			dmg := ticks * dpTickDamage * float64(priest.Talents.ImprovedDevouringPlague) * 0.15
			spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMagicCrit)
		},
	})

	core.MakePermanent(priest.RegisterAura(core.Aura{
		Label: "Improved Devouring Plague Talent",

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if priest.DevouringPlague != spell || !result.Landed() {
				return
			}

			impDPDamage.Cast(sim, result.Target)
		},
	}))
}

func (priest *Priest) applyPainAndSuffering() {
	if priest.Talents.PainAndSuffering == 0 {
		return
	}

	core.MakeProcTriggerAura(&priest.Unit, core.ProcTrigger{
		Name:           "Pain and Suffering",
		Callback:       core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		Outcome:        core.OutcomeLanded,
		ProcChance:     float64(priest.Talents.PainAndSuffering) * 0.3,
		ClassSpellMask: PriestSpellMindFlay,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			swp := priest.ShadowWordPain.Dot(result.Target)
			if swp.IsActive() {
				swp.Apply(sim)
			}
		},
	})
}

func (priest *Priest) applyMasochism() {
	if priest.Talents.Masochism == 0 {
		return
	}

	manaMetrics := priest.NewManaMetrics(core.ActionID{
		SpellID: []int32{0, 88894, 88995}[priest.Talents.Masochism],
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
				if spell.ClassSpellMask == PriestSpellShadowWordDeath && result.Landed() {
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

	mindMeltMod := priest.AddDynamicMod(core.SpellModConfig{
		ClassMask:  PriestSpellMindBlast,
		FloatValue: -0.5,
		Kind:       core.SpellMod_CastTime_Pct,
	})

	mindMeltSWDMod := priest.AddDynamicMod(core.SpellModConfig{
		ClassMask:  PriestSpellShadowWordDeath,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.3,
	})

	priest.RegisterResetEffect(func(s *core.Simulation) {
		mindMeltSWDMod.Deactivate()
		s.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
			if isExecute == 25 {
				mindMeltSWDMod.Activate()
			}
		})
	})

	procAura := priest.RegisterAura(core.Aura{
		Label:     "Mind Melt Proc",
		ActionID:  core.ActionID{SpellID: 87160},
		Duration:  time.Second * 6,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mindMeltMod.Activate()
		},
		OnStacksChange: func(_ *core.Aura, _ *core.Simulation, oldStacks int32, newStacks int32) {
			mindMeltMod.UpdateFloatValue(-0.5 * float64(newStacks))
		},

		OnExpire: func(_ *core.Aura, _ *core.Simulation) {
			mindMeltMod.Deactivate()
		},

		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask == PriestSpellMindBlast {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakeProcTriggerAura(&priest.Unit, core.ProcTrigger{
		Name:           "Mind Melt",
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ProcChance:     1,
		ClassSpellMask: PriestSpellMindSpike,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			procAura.Activate(sim)
			procAura.AddStack(sim)
		},
	})
}

func (priest *Priest) applySinAndPunishment() {
	if priest.Talents.SinAndPunishment == 0 {
		return
	}

	core.MakeProcTriggerAura(&priest.Unit, core.ProcTrigger{
		Name:           "Sin And Punishment",
		Callback:       core.CallbackOnPeriodicDamageDealt,
		Outcome:        core.OutcomeCrit,
		ProcChance:     1,
		ClassSpellMask: PriestSpellMindFlay,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.Log != nil {
				sim.Log("Sin and Punishment Proc. New CD: %d", sim.CurrentTime)
			}

			priest.Shadowfiend.CD.Set(priest.Shadowfiend.CD.ReadyAt() - time.Duration(priest.Talents.SinAndPunishment)*5*time.Second)
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
// 			core.TernaryFloat64(priest.CouldHaveSetBonus(ItemSetZabrasRaiment, 4), 1.1, 1),
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
// 				priest.Smite.Cost.Multiplier -= 100
// 				priest.Smite.BonusCritRating -= 100 * core.CritRatingPerCritChance
// 			}
// 			if priest.FlashHeal != nil {
// 				priest.FlashHeal.CastTimeMultiplier -= 1
// 				priest.FlashHeal.Cost.Multiplier -= 100
// 				priest.FlashHeal.BonusCritRating -= 100 * core.CritRatingPerCritChance
// 			}
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			if priest.Smite != nil {
// 				priest.Smite.CastTimeMultiplier += 1
// 				priest.Smite.Cost.Multiplier += 100
// 				priest.Smite.BonusCritRating += 100 * core.CritRatingPerCritChance
// 			}
// 			if priest.FlashHeal != nil {
// 				priest.FlashHeal.CastTimeMultiplier += 1
// 				priest.FlashHeal.Cost.Multiplier += 100
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
// 			aura.Unit.PseudoStats.CostMultiplier -= 100
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Unit.AddStatDynamic(sim, stats.SpellCrit, -25*core.CritRatingPerCritChance)
// 			aura.Unit.PseudoStats.CostMultiplier += 100
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
