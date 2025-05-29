package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (druid *Druid) registerMangleBearSpell() {
	maxHits := min(druid.Env.GetNumTargets(), 3)
	actionID := core.ActionID{SpellID: 33878}
	rageMetrics := druid.NewRageMetrics(actionID)

	druid.MangleBear = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: DruidSpellMangleBear,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 2.8,
		CritMultiplier:   druid.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 1,
		MaxRange:         core.MaxMeleeRange,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			numHits := core.TernaryInt32(druid.BerserkBearAura.IsActive(), maxHits, 1)
			curTarget := target

			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
				result := spell.CalcAndDealDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

				if result.Landed() {
					druid.AddRage(sim, 5, rageMetrics)
				}

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			if druid.BerserkBearAura.IsActive() {
				spell.CD.Reset()
			}
		},
	})
}

func (druid *Druid) registerMangleCatSpell() {
	flatBaseDamage := 0.07100000232 * druid.ClassSpellScaling // ~77.7265

	druid.MangleCat = druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 33876},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: DruidSpellMangleCat,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   35.0,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 5,
		CritMultiplier:   druid.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 1,
		MaxRange:         core.MaxMeleeRange,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := flatBaseDamage +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				druid.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				druid.ApplyBloodletting(target)
			} else {
				spell.IssueRefund(sim)
			}
		},

		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := flatBaseDamage + spell.Unit.AutoAttacks.MH().CalculateAverageWeaponDamage(spell.MeleeAttackPower())
			return spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMeleeWeaponSpecialHitAndCrit)
		},
	})
}

func (druid *Druid) CurrentMangleCatCost() float64 {
	return druid.MangleCat.Cost.GetCurrentCost()
}

func (druid *Druid) IsMangle(spell *core.Spell) bool {
	if druid.MangleBear != nil && druid.MangleBear.IsEqual(spell) {
		return true
	} else if druid.MangleCat != nil && druid.MangleCat.IsEqual(spell) {
		return true
	}
	return false
}
