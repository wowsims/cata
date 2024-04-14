package death_knight

import (
	"github.com/wowsims/cata/sim/core"
)

var ScourgeStrikeActionID = core.ActionID{SpellID: 55090}

// this is just a simple spell because it has no rune costs and is really just a wrapper.
func (dk *DeathKnight) registerScourgeStrikeShadowDamageSpell() *core.Spell {
	return dk.Unit.RegisterSpell(core.SpellConfig{
		ActionID:       ScourgeStrikeActionID.WithTag(2),
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellOrProc,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: DeathKnightSpellScourgeStrikeShadow,

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.LastScourgeStrikeDamage * dk.dkCountActiveDiseases(target) * 0.18
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
		},
	})
}

func (dk *DeathKnight) RegisterScourgeStrikeSpell() {
	shadowDamageSpell := dk.registerScourgeStrikeShadowDamageSpell()

	dk.ScourgeStrike = dk.RegisterSpell(core.SpellConfig{
		ActionID:       ScourgeStrikeActionID.WithTag(1),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellScourgeStrike,

		RuneCost: core.RuneCostOptions{
			UnholyRuneCost: 1,
			RunicPowerGain: 10,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,

		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassBaseScaling*0.55500000715 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			spell.SpendRefundableCost(sim, result)

			if result.Landed() && dk.DiseasesAreActive(target) {
				dk.LastScourgeStrikeDamage = result.Damage
				shadowDamageSpell.Cast(sim, target)
			}

			spell.DealDamage(sim, result)
		},
	})
}
