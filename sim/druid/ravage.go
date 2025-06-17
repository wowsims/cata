package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (druid *Druid) registerRavageSpell() {
	const weaponMultiplier = 9.5
	const highHpCritPercentBonus = 50.0
	flatDamageBonus := 0.07100000232 * druid.ClassSpellScaling

	druid.Ravage = druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 6785},
		SpellSchool:      core.SpellSchoolPhysical,
		ProcMask:         core.ProcMaskMeleeMHSpecial,
		ClassSpellMask:   DruidSpellRavage,
		Flags:            core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		DamageMultiplier: weaponMultiplier,
		CritMultiplier:   druid.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 1,
		MaxRange:         core.MaxMeleeRange,

		EnergyCost: core.EnergyCostOptions{
			Cost:   45,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},

			IgnoreHaste: true,
		},

		ExtraCastCondition: func(_ *core.Simulation, _ *core.Unit) bool {
			return druid.ProwlAura.IsActive() && !druid.PseudoStats.InFrontOfTarget && !druid.CannotShredTarget
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if sim.IsExecutePhase90() {
				spell.BonusCritPercent += highHpCritPercentBonus
			}

			baseDamage := flatDamageBonus + spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				druid.AddComboPoints(sim, 1, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}

			if sim.IsExecutePhase90() {
				spell.BonusCritPercent -= highHpCritPercentBonus
			}
		},
	})
}
