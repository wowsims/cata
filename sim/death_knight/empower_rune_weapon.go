package death_knight

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (dk *DeathKnight) registerEmpowerRuneWeaponSpell() {
	actionID := core.ActionID{SpellID: 47568}
	cdTimer := dk.NewTimer()
	cd := time.Minute * 5

	rpMetrics := dk.NewRunicPowerMetrics(actionID)
	dk.EmpowerRuneWeapon = dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dk.RegenAllRunes(sim)
			dk.AddRunicPower(sim, 25, rpMetrics)
		},
	})
}
