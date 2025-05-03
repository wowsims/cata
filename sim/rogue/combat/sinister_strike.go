package combat

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

func (comRogue *CombatRogue) registerSinisterStrikeSpell() {
	baseDamage := comRogue.ClassSpellScaling * 0.1780000031

	comRogue.SinisterStrike = comRogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1752},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | rogue.SpellFlagBuilder | rogue.SpellFlagColdBlooded | core.SpellFlagAPL,
		ClassSpellMask: rogue.RogueSpellSinisterStrike,

		EnergyCost: core.EnergyCostOptions{
			Cost:   50,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier:         1.04, // 84 * .73500001431 + 42
		DamageMultiplierAdditive: 1,
		CritMultiplier:           comRogue.CritMultiplier(true),
		ThreatMultiplier:         1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			comRogue.BreakStealth(sim)
			baseDamage := baseDamage +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				comRogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
