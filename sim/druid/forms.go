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

// Thick Hide contribution handled separately in talents code for cleanliness
// and UI stats display.
const BaseBearArmorMulti = 2.2

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

	statBonus := stats.Stats{
		stats.AttackPower: -20, // This offset is needed because the first 10 points of Agility do not contribute any Attack Power.
	}

	agiApDep := druid.NewDynamicStatDependency(stats.Agility, stats.AttackPower, 2)

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

			if !druid.Env.MeasuringStats {
				druid.AutoAttacks.SetMH(clawWeapon)
				druid.AutoAttacks.EnableAutoSwing(sim)
				druid.UpdateManaRegenRates()
				druid.MHAutoSpell.DamageMultiplier *= 2
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.form = Humanoid

			druid.PseudoStats.ThreatMultiplier /= 0.71
			druid.PseudoStats.SpiritRegenMultiplier /= AnimalSpiritRegenSuppression

			druid.AddStatsDynamic(sim, statBonus.Invert())
			druid.DisableBuildPhaseStatDep(sim, agiApDep)

			if !druid.Env.MeasuringStats {
				druid.AutoAttacks.SetMH(druid.WeaponFromMainHand(druid.DefaultCritMultiplier()))
				druid.AutoAttacks.EnableAutoSwing(sim)
				druid.UpdateManaRegenRates()
				druid.MHAutoSpell.DamageMultiplier /= 2
			}
		},
	})

	druid.CatFormAura.NewMovementSpeedEffect(0.25)
}

func (druid *Druid) registerCatFormSpell() {
	druid.CatForm = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID: core.ActionID{SpellID: 768},
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 3.7,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
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
	stamDep := druid.NewDynamicMultiplyStat(stats.Stamina, 1.4)
	critDep := druid.NewDynamicMultiplyStat(stats.CritRating, 1.5)
	hasteDep := druid.NewDynamicMultiplyStat(stats.HasteRating, 1.5)

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

			druid.PseudoStats.ThreatMultiplier *= 7
			druid.PseudoStats.SpiritRegenMultiplier *= AnimalSpiritRegenSuppression

			druid.AddStatsDynamic(sim, statBonus)
			druid.ApplyDynamicEquipScaling(sim, stats.Armor, BaseBearArmorMulti)
			druid.EnableBuildPhaseStatDep(sim, agiApDep)
			druid.EnableBuildPhaseStatDep(sim, critDep)
			druid.EnableBuildPhaseStatDep(sim, hasteDep)

			// Preserve fraction of max health when shifting
			healthFrac := druid.CurrentHealth() / druid.MaxHealth()
			druid.EnableBuildPhaseStatDep(sim, stamDep)

			if druid.GuardianLeatherSpecTracker.IsActive() {
				druid.EnableBuildPhaseStatDep(sim, druid.GuardianLeatherSpecDep)
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

			druid.PseudoStats.ThreatMultiplier /= 7
			druid.PseudoStats.SpiritRegenMultiplier /= AnimalSpiritRegenSuppression

			druid.AddStatsDynamic(sim, statBonus.Invert())
			druid.RemoveDynamicEquipScaling(sim, stats.Armor, BaseBearArmorMulti)
			druid.DisableBuildPhaseStatDep(sim, agiApDep)
			druid.DisableBuildPhaseStatDep(sim, critDep)
			druid.DisableBuildPhaseStatDep(sim, hasteDep)

			healthFrac := druid.CurrentHealth() / druid.MaxHealth()
			druid.DisableBuildPhaseStatDep(sim, stamDep)

			if druid.GuardianLeatherSpecTracker.IsActive() {
				druid.DisableBuildPhaseStatDep(sim, druid.GuardianLeatherSpecDep)
			}

			if !druid.Env.MeasuringStats {
				druid.RemoveHealth(sim, druid.CurrentHealth()-healthFrac*druid.MaxHealth())
				druid.AutoAttacks.SetMH(druid.WeaponFromMainHand(druid.DefaultCritMultiplier()))
				druid.AutoAttacks.EnableAutoSwing(sim)
				druid.UpdateManaRegenRates()
			}
		},
	})
}

func (druid *Druid) registerBearFormSpell() {
	actionID := core.ActionID{SpellID: 5487}
	rageMetrics := druid.NewRageMetrics(actionID)

	druid.BearForm = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 3.7,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			rageDelta := 10.0 - druid.CurrentRage()
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
