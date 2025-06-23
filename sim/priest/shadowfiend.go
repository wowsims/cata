package priest

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (priest *Priest) registerShadowfiendSpell() {
	if priest.Talents.Mindbender {
		return
	}

	actionID := core.ActionID{SpellID: 34433}

	// For timeline only
	priest.ShadowfiendAura = priest.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Shadowfiend",
		Duration: time.Second * 12.0,
	})

	priest.Shadowfiend = priest.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL | core.SpellFlagReadinessTrinket,
		ClassSpellMask: PriestSpellShadowFiend,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			priest.ShadowfiendPet.EnableWithTimeout(sim, priest.ShadowfiendPet, time.Second*12.0)
			priest.ShadowfiendAura.Activate(sim)
		},
	})
}
