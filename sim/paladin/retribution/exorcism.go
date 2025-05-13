package retribution

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/paladin"
)

func (ret *RetributionPaladin) registerExorcism() {
	apCoef := 0.67699998617
	scalingCoef := 6.09499979019
	variance := 0.1099999994

	exoHpActionID := core.ActionID{SpellID: 147715}
	ret.CanTriggerHolyAvengerHpGain(exoHpActionID)

	hasMassExorcism := ret.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfMassExorcism)

	ret.Exorcism = ret.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 879},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: paladin.SpellMaskExorcism,

		MaxRange: core.TernaryFloat64(hasMassExorcism, core.MaxMeleeRange, 30),

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
			baseDamage := ret.CalcAndRollDamageRange(sim, scalingCoef, variance) +
				apCoef*spell.MeleeAttackPower()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				ret.HolyPower.Gain(1, exoHpActionID, sim)
			}

			spell.DealOutcome(sim, result)
		},
	})
}
