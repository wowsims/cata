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
		demonology.DemonicFury.Gain(int32(fury), spell.ActionID, sim)
	})

	oldExtra := hellfire.ExtraCastCondition
	hellfire.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		if oldExtra != nil && !oldExtra(sim, target) {
			return false
		}

		return !demonology.IsInMeta()
	}

	demonology.Metamorphosis.RelatedSelfBuff.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
		demonology.Hellfire.SelfHot().Deactivate(sim)
	})

}
