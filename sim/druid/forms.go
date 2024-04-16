package druid

import (
	//"math"

	"github.com/wowsims/cata/sim/core"
	//"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

 type DruidForm uint8

 const (
 	Humanoid DruidForm = 1 << iota
 	Bear
 	Cat
 	Moonkin
 	Tree
 	Any = Humanoid | Bear | Cat | Moonkin | Tree
 )

// Converts from 0.009327 to 0.0085
const AnimalSpiritRegenSuppression = 0.911337

func (form DruidForm) Matches(other DruidForm) bool {
	return (form & other) != 0
}

// func (druid *Druid) GetForm() DruidForm {
// 	return druid.form
// }

func (druid *Druid) InForm(form DruidForm) bool {
	return druid.form.Matches(form)
}

func (druid *Druid) ClearForm(sim *core.Simulation) {
	if druid.InForm(Cat) {
		druid.CatFormAura.Deactivate(sim)
	} else if druid.InForm(Bear) {
		druid.BearFormAura.Deactivate(sim)
	} else if druid.InForm(Moonkin) {
		panic("cant clear moonkin form")
	}
	druid.form = Humanoid
	druid.SetCurrentPowerBar(core.ManaBar)
}

func (druid *Druid) GetCatWeapon() core.Weapon {
	unscaledWeapon := druid.WeaponFromMainHand(0)
	return core.Weapon{
		BaseDamageMin:        unscaledWeapon.BaseDamageMin / unscaledWeapon.SwingSpeed,
		BaseDamageMax:        unscaledWeapon.BaseDamageMax / unscaledWeapon.SwingSpeed,
		SwingSpeed:           1.0,
		NormalizedSwingSpeed: 1.0,
		CritMultiplier:       druid.MeleeCritMultiplier(1.0, 0.0),
		AttackPowerPerDPS:    core.DefaultAttackPowerPerDPS,
	}
}

// func (druid *Druid) GetBearWeapon() core.Weapon {
// 	return core.Weapon{
// 		BaseDamageMin:        109,
// 		BaseDamageMax:        165,
// 		SwingSpeed:           2.5,
// 		NormalizedSwingSpeed: 2.5,
// 		CritMultiplier:       druid.MeleeCritMultiplier(Bear),
// 		AttackPowerPerDPS:    core.DefaultAttackPowerPerDPS,
// 	}
// }

func (druid *Druid) registerCatFormSpell() {
	actionID := core.ActionID{SpellID: 768}

	srm := druid.getSavageRoarMultiplier()

	statBonus := stats.Stats{
		stats.AttackPower: -20, // This offset is needed because the first 10 points of Agility do not contribute any Attack Power.
		stats.MeleeCrit:   core.TernaryFloat64(druid.Talents.MasterShapeshifter, 4.0 * core.CritRatingPerCritChance, 0.0),
	}

	agiApDep := druid.NewDynamicStatDependency(stats.Agility, stats.AttackPower, 2)

	var hotwDep *stats.StatDependency
	if druid.Talents.HeartOfTheWild > 0 {
		hotwDep = druid.NewDynamicMultiplyStat(stats.AttackPower, 1.0+0.02*float64(druid.Talents.HeartOfTheWild))
	}

	clawWeapon := druid.GetCatWeapon()

	druid.CatFormAura = druid.RegisterAura(core.Aura{
		Label:      "Cat Form",
		ActionID:   actionID,
		Duration:   core.NeverExpires,
		BuildPhase: core.Ternary(druid.StartingForm.Matches(Cat), core.CharacterBuildPhaseBase, core.CharacterBuildPhaseNone),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if !druid.Env.MeasuringStats && druid.form != Humanoid {
				druid.ClearForm(sim)
			}
			druid.form = Cat
			druid.SetCurrentPowerBar(core.EnergyBar)

			druid.AutoAttacks.SetMH(clawWeapon)

			druid.PseudoStats.ThreatMultiplier *= 0.71
			druid.PseudoStats.SpiritRegenMultiplier *= AnimalSpiritRegenSuppression
			druid.PseudoStats.BaseDodge += 0.02 * float64(druid.Talents.FeralSwiftness)

			druid.AddStatsDynamic(sim, statBonus)
			druid.EnableDynamicStatDep(sim, agiApDep)
			if hotwDep != nil {
				druid.EnableDynamicStatDep(sim, hotwDep)
			}

			if !druid.Env.MeasuringStats {
				druid.AutoAttacks.SetReplaceMHSwing(nil)
				druid.AutoAttacks.EnableAutoSwing(sim)
				druid.UpdateManaRegenRates()

				// These buffs stay up, but corresponding changes don't
				if druid.SavageRoarAura.IsActive() {
					druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= srm
				}

				if druid.PredatoryInstinctsAura != nil {
					druid.PredatoryInstinctsAura.Activate(sim)
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.form = Humanoid

			druid.AutoAttacks.SetMH(druid.WeaponFromMainHand(druid.MeleeCritMultiplier(1.0, 0.0)))

			druid.PseudoStats.ThreatMultiplier /= 0.71
			druid.PseudoStats.SpiritRegenMultiplier /= AnimalSpiritRegenSuppression
			druid.PseudoStats.BaseDodge -= 0.02 * float64(druid.Talents.FeralSwiftness)

			druid.AddStatsDynamic(sim, statBonus.Invert())
			druid.DisableDynamicStatDep(sim, agiApDep)
			if hotwDep != nil {
				druid.DisableDynamicStatDep(sim, hotwDep)
			}

			if !druid.Env.MeasuringStats {
				druid.AutoAttacks.SetReplaceMHSwing(nil)
				druid.AutoAttacks.EnableAutoSwing(sim)
				druid.UpdateManaRegenRates()

				druid.TigersFuryAura.Deactivate(sim)

				// These buffs stay up, but corresponding changes don't
				if druid.SavageRoarAura.IsActive() {
					druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= srm
				}

				if druid.PredatoryInstinctsAura != nil {
					druid.PredatoryInstinctsAura.Deactivate(sim)
				}
			}
		},
	})

	energyMetrics := druid.NewEnergyMetrics(actionID)

	druid.CatForm = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.35,
			Multiplier: (1 - 0.2*float64(druid.Talents.KingOfTheJungle)) * (1 - 0.1*float64(druid.Talents.NaturalShapeshifter)),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			maxShiftEnergy := float64(20 * druid.Talents.Furor)

			energyDelta := maxShiftEnergy - druid.CurrentEnergy()

			if energyDelta < 0 {
				druid.SpendEnergy(sim, -energyDelta, energyMetrics)
			}
			druid.CatFormAura.Activate(sim)
		},
	})
}

