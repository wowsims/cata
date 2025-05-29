package paladin

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

/*
Inspire

-- Glyph of Devotion Aura --
yourself, granting you
-- else --
all party and raid members within 40 yards, granting them
----------

immunity to Silence and Interrupt effects and reducing all magic damage taken by 20%.
Lasts 6 sec.
*/
func (paladin *Paladin) registerDevotionAura() {
	devotionAura := core.DevotionAuraAura(&paladin.Character, 0)

	devotionAuraSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.DevotionAuraActionID,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: SpellMaskDevotionAura,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: core.DevotionAuraCD,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			devotionAura.Activate(sim)
		},
	})

	if paladin.Spec == proto.Spec_SpecProtectionPaladin {
		paladin.AddMajorCooldown(core.MajorCooldown{
			Spell: devotionAuraSpell,
			Type:  core.CooldownTypeSurvival,
		})
	}
}
