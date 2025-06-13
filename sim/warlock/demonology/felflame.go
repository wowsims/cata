package demonology

import "github.com/wowsims/mop/sim/core"

func (demo *DemonologyWarlock) registerFelFlame() {
	felFlame := demo.RegisterFelflame(func(_ []core.SpellResult, spell *core.Spell, sim *core.Simulation) {
		demo.DemonicFury.Gain(sim, 15, spell.ActionID)
	})

	// Is replaced within meta, can not use it when active
	felFlame.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		return !demo.Metamorphosis.RelatedSelfBuff.IsActive()
	}
}
