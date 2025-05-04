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
	}

	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfChainLightning) {

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

	//TODO verify in game
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfTelluricCurrents) {
		core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
			Name:           "Glyph of Telluric Currents",
			ClassSpellMask: SpellMaskLightningBolt | SpellMaskLightningBoltOverload,
			ProcChance:     1,
			Callback:       core.CallbackOnSpellHitDealt,
			Outcome:        core.OutcomeHit,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				amount := core.TernaryFloat64(shaman.Spec == proto.Spec_SpecElementalShaman, 0.02, 0.1)
				shaman.AddMana(sim, amount*shaman.MaxMana(), shaman.NewManaMetrics(core.ActionID{SpellID: 55453}))
			},
		})
	}
}
