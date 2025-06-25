package frost

import "github.com/wowsims/mop/sim/core"

/*
Permanently transforms your Blood Runes into Death Runes.
Death Runes count as a Blood, Frost, or Unholy Rune.
*/
func (fdk *FrostDeathKnight) registerBloodOfTheNorth() {
	core.MakePermanent(fdk.GetOrRegisterAura(core.Aura{
		Label:    "Blood of the North" + fdk.Label,
		ActionID: core.ActionID{SpellID: 54637},

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			fdk.SetPermanentDeathRunes([]int8{0, 1})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			fdk.SetPermanentDeathRunes([]int8{})
		},
	}))
}
