package subtlety

import (
	"time"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/rogue"

	"github.com/wowsims/cata/sim/core"
)

func (subRogue *SubtletyRogue) registerShadowDanceCD() {
	if !subRogue.Talents.ShadowDance {
		return
	}

	hasGlyph := subRogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfShadowDance)
	t13Bonus := subRogue.HasSetBonus(rogue.Tier13, 4)
	duration := core.TernaryDuration(hasGlyph, time.Second*8, time.Second*6) + core.TernaryDuration(t13Bonus, time.Second*2, 0)

	actionID := core.ActionID{SpellID: 51713}

	subRogue.ShadowDanceAura = subRogue.RegisterAura(core.Aura{
		Label:    "Shadow Dance",
		ActionID: actionID,
		Duration: duration,
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
