package hunter

import (
	"github.com/wowsims/cata/sim/core"
)

func (hunter *Hunter) registerHuntersMarkSpell() {
	hunter.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 1130},
		ProcMask: core.ProcMaskEmpty,
		Flags:    core.SpellFlagAPL,

		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			hm := core.HuntersMarkAura(target)
			hm.Activate(sim)
		},
	})
}
