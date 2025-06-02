package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (warrior *Warrior) registerShieldWall() {
	hasGlyph := warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfShieldWall)
	damageReductionMulti := 1 - core.TernaryFloat64(hasGlyph, 0.6, 0.4)
	cooldownDuration := core.TernaryDuration(hasGlyph, time.Minute*5, time.Minute*3)

	actionID := core.ActionID{SpellID: 871}
	aura := warrior.RegisterAura(core.Aura{
		Label:    "Shield Wall",
		ActionID: actionID,
		Duration: time.Second * 12,
	}).AttachMultiplicativePseudoStatBuff(
		&warrior.PseudoStats.DamageTakenMultiplier, damageReductionMulti,
	)

	spell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: SpellMaskShieldWall,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownDuration,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			aura.Activate(sim)
		},
		RelatedSelfBuff: aura,
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
	})
}
