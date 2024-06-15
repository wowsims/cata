package rogue

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func (rogue *Rogue) ApplyTalents() {
	rogue.ApplyArmorSpecializationEffect(stats.Agility, proto.ArmorType_ArmorTypeLeather)
	rogue.PseudoStats.MeleeSpeedMultiplier *= []float64{1, 1.02, 1.04, 1.06}[rogue.Talents.LightningReflexes]
	rogue.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*2*float64(rogue.Talents.Precision))
	rogue.AddStat(stats.SpellHit, core.SpellHitRatingPerHitChance*2*float64(rogue.Talents.Precision))

	if rogue.Talents.SavageCombat > 0 {
		rogue.MultiplyStat(stats.AttackPower, 1.0+0.03*float64(rogue.Talents.SavageCombat))
	}

	if rogue.Talents.Ruthlessness > 0 {
		rogue.ruthlessnessMetrics = rogue.NewComboPointMetrics(core.ActionID{SpellID: 14161})
	}

	if rogue.Talents.RelentlessStrikes > 0 {
		rogue.relentlessStrikesMetrics = rogue.NewEnergyMetrics(core.ActionID{SpellID: 14179})
	}
}

// DWSMultiplier returns the offhand damage multiplier
func (rogue *Rogue) DWSMultiplier() float64 {
	// DWS (Now named Ambidexterity) is now a Combat rogue passive
	return core.TernaryFloat64(rogue.Spec == proto.Spec_SpecCombatRogue, 1.75, 1)
}
