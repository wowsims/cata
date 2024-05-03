package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warrior *Warrior) RegisterCharge() {
	metrics := warrior.NewRageMetrics(core.ActionID{SpellID: 100})

	warrior.Charge = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 100},
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskCharge,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 15,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rage := float64(15 + 5*warrior.Talents.Blitz)
			warrior.AddRage(sim, rage, metrics)

		},
	})
}
