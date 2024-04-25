package rogue

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

const RuptureEnergyCost = 25.0
const RuptureSpellID = 1943

func (rogue *Rogue) registerRupture() {
	glyphTicks := core.TernaryInt32(rogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfRupture), 2, 0)

	rogue.Rupture = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: RuptureSpellID},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | SpellFlagFinisher | core.SpellFlagAPL,
		MetricSplits: 6,

		EnergyCost: core.EnergyCostOptions{
			Cost:          RuptureEnergyCost,
			Refund:        0.8,
			RefundMetrics: rogue.EnergyRefundMetrics,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(spell.Unit.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0
		},

		DamageMultiplier: 1,
		CritMultiplier:   rogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Rupture",
				Tag:   RogueBleedTag,
			},
			NumberOfTicks: 0, // Set dynamically
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotPhysical(target, rogue.ruptureDamage(rogue.ComboPoints()))
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
				dot.NumberOfTicks = 3 + rogue.ComboPoints() + glyphTicks
				dot.Apply(sim)
				// DealOutcome must come before ApplyFinisher, or it breaks calculations based on spent combo points.
				spell.DealOutcome(sim, result)
				rogue.ApplyFinisher(sim, spell)
			} else {
				spell.DealOutcome(sim, result)
				spell.IssueRefund(sim)
			}

		},
	})
}

func (rogue *Rogue) ruptureDamage(comboPoints int32) float64 {
	return 142 +
		20*float64(comboPoints) +
		[]float64{0, 0.06 / 4, 0.12 / 5, 0.18 / 6, 0.24 / 7, 0.30 / 8}[comboPoints]*rogue.Rupture.MeleeAttackPower()
}
