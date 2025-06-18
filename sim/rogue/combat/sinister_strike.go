package combat

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

func (comRogue *CombatRogue) registerSinisterStrikeSpell() {
	baseDamage := comRogue.GetBaseDamageFromCoefficient(0.22499999404)
	wepDamage := 1.92

	comRogue.SinisterStrike = comRogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1752},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | rogue.SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: rogue.RogueSpellSinisterStrike,

		EnergyCost: core.EnergyCostOptions{
			Cost:   40,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:    time.Second,
				GCDMin: time.Millisecond * 500,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier:         wepDamage,
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
				comRogue.AddComboPointsOrAnticipation(sim, 1, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
