package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (rogue *Rogue) registerAmbushSpell() {
	baseDamage := rogue.GetBaseDamageFromCoefficient(0.5)
	weaponDamage := 3.25
	daggerModifier := 1.447

	rogue.Ambush = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 8676},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: RogueSpellAmbush,

		EnergyCost: core.EnergyCostOptions{
			Cost:   60,
			Refund: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:    time.Second,
				GCDMin: time.Millisecond * 700,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !rogue.PseudoStats.InFrontOfTarget && (rogue.IsStealthed() || rogue.HasActiveAura("Shadowmeld") || rogue.HasActiveAura("Sleight of Hand"))
		},

		DamageMultiplier:         core.TernaryFloat64(rogue.HasDagger(core.MainHand), weaponDamage*daggerModifier, weaponDamage),
		DamageMultiplierAdditive: 1,
		CritMultiplier:           rogue.CritMultiplier(false),
		ThreatMultiplier:         1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := baseDamage +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)

			if result.Landed() {
				rogue.AddComboPointsOrAnticipation(sim, 2, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})

	rogue.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand}, func(s *core.Simulation, slot proto.ItemSlot) {
		// Recalculate Ambush's multiplier in case the MH weapon changed.
		rogue.Ambush.DamageMultiplier = core.TernaryFloat64(rogue.HasDagger(core.MainHand), weaponDamage*daggerModifier, weaponDamage)
	})
}
