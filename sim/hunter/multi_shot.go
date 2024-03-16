package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (hunter *Hunter) registerMultiShotSpell(timer *core.Timer) {
	numHits := hunter.Env.GetNumTargets() // Multi is uncapped in Cata

	hunter.MultiShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49048},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		FocusCost: core.FocusCostOptions{
			Cost: 40,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      time.Second,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
		},

		BonusCritRating: 0, //+4*core.CritRatingPerCritChance*float64(hunter.Talents.ImprovedBarrage),
		DamageMultiplierAdditive: 1,//.04*float64(hunter.Talents.Barrage),
		DamageMultiplier: 1,
		CritMultiplier:   hunter.CritMultiplier(true, false, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			sharedDmg := hunter.AutoAttacks.Ranged().BaseDamage(sim) +
				spell.BonusWeaponDamage() +
				408

			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := sharedDmg + 0.2*spell.RangedAttackPower(curTarget)
				spell.CalcAndDealDamage(sim, curTarget, baseDamage, spell.OutcomeRangedHitAndCrit)

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}
		},
	})
}
