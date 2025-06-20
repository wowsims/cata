package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) registerIcyVeinsCD() {

	actionID := core.ActionID{SpellID: 108978}
	mage.IcyVeinsAura = mage.RegisterAura(core.Aura{
		Label:    "Alter Time",
		ActionID: actionID,
		Duration: time.Second * 20,
	})

	mage.IcyVeins = mage.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: MageSpellIcyVeins,
		Flags:          core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 3,
			},
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			mage.IcyVeinsAura.Activate(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.IcyVeins,
		Type:  core.CooldownTypeDPS,
	})
}
