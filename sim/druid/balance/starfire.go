package balance

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/druid"
)

const (
	StarfireBonusCoeff = 2.166
	StarfireCoeff      = 4.456
	StarfireVariance   = 0.25
)

func (moonkin *BalanceDruid) registerStarfireSpell() {
	moonkin.Starfire = moonkin.RegisterSpell(druid.Humanoid|druid.Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 2912},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: druid.DruidSpellStarfire,
		Flags:          core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 15.5,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2700,
			},
		},

		BonusCoefficient: StarfireBonusCoeff,

		DamageMultiplier: 1,

		CritMultiplier: moonkin.DefaultCritMultiplier(),

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := moonkin.CalcAndRollDamageRange(sim, StarfireCoeff, StarfireVariance)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
