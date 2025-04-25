package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (hunter *Hunter) ApplyGlyphs() {
	// Prime Glyphs
	if hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfArcaneShot) {
		hunter.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  HunterSpellArcaneShot,
			FloatValue: 0.12,
		})
	}
	if hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfExplosiveShot) {
		hunter.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusCrit_Percent,
			ClassMask:  HunterSpellExplosiveShot,
			FloatValue: 6,
		})
	}
	if hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfSerpentSting) {
		hunter.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusCrit_Percent,
			ClassMask:  HunterSpellSerpentSting,
			FloatValue: 6,
		})
	}
	if hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfKillCommand) {
		hunter.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_PowerCost_Flat,
			ClassMask: HunterSpellKillCommand,
			IntValue:  -3,
		})
	}
	if hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfChimeraShot) {
		hunter.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_Cooldown_Flat,
			ClassMask: HunterSpellChimeraShot,
			TimeValue: -time.Second * 1,
		})
	}
	if hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfAimedShot) {
		focusMetrics := hunter.NewFocusMetrics(core.ActionID{SpellID: 42897})
		core.MakePermanent(hunter.RegisterAura(core.Aura{
			Label: "Glyph of Aimed Shot",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == hunter.AimedShot && result.DidCrit() {
					hunter.AddFocus(sim, 5, focusMetrics)
				}
			},
		}))
	}
	if hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfKillShot) {
		icd := core.Cooldown{
			Timer:    hunter.NewTimer(),
			Duration: time.Second * 6,
		}
		core.MakePermanent(hunter.RegisterAura(core.Aura{
			Label: "Kill Shot Glyph",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == hunter.KillShot {
					if icd.IsReady(sim) {
						icd.Use(sim)
						hunter.KillShot.CD.Reset()
					}
				}
			},
		}))
	}
	// Major Glyphs
	if hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfBestialWrath) && hunter.Talents.BestialWrath {
		hunter.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_Cooldown_Flat,
			ClassMask: HunterSpellBestialWrath,
			TimeValue: -time.Second * 20,
		})
	}
}
