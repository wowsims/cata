package fury

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *FuryWarrior) registerRagingBlow() {
	ragingBlowAura := war.RegisterAura(core.Aura{
		Label:     "Raging Blow!",
		ActionID:  core.ActionID{SpellID: 131116},
		Duration:  12 * time.Second,
		MaxStacks: 2,
		Icd: &core.Cooldown{
			Timer:    war.NewTimer(),
			Duration: time.Millisecond * 500,
		},
	})

	war.EnrageAura.ApplyOnGain(func(_ *core.Aura, sim *core.Simulation) {
		ragingBlowAura.Activate(sim)
		ragingBlowAura.AddStack(sim)
	})

	ragingBlowActionID := core.ActionID{SpellID: 85288}

	ohRagingBlow := war.RegisterSpell(core.SpellConfig{
		ActionID:       ragingBlowActionID.WithTag(3),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		ClassSpellMask: warrior.SpellMaskRagingBlowOH,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1.9 * 1.2, // 2013-09-23	[Bloodthirst]'s damage has been increased by 20%.
		CritMultiplier:   war.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			ohBaseDamage := spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			spell.CalcAndDealDamage(sim, target, ohBaseDamage, spell.OutcomeMeleeSpecialBlockAndCrit)
		},
	})

	mhRagingBlow := war.RegisterSpell(core.SpellConfig{
		ActionID:       ragingBlowActionID.WithTag(2),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: warrior.SpellMaskRagingBlowMH,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1.9 * 1.2, // 2013-09-23	[Bloodthirst]'s damage has been increased by 20%.
		CritMultiplier:   war.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			mhBaseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			spell.CalcAndDealDamage(sim, target, mhBaseDamage, spell.OutcomeMeleeSpecialBlockAndCrit)
		},
	})

	war.RegisterSpell(core.SpellConfig{
		ActionID:       ragingBlowActionID.WithTag(1),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: warrior.SpellMaskRagingBlow,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   10,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1.0,
		CritMultiplier:   war.DefaultCritMultiplier(),

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return war.HasMHWeapon() && war.HasOHWeapon() && ragingBlowAura.IsActive() && ragingBlowAura.GetStacks() > 0
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			//One roll for both mh and oh and for all meat cleaver targets.
			//Raging blow aura consumed always
			//Meat cleaver aura consumed on hit and consumed even if hitting only one target.

			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHitNoHitCounter)
			ragingBlowAura.RemoveStack(sim)

			if !result.Landed() {
				spell.IssueRefund(sim)
				return
			}

			meatCleaverStacks := int(war.MeatCleaverAura.GetStacks())
			if war.MeatCleaverAura.IsActive() && meatCleaverStacks > 0 {
				for index, mcTarget := range sim.Encounter.TargetUnits {
					if index <= meatCleaverStacks {
						mhRagingBlow.Cast(sim, mcTarget)
						ohRagingBlow.Cast(sim, mcTarget)
					}
				}
				war.MeatCleaverAura.Deactivate(sim)
			} else {
				mhRagingBlow.Cast(sim, result.Target)
				ohRagingBlow.Cast(sim, result.Target)
			}

		},
	})
}
