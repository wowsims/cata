package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core/proto"

	"github.com/wowsims/mop/sim/core"
)

/*
Surrounds the Death Knight in an Anti-Magic Shell, absorbing 75% of damage dealt by harmful spells (up to a maximum of 50% of the Death Knight's health) and preventing application of harmful magical effects.
Damage absorbed generates Runic Power.
Lasts 5 sec.
*/
func (dk *DeathKnight) registerAntiMagicShell() {
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
}
