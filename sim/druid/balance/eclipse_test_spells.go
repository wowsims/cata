package balance

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/druid"
)

func (balance *BalanceDruid) RegisterTestSpells() {

	testWrathlMetric := balance.NewLunarEnergyMetrics(core.ActionID{SpellID: 9739})
	balance.RegisterSpell(druid.Moonkin|druid.Humanoid, core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 9739},
		SpellSchool:  core.SpellSchoolNature,
		Flags:        core.SpellFlagAPL,
		ProcMask:     core.ProcMaskSpellDamage,
		MissileSpeed: 40,
		ManaCost: core.ManaCostOptions{
			BaseCost: 0.09,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
		},

		BonusCoefficient: 0.4,
		DamageMultiplier: 1,
		CritMultiplier:   balance.DefaultSpellCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, 200, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				balance.AddEclipseEnergy(13+1.0/3.0, core.LunarEnergy, sim, testWrathlMetric)
				spell.WaitTravelTime(sim, func(s *core.Simulation) {
					spell.DealDamage(sim, result)
				})
			}
		},
	})

	testSFMetric := balance.NewSolarEnergyMetric(core.ActionID{SpellID: 21668})
	balance.RegisterSpell(druid.Moonkin|druid.Humanoid, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 21668},
		SpellSchool: core.SpellSchoolNature,
		Flags:       core.SpellFlagAPL,
		ProcMask:    core.ProcMaskSpellDamage,
		ManaCost: core.ManaCostOptions{
			BaseCost: 0.09,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
		},

		BonusCoefficient: 0.4,
		DamageMultiplier: 1,
		CritMultiplier:   balance.DefaultSpellCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealDamage(sim, target, 200, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				balance.AddEclipseEnergy(13+1.0/3.0, core.SolarEnergy, sim, testSFMetric)
			}
		},
	})

	// Eclipse (Lunar) 48518
	// Eclipse (Solar) 48517

	lunarEclipse := balance.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 48518},
		Label:    "Eclipse (Lunar)",
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// Do stuff
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			// Do stuff
		},
	})

	solarEclipse := balance.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 48517},
		Label:    "Eclipse (Solar)",
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// Do stuff
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			// Do stuff
		},
	})

	balance.AddEclipseCallback(func(eclipse core.Eclipse, gained bool, sim *core.Simulation) {
		if eclipse == core.LunarEclipse {
			if gained {
				lunarEclipse.Activate(sim)
			} else {
				lunarEclipse.Deactivate(sim)
			}
		} else {
			if gained {
				solarEclipse.Activate(sim)
			} else {
				solarEclipse.Deactivate(sim)
			}
		}
	})
}
