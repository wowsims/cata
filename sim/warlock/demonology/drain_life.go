package demonology

import "github.com/wowsims/mop/sim/core"

func (demo *DemonologyWarlock) registerDrainLife() {
	demo.RegisterDrainLife(func(_ []core.SpellResult, spell *core.Spell, sim *core.Simulation) {
		if demo.IsInMeta() {
			if demo.DemonicFury.CanSpend(30) {
				demo.DemonicFury.Spend(30, spell.ActionID, sim)
			} else {
				demo.ChanneledDot.Deactivate(sim)
			}
		} else {
			demo.DemonicFury.Gain(10, spell.ActionID, sim)
		}
	})
}
