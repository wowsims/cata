package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

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
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
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
