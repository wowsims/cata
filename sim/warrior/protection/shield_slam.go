package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ProtectionWarrior) RegisterShieldSlam() {
	war.shieldSlam = war.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 23922},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial, // TODO: Is this right?
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: warrior.SpellMaskShieldSlam | warrior.SpellMaskSpecialAttack,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   20,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Second * 6,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return war.PseudoStats.CanBlock
		},

		DamageMultiplier: 1.0,
		CritMultiplier:   war.DefaultCritMultiplier(),
		ThreatMultiplier: 1.3,
		FlatThreatBonus:  770,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO: Reimplement using scaling coefficients and variance once those stats are available
			baseDamage := sim.Roll(1871.65, 1967.63) + spell.MeleeAttackPower()*0.6
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
