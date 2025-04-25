package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (dk *DeathKnight) registerIceboundFortitudeSpell() {
	actionID := core.ActionID{SpellID: 48792}

	dmgTakenMult := 0.8 - 0.15*float64(dk.Talents.SanguineFortitude)

	iceBoundFortituteAura := dk.RegisterAura(core.Aura{
		Label:    "Icebound Fortitude",
		ActionID: actionID,
		Duration: 12 * time.Second,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier *= dmgTakenMult
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier /= dmgTakenMult
		},
	})

	spell := dk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellIceboundFortitude,

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: 20,
		},
		Cast: core.CastConfig{
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
