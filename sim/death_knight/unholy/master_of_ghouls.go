package unholy

import "github.com/wowsims/mop/sim/core"

/*
The Ghoul summoned by your Raise Dead spell is considered a pet under your control.
Unlike normal Death Knight Ghouls, your pet does not have a limited duration.
Also reduces the cooldown of Raise Dead by 60 sec.
*/
func (uhdk *UnholyDeathKnight) registerMasterOfGhouls() {
	// This aura does nothing since we don't even register Raise Dead for Unholy
	core.MakePermanent(uhdk.RegisterAura(core.Aura{
		Label:    "Master of Ghouls" + uhdk.Label,
		ActionID: core.ActionID{SpellID: 52143},
	}))
}
