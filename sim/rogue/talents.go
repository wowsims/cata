package rogue

import (
	"math"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func (rogue *Rogue) ApplyTalents() {
	// 	rogue.applyMurder()
	// 	rogue.applySlaughterFromTheShadows()
	// 	rogue.applySealFate()
	// 	rogue.applyWeaponSpecializations()
	// 	rogue.applyCombatPotency()
	// 	rogue.applyFocusedAttacks()
	// 	rogue.applyInitiative()

	// 	rogue.AddStat(stats.Dodge, core.DodgeRatingPerDodgeChance*2*float64(rogue.Talents.LightningReflexes))
	// 	rogue.PseudoStats.MeleeSpeedMultiplier *= []float64{1, 1.03, 1.06, 1.10}[rogue.Talents.LightningReflexes]
	// 	rogue.AddStat(stats.Parry, core.ParryRatingPerParryChance*2*float64(rogue.Talents.Deflection))
	// 	rogue.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*1*float64(rogue.Talents.Malice))
	rogue.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*2*float64(rogue.Talents.Precision))
	rogue.AddStat(stats.SpellHit, core.SpellHitRatingPerHitChance*2*float64(rogue.Talents.Precision))
	// 	rogue.AddStat(stats.Expertise, core.ExpertisePerQuarterPercentReduction*5*float64(rogue.Talents.WeaponExpertise))
	// 	rogue.AddStat(stats.ArmorPenetration, core.ArmorPenPerPercentArmor*3*float64(rogue.Talents.SerratedBlades))
	// 	rogue.AutoAttacks.OHConfig().DamageMultiplier *= rogue.dwsMultiplier()

	// 	if rogue.Talents.Deadliness > 0 {
	// 		rogue.MultiplyStat(stats.AttackPower, 1.0+0.02*float64(rogue.Talents.Deadliness))
	// 	}

	// 	if rogue.Talents.SavageCombat > 0 {
	// 		rogue.MultiplyStat(stats.AttackPower, 1.0+0.02*float64(rogue.Talents.SavageCombat))
	// 	}

	// 	if rogue.Talents.SinisterCalling > 0 {
	// 		rogue.MultiplyStat(stats.Agility, 1.0+0.03*float64(rogue.Talents.SinisterCalling))
	// 	}

	// 	rogue.registerOverkill()
	// 	rogue.registerHungerForBlood()
	// 	rogue.registerColdBloodCD()
	// 	rogue.registerBladeFlurryCD()
	// 	rogue.registerAdrenalineRushCD()
	// 	rogue.registerKillingSpreeCD()
	// 	rogue.registerShadowstepCD()
	// 	rogue.registerShadowDanceCD()
	// 	rogue.registerMasterOfSubtletyCD()
	// 	rogue.registerPreparationCD()
	// 	rogue.registerPremeditation()
	// 	rogue.registerGhostlyStrikeSpell()
	// 	rogue.registerDirtyDeeds()
	// 	rogue.registerHonorAmongThieves()
}

// dwsMultiplier returns the offhand damage multiplier
func (rogue *Rogue) dwsMultiplier() float64 {
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
