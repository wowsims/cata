package protection

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/warrior"
)

func (war *ProtectionWarrior) ApplyGlyphs() {
	war.Warrior.ApplyGlyphs()

	war.applyGlyphOfDevastate()
	war.applyGlyphOfShieldSlam()
	war.applyGlyphofShockwave()
}

func (war *ProtectionWarrior) applyGlyphOfDevastate() {
	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskDevastate,
		Kind:       core.SpellMod_BonusCrit_Rating,
		FloatValue: 5 * core.CritRatingPerCritChance,
	})
}

func (war *ProtectionWarrior) applyGlyphOfShieldSlam() {
	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskShieldSlam,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.1,
	})
}

func (war *ProtectionWarrior) applyGlyphofShockwave() {
	war.AddStaticMod(core.SpellModConfig{
		ClassMask: warrior.SpellMaskShockwave,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -3 * time.Second,
	})
}
