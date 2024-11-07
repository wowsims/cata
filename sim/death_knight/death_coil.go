package death_knight

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

var DeathCoilActionID = core.ActionID{SpellID: 47541}

func (dk *DeathKnight) registerDeathCoilSpell() {
	rpMetrics := dk.NewRunicPowerMetrics(DeathCoilActionID)
	hasGlyphOfDeathsEmbrace := dk.HasMinorGlyph(proto.DeathKnightMinorGlyph_GlyphOfDeathSEmbrace)

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       DeathCoilActionID,
		Flags:          core.SpellFlagAPL,
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
			var baseDamage float64
			if hasGlyphOfDeathsEmbrace && sim.CurrentTime < 0 {
				baseDamage = 0
			} else {
				baseDamage = dk.ClassSpellScaling*0.87599998713 + spell.MeleeAttackPower()*0.23
			}

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			// Instead of actually healing the ghoul, we just add the runic power
			// to not have to deal with healing metrics and other weird stuff.
			// Damage doesn't count before 0 anyway.
			if hasGlyphOfDeathsEmbrace && sim.CurrentTime < 0 {
				dk.AddRunicPower(sim, 20, rpMetrics)
			}
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
