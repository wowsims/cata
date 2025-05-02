package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (druid *Druid) registerWrathSpell() {
	druid.SetSpellEclipseEnergy(DruidSpellWrath, WrathBaseEnergyGain, WrathBaseEnergyGain)

	druid.Wrath = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 5176},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DruidSpellWrath,
		Flags:          core.SpellFlagAPL | SpellFlagOmenTrigger,
		MissileSpeed:   20,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 9,
			PercentModifier: 100,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2500,
			},
		},

		BonusCoefficient: 0.879,

		// TODO: Was the value of 1 here incorrect to begin with?
		BonusCritPercent: 1 / core.CritRatingPerCritPercent,

		DamageMultiplier: 1,

		CritMultiplier: druid.DefaultCritMultiplier(),

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			min, max := core.CalcScalingSpellEffectVarianceMinMax(proto.Class_ClassDruid, 0.896, 0.12)
			baseDamage := sim.Roll(min, max)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				spell.WaitTravelTime(sim, func(sim *core.Simulation) {
					spell.DealDamage(sim, result)
				})
			}
		},
	})
}
