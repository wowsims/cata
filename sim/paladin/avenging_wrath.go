package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

/*
Increases all damage and healing caused by 20% for 20 sec.

-- Glyph of Avenging Wrath --
You heal for 1% of your maximum health every 2 sec.
-- /Glyph of Avenging Wrath --

-- Glyph of the Falling Avenger --
Your falling speed is slowed.
-- /Glyph of the Falling Avenger --
*/
func (paladin *Paladin) registerAvengingWrath() {
	actionID := core.ActionID{SpellID: 31884}

	paladin.AvengingWrathAura = paladin.RegisterAura(core.Aura{
		Label:    "Avenging Wrath" + paladin.Label,
		ActionID: actionID,
		Duration: core.DurationFromSeconds(core.TernaryFloat64(paladin.Talents.SanctifiedWrath, 30, 20)),
	}).AttachMultiplicativePseudoStatBuff(&paladin.Unit.PseudoStats.DamageDealtMultiplier, 1.2)
	core.RegisterPercentDamageModifierEffect(paladin.AvengingWrathAura, 1.2)

	avengingWrath := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: SpellMaskAvengingWrath,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: 3 * time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},

		RelatedSelfBuff: paladin.AvengingWrathAura,
	})

	paladin.AddMajorCooldown(core.MajorCooldown{
		Spell: avengingWrath,
		Type:  core.CooldownTypeDPS,
	})
}
