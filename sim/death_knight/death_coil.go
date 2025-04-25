package death_knight

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

var DeathCoilActionID = core.ActionID{SpellID: 47541}

func (dk *DeathKnight) registerDeathCoilSpell() {
	rpMetrics := dk.NewRunicPowerMetrics(core.ActionID{SpellID: 58679})
	hasGlyphOfDeathsEmbrace := dk.HasMinorGlyph(proto.DeathKnightMinorGlyph_GlyphOfDeathSEmbrace)

	// Death Coil Heal
	dk.RegisterSpell(core.SpellConfig{
		ActionID:       DeathCoilActionID.WithTag(2),
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellHealing,
		ClassSpellMask: DeathKnightSpellDeathCoilHeal,
		Flags:          core.SpellFlagAPL | core.SpellFlagPrepullOnly | core.SpellFlagNoMetrics,

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: 40,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 3.5,
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1.0,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return dk.Inputs.Spec == proto.Spec_SpecUnholyDeathKnight
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseHealing := dk.ClassSpellScaling*0.87599998713 + spell.MeleeAttackPower()*0.23
			spell.CalcAndDealHealing(sim, &dk.Ghoul.Unit, baseHealing, spell.OutcomeHealingCrit)

			if hasGlyphOfDeathsEmbrace {
				dk.AddRunicPower(sim, 20, rpMetrics)
			}
		},
	})

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       DeathCoilActionID,
		Flags:          core.SpellFlagAPL | core.SpellFlagEncounterOnly,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DeathKnightSpellDeathCoil,

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: 40,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.87599998713 + spell.MeleeAttackPower()*0.23
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (dk *DeathKnight) registerDrwDeathCoilSpell() *core.Spell {
	return dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    DeathCoilActionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.87599998713 + spell.MeleeAttackPower()*0.23
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
