package fury

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *FuryWarrior) registerWildStrike() {
	actionID := core.ActionID{SpellID: 100130}
	war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: warrior.SpellMaskWildStrike,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   30,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 2.3,
		CritMultiplier:   war.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return war.HasOHWeapon()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if war.BloodsurgeAura.IsActive() && war.BloodsurgeAura.GetStacks() > 0 {
				war.BloodsurgeAura.RemoveStack(sim)
			}

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
