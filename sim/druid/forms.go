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
		CritMultiplier:       druid.DefaultMeleeCritMultiplier(),
		AttackPowerPerDPS:    core.DefaultAttackPowerPerDPS,
		MaxRange:             core.MaxMeleeRange,
	}
}

func (druid *Druid) GetBearWeapon() core.Weapon {
	unscaledWeapon := druid.WeaponFromMainHand(0)
	return core.Weapon{
		BaseDamageMin:        unscaledWeapon.BaseDamageMin / unscaledWeapon.SwingSpeed * 2.5,
		BaseDamageMax:        unscaledWeapon.BaseDamageMax / unscaledWeapon.SwingSpeed * 2.5,
		SwingSpeed:           2.5,
		NormalizedSwingSpeed: 2.5,
		CritMultiplier:       druid.DefaultMeleeCritMultiplier(),
		AttackPowerPerDPS:    core.DefaultAttackPowerPerDPS,
		MaxRange:             core.MaxMeleeRange,
	}
}

func (druid *Druid) registerCatFormSpell() {
	actionID := core.ActionID{SpellID: 768}

	srm := druid.GetSavageRoarMultiplier()

	statBonus := stats.Stats{
		stats.AttackPower:         -20, // This offset is needed because the first 10 points of Agility do not contribute any Attack Power.
		stats.PhysicalCritPercent: core.TernaryFloat64(druid.Talents.MasterShapeshifter, 4, 0),
	}

	agiApDep := druid.NewDynamicStatDependency(stats.Agility, stats.AttackPower, 2)

	var hotwDep *stats.StatDependency
	if druid.Talents.HeartOfTheWild > 0 {
		hotwDep = druid.NewDynamicMultiplyStat(stats.AttackPower, []float64{1.0, 1.03, 1.07, 1.1}[druid.Talents.HeartOfTheWild])
	}

	var leatherSpecDep *stats.StatDependency
	if druid.LeatherSpecActive {
		leatherSpecDep = druid.NewDynamicMultiplyStat(stats.Agility, 1.05)
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
			druid.PseudoStats.BaseDodgeChance += 0.02 * float64(druid.Talents.FeralSwiftness)

			druid.AddStatsDynamic(sim, statBonus)
			druid.EnableDynamicStatDep(sim, agiApDep)
			if hotwDep != nil {
				druid.EnableDynamicStatDep(sim, hotwDep)
			}
			if leatherSpecDep != nil {
				druid.EnableDynamicStatDep(sim, leatherSpecDep)
			}

			if !druid.Env.MeasuringStats {
				druid.AutoAttacks.SetReplaceMHSwing(nil)
				druid.AutoAttacks.EnableAutoSwing(sim)
				druid.UpdateManaRegenRates()

				// These buffs stay up, but corresponding changes don't
				if druid.SavageRoarAura.IsActive() {
					druid.MHAutoSpell.DamageMultiplier *= srm
				}

				if druid.PredatoryInstinctsAura != nil {
					druid.PredatoryInstinctsAura.Activate(sim)
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.form = Humanoid

			druid.AutoAttacks.SetMH(druid.WeaponFromMainHand(druid.DefaultMeleeCritMultiplier()))

			druid.PseudoStats.ThreatMultiplier /= 0.71
			druid.PseudoStats.SpiritRegenMultiplier /= AnimalSpiritRegenSuppression
			druid.PseudoStats.BaseDodgeChance -= 0.02 * float64(druid.Talents.FeralSwiftness)

			druid.AddStatsDynamic(sim, statBonus.Invert())
			druid.DisableDynamicStatDep(sim, agiApDep)
			if hotwDep != nil {
				druid.DisableDynamicStatDep(sim, hotwDep)
			}
			if leatherSpecDep != nil {
				druid.DisableDynamicStatDep(sim, leatherSpecDep)
			}

			if !druid.Env.MeasuringStats {
				druid.AutoAttacks.SetReplaceMHSwing(nil)
				druid.AutoAttacks.EnableAutoSwing(sim)
				druid.UpdateManaRegenRates()

				druid.TigersFuryAura.Deactivate(sim)

				// These buffs stay up, but corresponding changes don't
				if druid.SavageRoarAura.IsActive() {
					druid.MHAutoSpell.DamageMultiplier /= srm
				}

				if druid.PredatoryInstinctsAura != nil {
					druid.PredatoryInstinctsAura.Deactivate(sim)
				}

				if druid.StrengthOfThePantherAura.IsActive() {
					druid.StrengthOfThePantherAura.Deactivate(sim)
				}
			}
		},
	})

	if druid.Talents.FeralSwiftness > 0 {
		druid.CatFormAura.NewMovementSpeedEffect(0.15 * float64(druid.Talents.FeralSwiftness))
	}

	energyMetrics := druid.NewEnergyMetrics(actionID)

	druid.CatForm = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.05,
			Multiplier: 1 - 0.1*float64(druid.Talents.NaturalShapeshifter),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			maxShiftEnergy := float64(100*druid.Talents.Furor) / 3.0

			energyDelta := maxShiftEnergy - druid.CurrentEnergy()

			if energyDelta < 0 {
				druid.SpendEnergy(sim, -energyDelta, energyMetrics)
			}
			druid.CatFormAura.Activate(sim)
		},
	})
}

