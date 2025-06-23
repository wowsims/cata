package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

// The Death Knight freezes his blood to become immune to Stun effects and reduce all damage taken by 20% for 12 sec.
func (dk *DeathKnight) registerIceboundFortitude() {
	actionID := core.ActionID{SpellID: 48792}

	dmgTakenMult := 0.8 - core.TernaryFloat64(dk.Spec == proto.Spec_SpecBloodDeathKnight, 0.3, 0)

	iceBoundFortituteAura := dk.RegisterAura(core.Aura{
		Label:    "Icebound Fortitude" + dk.Label,
		ActionID: actionID,
		Duration: 12 * time.Second,
	}).AttachMultiplicativePseudoStatBuff(&dk.PseudoStats.DamageTakenMultiplier, dmgTakenMult)

	spell := dk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL | core.SpellFlagReadinessTrinket,
		ClassSpellMask: DeathKnightSpellIceboundFortitude,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},

		RelatedSelfBuff: iceBoundFortituteAura,
	})

	if !dk.Inputs.IsDps {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeSurvival,
			ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
				return dk.CurrentHealthPercent() < 0.2
			},
		})
	}
}
