package retribution

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/paladin"
)

/*
Forcefully attempt to expel the evil from the target with a blast of Holy Light.
Causes (<6577-7343> + 0.677 * <AP>) Holy damage

-- Glyph of Mass Exorcism --
to the target and 25% of that to other enemies within 8 yards
-- /Glyph of Mass Exorcism --

and generates a charge of Holy Power.
*/
func (ret *RetributionPaladin) registerExorcism() {
	exoHpActionID := core.ActionID{SpellID: 147715}
	ret.CanTriggerHolyAvengerHpGain(exoHpActionID)

	hasGlyphOfMassExorcism := ret.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfMassExorcism)

	ret.Exorcism = ret.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 879},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: paladin.SpellMaskExorcism,

		MaxRange: core.TernaryFloat64(hasGlyphOfMassExorcism, core.MaxMeleeRange, 30),

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 4,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    ret.NewTimer(),
				Duration: time.Second * 15,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   ret.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := ret.CalcAndRollDamageRange(sim, 6.09499979019, 0.1099999994) +
				0.67699998617*spell.MeleeAttackPower()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				ret.HolyPower.Gain(sim, 1, exoHpActionID)
			}

			spell.DealDamage(sim, result)
		},
	})
}
