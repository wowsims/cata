package blood

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

// Converts 1 Blood Rune into 10% of your maximum health.
func (bdk *BloodDeathKnight) registerRuneTap() {
	bdk.RuneTapSpell = bdk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 48982},
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagAPL | core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,
		ProcMask:       core.ProcMaskSpellHealing,
		ClassSpellMask: death_knight.DeathKnightSpellRuneTap,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    bdk.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealHealing(sim, &bdk.Unit, bdk.MaxHealth()*0.1, spell.OutcomeHealing)
		},
	})

	bdk.AddMajorCooldown(core.MajorCooldown{
		Spell: bdk.RuneTapSpell,
		Type:  core.CooldownTypeSurvival,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return bdk.CurrentHealthPercent() < 0.7
		},
	})
}
