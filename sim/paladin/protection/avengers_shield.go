package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

/*
Hurls your shield at an enemy target, dealing (<6058-7405> + 0.8175 * <AP> + 0.315 * <SP>) Holy damage,

-- Glyph of Dazing Shield --
dazing,
-- /Glyph of Dazing Shield --

silencing and interrupting spellcasting for 3 sec, and then jumping to additional nearby enemies.

Affects 3 total targets.
*/
func (prot *ProtectionPaladin) registerAvengersShieldSpell() {
	hasGlyphOfFocusedShield := prot.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfFocusedShield)

	// Glyph to single target, OR apply to up to 3 targets
	numTargets := core.TernaryInt32(hasGlyphOfFocusedShield, 1, min(3, prot.Env.GetNumTargets()))
	results := make([]*core.SpellResult, numTargets)

	prot.AvengersShield = prot.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 31935},
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
			return prot.PseudoStats.CanBlock
		},

		DamageMultiplier: 1,
		CritMultiplier:   prot.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			bonusDamage := 0.31499999762*spell.SpellPower() + 0.81749999523*spell.MeleeAttackPower()

			for idx := range numTargets {
				baseDamage := prot.CalcAndRollDamageRange(sim, 5.89499998093, 0.20000000298) + bonusDamage
				results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				target = sim.Environment.NextTargetUnit(target)
			}

			for idx := range numTargets {
				spell.DealDamage(sim, results[idx])
			}
		},
	})
}
