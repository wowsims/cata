package shadow

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/priest"
)

// impact spell
const dpImpactScale = 1.566
const dpImpactCoeff = 0.786

// dot spell
const DpDotScale = 0.261
const DpDotCoeff = 0.131

func (shadow *ShadowPriest) registerDevouringPlagueSpell() {
	actionID := core.ActionID{SpellID: 2944, Tag: 0}
	shadow.DevouringPlague = shadow.RegisterSpell(core.SpellConfig{
		ActionID:                 actionID,
		SpellSchool:              core.SpellSchoolShadow,
		ProcMask:                 core.ProcMaskSpellDamage,
		Flags:                    core.SpellFlagDisease | core.SpellFlagAPL,
		ClassSpellMask:           priest.PriestSpellDevouringPlague,
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,
		CritMultiplier:           shadow.DefaultCritMultiplier(),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		BonusCoefficient: dpImpactCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			shadow.orbsConsumed = shadow.ShadowOrbs.Value()
			spell.DamageMultiplier *= float64(shadow.orbsConsumed)
			result := spell.CalcDamage(sim, target, shadow.CalcScalingSpellDmg(dpImpactScale), spell.OutcomeMagicHitAndCrit)
			spell.DamageMultiplier /= float64(shadow.orbsConsumed)
			if result.Landed() {
				shadow.ShadowOrbs.Spend(sim, shadow.orbsConsumed, actionID)
				spell.RelatedDotSpell.Cast(sim, target)
			}

			spell.DealOutcome(sim, result)
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {

			// At least 1 shadow orb needs to be present
			return shadow.ShadowOrbs.CanSpend(1)
		},
	})

	shadow.DevouringPlague.RelatedDotSpell = shadow.RegisterSpell(core.SpellConfig{
		ActionID:                 actionID.WithTag(1),
		SpellSchool:              core.SpellSchoolShadow,
		ProcMask:                 core.ProcMaskSpellDamage,
		Flags:                    core.SpellFlagDisease | core.SpellFlagPassiveSpell,
		ClassSpellMask:           priest.PriestSpellDevouringPlagueDoT,
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,
		CritMultiplier:           shadow.DefaultCritMultiplier(),
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Devouring Plague",
			},
			NumberOfTicks:       6,
			TickLength:          time.Second,
			AffectedByCastSpeed: true,
			BonusCoefficient:    DpDotCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Spell.DamageMultiplier *= float64(shadow.orbsConsumed)
				dot.Snapshot(target, shadow.CalcScalingSpellDmg(DpDotScale))
				dot.Spell.DamageMultiplier /= float64(shadow.orbsConsumed)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})
}
