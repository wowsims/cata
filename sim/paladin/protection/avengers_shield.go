package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (prot *ProtectionPaladin) registerAvengersShieldSpell() {
	actionId := core.ActionID{SpellID: 31935}
	asMinDamage, asMaxDamage := core.CalcScalingSpellEffectVarianceMinMax(proto.Class_ClassPaladin, 3.02399992943, 0.20000000298)
	glyphedSingleTargetAS := prot.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfFocusedShield)

	// Glyph to single target, OR apply to up to 3 targets
	numTargets := core.TernaryInt32(glyphedSingleTargetAS, 1, min(3, prot.Env.GetNumTargets()))
	results := make([]*core.SpellResult, numTargets)

	prot.AvengersShield = prot.RegisterSpell(core.SpellConfig{
		ActionID:    actionId,
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		MaxRange:     30,
		MissileSpeed: 35,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 6,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    prot.NewTimer(),
				Duration: time.Second * 15,
			},
		},

		DamageMultiplier: core.TernaryFloat64(glyphedSingleTargetAS, 1.3, 1),
		CritMultiplier:   prot.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			constBaseDamage := 0.20999999344*spell.SpellPower() + 0.41899999976*spell.MeleeAttackPower()

			for idx := int32(0); idx < numTargets; idx++ {
				baseDamage := constBaseDamage + sim.RollWithLabel(asMinDamage, asMaxDamage, "Avengers Shield"+prot.Label)

				currentTarget := sim.Environment.GetTargetUnit(idx)
				results[idx] = spell.CalcDamage(sim, currentTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			}

			for idx := int32(0); idx < numTargets; idx++ {
				spell.DealDamage(sim, results[idx])
			}

			if prot.GrandCrusaderAura.IsActive() {
				prot.HolyPower.Gain(1, actionId, sim)
			}

		},
	})
}
