package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

var RuneTapActionID = core.ActionID{SpellID: 48982}

func (dk *DeathKnight) registerRuneTapSpell() {
	if !dk.Talents.RuneTap {
		return
	}

	spell := dk.RegisterSpell(core.SpellConfig{
		ActionID:       RuneTapActionID,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ProcMask:       core.ProcMaskSpellHealing,
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: DeathKnightSpellRuneTap,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost: 1,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Second * 30,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealHealing(sim, &dk.Unit, dk.MaxHealth()*0.1, spell.OutcomeHealing)
		},
	})

	dk.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
		ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
			return dk.CurrentHealthPercent() < 0.7
		},
	})
}
