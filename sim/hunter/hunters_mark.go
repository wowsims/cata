package hunter

import (
	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) registerHuntersMarkSpell() {
	enemyHuntersMarks := hunter.NewEnemyAuraArray(core.HuntersMarkAura)

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
			for _, aura := range enemyHuntersMarks {
				if aura.IsActive() {
					aura.Deactivate(sim)
				}
			}
			// Activating Hunters Mark for the new target
			enemyHuntersMarks.Get(target).Activate(sim)
		},
	})
}
