package paladin

import (
	"github.com/wowsims/mop/sim/core"
)

func (paladin *Paladin) registerTemplarsVerdict() {
	actionID := core.ActionID{SpellID: 85256}
	bonusDamage := paladin.CalcScalingSpellDmg(0.55000001192)

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskTemplarsVerdict,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return paladin.HolyPower.CanSpend(3)
		},

		DamageMultiplier: 2.75,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := paladin.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) + bonusDamage

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				paladin.HolyPower.Spend(3, actionID, sim)
			}

			spell.DealOutcome(sim, result)
		},
	})
}
