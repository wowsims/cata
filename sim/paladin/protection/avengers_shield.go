package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (prot *ProtectionPaladin) registerAvengersShieldSpell() {
	actionID := core.ActionID{SpellID: 31935}
	scalingCoef := 5.89499998093
	variance := 0.20000000298
	spCoef := 0.31499999762
	apCoef := 0.81749999523
	glyphedSingleTargetAS := prot.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfFocusedShield)

	// Glyph to single target, OR apply to up to 3 targets
	numTargets := core.TernaryInt32(glyphedSingleTargetAS, 1, min(3, prot.Env.GetNumTargets()))
	results := make([]*core.SpellResult, numTargets)

	prot.AvengersShield = prot.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		MaxRange:     30,
		MissileSpeed: 35,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 7,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    prot.NewTimer(),
				Duration: time.Second * 15,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return prot.OffHand().WeaponType == proto.WeaponType_WeaponTypeShield
		},

		DamageMultiplier: core.TernaryFloat64(glyphedSingleTargetAS, 1.3, 1),
		CritMultiplier:   prot.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			constBaseDamage := spCoef*spell.SpellPower() + apCoef*spell.MeleeAttackPower()

			for idx := int32(0); idx < numTargets; idx++ {
				baseDamage := constBaseDamage + prot.CalcAndRollDamageRange(sim, scalingCoef, variance)
				results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				target = sim.Environment.NextTargetUnit(target)
			}

			for idx := int32(0); idx < numTargets; idx++ {
				spell.DealDamage(sim, results[idx])
			}
		},
	})
}
