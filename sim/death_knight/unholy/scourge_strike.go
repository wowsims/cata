package unholy

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

var ScourgeStrikeActionID = core.ActionID{SpellID: 55090}

/*
An unholy strike that deals 135% of weapon damage as Physical damage plus 597.
In addition, for each of your diseases on your target, you deal an additional 25% of the Physical damage done as Shadow damage.
*/
func (dk *UnholyDeathKnight) registerScourgeStrikeShadowDamage() *core.Spell {
	return dk.Unit.RegisterSpell(core.SpellConfig{
		ActionID:       ScourgeStrikeActionID.WithTag(2), // actually 70890
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamageProc,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreModifiers,
		ClassSpellMask: death_knight.DeathKnightSpellScourgeStrikeShadow,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.lastScourgeStrikeDamage * dk.GetDiseaseMulti(target, 0.0, 0.25)
			if target.HasActiveAuraWithTag(core.SpellDamageEffectAuraTag) {
				baseDamage *= 1.05
			}
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
		},
	})
}

func (dk *UnholyDeathKnight) registerScourgeStrike() {
	shadowDamageSpell := dk.registerScourgeStrikeShadowDamage()

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       ScourgeStrikeActionID.WithTag(1),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: death_knight.DeathKnightSpellScourgeStrike,

		MaxRange: core.MaxMeleeRange,

		RuneCost: core.RuneCostOptions{
			UnholyRuneCost: 1,
			RunicPowerGain: 10,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
		},

		DamageMultiplier: 1.35,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.CalcScalingSpellDmg(0.47900000215) +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			spell.SpendRefundableCost(sim, result)

			if result.Landed() && dk.DiseasesAreActive(target) {
				dk.lastScourgeStrikeDamage = result.Damage
				shadowDamageSpell.Cast(sim, target)
			}

			spell.DealDamage(sim, result)
		},
	})
}
