package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

var OutbreakActionID = core.ActionID{SpellID: 77575}

// Instantly applies Blood Plague and Frost Fever to the target enemy.
func (dk *DeathKnight) registerOutbreak() {
	dk.RegisterSpell(core.SpellConfig{
		ActionID:       OutbreakActionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellOutbreak,

		MaxRange: 30,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				dk.FrostFeverSpell.Cast(sim, target)
				dk.BloodPlagueSpell.Cast(sim, target)
			}
		},
	})
}

func (dk *DeathKnight) registerDrwOutbreak() *core.Spell {
	return dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:       OutbreakActionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellOutbreak,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				dk.RuneWeapon.FrostFeverSpell.Cast(sim, target)
				dk.RuneWeapon.BloodPlagueSpell.Cast(sim, target)
			}
		},
	})
}
