package assassination

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

func (sinRogue *AssassinationRogue) registerDispatch() *core.Spell {
	addedDamage := 655
	weaponDamage := 645.0

	return sinRogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 111240},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | rogue.SpellFlagBuilder | rogue.SpellFlagColdBlooded,
		ClassSpellMask: rogue.RogueSpellDispatch,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           sinRogue.CritMultiplier(true),
		ThreatMultiplier:         1,

		EnergyCost: core.EnergyCostOptions{
			Cost:   30,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return target.CurrentHealthPercent() <= 35.0 || false // Blindside Aura Check
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := weaponDamage*spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) + float64(addedDamage)
			outcome := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMeleeSpecialHitAndCrit)
			if outcome.Landed() {
				sinRogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
