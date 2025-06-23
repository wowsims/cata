package unholy

import "github.com/wowsims/mop/sim/core"

/*
Whenever you hit with Blood Strike, Pestilence, Festering Strike, Icy Touch, or Blood Boil, the Runes spent will become Death Runes when they activate.
Death Runes count as a Blood, Frost or Unholy Rune.
(100ms cooldown)
*/
func (uhdk *UnholyDeathKnight) registerReaping() {
	core.MakePermanent(uhdk.RegisterAura(core.Aura{
		Label:    "Reaping" + uhdk.Label,
		ActionID: core.ActionID{SpellID: 56835},
	}))
}
