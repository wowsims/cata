package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ProtectionWarrior) ApplyGlyphs() {
	war.Warrior.ApplyGlyphs()

	war.applyGlyphOfDevastate()
	war.applyGlyphOfShieldSlam()
	war.applyGlyphofShockwave()
}

func (war *ProtectionWarrior) applyGlyphOfDevastate() {
	if !war.HasPrimeGlyph(proto.WarriorPrimeGlyph_GlyphOfDevastate) {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskDevastate,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 5,
	})
}

func (war *ProtectionWarrior) applyGlyphOfShieldSlam() {
	if !war.HasPrimeGlyph(proto.WarriorPrimeGlyph_GlyphOfShieldSlam) {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskShieldSlam,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.1,
	})
}

func (war *ProtectionWarrior) applyGlyphofShockwave() {
	if !war.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfShockwave) {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask: warrior.SpellMaskShockwave,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -3 * time.Second,
	})
}
