package retribution

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/paladin"
)

// A powerful weapon strike that consumes 3 charges of Holy Power to deal 275% weapon damage plus 628.
func (ret *RetributionPaladin) registerTemplarsVerdict() {
	actionID := core.ActionID{SpellID: 85256}

	ret.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: paladin.SpellMaskTemplarsVerdict,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return ret.HolyPower.CanSpend(3)
		},

		DamageMultiplier: 2.75,
		CritMultiplier:   ret.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if ret.T15Ret4pc.IsActive() {
				ret.T15Ret4pcTemplarsVerdict.Cast(sim, target)
				spell.SpellMetrics[target.UnitIndex].Casts--
				return
			}

			baseDamage := ret.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) + ret.CalcScalingSpellDmg(0.55000001192)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				ret.HolyPower.Spend(sim, 3, actionID)
			}

			spell.DealDamage(sim, result)
		},
	})
}
