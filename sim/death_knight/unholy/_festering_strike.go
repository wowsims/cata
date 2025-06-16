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

// An instant attack that deals 200% weapon damage plus 540 and increases the duration of your Blood Plague, Frost Fever, and Chains of Ice effects on the target by up to 6 sec.
func (dk *DeathKnight) registerFesteringStrike() {
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

			spell.SpendRefundableCostAndConvertBloodOrFrostRune(sim, result.Landed())

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
