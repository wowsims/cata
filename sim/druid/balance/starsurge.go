package balance

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/druid"
)

const (
	StarsurgeVariance   = 0.319
	StarsurgeCoeff      = 4.54
	StarsurgeBonusCoeff = 2.388
)

func (moonkin *BalanceDruid) registerStarsurgeSpell() {
	//druid.SetSpellEclipseEnergy(DruidSpellStarsurge, StarsurgeBaseEnergyGain, StarsurgeBaseEnergyGain)

	moonkin.Starsurge = moonkin.RegisterSpell(druid.Humanoid|druid.Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 78674},
		SpellSchool:    core.SpellSchoolArcane | core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: druid.DruidSpellStarsurge,
		Flags:          core.SpellFlagAPL,
		MissileSpeed:   20,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           moonkin.DefaultCritMultiplier(),
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 15.5,
			PercentModifier: 100,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
			CD: core.Cooldown{
				Timer:    moonkin.NewTimer(),
				Duration: time.Second * 15,
			},
		},

		ThreatMultiplier: 1,

		BonusCoefficient: StarsurgeBonusCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {

			baseDamage := moonkin.CalcAndRollDamageRange(sim, StarsurgeCoeff, StarsurgeVariance)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})

			if result.Landed() && result.DidCrit() {
				sunfireDot := moonkin.Sunfire.Dot(target)
				moonfireDot := moonkin.Moonfire.Dot(target)

				tryExtendDot(moonfireDot)
				tryExtendDot(sunfireDot)
			}
		},
	})
}