func (druid *Druid) registerBearFormSpell() {
	actionID := core.ActionID{SpellID: 5487}
	healthMetrics := druid.NewHealthMetrics(actionID)

	statBonus := stats.Stats{
		stats.AttackPower: -20, // This offset is needed because the first 10 points of Agility do not contribute any Attack Power.
	}

	agiApDep := druid.NewDynamicStatDependency(stats.Agility, stats.AttackPower, 2)
	stamDep := druid.NewDynamicMultiplyStat(stats.Stamina, 1.2)

	var hotwDep *stats.StatDependency
	if druid.Talents.HeartOfTheWild > 0 {
		hotwDep = druid.NewDynamicMultiplyStat(stats.Stamina, 1.0+0.02*float64(druid.Talents.HeartOfTheWild))
	}

	var leatherSpecDep *stats.StatDependency
	if druid.LeatherSpecActive {
		leatherSpecDep = druid.NewDynamicMultiplyStat(stats.Stamina, 1.05)
	}

	nrdtm := 1 - 0.09*float64(druid.Talents.NaturalReaction)

	clawWeapon := druid.GetBearWeapon()

	druid.BearFormAura = druid.RegisterAura(core.Aura{
		Label:      "Bear Form",
		ActionID:   actionID,
		Duration:   core.NeverExpires,
		BuildPhase: core.Ternary(druid.StartingForm.Matches(Bear), core.CharacterBuildPhaseBase, core.CharacterBuildPhaseNone),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if !druid.Env.MeasuringStats && druid.form != Humanoid {
				druid.ClearForm(sim)
			}
			druid.form = Bear
			druid.SetCurrentPowerBar(core.RageBar)

			druid.AutoAttacks.SetMH(clawWeapon)

			druid.PseudoStats.ThreatMultiplier *= 5
			druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= core.TernaryFloat64(druid.Talents.MasterShapeshifter, 1.04, 1.0)
			druid.PseudoStats.DamageTakenMultiplier *= nrdtm
			druid.PseudoStats.SpiritRegenMultiplier *= AnimalSpiritRegenSuppression
			druid.PseudoStats.BaseDodgeChance += 0.02*float64(druid.Talents.FeralSwiftness) + 0.03*float64(druid.Talents.NaturalReaction)

			druid.AddStatsDynamic(sim, statBonus)
			druid.ApplyDynamicEquipScaling(sim, stats.Armor, druid.BearArmorMultiplier())
			druid.EnableDynamicStatDep(sim, agiApDep)

			// Preserve fraction of max health when shifting
			healthFrac := druid.CurrentHealth() / druid.MaxHealth()
			druid.EnableDynamicStatDep(sim, stamDep)
			if hotwDep != nil {
				druid.EnableDynamicStatDep(sim, hotwDep)
			}
			if leatherSpecDep != nil {
				druid.EnableDynamicStatDep(sim, leatherSpecDep)
			}

			if !druid.Env.MeasuringStats {
				druid.GainHealth(sim, healthFrac*druid.MaxHealth()-druid.CurrentHealth(), healthMetrics)
				druid.AutoAttacks.EnableAutoSwing(sim)
				druid.UpdateManaRegenRates()
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.form = Humanoid
			druid.AutoAttacks.SetMH(druid.WeaponFromMainHand(druid.DefaultMeleeCritMultiplier()))

			druid.PseudoStats.ThreatMultiplier /= 5
			druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= core.TernaryFloat64(druid.Talents.MasterShapeshifter, 1.04, 1.0)
			druid.PseudoStats.DamageTakenMultiplier /= nrdtm
			druid.PseudoStats.SpiritRegenMultiplier /= AnimalSpiritRegenSuppression
			druid.PseudoStats.BaseDodgeChance -= 0.02*float64(druid.Talents.FeralSwiftness) + 0.03*float64(druid.Talents.NaturalReaction)

			druid.AddStatsDynamic(sim, statBonus.Invert())
			druid.RemoveDynamicEquipScaling(sim, stats.Armor, druid.BearArmorMultiplier())
			druid.DisableDynamicStatDep(sim, agiApDep)

			healthFrac := druid.CurrentHealth() / druid.MaxHealth()
			druid.DisableDynamicStatDep(sim, stamDep)
			if hotwDep != nil {
				druid.DisableDynamicStatDep(sim, hotwDep)
			}
			if leatherSpecDep != nil {
				druid.DisableDynamicStatDep(sim, leatherSpecDep)
			}

			if !druid.Env.MeasuringStats {
				druid.RemoveHealth(sim, druid.CurrentHealth()-healthFrac*druid.MaxHealth())
				druid.AutoAttacks.EnableAutoSwing(sim)
				druid.UpdateManaRegenRates()
				druid.EnrageAura.Deactivate(sim)

				if druid.PulverizeAura.IsActive() {
					druid.PulverizeAura.Deactivate(sim)
				}
			}
		},
	})

	rageMetrics := druid.NewRageMetrics(actionID)

	furorProcChance := float64(druid.Talents.Furor) / 3.0

	druid.BearForm = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.05,
			Multiplier: 1 - 0.1*float64(druid.Talents.NaturalShapeshifter),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rageDelta := 0 - druid.CurrentRage()
			if sim.Proc(furorProcChance, "Furor") {
				rageDelta += 10
			}
			if rageDelta > 0 {
				druid.AddRage(sim, rageDelta, rageMetrics)
			} else if rageDelta < 0 {
				druid.SpendRage(sim, -rageDelta, rageMetrics)
			}
			druid.BearFormAura.Activate(sim)
		},
	})
}

// func (druid *Druid) applyMoonkinForm() {
// 	if !druid.InForm(Moonkin) || !druid.Talents.MoonkinForm {
// 		return
// 	}

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
