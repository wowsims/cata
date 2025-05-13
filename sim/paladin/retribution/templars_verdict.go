package retribution

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/paladin"
)

func (retPaladin *RetributionPaladin) RegisterTemplarsVerdict() {
	actionId := core.ActionID{SpellID: 85256}

	retPaladin.TemplarsVerdict = retPaladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
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
			return retPaladin.HolyPower.CanSpend(1)
		},

		DamageMultiplier: 1,
		CritMultiplier:   retPaladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			holyPower := int32(retPaladin.HolyPower.Value())

			multiplier := []float64{0, 0.3, 0.9, 2.35}[holyPower]

			spell.DamageMultiplier *= multiplier
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			spell.DamageMultiplier /= multiplier

			if result.Landed() {
				retPaladin.HolyPower.SpendUpTo(3, actionId, sim)
			}

			spell.DealOutcome(sim, result)
		},
	})
}
