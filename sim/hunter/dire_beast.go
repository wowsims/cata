package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) RegisterDireBeastSpell() {
	if !hunter.Talents.DireBeast {
		return
	}
	actionID := core.ActionID{SpellID: 120679}
	direBeastSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 30,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.DireBeastPet.EnableWithTimeout(sim, hunter.DireBeastPet, time.Second*15)
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: direBeastSpell,
		Type:  core.CooldownTypeDPS,
	})
}
