package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

/*
Raises a ghoul to fight by your side.
You can have a maximum of one ghoul at a time.
Lasts 1 min.
*/
func (dk *DeathKnight) registerRaiseDead() {
	dk.RaiseDeadAura = dk.RegisterAura(core.Aura{
		Label:    "Raise Dead" + dk.Label,
		ActionID: core.ActionID{SpellID: 46584},
		Duration: time.Minute * 1,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.Ghoul.Enable(sim, dk.Ghoul)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.Ghoul.Pet.Disable(sim)
		},
	})

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 46584},
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellRaiseDead,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},

		RelatedSelfBuff: dk.RaiseDeadAura,
	})
}
