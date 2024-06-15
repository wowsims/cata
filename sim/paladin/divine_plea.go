package paladin

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (paladin *Paladin) registerDivinePleaSpell() {
	actionID := core.ActionID{SpellID: 54428}
	manaMetrics := paladin.NewManaMetrics(actionID)

	manaReturn := 0.12
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfDivinePlea) {
		manaReturn += 0.06
	}
	manaReturn /= 3

	paladin.DivinePleaAura = paladin.RegisterAura(core.Aura{
		Label:    "Divine Plea",
		ActionID: actionID,
		Duration: 9 * time.Second,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   3 * time.Second,
				NumTicks: 3,

				OnAction: func(sim *core.Simulation) {
					paladin.AddMana(sim, math.Round(manaReturn*paladin.MaxMana()), manaMetrics)
				},
			})
		},
	})

	paladin.DivinePlea = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskDivinePlea,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: 2 * time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			paladin.DivinePleaAura.Activate(sim)
		},
	})
}
