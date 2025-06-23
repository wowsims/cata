package monk

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (monk *Monk) registerFortifyingBrew() {
	actionID := core.ActionID{SpellID: 126456}
	healthMetrics := monk.NewHealthMetrics(actionID)

	hasGlyphOfFortifyingBrew := monk.HasMajorGlyph(proto.MonkMajorGlyph_GlyphOfFortifyingBrew)
	healthModifier := core.TernaryFloat64(hasGlyphOfFortifyingBrew, 0.10, 0.20)
	damageTakenModifier := core.TernaryFloat64(hasGlyphOfFortifyingBrew, 0.75, 0.8)

	var bonusHealth float64
	monk.FortifyingBrewAura = monk.RegisterAura(core.Aura{
		Label:    "Fortifying Brew" + monk.Label,
		ActionID: actionID,
		Duration: time.Second * 20,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			bonusHealth = monk.MaxHealth() * healthModifier
			monk.PseudoStats.DamageTakenMultiplier *= damageTakenModifier
			monk.UpdateMaxHealth(sim, bonusHealth, healthMetrics)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			monk.PseudoStats.DamageTakenMultiplier /= damageTakenModifier
			monk.UpdateMaxHealth(sim, -bonusHealth, healthMetrics)
		},
	})

	spell := monk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL | core.SpellFlagReadinessTrinket,
		ClassSpellMask: MonkSpellFortifyingBrew,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    monk.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			monk.FortifyingBrewAura.Activate(sim)
		},
	})

	monk.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
		ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
			return monk.CurrentHealthPercent() < 0.4
		},
	})
}
