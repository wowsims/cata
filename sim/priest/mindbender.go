package priest

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (priest *Priest) registerMindbenderSpell() {
	if !priest.Talents.Mindbender {
		return
	}

	actionID := core.ActionID{SpellID: 123040}

	// For timeline only
	priest.MindbenderAura = priest.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Mindbender",
		Duration: time.Second * 15.0,
	})

	priest.MindBender = priest.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: PriestSpellMindBender,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			priest.MindbenderPet.EnableWithTimeout(sim, priest.MindbenderPet, time.Second*15.0)
			priest.MindbenderAura.Activate(sim)
		},
	})
}
