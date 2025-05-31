package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

/*
Hammer the current target for 20% weapon damage, causing a wave of light that hits all targets within 8 yards for 35% Holy weapon damage and applying the Weakened Blows effect.
Grants a charge of Holy Power.

Weakened Blows
Demoralizes the target, reducing their physical damage dealt by 10% for 30 sec.
*/
func (paladin *Paladin) registerHammerOfTheRighteous() {
	numTargets := paladin.Env.GetNumTargets()
	actionID := core.ActionID{SpellID: 53595}
	paladin.CanTriggerHolyAvengerHpGain(actionID)
	auraArray := paladin.NewEnemyAuraArray(core.WeakenedBlowsAura)
	hasGlyphOfHammerOfTheRighteous := paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfHammerOfTheRighteous)

	hammerOfTheRighteousAoe := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 88263},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell | core.SpellFlagAoE,
		ClassSpellMask: SpellMaskHammerOfTheRighteousAoe,

		MaxRange: 8,

		DamageMultiplier: 0.35,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results := make([]*core.SpellResult, numTargets)

			for idx := range numTargets {
				currentTarget := sim.Environment.GetTargetUnit(idx)
				baseDamage := paladin.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
				results[idx] = spell.CalcDamage(sim, currentTarget, baseDamage, spell.OutcomeMagicCrit)
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				for idx := range numTargets {
					spell.DealDamage(sim, results[idx])
					aura := auraArray.Get(results[idx].Target)
					if hasGlyphOfHammerOfTheRighteous && aura.Duration != core.NeverExpires {
						aura.Duration = core.DurationFromSeconds(core.WeakenedBlowsDuration.Seconds() * 1.5)
					}
					aura.Activate(sim)
				}
			})
		},
	})

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskHammerOfTheRighteousMelee,

		MaxRange: core.MaxMeleeRange,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 3,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: paladin.Spec == proto.Spec_SpecHolyPaladin,
			CD: core.Cooldown{
				Timer:    paladin.BuilderCooldown(),
				Duration: time.Millisecond * 4500,
			},
		},

		DamageMultiplier: 0.2,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				paladin.HolyPower.Gain(sim, 1, actionID)
				hammerOfTheRighteousAoe.Cast(sim, target)
			}

			spell.DealDamage(sim, result)
		},
	})
}