// func (druid *Druid) registerBearFormSpell() {
// 	actionID := core.ActionID{SpellID: 9634}
// 	healthMetrics := druid.NewHealthMetrics(actionID)

// 	statBonus := druid.GetFormShiftStats().Add(stats.Stats{
// 		stats.AttackPower: 3 * float64(druid.Level),
// 	})

// 	stamDep := druid.NewDynamicMultiplyStat(stats.Stamina, 1.25)

// 	var potpDep *stats.StatDependency
// 	if druid.Talents.ProtectorOfThePack > 0 {
// 		potpDep = druid.NewDynamicMultiplyStat(stats.AttackPower, 1.0+0.02*float64(druid.Talents.ProtectorOfThePack))
// 	}

// 	var hotwDep *stats.StatDependency
// 	if druid.Talents.HeartOfTheWild > 0 {
// 		hotwDep = druid.NewDynamicMultiplyStat(stats.Stamina, 1.0+0.02*float64(druid.Talents.HeartOfTheWild))
// 	}

// 	potpdtm := 1 - 0.04*float64(druid.Talents.ProtectorOfThePack)

// 	clawWeapon := druid.GetBearWeapon()
// 	predBonus := stats.Stats{}

// 	druid.BearFormAura = druid.RegisterAura(core.Aura{
// 		Label:      "Bear Form",
// 		ActionID:   actionID,
// 		Duration:   core.NeverExpires,
// 		BuildPhase: core.Ternary(druid.StartingForm.Matches(Bear), core.CharacterBuildPhaseBase, core.CharacterBuildPhaseNone),
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			if !druid.Env.MeasuringStats && druid.form != Humanoid {
// 				druid.ClearForm(sim)
// 			}
// 			druid.form = Bear
// 			druid.SetCurrentPowerBar(core.RageBar)

// 			druid.AutoAttacks.SetMH(clawWeapon)

// 			druid.PseudoStats.ThreatMultiplier *= 2.1021
// 			druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.0 + 0.02*float64(druid.Talents.MasterShapeshifter)
// 			druid.PseudoStats.DamageTakenMultiplier *= potpdtm
// 			druid.PseudoStats.SpiritRegenMultiplier *= AnimalSpiritRegenSuppression
// 			druid.PseudoStats.BaseDodge += 0.02 * float64(druid.Talents.FeralSwiftness+druid.Talents.NaturalReaction)

// 			predBonus = druid.GetDynamicPredStrikeStats()
// 			druid.AddStatsDynamic(sim, predBonus)
// 			druid.AddStatsDynamic(sim, statBonus)
// 			druid.ApplyDynamicEquipScaling(sim, stats.Armor, druid.BearArmorMultiplier())
// 			if potpDep != nil {
// 				druid.EnableDynamicStatDep(sim, potpDep)
// 			}

// 			// Preserve fraction of max health when shifting
// 			healthFrac := druid.CurrentHealth() / druid.MaxHealth()
// 			druid.EnableDynamicStatDep(sim, stamDep)
// 			if hotwDep != nil {
// 				druid.EnableDynamicStatDep(sim, hotwDep)
// 			}
// 			druid.GainHealth(sim, healthFrac*druid.MaxHealth()-druid.CurrentHealth(), healthMetrics)

// 			if !druid.Env.MeasuringStats {
// 				druid.AutoAttacks.SetReplaceMHSwing(druid.ReplaceBearMHFunc)
// 				druid.AutoAttacks.EnableAutoSwing(sim)

// 				druid.manageCooldownsEnabled()
// 				druid.UpdateManaRegenRates()
// 			}
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			druid.form = Humanoid
// 			druid.AutoAttacks.SetMH(druid.WeaponFromMainHand(druid.MeleeCritMultiplier(Humanoid)))

// 			druid.PseudoStats.ThreatMultiplier /= 2.1021
// 			druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.0 + 0.02*float64(druid.Talents.MasterShapeshifter)
// 			druid.PseudoStats.DamageTakenMultiplier /= potpdtm
// 			druid.PseudoStats.SpiritRegenMultiplier /= AnimalSpiritRegenSuppression
// 			druid.PseudoStats.BaseDodge -= 0.02 * float64(druid.Talents.FeralSwiftness+druid.Talents.NaturalReaction)

// 			druid.AddStatsDynamic(sim, predBonus.Invert())
// 			druid.AddStatsDynamic(sim, statBonus.Invert())
// 			druid.RemoveDynamicEquipScaling(sim, stats.Armor, druid.BearArmorMultiplier())
// 			if potpDep != nil {
// 				druid.DisableDynamicStatDep(sim, potpDep)
// 			}

// 			healthFrac := druid.CurrentHealth() / druid.MaxHealth()
// 			druid.DisableDynamicStatDep(sim, stamDep)
// 			if hotwDep != nil {
// 				druid.DisableDynamicStatDep(sim, hotwDep)
// 			}
// 			druid.RemoveHealth(sim, druid.CurrentHealth()-healthFrac*druid.MaxHealth())

// 			if !druid.Env.MeasuringStats {
// 				druid.AutoAttacks.SetReplaceMHSwing(nil)
// 				druid.AutoAttacks.EnableAutoSwing(sim)

// 				druid.manageCooldownsEnabled()
// 				druid.UpdateManaRegenRates()
// 				druid.EnrageAura.Deactivate(sim)
// 				druid.MaulQueueAura.Deactivate(sim)
// 			}
// 		},
// 	})

// 	rageMetrics := druid.NewRageMetrics(actionID)

// 	furorProcChance := []float64{0, 0.2, 0.4, 0.6, 0.8, 1}[druid.Talents.Furor]

// 	druid.BearForm = druid.RegisterSpell(Any, core.SpellConfig{
// 		ActionID: actionID,
// 		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

// 		ManaCost: core.ManaCostOptions{
// 			BaseCost:   0.35,
// 			Multiplier: (1 - 0.2*float64(druid.Talents.KingOfTheJungle)) * (1 - 0.1*float64(druid.Talents.NaturalShapeshifter)),
// 		},
// 		Cast: core.CastConfig{
// 			DefaultCast: core.Cast{
// 				GCD: core.GCDDefault,
// 			},
// 			IgnoreHaste: true,
// 		},

// 		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
// 			rageDelta := 0 - druid.CurrentRage()
// 			if sim.Proc(furorProcChance, "Furor") {
// 				rageDelta += 10
// 			}
// 			if rageDelta > 0 {
// 				druid.AddRage(sim, rageDelta, rageMetrics)
// 			} else if rageDelta < 0 {
// 				druid.SpendRage(sim, -rageDelta, rageMetrics)
// 			}
// 			druid.BearFormAura.Activate(sim)
// 		},
// 	})
// }

// func (druid *Druid) applyMoonkinForm() {
// 	if !druid.InForm(Moonkin) || !druid.Talents.MoonkinForm {
// 		return
// 	}

// 	druid.MultiplyStat(stats.Intellect, 1+(0.02*float64(druid.Talents.Furor)))
// 	druid.PseudoStats.DamageDealtMultiplier *= 1 + (float64(druid.Talents.MasterShapeshifter) * 0.02)
// 	if druid.Talents.ImprovedMoonkinForm > 0 {
// 		druid.AddStatDependency(stats.Spirit, stats.SpellPower, 0.1*float64(druid.Talents.ImprovedMoonkinForm))
// 	}

// 	manaMetrics := druid.NewManaMetrics(core.ActionID{SpellID: 24858})
// 	druid.RegisterAura(core.Aura{
// 		Label:    "Moonkin Form",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if result.DidCrit() {
// 				if druid.Moonfire.IsEqual(spell) || druid.Starfire.IsEqual(spell) || druid.Wrath.IsEqual(spell) {
// 					druid.AddMana(sim, 0.02*druid.MaxMana(), manaMetrics)
// 				}
// 			}
// 		},
// 	})
// }
