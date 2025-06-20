package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (rogue *Rogue) registerStealthAura() {
	extraDuration := core.Ternary(rogue.Talents.Subterfuge, time.Second*3, 0)

	rogue.StealthAura = rogue.RegisterAura(core.Aura{
		Label:    "Stealth",
		ActionID: core.ActionID{SpellID: 1784},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// Stealth triggered auras
			if rogue.Spec == proto.Spec_SpecSubtletyRogue {
				rogue.MasterOfSubtletyAura.Duration = core.NeverExpires
				rogue.MasterOfSubtletyAura.Activate(sim)
			}
			if rogue.Talents.Nightstalker {
				rogue.NightstalkerMod.Activate()
			}
			if rogue.Talents.ShadowFocus {
				rogue.ShadowFocusMod.Activate()
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if rogue.Spec == proto.Spec_SpecSubtletyRogue {
				rogue.MasterOfSubtletyAura.Deactivate(sim)
				rogue.MasterOfSubtletyAura.Duration = time.Second*6 + extraDuration
				rogue.MasterOfSubtletyAura.Activate(sim)
			}
			if rogue.Talents.Subterfuge {
				rogue.SubterfugeAura.Activate(sim)
			}
			if rogue.Talents.Nightstalker {
				rogue.NightstalkerMod.Deactivate()
			}
			if rogue.Talents.ShadowFocus {
				rogue.ShadowFocusMod.Deactivate()
			}
		},
		// Stealth breaks on damage taken (if not absorbed)
		// This may be desirable later, but not applicable currently
	})
}
