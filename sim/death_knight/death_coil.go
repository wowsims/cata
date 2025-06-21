package death_knight

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

var DeathCoilActionID = core.ActionID{SpellID: 47541}

/*
Fire a blast of unholy energy, causing (929 + <AP> * 0.514) Shadow damage to an enemy target or healing ((929 + <AP> * 0.514) * 3.5) damage on a friendly Undead target.

-- Glyph of Death's Embrace --

# Refunds 20 Runic Power when used to heal

-- /Glyph of Death's Embrace --
*/
func (dk *DeathKnight) registerDeathCoil() {
	dk.registerDeathCoilHeal()

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       DeathCoilActionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagEncounterOnly,
		ClassSpellMask: DeathKnightSpellDeathCoil,

		MaxRange: 30,

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: 40,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.CalcScalingSpellDmg(0.74544) + spell.MeleeAttackPower()*0.514
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (dk *DeathKnight) registerDeathCoilHeal() {
	if dk.Spec != proto.Spec_SpecUnholyDeathKnight {
		return
	}

	rpMetrics := dk.NewRunicPowerMetrics(core.ActionID{SpellID: 58679})
	hasGlyphOfDeathsEmbrace := dk.HasMinorGlyph(proto.DeathKnightMinorGlyph_GlyphOfDeathsEmbrace)

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       DeathCoilActionID.WithTag(2),
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellHealing,
		ClassSpellMask: DeathKnightSpellDeathCoilHeal,
		Flags:          core.SpellFlagAPL | core.SpellFlagPrepullOnly | core.SpellFlagNoMetrics | core.SpellFlagHelpful,

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: 40,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
		},

		DamageMultiplier: 3.5,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseHealing := dk.CalcScalingSpellDmg(0.74544) + spell.MeleeAttackPower()*0.514
			spell.CalcAndDealHealing(sim, &dk.Ghoul.Unit, baseHealing, spell.OutcomeHealingCrit)

			if hasGlyphOfDeathsEmbrace {
				dk.AddRunicPower(sim, 20, rpMetrics)
			}
		},
	})
}

func (dk *DeathKnight) registerDrwDeathCoil() *core.Spell {
	return dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    DeathCoilActionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.CalcScalingSpellDmg(0.74544) + spell.MeleeAttackPower()*0.514
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
