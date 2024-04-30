package death_knight

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (dk *DeathKnight) registerBloodTapSpell() {
	actionID := core.ActionID{SpellID: 45529}

	bloodMetrics := dk.NewBloodRuneMetrics(actionID)
	deathMetrics := dk.NewDeathRuneMetrics(actionID)
	aura := dk.RegisterAura(core.Aura{
		Label:    "Blood Tap",
		ActionID: actionID,
		Duration: time.Second * 20,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.BloodTapConversion(sim, bloodMetrics, deathMetrics)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.CancelBloodTap(sim)
		},
	})

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellBloodTap,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			aura.Activate(sim)
		},
	})
}
