package rogue

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func (rogue *Rogue) ApplyTalents() {
	rogue.ApplyArmorSpecializationEffect(stats.Agility, proto.ArmorType_ArmorTypeLeather, 87504)
}

// DWSMultiplier returns the offhand damage multiplier
func (rogue *Rogue) DWSMultiplier() float64 {
	// DWS (Now named Ambidexterity) is now a Combat rogue passive
	return core.TernaryFloat64(rogue.Spec == proto.Spec_SpecCombatRogue, 1.75, 1)
}

func (rogue *Rogue) ApplyCutToTheChase(sim *core.Simulation) {
	if rogue.Spec == proto.Spec_SpecAssassinationRogue && rogue.SliceAndDiceAura.IsActive() {
		rogue.SliceAndDiceAura.Duration = rogue.sliceAndDiceDurations[5]
		rogue.SliceAndDiceAura.Activate(sim)
	}
}
