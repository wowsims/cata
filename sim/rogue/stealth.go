package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (rogue *Rogue) registerStealthAura() {
	rogue.StealthAura = rogue.RegisterAura(core.Aura{
		Label:    "Stealth",
		ActionID: core.ActionID{SpellID: 1784},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// Stealth triggered auras
			if rogue.Talents.Overkill {
				rogue.OverkillAura.Duration = core.NeverExpires
				rogue.OverkillAura.Activate(sim)
			}
			if rogue.Spec == proto.Spec_SpecSubtletyRogue {
				rogue.MasterOfSubtletyAura.Duration = core.NeverExpires
				rogue.MasterOfSubtletyAura.Activate(sim)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if rogue.Talents.Overkill {
				rogue.OverkillAura.Deactivate(sim)
				rogue.OverkillAura.Duration = time.Second * 20
				rogue.OverkillAura.Activate(sim)
			}
			if rogue.Spec == proto.Spec_SpecSubtletyRogue {
				rogue.MasterOfSubtletyAura.Deactivate(sim)
				rogue.MasterOfSubtletyAura.Duration = time.Second * 6
				rogue.MasterOfSubtletyAura.Activate(sim)
			}
		},
		// Stealth breaks on damage taken (if not absorbed)
		// This may be desirable later, but not applicable currently
	})
}
