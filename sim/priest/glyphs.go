package priest

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (priest *Priest) ApplyGlyphs() {
	// Glphys
	if priest.HasGlyph(int32(proto.PriestPrimeGlyph_GlyphOfShadowWordPain)) {
		priest.AddStaticMod(core.SpellModConfig{
			FloatValue: 0.1,
			ClassMask:  int64(PriestSpellShadowWordPain),
			Kind:       core.SpellMod_DamageDone_Flat,
		})
	}

	if priest.HasGlyph(int32(proto.PriestPrimeGlyph_GlyphOfMindFlay)) {
		priest.AddStaticMod(core.SpellModConfig{
			ClassMask:  int64(PriestSpellMindFlay),
			FloatValue: 0.1,
			Kind:       core.SpellMod_DamageDone_Flat,
		})
	}

	if priest.HasGlyph(int32(proto.PriestPrimeGlyph_GlyphOfDispersion)) {
		priest.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: time.Second * -45,
			ClassMask: int64(PriestSpellDispersion),
		})
	}

	if priest.HasGlyph(int32(proto.PriestPrimeGlyph_GlyphOfShadowWordDeath)) {
		priest.RegisterAura(core.Aura{
			Label:    "Glyph of Shadow Word: Death",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},

			Icd: &core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 6,
			},

			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if priest.ShadowWordDeath == spell && sim.IsExecutePhase25() && aura.Icd.IsReady(sim) {
					aura.Icd.Use(sim)
					priest.ShadowWordDeath.CD.Reset()
				}
			},
		})
	}
}
