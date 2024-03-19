package marksmanship

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (hunter *MarksmanshipHunter) registerAimedShotSpell(timer *core.Timer) {

	hunter.AimedShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 19434},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		FocusCost: core.FocusCostOptions{
			Cost: 50 - (float64(hunter.Talents.Efficiency) * 2),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      time.Second,
				CastTime: time.Second * 3,
			},
			IgnoreHaste: true,
			ModifyCast: func(_ *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.CastTime = spell.CastTime()
			},

			CastTime: func(spell *core.Spell) time.Duration {
				return time.Duration(float64(spell.DefaultCast.CastTime) / hunter.RangedSwingSpeed())
			},
		},
		BonusCritRating: 0 +
			core.TernaryFloat64(hunter.Talents.TrueshotAura, 10*core.CritRatingPerCritChance, 0),
		DamageMultiplierAdditive: 1,
		DamageMultiplier: 1,
		CritMultiplier:   hunter.CritMultiplier(true, true, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			wepDmg := hunter.AutoAttacks.Ranged().CalculateWeaponDamage(sim, spell.RangedAttackPower(target))
			rap := spell.RangedAttackPower(target) * 0.724 + 766
			baseDamage := ((wepDmg + rap) * 1.6) + 100
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
		},
	})
}
