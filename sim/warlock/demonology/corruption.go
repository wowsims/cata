package demonology

import "github.com/wowsims/mop/sim/core"

func (demonology *DemonologyWarlock) registerCorruption() {
	corruption := demonology.RegisterCorruption(func(resultList []core.SpellResult, spell *core.Spell, sim *core.Simulation) {
		if resultList[0].Landed() {
			demonology.DemonicFury.Gain(sim, 4, spell.ActionID)
		}
	})

	// replaced by doom in meta
	corruption.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		return !demonology.IsInMeta()
	}
}
