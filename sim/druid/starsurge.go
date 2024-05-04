package druid

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (druid *Druid) registerStarsurgeSpell() {
	solarMetric := druid.NewSolarEnergyMetric(core.ActionID{SpellID: 78674})
	lunarMetric := druid.NewLunarEnergyMetrics(core.ActionID{SpellID: 78674})

	druid.Starsurge = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 78674},
		SpellSchool:    core.SpellSchoolArcane | core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DruidSpellStarsurge,
		Flags:          core.SpellFlagAPL,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           druid.DefaultSpellCritMultiplier(),
		ManaCost: core.ManaCostOptions{
			BaseCost:   0.11,
			Multiplier: 1,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 15,
			},
		},

		ThreatMultiplier: 1,

		BonusCoefficient: 1.228,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			min, max := core.CalcScalingSpellEffectVarianceMinMax(proto.Class_ClassDruid, 1.228, 0.32)
			baseDamage := sim.Roll(min, max)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				if druid.CanGainEnergy(SolarEnergy) {
					druid.AddEclipseEnergy(15, SolarEnergy, sim, solarMetric)
				} else {
					druid.AddEclipseEnergy(15, LunarEnergy, sim, lunarMetric)
				}

				spell.WaitTravelTime(sim, func(sim *core.Simulation) {
					spell.DealDamage(sim, result)
				})
			}
		},
	})
}
