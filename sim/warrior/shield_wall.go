package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (war *Warrior) registerShieldWall() {
	hasGlyph := war.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfShieldWall)
	damageReductionMulti := 1 - core.TernaryFloat64(hasGlyph, 0.6, 0.4)
	cooldownDuration := core.TernaryDuration(hasGlyph, time.Minute*5, time.Minute*3)

	actionID := core.ActionID{SpellID: 871}
	aura := war.RegisterAura(core.Aura{
		Label:    "Shield Wall",
		ActionID: actionID,
		Duration: time.Second * 12,
	}).AttachMultiplicativePseudoStatBuff(
		&war.PseudoStats.DamageTakenMultiplier, damageReductionMulti,
	)

	spell := war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: SpellMaskShieldWall,
		Flags:          core.SpellFlagReadinessTrinket,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: cooldownDuration,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			aura.Activate(sim)
		},
		RelatedSelfBuff: aura,
	})

	war.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
		ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
			return war.CurrentHealthPercent() < 0.4
		},
	})
}
