package retribution

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/paladin"
)

func (retPaladin *RetributionPaladin) RegisterTemplarsVerdict() {
	actionId := core.ActionID{SpellID: 85256}
	hpMetrics := retPaladin.NewHolyPowerMetrics(actionId)

	retPaladin.TemplarsVerdict = retPaladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,
		ClassSpellMask: paladin.SpellMaskTemplarsVerdict,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return retPaladin.GetHolyPowerValue() > 0
		},

		DamageMultiplier: 1,
		CritMultiplier:   retPaladin.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			holyPower := retPaladin.GetHolyPowerValue()

			multiplier := []float64{0, 0.3, 0.9, 2.35}[holyPower]

			spell.ApplyDamageMultiplierMultiplicative(multiplier)
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			spell.ApplyDamageMultiplierMultiplicative(1 / multiplier)

			if result.Landed() {
				retPaladin.SpendHolyPower(sim, hpMetrics)
			}

			spell.DealOutcome(sim, result)
		},
	})
}
