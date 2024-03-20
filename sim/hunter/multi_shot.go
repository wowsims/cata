package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (hunter *Hunter) registerMultiShotSpell() {
	numHits := hunter.Env.GetNumTargets() // Multi is uncapped in Cata

	hunter.MultiShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 2643},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		FocusCost: core.FocusCostOptions{
			Cost: 40,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		BonusCritRating:          0,
		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1.2,
		CritMultiplier:           hunter.CritMultiplier(true, false, false),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			sharedDmg := hunter.AutoAttacks.Ranged().BaseDamage(sim) +
				spell.BonusWeaponDamage() //

			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := sharedDmg + 0.2*spell.RangedAttackPower(curTarget)
				spell.CalcAndDealDamage(sim, curTarget, baseDamage, spell.OutcomeRangedHitAndCrit)

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}
		},
	})
}
