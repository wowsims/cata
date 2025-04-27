package druid

import (
	//"math"

	"github.com/wowsims/mop/sim/core"
	//"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
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
		CritMultiplier:       druid.DefaultCritMultiplier(),
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
		CritMultiplier:       druid.DefaultCritMultiplier(),
		AttackPowerPerDPS:    core.DefaultAttackPowerPerDPS,
		MaxRange:             core.MaxMeleeRange,
	}
}

func (druid *Druid) RegisterCatFormAura() {
	actionID := core.ActionID{SpellID: 768}

	// TODO: Fix this to work with the new talent system.
	// srm := druid.GetSavageRoarMultiplier()
	srm := 1.8

	statBonus := stats.Stats{
		stats.AttackPower: -20, // This offset is needed because the first 10 points of Agility do not contribute any Attack Power.
	}

	agiApDep := druid.NewDynamicStatDependency(stats.Agility, stats.AttackPower, 2)
	leatherSpecDep := druid.NewDynamicMultiplyStat(stats.Agility, 1.05)

	// Need redundant enabling/disabling of the dep both here and below
	// because we don't know whether the leather spec tracker or Cat Form will
	// activate first.
	druid.LeatherSpec.ApplyOnGain(func(_ *core.Aura, sim *core.Simulation) {
		if druid.InForm(Cat) {
			druid.EnableBuildPhaseStatDep(sim, leatherSpecDep)
		}
	})

	druid.LeatherSpec.ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
		if druid.InForm(Cat) {
			druid.DisableBuildPhaseStatDep(sim, leatherSpecDep)
		}
	})

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

			druid.PseudoStats.ThreatMultiplier *= 0.71
			druid.PseudoStats.SpiritRegenMultiplier *= AnimalSpiritRegenSuppression

			druid.AddStatsDynamic(sim, statBonus)
			druid.EnableBuildPhaseStatDep(sim, agiApDep)
			if druid.HotWCatDep != nil {
				druid.EnableBuildPhaseStatDep(sim, druid.HotWCatDep)
			}
			if druid.LeatherSpec.IsActive() {
				druid.EnableBuildPhaseStatDep(sim, leatherSpecDep)
			}

			if !druid.Env.MeasuringStats {
				druid.AutoAttacks.SetMH(clawWeapon)
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

			druid.PseudoStats.ThreatMultiplier /= 0.71
			druid.PseudoStats.SpiritRegenMultiplier /= AnimalSpiritRegenSuppression

			druid.AddStatsDynamic(sim, statBonus.Invert())
			druid.DisableBuildPhaseStatDep(sim, agiApDep)
			if druid.HotWCatDep != nil {
				druid.DisableBuildPhaseStatDep(sim, druid.HotWCatDep)
			}
			if druid.LeatherSpec.IsActive() {
				druid.DisableBuildPhaseStatDep(sim, leatherSpecDep)
			}

			if !druid.Env.MeasuringStats {
				druid.AutoAttacks.SetMH(druid.WeaponFromMainHand(druid.DefaultCritMultiplier()))
				druid.AutoAttacks.EnableAutoSwing(sim)
				druid.UpdateManaRegenRates()

				// druid.TigersFuryAura.Deactivate(sim)

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

	// if druid.Talents.FeralSwiftness > 0 {
	// 	druid.CatFormAura.NewMovementSpeedEffect(0.15 * float64(druid.Talents.FeralSwiftness))
	// }
}

func (druid *Druid) registerCatFormSpell() {
	actionID := core.ActionID{SpellID: 768}
	energyMetrics := druid.NewEnergyMetrics(actionID)

	druid.CatForm = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 5,
			// TODO: Fix this to work with the new talent system.
			// PercentModifier: 100 - (10 * druid.Talents.NaturalShapeshifter),
			PercentModifier: 100,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			// TODO: Fix this to work with the new talent system.
			// maxShiftEnergy := float64(100*druid.Talents.Furor) / 3.0
			maxShiftEnergy := 100 / 3.0

			energyDelta := maxShiftEnergy - druid.CurrentEnergy()

			if energyDelta < 0 {
				druid.SpendEnergy(sim, -energyDelta, energyMetrics)
			}
			druid.CatFormAura.Activate(sim)
		},
	})
}

