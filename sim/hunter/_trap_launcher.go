package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (hunter *Hunter) registerTrapLauncher() {
	actionID := core.ActionID{SpellID: 77769}

	hunter.TrapLauncherAura = hunter.RegisterAura(core.Aura{
		Label:    "Trap Launcher",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == hunter.ExplosiveTrap {
				aura.Deactivate(sim)
			}
		},
	})

	hunter.TrapLauncher = hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		FocusCost: core.FocusCostOptions{
			Cost: 20 - core.TernaryInt32(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfTrapLauncher), 10, 0),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Millisecond * 1500,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.TrapLauncherAura.Activate(sim)
		},
	})
}
