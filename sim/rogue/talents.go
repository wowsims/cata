package rogue

import (
	"math"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func (rogue *Rogue) ApplyTalents() {
	rogue.ApplyArmorSpecializationEffect()
	rogue.PseudoStats.MeleeSpeedMultiplier *= []float64{1, 1.02, 1.04, 1.06}[rogue.Talents.LightningReflexes]
	rogue.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*2*float64(rogue.Talents.Precision))
	rogue.AddStat(stats.SpellHit, core.SpellHitRatingPerHitChance*2*float64(rogue.Talents.Precision))

	if rogue.Talents.SavageCombat > 0 {
		rogue.MultiplyStat(stats.AttackPower, 1.0+0.03*float64(rogue.Talents.SavageCombat))
	}
}

// DWSMultiplier returns the offhand damage multiplier
func (rogue *Rogue) DWSMultiplier() float64 {
	// DWS (Now named Ambidexterity) is now a Combat rogue passive
	return core.TernaryFloat64(rogue.Spec == proto.Spec_SpecCombatRogue, 1.75, 1)
}

func (rogue *Rogue) makeFinishingMoveEffectApplier() func(sim *core.Simulation, numPoints int32) {
	ruthlessnessMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 14161})
	relentlessStrikesMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 14179})
	var mayhemMetrics *core.ResourceMetrics
	if rogue.HasSetBonus(Tier10, 4) {
		mayhemMetrics = rogue.NewComboPointMetrics(core.ActionID{SpellID: 70802})
	}
	return func(sim *core.Simulation, numPoints int32) {
		if rogue.Talents.Ruthlessness > 0 {
			procChance := 0.2 * float64(rogue.Talents.Ruthlessness)
			if sim.Proc(procChance, "Ruthlessness") {
				rogue.AddComboPoints(sim, 1, ruthlessnessMetrics)
			}
		}
		if rogue.Talents.RelentlessStrikes > 0 {
			procChance := []float64{0.0, 0.07, 0.14, 0.2}[rogue.Talents.RelentlessStrikes] * float64(numPoints)
			if sim.Proc(procChance, "Relentless Strikes") {
				rogue.AddEnergy(sim, 25, relentlessStrikesMetrics)
			}
		}
		if mayhemMetrics != nil {
			if sim.RandomFloat("Mayhem") < 0.13 {
				rogue.AddComboPoints(sim, 3, mayhemMetrics)
			}
		}
	}
}

func (rogue *Rogue) makeGeneratorCostModifier() func(baseCost float64) float64 {
	if rogue.HasSetBonus(Tier7, 4) {
		return func(baseCost float64) float64 {
			return math.RoundToEven(0.95 * baseCost)
		}
	}
	return func(baseCost float64) float64 {
		return baseCost
	}
}
