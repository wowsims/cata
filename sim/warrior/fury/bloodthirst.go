package fury

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *FuryWarrior) registerBloodthirst() {
	actionID := core.ActionID{SpellID: 23881}
	rageMetrics := war.NewRageMetrics(actionID)
	war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: warrior.SpellMaskBloodthirst,
		MaxRange:       core.MaxMeleeRange,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Millisecond * 4500,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 0.9 * 1.2, // 2013-09-23	[Bloodthirst]'s damage has been increased by 20%.
		CritMultiplier:   war.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := war.CalcScalingSpellDmg(1) + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			bonusCritPercent := spell.Unit.GetStat(stats.PhysicalCritPercent)
			spell.BonusCritPercent += bonusCritPercent
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			spell.BonusCritPercent -= bonusCritPercent

			if result.Landed() {
				war.AddRage(sim, 10, rageMetrics)
			}
		},
	})
}
