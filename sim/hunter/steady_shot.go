package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (hunter *Hunter) registerSteadyShotSpell() {

	ssMetrics := hunter.NewFocusMetrics(core.ActionID{SpellID: 56641})

	hunter.SteadyShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 56641},
		SpellSchool: core.SpellSchoolPhysical,
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
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
			ModifyCast: func(_ *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.CastTime = spell.CastTime()
			},

			CastTime: func(spell *core.Spell) time.Duration {
				return time.Duration(float64(spell.DefaultCast.CastTime) / hunter.RangedSwingSpeed())
			},
		},

		// BonusCritRating: 0 +
		// 	2*core.CritRatingPerCritChance*float64(hunter.Talents.SurvivalInstincts),
		DamageMultiplierAdditive: 1 + core.TernaryFloat64(hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfSteadyShot), 0.1, 0),
		DamageMultiplier: 0.62,
		// 	hunter.markedForDeathMultiplier(),
		CritMultiplier:1,//   hunter.critMultiplier(true, true, false), // what is this
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if hunter.Talents.ImprovedSteadyShot > 0 {
				if !hunter.ImprovedSteadyShotAuraCounter.IsActive() {
					hunter.ImprovedSteadyShotAuraCounter.Activate(sim)
				}
			}
			baseDamage := hunter.AutoAttacks.Ranged().CalculateWeaponDamage(sim, spell.RangedAttackPower(target)) + (280 + spell.RangedAttackPower(target) * 0.021)
			hunter.AddFocus(sim, 9, ssMetrics)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
			spell.DealDamage(sim, result)
		},
	})
}
