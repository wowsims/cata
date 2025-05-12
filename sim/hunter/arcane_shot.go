package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) registerArcaneShotSpell() {
	hunter.ArcaneShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 3044},
		SpellSchool:    core.SpellSchoolArcane,
		ClassSpellMask: HunterSpellArcaneShot,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MissileSpeed:   40,
		MinRange:       5,
		MaxRange:       40,
		FocusCost: core.FocusCostOptions{
			Cost: 25,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		BonusCritPercent:         0,
		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		CritMultiplier:           hunter.DefaultCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			wepDmg := hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower())
			baseDamage := wepDmg + (0.0483 * spell.RangedAttackPower()) + 290

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
