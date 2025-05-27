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
	numTargets := core.TernaryInt32(hasMassExorcism, ret.Env.GetNumTargets(), 1)

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
			results := make([]*core.SpellResult, numTargets)

			currentTarget := target
			for idx := int32(0); idx < numTargets; idx++ {
				baseDamage := ret.CalcAndRollDamageRange(sim, scalingCoef, variance) +
					apCoef*spell.MeleeAttackPower()

				damageMultiplier := spell.DamageMultiplier
				if currentTarget != target {
					spell.DamageMultiplier *= 0.25
				}

				results[idx] = spell.CalcDamage(sim, currentTarget, baseDamage, spell.OutcomeMagicHitAndCrit)

				spell.DamageMultiplier = damageMultiplier

				currentTarget = sim.Environment.NextTargetUnit(currentTarget)
			}

			if results[0].Landed() {
				ret.HolyPower.Gain(1, exoHpActionID, sim)
			}

			for idx := int32(0); idx < numTargets; idx++ {
				spell.DealDamage(sim, results[idx])
			}
		},
	})
}
