package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (rogue *Rogue) registerAmbushSpell() {
	baseDamage := rogue.ClassSpellScaling * 0.32699999213

	rogue.Ambush = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 8676},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagBuilder | SpellFlagColdBlooded | core.SpellFlagAPL,
		ClassSpellMask: RogueSpellAmbush,

		EnergyCost: core.EnergyCostOptions{
			Cost:   60,
			Refund: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !rogue.PseudoStats.InFrontOfTarget && (rogue.IsStealthed() || rogue.HasActiveAura("Shadowmeld"))
		},

		BonusCritPercent: 20 * float64(rogue.Talents.ImprovedAmbush),
		DamageMultiplier: core.TernaryFloat64(rogue.HasDagger(core.MainHand), 2.86, 1.97), // 77 * 1.38999998569 + 90 (*1.45 for Dagger)
		// Imp Ambush also Additive
		DamageMultiplierAdditive: 1 +
			0.05*float64(rogue.Talents.ImprovedAmbush) +
			0.1*float64(rogue.Talents.Opportunity),
		CritMultiplier:   rogue.CritMultiplier(false),
		ThreatMultiplier: 1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := baseDamage +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 2, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})

	rogue.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand}, func(s *core.Simulation, slot proto.ItemSlot) {
		// Recalculate Ambush's multiplier in case the MH weapon changed.
		rogue.Ambush.DamageMultiplier = core.TernaryFloat64(rogue.HasDagger(core.MainHand), 2.86, 1.97)
	})
}
