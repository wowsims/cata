package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (dk *DeathKnight) registerRaiseDeadSpell() {
	// If talented as permanent pet skip this spell
	if dk.Inputs.Spec == proto.Spec_SpecUnholyDeathKnight {
		return
	}

	raiseDeadAura := dk.RegisterAura(core.Aura{
		Label:    "Raise Dead",
		ActionID: core.ActionID{SpellID: 46584},
		Duration: time.Minute * 1,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.Ghoul.Enable(sim, dk.Ghoul)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.Ghoul.Pet.Disable(sim)
		},
	})

	spell := dk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 46584},
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellRaiseDead,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			raiseDeadAura.Activate(sim)
		},
	})

	dk.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
			return dk.HasActiveAuraWithTag(core.UnholyFrenzyAuraTag)
		},
	})
}
