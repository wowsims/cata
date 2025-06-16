package unholy

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

var FesteringStrikeActionID = core.ActionID{SpellID: 85948}

func festeringExtendHandler(aura *core.Aura) {
	aura.UpdateExpires(aura.ExpiresAt() + time.Second*6)
}

// An instant attack that deals 200% weapon damage plus 540 and increases the duration of your Blood Plague, Frost Fever, and Chains of Ice effects on the target by up to 6 sec.
func (uhdk *UnholyDeathKnight) registerFesteringStrike() {
	uhdk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       FesteringStrikeActionID.WithTag(1),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: death_knight.DeathKnightSpellFesteringStrike,

		MaxRange: core.MaxMeleeRange,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost:  1,
			FrostRuneCost:  1,
			RunicPowerGain: 20,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
		},

		DamageMultiplier: 2.0,
		CritMultiplier:   uhdk.DefaultCritMultiplier(),
		ThreatMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := uhdk.CalcScalingSpellDmg(0.43299999833) +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			spell.SpendRefundableCostAndConvertBloodOrFrostRune(sim, result.Landed())

			if result.Landed() {
				if uhdk.FrostFeverSpell.Dot(target).IsActive() {
					festeringExtendHandler(uhdk.FrostFeverSpell.Dot(target).Aura)
				}
				if uhdk.BloodPlagueSpell.Dot(target).IsActive() {
					festeringExtendHandler(uhdk.BloodPlagueSpell.Dot(target).Aura)
					for _, relatedAura := range uhdk.BloodPlagueSpell.RelatedAuraArrays {
						festeringExtendHandler(relatedAura.Get(target))
					}
				}
			}

			spell.DealDamage(sim, result)
		},
	})
}
