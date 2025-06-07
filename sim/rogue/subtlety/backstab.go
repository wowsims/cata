package subtlety

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

func (subRogue *SubtletyRogue) registerBackstabSpell() {
	baseDamage := subRogue.GetBaseDamageFromCoefficient(0.36800000072)
	weaponDamage := 3.8

	subRogue.Backstab = subRogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 53},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | rogue.SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: rogue.RogueSpellBackstab,

		EnergyCost: core.EnergyCostOptions{
			Cost:   35,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:    time.Second,
				GCDMin: time.Millisecond * 700,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !subRogue.PseudoStats.InFrontOfTarget && subRogue.HasDagger(core.MainHand)
		},

		DamageMultiplierAdditive: weaponDamage,
		DamageMultiplier:         1,
		CritMultiplier:           subRogue.CritMultiplier(true),
		ThreatMultiplier:         1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			subRogue.BreakStealth(sim)
			baseDamage := baseDamage +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				subRogue.AddComboPointsOrAnticipation(sim, 1, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
