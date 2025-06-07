package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (dk *DeathKnight) registerVampiricBloodSpell() {
	if !dk.Talents.VampiricBlood {
		return
	}

	actionID := core.ActionID{SpellID: 55233}
	healthMetrics := dk.NewHealthMetrics(actionID)

	// Implemented here for ease should maybe be moved to glyphs
	hasGlyph := dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfVampiricBlood)
	healBonus := core.TernaryFloat64(hasGlyph, 1.40, 1.25)

	var bonusHealth float64
	aura := dk.RegisterAura(core.Aura{
		Label:    "Vampiric Blood",
		ActionID: actionID,
		Duration: time.Second * 10,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.PseudoStats.HealingTakenMultiplier *= healBonus
			if !hasGlyph {
				bonusHealth = dk.MaxHealth() * 0.15
				dk.UpdateMaxHealth(sim, bonusHealth, healthMetrics)
			}

		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.PseudoStats.HealingTakenMultiplier /= healBonus
			if !hasGlyph {
				dk.UpdateMaxHealth(sim, -bonusHealth, healthMetrics)
			}
		},
	})

	spell := dk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: DeathKnightSpellVampiricBlood,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			aura.Activate(sim)
		},
	})

	dk.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
		ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
			return dk.CurrentHealthPercent() < 0.4
		},
	})
}
