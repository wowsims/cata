package subtlety

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/rogue"
)

func (subRogue *SubtletyRogue) registerShadowDanceCD() {
	if !subRogue.Talents.ShadowDance {
		return
	}

	hasGlyph := subRogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfShadowDance)
	getDuration := func() time.Duration {
		return core.TernaryDuration(hasGlyph, time.Second*8, time.Second*6) + core.TernaryDuration(subRogue.Has4pcT13(), time.Second*2, 0)
	}
	actionID := core.ActionID{SpellID: 51713}

	subRogue.ShadowDanceAura = subRogue.RegisterAura(core.Aura{
		Label:    "Shadow Dance",
		ActionID: actionID,
		Duration: getDuration(),
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
			subRogue.ShadowDanceAura.Duration = getDuration()
			subRogue.ShadowDanceAura.Activate(sim)
		},
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
