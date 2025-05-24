package destruction

import (
	"github.com/wowsims/mop/sim/core"
)

func (destruction DestructionWarlock) registerFelflame() {
	destruction.RegisterFelflame(func(resultList []core.SpellResult, spell *core.Spell, sim *core.Simulation) {
		destruction.BurningEmbers.Gain(core.TernaryInt32(resultList[0].DidCrit(), 2, 1), spell.ActionID, sim)
	})
}
