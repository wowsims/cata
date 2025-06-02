package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

const cdDuration = time.Millisecond * 1500

func (war *Warrior) registerHeroicStrikeSpell() {
	getHSDamageMultiplier := func() float64 {
		has1H := war.MainHand().HandType != proto.HandType_HandTypeTwoHand
		return core.TernaryFloat64(has1H, 1.4, 1)
	}

	weaponDamageMod := war.AddDynamicMod(core.SpellModConfig{
		ClassMask:  SpellMaskHeroicStrike,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: getHSDamageMultiplier(),
	})

	war.RegisterItemSwapCallback(core.AllWeaponSlots(), func(_ *core.Simulation, _ proto.ItemSlot) {
		weaponDamageMod.UpdateFloatValue(getHSDamageMultiplier())
	})

	war.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 78},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskHeroicStrike,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   30,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    war.sharedHSCleaveCD,
				Duration: cdDuration,
			},
		},

		DamageMultiplier: 1.1,
		ThreatMultiplier: 1,
		CritMultiplier:   war.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := war.CalcScalingSpellDmg(0.40000000596)*getHSDamageMultiplier() + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}

func (war *Warrior) registerCleaveSpell() {
	maxTargets := 2
	results := make([]*core.SpellResult, maxTargets)

	getCleaveDamageMultiplier := func() float64 {
		has1H := war.MainHand().HandType != proto.HandType_HandTypeTwoHand
		return core.TernaryFloat64(has1H, 1.15, 0.8)
	}

	weaponDamageMod := war.AddDynamicMod(core.SpellModConfig{
		ClassMask:  SpellMaskCleave,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: getCleaveDamageMultiplier(),
	})

	war.RegisterItemSwapCallback(core.AllWeaponSlots(), func(_ *core.Simulation, _ proto.ItemSlot) {
		weaponDamageMod.UpdateFloatValue(getCleaveDamageMultiplier())
	})

	war.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 845},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskCleave,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost: 30,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    war.sharedHSCleaveCD,
				Duration: cdDuration,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   war.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			for idx, target := range sim.Encounter.TargetUnits {
				if idx > maxTargets {
					break
				}
				results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	})

}
