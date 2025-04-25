package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (dk *DeathKnight) registerHornOfWinterSpell() {
	actionID := core.ActionID{SpellID: 57330}
	rpMetrics := dk.NewRunicPowerMetrics(actionID)

	// hornAura := dk.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
	// 	return core.HornOfWinterAura(unit, false, dk.HasMinorGlyph(proto.DeathKnightMinorGlyph_GlyphOfHornOfWinter))
	// })

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellHornOfWinter,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: 20 * time.Second,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dk.AddRunicPower(sim, 10, rpMetrics)

			// for _, unit := range sim.Raid.GetActiveAllyUnits() {
			// 	// Horn of Winter doesnt apply to pets
			// 	if unit.Type != core.PetUnit {
			// 		hornAura.Get(unit).Activate(sim)
			// 	}
			// }
		},
	})
}
