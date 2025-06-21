package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

/*
The Death Knight blows the Horn of Winter, which generates 10 Runic Power and increases attack power of all party and raid members within 100 yards by 10%.
Lasts 5 min.
*/
func (dk *DeathKnight) registerHornOfWinter() {
	actionID := core.ActionID{SpellID: 57330}
	rpMetrics := dk.NewRunicPowerMetrics(actionID)
	hasGlyphOfTheLongWinter := dk.HasMinorGlyph(proto.DeathKnightMinorGlyph_GlyphOfTheLongWinter)

	hornArray := dk.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		aura := core.HornOfWinterAura(unit, false)
		if hasGlyphOfTheLongWinter {
			aura.Duration = time.Hour * 1
		}
		return aura
	})

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: DeathKnightSpellHornOfWinter,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: 20 * time.Second,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, unit := range sim.Raid.AllPlayerUnits {
				hornArray.Get(unit).Activate(sim)
			}

			dk.AddRunicPower(sim, 10, rpMetrics)
		},
	})
}
