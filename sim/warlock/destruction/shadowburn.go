package destruction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

var shadowBurnScale = 3.5
var shadowBurnVariance = 0.2
var shadowBurnCoeff = 3.5

func (destruction *DestructionWarlock) registerShadowBurnSpell() {
	manaMetric := destruction.NewManaMetrics(core.ActionID{SpellID: 17877})
	destruction.Shadowburn = destruction.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 17877},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellShadowBurn,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return sim.IsExecutePhase20() && destruction.BurningEmbers.CanSpend(core.TernaryInt32(destruction.T15_2pc.IsActive(), 8, 10))
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           destruction.DefaultCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         shadowBurnCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := destruction.CalcAndRollDamageRange(sim, shadowBurnScale, shadowBurnVariance)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				destruction.BurningEmbers.Spend(sim, core.TernaryInt32(destruction.T15_2pc.IsActive(), 8, 10), spell.ActionID)
			}

			pa := sim.GetConsumedPendingActionFromPool()
			pa.NextActionAt = sim.CurrentTime + time.Second*5
			pa.Priority = core.ActionPriorityAuto

			pa.OnAction = func(sim *core.Simulation) {
				destruction.AddMana(sim, destruction.MaxMana()*0.15, manaMetric)
			}

			sim.AddPendingAction(pa)
		},
	})
}
