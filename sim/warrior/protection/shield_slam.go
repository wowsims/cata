package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ProtectionWarrior) registerShieldSlam() {
	actionID := core.ActionID{SpellID: 23922}
	rageMetrics := war.NewRageMetrics(actionID)

	war.ShieldSlam = war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: warrior.SpellMaskShieldSlam,
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

		DamageMultiplier: 1,
		CritMultiplier:   war.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := war.CalcAndRollDamageRange(sim, 11.25, 0.05000000075) + spell.MeleeAttackPower()*1.5
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			additionalRage := core.TernaryFloat64(war.SwordAndBoardAura.IsActive(), 5, 0)

			if result.Landed() {
				war.AddRage(sim, (20+additionalRage)*war.GetRageMultiplier(target), rageMetrics)
			}
		},
	})
}
