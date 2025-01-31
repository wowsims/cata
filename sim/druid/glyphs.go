package druid

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (druid *Druid) ApplyGlyphs() {

	if druid.HasPrimeGlyph(proto.DruidPrimeGlyph_GlyphOfMoonfire) {
		druid.AddStaticMod(core.SpellModConfig{
			ClassMask: DruidSpellMoonfireDoT | DruidSpellSunfireDoT,
			Kind:      core.SpellMod_DamageDone_Flat,
			IntValue:  20,
		})
	}

	if druid.HasPrimeGlyph(proto.DruidPrimeGlyph_GlyphOfInsectSwarm) {
		druid.AddStaticMod(core.SpellModConfig{
			ClassMask: DruidSpellInsectSwarm,
			Kind:      core.SpellMod_DamageDone_Flat,
			IntValue:  30,
		})
	}

	if druid.HasPrimeGlyph(proto.DruidPrimeGlyph_GlyphOfWrath) {
		druid.AddStaticMod(core.SpellModConfig{
			ClassMask: DruidSpellWrath,
			Kind:      core.SpellMod_DamageDone_Flat,
			IntValue:  10,
		})
	}

	if druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfStarfall) {
		druid.AddStaticMod(core.SpellModConfig{
			ClassMask: DruidSpellStarfall,
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: time.Second * -30,
		})
	}

	if druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfFocus) {
		druid.AddStaticMod(core.SpellModConfig{
			ClassMask: DruidSpellStarfall,
			Kind:      core.SpellMod_DamageDone_Flat,
			IntValue:  10,
		})

		// range mod?
	}

	if druid.HasPrimeGlyph(proto.DruidPrimeGlyph_GlyphOfStarsurge) {
		druid.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 62971},
			Label:    "Glyph of Starsurge",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},

			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.ClassSpellMask == DruidSpellStarsurge && !druid.Starfall.CD.IsReady(sim) {
					druid.Starfall.CD.Reduce(time.Second * 5)
				}
			},
		})
	}
}
