package blood

import (
	"github.com/wowsims/mop/sim/core"
)

// Your Icebound Fortitude reduces damage taken by an additional 30%.
func (bdk *BloodDeathKnight) registerSanguineFortitude() {
	core.MakePermanent(bdk.RegisterAura(core.Aura{
		Label:    "Sanguine Fortitude" + bdk.Label,
		ActionID: core.ActionID{SpellID: 81127},
	}))
}
