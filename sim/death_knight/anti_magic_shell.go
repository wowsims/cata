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

	// runicPowerMetrics := dk.NewRunicPowerMetrics(core.ActionID{SpellID: 49088})
	currentShield := 0.0
	damageReductionMultiplier := 0.75

	if dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfAntiMagicShell) {
		damageReductionMultiplier += 0.25
	}

	var antiMagicShellAura *core.DamageAbsorptionAura
	antiMagicShellAura = dk.NewDamageAbsorptionAura(core.AbsorptionAuraConfig{
		Aura: core.Aura{
			Label:    "Anti-Magic Shell" + dk.Label,
			ActionID: actionID,
			Duration: time.Second * 5,
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if antiMagicShellAura.ShieldStrength > 0 {
					// TODO: Reduce CD
				}
			},
		},

		DamageMultiplier: damageReductionMultiplier,

		OnDamageAbsorbed: func(sim *core.Simulation, aura *core.DamageAbsorptionAura, result *core.SpellResult, absorbedDamage float64) {
			// TODO: RP return
		},
		ShouldApplyToResult: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult, isPeriodic bool) bool {
			return !spell.SpellSchool.Matches(core.SpellSchoolPhysical)
		},
		ShieldStrengthCalculator: func(unit *core.Unit) float64 {
			return currentShield
		},
	})

	dk.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		ProcMask:    core.ProcMaskSpellHealing,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagAPL | core.SpellFlagHelpful | core.SpellFlagReadinessTrinket,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Second * 45,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			currentShield = dk.MaxHealth() * 0.5
			antiMagicShellAura.Activate(sim)
		},
	})
}
