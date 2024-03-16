package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (hunter *Hunter) registerKillShotSpell() {
	hunter.KillShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53351},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second*15,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return sim.IsExecutePhase20()
		},

		BonusCritRating: 0 +
			5*core.CritRatingPerCritChance*float64(hunter.Talents.SniperTraining),
		DamageMultiplier: 1, //
		CritMultiplier: 1,//  hunter.critMultiplier(true, true, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// 0.2 rap from normalized weapon (2.8/14) and 0.2 from bonus ratio
			baseDamage := 0.4*spell.RangedAttackPower(target) +
				hunter.AutoAttacks.Ranged().BaseDamage(sim) +
				spell.BonusWeaponDamage() +
				325
			baseDamage *= 2
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
		},
	})
}
