package demonology

import (
	"github.com/wowsims/mop/sim/core"
)

func (demonology *DemonologyWarlock) registerHellfire() {
	hellfire := demonology.RegisterHellfire(func(resultList []core.SpellResult, spell *core.Spell, sim *core.Simulation) {
		if demonology.IsInMeta() {
			return
		}

		// 10 for primary, 3 for every other target
		fury := 10 + (len(resultList)) - 1*3
		demonology.DemonicFury.Gain(sim, int32(fury), spell.ActionID)
	})

	oldExtra := hellfire.ExtraCastCondition
	hellfire.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		if oldExtra != nil && !oldExtra(sim, target) {
			return false
		}

		return !demonology.IsInMeta()
	}

	hellfire.DamageMultiplier *= 1.25 // 2025.06.13 Changes to Beta - Hellfire and Immolation Aura increased by 25%
	demonology.Metamorphosis.RelatedSelfBuff.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
		demonology.Hellfire.SelfHot().Deactivate(sim)
	})

}
