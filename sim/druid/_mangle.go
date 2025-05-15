package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (druid *Druid) registerMangleBearSpell() {
	mangleAuras := druid.NewEnemyAuraArray(core.MangleAura)
	glyphBonus := core.TernaryFloat64(druid.HasPrimeGlyph(proto.DruidPrimeGlyph_GlyphOfMangle), 1.1, 1.0)
	maxHits := min(druid.Env.GetNumTargets(), 3)

	druid.MangleBear = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 33878},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: DruidSpellMangleBear,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   15,
			Refund: 0.8,
		},
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

		DamageMultiplier: 1.9 * glyphBonus,
		CritMultiplier:   druid.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 1,
		MaxRange:         core.MaxMeleeRange,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			numHits := core.TernaryInt32(druid.BerserkAura.IsActive(), maxHits, 1)
			curTarget := target

			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := 3306.0/1.9 + spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
				result := spell.CalcAndDealDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

				if result.Landed() {
					mangleAuras.Get(curTarget).Activate(sim)
				} else {
					spell.IssueRefund(sim)
				}

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			if druid.BerserkAura.IsActive() {
				spell.CD.Reset()
			}

			// Preferentially consume Berserk procs over Clearcasting procs
			if druid.BerserkProcAura.IsActive() {
				druid.BerserkProcAura.Deactivate(sim)
			} else if druid.ClearcastingAura.IsActive() {
				druid.ClearcastingAura.Deactivate(sim)
			}
		},

		RelatedAuraArrays: mangleAuras.ToMap(),
	})
}

func (druid *Druid) registerMangleCatSpell() {
	mangleAuras := druid.NewEnemyAuraArray(core.MangleAura)
	glyphBonus := core.TernaryFloat64(druid.HasPrimeGlyph(proto.DruidPrimeGlyph_GlyphOfMangle), 1.1, 1.0)
	hasBloodletting := druid.HasPrimeGlyph(proto.DruidPrimeGlyph_GlyphOfBloodletting)

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

		DamageMultiplier: 5.4 * glyphBonus,
		CritMultiplier:   druid.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 1,
		MaxRange:         core.MaxMeleeRange,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 302.0/5.4 +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				druid.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				mangleAuras.Get(target).Activate(sim)

				// Mangle (Cat) can also extend Rip in Cata
				if hasBloodletting {
					druid.ApplyBloodletting(target)
				}

				// 4pT11
				if druid.T11Feral4pBonus.IsActive() {
					aura := druid.StrengthOfThePantherAura

					if aura.IsActive() {
						aura.Refresh(sim)
						aura.AddStack(sim)
					} else {
						aura.Activate(sim)
						aura.SetStacks(sim, 1)
					}
				}
			} else {
				spell.IssueRefund(sim)
			}
		},

		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := 302.0/5.4 + spell.Unit.AutoAttacks.MH().CalculateAverageWeaponDamage(spell.MeleeAttackPower())
			return spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMeleeWeaponSpecialHitAndCrit)
		},

		RelatedAuraArrays: mangleAuras.ToMap(),
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
