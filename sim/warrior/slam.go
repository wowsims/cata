package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warrior *Warrior) RegisterSlamSpell() {

	spellScallingConfig := core.SpellScalingConfig{
		BaseWeapon_Pct:    1.1,
		Coefficient:       0.3829999864,
		EffectPerLevel:    0.97299998999,
		BaseSpellLevel:    44,
		MaxSpellLevel:     80,
		ClassSpellScaling: warrior.ClassSpellScaling,
	}

	warrior.Slam = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1464},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskSlam | SpellMaskSpecialAttack,

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

		DamageMultiplier: spellScallingConfig.CalcSpellDamagePct(),
		CritMultiplier:   warrior.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,
		FlatThreatBonus:  140,

		BonusCoefficient: 1,

		// TODO: check if the OH SMF hit is on a separate attack table roll
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spellScallingConfig.CalcAddedSpellDamage() + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			if !result.Landed() {
				spell.IssueRefund(sim)
			}

			// SMF adds an OH hit to slam
			if warrior.Talents.SingleMindedFury {
				baseDamage := spellScallingConfig.CalcAddedSpellDamage() + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			}
		},
	})
}
