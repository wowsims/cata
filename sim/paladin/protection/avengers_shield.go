package protection

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (prot *ProtectionPaladin) registerAvengersShieldSpell() {
	actionId := core.ActionID{SpellID: 31935}
	hpMetrics := prot.NewHolyPowerMetrics(actionId)
	asMinDamage, asMaxDamage := core.CalcScalingSpellEffectVarianceMinMax(proto.Class_ClassPaladin, 3.02399992943, 0.20000000298)
	glyphedSingleTargetAS := prot.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfFocusedShield)

	// Glyph to single target, OR apply to up to 3 targets
	maxTargets := core.TernaryInt32(glyphedSingleTargetAS, 1, min(3, prot.Env.TotalTargetCount()))
	results := make([]*core.SpellResult, maxTargets)

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
		CritMultiplier:   prot.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			constBaseDamage := 0.20999999344*spell.SpellPower() + 0.41899999976*spell.MeleeAttackPower()
			numTargets := min(maxTargets, sim.Environment.ActiveTargetCount())

			for idx := int32(0); idx < numTargets; idx++ {
				baseDamage := constBaseDamage + sim.RollWithLabel(asMinDamage, asMaxDamage, "Avengers Shield"+prot.Label)
				results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				target = sim.Environment.NextActiveTargetUnit(target)
			}

			for idx := int32(0); idx < numTargets; idx++ {
				spell.DealDamage(sim, results[idx])
			}

			if prot.GrandCrusaderAura.IsActive() {
				prot.GainHolyPower(sim, 1, hpMetrics)
			}

		},
	})
}
