package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core/proto"

	"github.com/wowsims/mop/sim/core"
)

func (dk *DeathKnight) registerAntiMagicShellSpell() {
	actionID := core.ActionID{SpellID: 48707}

	dmgReductionModifier := 0.75 + []float64{0, 0.08, 0.16, 0.25}[dk.Talents.MagicSuppression]
	currentShield := 0.0

	var shieldSpell *core.Spell
	shieldSpell = dk.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		ProcMask:    core.ProcMaskSpellHealing,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Second * 45,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Shield: core.ShieldConfig{
			SelfOnly: true,
			Aura: core.Aura{
				Label:    "Anti-Magic Shell",
				ActionID: actionID,
				Duration: time.Second*5 + core.TernaryDuration(dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfAntiMagicShell), 2*time.Second, 0),

				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.SpellSchool.Matches(core.SpellSchoolPhysical) {
						return
					}

					if currentShield <= 0 || result.Damage <= 0 {
						return
					}

					damageReduced := min(result.Damage*dmgReductionModifier, currentShield)
					currentShield -= damageReduced

					dk.GainHealth(sim, damageReduced, shieldSpell.HealthMetrics(result.Target))

					if currentShield <= 0 {
						shieldSpell.SelfShield().Deactivate(sim)
					}
				},
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			currentShield = dk.MaxHealth() * 0.5
			spell.SelfShield().Apply(sim, currentShield)
		},
	})

	// TODO:
	if dk.Inputs.UseAMS {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell:    shieldSpell,
			Type:     core.CooldownTypeDPS,
			Priority: core.CooldownPriorityLow,
		})
	}
}
