package subtlety

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/rogue"
)

func (subRogue *SubtletyRogue) registerShadowDanceCD() {
	if !subRogue.Talents.ShadowDance {
		return
	}

	hasGlyph := subRogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfShadowDance)
	actionID := core.ActionID{SpellID: 51713}

	subRogue.ShadowDanceAura = subRogue.RegisterAura(core.Aura{
		Label:    "Shadow Dance",
		ActionID: actionID,
		Duration: core.TernaryDuration(hasGlyph, time.Second*8, time.Second*6),
		// Can now cast opening abilities outside of stealth
		// Covered in rogue.go by IsStealthed()
	})

	subRogue.ShadowDance = subRogue.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: rogue.RogueSpellShadowDance,

		Cast: core.CastConfig{
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    subRogue.NewTimer(),
				Duration: time.Minute,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			subRogue.BreakStealth(sim)
			subRogue.ShadowDanceAura.Activate(sim)
		},
		RelatedSelfBuff: subRogue.ShadowDanceAura,
	})

	subRogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    subRogue.ShadowDance,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault,
		ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
			return subRogue.GCD.IsReady(s) && subRogue.CurrentEnergy() >= 80 && subRogue.SliceAndDiceAura.IsActive() && subRogue.RecuperateAura.IsActive()
		},
	})
}
