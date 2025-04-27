package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

var FesteringStrikeActionID = core.ActionID{SpellID: 85948}

func festeringExtendHandler(aura *core.Aura) {
	aura.UpdateExpires(aura.ExpiresAt() + time.Second*6)
}

func (dk *DeathKnight) registerFesteringStrikeSpell() {
	ohSpell := dk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       FesteringStrikeActionID.WithTag(2),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: DeathKnightSpellFesteringStrike,

		DamageMultiplier: 1.5,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.24899999797 +
				spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	hasReaping := dk.Inputs.Spec == proto.Spec_SpecUnholyDeathKnight

	dk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       FesteringStrikeActionID.WithTag(1),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellFesteringStrike,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost:  1,
			FrostRuneCost:  1,
			RunicPowerGain: 20,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1.5,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.49799999595 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if hasReaping {
				spell.SpendRefundableCostAndConvertBloodOrFrostRune(sim, result, 1)
			} else {
				spell.SpendRefundableCost(sim, result)
			}
			dk.ThreatOfThassarianProc(sim, result, ohSpell)

			if result.Landed() {
				if dk.FrostFeverSpell.Dot(target).IsActive() {
					festeringExtendHandler(dk.FrostFeverSpell.Dot(target).Aura)
				}
				if dk.BloodPlagueSpell.Dot(target).IsActive() {
					festeringExtendHandler(dk.BloodPlagueSpell.Dot(target).Aura)
				}
				if dk.Talents.EbonPlaguebringer > 0 && dk.EbonPlagueAura.Get(target).IsActive() {
					festeringExtendHandler(dk.EbonPlagueAura.Get(target))
				}
			}

			spell.DealDamage(sim, result)
		},
	})
}

func (dk *DeathKnight) registerDrwFesteringStrikeSpell() *core.Spell {
	return dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    FesteringStrikeActionID.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.49799999595 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				if dk.RuneWeapon.FrostFeverSpell.Dot(target).IsActive() {
					festeringExtendHandler(dk.RuneWeapon.FrostFeverSpell.Dot(target).Aura)
				}
				if dk.RuneWeapon.BloodPlagueSpell.Dot(target).IsActive() {
					festeringExtendHandler(dk.RuneWeapon.BloodPlagueSpell.Dot(target).Aura)
				}
			}

			spell.DealDamage(sim, result)
		},
	})
}
