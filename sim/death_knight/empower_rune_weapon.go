package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

// Empower your rune weapon, immediately activating all your runes and generating 25 Runic Power.
func (dk *DeathKnight) registerEmpowerRuneWeapon() {
	actionId := core.ActionID{SpellID: 47568}
	metrics := []*core.ResourceMetrics{
		dk.NewBloodRuneMetrics(actionId),
		dk.NewFrostRuneMetrics(actionId),
		dk.NewUnholyRuneMetrics(actionId),
		dk.NewDeathRuneMetrics(actionId),
		dk.NewRunicPowerMetrics(actionId),
	}

	spell := dk.RegisterSpell(core.SpellConfig{
		ActionID: actionId,
		Flags:    core.SpellFlagAPL | core.SpellFlagHelpful | core.SpellFlagNoOnCastComplete | core.SpellFlagReadinessTrinket,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute * 5,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dk.AddRunicPower(sim, 25, metrics[4])
			dk.RegenAllRunes(sim, metrics)
		},
	})

	dk.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, _ *core.Character) bool {
			return dk.AllRunesSpent()
		},
	})
}
