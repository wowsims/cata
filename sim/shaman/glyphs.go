package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (shaman *Shaman) ApplyGlyphs() {
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfFireElementalTotem) {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask: SpellMaskFireElementalTotem,
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: time.Second * -150,
		})
		shaman.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellID == 2894 {
				shaman.AddStaticMod(core.SpellModConfig{
					Kind:      core.SpellMod_BuffDuration_Flat,
					TimeValue: -time.Second * 30,
					ClassMask: SpellMaskFireElementalTotem,
				})
			}
		})

	}

	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfChainLightning) {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask:  SpellMaskChainLightning | SpellMaskChainLightningOverload,
			Kind:       core.SpellMod_DamageDone_Pct,
			FloatValue: -0.1,
		})
	}

	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfFrostShock) {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask: SpellMaskFrostShock,
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: time.Second * -2,
		})
	}

	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfSpiritwalkersGrace) {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask: SpellMaskSpiritwalkersGrace,
			Kind:      core.SpellMod_BuffDuration_Flat,
			TimeValue: time.Second * 5,
		})
	}

	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfTelluricCurrents) {
		metric := shaman.NewManaMetrics(core.ActionID{SpellID: 55453})
		core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
			Name:           "Glyph of Telluric Currents",
			ClassSpellMask: SpellMaskLightningBolt | SpellMaskLightningBoltOverload,
			ProcChance:     1,
			Callback:       core.CallbackOnSpellHitDealt,
			Outcome:        core.OutcomeLanded,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				amount := core.TernaryFloat64(shaman.Spec == proto.Spec_SpecElementalShaman, 0.02, 0.1)
				shaman.AddMana(sim, amount*shaman.MaxMana(), metric)
			},
		})
	}
}
