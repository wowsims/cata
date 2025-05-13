package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (rogue *Rogue) registerGougeSpell() {
	hasGlyph := rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfGouge)
	baseDamage := rogue.ClassSpellScaling * 0.10400000215

	rogue.Gouge = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1776},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: RogueSpellGouge,

		EnergyCost: core.EnergyCostOptions{
			Cost:   45 - 15*rogue.Talents.ImprovedGouge,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * 10,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.PseudoStats.InFrontOfTarget || hasGlyph
		},

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           rogue.CritMultiplier(false),
		ThreatMultiplier:         1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			calcBaseDamage := baseDamage +
				0.20999999344*spell.MeleeAttackPower()

			// Gouge has a Suppress Weapon Procs flag, but respects bonuses to MH such as expertise.
			// This models that effect without introducing a new spell flag/proc mask for this specific case
			spell.ProcMask = core.ProcMaskEmpty
			result := spell.CalcAndDealDamage(sim, target, calcBaseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			// Gouge disables auto attacks, requiring a macro to re-enable, retaining whatever the remaining swing timer is.
			// By pausing auto attacks for a short time, we can safely model Gouge usage without potentially over-valuing it.
			rogue.AutoAttacks.PauseMeleeBy(sim, rogue.ReactionTime)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
