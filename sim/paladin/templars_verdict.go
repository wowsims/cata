package paladin

import (
	"github.com/wowsims/mop/sim/core"
)

// A powerful weapon strike that consumes 3 charges of Holy Power to deal 275% weapon damage plus 628.
func (paladin *Paladin) registerTemplarsVerdict() {
	actionID := core.ActionID{SpellID: 85256}

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
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return paladin.HolyPower.CanSpend(3)
		},

		DamageMultiplier: 2.75,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if paladin.T15Ret4pc.IsActive() {
				paladin.T15Ret4pcTemplarsVerdict.Cast(sim, target)
				spell.SpellMetrics[target.UnitIndex].Casts--
				return
			}

			baseDamage := paladin.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) + paladin.CalcScalingSpellDmg(0.55000001192)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				paladin.HolyPower.Spend(sim, 3, actionID)
			}

			spell.DealDamage(sim, result)
		},
	})
}
