package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ProtectionWarrior) registerShieldBarrier() {
	actionID := core.ActionID{SpellID: 112048}
	rageMetrics := war.NewRageMetrics(actionID)
	maxRageSpent := 60.0
	rageSpent := 20.0
	apScaling := 2.0 // Beta changes 2025-06-16: Shield Barrier’s attack power modifier increased to 2.0 (was 1.8). [5.2 Revert]
	staminaScaling := 2.50
	newAbsorb := 0.0

	war.ShieldBarrierAura = war.NewDamageAbsorptionAura(core.AbsorptionAuraConfig{
		Aura: core.Aura{
			Label:    "Shield Barrier",
			ActionID: actionID,
			Duration: 6 * time.Second,
		},
		ShieldStrengthCalculator: func(_ *core.Unit) float64 {
			return newAbsorb
		},
	})

	war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: warrior.SpellMaskShieldBarrier,

		RageCost: core.RageCostOptions{
			Cost: 20,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Millisecond * 1500,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   war.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			additionalRage := min(40, war.CurrentRage())
			war.SpendRage(sim, additionalRage, rageMetrics)
			rageSpent = float64(spell.Cost.BaseCost) + additionalRage

			absorbMultiplier := core.TernaryFloat64(war.T14Tank2P != nil && war.T14Tank2P.IsActive(), 1.05, 1)
			newAbsorb = (max(
				apScaling*(war.GetStat(stats.AttackPower)-war.GetStat(stats.Strength)*2),
				war.GetStat(stats.Stamina)*staminaScaling,
			) * rageSpent / maxRageSpent) * absorbMultiplier

			if !war.ShieldBarrierAura.Aura.IsActive() || (war.ShieldBarrierAura.Aura.IsActive() && newAbsorb < war.ShieldBarrierAura.ShieldStrength) {
				war.ShieldBarrierAura.Deactivate(sim)
				war.ShieldBarrierAura.Activate(sim)
			}
		},

		RelatedSelfBuff: war.ShieldBarrierAura.Aura,
	})

}
