package druid

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (druid *Druid) registerPulverizeSpell() {
	if !druid.Talents.Pulverize {
		return
	}

	statBonusPerStack := stats.Stats{stats.MeleeCrit: 3.0 * core.CritRatingPerCritChance}

	druid.PulverizeAura = druid.RegisterAura(core.Aura{
		Label:     "Pulverize",
		ActionID:  core.ActionID{SpellID: 80951},
		MaxStacks: 3,
		Duration:  core.DurationFromSeconds(10.0 + 4.0*float64(druid.Talents.EndlessCarnage)),

		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			druid.AddStatsDynamic(sim, statBonusPerStack.Multiply(float64(newStacks-oldStacks)))
		},
	})

	druid.Pulverize = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 80313},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   15,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 0.6,
		CritMultiplier:   druid.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			lacerateDot := druid.Lacerate.Dot(target)
			lacerateStacksConsumed := core.TernaryInt32(lacerateDot.IsActive(), lacerateDot.GetStacks(), 0)
			flatDamage := 1623.6 * float64(lacerateStacksConsumed)
			baseDamage := flatDamage/0.6 + spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				lacerateDot.Cancel(sim)
				druid.PulverizeAura.Activate(sim)
				druid.PulverizeAura.SetStacks(sim, lacerateStacksConsumed)
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