func (druid *Druid) RegisterBearFormAura() {
	actionID := core.ActionID{SpellID: 5487}
	healthMetrics := druid.NewHealthMetrics(actionID)

	statBonus := stats.Stats{
		stats.AttackPower: -20, // This offset is needed because the first 10 points of Agility do not contribute any Attack Power.
	}

	agiApDep := druid.NewDynamicStatDependency(stats.Agility, stats.AttackPower, 2)
	stamDep := druid.NewDynamicMultiplyStat(stats.Stamina, 1.2)
	leatherSpecDep := druid.NewDynamicMultiplyStat(stats.Stamina, 1.05)

	// Need redundant enabling/disabling of the dep both here and below
	// because we don't know whether the leather spec tracker or Bear Form
	// will activate first.
	druid.LeatherSpec.ApplyOnGain(func(_ *core.Aura, sim *core.Simulation) {
		if druid.InForm(Bear) {
			druid.EnableBuildPhaseStatDep(sim, leatherSpecDep)
		}
	})

	druid.LeatherSpec.ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
		if druid.InForm(Bear) {
			druid.DisableBuildPhaseStatDep(sim, leatherSpecDep)
		}
	})

	clawWeapon := druid.GetBearWeapon()
	baseBearArmorMulti := 2.2 // Thick Hide contribution handled separately in talents code for cleanliness and UI stats display.

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

			druid.PseudoStats.ThreatMultiplier *= 5
			druid.PseudoStats.SpiritRegenMultiplier *= AnimalSpiritRegenSuppression

			druid.AddStatsDynamic(sim, statBonus)
			druid.ApplyDynamicEquipScaling(sim, stats.Armor, baseBearArmorMulti)
			druid.EnableBuildPhaseStatDep(sim, agiApDep)

			// Preserve fraction of max health when shifting
			healthFrac := druid.CurrentHealth() / druid.MaxHealth()
			druid.EnableBuildPhaseStatDep(sim, stamDep)
			if druid.HotWBearDep != nil {
				druid.EnableBuildPhaseStatDep(sim, druid.HotWBearDep)
			}
			if druid.LeatherSpec.IsActive() {
				druid.EnableBuildPhaseStatDep(sim, leatherSpecDep)
			}

			if !druid.Env.MeasuringStats {
				druid.GainHealth(sim, healthFrac*druid.MaxHealth()-druid.CurrentHealth(), healthMetrics)
				druid.AutoAttacks.SetMH(clawWeapon)
				druid.AutoAttacks.EnableAutoSwing(sim)
				druid.UpdateManaRegenRates()
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.form = Humanoid

			druid.PseudoStats.ThreatMultiplier /= 5
			druid.PseudoStats.SpiritRegenMultiplier /= AnimalSpiritRegenSuppression

			druid.AddStatsDynamic(sim, statBonus.Invert())
			druid.RemoveDynamicEquipScaling(sim, stats.Armor, baseBearArmorMulti)
			druid.DisableBuildPhaseStatDep(sim, agiApDep)

			healthFrac := druid.CurrentHealth() / druid.MaxHealth()
			druid.DisableBuildPhaseStatDep(sim, stamDep)
			if druid.HotWBearDep != nil {
				druid.DisableBuildPhaseStatDep(sim, druid.HotWBearDep)
			}
			if druid.LeatherSpec.IsActive() {
				druid.DisableBuildPhaseStatDep(sim, leatherSpecDep)
			}

			if !druid.Env.MeasuringStats {
				druid.RemoveHealth(sim, druid.CurrentHealth()-healthFrac*druid.MaxHealth())
				druid.AutoAttacks.SetMH(druid.WeaponFromMainHand(druid.DefaultCritMultiplier()))
				druid.AutoAttacks.EnableAutoSwing(sim)
				druid.UpdateManaRegenRates()
				druid.EnrageAura.Deactivate(sim)

				if druid.PulverizeAura.IsActive() {
					druid.PulverizeAura.Deactivate(sim)
				}
			}
		},
	})
}

func (druid *Druid) registerBearFormSpell() {
	actionID := core.ActionID{SpellID: 5487}
	rageMetrics := druid.NewRageMetrics(actionID)

	// TODO: Fix this to work with the new talent system.
	// furorProcChance := float64(druid.Talents.Furor) / 3.0
	furorProcChance := 0 / 3.0

	druid.BearForm = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 5,
			// TODO: Fix this to work with the new talent system.
			// PercentModifier: 100 - (10 * druid.Talents.NaturalShapeshifter),
			PercentModifier: 100,
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
