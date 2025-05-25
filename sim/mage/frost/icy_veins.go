package frost

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/mage"
)

func (frostMage *FrostMage) registerIcyVeinsCD() {
	icyVeinsMod := frostMage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  mage.MageSpellsAll,
		FloatValue: -0.2,
		Kind:       core.SpellMod_CastTime_Pct,
	})

	actionID := core.ActionID{SpellID: 12472}
	frostMage.icyVeinsAura = frostMage.RegisterAura(core.Aura{
		Label:    "Icy Veins",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if !frostMage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfIcyVeins) {
				icyVeinsMod.Activate()
			}
		},
		OnExpire: func(_ *core.Aura, sim *core.Simulation) {
			if !frostMage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfIcyVeins) {
				icyVeinsMod.Deactivate()
			}
		},
	})

	frostMage.IcyVeins = frostMage.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: mage.MageSpellIcyVeins,
		Flags:          core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    frostMage.NewTimer(),
				Duration: time.Minute * 3,
			},
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			frostMage.icyVeinsAura.Activate(sim)
		},
	})

	frostMage.AddMajorCooldown(core.MajorCooldown{
		Spell: frostMage.IcyVeins,
		Type:  core.CooldownTypeDPS,
	})
}
