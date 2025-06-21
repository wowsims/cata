package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (war *Warrior) registerRecklessness() {
	actionID := core.ActionID{SpellID: 1719}

	reckAura := war.RegisterAura(core.Aura{
		Label:    "Recklessness",
		ActionID: actionID,
		Duration: time.Second * 12,
	}).AttachSpellMod(core.SpellModConfig{
		ProcMask:   core.ProcMaskMeleeSpecial,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 30,
	})

	spell := war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL | core.SpellFlagReadinessTrinket,
		ClassSpellMask: SpellMaskRecklessness,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Minute * 3,
			},

			SharedCD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: 12 * time.Second,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			reckAura.Activate(sim)
		},

		RelatedSelfBuff: reckAura,
	})

	war.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}
