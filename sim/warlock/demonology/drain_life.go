package demonology

import "github.com/wowsims/mop/sim/core"

func (demo *DemonologyWarlock) registerDrainLife() {
	demo.RegisterDrainLife(func(_ []core.SpellResult, spell *core.Spell, sim *core.Simulation) {
		demo.DemonicFury.Gain(10, spell.ActionID, sim)
	})
}
