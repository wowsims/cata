package death_knight

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (dk *DeathKnight) registerGlyphs() {
	// Major glyphs
	if dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfAntiMagicShell) {
		// Causes your Anti-Magic Shell to absorb all incoming magical damage, up to the absorption limit.
		// Handled in anti_magic_shell.go
	}
	if dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfDancingRuneWeapon) {
		dk.registerGlyphOfDancingRuneWeapon()
	}
	if dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfDeathCoil) {
		// Your Death Coil spell is now usable on all allies.
		// When cast on a non-undead ally, Death Coil shrouds them with a protective barrier that absorbs up to 168 damage.
		// TODO: Handle in death_coil.go
	}
	if dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfFesteringBlood) {
		// Blood Boil will now treat all targets as though they have Blood Plague or Frost Fever applied.
		// Handled in blood_boil.go
	}
	if dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfIceboundFortitude) {
		dk.registerGlyphOfIceboundFortitude()
	}
	if dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfLoudHorn) {
		dk.registerGlyphOfTheLoudHorn()
	}
	if dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfOutbreak) {
		// Your Outbreak spell no longer has a cooldown, but now costs 30 Runic Power.
		// Handled in outbreak.go
	}
	if dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfPestilence) {
		// Increases the radius of your Pestilence effect by 5 yards.
		// Handled in pestilence.go
	}
	if dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfRegenerativeMagic) {
		// If Anti-Magic Shell expires after its full duration, the cooldown is reduced by up to 50%, based on the amount of damage absorbtion remaining
		// TODO Handle in anti_magic_shell.go
	}
	if dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfShiftingPresences) {
		// You retain 70% of your Runic Power when switching Presences.
		// Handled in presences.go
	}
	if dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfUnholyFrenzy) {
		// Causes your Unholy Frenzy to no longer deal damage to the affected target.
		// TODO Handle in sim/core/buffs.go
	}
	if dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfVampiricBlood) {
		// Increases the bonus healing received while your Vampiric Blood is active by an additional 15%, but your Vampiric Blood no longer grants you health.
		// Handled in blood/vampiric_blood.go
	}

	// Minor glyphs
	if dk.HasMinorGlyph(proto.DeathKnightMinorGlyph_GlyphOfDeathsEmbrace) {
		// Your Death Coil refunds 20 Runic Power when used to heal an allied minion, but will no longer trigger Blood Tap when used this way.
		// Handled in death_coil.go
	}
	if dk.HasMinorGlyph(proto.DeathKnightMinorGlyph_GlyphOfTheLongWinter) {
		// The effect of your Horn of Winter now lasts for 1 hour.
		// Handled in horn_of_winter.go
	}
}

// Increases your threat generation by 100% while your Dancing Rune Weapon is active, but reduces its damage dealt by 25%.
func (dk *DeathKnight) registerGlyphOfDancingRuneWeapon() {
	if dk.Spec != proto.Spec_SpecBloodDeathKnight {
		return
	}

	dk.OnSpellRegistered(func(spell *core.Spell) {
		if !spell.Matches(DeathKnightSpellDancingRuneWeapon) {
			return
		}

		glyphAura := dk.RegisterAura(core.Aura{
			Label:    "Glyph of Dancing Rune Weapon" + dk.Label,
			ActionID: core.ActionID{SpellID: 63330},
			Duration: core.NeverExpires,
		}).AttachMultiplicativePseudoStatBuff(
			&dk.PseudoStats.ThreatMultiplier, 2.0,
		).AttachMultiplicativePseudoStatBuff(
			&dk.RuneWeapon.PseudoStats.DamageDealtMultiplier, 0.75,
		)

		spell.RelatedSelfBuff.AttachDependentAura(glyphAura)
	})
}

// Reduces the cooldown of your Icebound Fortitude by 50%, but also reduces its duration by 75%.
func (dk *DeathKnight) registerGlyphOfIceboundFortitude() {
	core.MakePermanent(dk.RegisterAura(core.Aura{
		Label:    "Glyph of Icebound Fortitude" + dk.Label,
		ActionID: core.ActionID{SpellID: 58673},
	})).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_Cooldown_Multiplier,
		ClassMask:  DeathKnightSpellIceboundFortitude,
		FloatValue: 0.5,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_BuffDuration_Flat,
		ClassMask: DeathKnightSpellIceboundFortitude,
		TimeValue: -core.DurationFromSeconds(12 * 0.75),
	})
}

// Your Horn of Winter now generates an additional 10 Runic Power, but the cooldown is increased by 100%.
func (dk *DeathKnight) registerGlyphOfTheLoudHorn() {
	rpMetrics := dk.NewRunicPowerMetrics(core.ActionID{SpellID: 147078})

	core.MakePermanent(dk.RegisterAura(core.Aura{
		Label:    "Glyph of the Loud Horn" + dk.Label,
		ActionID: core.ActionID{SpellID: 146646},
	})).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_Cooldown_Multiplier,
		ClassMask:  DeathKnightSpellHornOfWinter,
		FloatValue: 2.0,
	}).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: DeathKnightSpellHornOfWinter,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			dk.AddRunicPower(sim, 10, rpMetrics)
		},
	})
}
