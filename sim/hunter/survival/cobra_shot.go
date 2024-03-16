package survival

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (hunter *SurvivalHunter) registerCobraShotSpell() {
	
	csMetrics := hunter.NewFocusMetrics(core.ActionID{SpellID: 77767})

	hunter.Hunter.CobraShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 77767},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,
		FocusCost: core.FocusCostOptions{
			
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      time.Second,
				CastTime: time.Second * 2,
			},
			IgnoreHaste: true,
			CastTime: func(spell *core.Spell) time.Duration {
				return time.Duration(float64(spell.DefaultCast.CastTime) / hunter.RangedSwingSpeed())
			},
		},

		// BonusCritRating: 0 +
		// 	2*core.CritRatingPerCritChance*float64(hunter.Talents.SurvivalInstincts),
		DamageMultiplierAdditive: 1,
		// DamageMultiplier: 1 *
		// 	hunter.markedForDeathMultiplier(),
		CritMultiplier:1,//   hunter.critMultiplier(true, true, false), // what is this
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 
				hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower(target)) * 0.246 + 277.21
			hunter.AddFocus(sim, 9, csMetrics)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
			spell.DealDamage(sim, result)
		},
	})
}