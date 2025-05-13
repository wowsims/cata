package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (warrior *Warrior) RegisterSlamSpell() {

	weaponDamageConfig := SpellEffectWeaponDmgPctConfig{
		BaseWeapon_Pct:    1.1,
		Coefficient:       0.3829999864,
		EffectPerLevel:    0.97299998999,
		BaseSpellLevel:    44,
		MaxSpellLevel:     80,
		ClassSpellScaling: warrior.ClassSpellScaling,
	}

	slamActionID := core.ActionID{SpellID: 1464}

	ohSlam := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       slamActionID.WithTag(2),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		ClassSpellMask: SpellMaskSlam | SpellMaskSpecialAttack,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: weaponDamageConfig.CalcSpellDamagePct(),
		CritMultiplier:   warrior.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			ohBaseDamage := weaponDamageConfig.CalcAddedSpellDamage() + spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, ohBaseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	})

	warrior.Slam = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       slamActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskSlam | SpellMaskSpecialAttack,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   15,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {

				// Slam now has a "Haste Affects Melee Ability Casttime" flag in cata, which doesn't affect the gcd
				warrior.Slam.CurCast.CastTime = spell.Unit.ApplyCastSpeedForSpell(spell.CurCast.CastTime, spell)

				if cast.CastTime > 0 {
					warrior.AutoAttacks.DelayMeleeBy(sim, cast.CastTime)
				}
			},
		},

		DamageMultiplier: weaponDamageConfig.CalcSpellDamagePct(),
		CritMultiplier:   warrior.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		FlatThreatBonus:  140,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := weaponDamageConfig.CalcAddedSpellDamage() + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			//OH SMF hit is on a separate attack table roll, so we continue if the main hand hit didn't land.
			if !result.Landed() {
				spell.IssueRefund(sim)
			}

			// SMF adds an OH hit to slam
			if warrior.Talents.SingleMindedFury && warrior.AutoAttacks.IsDualWielding {
				ohSlam.Cast(sim, target)
			}
		},
	})
}
