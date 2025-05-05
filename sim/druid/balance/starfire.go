package balance

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/druid"
)

const variance = 0.25
const coeff = 4.456
const bonusCoeff = 2.166

func (moonkin *BalanceDruid) registerStarfireSpell() {
	//druid.SetSpellEclipseEnergy(moonkin.Starfire, moonkin.StarfireBaseEnergyGain, moonkin.StarfireBaseEnergyGain)

	moonkin.Starfire = moonkin.RegisterSpell(druid.Humanoid|druid.Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 2912},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: druid.DruidSpellStarfire,
		Flags:          core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 15.5,
			PercentModifier: 100,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2700,
			},
		},

		BonusCoefficient: bonusCoeff,

		DamageMultiplier: 1,

		CritMultiplier: moonkin.DefaultCritMultiplier(),

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			min, max := core.CalcScalingSpellEffectVarianceMinMax(proto.Class_ClassDruid, 1.383, 0.22)
			baseDamage := sim.Roll(min, max)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				// starfire crits extend moonfire by one tick
				// if result.DidCrit() {
				// 	moonfireDot := druid.Moonfire.Dot(target)
				// 	tryExtendDot(moonfireDot)
				// }

				spell.DealDamage(sim, result)
			}
		},
	})
}
