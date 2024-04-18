package fury

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/warrior"
)

func (war *FuryWarrior) RegisterRagingBlow() {

	actionID := core.ActionID{SpellID: 85288}
	rbOffhand := war.RegisterSpell(core.SpellConfig{
		ActionID:    actionID.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1.0,
		CritMultiplier:   war.DefaultMeleeCritMultiplier(),
	})

	war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskRagingBlow | warrior.SpellMaskSpecialAttack,

		RageCost: core.RageCostOptions{
			Cost:   20,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},

			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: 6 * time.Second,
			},
		},

		DamageMultiplier: 1.0,
		CritMultiplier:   war.DefaultMeleeCritMultiplier(),

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return war.StanceMatches(warrior.BerserkerStance) && war.HasActiveAuraWithTag(warrior.EnrageTag) && war.AutoAttacks.OH() != nil
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
				return
			}

			// TODO: check if these two hits should use the same attack rolls and proc masks
			// Mirroring Stormstrike for now
			mhBaseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, mhBaseDamage*war.EnrageEffectMultiplier, spell.OutcomeMeleeSpecialCritOnly)

			ohBaseDamage := spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			rbOffhand.CalcAndDealDamage(sim, target, ohBaseDamage*war.EnrageEffectMultiplier, spell.OutcomeMeleeSpecialCritOnly)
		},
	})
}
