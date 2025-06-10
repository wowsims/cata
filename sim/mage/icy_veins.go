package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (mage *Mage) registerIcyVeinsCD() {
	if mage.Spec != proto.Spec_SpecFrostMage {
		return
	}

	hasGlyph := mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfIcyVeins)

	icyVeinsMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll,
		FloatValue: -0.2,
		Kind:       core.SpellMod_CastTime_Pct,
	})

	actionID := core.ActionID{SpellID: 12472}
	mage.IcyVeinsAura = mage.RegisterAura(core.Aura{
		Label:    "Icy Veins",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if !hasGlyph {
				icyVeinsMod.Activate()
			}
		},
		OnExpire: func(_ *core.Aura, sim *core.Simulation) {
			if !hasGlyph {
				icyVeinsMod.Deactivate()
			}
		},
	})

	mage.IcyVeins = mage.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: MageSpellIcyVeins,
		Flags:          core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 3,
			},
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			mage.IcyVeinsAura.Activate(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.IcyVeins,
		Type:  core.CooldownTypeDPS,
	})
}
