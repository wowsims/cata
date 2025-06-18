package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

const RuptureEnergyCost = 25.0
const RuptureSpellID = 1943

func (rogue *Rogue) registerRupture() {
	coefficient := 0.18500000238
	resourceCoefficient := 0.02600000054

	baseDamage := rogue.GetBaseDamageFromCoefficient(coefficient)
	damagePerComboPoint := rogue.GetBaseDamageFromCoefficient(resourceCoefficient)

	rogue.Rupture = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: RuptureSpellID},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagFinisher | core.SpellFlagAPL,
		MetricSplits:   6,
		ClassSpellMask: RogueSpellRupture,

		EnergyCost: core.EnergyCostOptions{
			Cost:          RuptureEnergyCost,
			Refund:        0.8,
			RefundMetrics: rogue.EnergyRefundMetrics,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:    time.Second,
				GCDMin: time.Millisecond * 500,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(rogue.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0
		},

		DamageMultiplier: 1,
		CritMultiplier:   rogue.CritMultiplier(false),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Rupture",
				Tag:   RogueBleedTag,
			},
			NumberOfTicks: 0, // Set dynamically
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotPhysical(target, rogue.ruptureDamage(rogue.ComboPoints(), baseDamage, damagePerComboPoint))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				dot := spell.Dot(target)
				dot.BaseTickCount = 2 + (2 * rogue.ComboPoints())
				if rogue.Has2PT15 {
					dot.BaseTickCount += 2
				}
				dot.Apply(sim)
				rogue.ApplyFinisher(sim, spell)
				spell.DealOutcome(sim, result)
			} else {
				spell.DealOutcome(sim, result)
				spell.IssueRefund(sim)
			}

		},
	})
}

func (rogue *Rogue) ruptureDamage(comboPoints int32, baseDamage float64, damagePerComboPoint float64) float64 {
	return baseDamage +
		damagePerComboPoint*float64(comboPoints) +
		[]float64{0, 0.025, 0.04, 0.05, 0.056, 0.062}[comboPoints]*rogue.Rupture.MeleeAttackPower()
}
