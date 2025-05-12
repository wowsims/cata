package fury

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *FuryWarrior) RegisterBloodthirst() {
	war.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 23881},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: warrior.SpellMaskBloodthirst | warrior.SpellMaskSpecialAttack,
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
				Duration: time.Second * 3,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   war.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			const effectBasePoints = 43
			const baseLevel = 10
			const maxLevel = 80
			const effectRealPointsPerLevel = 0.52899998426
			multiplier := effectRealPointsPerLevel*(maxLevel-baseLevel) + effectBasePoints
			baseDamage := math.Floor(multiplier) / 100 * spell.MeleeAttackPower() //multiplier is floored, tooltip is incorrect
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
