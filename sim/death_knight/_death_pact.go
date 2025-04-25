package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (dk *DeathKnight) registerDeathPactSpell() {
	actionID := core.ActionID{SpellID: 48743}

	spell := dk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ProcMask:       core.ProcMaskSpellHealing,
		SpellSchool:    core.SpellSchoolShadow,
		ClassSpellMask: DeathKnightSpellDeathPact,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: 40,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute * 2,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return dk.Ghoul.Pet.IsEnabled()
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			healthGain := 0.25 * dk.Ghoul.MaxHealth()
			spell.CalcAndDealHealing(sim, spell.Unit, healthGain, spell.OutcomeHealing)
			dk.GetAura("Raise Dead").Deactivate(sim)
		},
	})

	if !dk.Inputs.IsDps {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeSurvival,
		})
	}
}
