package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

var OutbreakActionID = core.ActionID{SpellID: 77575}

// Instantly applies Blood Plague and Frost Fever to the target enemy.
func (dk *DeathKnight) registerOutbreak() {
	hasGlyphOfOutbreak := dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfOutbreak)
	cost := core.RuneCostOptions{}
	cast := core.CastConfig{
		DefaultCast: core.Cast{
			GCD: core.GCDMin,
		},
	}

	if hasGlyphOfOutbreak {
		cost.RunicPowerCost = 30
	} else {
		cast.CD = core.Cooldown{
			Timer:    dk.NewTimer(),
			Duration: time.Minute,
		}
	}

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       OutbreakActionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagReadinessTrinket,
		ClassSpellMask: DeathKnightSpellOutbreak,

		MaxRange: 30,

		RuneCost: cost,

		Cast: cast,

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
