package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (druid *Druid) registerRavageSpell() {
	weaponMultiplier := 9.5
	flatDamageBonus := 532.0 / weaponMultiplier
	highHpCritPercentBonus := 25.0 * float64(druid.Talents.PredatoryStrikes)

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
			Cost:   60,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Pre-pull stealth is not supported currently (and is never a DPS
			// gain anyway), so require a Stampede proc to cast in combat.
			return druid.StampedeCatAura.IsActive()
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

			druid.StampedeCatAura.Deactivate(sim)

			if sim.IsExecutePhase90() {
				spell.BonusCritPercent -= highHpCritPercentBonus
			}
		},
	})
}
