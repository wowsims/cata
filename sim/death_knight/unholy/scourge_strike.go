package unholy

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/death_knight"
)

var scourgeStrikeActionID = core.ActionID{SpellID: 55090}

// this is just a simple spell because it has no rune costs and is really just a wrapper.
func (dk *UnholyDeathKnight) registerScourgeStrikeShadowDamageSpell() *core.Spell {
	return dk.Unit.RegisterSpell(core.SpellConfig{
		ActionID:       scourgeStrikeActionID.WithTag(2),
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellOrProc,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: death_knight.DeathKnightSpellScourgeStrikeShadow,

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.lastScourgeStrikeDamage * dk.GetDiseaseMulti(target, 0.0, 0.18)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
		},
	})
}

func (dk *UnholyDeathKnight) registerScourgeStrikeSpell() {
	shadowDamageSpell := dk.registerScourgeStrikeShadowDamageSpell()

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       scourgeStrikeActionID.WithTag(1),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: death_knight.DeathKnightSpellScourgeStrike,

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
			baseDamage := dk.ClassSpellScaling*0.55500000715 +
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
